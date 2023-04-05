package rest

import (
	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateService(ctx *gin.Context)
	DeletService(ctx *gin.Context)
	UpdateServiceStatus(ctx *gin.Context)
}
type Domain interface {
	CreateDomain(ctx *gin.Context)
	DeleteDomain(ctx *gin.Context)
}
type Permission interface {
	CreatePermission(ctx *gin.Context)
	BulkCreatePermission(ctx *gin.Context)
	ListPermissions(ctx *gin.Context)
	CreatePermissionDependency(ctx *gin.Context)
	DeletePermission(ctx *gin.Context)
	GetPermission(ctx *gin.Context)
	UpdatePermissionStatus(ctx *gin.Context)
}

type Tenant interface {
	CreateTenant(ctx *gin.Context)
	RegisterTenantPermission(ctx *gin.Context)
	UpdateTenantStatus(ctx *gin.Context)
	GetUsersWithTheirRoles(ctx *gin.Context)
}

type User interface {
	RegisterUser(ctx *gin.Context)
	UpdateUserStatus(ctx *gin.Context)
	GetPermissionWithInTenant(ctx *gin.Context)
	UpdateUserRoleStatus(ctx *gin.Context)
	SystemUpdateUserRoleStatus(ctx *gin.Context)
	GetPermissionWithInDomain(ctx *gin.Context)
}

type Role interface {
	CreateRole(ctx *gin.Context)
	AssignRole(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	RevokeRole(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
	ListRoles(ctx *gin.Context)
	UpdateRoleStatus(ctx *gin.Context)
	GetRole(ctx *gin.Context)
	SystemAssignRole(ctx *gin.Context)
}

type Opa interface {
	Authorize(ctx *gin.Context)
}
