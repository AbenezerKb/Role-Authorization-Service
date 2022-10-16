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
	CreateDomain(ctx context.Context, param dto.CreateDomain) (*dto.Domain, error)
	DeleteDomain(ctx context.Context, param dto.Domain) error
}

type Permission interface {
	CreatePermission(ctx context.Context, param dto.CreatePermission) error
}

type Tenant interface {
	CreateTenant(ctx context.Context, param dto.CreateTenent) error
}
type User interface {
	RegisterUser(ctx context.Context, param dto.RegisterUser) error
}
