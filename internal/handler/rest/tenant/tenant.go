package tenant

import (
	"2f-authorization/internal/constants"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/handler/rest"
	"2f-authorization/internal/module"
	"2f-authorization/platform/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
	"go.uber.org/zap"
)

type tenant struct {
	logger       logger.Logger
	tenantModule module.Tenant
}

func Init(log logger.Logger, tenantModule module.Tenant) rest.Tenant {
	return &tenant{
		logger:       log,
		tenantModule: tenantModule,
	}
}

// CreateTenant is used to create tenant.
// @Summary      create tenant.
// @Description  this function create tenant if it is not exist in the service.
// @Tags         tenants
// @Accept       json
// @Produce      json
// @param 		 createtenant body dto.CreateTenent true "create tenant request body"
// @Success      201  boolean true "successfully create new tenant"
// @Failure      400  {object}  model.ErrorResponse "required field error,bad request error"
// @Router       /tenants [post]
// @security 	 BasicAuth
func (t *tenant) CreateTenant(ctx *gin.Context) {

	tenant := dto.CreateTenent{}
	if err := ctx.ShouldBind(&tenant); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.logger.Info(ctx, "couldn't bind to dto.Tenant body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	err := t.tenantModule.CreateTenant(ctx, tenant)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	constants.SuccessResponse(ctx, http.StatusCreated, nil, nil)

}

// RegisterTenantPermission is used to register new permissions under the tenant.
// @Summary      register a new permission.
// @Tags         tenants
// @Accept       json
// @Produce      json
// @param 		 creatnewpermission body dto.RegisterTenantPermission true "new permission request body"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      201  {object} dto.Permission "successfully registered the permission"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized service"
// @Failure      403  {object}  model.ErrorResponse "service is not active"
// @Router       /tenants/permissions [post]
// @security 	 BasicAuth
func (t *tenant) RegisterTenantPermission(ctx *gin.Context) {
	permission := dto.RegisterTenantPermission{}
	err := ctx.ShouldBind(&permission)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.logger.Info(ctx, "couldn't bind to dto.Service body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	result, err := t.tenantModule.RegsiterTenantPermission(ctx, permission)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusCreated, result, nil)
}

// UpdateTenantStatus updates tenant status
// @Summary      changes tenant status
// @Tags         tenants
// @Accept       json
// @Produce      json
// @param status body dto.UpdateTenantStatus true "status"
// @param 		 id 	path string true "tenant id"
// @Success      200 boolean true "successfully updated the tenant's status"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Router       /tenants/{id}/status [patch]
// @Security	 BasicAuth
func (t *tenant) UpdateTenantStatus(ctx *gin.Context) {
	updateStatusParam := dto.UpdateTenantStatus{}
	err := ctx.ShouldBind(&updateStatusParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.logger.Info(ctx, "unable to bind tenant status", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	if err := t.tenantModule.UpdateTenantStatus(ctx, updateStatusParam, ctx.Param("id")); err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// GetUsersWithTheirRoles get tenant users with their roles
// @Summary      get Tenant Users with their roles
// @Tags         tenants
// @Produce      json
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200 boolean true "successfully get user's roles"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Router       /users [get]
// @Security	 BasicAuth
func (t *tenant) GetUsersWithTheirRoles(ctx *gin.Context) {
	param := db_pgnflt.PgnFltQueryParams{}

	err := ctx.BindQuery(&param)
	if err != nil {
		er := errors.ErrInvalidUserInput.Wrap(err, "invalid user input")
		t.logger.Info(ctx, "unable to bind tenant query", zap.Error(err))
		_ = ctx.Error(er)
		return
	}
	tenantUserRoles, metadata, err := t.tenantModule.GetTenantUsersWithRoles(ctx, param)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	constants.SuccessResponse(ctx, http.StatusOK, tenantUserRoles, metadata)

}

// UpdateCorporateUserRoleStatus updates corporate user's role status
// @Summary      updates corporate user's role status
// @Tags         users
// @Accept       json
// @Produce      json
// @param status body dto.UpdateUserRoleStatus true "status"
// @param 		 corporate-id 	path string true "role id"
// @param 		 role-id 	path string true "role id"
// @param 		 id 	path string true "user id"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200 boolean true "successfully updates the user's role status"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       corporate/{corporate-id}/{user-id}/roles/{role-id}/status [patch]
// @Security	 BasicAuth
func (u *tenant) UpdateCorporateUserRoleStatus(ctx *gin.Context) {
	updateStatusParam := dto.UpdateUserRoleStatus{}
	err := ctx.ShouldBind(&updateStatusParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		u.logger.Info(ctx, "unable to bind role status", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	roleId, err := uuid.Parse(ctx.Param("role-id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid role id")
		u.logger.Info(ctx, "invalid role id", zap.Error(err), zap.Any("role id", ctx.Param("role-id")))
		_ = ctx.Error(err)
		return
	}
	userId, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid user id")
		u.logger.Info(ctx, "invalid role id", zap.Error(err), zap.Any("user id", ctx.Param("id")))
		_ = ctx.Error(err)
		return
	}
	if err := u.tenantModule.UpdateCorporateUserRoleStatus(ctx, updateStatusParam, ctx.Param("corporate-id"), roleId, userId); err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}
