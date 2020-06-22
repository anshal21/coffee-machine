package expressions

import "github.com/anshal21/coffee-machine/lib/models"

// EvaluationResponse contains the response of an expression evaluation
// It contains the resultant value of the expression and the DataType associated
// with it
type EvaluationResponse struct {
	Value models.Value
	Type  models.DataType
}
