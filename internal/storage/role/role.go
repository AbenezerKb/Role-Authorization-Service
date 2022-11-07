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

	"github.com/google/uuid"
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
		Status:    string(role.Status),
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

func (r *role) AssignRole(ctx context.Context, param dto.TenantUsersRole) error {

	err := r.db.AssignRole(ctx, db.AssignRoleParams{
		TenantName: param.TenantName,
		RoleID:     param.RoleID,
		UserID:     param.UserID,
	})
	if err != nil {
		err := errors.ErrWriteError.Wrap(err, "could not assign role")
		r.log.Error(ctx, "unable to assign role ", zap.Error(err), zap.Any("role", param))
		return err
	}
	return nil

}
func (r *role) IsRoleAssigned(ctx context.Context, param dto.TenantUsersRole) (bool, error) {
	count, err := r.db.IsRoleAssigned(ctx, db.IsRoleAssignedParams{
		TenantName: param.TenantName,
		UserID:     param.UserID,
		RoleID:     param.RoleID,
	})
	if err != nil {
		err := errors.ErrReadError.Wrap(err, "could not  read role")
		r.log.Error(ctx, "unable to read the role", zap.Error(err), zap.Any("role ", param))
		return false, err
	}

	if count.(int64) > 0 {
		return true, nil
	}
	return false, nil
}

func (r *role) UpdateRole(ctx context.Context, param dto.UpdateRole) error {
	if err := r.db.UpdateRole(ctx, db.UpdateRoleParams{
		RoleID:  param.RoleID,
		Column2: param.PermissionsID,
	}); err != nil {
		err := errors.ErrWriteError.Wrap(err, "error updating the role")
		r.log.Error(ctx, "unable to update the role", zap.Error(err), zap.String("role-id", param.RoleID.String()), zap.Any("permissions", param.PermissionsID))
		return err
	}
	return nil
}

func (r *role) RemovePermissionsFromRole(ctx context.Context, param dto.UpdateRole) error {
	if err := r.db.RemovePermissionsFromRole(ctx, db.RemovePermissionsFromRoleParams{
		RoleID:  param.RoleID,
		Column2: param.PermissionsID,
	}); err != nil {
		err := errors.ErrDBDelError.Wrap(err, "error removing the permissions from the role")
		r.log.Error(ctx, "unable to remove the permissions", zap.Error(err), zap.String("role-id", param.RoleID.String()), zap.Any("permissions", param.PermissionsID))
		return err
	}
	return nil
}
func (r *role) RevokeRole(ctx context.Context, param dto.TenantUsersRole) error {

	err := r.db.RevokeUserRole(ctx, db.RevokeUserRoleParams{
		TenantName: param.TenantName,
		UserID:     param.UserID,
		RoleID:     param.RoleID,
	})
	if err != nil {
		err := errors.ErrWriteError.Wrap(err, "could not revoke role")
		r.log.Error(ctx, "unable to revoke role", zap.Error(err), zap.Any("role", param))
		return err
	}
	return nil

}

func (r *role) DeleteRole(ctx context.Context, roleId uuid.UUID) (*dto.Role, error) {
	role, err := r.db.DeleteRole(ctx, roleId)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "role does not exists")
			r.log.Warn(ctx, "unable to find the role", zap.Error(err), zap.String(("role-id"), roleId.String()))
			return nil, err
		}
		err := errors.ErrDBDelError.Wrap(err, "unable to delete the role")
		r.log.Warn(ctx, "unable to delete the role", zap.Error(err), zap.String("role-id", roleId.String()))
		return nil, err
	}
	return &dto.Role{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}, nil
}

func (r *role) ListAllRoles(ctx context.Context, param dto.GetAllRolesReq) ([]dto.Role, error) {
	roles, err := r.db.ListRoles(ctx, dbinstance.ListRolesParams{
		ServiceID:  param.ServiceID,
		TenantName: param.TenantName,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no roles found")
			r.log.Info(ctx, "no roles were found", zap.Error(err), zap.String("tenany-name", param.TenantName), zap.String("service-id", param.ServiceID.String()))
			return []dto.Role{}, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading roles")
			r.log.Error(ctx, "error reading roles", zap.Error(err), zap.String("tenany-name", param.TenantName), zap.String("service-id", param.ServiceID.String()))
			return []dto.Role{}, err
		}
	}
	return roles, nil
}

func (r *role) UpdateRoleStatus(ctx context.Context, param dto.UpdateRoleStatus, roleId, serviceId uuid.UUID, tenant string) error {
	_, err := r.db.UpdateRoleStatus(ctx, db.UpdateRoleStatusParams{
		ID:         roleId,
		TenantName: tenant,
		ServiceID:  serviceId,
		Status:     db.Status(param.Status),
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "role not found")
			r.log.Error(ctx, "error changing role's status", zap.Error(err), zap.String("service", serviceId.String()), zap.String("role-status", param.Status), zap.String("role-id", roleId.String()), zap.String("tenant", tenant))
			return err
		}

		err = errors.ErrUpdateError.Wrap(err, "error changing role's status")
		r.log.Error(ctx, "error changing role's status", zap.Error(err), zap.String("service", serviceId.String()), zap.String("role-status", param.Status), zap.String("role-id", roleId.String()), zap.String("tenant", tenant))
		return err
	}
	return nil
}

func (r *role) GetRole(ctx context.Context, param uuid.UUID, serviceId uuid.UUID) (*dto.Role, error) {
	role, err := r.db.GetRoleById(ctx, db.GetRoleByIdParams{
		ServiceID: serviceId,
		ID:        param,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "role not found")
			r.log.Error(ctx, "error getting the role's data", zap.Error(err), zap.String("role-id", param.String()))
			return nil, err
		}

		err = errors.ErrUpdateError.Wrap(err, "error getting the role's data")
		r.log.Error(ctx, "error getting the role's data", zap.Error(err), zap.String("role-id", param.String()))
	}
	return &dto.Role{
		Name:        role.Name,
		Status:      string(role.Status),
		ID:          role.ID,
		CreatedAt:   role.CreatedAt,
		Permissions: role.Permission,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}
