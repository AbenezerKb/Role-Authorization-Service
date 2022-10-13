package initiator

import (
	"2f-authorization/internal/constants/dbinstance"
	"2f-authorization/internal/storage"
	"2f-authorization/internal/storage/domain"
	"2f-authorization/internal/storage/opa"
	"2f-authorization/internal/storage/service"
	"2f-authorization/platform/logger"
)

type Persistence struct {
	opa     storage.Policy
	service storage.Service
	doamin  storage.Domain
}

func InitPersistence(db dbinstance.DBInstance, log logger.Logger) Persistence {
	return Persistence{
		opa:     opa.Init(db, log),
		service: service.Init(db, log),
		doamin:  domain.Init(db, log.Named("domain-persistant")),
	}
}
