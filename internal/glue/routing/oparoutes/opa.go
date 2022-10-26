package oparoutes

import (
	"2f-authorization/internal/glue/routing"
	"2f-authorization/internal/handler/middleware"
	"2f-authorization/internal/handler/rest"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, opa rest.Opa, authMiddleware middleware.AuthMiddeleware) {
	opaGrp := group.Group("authorize")
	opaRoutes := []routing.Router{
		{
			Method:      http.MethodPost,
			Path:        "",
			Handler:     opa.Authorize,
			UnAuthorize: true,
			Middlewares: []gin.HandlerFunc{
				authMiddleware.BasicAuth(),
			},
		},
	}
	routing.RegisterRoutes(opaGrp, opaRoutes)
}
