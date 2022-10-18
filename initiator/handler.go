package initiator

import (
	"2f-authorization/internal/handler/rest"
	"2f-authorization/internal/handler/rest/domain"
	"2f-authorization/internal/handler/rest/permission"
	"2f-authorization/internal/handler/rest/role"
	"2f-authorization/internal/handler/rest/service"
	"2f-authorization/internal/handler/rest/tenant"
	"2f-authorization/internal/handler/rest/user"
	"2f-authorization/platform/logger"
)

type Handler struct {
	service    rest.Service
	domain     rest.Domain
	permission rest.Permission
	tenant     rest.Tenant
	user       rest.User
	role       rest.Role
}

func InitHandler(module Module, log logger.Logger) Handler {
	return Handler{
		service:    service.Init(log, module.service),
		domain:     domain.Init(log.Named("domain handler"), module.domain),
		permission: permission.Init(log.Named("permission-handler"), module.permission),
		tenant:     tenant.Init(log.Named("tenant-handler"), module.tenant),
		user:       user.Init(log.Named("user-handler"), module.user),
		role:       role.Init(log.Named("role-handler"), module.role),
	}
}
