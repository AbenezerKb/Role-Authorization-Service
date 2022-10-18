package dto

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type CreateRole struct {
	// Name is the name of the role.
	Name string `json:"name"`
	// TenantName is the name of the tenant the role is in.
	TenantName string `json:"tenant_name"`
	// PermissionID is the list of permissions id's.
	PermissionID []uuid.UUID `json:"permissions_id"`
	// ServiceID is the id of the service the tenant belongs to.
	ServiceID uuid.UUID `json:"service_id"`
}

func (r CreateRole) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("role name is required")),
		validation.Field(&r.TenantName, validation.Required.Error("tenant name is required")),
		validation.Field(&r.PermissionID, validation.By(validatePermissions)),
	)
}
func validatePermissions(id interface{}) error {
	p_id := id.([]uuid.UUID)
	if len(p_id) == 0 {
		return fmt.Errorf("atleast one permission is required")
	}

	for _, id := range p_id {
		if err := id == uuid.Nil; err {
			return fmt.Errorf("permission-id is required")
		}
	}

	return nil
}

type Role struct {
	// ID is the unique identifier for the service.
	// It is automatically generated when the role is created.
	ID uuid.UUID `json:"id,omitempty"`
	// Name is the name of the role.
	Name string `json:"name,omitempty"`
	// Status is the status of the role.
	Status string `json:"status,omitempty"`
	// DeletedAt is the time this service was created.
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	// CreatedAt is the time this service was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt is the time this service was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
