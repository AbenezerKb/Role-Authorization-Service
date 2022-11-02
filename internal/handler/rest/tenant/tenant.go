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
// @Tags         tenant
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
// @Success      201  {object} dto.Permission "successfully register the permission"
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
