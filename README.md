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




How do I create a rule-set?
--



How can I update a rule-set during the run-time?
--

Benchmarks
--