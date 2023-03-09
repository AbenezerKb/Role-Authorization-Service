package opa

import (
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model"
	"2f-authorization/platform/logger"
	"context"

	dbstore "2f-authorization/internal/storage"

	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
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
	db     dbstore.Policy
	store  storage.Store
	policy string
	Query  string
	log    logger.Logger
	query  rego.PreparedEvalQuery
}

func Init(policy string, policyDb dbstore.Policy, log logger.Logger) Opa {
	return &opa{
		policy: policy,
		db:     policyDb,
		log:    log,
	}
}

func (o *opa) Allow(ctx context.Context, req model.Request) (bool, error) {

	input := map[string]interface{}{
		"subject":  req.Subject,
		"resource": req.Resource,
		"action":   req.Action,
		"service":  req.Service,
		"tenant":   req.Tenant,
		"fields":   req.Fields,
	}

	results, err := o.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		err := errors.ErrOpaEvalError.Wrap(err, "can not evaluate the user")
		o.log.Error(ctx, "error evaluating the user", zap.Error(err), zap.Any("input", input))
		return false, err
	}

	return results.Allowed(), nil
}

func (o *opa) Refresh(ctx context.Context, reason string) error {
	o.log.Info(ctx, reason)
	if err := o.GetData(ctx); err != nil {
		return err
	}

	qr, err := rego.New(
		rego.Query("data.authz.allow"),
		rego.Store(o.store),
		rego.Module("authz.rego", o.policy),
	).PrepareForEval(ctx)

	if err != nil {
		err := errors.ErrOpaPrepareEvalError.Wrap(err, "error preparing for evaluation")
		o.log.Error(ctx, "error preparing the rego for eval", zap.Error(err))
		return err
	}

	o.query = qr
	o.log.Info(ctx, "successfully triggered policy data update")
	return nil
}

func (o *opa) GetData(ctx context.Context) error {
	data, err := o.db.GetOpaData(ctx)
	if err != nil {
		return err
	}
	var services map[string]interface{}
	_ = util.UnmarshalJSON(data, &services)
	o.store = inmem.NewFromObject(map[string]interface{}{
		"services": services,
	})
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
