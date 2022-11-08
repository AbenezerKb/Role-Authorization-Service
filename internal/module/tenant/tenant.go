package tenant

import (
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/module"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type tenant struct {
	tenantPersistant storage.Tenant
	log              logger.Logger
}

func Init(log logger.Logger, tenantPersistant storage.Tenant) module.Tenant {
	return &tenant{
		tenantPersistant: tenantPersistant,
		log:              log,
	}
}

func (t *tenant) CreateTenant(ctx context.Context, param dto.CreateTenent) error {

	var err error
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	if err = param.Validate(); err != nil {

		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}
	isTenantExist, err := t.tenantPersistant.CheckIfTenantExists(ctx, param)
	if err != nil {
		return err
	}

	if isTenantExist {
		t.log.Info(ctx, "tenant already exists", zap.String("name", param.TenantName))
		return errors.ErrDataExists.New("tenant with this name already exists")
	}
	return t.tenantPersistant.CreateTenant(ctx, param)
}

func (t *tenant) RegsiterTenantPermission(ctx context.Context, param dto.RegisterTenantPermission) (*dto.Permission, error) {
	var err error
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.log.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	tenant, ok := ctx.Value("x-tenant").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid input")
		t.log.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.log.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	if len(param.Statement.Fields) == 0 {
		param.Statement.Fields = []string{"*"}
	}

	exists, err := t.tenantPersistant.CheckIfPermissionExistsInTenant(ctx, tenant, param)
	if err != nil {
		return nil, err
	}

	if exists {
		t.log.Info(ctx, "permission already exists", zap.Any("param", param), zap.String("tenant", tenant))
		return nil, errors.ErrDataExists.New("permission with this name already exists")
	}

	permission, err := t.tenantPersistant.RegsiterTenantPermission(ctx, tenant, param)
	if err != nil {
		return nil, err
	}

	return permission, nil
}
