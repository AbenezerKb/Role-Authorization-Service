package domain

import (
	"2f-authorization/internal/glue/routing"
	"2f-authorization/internal/handler/middleware"
	"2f-authorization/internal/handler/rest"
	"2f-authorization/platform/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, domain rest.Domain, log logger.Logger, authMiddleware middleware.AuthMiddeleware) {
	domains := group.Group("/domains")
	domainRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "",
			Handler:     domain.CreateDomain,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				// authMiddleware.Authorize(),
			},
		},
	}
	routing.RegisterRoutes(domains, domainRoutes)
}
