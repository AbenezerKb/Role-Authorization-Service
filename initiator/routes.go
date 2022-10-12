package initiator

import (
	"2f-authorization/internal/glue/routing/services"
	"2f-authorization/internal/handler/middleware"
	"2f-authorization/platform/logger"
	"2f-authorization/platform/opa"

	"github.com/gin-gonic/gin"
)

func InitRouter(group *gin.RouterGroup, handler Handler, persistence Persistence, log logger.Logger, opa opa.Opa) {
	authmiddleware := middleware.InitAuthMiddleware(log.Named("auth-middleware"), persistence.service,opa)
	services.InitRoute(group, handler.service, log, authmiddleware)
}
