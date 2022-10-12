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
	{
		ErrorCode: http.StatusNotFound,
		ErrorType: ErrNoRecordFound,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrDBDelError,
	},
	{
		ErrorCode: http.StatusUnauthorized,
		ErrorType: ErrAuthError,
	},
	{
		ErrorCode: http.StatusForbidden,
		ErrorType: ErrAcessError,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrInternalServerError,
	},
}

var (
	opa          = errorx.NewNamespace("opa error")
	db           = errorx.NewNamespace("db error")
	invalidInput = errorx.NewNamespace("validation error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	dbError      = errorx.NewNamespace("db error")
	duplicate    = errorx.NewNamespace("duplicate").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	dataNotFound = errorx.NewNamespace("data not found").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	AccessDenied = errorx.RegisterTrait("You are not authorized to perform the action")
	unauthorized = errorx.NewNamespace("unauthorized").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	serverError  = errorx.NewNamespace("server error")
)

var (
	ErrOpaEvalError         = errorx.NewType(opa, "error evaluating")
	ErrOpaUpdatePolicyError = errorx.NewType(db, "error updating policy data")
	ErrOpaPrepareEvalError  = errorx.NewType(opa, "error preparing for evaluation")
	ErrInvalidUserInput     = errorx.NewType(invalidInput, "invalid user input")
	ErrWriteError           = errorx.NewType(dbError, "could not write to db")
	ErrReadError            = errorx.NewType(dbError, "could not read data from db")
	ErrDataExists           = errorx.NewType(duplicate, "data already exists")
	ErrDBDelError           = errorx.NewType(dbError, "could not delete record")
	ErrNoRecordFound        = errorx.NewType(dataNotFound, "no record found")
	ErrAuthError            = errorx.NewType(unauthorized, "you are not authorized.")
	ErrAcessError           = errorx.NewType(unauthorized, "Unauthorized", AccessDenied)
	ErrInternalServerError  = errorx.NewType(serverError, "internal server error")
)
