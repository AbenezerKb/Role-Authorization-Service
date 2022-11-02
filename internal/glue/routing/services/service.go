package services

import (
	"2f-authorization/internal/glue/routing"
	"2f-authorization/internal/handler/middleware"
	"2f-authorization/internal/handler/rest"
	"2f-authorization/platform/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, service rest.Service, log logger.Logger, authMiddleware middleware.AuthMiddeleware) {
	services := group.Group("/services")
	servicesRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "",
			Handler:     service.CreateService,
			UnAuthorize: true,
		},
		{
			Method:      http.MethodDelete,
			Path:        "",
			Handler:     service.DeletService,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/status",
			Handler:     service.UpdateServiceStatus,
			UnAuthorize: true,
		},
	}
	routing.RegisterRoutes(services, servicesRoutes)
}
