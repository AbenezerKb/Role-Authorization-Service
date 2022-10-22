package dto

import (
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
