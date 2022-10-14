package initiator

import (
	"2f-authorization/internal/glue/routing/domain"
	"2f-authorization/internal/glue/routing/permission"
	"2f-authorization/internal/glue/routing/services"
	"2f-authorization/internal/handler/middleware"
	"2f-authorization/platform/logger"
	"2f-authorization/platform/opa"

	"2f-authorization/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(group *gin.RouterGroup, handler Handler, persistence Persistence, log logger.Logger, opa opa.Opa) {
	authmiddleware := middleware.InitAuthMiddleware(log.Named("auth-middleware"), persistence.service, opa)

	docs.SwaggerInfo.BasePath = "/v1"
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	services.InitRoute(group, handler.service, log, authmiddleware)
	domain.InitRoute(group, handler.domain, log, authmiddleware)
	permission.InitRoute(group, handler.permission, log, authmiddleware)
}
