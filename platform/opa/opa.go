package opa

import (
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model"
	"2f-authorization/platform/logger"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	dbstore "2f-authorization/internal/storage"

	"github.com/goccy/go-json"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/util"
	"go.uber.org/zap"
)

type Opa interface {
	Refresh(ctx context.Context, reason string) error
	GetData(ctx context.Context) error
	Allow(ctx context.Context, req model.Request) (bool, error)
	AllowedPermissions(ctx context.Context, input map[string]interface{}) (interface{}, error)
}

type opa struct {
	db            dbstore.Policy
	store         storage.Store
	policy        string
	Query         string
	log           logger.Logger
	filepath      string
	regopath      string
	server        string
	query         rego.PreparedEvalQuery
	evaluatorPort int
}

func Init(policy string, policyDb dbstore.Policy, filepath, regopath, server string, port int, log logger.Logger) Opa {

	go func() {
		cmd := exec.Command(server, "run", "--server", "--watch", "--addr", fmt.Sprintf("localhost:%d", port), regopath, filepath)
		defer func() {
			err := cmd.Process.Kill()
			log.Error(context.Background(), "error while killing somebody", zap.Error(err))
		}()
		output, err := cmd.CombinedOutput()
		if err != nil {
			err := errors.ErrOpaPrepareEvalError.Wrap(err, "error  Initializing OPA  Server")
			log.Fatal(context.Background(), "error preparing the rego for eval", zap.Error(err), zap.String("combined-output", string(output)))
		}
	}()

	return &opa{
		policy:        policy,
		db:            policyDb,
		filepath:      filepath,
		regopath:      regopath,
		server:        server,
		log:           log,
		evaluatorPort: port,
	}

	go func() {

		err := exec.Command(server, "run", "--server", "--watch", regopath, filepath).Run()
		if err != nil {
			err := errors.ErrOpaPrepareEvalError.Wrap(err, "error  Initializing OPA  Server")
			log.Error(context.Background(), "error preparing the rego for eval", zap.Error(err))
		}
	}()

	return &opa{
		policy:   policy,
		db:       policyDb,
		filepath: filepath,
		regopath: regopath,
		server:   server,
		log:      log,
	}
}

type responseBody struct {
	Response bool `json:"result"`
}
type RequestBody struct {
	Input model.Request `json:"input"`
}

func (o *opa) Allow(ctx context.Context, req model.Request) (bool, error) {
	posturl := fmt.Sprintf("http://localhost:%d/v1/data/authz/allow", o.evaluatorPort)
	reqst := RequestBody{
		Input: req,
	}
	resp := responseBody{}
	js, err := json.Marshal(reqst)
	if err != nil {
		err := errors.ErrOpaPrepareEvalError.Wrap(err, "error while preparing evaluation to json")
		o.log.Error(ctx, "error preparing the opa data to json", zap.Error(err))
		return false, err
	}
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(js))

	if err != nil {
		err := errors.ErrOpaPrepareEvalError.Wrap(err, "error while preparing evaluation")
		o.log.Error(ctx, "error preparing the opa data", zap.Error(err))
		return false, err
	}
	httpCli := &http.Client{}
	res, err := httpCli.Do(r)
	if err != nil {
		err := errors.ErrOpaPrepareEvalError.Wrap(err, "error while getting response from opa server")
		o.log.Error(ctx, "error while getting response from opa server", zap.Error(err))
		return false, err
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&resp)

	return resp.Response, nil
}

func (o *opa) Refresh(ctx context.Context, reason string) error {
	o.log.Info(ctx, reason)
	if err := o.GetData(ctx); err != nil {
		return err
	}

	return nil
}

func (o *opa) GetData(ctx context.Context) error {
	data, err := o.db.GetOpaData(ctx)
	if err != nil {
		return err
	}
	var services map[string]interface{}
	err = util.UnmarshalJSON(data, &services)
	if err != nil {
		err := errors.ErrOpaUpdatePolicyError.Wrap(err, "error while preparing opa data to json")
		o.log.Error(ctx, "error while updating  opa data", zap.Error(err))
		return err
	}

	serv := map[string]interface{}{
		"services": services,
	}
	servicedata, err := json.Marshal(serv)
	if err != nil {
		err := errors.ErrOpaUpdatePolicyError.Wrap(err, "error while preparing opa service data to json")
		o.log.Error(ctx, "error while updating  opa service data", zap.Error(err))
		return err
	}

	f, err := os.OpenFile(o.filepath, os.O_WRONLY|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err := errors.ErrOpaUpdatePolicyError.Wrap(err, "can not read  opa data")
		o.log.Error(ctx, "error while reading opa data", zap.Error(err))
		return err
	}

	defer f.Close()
	if _, err := f.WriteString(string(servicedata)); err != nil {
		err := errors.ErrOpaUpdatePolicyError.Wrap(err, "can not write new opa data")
		o.log.Error(ctx, "error while updating opa data", zap.Error(err))
		return err

	}
	time.Sleep(time.Second)
	o.log.Info(ctx, "Opa rule refreshed")
	return nil
}

func (o *opa) AllowedPermissions(ctx context.Context, input map[string]interface{}) (interface{}, error) {

	results, err := o.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		err := errors.ErrOpaEvalError.Wrap(err, "can not evaluate the user")
		o.log.Error(ctx, "error evaluating the user", zap.Error(err), zap.Any("input", input))
		return rego.ResultSet{}, err
	}
	return results[0].Expressions[0].Value, nil
}
