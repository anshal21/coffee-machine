package coffeemachine

import "context"

type Evaluator interface {
	Evaluate(ctx context.Context, graph *RuleGraph) (*RuleEngineRunResponse, error)
}

type evaluator struct {
}

func (e *Evaluator) Evaluate(ctx context.Context, graph *RuleGraph) (*RuleEngineRunResponse, error) {

}
