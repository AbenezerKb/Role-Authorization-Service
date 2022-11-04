package dto

import (
	"2f-authorization/internal/constants"
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

type UpdateUserStatus struct {
	// Status is new status that will replace old status of the user
	Status string `json:"status"`
	// ServiceID is the unique identifier for the service.
	ServiceID uuid.UUID `json:"service"`
	// UserId is the id of the user.
	UserID uuid.UUID `json:"user"`
}

func (u UpdateUserStatus) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Status, validation.Required.Error("status is required"), validation.In(constants.Active, constants.InActive).Error("invalid status")),
		validation.Field(&u.ServiceID, validation.NotIn(uuid.Nil.String()).Error("service id is required")),
		validation.Field(&u.UserID, validation.NotIn(uuid.Nil.String()).Error("user id is required")),
	)
}
