package initiator

import (
	"2f-authorization/internal/glue/routing/services"
	"2f-authorization/platform/logger"

	"github.com/gin-gonic/gin"
)

func InitRouter(group *gin.RouterGroup, handler Handler, module Module, log logger.Logger) {
	services.InitRoute(group, handler.service,log)
}
