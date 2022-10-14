package storage

import (
	"2f-authorization/internal/constants/model/dto"
	"context"
)

type Policy interface {
	GetOpaData(ctx context.Context) ([]byte, error)
}

type Service interface {
	CreateService(ctx context.Context, param dto.CreateService) (*dto.CreateServiceResponse, error)
	IsServiceExist(ctx context.Context, param dto.CreateService) (bool, error)
	SoftDeleteService(ctx context.Context, param dto.Service) error
	GetServiceById(ctx context.Context, param dto.Service) (*dto.Service, error)
}

type Domain interface {
	CreateDomain(ctx context.Context, param dto.Domain) (*dto.Domain, error)
	IsDomainExist(ctx context.Context, param dto.Domain) (bool, error)
	SoftDeleteDomain(ctx context.Context, param dto.Domain) error
}
