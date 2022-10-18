package role

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

type role struct {
	db  dbinstance.DBInstance
	log logger.Logger
}

func Init(db dbinstance.DBInstance, log logger.Logger) storage.Role {
	return &role{
		db:  db,
		log: log,
	}
}

func (r *role) CreateRole(ctx context.Context, param dto.CreateRole) (*dto.Role, error) {

	role, err := r.db.CreateRole(ctx, db.CreateRoleParams{
		Name:       param.Name,
		TenantName: param.TenantName,
		ServiceID:  param.ServiceID,
		Column4:    param.PermissionID,
	})
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not create role")
		r.log.Error(ctx, "unable to create role", zap.Error(err), zap.Any("role", param))
		return nil, err
	}
	return &dto.Role{
		Name:      role.Name,
		ID:        role.RoleID,
		CreatedAt: role.CreatedAt,
		Status:    role.Status,
	}, nil
}
func (r *role) IsRoleExist(ctx context.Context, param dto.CreateRole) (bool, error) {
	_, err := r.db.GetRoleByNameAndTenantName(ctx, db.GetRoleByNameAndTenantNameParams{
		Name:       param.Name,
		TenantName: param.TenantName,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return false, nil
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading the role data")
			r.log.Error(ctx, "unable to get role data", zap.Error(err), zap.Any("tenant-name", param.TenantName), zap.Any("role", param.Name))
			return false, err
		}
	}
	return true, nil
}
