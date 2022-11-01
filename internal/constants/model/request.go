package model

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Request struct {
	// Subject is the user id who is trying to take action on the resource.
	Subject string `json:"subject"`
	// Resource is the urn for the resource the user is trying to take action on.
	Resource string `json:"resource"`
	// Tenant is the scope the user is operating.
	// It is set to "administrator" if it is not provided.
	Tenant string `json:"tenant"`
	// Service is the id of the service.
	// It is set by the server after authenticating the service.
	Service string `json:"service"`
	// Action  is the urn of the action the user is taking on the resource.
	Action string `json:"action"`
	// Fields are the attributes of the entity.
	Fields []string `json:"fields"`
}

func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Subject, validation.Required.Error("subject is required")),
		validation.Field(&r.Action, validation.Required.Error("action is required")),
		validation.Field(&r.Resource, validation.Required.Error("resource is required")),
		validation.Field(&r.Tenant, validation.Required.Error("tenant is required")),
	)
}
