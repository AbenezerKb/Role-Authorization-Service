package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/google/uuid"
)

type CreateDomain struct {
	//Name is the name of the domain
	Name string `json:"name"`
	// ServiceID is the id of the service which own the domain.
	ServiceID uuid.UUID `json:"service_id"`
}

func (d CreateDomain) Validated() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required.Error("domain name can not be blank")),
		validation.Field(&d.ServiceID, validation.Required.Error("service id is required"), is.UUID),
	)

}

type Domain struct {
	// ID is the unique identifier for the domain
	// It is automatically generated when the domain is created.
	ID uuid.UUID `json:"id"`
	//Name is the name of the domain
	Name string `json:"name"`
	// DeletedAt is the time this domain was deleted.
	DeletedAt time.Time `json:"deleted_at"`
	// ServiceID is the id of the service which own the domain.
	ServiceID uuid.UUID `json:"service_id"`
	// CreatedAt is the  time this domain created at.
	CreatedAt time.Time `json:"created_at"`
	// CreatedAt is the  time this domain updated at.
	UpdatedAt time.Time `json:"updated_at"`
}

func (d Domain) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required.Error("domain name can not be blank")),
		validation.Field(&d.ServiceID, validation.Required.Error("service id is required"), is.UUID),
	)

}
