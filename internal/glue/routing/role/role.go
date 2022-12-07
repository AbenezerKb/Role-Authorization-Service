package role

import (
	"2f-authorization/internal/glue/routing"
	"2f-authorization/internal/handler/middleware"
	"2f-authorization/internal/handler/rest"
	"2f-authorization/platform/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, role rest.Role, log logger.Logger, authMiddleware middleware.AuthMiddeleware) {
	// roles := group.Group("/roles")
	roleRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "roles",
			Handler:     role.CreateRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "roles/:id/users/:userid",
			Handler:     role.AssignRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPatch,
			Path:        "roles/:id/users/:userid",
			Handler:     role.RevokeRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPut,
			Path:        "roles/:id",
			Handler:     role.UpdateRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodDelete,
			Path:        "roles/:id",
			Handler:     role.DeleteRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "roles",
			Handler:     role.ListRoles,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPatch,
			Path:        "roles/:id/status",
			Handler:     role.UpdateRoleStatus,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "roles/:id",
			Handler:     role.GetRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "system/users/:id/roles",
			Handler:     role.SystemAssignRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
			},
		},
	}
	routing.RegisterRoutes(group, roleRoutes)
}
