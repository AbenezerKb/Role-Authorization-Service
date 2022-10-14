package initiator

import (
	"2f-authorization/internal/handler/rest"
	"2f-authorization/internal/handler/rest/domain"
	"2f-authorization/internal/handler/rest/permission"
	"2f-authorization/internal/handler/rest/service"
	"2f-authorization/platform/logger"
)

type Handler struct {
	service    rest.Service
	domain     rest.Domain
	permission rest.Permission
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		service:    service.Init(log, module.service),
		domain:     domain.Init(log.Named("domain handler"), module.domain),
		permission: permission.Init(log.Named("permission-handler"), module.permission),
	}
}
