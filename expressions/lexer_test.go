package expressions

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/trace"
	"testing"

	"github.com/anshal21/coffee-machine/lib"
	"github.com/anshal21/coffee-machine/lib/models"
	"github.com/stretchr/testify/assert"
)

func Test_Lex(t *testing.T) {
	//lexer := lexer{}
	//expr := ` age  + 20 / 10 / age `
	expr := `a * a + b * b - 2 * a * b `
	//expr := ` evaluated == false`
	evalautor, err := New(expr)
	if err != nil {
		fmt.Println(err)
		assert.Equal(t, 1, 2)
		return
	}
	f, _ := os.Create("./trace.out")
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	evalautor.Visualise()
	res, err := evalautor.Evaluate(&EvaluationRequest{
		Variables: map[string]interface{}{
			"a": 10,
			"b": 2,
		},
	})

	expectedRes := &EvaluationResponse{
		Type: models.DataTypeNumber,
		Value: models.Value{
			Number: lib.Float64Ptr(64),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expectedRes, res)

}

func printJSON(val interface{}) {
	v, _ := json.Marshal(val)
	fmt.Println(string(v))
}

func Benchmark_TestLexer(b *testing.B) {
	expr := "(requests_made * requests_succeeded / 100) >= 90"
	evalautor, err := New(expr)
	if err != nil {
		//fmt.Println(err)
		return
	}
	parameters := map[string]interface{}{
		"requests_made":      99.0,
		"requests_succeeded": 90.0,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evalautor.Evaluate(&EvaluationRequest{
			Variables: parameters,
		})
	}

}
