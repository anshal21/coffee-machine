package coffeemachine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testRuleSet = `
{
    "id": "some_ruleset",
	"predicates": {
		"P1": "a > b"
	},
	"rules": {
		"R1": {
			"predicate": "Predicate:P1",
			"post_evals": [
				{
					"id": "output_1",
					"type": "EXPR",
					"value": "a + b"
				},
				{
					"id": "output_1",
					"type": "CONST",
					"value": "action_1"
				}
			]
		}
	}   
}
`

func Test_Parse(t *testing.T) {
	p := &parser{}
	rg, err := p.Parse(bytes.NewReader([]byte(testRuleSet)))
	fmt.Println(rg, err)
	eval := &evaluator{}
	res, err := eval.Evaluate(context.Background(), &RuleEngineRequest{
		Variables: map[string]interface{}{
			"a": 10,
			"b": 8,
		},
	}, rg)
	fmt.Println(res, err)
	printJSON(res)
	assert.Equal(t, 1, 2)
}

func printJSON(v interface{}) {
	val, _ := json.Marshal(v)

	fmt.Println(string(val))
}
