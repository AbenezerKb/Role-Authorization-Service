package role

import (
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/module"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"2f-authorization/platform/opa"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type role struct {
	log             logger.Logger
	rolePersistence storage.Role
	opa             opa.Opa
}

func Init(log logger.Logger, rolePersistence storage.Role, opa opa.Opa) module.Role {
	return &role{
		log:             log,
		rolePersistence: rolePersistence,
		opa:             opa,
	}
}

func (r *role) CreateRole(ctx context.Context, param dto.CreateRole) (*dto.Role, error) {
	var err error
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service", ctx.Value("x-service-id")))
		return nil, err
	}
	
	var ok bool
	param.TenantName, ok = ctx.Value("x-tenant").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("tenant", ctx.Value("x-tenant")))
		return nil, err
	}

	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("input", param))
		return nil, err
	}

	exists, err := r.rolePersistence.IsRoleExist(ctx, param)
	if err != nil {
		return nil, err
	}
	if exists {
		err := errors.ErrDataExists.Wrap(err, "role with this name already exists")
		r.log.Info(ctx, "role with this name already exists", zap.Error(err), zap.Any("role", param))
		return nil, err
	}

	return r.rolePersistence.CreateRole(ctx, param)
}
