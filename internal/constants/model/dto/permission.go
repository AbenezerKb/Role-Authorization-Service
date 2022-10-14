package dto

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type CreatePermission struct {
	// Name is the name of the permission being created
	Name string `json:"name"`
	// Description is the description of the permission being created
	Description string `json:"description"`
	// Statement is an object that holds the action, resource and effect of the permission being created
	Statement `json:"statement"`
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
	Action   string `json:"action"`
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
