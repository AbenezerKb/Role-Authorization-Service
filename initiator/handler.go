package initiator

import (
	"2f-authorization/internal/handler/rest"
	"2f-authorization/internal/handler/rest/service"
	"2f-authorization/platform/logger"
)

type Handler struct {
	service rest.Service
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		service: service.Init(log, module.service),
	}
}
