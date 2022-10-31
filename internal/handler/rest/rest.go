package rest

import (
	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateService(ctx *gin.Context)
	DeletService(ctx *gin.Context)
}
type Domain interface {
	CreateDomain(ctx *gin.Context)
	DeleteDomain(ctx *gin.Context)
}
type Permission interface {
	CreatePermission(ctx *gin.Context)
	ListPermissions(ctx *gin.Context)
}

type Tenant interface {
	CreateTenant(ctx *gin.Context)
}

type User interface {
	RegisterUser(ctx *gin.Context)
}

type Role interface {
	CreateRole(ctx *gin.Context)
	AssignRole(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	RevokeRole(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
}

type Opa interface {
	Authorize(ctx *gin.Context)
}
