package expressions

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Lex(t *testing.T) {
	//lexer := lexer{}
	expr := ` age  + 20 / 10 / age `
	//expr := ` evaluated == false`
	evalautor, err := New(expr)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := evalautor.Evaluate(EvaluationRequest{
		Variables: map[string]interface{}{
			"age":       1,
			"evaluated": true,
		},
	})
	fmt.Println(err)
	printJSON(res)
	evalautor.Visualise()
	// tokens, err := lexer.Lex(expr)
	// fmt.Println(err)
	// parser := &parser{}
	// res, err := parser.Parse(tokens)
	// fmt.Println(res, err)
	// EulerTraversal(res.Root)
	// r, _ := res.Evaluate()

	// printJSON(r)
}

func printJSON(val interface{}) {
	v, _ := json.Marshal(val)
	fmt.Println(string(v))
}

func Benchmark_TestLexer(b *testing.B) {
	expr := "(requests_made * requests_succeeded / 100) >= 90"
	evalautor, err := New(expr)
	if err != nil {
		fmt.Println(err)
		return
	}
	parameters := map[string]interface{}{
		"requests_made":      99.0,
		"requests_succeeded": 90.0,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evalautor.Evaluate(EvaluationRequest{
			Variables: parameters,
		})
	}

}
