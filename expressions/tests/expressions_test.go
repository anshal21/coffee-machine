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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedRes := getExpectedResponse(test.outputType, test.outputValue)
			evalautor, err := expressions.New(test.expression)
			if test.err == nil {
				assert.NoError(t, err)
				res, err := evalautor.Evaluate(&expressions.EvaluationRequest{
					Variables: test.variables,
				})
				if test.err != nil {
					assert.Error(t, err)
				} else {
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
