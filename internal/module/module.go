package module

import (
	"2f-authorization/internal/constants/model/dto"
	"context"
)

type Service interface {
	CreateService(ctx context.Context, param dto.CreateService) (*dto.CreateServiceResponse, error)
	DeleteService(ctx context.Context, param dto.Service) error
}

type Domain interface {
	CreateDomain(ctx context.Context, param dto.Domain) (*dto.Domain, error)
	DeleteDomain(ctx context.Context, param dto.Domain) error
}
