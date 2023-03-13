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
	"fmt"

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

	exists, err := s.servicePersistence.CheckIfServiceExists(ctx, param)
	if err != nil {
		return nil, err
	}

	if exists {
		s.log.Info(ctx, "service already exists", zap.String("name", param.Name))
		return nil, errors.ErrDataExists.Wrap(err, "service with this name already exists")
	}

	generatedPass := utils.GenerateRandomString(20, true)
	param.Password = generatedPass

	service, err := s.servicePersistence.CreateService(ctx, param)
	if err != nil {
		return nil, err
	}
	service.Password = generatedPass
	return service, nil
}

func (s *service) DeleteService(ctx context.Context, param dto.Service) error {
	if err := s.servicePersistence.SoftDeleteService(ctx, param); err != nil {
		return err
	}

	if err := s.opa.Refresh(ctx, fmt.Sprintf("Removed service with id - [%v]", param.ID)); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateServiceStatus(ctx context.Context, param dto.UpdateServiceStatus) error {

	if err := param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		s.log.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	if err := s.servicePersistence.UpdateServicePersistence(ctx, param); err != nil {
		return err
	}

	if err := s.opa.Refresh(ctx, fmt.Sprintf("Updating service [%v] with status [%v]", param.ServiceID.String(), param.Status)); err != nil {
		return err
	}

	return nil
}
