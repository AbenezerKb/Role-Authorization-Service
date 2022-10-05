package initiator

import (
	"2f-authorization/platform/logger"
)

type Handler struct {
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{}
}
