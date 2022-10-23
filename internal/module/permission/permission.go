package permission

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

type permission struct {
	log                   logger.Logger
	permissionPersistence storage.Permission
	opa                   opa.Opa
}

func Init(log logger.Logger, permissionPersistence storage.Permission, opa opa.Opa) module.Permission {
	return &permission{
		log:                   log,
		permissionPersistence: permissionPersistence,
		opa:                   opa,
	}
}

func (p *permission) CreatePermission(ctx context.Context, param dto.CreatePermission) error {

	var err error
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service id", ctx.Value("x-service-id")))
		return err
	}

	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}
	_, err = p.permissionPersistence.CreatePermission(ctx, param)
	if err != nil {
		return err
	}
	
	return nil
}

func (p *permission) ListPermissions(ctx context.Context) ([]dto.Permission, error) {

	var err error
	param := dto.GetAllPermissionsReq{}
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service id", ctx.Value("x-service-id")))
		return nil, err
	}
	var ok bool
	param.TenantName, ok = ctx.Value("x-tenant").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("tenant", ctx.Value("x-tenant")))
		return nil, err
	}

	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	permission, err := p.permissionPersistence.ListAllPermission(ctx, param)
	if err != nil {
		return nil, err
	}
	return permission, nil
}
