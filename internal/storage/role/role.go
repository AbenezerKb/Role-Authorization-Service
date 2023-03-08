package role

import (
	"2f-authorization/internal/constants/dbinstance"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/error/sqlcerr"
	"2f-authorization/internal/constants/model"
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"context"

	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"

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
		CreatedAt: &role.CreatedAt,
		Status:    string(role.Status),
	}, nil
}
func (r *role) CheckIfRoleExists(ctx context.Context, param dto.CreateRole) (bool, error) {
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

func (r *role) AssignRole(ctx context.Context, serviceID uuid.UUID, param dto.TenantUsersRole) error {

	err := r.db.AssignRole(ctx, db.AssignRoleParams{
		TenantName: param.TenantName,
		ID:         param.RoleID,
		Name:       param.RoleName,
		UserID:     param.UserID,
		ServiceID:  serviceID,
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

	return count.(int64) > 0, nil
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
		CreatedAt: &role.CreatedAt,
	}, nil
}

func (r *role) ListAllRoles(ctx context.Context, filter db_pgnflt.FilterParams, param dto.GetAllRolesReq) ([]dto.Role, *model.MetaData, error) {
	roles, metaData, err := r.db.ListRoles(ctx, filter, dbinstance.ListRolesParams{
		ServiceID:  param.ServiceID,
		TenantName: param.TenantName,
	})
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no roles found")
			r.log.Info(ctx, "no roles were found", zap.Error(err), zap.String("tenany-name", param.TenantName), zap.String("service-id", param.ServiceID.String()))
			return []dto.Role{}, nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading roles")
			r.log.Error(ctx, "error reading roles", zap.Error(err), zap.String("tenany-name", param.TenantName), zap.String("service-id", param.ServiceID.String()))
			return []dto.Role{}, nil, err
		}
	}
	return roles, metaData, nil
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
		CreatedAt:   &role.CreatedAt,
		Permissions: role.Permission,
	}, nil
}

func (r *role) RevokeAdminRole(ctx context.Context, tenantID uuid.UUID) error {
	err := r.db.RevokeAdminRole(ctx, tenantID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "admin role not found")
			r.log.Error(ctx, "admin role not found", zap.Error(err),
				zap.String("tenant", tenantID.String()))
			return err
		}

		err = errors.ErrUpdateError.Wrap(err, "error revoking admin role's status")
		r.log.Error(ctx, "error revoking admin role's status",
			zap.Error(err),
			zap.String("tenant", tenantID.String()))
		return err
	}
	return nil
}
