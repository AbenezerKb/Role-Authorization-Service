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

	"github.com/jackc/pgtype"
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
		DomainID:   param.DomainID,
	})
	if err != nil {
		err := errors.ErrWriteError.Wrap(err, "could not create tenant")
		t.log.Error(ctx, "unable to create tenant ", zap.Error(err), zap.Any("tenant", param))
		return err
	}
	return nil

}
func (t *tenant) CheckIfTenantExists(ctx context.Context, param dto.CreateTenent) (bool, error) {
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

func (t *tenant) RegsiterTenantPermission(ctx context.Context, tenant string, param dto.RegisterTenantPermission) (*dto.Permission, error) {
	statement, _ := param.Statement.Value()
	permission, err := t.db.TenantRegisterPermission(ctx, db.TenantRegisterPermissionParams{
		Name:        param.Name,
		ServiceID:   param.ServiceID,
		Description: param.Description,
		TenantName:  tenant,
		Statement: pgtype.JSON{
			Bytes:  statement,
			Status: pgtype.Present,
		},
		Column6: param.InheritedPermissions,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "tenant does not exists")
			t.log.Warn(ctx, "unable to find the tenant", zap.Error(err), zap.String(("tenant"), tenant))
			return nil, err
		}
		err := errors.ErrWriteError.Wrap(err, "unable to register the pemrission")
		t.log.Warn(ctx, "error registering the permission", zap.Error(err), zap.Any("permission", param), zap.String("tenant", tenant))
		return nil, err
	}

	st := dto.Statement{}
	if err := st.Scan(permission.Statement.Bytes); err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "unable to unmarshall the permission statement")
		t.log.Warn(ctx, "unable to unmarshall the permission statement", zap.Error(err), zap.Any("statement", permission.Statement), zap.String("tenant", tenant))
		return nil, err
	}

	return &dto.Permission{
		ID:          permission.ID,
		Description: permission.Description,
		Name:        permission.Name,
		Statement:   st,
		Tenant:      permission.Tenant,
		CreatedAt:   &permission.CreatedAt,
		ServiceID:   &permission.ServiceID,
	}, nil
}

func (t *tenant) CheckIfPermissionExistsInTenant(ctx context.Context, tenant string, param dto.RegisterTenantPermission) (bool, error) {
	count, err := t.db.CheckIfPermissionExistsInTenant(ctx, db.CheckIfPermissionExistsInTenantParams{
		Name:       param.Name,
		ServiceID:  param.ServiceID,
		TenantName: tenant,
	})
	if err != nil {
		err := errors.ErrReadError.Wrap(err, "could not read permission data")
		t.log.Error(ctx, "unable to read the permission data", zap.Error(err), zap.Any("param", param), zap.String("tenant", tenant))
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}
