package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Service struct {
	ID        uuid.UUID `json:"id"`
	Status    bool      `json:"status"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	DeletedAt time.Time `json:"deleted_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s Service) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required.Error("service name is required"), validation.Length(3, 32).Error("name must be between 3 and 32 characters")),
	)
}
