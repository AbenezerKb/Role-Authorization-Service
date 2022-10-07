package initiator

import (
	"2f-authorization/platform/logger"
	opa_platform "2f-authorization/platform/opa"
)

type Module struct {
}

func InitModule(persistence Persistence, log logger.Logger, opa opa_platform.Opa) Module {
	return Module{}
}
