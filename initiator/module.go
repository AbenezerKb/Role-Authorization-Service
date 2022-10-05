package initiator

import (
	"2f-authorization/platform/logger"
)

type Module struct {
}

func InitModule(persistence Persistence, log logger.Logger) Module {
	return Module{}
}
