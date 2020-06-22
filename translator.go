package coffeemachine

import (
	"encoding/json"
	"fmt"
)

type Translator interface {
	Translate(intermediateRepresentation []byte) (*RuleSet, error)
}

type translator struct {
}

func (t *translator) Translate(intermediateRepresentation []byte) (*RuleSet, error) {
	ruleSet := struct {
		ID         string            `json:"id"`
		Predicates map[string]string `json:"predicates"`
		Rules      map[string]struct {
			Predicate string `json:"predicate"`
			PostEvals map[string]struct {
				Expr string `json:"expr"`
				Echo bool   `json:"echo"`
			} `json:"postEvals"`
		} `json:"rules"`
		Relations struct {
			Edges []struct {
				From          string `json:"from"`
				To            string `json:"to"`
				ForwardOutput bool   `json:"forwardOuput"`
			} `json:"edges"`
		} `json:"relations"`
	}{}
	err := json.Unmarshal(intermediateRepresentation, &ruleSet)
	if err != nil {
		return nil, err
	}
	v, _ := json.MarshalIndent(ruleSet, "", "   ")
	fmt.Println(string(v))
	return nil, nil
}
