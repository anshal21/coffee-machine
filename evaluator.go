package coffeemachine

import (
	"context"
	"fmt"

	"github.com/anshal21/coffee-machine/expressions"
	"github.com/anshal21/coffee-machine/lib/models"
)

type Evaluator interface {
	Evaluate(ctx context.Context, req *RuleEngineRequest, graph *RuleGraph) (*RuleEngineResponse, error)
}

type evaluator struct {
}

func (e *evaluator) Evaluate(ctx context.Context, req *RuleEngineRequest, graph *RuleGraph) (*RuleEngineResponse, error) {
	response := &RuleEngineResponse{}

	outCh := make(chan *RuleOutput, 100)
	err := e.dfs(req, graph.Root, outCh, &evaluationStats{})
	if err != nil {
		return nil, err
	}
	close(outCh)

	for ruleOutput := range outCh {
		response.Outputs = append(response.Outputs, ruleOutput)
	}

	return response, nil
}

type evaluationStats struct {
	evaluated      int
	evaluatedTrue  int
	evaluatedRules []string
}

func (e *evaluator) dfs(req *RuleEngineRequest, node *Node, outCh chan<- *RuleOutput, stats *evaluationStats) error {
	fmt.Println("Rule ID: ", node.Rule.ID)
	printJSON(node)
	res, err := node.Rule.Predicate.Evaluate(&expressions.EvaluationRequest{
		Variables: req.Variables,
	})

	if err != nil {
		return err
	}

	if res.Type != models.DataTypeBool {
		return fmt.Errorf("rule %v, does not have a boolean expression", node.Rule.ID)
	}

	if *res.Value.Bool {
		postEvals, err := e.evaluatePostEvals(req, node.Rule.PostEvals)
		if err != nil {
			return err
		}
		postEvals.ID = node.Rule.ID
		outCh <- postEvals

		for _, edge := range node.Relations {
			err = e.dfs(req, edge.Destination, outCh, stats)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *evaluator) evaluatePostEvals(req *RuleEngineRequest, postEvals []*RulePostEval) (*RuleOutput, error) {

	ruleOutput := &RuleOutput{
		PostEvals: make([]*EvaluationOutput, 0, len(postEvals)),
	}

	for _, postEval := range postEvals {
		switch postEval.Type {
		case OutputTypeExpression:
			res, err := postEval.Evaluable.Evaluate(&expressions.EvaluationRequest{
				Variables: req.Variables,
			})
			if err != nil {
				return nil, err
			}

			ruleOutput.PostEvals = append(ruleOutput.PostEvals, &EvaluationOutput{
				ID:    postEval.ID,
				Value: res.Value,
				Type:  res.Type,
			})

		case OutputTypeConstant:
			ruleOutput.PostEvals = append(ruleOutput.PostEvals, &EvaluationOutput{
				ID: postEval.ID,
				Value: models.Value{
					String: &postEval.Const,
				},
				Type: models.DataTypeString,
			})
		default:
			return nil, fmt.Errorf("invalid output type %v in the rule", postEval.Type)
		}
	}
	return ruleOutput, nil
}
