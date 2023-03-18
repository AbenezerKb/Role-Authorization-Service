package permission

import (
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/module"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"2f-authorization/platform/opa"
	"context"
	"fmt"

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

	serviceID, err := uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service id", ctx.Value("x-service-id")))
		return err
	}

	if len(param.Statement.Fields) == 0 {
		param.Statement.Fields = []string{"*"}
	}

	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}
	_, err = p.permissionPersistence.CreatePermission(ctx, param, serviceID)
	if err != nil {
		return err
	}

	return nil
}

func (p *permission) BulkCreatePermission(ctx context.Context, param []dto.CreatePermission) error {

	serviceID, err := uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service id", ctx.Value("x-service-id")))
		return err
	}

	for _, per := range param {

		if len(per.Statement.Fields) == 0 {
			per.Statement.Fields = []string{"*"}
		}

		if err = per.Validate(); err != nil {
			err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
			p.log.Info(ctx, "invalid input", zap.Error(err))
			return err
		}
		_, err = p.permissionPersistence.CreatePermission(ctx, per, serviceID)
		if err != nil {
			return err
		}

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

func (p *permission) CreatePermissionDependency(ctx context.Context, param []dto.CreatePermissionDependency) error {
	var err error
	serviceId, err := uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service id", ctx.Value("x-service-id")))
		return err
	}

	for _, v := range param {
		if err = v.Validate(); err != nil {
			err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
			p.log.Info(ctx, "invalid input", zap.Error(err))
			return err
		}

		if err := p.permissionPersistence.CreatePermissionDependency(ctx, v, serviceId); err != nil {
			return err
		}
	}

	if err := p.opa.Refresh(ctx, "Created an inheritance between permissions"); err != nil {
		return err
	}

	return nil
}

func (p *permission) DeletePermission(ctx context.Context, param string) error {

	serviceId, err := uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service id", ctx.Value("x-service-id")))
		return err
	}

	permissionId, err := uuid.Parse(param)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("permission-id", param))
		return err
	}

	tenantName, ok := ctx.Value("x-tenant").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid tenant")
		p.log.Info(ctx, "invalid tenant", zap.Error(err), zap.Any("tenant", ctx.Value("x-tenant")))
		return err
	}

	ok, err = p.permissionPersistence.CanBeDeletedOrUpdated(ctx, permissionId, serviceId)
	if err != nil {
		return err
	}

	if !ok {
		err := errors.ErrDBDelError.Wrap(err, "you can not delete this permission")
		p.log.Info(ctx, "unable to delete permission", zap.Error(err), zap.Any("permission-id", permissionId))
		return err
	}

	if err := p.permissionPersistence.DeletePermission(ctx, serviceId, permissionId, tenantName); err != nil {
		return err
	}
	return nil
}

func (p *permission) GetPermission(ctx context.Context, param uuid.UUID) (*dto.Permission, error) {
	serviceID, err := uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid service id")
		p.log.Info(ctx, "invalid service id", zap.Error(err), zap.Any("service-id", ctx.Value("x-service-id")))
		return nil, err
	}
	tenantName, ok := ctx.Value("x-tenant").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid tenant")
		p.log.Info(ctx, "invalid tenant", zap.Error(err), zap.Any("tenant", ctx.Value("x-tenant")))
		return nil, err
	}

	return p.permissionPersistence.GetPermission(ctx, param, serviceID, tenantName)
}

func (p *permission) UpdatePermissionStatus(ctx context.Context, param dto.UpdatePermissionStatus, permissionId uuid.UUID) error {
	serviceId, err := uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid service id")
		p.log.Info(ctx, "invalid service id", zap.Error(err), zap.Any("service-id", ctx.Value("x-service-id")))
		return err
	}

	tenant, ok := ctx.Value("x-tenant").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid tenant")
		p.log.Info(ctx, "invalid tenant", zap.Error(err), zap.Any("tenant", ctx.Value("x-tenant")))
		return err
	}

	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	ok, err = p.permissionPersistence.CanBeDeletedOrUpdated(ctx, permissionId, serviceId)
	if err != nil {
		return err
	}

	if !ok {
		err := errors.ErrDBDelError.Wrap(err, "you can not update this permission status")
		p.log.Info(ctx, "unable to update permission status", zap.Error(err), zap.Any("permission-id", permissionId))
		return err
	}

	if err = p.permissionPersistence.UpdatePermissionStatus(ctx, param, permissionId, serviceId, tenant); err != nil {
		return err
	}

	return p.opa.Refresh(ctx, fmt.Sprintf("Updating permission [%v] in tenant [%v] with status [%v]", permissionId.String(), tenant, param.Status))
}
