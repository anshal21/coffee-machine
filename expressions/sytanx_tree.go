package expressions

import (
	"fmt"
	"sync"

	"github.com/anshal21coffee-machine/lib"
	"github.com/anshal21coffee-machine/lib/errors"
)

type node struct {
	Token    *Token
	Children []*node
}

type syntaxTree struct {
	Root *node
}

type evaluationResult struct {
	Type  DataType
	Value *Value
}

type Value struct {
	Number *float64
	String *string
	Bool   *bool
}

const (
	_treeLevelPrefix = "----"
)

func (t *syntaxTree) Evaluate(values map[string]interface{}) (*evaluationResult, error) {
	return t.evaluteHelper(t.Root, values)
}

func (t *syntaxTree) Print() {
	fmt.Printf("_\n")
	inorderTraversal(t.Root, _treeLevelPrefix, 1)
}

func inorderTraversal(node *node, prefix string, level int) {
	if node == nil {
		return
	}
	nextPrefix := prefix
	for i := 0; i < level; i++ {
		nextPrefix = nextPrefix + _treeLevelPrefix
	}
	if node.Token.Type == Operator {
		inorderTraversal(node.Children[0], nextPrefix, level+1)
	}
	fmt.Printf("|\n|%v> %v [%v]\n", prefix, node.Token.Value, node.Token.Type)
	if node.Token.Type == Operator {
		inorderTraversal(node.Children[1], nextPrefix, level+1)
	}
}

func (t *syntaxTree) evaluteHelper(curr *node, values map[string]interface{}) (*evaluationResult, error) {
	switch curr.Token.Type {
	case String:
		return &evaluationResult{
			Type: DataTypeString,
			Value: &Value{
				String: lib.StrPtr(curr.Token.Value.(string)),
			},
		}, nil
	case Number:
		return &evaluationResult{
			Type: DataTypeNumber,
			Value: &Value{
				Number: lib.Float64Ptr(curr.Token.Value.(float64)),
			},
		}, nil
	case Bool:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(curr.Token.Value.(bool)),
			},
		}, nil
	case Variable:
		return resolveVariableValue(curr.Token, values)
	case Operator:
		if len(curr.Children) != 2 {
			return nil, fmt.Errorf("Not enough operand present for operator %v", curr.Token.Value)
		}
		wg := sync.WaitGroup{}
		wg.Add(2)
		var err1, err2 error
		var res1, res2 *evaluationResult
		func() {
			res1, err1 = t.evaluteHelper(curr.Children[0], values)
			wg.Done()
		}()
		func() {
			res2, err2 = t.evaluteHelper(curr.Children[1], values)
			wg.Done()
		}()
		wg.Wait()
		if err1 != nil {
			return nil, err1
		}
		if err2 != nil {
			return nil, err2
		}

		return applyOperator(res1, res2, curr.Token)
	}
	return nil, fmt.Errorf("unsupported token type %v", curr.Token.Type)
}

func resolveVariableValue(token *Token, values map[string]interface{}) (*evaluationResult, error) {
	val, ok := values[token.Value.(string)]
	if !ok {
		return nil, fmt.Errorf("error value not provided for variable %v", token.Value)
	}
	// TODO: Add nested feature resolution here
	switch val.(type) {
	case string:
		return &evaluationResult{
			Type: DataTypeString,
			Value: &Value{
				String: lib.StrPtr(val.(string)),
			},
		}, nil
	case float64:
		return &evaluationResult{
			Type: DataTypeNumber,
			Value: &Value{
				Number: lib.Float64Ptr(val.(float64)),
			},
		}, nil
	case int:
		return &evaluationResult{
			Type: DataTypeNumber,
			Value: &Value{
				Number: lib.Float64Ptr(float64(val.(int))),
			},
		}, nil
	case bool:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(val.(bool)),
			},
		}, nil
	}
	return nil, fmt.Errorf("invalid variable type %v", val)
}

func applyOperator(res1 *evaluationResult, res2 *evaluationResult, operation *Token) (res *evaluationResult, err error) {
	defer func() {
		if err != nil {
			if e := err.(*errors.Error); e.Code == ErrIncompatibleOperation {
				err = errors.New(ErrIncompatibleOperation, fmt.Errorf("%v at position %v", e.Msg, operation.Index))
			}
		}
	}()

	if res1.Type != res2.Type {
		return nil, errors.New(ErrIncompatibleOperation, fmt.Errorf("cannot apply '%v' operation on type '%v' and '%v' at position '%v'", operation.Value, res1.Type, res2.Type, operation.Index))
	}

	switch operation.Value.(string) {
	case "+":
		return add(res1, res2)
	case "-":
		return sub(res1, res2)
	case "*":
		return mul(res1, res2)
	case "/":
		return div(res1, res2)
	case "<":
		return lt(res1, res2)
	case ">":
		return gt(res1, res2)
	case "<=":
		return lte(res1, res2)
	case ">=":
		return gte(res1, res2)
	case "==":
		return equal(res1, res2)
	}
	return nil, errors.New(ErrUnsupportedOperation, fmt.Errorf("unsupported operator %v at position %v", operation.Value, operation.Index))
}

func add(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case DataTypeString:
		return &evaluationResult{
			Type: DataTypeString,
			Value: &Value{
				String: lib.StrPtr(fmt.Sprintf("%v%v", *res1.Value.String, *res2.Value.String)),
			},
		}, nil
	case DataTypeNumber:
		return &evaluationResult{
			Type: DataTypeNumber,
			Value: &Value{
				Number: lib.Float64Ptr(*res1.Value.Number + *res2.Value.Number),
			},
		}, nil
	}
	return nil, incompatibleOperationError("+", res1.Type)
}

func sub(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case DataTypeNumber:
		return &evaluationResult{
			Type: DataTypeNumber,
			Value: &Value{
				Number: lib.Float64Ptr(*res1.Value.Number - *res2.Value.Number),
			},
		}, nil
	}
	return nil, incompatibleOperationError("-", res1.Type)
}

func mul(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case DataTypeNumber:
		return &evaluationResult{
			Type: DataTypeNumber,
			Value: &Value{
				Number: lib.Float64Ptr(*res1.Value.Number * *res2.Value.Number),
			},
		}, nil
	}
	return nil, incompatibleOperationError("*", res1.Type)
}

func div(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case DataTypeNumber:
		if *res2.Value.Number == 0 {
			return nil, fmt.Errorf("encounter 0 value as denominatior")
		}
		return &evaluationResult{
			Type: DataTypeNumber,
			Value: &Value{
				Number: lib.Float64Ptr(*res1.Value.Number / *res2.Value.Number),
			},
		}, nil
	}
	return nil, incompatibleOperationError("/", res1.Type)
}

func lt(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case DataTypeString:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(*res1.Value.String < *res2.Value.String),
			},
		}, nil
	case DataTypeNumber:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(*res1.Value.Number < *res2.Value.Number),
			},
		}, nil
	}
	return nil, incompatibleOperationError("<", res1.Type)
}

func gt(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	res, err := lt(res2, res1)
	if err != nil {
		if err.(*errors.Error).Code == ErrIncompatibleOperation {
			err = incompatibleOperationError(">", res1.Type)
		}
		return nil, err
	}
	return res, nil
}

func lte(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case DataTypeString:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(*res1.Value.String <= *res2.Value.String),
			},
		}, nil
	case DataTypeNumber:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(*res1.Value.Number <= *res2.Value.Number),
			},
		}, nil
	}
	return nil, incompatibleOperationError("<=", res1.Type)
}

func gte(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	res, err := lte(res2, res1)
	if err != nil {
		if err.(*errors.Error).Code == ErrIncompatibleOperation {
			err = incompatibleOperationError(">=", res1.Type)
		}
		return nil, err
	}
	return res, nil
}

func equal(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case DataTypeString:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(*res1.Value.String == *res2.Value.String),
			},
		}, nil
	case DataTypeNumber:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(*res1.Value.Number == *res2.Value.Number),
			},
		}, nil
	case DataTypeBool:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(*res1.Value.Bool == *res2.Value.Bool),
			},
		}, nil
	}
	return nil, incompatibleOperationError("==", res1.Type)
}

func or(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case DataTypeBool:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(*res1.Value.Bool || *res2.Value.Bool),
			},
		}, nil
	}
	return nil, incompatibleOperationError("||", res1.Type)
}

func and(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case DataTypeBool:
		return &evaluationResult{
			Type: DataTypeBool,
			Value: &Value{
				Bool: lib.BoolPtr(*res1.Value.Bool && *res2.Value.Bool),
			},
		}, nil
	}
	return nil, incompatibleOperationError("&&", res1.Type)
}

func incompatibleOperationError(op string, operandType DataType) *errors.Error {
	return errors.New(ErrIncompatibleOperation, fmt.Errorf("operation '%v' is not compatible with '%v' type", op, operandType))
}
