package initiator

import (
	"2f-authorization/internal/constants/dbinstance"
	"2f-authorization/internal/storage"
	"2f-authorization/internal/storage/domain"
	"2f-authorization/internal/storage/opa"
	"2f-authorization/internal/storage/permission"
	"2f-authorization/internal/storage/role"
	"2f-authorization/internal/storage/service"
	"2f-authorization/internal/storage/tenant"
	"2f-authorization/internal/storage/user"
	"2f-authorization/platform/logger"
)

type Persistence struct {
	opa        storage.Policy
	service    storage.Service
	doamin     storage.Domain
	permission storage.Permission
	tenant     storage.Tenant
	user       storage.User
	role       storage.Role
}

func InitPersistence(db dbinstance.DBInstance, log logger.Logger) Persistence {
	return Persistence{
		opa:        opa.Init(db, log),
		service:    service.Init(db, log),
		doamin:     domain.Init(db, log.Named("domain-persistant")),
		permission: permission.Init(db, log.Named("permission-persistant")),
		tenant:     tenant.Init(db, log.Named("tenant-persistant")),
		user:       user.Init(db, log.Named("user-persistant")),
		role:       role.Init(db, log.Named("role-persistant")),
	}
}
