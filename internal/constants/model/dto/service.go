package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Service struct {
	// ID is the unique identifier for the service.
	// It is automatically generated when the service is created.
	ID uuid.UUID `json:"id"`
	// Status is the current status of the service.
	// It is set to false by default.
	Status bool `json:"status"`
	// Name is the name of the service.
	Name string `json:"name"`
	// Password is the secret the service uses to authenticate itself.
	// It is automatically generated when the service is created.
	Password string `json:"password"`
	// CreatedAt is the time this service was deleted.
	DeletedAt time.Time `json:"deleted_at"`
	// CreatedAt is the time this service was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time this service was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateService struct {
	// Name is the name of the service.
	Name string `json:"name"`
	// UserId is the id of the user assigned as the super admin for the created service.
	UserId string `json:"user_id"`
	// Password is the secret the service uses to authenticate itself.
	// It is automatically generated when the service is created.
	Password string `json:"password"`
}
type CreateServiceResponse struct {
	// ServiceID is the unique identifier for the created service.
	ServiceID uuid.UUID `json:"service_id"`
	// Password is the secret the service uses to authenticate itself.
	// It is automatically generated when the service is created.
	Password string `json:"password"`
	// Service is the name of the service.
	Service string `json:"service"`
	// ServiceStatus is the status of the service.
	// It is set to false when the service is created.
	ServiceStatus bool `json:"service_status"`
	// Tenant is the domain the super admin is in.
	// It is automatically created upon the creation of the service.
	Tenant string `json:"tenant"`
}

func (s CreateService) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required.Error("service name is required"), validation.Length(3, 32).Error("name must be between 3 and 32 characters")),
		validation.Field(&s.UserId, validation.Required.Error("user id is required")),
	)
}
