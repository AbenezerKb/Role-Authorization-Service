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

func (s *service) CreateService(ctx *gin.Context) {
	service := dto.Service{}
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
