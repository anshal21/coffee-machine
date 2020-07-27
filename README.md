[![Build Status](https://travis-ci.org/anshal21/coffee-machine.svg?branch=master)](https://travis-ci.org/anshal21/coffee-machine)

coffee-machine
====

What is coffee-machine?
--
coffee-machine is a rule-engine written in GoLang. It enables the separation of business rules from the code. Business rules can be provided to the system as JSON / YAML files and can be evaluated against input parameters in runtime. It comes with a watcher that reflect the any changes to the rule in near realtime, without any need for a restart

What is coffee-machine/expressions
--
expressions is a library to evaluate logical and mathematical expressions provided as strings. It is used by coffee-machine to evaluate the rules, but is not limited to it, the library can directly be used for expression evaluation, without any integration with the rule-engine

```go
  expr := `a * a + b * b - 2 * a * b `
  evaluable, _ := New(expr)
  evaluable.Evaluate(&EvaluationRequest{
		Variables: map[string]interface{}{
			"a": 10,
			"b": 2,
		},
	})
```
```
Output:
{
  "Value": {
    "Number": 64,
    "String": null,
    "Bool": null
  },
  "Type": 2
}
```

```go
  evaluable.Visualize()
```
```
Output:
_
|
|----------------------------> a [Variable]
|
|----------------> * [Operator]
|
|----------------------------> a [Variable]
|
|--------> + [Operator]
|
|----------------------------> b [Variable]
|
|----------------> * [Operator]
|
|----------------------------> b [Variable]
|
|----> - [Operator]
|
|----------------------------> 2 [Number]
|
|----------------> * [Operator]
|
|----------------------------> a [Variable]
|
|--------> * [Operator]
|
|----------------> b [Variable]
```


When do I need a rule-engine?
--
A rule-engine can be leveraged by the systems, that
- Depend on business rules/logic that require changes very often
- Deal with huge number of business rules
- Contain business-rules that deal with large number or conditional branching making
 them hard to change and maintain in the code
- Have a requirement for non-developers to be able to configure and change rule-sets

Concepts
--
#### predicates
It's a list of predicates. A predicate is a logical expression, It is used by rules to represent the associated condition

#### rules
It is a list of different business rules. A rule is made of two components
###### predicate
predicate holds the business condition for the rule
###### post_evals
post_evals is a set of output that rule is supposed to return if the associated condition evaluates to true
A post_evals can either be a CONST ( constant string ) or an EXPR ( logical or mathematical expression ) in itself


#### realtions
It is a list of dependency edges between different rules. Each has some property associated, currently one such provided property is
`forward_output` which enables the system to feed the response from one rule as an input variables to the other rules


How do I create a rule-set?
--

An example rule-set
```json
{
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
}
```


checkout tests package for more


How can I update a rule-set during the run-time?
--

Benchmarks
--