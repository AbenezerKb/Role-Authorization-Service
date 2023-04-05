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
	// tenants := group.Group("/users")
	tenantRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "users",
			Handler:     user.RegisterUser,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPatch,
			Path:        "users/status",
			Handler:     user.UpdateUserStatus,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "users/:id/tenants/:tenant-id/permissions",
			Handler:     user.GetPermissionWithInTenant,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
			},
		},
		{
			Method:      http.MethodPatch,
			Path:        "users/:id/roles/:role-id/status",
			Handler:     user.UpdateUserRoleStatus,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "users/:id/domains/:domain-id/permissions",
			Handler:     user.GetPermissionWithInDomain,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
			},
		},
		{
			Method:      http.MethodPatch,
			Path:        "system/tenants/:tenant_id/users/:user_id/roles/:role_id/status",
			Handler:     user.SystemUpdateUserRoleStatus,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
	}
	routing.RegisterRoutes(group, tenantRoutes)
}
