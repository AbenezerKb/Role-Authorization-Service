package errors

import (
	"net/http"

	"github.com/joomcode/errorx"
)

type ErrorType struct {
	ErrorCode int
	ErrorType *errorx.Type
}

var Error = []ErrorType{
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrOpaEvalError,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrOpaUpdatePolicyError,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrOpaPrepareEvalError,
	},
}

var (
	opa = errorx.NewNamespace("opa error")
	db  = errorx.NewNamespace("db error")
)

var (
	ErrOpaEvalError         = errorx.NewType(opa, "error evaluating")
	ErrOpaUpdatePolicyError = errorx.NewType(db, "error updating policy data")
	ErrOpaPrepareEvalError  = errorx.NewType(opa, "error preparing for evaluation")
)
