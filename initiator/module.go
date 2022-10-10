package initiator

import (
	"2f-authorization/internal/module"
	"2f-authorization/internal/module/service"
	"2f-authorization/platform/logger"
	opa_platform "2f-authorization/platform/opa"
)

type Module struct {
	service module.Service
}

func InitModule(persistence Persistence, log logger.Logger, opa opa_platform.Opa) Module {
	return Module{
		service:  service.Init(log,persistence.service,opa),
	}
}
