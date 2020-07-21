coffee-machine
====

What is coffee-machine?
--
coffee-machine is a rule-engine written in GoLang. It enables the separation of business rules from the code. Business rules can be provided to the system as JSON / YAML files and can be evaluated against input parameters in runtime. It comes with a watcher that reflect the any changes to the rule in near realtime, without any need for a restart

what is coffee-machine/expressions
--
expressions is a library to evaluate logical and mathematical expressions provided as strings. It is used by coffee-machine to evaluate the rules, but is not limited to it, the library can directly be used for expression evaluation, without any integration with the rule-engine

```go
  expr := `age  + 20`
  evaluable, _ := New(expr)
  evaluable.Evaluate()
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



How do I create a rule-set?
--



How can I update a rule-set during the run-time?
--

Benchmarks
--