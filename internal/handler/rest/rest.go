package rest

import "github.com/gin-gonic/gin"

type Service interface {
	CreateService(ctx *gin.Context)
}
