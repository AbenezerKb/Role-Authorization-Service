package dto

import (
	"encoding/json"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Permission struct {
	// ID is the unique identifier for the service.
	// It is automatically generated when the permission is registered.
	ID uuid.UUID `json:"id"`
	// Name is the name of the permission being created
	Name string `json:"name,omitempty"`
	// Description is the description of the permission being created
	Description string `json:"description,omitempty"`
	// Statement is an object that holds the action, resource and effect of the permission being created
	Statement `json:"statement,omitempty"`
	// ServiceID is the id of the service the permission belongs to
	ServiceID uuid.UUID `json:"service_id,omitempty"`
	// Domain is an array that holds the id of the domains the permission is accessible at
	Domain []string `json:"domains,omitempty"`
	// Status is the status of the permission.
	Status string `json:"status,omitempty"`
	// DeletedAt is the time this permission was deleted.
	DeletedAt time.Time `json:"-"`
	// CreatedAt is the time this permission was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt is the time this permission was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type CreatePermission struct {
	// Name is the name of the permission being created
	Name string `json:"name"`
	// Description is the description of the permission being created
	Description string `json:"description"`
	// Statement is an object that holds the action, resource and effect of the permission being created
	Statement Statement `json:"statement"`
	// ServiceID is the id of the service the permission belongs to
	ServiceID uuid.UUID `json:"service_id"`
	// Domain is an array that holds the id of the domains the permission is accessible at
	Domain []string `json:"domains"`
}

type Statement struct {
	// Effect is the effect that's taken on the permission
	// It is either allow or deny
	Effect string `json:"effect"`
	// Resource is the urn for the path that is being accessed
	Resource string `json:"resource"`
	// Action is the urn for the action(method) the user is taking on the resource
	Action string `json:"action"`
}

func (a Statement) Value() ([]byte, error) {
	return json.Marshal(a)
}

func (c CreatePermission) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required.Error("permission name is required")),
		validation.Field(&c.Description, validation.Required.Error("permission description is required")),
		validation.Field(&c.Statement),
	)
}

func (s Statement) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Resource, validation.Required.Error("statement resource is required")),
		validation.Field(&s.Effect, validation.Required.Error("statement effect is required")),
		validation.Field(&s.Action, validation.Required.Error("statement action is required")),
	)
}

type GetAllPermissionsReq struct {
	//TenantName is the name of the tenant
	TenantName string `json:"tenant_name"`
	//ServiceID  is the service id of service.
	ServiceID uuid.UUID `json:"service_id"`
}

func (g GetAllPermissionsReq) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.TenantName, validation.Required.Error("tenant name can not be blank")),
		validation.Field(&g.ServiceID, validation.NotIn(uuid.Nil.String()).Error("service id is required")))
}
