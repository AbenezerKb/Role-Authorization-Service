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

func (t *tenant) CreateTenant(ctx *gin.Context) {

	tenant := dto.CreateTenent{}
	if err := ctx.ShouldBind(&tenant); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.logger.Info(ctx, "couldn't bind to dto.Tenant body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	serviceId, err := uuid.Parse(ctx.GetString("x-service-id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		t.logger.Info(ctx, "invalid input", zap.Error(err), zap.String("service id", ctx.GetString("x-service-id")))
		_ = ctx.Error(err)
		return
	}
	tenant.ServiceID = serviceId

	err = t.tenantModule.CreateTenant(ctx, tenant)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	constants.SuccessResponse(ctx, http.StatusCreated, nil, nil)

}
