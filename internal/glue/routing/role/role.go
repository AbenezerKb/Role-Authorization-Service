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
	roles := group.Group("/roles")
	roleRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "",
			Handler:     role.CreateRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/:roleid/users/:userid",
			Handler:     role.AssignRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPatch,
			Path:        "/:roleid/users/:userid",
			Handler:     role.RevokeRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodPut,
			Path:        "/:id",
			Handler:     role.UpdateRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/:id",
			Handler:     role.DeleteRole,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
				authMiddleware.Authorize(),
			},
		},
	}
	routing.RegisterRoutes(roles, roleRoutes)
}
