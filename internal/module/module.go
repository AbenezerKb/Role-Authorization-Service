package module

import (
	"2f-authorization/internal/constants/model"
	"2f-authorization/internal/constants/model/dto"
	"context"

	"github.com/google/uuid"
)

type Service interface {
	CreateService(ctx context.Context, param dto.CreateService) (*dto.CreateServiceResponse, error)
	DeleteService(ctx context.Context, param dto.Service) error
	UpdateServiceStatus(ctx context.Context, param dto.UpdateServiceStatus) error
}

type Domain interface {
	CreateDomain(ctx context.Context, param dto.CreateDomain) (*dto.Domain, error)
	DeleteDomain(ctx context.Context, param dto.DeleteDomain) error
}

type Permission interface {
	CreatePermission(ctx context.Context, param dto.CreatePermission) error
	ListPermissions(ctx context.Context) ([]dto.Permission, error)
	CreatePermissionDependency(ctx context.Context, param []dto.CreatePermissionDependency) error
}

type Tenant interface {
	CreateTenant(ctx context.Context, param dto.CreateTenent) error
	RegsiterTenantPermission(ctx context.Context, param dto.RegisterTenantPermission) (*dto.Permission, error)
}
type User interface {
	RegisterUser(ctx context.Context, param dto.RegisterUser) error
	UpdateUserStatus(ctx context.Context, param dto.UpdateUserStatus) error
	GetPermissionWithInTenant(ctx context.Context, tenant string, userId uuid.UUID) ([]dto.Permission, error)
}

type Role interface {
	CreateRole(ctx context.Context, param dto.CreateRole) (*dto.Role, error)
	UpdateRole(ctx context.Context, param dto.UpdateRole) error
	AssignRole(ctx context.Context, param dto.TenantUsersRole) error
	RevokeRole(ctx context.Context, param dto.TenantUsersRole) error
	DeleteRole(ctx context.Context, param string) (*dto.Role, error)
	ListRoles(ctx context.Context) ([]dto.Role, error)
	UpdateRoleStatus(ctx context.Context, param dto.UpdateRoleStatus, roleId uuid.UUID) error
	GetRole(ctx context.Context, param uuid.UUID) (*dto.Role, error)
}

type Opa interface {
	Authorize(ctx context.Context, req model.Request) (bool, error)
}
