package expressions

type Expression interface {
	Evaluate(request EvaluationRequest) (*Response, error)
	Visualise() error
}

type expression struct {
	infix               string
	abstractSyntaxtTree *syntaxTree
}

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

func (e *expression) Evaluate(request EvaluationRequest) (*Response, error) {
	res, err := e.abstractSyntaxtTree.Evaluate(request.Variables)
	if err != nil {
		return nil, err
	}
	return &Response{
		Value: *res.Value,
		Type:  res.Type,
	}, nil
}

func (e *expression) Visualise() error {
	e.abstractSyntaxtTree.Print()
	return nil
}
