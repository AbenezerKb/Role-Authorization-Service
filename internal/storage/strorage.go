package storage

import (
	"2f-authorization/internal/constants/model/dto"
	"context"

	"github.com/google/uuid"
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
	CreateDomain(ctx context.Context, param dto.CreateDomain) (*dto.Domain, error)
	IsDomainExist(ctx context.Context, param dto.CreateDomain) (bool, error)
	SoftDeleteDomain(ctx context.Context, param dto.DeleteDomain) error
}

type Permission interface {
	CreatePermission(ctx context.Context, param dto.CreatePermission) (uuid.UUID, error)
	AddToDomain(ctx context.Context, permissionId, domain uuid.UUID) error
	ListAllPermission(ctx context.Context, param dto.GetAllPermissionsReq) ([]dto.Permission, error)
}

type Tenant interface {
	CreateTenant(ctx context.Context, param dto.CreateTenent) error
	IsTenantExist(ctx context.Context, param dto.CreateTenent) (bool, error)
}
type User interface {
	RegiseterUser(ctx context.Context, param dto.RegisterUser) error
	IsUserExist(ctx context.Context, param dto.RegisterUser) (bool, error)
}

type Role interface {
	CreateRole(ctx context.Context, param dto.CreateRole) (*dto.Role, error)
	UpdateRole(ctx context.Context, param dto.UpdateRole) error
	RemovePermissionsFromRole(ctx context.Context, param dto.UpdateRole) error
	IsRoleExist(ctx context.Context, param dto.CreateRole) (bool, error)
	IsRoleAssigned(ctx context.Context, param dto.TenantUsersRole) (bool, error)
	AssignRole(ctx context.Context, param dto.TenantUsersRole) error
	RevokeRole(ctx context.Context, param dto.TenantUsersRole) error
	DeleteRole(ctx context.Context, roleId uuid.UUID) (*dto.Role, error)
}
