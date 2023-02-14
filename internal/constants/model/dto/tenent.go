package dto

import (
	"2f-authorization/internal/constants"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type CreateTenent struct {
	//TenantName is the name of the tenant
	TenantName string `json:"tenant_name"`
	//ServiceID  is the service id of service.
	ServiceID uuid.UUID `json:"service_id"`
	// DomainID is the id of the domain the tenant is in.
	DomainID uuid.UUID `json:"domain_id"`
}

func (d CreateTenent) Validate() error {

	return validation.ValidateStruct(&d,
		validation.Field(&d.TenantName, validation.Required.Error("tenant name can not be blank")),
		validation.Field(&d.ServiceID, validation.Required.Error("service id is required"), is.UUID),
		validation.Field(&d.DomainID, validation.NotIn(uuid.Nil.String()).Error("domain id is required")),
	)
}

type GetTenantUsersRequest struct {
	//TenantName is the name of the tenant
	TenantName string `json:"tenant_name"`
	//ServiceID  is the service id of service.
	ServiceID uuid.UUID `json:"service_id"`
}

type UpdateTenantStatus struct {
	// Status is new status of the tenant that is going to replace the old one
	Status string `json:"status"`
}

func (u UpdateTenantStatus) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Status, validation.Required.Error("status is required"), validation.In(constants.Active, constants.InActive).Error("invalid status")),
	)
}

type UserRole struct {
	RoleName string `json:"role_name"`
	ID   string `json:"id"`
	Status   string `json:"status"`
}
type TenantUserRoles struct {
	UserId uuid.UUID  `json:"user_id"`
	Roles  []UserRole `json:"roles"`
}
