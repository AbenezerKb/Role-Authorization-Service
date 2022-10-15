package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type TenantResponse struct {
	ID         uuid.UUID `json:"id"`
	Status     bool      `json:"status"`
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
	DeletedAt  time.Time `json:"deleted_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateTenent struct {
	TenantName string    `json:"tenant_name"`
	ServiceID  uuid.UUID `json:"service_id"`
}

func (d CreateTenent) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.TenantName, validation.Required.Error("tenant name can not be blank")),
		validation.Field(&d.ServiceID, validation.Required.Error("service id is required"), is.UUID),
	)

}
