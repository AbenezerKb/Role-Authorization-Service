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

	"github.com/google/uuid"
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

func (s *service) CreateService(ctx context.Context, param dto.CreateService) (*dto.CreateServiceResponse, error) {
	service, err := s.db.CreateService(ctx, db.CreateServiceParams{
		Name:     param.Name,
		Password: param.Password,
		UserID:   uuid.MustParse(param.UserId),
	})
	if err != nil {
		err = errors.ErrWriteError.Wrap(err, "could not create service")
		s.log.Error(ctx, "unable to create service", zap.Error(err), zap.Any("service", param))
		return nil, err
	}
	return &dto.CreateServiceResponse{
		ServiceID:     service.ServiceID,
		Tenant:        service.Tenant,
		ServiceStatus: service.ServiceStatus,
		Service:       service.Service,
	}, nil
}

func (s *service) IsServiceExist(ctx context.Context, param dto.CreateService) (bool, error) {
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

func (s *service) SoftDeleteService(ctx context.Context, param dto.Service) error {
	if _, err := s.db.SoftDeleteService(ctx, param.ID); err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no record of service found")
			s.log.Info(ctx, "service not found", zap.Error(err), zap.String("service-id", param.ID.String()))
			return err
		}
		err = errors.ErrDBDelError.Wrap(err, "error deleting service")
		s.log.Error(ctx, "error deleting service", zap.Error(err), zap.String("service-id", param.ID.String()))
		return err
	}
	return nil
}

func (s *service) GetServiceById(ctx context.Context, param dto.Service) (*dto.Service, error) {
	service, err := s.db.GetServiceById(ctx, param.ID)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			err := errors.ErrNoRecordFound.Wrap(err, "no record of service found")
			s.log.Warn(ctx, "service not found", zap.Error(err), zap.String("service-id", param.ID.String()))
			return nil, err
		} else {
			err = errors.ErrReadError.Wrap(err, "error reading the service data")
			s.log.Error(ctx, "unable to get service data", zap.Error(err), zap.Any("service-name", param.Name))
			return nil, err
		}
	}
	return &dto.Service{
		ID:        service.ID,
		Status:    service.Status,
		Name:      service.Name,
		Password:  service.Password,
	}, nil
}
