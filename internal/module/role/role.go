package role

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

func (r *role) AssignRole(ctx context.Context, param dto.TenantUsersRole) error {

	var err error
	param.TenantName = ctx.Value("x-tenant").(string)

	if err = param.Validate(); err != nil {

		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}
	isExist, err := r.rolePersistence.IsRoleAssigned(ctx, param)
	if err != nil {
		return err
	}

	if isExist {
		r.log.Info(ctx, "role already exists", zap.String("name", param.RoleID.String()))
		return errors.ErrDataExists.Wrap(err, "user  with this role  already exists")
	}
	if err := r.rolePersistence.AssignRole(ctx, param); err != nil {
		return err
	}

	return r.opa.Refresh(ctx, fmt.Sprintf("Assigning [%v]  role  to user  - [%v]", param.RoleID, param.UserID))
}

func (r *role) RevokeRole(ctx context.Context, param dto.TenantUsersRole) error {

	var err error
	param.TenantName = ctx.Value("x-tenant").(string)

	if err = param.Validate(); err != nil {

		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}
	isExist, err := r.rolePersistence.IsRoleAssigned(ctx, param)
	if err != nil {
		return err
	}

	if !isExist {
		r.log.Info(ctx, "role does not exists", zap.String("role id ", param.RoleID.String()))
		return errors.ErrDataExists.Wrap(err, "user  with this role  does not  exists")
	}
	err = r.rolePersistence.RevokeRole(ctx, param)
	if err != nil {
		return err
	}
	if err := r.opa.Refresh(ctx, fmt.Sprintf("Revoking user role with role id  [%v]  and user id  - [%v]", param.RoleID, param.UserID)); err != nil {
		return err
	}
	return nil
}

func (r *role) UpdateRole(ctx context.Context, param dto.UpdateRole) error {

	if err := param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	if err := r.rolePersistence.RemovePermissionsFromRole(ctx, param); err != nil {
		return err
	}

	if err := r.rolePersistence.UpdateRole(ctx, param); err != nil {
		return err
	}

	if err := r.opa.Refresh(ctx, fmt.Sprintf("Updating [%v]  role", param.RoleID)); err != nil {
		return err
	}

	return nil
}

func (r *role) DeleteRole(ctx context.Context, param string) (*dto.Role, error) {
	roleId, err := uuid.Parse(param)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid id")
		r.log.Warn(ctx, "invalid input", zap.Error(err), zap.String("role-id", param))
		return nil, err
	}

	role, err := r.rolePersistence.DeleteRole(ctx, roleId)
	if err != nil {
		return nil, err
	}

	if err := r.opa.Refresh(ctx, fmt.Sprintf("Deleting role with id [%v]", param)); err != nil {
		return nil, err
	}

	return role, nil
}

func (r *role) ListRoles(ctx context.Context) ([]dto.Role, error) {
	var err error
	param := dto.GetAllRolesReq{}
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service id", ctx.Value("x-service-id")))
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
		r.log.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	roles, err := r.rolePersistence.ListAllRoles(ctx, param)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *role) UpdateRoleStatus(ctx context.Context, param dto.UpdateRoleStatus, roleId uuid.UUID) error {
	serviceID, err := uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid service id")
		r.log.Info(ctx, "invalid service id", zap.Error(err), zap.Any("service-id", ctx.Value("x-service-id")))
		return err
	}

	tenant, ok := ctx.Value("x-tenant").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid tenant")
		r.log.Info(ctx, "invalid tenant", zap.Error(err), zap.Any("tenant", ctx.Value("x-tenant")))
		return err
	}

	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	if err = r.rolePersistence.UpdateRoleStatus(ctx, param, roleId, serviceID, tenant); err != nil {
		return err
	}

	return r.opa.Refresh(ctx, fmt.Sprintf("Updating role [%v] in tenant [%v] with status [%v]", roleId.String(), tenant, param.Status))
}

func (r *role) GetRole(ctx context.Context, param uuid.UUID) (*dto.Role, error) {
	return r.rolePersistence.GetRole(ctx, param)
}
