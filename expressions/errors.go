package expressions

import (
	"github.com/anshal21/coffee-machine/lib/errors"
)

// Set of error codes thrown by the expression evaluation
const (
	ErrInvalidExpression     errors.ErrCode = "InvalidExpression"
	ErrMissingVariableValue  errors.ErrCode = "MissingVariableValue"
	ErrIncompatibleOperation errors.ErrCode = "IncompatibleOperation"
	ErrUnsupportedOperation  errors.ErrCode = "UnsupportedOperation"
)
