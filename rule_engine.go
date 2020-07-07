package coffeemachine

import "context"

type RuleEngine interface {
	Run(ctx context.Context, req *RuleEngineRequest) (*RuleEngineResponse, error)
}

type ruleengine struct {
	ruleGraph *RuleGraph
}
