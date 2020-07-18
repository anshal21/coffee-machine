package expressions

import (
	"fmt"
	"sync"

	"github.com/anshal21/coffee-machine/lib/errors"
	"github.com/anshal21/coffee-machine/lib/models"
)

// Node represents a node of a syntax tree
// Node can either be an Operand or Operator node
type node struct {
	Token      *Token
	LeftChild  *node
	RightChild *node
}

// SyntaxTree represents the AST composed of nodes
type syntaxTree struct {
	Root *node
}

type evaluationResult struct {
	Type  models.DataType
	Value *models.Value
}

const (
	_treeLevelPrefix = "----"
)

type Evaluator interface {
	// Evaluate traverses through the sytanx tree with provided set of variable values
	// and evaluate the expression value
	Evaluate(tree *syntaxTree, values map[string]interface{}) (*evaluationResult, error)
}

type evaluator struct {
	resultPool sync.Pool
}

func NewEvaluator() Evaluator {
	return &evaluator{
		resultPool: sync.Pool{
			New: func() interface{} {
				return &evaluationResult{
					Value: &models.Value{},
				}
			},
		},
	}
}

func (e *evaluator) Evaluate(tree *syntaxTree, values map[string]interface{}) (*evaluationResult, error) {
	return e.evaluteHelper(tree.Root, values)
}

func (e *evaluator) stringEvaluationResult(val string) *evaluationResult {
	res := e.resultPool.Get().(*evaluationResult)
	res.Type = models.DataTypeString
	res.Value.String = &val
	return res
}

func (e *evaluator) numberEvaluationResult(val float64) *evaluationResult {
	res := e.resultPool.Get().(*evaluationResult)
	res.Type = models.DataTypeNumber
	res.Value.Number = &val
	return res
}

func (e *evaluator) boolEvaluationResult(val bool) *evaluationResult {
	res := e.resultPool.Get().(*evaluationResult)
	res.Type = models.DataTypeBool
	res.Value.Bool = &val
	return res
}

func (e *evaluator) returnResultToPool(res *evaluationResult) {
	res.Value.String = nil
	res.Value.Number = nil
	res.Value.Bool = nil
	e.resultPool.Put(res)
}

func (e *evaluator) evaluteHelper(curr *node, values map[string]interface{}) (*evaluationResult, error) {
	switch curr.Token.Type {
	case Variable:
		return e.resolveVariableValue(curr.Token, values)
	case String:
		return e.stringEvaluationResult(curr.Token.Value.(string)), nil
	case Number:
		return e.numberEvaluationResult(curr.Token.Value.(float64)), nil
	case Bool:
		return e.boolEvaluationResult(curr.Token.Value.(bool)), nil
	case Operator:
		res1, err := e.evaluteHelper(curr.LeftChild, values)
		if err != nil {
			return nil, err
		}

		res2, err := e.evaluteHelper(curr.RightChild, values)
		if err != nil {
			return nil, err
		}

		res, err := e.applyOperator(res1, res2, curr.Token)
		if err != nil {
			return nil, err
		}

		e.returnResultToPool(res1)
		e.returnResultToPool(res2)

		return res, nil
	}
	return nil, fmt.Errorf("unsupported token type %v", curr.Token.Type)
}

func (e *evaluator) resolveVariableValue(token *Token, values map[string]interface{}) (*evaluationResult, error) {
	val, ok := values[token.Value.(string)]
	if !ok {
		return nil, fmt.Errorf("error value not provided for variable %v", token.Value)
	}
	// TODO: Add nested feature resolution here
	switch val.(type) {
	case string:
		return e.stringEvaluationResult(val.(string)), nil
	case float64:
		return e.numberEvaluationResult(val.(float64)), nil
	case int:
		return e.numberEvaluationResult(float64(val.(int))), nil
	case bool:
		return e.boolEvaluationResult(val.(bool)), nil
	}
	return nil, fmt.Errorf("invalid variable type %v", val)
}

// Print is a utility method to visualize the AST
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
		inorderTraversal(node.LeftChild, nextPrefix, level+1)
	}
	fmt.Printf("|\n|%v> %v [%v]\n", prefix, node.Token.Value, node.Token.Type)
	if node.Token.Type == Operator {
		inorderTraversal(node.RightChild, nextPrefix, level+1)
	}
}

// Operations
func (e *evaluator) applyOperator(res1 *evaluationResult, res2 *evaluationResult, operation *Token) (res *evaluationResult, err error) {
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
		return e.add(res1, res2)
	case "-":
		return e.sub(res1, res2)
	case "*":
		return e.mul(res1, res2)
	case "/":
		return e.div(res1, res2)
	case "<":
		return e.lt(res1, res2)
	case ">":
		return e.gt(res1, res2)
	case "<=":
		return e.lte(res1, res2)
	case ">=":
		return e.gte(res1, res2)
	case "==":
		return e.equal(res1, res2)
	case "||":
		return e.or(res1, res2)
	case "&&":
		return e.and(res1, res2)
	}
	return nil, errors.New(ErrUnsupportedOperation, fmt.Errorf("unsupported operator %v at position %v", operation.Value, operation.Index))
}

func (e *evaluator) add(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case models.DataTypeString:
		return e.stringEvaluationResult(fmt.Sprintf("%v%v", *res1.Value.String, *res2.Value.String)), nil
	case models.DataTypeNumber:
		return e.numberEvaluationResult(*res1.Value.Number + *res2.Value.Number), nil
	}
	return nil, incompatibleOperationError("+", res1.Type)
}

func (e *evaluator) sub(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case models.DataTypeNumber:
		return e.numberEvaluationResult(*res1.Value.Number - *res2.Value.Number), nil
	}
	return nil, incompatibleOperationError("-", res1.Type)
}

func (e *evaluator) mul(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case models.DataTypeNumber:
		return e.numberEvaluationResult(*res1.Value.Number * *res2.Value.Number), nil
	}
	return nil, incompatibleOperationError("*", res1.Type)
}

func (e *evaluator) div(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case models.DataTypeNumber:
		if *res2.Value.Number == 0 {
			return nil, fmt.Errorf("encountered 0 value as denominatior")
		}
		return e.numberEvaluationResult(*res1.Value.Number / *res2.Value.Number), nil
	}
	return nil, incompatibleOperationError("/", res1.Type)
}

func (e *evaluator) lt(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case models.DataTypeString:
		return e.boolEvaluationResult(*res1.Value.String < *res2.Value.String), nil
	case models.DataTypeNumber:
		return e.boolEvaluationResult(*res1.Value.Number < *res2.Value.Number), nil
	}
	return nil, incompatibleOperationError("<", res1.Type)
}

func (e *evaluator) gt(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	res, err := e.lt(res2, res1)
	if err != nil {
		if err.(*errors.Error).Code == ErrIncompatibleOperation {
			err = incompatibleOperationError(">", res1.Type)
		}
		return nil, err
	}
	return res, nil
}

func (e *evaluator) lte(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case models.DataTypeString:
		return e.boolEvaluationResult(*res1.Value.String <= *res2.Value.String), nil
	case models.DataTypeNumber:
		return e.boolEvaluationResult(*res1.Value.Number <= *res2.Value.Number), nil
	}
	return nil, incompatibleOperationError("<=", res1.Type)
}

func (e *evaluator) gte(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	res, err := e.lte(res2, res1)
	if err != nil {
		if err.(*errors.Error).Code == ErrIncompatibleOperation {
			err = incompatibleOperationError(">=", res1.Type)
		}
		return nil, err
	}
	return res, nil
}

func (e *evaluator) equal(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case models.DataTypeString:
		return e.boolEvaluationResult(*res1.Value.String == *res2.Value.String), nil
	case models.DataTypeNumber:
		return e.boolEvaluationResult(*res1.Value.Number == *res2.Value.Number), nil
	case models.DataTypeBool:
		return e.boolEvaluationResult(*res1.Value.Bool == *res2.Value.Bool), nil
	}
	return nil, incompatibleOperationError("==", res1.Type)
}

func (e *evaluator) or(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case models.DataTypeBool:
		return e.boolEvaluationResult(*res1.Value.Bool || *res2.Value.Bool), nil
	}
	return nil, incompatibleOperationError("||", res1.Type)
}

func (e *evaluator) and(res1 *evaluationResult, res2 *evaluationResult) (*evaluationResult, error) {
	switch res1.Type {
	case models.DataTypeBool:
		return e.boolEvaluationResult(*res1.Value.Bool && *res2.Value.Bool), nil
	}
	return nil, incompatibleOperationError("&&", res1.Type)
}

func incompatibleOperationError(op string, operandType models.DataType) *errors.Error {
	return errors.New(ErrIncompatibleOperation, fmt.Errorf("operation '%v' is not compatible with '%v' type", op, operandType))
}
