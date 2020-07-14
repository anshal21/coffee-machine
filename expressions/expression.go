package expressions

// Expression is an interface to represent an expression
// It exposes Evaluate method to evaluate an expression
// and a Visualise method to display the execution plan
type Expression interface {
	Evaluate(request *EvaluationRequest) (*EvaluationResponse, error)
	Visualise() error
}

type expression struct {
	infix               string
	abstractSyntaxtTree *syntaxTree
}

// New is a constructor to instantiate a new Expression
// example usage:
// expr, err := New("a > b")
func New(expr string) (Expression, error) {
	lexer := NewLexer()
	tokens, err := lexer.Lex(expr)
	if err != nil {
		return nil, err
	}
	parser := NewParser()
	ast, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return &expression{
		infix:               expr,
		abstractSyntaxtTree: ast,
	}, nil
}

func (e *expression) Evaluate(request *EvaluationRequest) (*EvaluationResponse, error) {
	res, err := e.abstractSyntaxtTree.Evaluate(request.Variables)
	if err != nil {
		return nil, err
	}
	return &EvaluationResponse{
		Value: *res.Value,
		Type:  res.Type,
	}, nil
}

func (e *expression) Visualise() error {
	e.abstractSyntaxtTree.Print()
	return nil
}
