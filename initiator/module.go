package initiator

import (
	"2f-authorization/internal/module"
	"2f-authorization/internal/module/domain"
	"2f-authorization/internal/module/opamodule"
	"2f-authorization/internal/module/permission"
	"2f-authorization/internal/module/role"
	"2f-authorization/internal/module/service"
	"2f-authorization/internal/module/tenant"
	"2f-authorization/internal/module/user"

	"2f-authorization/platform/logger"
	opa_platform "2f-authorization/platform/opa"
)

type Module struct {
	service    module.Service
	domain     module.Domain
	permission module.Permission
	tenant     module.Tenant
	user       module.User
	role       module.Role
	opa        module.Opa
}

func InitModule(persistence Persistence, log logger.Logger, opa opa_platform.Opa) Module {
	return Module{
		service:    service.Init(log, persistence.service, opa),
		domain:     domain.Init(log.Named("domain module"), persistence.doamin),
		permission: permission.Init(log.Named("permission module"), persistence.permission, opa),
		tenant:     tenant.Init(log.Named("tenant module"), persistence.tenant, opa),
		user:       user.Init(log.Named("user-module"), persistence.user, persistence.role, opa),
		role:       role.Init(log.Named("role-module"), persistence.role, opa),
		opa:        opamodule.Init(log.Named("opa-module"), opa),
	}
}
