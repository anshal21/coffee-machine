package coffeemachine

import (
	"io"
)

// RuleEngine is an interface that exposes a Run method to evaluate a RuleSet against
// provided parameters
type RuleEngine interface {
	Run(req *RuleEngineRequest) (*RuleEngineResponse, error)
}

type ruleengine struct {
	evaluator Evaluator
}

// NewRuleEngine is a constructor for RuleEngine
// It accpets an io.Reader to read the rule-set
// and returns an instance of engine which can be run with
// different set of parameter values
func NewRuleEngine(ruleSet io.Reader) (RuleEngine, error) {
	ruleGraph, err := NewParser().Parse(ruleSet)
	if err != nil {
		return nil, err
	}

	return &ruleengine{
		evaluator: NewEvaluator(ruleGraph),
	}, nil

}

// Run method accepts a request containing parameter value and output
// criteria
// It evaluates the rule-set against provided values and returns the response
// as per the output criteria
func (r *ruleengine) Run(req *RuleEngineRequest) (*RuleEngineResponse, error) {
	return r.evaluator.Evaluate(req)
}
