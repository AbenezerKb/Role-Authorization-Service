package dto

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type RegisterUser struct {
	// UserId is the id of the user.
	UserId uuid.UUID `json:"user_id"`
	// ServiceID is the id of the service the user belongs to.
	ServiceID uuid.UUID `json:"service_id"`
}

func (d RegisterUser) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.UserId, validation.By(validateUUID)),
		validation.Field(&d.ServiceID, validation.Required.Error("service id is required"), is.UUID),
	)
}

func validateUUID(id interface{}) error {
	if err := id == uuid.Nil; err {
		return fmt.Errorf("user-id is required")
	}

	return nil
}
