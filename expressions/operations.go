package expressions

import (
	"fmt"

	"github.com/anshal21/coffee-machine/lib"
	"github.com/anshal21/coffee-machine/lib/errors"
	"github.com/anshal21/coffee-machine/lib/models"
)

type OperationResult = evaluationResult

type UDF struct {
	Token    string
	BinaryOp OperatorFunc
}

type OperatorFunc func(operandA, operandB, output *OperationResult) error

type OperatorFactory interface {
	Get(operatorToken string) (OperatorFunc, error)
}

func NewOperatorFactory() OperatorFactory {
	return &operatorFactory{}
}

func NewOperatorFactoryWithUDFs(ops ...UDF) OperatorFactory {
	customOperator := make(map[string]OperatorFunc)

	for _, op := range ops {
		customOperator[op.Token] = op.BinaryOp
	}

	return &operatorFactory{
		customOperator: customOperator,
	}
}

type operatorFactory struct {
	customOperator map[string]OperatorFunc
}

func (o *operatorFactory) Get(operatorToken string) (OperatorFunc, error) {
	if op, ok := o.customOperator[operatorToken]; ok {
		return op, nil
	}

	op, err := getOperator(operatorToken)
	if err == nil {
		return op, nil
	}
	return nil, errors.New(ErrUnsupportedOperation, fmt.Errorf("unsupported operator"))
}

func getOperator(operatorToken string) (OperatorFunc, error) {
	switch operatorToken {
	case "+":
		return add, nil
	case "-":
		return sub, nil
	case "*":
		return mul, nil
	case "/":
		return div, nil
	case "<":
		return lt, nil
	case ">":
		return gt, nil
	case "<=":
		return lte, nil
	case ">=":
		return gte, nil
	case "==":
		return equal, nil
	case "||":
		return or, nil
	case "&&":
		return and, nil
	default:
		return nil, errors.New(ErrUnsupportedOperation, fmt.Errorf("unsupported operator"))
	}
}

func sameOperand(op OperatorFunc, operatorToken string) OperatorFunc {
	return func(operandA, operandB, output *OperationResult) error {
		if operandA.Type != operandB.Type {
			return errors.New(ErrIncompatibleOperation, fmt.Errorf("cannot apply '%v' operation on type '%v' and '%v'", operatorToken, operandA.Type, operandB.Type))
		}

		return op(operandA, operandB, output)
	}
}

func add(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	res.Type = operand1.Type
	switch operand1.Type {
	case models.DataTypeString:
		res.Value.String = lib.StrPtr(fmt.Sprintf("%v%v", *operand1.Value.String, *operand2.Value.String))
		return nil
	case models.DataTypeNumber:
		res.Value.Number = lib.Float64Ptr(*operand1.Value.Number + *operand2.Value.Number)
		return nil
	}
	return incompatibleOperationError("+", operand1.Type)
}

func sub(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	res.Type = operand1.Type
	switch operand1.Type {
	case models.DataTypeNumber:
		res.Value.Number = lib.Float64Ptr(*operand1.Value.Number - *operand2.Value.Number)
		return nil
	}
	return incompatibleOperationError("-", operand1.Type)
}

func mul(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	res.Type = operand1.Type
	switch operand1.Type {
	case models.DataTypeNumber:
		res.Value.Number = lib.Float64Ptr(*operand1.Value.Number * *operand2.Value.Number)
		return nil
	}
	return incompatibleOperationError("*", operand1.Type)
}

func div(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	res.Type = operand1.Type
	switch operand1.Type {
	case models.DataTypeNumber:
		if *operand2.Value.Number == 0 {
			return fmt.Errorf("encountered 0 value as denominatior")
		}
		res.Value.Number = lib.Float64Ptr(*operand1.Value.Number / *operand2.Value.Number)
		return nil
	}
	return incompatibleOperationError("/", operand1.Type)
}

func lt(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	res.Type = models.DataTypeBool
	switch operand1.Type {
	case models.DataTypeString:
		res.Value.Bool = lib.BoolPtr(*operand1.Value.String < *operand2.Value.String)
		return nil
	case models.DataTypeNumber:
		res.Value.Bool = lib.BoolPtr(*operand1.Value.Number < *operand2.Value.Number)
		return nil
	}
	return incompatibleOperationError("<", operand1.Type)
}

func gt(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	err := lt(operand2, operand1, res)
	if err != nil {
		if err.(*errors.Error).Code == ErrIncompatibleOperation {
			err = incompatibleOperationError(">", operand1.Type)
		}
		return err
	}
	return nil
}

func lte(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	res.Type = models.DataTypeBool
	switch operand1.Type {
	case models.DataTypeString:
		res.Value.Bool = lib.BoolPtr(*operand1.Value.String <= *operand2.Value.String)
		return nil
	case models.DataTypeNumber:
		res.Value.Bool = lib.BoolPtr(*operand1.Value.Number <= *operand2.Value.Number)
		return nil
	}
	return incompatibleOperationError("<=", operand1.Type)
}

func gte(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	res.Type = models.DataTypeBool
	err := lte(operand2, operand1, res)
	if err != nil {
		if err.(*errors.Error).Code == ErrIncompatibleOperation {
			err = incompatibleOperationError(">=", operand1.Type)
		}
		return err
	}
	return nil
}

func equal(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	res.Type = models.DataTypeBool
	switch operand1.Type {
	case models.DataTypeString:
		res.Value.Bool = lib.BoolPtr(*operand1.Value.String == *operand2.Value.String)
		return nil
	case models.DataTypeNumber:
		res.Value.Bool = lib.BoolPtr(*operand1.Value.Number == *operand2.Value.Number)
		return nil
	case models.DataTypeBool:
		res.Value.Bool = lib.BoolPtr(*operand1.Value.Bool == *operand2.Value.Bool)
		return nil
	}
	return incompatibleOperationError("==", operand1.Type)
}

func or(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	res.Type = models.DataTypeBool
	switch operand1.Type {
	case models.DataTypeBool:
		res.Value.Bool = lib.BoolPtr(*operand1.Value.Bool || *operand2.Value.Bool)
		return nil
	}
	return incompatibleOperationError("||", operand1.Type)
}

func and(operand1 *evaluationResult, operand2 *evaluationResult, res *evaluationResult) error {
	res.Type = models.DataTypeBool
	switch operand1.Type {
	case models.DataTypeBool:
		res.Value.Bool = lib.BoolPtr(*operand1.Value.Bool && *operand2.Value.Bool)
		return nil
	}
	return incompatibleOperationError("&&", operand1.Type)
}

// func incompatibleOperationError(op string, operandType models.DataType) *errors.Error {
// 	return errors.New(ErrIncompatibleOperation, fmt.Errorf("operation '%v' is not compatible with '%v' type", op, operandType))
// }
