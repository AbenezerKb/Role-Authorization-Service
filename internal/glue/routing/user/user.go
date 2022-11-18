package user

import (
	"2f-authorization/internal/glue/routing"
	"2f-authorization/internal/handler/middleware"
	"2f-authorization/internal/handler/rest"
	"2f-authorization/platform/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, user rest.User, log logger.Logger, authMiddleware middleware.AuthMiddeleware) {
	tenants := group.Group("/users")
	tenantRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "",
			Handler:     user.RegisterUser,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/status",
			Handler:     user.UpdateUserStatus,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/tenants/:tenant-id/permissions",
			Handler:     user.GetPermissionWithInTenant,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
			},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/:id/roles/:role-id/status",
			Handler:     user.UpdateUserRoleStatus,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
	}
	routing.RegisterRoutes(tenants, tenantRoutes)
}
