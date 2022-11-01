package opamodule

import (
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model"
	"2f-authorization/internal/module"
	"2f-authorization/platform/logger"
	"2f-authorization/platform/opa"

	"context"

	"go.uber.org/zap"
)

type opamodule struct {
	logger logger.Logger
	opa    opa.Opa
}

func Init(log logger.Logger, opa opa.Opa) module.Opa {
	return &opamodule{
		logger: log,
		opa:    opa,
	}
}

func (o *opamodule) Authorize(ctx context.Context, req model.Request) (bool, error) {
	var ok bool
	req.Service, ok = ctx.Value("x-service-id").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err), zap.Any("service id", ctx.Value("x-service-id")))
		return false, err
	}

	if len(req.Fields) == 0 {
		req.Fields = []string{"*"}
	}

	if err := req.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return false, err
	}
	
	ok, err := o.opa.Allow(ctx, req)
	if err != nil {
		err := errors.ErrAcessError.Wrap(err, "unable to perform operation")
		o.logger.Error(ctx, "error while enforcing policy", zap.Error(err), zap.String("service-id", req.Service), zap.String("user-id", req.Subject))
		return false, err
	}

	return ok, nil
}
