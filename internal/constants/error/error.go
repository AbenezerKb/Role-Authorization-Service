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
	{
		ErrorCode: http.StatusBadRequest,
		ErrorType: ErrInvalidUserInput,
	},
	{
		ErrorCode: http.StatusBadRequest,
		ErrorType: ErrDataExists,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrReadError,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrWriteError,
	},
}

var (
	opa          = errorx.NewNamespace("opa error")
	db           = errorx.NewNamespace("db error")
	invalidInput = errorx.NewNamespace("validation error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	dbError      = errorx.NewNamespace("db error")
	duplicate    = errorx.NewNamespace("duplicate").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
)

var (
	ErrOpaEvalError         = errorx.NewType(opa, "error evaluating")
	ErrOpaUpdatePolicyError = errorx.NewType(db, "error updating policy data")
	ErrOpaPrepareEvalError  = errorx.NewType(opa, "error preparing for evaluation")
	ErrInvalidUserInput     = errorx.NewType(invalidInput, "invalid user input")
	ErrWriteError           = errorx.NewType(dbError, "could not write to db")
	ErrReadError            = errorx.NewType(dbError, "could not read data from db")
	ErrDataExists           = errorx.NewType(duplicate, "data already exists")
)
