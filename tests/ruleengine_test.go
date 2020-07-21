package tests

import (
	"bytes"
	"testing"

	coffeemachine "github.com/anshal21/coffee-machine"
	"github.com/anshal21/coffee-machine/lib"
	"github.com/anshal21/coffee-machine/lib/models"
	"github.com/stretchr/testify/assert"
)

func Test_SampleRule(t *testing.T) {

	tests := []struct {
		name    string
		ruleSet string
		request *coffeemachine.RuleEngineRequest
		res     *coffeemachine.RuleEngineResponse
		err     error
	}{
		{
			name:    "valid rule-set | simple rule-set",
			ruleSet: _simpleRuleSet,
			request: &coffeemachine.RuleEngineRequest{
				Variables: map[string]interface{}{
					"a": 10,
					"b": 8,
				},
			},
			res: &coffeemachine.RuleEngineResponse{
				Outputs: []*coffeemachine.RuleOutput{
					&coffeemachine.RuleOutput{
						ID: "R1",
						PostEvals: []*coffeemachine.EvaluationOutput{
							&coffeemachine.EvaluationOutput{
								ID:   "output_1",
								Type: models.DataTypeNumber,
								Value: models.Value{
									Number: lib.Float64Ptr(18),
								},
							},
							&coffeemachine.EvaluationOutput{
								ID:   "output_2",
								Type: models.DataTypeString,
								Value: models.Value{
									String: lib.StrPtr("action_1"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "valid rule-set | simple dependency rule-set",
			ruleSet: _simpleDependencyRuleSet,
			request: &coffeemachine.RuleEngineRequest{
				Variables: map[string]interface{}{
					"a": 10,
					"b": 8,
					"c": 6,
				},
			},
			res: &coffeemachine.RuleEngineResponse{
				Outputs: []*coffeemachine.RuleOutput{
					&coffeemachine.RuleOutput{
						ID: "R1",
						PostEvals: []*coffeemachine.EvaluationOutput{
							&coffeemachine.EvaluationOutput{
								ID:   "output_1",
								Type: models.DataTypeNumber,
								Value: models.Value{
									Number: lib.Float64Ptr(18),
								},
							},
							&coffeemachine.EvaluationOutput{
								ID:   "output_2",
								Type: models.DataTypeString,
								Value: models.Value{
									String: lib.StrPtr("action_1"),
								},
							},
						},
					},
					&coffeemachine.RuleOutput{
						ID: "R2",
						PostEvals: []*coffeemachine.EvaluationOutput{
							&coffeemachine.EvaluationOutput{
								ID:   "output_1",
								Type: models.DataTypeNumber,
								Value: models.Value{
									Number: lib.Float64Ptr(24),
								},
							},
						},
					},
					&coffeemachine.RuleOutput{
						ID: "R3",
						PostEvals: []*coffeemachine.EvaluationOutput{
							&coffeemachine.EvaluationOutput{
								ID:   "output_1",
								Type: models.DataTypeNumber,
								Value: models.Value{
									Number: lib.Float64Ptr(10),
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "valid rule-set | simple dependency rule-set | partial rules evaluated",
			ruleSet: _simpleDependencyRuleSet,
			request: &coffeemachine.RuleEngineRequest{
				Variables: map[string]interface{}{
					"a": 8,
					"b": 10,
					"c": 6,
				},
			},
			res: &coffeemachine.RuleEngineResponse{
				Outputs: []*coffeemachine.RuleOutput{
					&coffeemachine.RuleOutput{
						ID: "R3",
						PostEvals: []*coffeemachine.EvaluationOutput{
							&coffeemachine.EvaluationOutput{
								ID:   "output_1",
								Type: models.DataTypeNumber,
								Value: models.Value{
									Number: lib.Float64Ptr(8),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine, err := coffeemachine.NewRuleEngine(bytes.NewReader([]byte(test.ruleSet)))
			assert.NoError(t, err)
			res, err := engine.Run(test.request)
			if test.err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.res, res)
			}
		})
	}

}
