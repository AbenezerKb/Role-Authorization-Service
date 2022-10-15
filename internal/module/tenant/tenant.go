package tenant

import (
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/module"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"context"

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
	if err = param.Validate(); err != nil {

		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}
	isTenantExist, err := t.tenantPersistant.IsTenantExist(ctx, param)
	if err != nil {
		return err
	}

	if isTenantExist {
		t.log.Info(ctx, "tenant already exists", zap.String("name", param.TenantName))
		return errors.ErrDataExists.Wrap(err, "tenant with this name already exists")
	}
	return t.tenantPersistant.CreateTenant(ctx, param)
}
