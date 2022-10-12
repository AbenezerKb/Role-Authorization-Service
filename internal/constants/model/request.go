package model

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Request struct {
	Subject  string `json:"subject"`
	Resource string `json:"resource"`
	Tenant   string `json:"tenant"`
	Service  string `json:"service"`
	Action   string `json:"action"`
}

func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Subject, validation.Required.Error("subject is required")),
		validation.Field(&r.Action, validation.Required.Error("action is required")),
		validation.Field(&r.Resource, validation.Required.Error("resource is required")),
	)
}
