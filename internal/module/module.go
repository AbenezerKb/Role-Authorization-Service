package module

import (
	"2f-authorization/internal/constants/model"
	"2f-authorization/internal/constants/model/dto"
	"context"
)

type Service interface {
	CreateService(ctx context.Context, param dto.CreateService) (*dto.CreateServiceResponse, error)
	DeleteService(ctx context.Context, param dto.Service) error
}

type Domain interface {
	CreateDomain(ctx context.Context, param dto.CreateDomain) (*dto.Domain, error)
	DeleteDomain(ctx context.Context, param dto.DeleteDomain) error
}

type Permission interface {
	CreatePermission(ctx context.Context, param dto.CreatePermission) error
	ListPermissions(ctx context.Context) ([]dto.Permission, error)
}

type Tenant interface {
	CreateTenant(ctx context.Context, param dto.CreateTenent) error
}
type User interface {
	RegisterUser(ctx context.Context, param dto.RegisterUser) error
}

type Role interface {
	CreateRole(ctx context.Context, param dto.CreateRole) (*dto.Role, error)
	UpdateRole(ctx context.Context, param dto.UpdateRole) error
	AssignRole(ctx context.Context, param dto.TenantUsersRole) error
	RevokeRole(ctx context.Context, param dto.TenantUsersRole) error
	DeleteRole(ctx context.Context, param string) (*dto.Role, error)
}

type Opa interface {
	Authorize(ctx context.Context, req model.Request) (bool, error)
}
