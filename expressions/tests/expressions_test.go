package tests

import (
	"testing"

	"github.com/anshal21/coffee-machine/expressions"
	"github.com/anshal21/coffee-machine/lib"
	"github.com/anshal21/coffee-machine/lib/models"
	"github.com/stretchr/testify/assert"
)

func Test_Expressions(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		variables   map[string]interface{}
		outputValue interface{}
		udfs        []expressions.UDF
		outputType  models.DataType
		err         error
	}{
		{
			name:       "mathematical | simple addition",
			expression: "a + b + c",
			variables: map[string]interface{}{
				"a": 10,
				"b": 100,
				"c": 125,
			},
			outputValue: float64(235),
			outputType:  models.DataTypeNumber,
		},
		{
			name:       "mathematical | BODMAS",
			expression: "(a + b) * c / d",
			variables: map[string]interface{}{
				"a": 10,
				"b": 100,
				"c": 8,
				"d": 4,
			},
			outputValue: float64(220),
			outputType:  models.DataTypeNumber,
		},
		{
			name:       "user_defined_function",
			expression: "a MY_OP b",
			variables: map[string]interface{}{
				"a": 10,
				"b": 20,
			},
			udfs: []expressions.UDF{
				{
					Token: "MY_OP",
					BinaryOp: func(operandA, operandB, output *expressions.OperationResult) error {
						output.Type = operandA.Type
						output.Value.Number = lib.Float64Ptr(10*(*operandA.Value.Number) + 2*(*operandB.Value.Number))
						return nil
					},
				},
			},
			outputValue: float64(140),
			outputType:  models.DataTypeNumber,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedRes := getExpectedResponse(test.outputType, test.outputValue)
			evalautor, err := expressions.New(test.expression)

			if test.udfs != nil {
				evalautor, err = expressions.NewExpressionsWithUDFs(test.expression, test.udfs...)
				assert.NoError(t, err)
			}

			if test.err == nil {
				assert.NoError(t, err)
				res, err := evalautor.Evaluate(&expressions.EvaluationRequest{
					Variables: test.variables,
				})
				if test.err != nil {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, expectedRes, res)
				}
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func getExpectedResponse(dataType models.DataType, value interface{}) *expressions.EvaluationResponse {
	switch dataType {
	case models.DataTypeNumber:
		return &expressions.EvaluationResponse{
			Type: dataType,
			Value: models.Value{
				Number: lib.Float64Ptr(value.(float64)),
			},
		}
	default:
		return nil
	}
}
