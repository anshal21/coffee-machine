package coffeemachine

import (
	"fmt"

	"github.com/anshal21/coffee-machine/expressions"
	"github.com/anshal21/coffee-machine/lib/models"
)


// Evaluator is an interface to evaluate a rule-graph against
// given set of parameter values
// It exposes and Evaluate method that evaluates the rule-graph
type Evaluator interface {
	Evaluate(req *RuleEngineRequest) (*RuleEngineResponse, error)
}

// NewEvaluator is a constructor for Evaluator
// It takes a rule-graph as an input and returns an instance of Evaluator
func NewEvaluator(ruleGraph *RuleGraph) Evaluator {
	return &evaluator{
		ruleGraph: ruleGraph,
	}
}

type evaluator struct {
	ruleGraph *RuleGraph
}

// Evaluate method evaluates the rule-graph against the rule-graph
// It traverses the graph in the dependency order and evaluates the nodes
// for the expressions associated
// It early stops, i.e doesn't evaluate a node if some of it dependecy evaluates to false
func (e *evaluator) Evaluate(req *RuleEngineRequest) (*RuleEngineResponse, error) {
	response := &RuleEngineResponse{}

	outCh := make(chan *RuleOutput, 100)
	// TODO: this can be improved by pre-computing the execution order using topo-sort
	err := e.dfs(req, e.ruleGraph.Root, outCh, &evaluationStats{})
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
		if node.Rule.ID != _rootNodeID {
			postEvals, err := e.evaluatePostEvals(req, node.Rule.PostEvals)
			if err != nil {
				return err
			}
			postEvals.ID = node.Rule.ID
			outCh <- postEvals
		}

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
