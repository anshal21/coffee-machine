package expressions

// EvaluationRequest is a request object provided for an expression evaluation
// It contains the set of values for variables in the implementation
type EvaluationRequest struct {
	Variables map[string]interface{}
}
