package rest

import "github.com/gin-gonic/gin"

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
}
