package coffeemachine

import (
	"github.com/anshal21/coffee-machine/expressions"
	"github.com/anshal21/coffee-machine/lib/models"
)

const (
	OutputTypeExpression = "EXPR"
	OutputTypeConstant   = "CONST"
)

// Rule is an in-memory represent of a rule
// defined in the rule-set
// It consists of rule-id, associated predicate
// ans a list of post-evals
type Rule struct {
	ID        string
	Predicate expressions.Expression
	PostEvals []*RulePostEval
}

// RulePostEval is an expression or constant
// that is calculated if the associated rule
// evaluates to true
// The calculated value is returned in the
// response
type RulePostEval struct {
	ID        string
	Type      string
	Const     string
	Evaluable expressions.Expression
	Echo      bool
}

// Node represents a node in the rule-graph
// it consists of a rule and any edges / relations
// with other rules
type Node struct {
	Rule      *Rule
	Relations []*Edge
}

// Edge is a struct to represent a relation
// between two rules in the dependency graph
type Edge struct {
	Destination   *Node
	ForwardOutput bool
}

// RuleGraph is a dependency graph representation
// of the rule-set
type RuleGraph struct {
	ID        string
	Root      *Node
	Constants []interface{}
}

// RuleEngineRequest is a struct that holds
// the input parameters required for the
// evaluation and the response criteria
type RuleEngineRequest struct {
	// map[string]interface{} contains the values for the variables
	// in the predicates defined in the rule-engine
	Variables map[string]interface{}
	// EvaluatedCount, If true, response contains the count of rules that were
	// evaluated
	EvaluatedCount bool
	// EvaluatedTrueCount, If true, response contais the count of rules that
	// evaluated to true
	EvaluatedTrueCount bool
	// EvaluatedRules, If true, output contains the rule-ids of the evaluated rules
	EvaluatedRules bool
}

// EvaluationOutput contains evaluation output for a expression
// It is used to hold the values of the PostEvals calculations
type EvaluationOutput struct {
	ID    string
	Value models.Value
	Type  models.DataType
}

// RuleOutput is a struct to hold the result of a rule evaluation
// It contains the RuleID and evaluation response for all
// the associated post-evals
type RuleOutput struct {
	ID        string
	PostEvals []*EvaluationOutput
}

// RuleEngineResponse is a struct that holds the response for a
// rule-engine Run
type RuleEngineResponse struct {
	RulesEvaluated     int
	RulesEvaluatedTrue int
	Outputs            []*RuleOutput
	EvaluatedRules     []string
}
