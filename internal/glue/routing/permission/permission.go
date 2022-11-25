package permission

import (
	"2f-authorization/internal/glue/routing"
	"2f-authorization/internal/handler/middleware"
	"2f-authorization/internal/handler/rest"
	"2f-authorization/platform/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, permission rest.Permission, log logger.Logger, authMiddleware middleware.AuthMiddeleware) {
	permissions := group.Group("permissions")
	permissionRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "",
			Handler:     permission.CreatePermission,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "",
			Handler:     permission.ListPermissions,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/inherit",
			Handler:     permission.CreatePermissionDependency,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
			},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/:id",
			Handler:     permission.DeletePermission,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id",
			Handler:     permission.GetPermission,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
	}
	routing.RegisterRoutes(permissions, permissionRoutes)
}
