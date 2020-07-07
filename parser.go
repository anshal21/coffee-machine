package coffeemachine

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"io/ioutil"

	"github.com/anshal21/coffee-machine/lib/errors"
)

const (
	_rootNodeID = "root"
	_rootNodePredicate = "1 == 1"
)

type Parser interface {
	Parse(reader io.Reader) (*RuleGraph, error)
}

type parser struct{}

func (p *parser) Parse(reader io.Reader) (*RuleGraph, error) {
	data := struct {
		ID         string            `json:"id"`
		Predicates map[string]string `json:"predicates`
		Rules      map[string]struct {
			Predicate string `json:"predicate"`
			PostEvals []map[string]struct {
				Type  string `json:"type"`
				Value string `json:"value"`
				Echo  bool   `json:"echo"`
			}
		}
		Relations []struct {
			From          string `json:"from"`
			To            string `json:"to"`
			ForwardOutput string `json:"forward_output"`
		}
	}{}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	rulesIDToNode := make(map[string]*Node)
	for ruleID, ruleDef := range data.Rules {
		if _, ok := rulesIDToStruct[ruleID]; ok {
			return nil, errors.New(ErrInvalidRuleSet, fmt.Errorf("rule id %v has been used already", ruleID))
		}
		postEvals := make([]*RuleOutput, 0)
		for key, postEval := range ruleDef.PostEvals {
			
			output := &RuleOutput{
				ID: key,
				Type: postEval.Type,
				Echo: postEval.Echo,
			}
			
			switch postEval.Type {
			case OutputTypeExpression:
				predicate, err := resolvePredicate(data.Predicates, postEval.Value)
				if err != nil {
					return nil, errors.New(ErrInvalidRuleSet, err)
				}
				expr, err := expressions.New(predicate)
				if err != nil {
					return nil, errors.New(ErrInvalidRuleSet, fmt.Errorf("rule %v has invalid predicate for output %v, %v", ruleID, key, err.Error()))
				}
				output.Evaluable = expr	
		
			case OutputTypeConstant:
				output.Const = postEval.Value
				
			default:
				return nil, errors.New(ErrInvalidRuleSet, fmt.Errorf("invalid output type used for output %v in rule %v", key, ruleID))
			}
			postEvals = append(postEvals, output)
		}

		predicate, err := resolvePredicate(data.Predicates, ruleDef.Predicate)
		if err != nil {
			return nil, errors.New(ErrInvalidRuleSet, err)
		}
		expr, err := expressions.New(predicate)
		if err != nil {
			return nil, errors.New(ErrInvalidRuleSet, fmt.Errorf("rule %v has invalid predicate, %v", ruleID, err.Error()))
		}

		rule := &Rule{
			ID: ruleID,
			Predicate: expr,
			PostEvals: postEvals,
		}

		rulesIDToStruct[ruleID] = &Node{
			Rule: rule,
		}
	}

	rootNode = &Node{
		Rule: &Rule{
			ID: _rootNodeID,
			Predicate: _rootNodePredicate,
		},
	}

	indegree := make(map[*Node]int)

	for _, relation := range data.Relations {
		if _, ok := rulesIDToNode[relation.From]; !ok {
			return errors.New(ErrInvalidRuleSet, fmt.Errorf("invalid rule id %v used for relation", relation.From))
		}
		if _, ok := rulesIDToNode[relation.To]; !ok {
			return errors.New(ErrInvalidRuleSet, fmt.Errorf("invalid rule id %v used for relation", relation.To))
		}

		fromNode := rulesIDToNode[relation.From]
		toNode := rulesIDToNode[relation.To]
		fromNode.Relations = append(fromNode.Relations, &Edge{
			Destination: toNode,
			ForwardOutput: relation.ForwardOutput,
		})

		indegree[toNode] = indegree[toNode] + 1
	}
	
	for key, val := range indegree {
		if val == 0 {
			root.Relations = append(root.Relations, &Edge{
				Destination: key,
			})
		}
	}

	return &RuleGraph{
		ID: data.ID,
		Root: root,
	}, nil
}


func resolvePredicate(predicates map[string]string, expression string) (string, error) {

	isPredicate := func(s string) {
		pattern := "Predicate:"
		
		if len(s) <= len(pattern) {
			return false
		} 

		for index := range pattern {
			if pattern[index] != s[index] {
				return false
			}
		}
		return true
	}

	tokens := strings.Split(expression, " ")
	for index := range tokens {
		if isPredicate(tokens[index]) {
			predicateID := strings.Split(tokens[index], ":")[1]
			if _, ok := predicates[predicateID] {
				return fmt.Errorf("reference to invalid predicate %v", predicateID)
			}
			tokens[index] = predicates[predicateID]
		}
	}

	return string.Join(tokens, " "), nil
}
