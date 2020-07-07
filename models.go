package coffeemachine

import (
	"github.com/anshal21/coffee-machine/expressions"
	"github.com/anshal21/coffee-machine/models"
)

const (
	OutputTypeExpression = "EXPR"
	OutputTypeConstant   = "CONST"
)

type Rule struct {
	ID        string
	Predicate expressions.Expression
	PostEvals []*RuleOutput
}

type RuleOutput struct {
	ID        string
	Type      string
	Const     string
	Evaluable expressions.Expression
	Echo      bool
}

type Node struct {
	Rule      *Rule
	Relations []*Edge
}

type Edge struct {
	Destination   *Node
	ForwardOutput bool
}

type RuleGraph struct {
	ID        string
	Root      *Node
	Constants []interface{}
}

type RuleEngineRequest struct {
	Variables          map[string]interface{}
	EvaluatedCount     bool
	EvaluatedTrueCount bool
	EvaluatedRules     bool
}

type EvaluationOutput struct {
	ID    string
	Value models.Value
}

type RuleEngineResponse struct {
	RulesEvaluated     int
	RulesEvaluatedTrue int
	Outputs            map[string]EvaluationOutput
	EvaluatedRules     []string
}
