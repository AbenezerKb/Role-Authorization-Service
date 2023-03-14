package tenant

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

type tenant struct {
	tenantPersistant storage.Tenant
	log              logger.Logger
	opa              opa.Opa
}

func Init(log logger.Logger, tenantPersistant storage.Tenant, opa opa.Opa) module.Tenant {
	return &tenant{
		tenantPersistant: tenantPersistant,
		log:              log,
		opa:              opa,
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
func (t *tenant) UpdateTenantStatus(ctx context.Context, param dto.UpdateTenantStatus, tenant string) error {
	if err := param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	serviceID, err := uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid service id")
		t.log.Info(ctx, "invalid service id", zap.Error(err), zap.Any("service-id", ctx.Value("x-service-id")))
		return err
	}

	if err = t.tenantPersistant.UpdateTenantStatus(ctx, param, serviceID, tenant); err != nil {
		return err
	}

	return t.opa.Refresh(ctx, fmt.Sprintf("Updating tenant [%v] with status [%v]", tenant, param.Status))
}

func (t *tenant) GetTenantUsersWithRoles(ctx context.Context, query db_pgnflt.PgnFltQueryParams) ([]dto.TenantUserRoles, *model.MetaData, error) {
	param := dto.GetTenantUsersRequest{}
	var err error
	param.TenantName = ctx.Value("x-tenant").(string)
	param.ServiceID, err = uuid.Parse(ctx.Value("x-service-id").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("service id", ctx.Value("x-service-id")))
		return nil, nil, err
	}

	userID, err := uuid.Parse(ctx.Value("x-user").(string))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.log.Info(ctx, "invalid input", zap.Error(err), zap.Any("user id", ctx.Value("x-user")))
		return nil, nil, err
	}
	filter, err := query.ToFilterParams([]db_pgnflt.FieldType{
		{
			Name:   "user_id",
			Type:   db_pgnflt.String,
			DBName: "u.user_id",
		},
		{
			Name:   "created_at",
			Type:   db_pgnflt.Time,
			DBName: "u.created_at",
		},
		{
			Name:   "role_name",
			Type:   db_pgnflt.String,
			DBName: "rl.name",
		},
		{
			Name:   "role_status",
			Type:   db_pgnflt.Enum,
			Values: []string{constants.Active, constants.InActive},
			DBName: "tur.status",
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
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid filter params provided for tenant users roles")
		t.log.Warn(ctx, "invalid filter input", zap.Error(err))
		return nil, nil, err
	}

	tenantUserRoles, metadata, err := t.tenantPersistant.GetUsersWithTheirRoles(ctx, filter, param, userID)
	if err != nil {
		return nil, nil, err
	}
	return tenantUserRoles, metadata, nil

}
