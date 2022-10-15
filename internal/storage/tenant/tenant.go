package tenant

import (
	"2f-authorization/internal/constants/dbinstance"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/error/sqlcerr"
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"context"

	"go.uber.org/zap"
)

type tenant struct {
	db  dbinstance.DBInstance
	log logger.Logger
}

func Init(db dbinstance.DBInstance, log logger.Logger) storage.Tenant {
	return &tenant{
		db:  db,
		log: log,
	}
}
func (t *tenant) CreateTenant(ctx context.Context, param dto.CreateTenent) error {

	err := t.db.CreateTenent(ctx, db.CreateTenentParams{
		TenantName: param.TenantName,
		ServiceID:  param.ServiceID,
	})
	if err != nil {
		err := errors.ErrWriteError.Wrap(err, "could not create tenant")
		t.log.Error(ctx, "unable to create tenant ", zap.Error(err), zap.Any("tenant", param))
		return err
	}
	return nil

}
func (t *tenant) IsTenantExist(ctx context.Context, param dto.CreateTenent) (bool, error) {
	_, err := t.db.GetTenentWithNameAndServiceId(ctx, db.GetTenentWithNameAndServiceIdParams{
		TenantName: param.TenantName,
		ServiceID:  param.ServiceID,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return false, nil
		} else {
			err := errors.ErrReadError.Wrap(err, "could not  read tenant")
			t.log.Error(ctx, "unable to read the tenant", zap.Error(err), zap.Any("tenant ", param))
			return false, err

		}

	}
	return true, nil
}
