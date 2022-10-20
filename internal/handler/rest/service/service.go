package service

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

type service struct {
	logger        logger.Logger
	serviceModule module.Service
}

func Init(logger logger.Logger, serviceModule module.Service) rest.Service {
	return &service{
		logger:        logger,
		serviceModule: serviceModule,
	}
}

// CreateService is used to create new service.
// @Summary      create new service.
// @Description  this function creates new service if it does not already exist.
// @Tags         service
// @Accept       json
// @Produce      json
// @param 		 createservice body dto.CreateService true "create service request body"
// @Success      201 {object} dto.CreateServiceResponse "successfully create new service"
// @Failure      400  {object}  model.ErrorResponse "required field error,bad request error"
// @Router       /services [post]
func (s *service) CreateService(ctx *gin.Context) {
	service := dto.CreateService{}
	err := ctx.ShouldBind(&service)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		s.logger.Info(ctx, "couldn't bind to dto.Service body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	createdService, err := s.serviceModule.CreateService(ctx, service)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusCreated, createdService, nil)
}

// DeletService is used to delete service.
// @Summary      deletes the service.
// @Description  this function deletes the service if it does already exist.
// @Tags         service
// @Accept       json
// @Produce      json
// @param 		 x-subject header string true "user id"
// @param 		 x-action header string true "action"
// @param 		 x-tenant header string true "tenant"
// @param 		 x-resource header string true "resource"
// @Success      200 boolean true "successfully deletes the service"
// @Failure      400  {object}  model.ErrorResponse "required field error"
// @Failure      404  {object}  model.ErrorResponse "service not found"
// @Failure      403  {object}  model.ErrorResponse "unauthorized"
// @Failure      401  {object}  model.ErrorResponse "service is not active"
// @Failure      500  {object}  model.ErrorResponse "invalid input"
// @Router       /services [delete]
// @security 	 BasicAuth
func (s *service) DeletService(ctx *gin.Context) {

	serviceId, err := uuid.Parse(ctx.GetString("x-service-id"))
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		s.logger.Info(ctx, "invalid input", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	if err := s.serviceModule.DeleteService(ctx, dto.Service{ID: serviceId}); err != nil {
		_ = ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, nil, nil)
}
