package service

import (
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/module"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"2f-authorization/platform/opa"
	"2f-authorization/platform/utils"
	"context"

	"go.uber.org/zap"
)

type service struct {
	log                logger.Logger
	servicePersistence storage.Service
	opa                opa.Opa
}

func Init(log logger.Logger, servicePersistence storage.Service, opa opa.Opa) module.Service {
	return &service{
		log:                log,
		servicePersistence: servicePersistence,
		opa:                opa,
	}
}

func (s *service) CreateService(ctx context.Context, param dto.CreateService) (*dto.CreateServiceResponse, error) {
	var err error
	if err = param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		s.log.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	exists, err := s.servicePersistence.IsServiceExist(ctx, param)
	if err != nil {
		return nil, err
	}

	if exists {
		s.log.Info(ctx, "service already exists", zap.String("name", param.Name))
		return nil, errors.ErrDataExists.Wrap(err, "service with this name already exists")
	}

	hashpass := utils.GenerateRandomString(10, true)
	if param.Password, err = utils.HashAndSalt(ctx, []byte(hashpass), s.log); err != nil {
		return nil, err
	}

	service, err := s.servicePersistence.CreateService(ctx, param)
	if err != nil {
		return nil, err
	}
	service.Password = hashpass
	return service, nil
}
