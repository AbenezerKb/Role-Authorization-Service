package errors

import (
	"github.com/joomcode/errorx"
)

type ErrorType struct {
	ErrorCode int
	ErrorType *errorx.Type
}

var Error = []ErrorType{}
