package tests

var _simpleDependencyRuleSet = `{
  "id": "some_ruleset",
  "predicates": {
    "P1": "a > b",
    "P2": "a + b > c",
    "P3": "b > c"
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
    },
    "R2": {
      "predicate": "Predicate:P2",
      "post_evals": [
        {
          "id": "output_1",
          "type": "EXPR",
          "value": "a + b + c"
        }
      ]
    },
    "R3": {
      "predicate": "Predicate:P3",
      "post_evals": [
        {
          "id": "output_1",
          "type": "EXPR",
          "value": "a"
        }
      ]
    }
  },
  "relations": [
    {
      "from": "R1",
      "to": "R2"
    }
  ]
}`
