package tenant

import (
	"2f-authorization/internal/glue/routing"
	"2f-authorization/internal/handler/middleware"
	"2f-authorization/internal/handler/rest"
	"2f-authorization/platform/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, tenant rest.Tenant, log logger.Logger, authMiddleware middleware.AuthMiddeleware) {
	tenants := group.Group("/tenants")
	tenantRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "",
			Handler:     tenant.CreateTenant,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				// authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/permissions",
			Handler:     tenant.RegisterTenantPermission,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/:id/status",
			Handler:     tenant.UpdateTenantStatus,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				// authMiddleware.Authorize(),
			},
		}, {
			Method:      http.MethodGet,
			Path:        "/users",
			Handler:     tenant.GetUsersWithTheirRoles,
			UnAuthorize: false,
			Middlewares: []gin.HandlerFunc{
				// authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
	}
	routing.RegisterRoutes(tenants, tenantRoutes)
}
