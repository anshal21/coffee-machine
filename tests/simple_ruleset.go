package tests

var _simpleRuleSet = `{
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
					"id": "output_2",
					"type": "CONST",
					"value": "action_1"
				}
			]
		}
	}
}
`
