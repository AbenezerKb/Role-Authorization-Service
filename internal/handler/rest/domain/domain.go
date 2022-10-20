package domain

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

type domain struct {
	logger       logger.Logger
	domainModule module.Domain
}

func Init(log logger.Logger, domainModule module.Domain) rest.Domain {
	return &domain{
		logger:       log,
		domainModule: domainModule,
	}
}

// CreateDomain is used to create new domain within the service.
// @Summary      create new domain.
// @Description  this function create new domain within the service if not exist.
// @Tags         domain
// @Accept       json
// @Produce      json
// @param 		 createdomain body dto.Domain true "create domain request body"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      201 {object} dto.Domain "successfully create new domain"
// @Failure      400  {object}  model.ErrorResponse "required field error,bad request error"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /domains [post]
// @security 	 BasicAuth
func (d *domain) CreateDomain(ctx *gin.Context) {

	domain := dto.CreateDomain{}
	if err := ctx.ShouldBind(&domain); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		d.logger.Info(ctx, "couldn't bind to dto.Domain body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	serviceId, err := uuid.Parse(ctx.GetString("x-service-id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		d.logger.Info(ctx, "invalid input", zap.Error(err), zap.String("service id", ctx.GetString("x-service-id")))
		_ = ctx.Error(err)
		return
	}
	domain.ServiceID = serviceId
	createdDomain, err := d.domainModule.CreateDomain(ctx, domain)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	constants.SuccessResponse(ctx, http.StatusCreated, createdDomain, nil)

}

// DeleteDomain is used to delete domain.
// @Summary      deletes the domain.
// @Description  this function deletes the domain if it does already exist.
// @Tags         domain
// @Accept       json
// @Produce      json
// @param 		 deletedomain body dto.DeleteDomain true "delete domain request body"
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200 boolean true "successfully deletes the service"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      404  {object}  model.ErrorResponse "service not found"
// @Failure      401  {object}  model.ErrorResponse "unauthorized"
// @Failure      403  {object}  model.ErrorResponse "access denied"
// @Router       /domains [delete]
// @security 	 BasicAuth
func (d *domain) DeleteDomain(ctx *gin.Context) {
	domain := dto.DeleteDomain{}

	if err := ctx.ShouldBind(&domain); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		d.logger.Info(ctx, "couldn't bind to dto.Domain body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	if err := d.domainModule.DeleteDomain(ctx, domain); err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}
