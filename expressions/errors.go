package expressions

import (
	"github.com/anshal21coffee-machine/lib/errors"
)

const (
	ErrInvalidExpression     errors.ErrCode = "InvalidExpression"
	ErrMissingVariableValue  errors.ErrCode = "MissingVariableValue"
	ErrIncompatibleOperation errors.ErrCode = "IncompatibleOperation"
	ErrUnsupportedOperation  errors.ErrCode = "UnsupportedOperation"
)
