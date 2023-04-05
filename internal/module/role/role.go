package role

import (
	"2f-authorization/internal/constants"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/module"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"2f-authorization/platform/opa"
	"context"
	"fmt"

	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"

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

	exists, err := r.rolePersistence.CheckIfRoleExists(ctx, param)
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

	serviceID, err := uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service", ctx.Value("x-service-id")))
		return err
	}

	tenant, ok := ctx.Value("x-tenant").(string)
	if ok {
		param.TenantName = tenant
		role, err := r.rolePersistence.GetRole(ctx, param.RoleID, serviceID)
		if err != nil {
			return err
		}

		if role.Name == "admin" {
			err := errors.ErrAcessError.New("Access denied")
			r.log.Info(ctx, "Access denied, Not eligible to assign admin role", zap.Error(err), zap.Any("tenant", ctx.Value("x-tenant")))
			return err
		}

	}

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
		r.log.Info(ctx, "role already exists", zap.String("role id", param.RoleID.String()), zap.String("role name", param.RoleName))
		return errors.ErrDataExists.Wrap(err, "user  with this role  already exists")
	}
	if param.RoleName == "admin" {
		if err := r.rolePersistence.RevokeAdminRole(ctx, param.TenantName); err != nil {
			return err
		}
	} else if param.RoleName != "admin" && param.RoleName != "" {
		err := errors.ErrAcessError.New("Access denied")
		r.log.Info(ctx, "Access denied, Not eligible to assign admin role", zap.Error(err), zap.Any("tenant", ctx.Value("x-tenant")))
		return err

	}

	if err := r.rolePersistence.AssignRole(ctx, serviceID, param); err != nil {
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

	return r.opa.Refresh(ctx, fmt.Sprintf("Updating [%v]  role", param.RoleID))
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

func (r *role) ListRoles(ctx context.Context, query db_pgnflt.PgnFltQueryParams) ([]dto.Role, *model.MetaData, error) {
	var err error
	param := dto.GetAllRolesReq{}
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service id", ctx.Value("x-service-id")))
		return nil, nil, err
	}
	var ok bool
	param.TenantName, ok = ctx.Value("x-tenant").(string)
	if !ok {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("tenant", ctx.Value("x-tenant")))
		return nil, nil, err
	}

	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		r.log.Info(ctx, "invalid input", zap.Error(err))
		return nil, nil, err
	}

	filter, err := query.ToFilterParams([]db_pgnflt.FieldType{
		{
			Name: "name",
			Type: db_pgnflt.String,
		},
		{
			Name: "created_at",
			Type: db_pgnflt.Time,
		},
		{
			Name: "updated_at",
			Type: db_pgnflt.Time,
		},
		{
			Name:   "status",
			Type:   db_pgnflt.Enum,
			Values: []string{constants.Active, constants.InActive},
		},
	}, db_pgnflt.Defaults{
		Sort: []db_pgnflt.Sort{
			{
				Field: "created_at",
				Sort:  db_pgnflt.SortDesc,
			},
		},
		PerPage: 10,
	})
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid filter params provided")
		r.log.Warn(ctx, "invalid filter input", zap.Error(err))
		return nil, nil, err
	}
	roles, metaData, err := r.rolePersistence.ListAllRoles(ctx, filter, param)
	if err != nil {
		return nil, nil, err
	}
	return roles, metaData, nil
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
	serviceID, err := uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid service id")
		r.log.Info(ctx, "invalid service id", zap.Error(err), zap.Any("service-id", ctx.Value("x-service-id")))
		return nil, err
	}
	return r.rolePersistence.GetRole(ctx, param, serviceID)
}
