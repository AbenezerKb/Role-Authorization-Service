package storage

import (
	"2f-authorization/internal/constants/model"
	"2f-authorization/internal/constants/model/dto"
	"context"

	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"

	"github.com/google/uuid"
)

type Policy interface {
	GetOpaData(ctx context.Context) ([]byte, error)
}

type Service interface {
	CreateService(ctx context.Context, param dto.CreateService) (*dto.CreateServiceResponse, error)
	CheckIfServiceExists(ctx context.Context, param dto.CreateService) (bool, error)
	SoftDeleteService(ctx context.Context, param dto.Service) error
	GetServiceById(ctx context.Context, param dto.Service) (*dto.Service, error)
	UpdateServicePersistence(ctx context.Context, param dto.UpdateServiceStatus) error
}

type Domain interface {
	CreateDomain(ctx context.Context, param dto.CreateDomain) (*dto.Domain, error)
	CheckIfDomainExists(ctx context.Context, param dto.CreateDomain) (bool, error)
	SoftDeleteDomain(ctx context.Context, param dto.DeleteDomain) error
}

type Permission interface {
	CreatePermission(ctx context.Context, param dto.CreatePermission) (uuid.UUID, error)
	AddToDomain(ctx context.Context, permissionId, domain uuid.UUID) error
	ListAllPermission(ctx context.Context, param dto.GetAllPermissionsReq) ([]dto.Permission, error)
	CreatePermissionDependency(ctx context.Context, param dto.CreatePermissionDependency, serviceId uuid.UUID) error
	DeletePermission(ctx context.Context, serviceId, permissionId uuid.UUID, tenantName string) error
	CanBeDeletedOrUpdated(ctx context.Context, permissionId, serviceId uuid.UUID) (bool, error)
	GetPermission(ctx context.Context, permissionId, serviceID uuid.UUID, tenantName string) (*dto.Permission, error)
	UpdatePermissionStatus(ctx context.Context, param dto.UpdatePermissionStatus, permissionId, serviceId uuid.UUID, tenant string) error
}

type Tenant interface {
	CreateTenant(ctx context.Context, param dto.CreateTenent) error
	CheckIfTenantExists(ctx context.Context, param dto.CreateTenent) (bool, error)
	CheckIfPermissionExistsInTenant(ctx context.Context, tenant string, param dto.RegisterTenantPermission) (bool, error)
	RegsiterTenantPermission(ctx context.Context, tenant string, param dto.RegisterTenantPermission) (*dto.Permission, error)
	UpdateTenantStatus(ctx context.Context, param dto.UpdateTenantStatus, serviceId uuid.UUID, tenant string) error
	GetUsersWithTheirRoles(ctx context.Context, filter db_pgnflt.FilterParams, param dto.GetTenantUsersRequest) ([]dto.TenantUserRoles, *model.MetaData, error)
}
type User interface {
	RegiseterUser(ctx context.Context, param dto.RegisterUser) error
	CheckIfUserExists(ctx context.Context, param dto.RegisterUser) (bool, error)
	UpdateUserStatus(ctx context.Context, param dto.UpdateUserStatus) error
	GetPermissionWithInTenant(ctx context.Context, tenant string, userId, serviceID uuid.UUID) ([]dto.Permission, error)
	GetPermissionWithInDomain(ctx context.Context, domain, userId, serviceID uuid.UUID) ([]dto.DomainPermissions, error)
	UpdateUserRoleStatus(ctx context.Context, param dto.UpdateUserRoleStatus, roleId, userId, serviceId uuid.UUID, tenant string) error
}

type Role interface {
	CreateRole(ctx context.Context, param dto.CreateRole) (*dto.Role, error)
	UpdateRole(ctx context.Context, param dto.UpdateRole) error
	RemovePermissionsFromRole(ctx context.Context, param dto.UpdateRole) error
	CheckIfRoleExists(ctx context.Context, param dto.CreateRole) (bool, error)
	IsRoleAssigned(ctx context.Context, param dto.TenantUsersRole) (bool, error)
	AssignRole(ctx context.Context, serviceID uuid.UUID, param dto.TenantUsersRole) error
	RevokeRole(ctx context.Context, param dto.TenantUsersRole) error
	DeleteRole(ctx context.Context, roleId uuid.UUID) (*dto.Role, error)
	ListAllRoles(ctx context.Context, filter db_pgnflt.FilterParams, param dto.GetAllRolesReq) ([]dto.Role, *model.MetaData, error)
	UpdateRoleStatus(ctx context.Context, param dto.UpdateRoleStatus, roleId, serviceId uuid.UUID, tenant string) error
	GetRole(ctx context.Context, param uuid.UUID, serviceId uuid.UUID) (*dto.Role, error)
	RevokeAdminRole(ctx context.Context, param uuid.UUID) error
}
