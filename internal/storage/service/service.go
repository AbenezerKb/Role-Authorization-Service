package service

import (
	"2f-authorization/internal/constants/dbinstance"
	errors "2f-authorization/internal/constants/error"
	"2f-authorization/internal/constants/error/sqlcerr"
	"2f-authorization/internal/constants/model/db"
	"2f-authorization/internal/constants/model/dto"
	"2f-authorization/internal/storage"
	"2f-authorization/platform/logger"
	"context"

	"go.uber.org/zap"
)

type service struct {
	db  dbinstance.DBInstance
	log logger.Logger
}

func Init(db dbinstance.DBInstance, log logger.Logger) storage.Service {
	return &service{
		db:  db,
		log: log,
	}
}

func (s *service) CreateService(ctx context.Context, param dto.Service) (*dto.Service, error) {
	service, err := s.db.CreateService(ctx, db.CreateServiceParams{
		Name:     param.Name,
		Password: param.Password,
	})
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not create service")
		s.log.Error(ctx, "unable to create service", zap.Error(err), zap.Any("service", param))
		return nil, err
	}
	return &dto.Service{
		ID:        service.ID,
		Name:      service.Name,
		Password:  service.Password,
		CreatedAt: service.CreatedAt,
		UpdatedAt: service.UpdatedAt,
		Status:    service.Status,
	}, nil
}

func (s *service) IsServiceExist(ctx context.Context, param dto.Service) (bool, error) {
	_, err := s.db.GetServiceByName(ctx, param.Name)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return false, nil
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading the service data")
			s.log.Error(ctx, "unable to get service data", zap.Error(err), zap.Any("service-name", param.Name))
			return false, err
		}
	}
	return true, nil
}
