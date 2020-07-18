package expressions

import (
	"fmt"

	"github.com/anshal21/coffee-machine/lib/errors"
)

// Parser interface exposes Parse function, that parses
// a stream of token and generates an AST
type Parser interface {
	Parse(tokens []*Token) (*syntaxTree, error)
}

// NewParser is a constructor to instantiate a Parser
func NewParser() Parser {
	return &parser{}
}

type parser struct {
}

func (p *parser) Parse(tokens []*Token) (*syntaxTree, error) {
	root := &node{}
	pos := 0
	err := parserHelper(tokens, &pos, root)
	if err != nil {
		return nil, err
	}
	return &syntaxTree{
		Root: root,
	}, nil
}

func parserHelper(tokens []*Token, pos *int, root *node) error {
	operandStack := newStack(len(tokens))
	operatorStack := newStack(len(tokens))

	buildExpr := func(op *Token) error {
		operand1 := toNode(operandStack.Top())
		operandStack.Pop()
		operand2 := toNode(operandStack.Top())
		operandStack.Pop()
		if operand1 == nil || operand2 == nil {
			return errors.New(ErrInvalidExpression,
				fmt.Errorf("missing operands for operator %v at position %v", op.Value, op.Index))
		}
		operandStack.Push(&node{
			Token:      op,
			LeftChild:  operand2,
			RightChild: operand1,
		})
		return nil
	}

OuterLoop:
	for _, val := range tokens {
		switch val.Type {
		case LeftParenthesis:
			operatorStack.Push(val)
		case Variable, String, Number, Bool:
			operandStack.Push(&node{
				Token: val,
			})
		case Operator:
			for {
				topEle := toToken(operatorStack.Top())
				if topEle == nil || topEle.Type == LeftParenthesis || operatorPrecedence(topEle.Value.(string)) < operatorPrecedence(val.Value.(string)) {
					operatorStack.Push(val)
					continue OuterLoop
				}
				operatorStack.Pop()
				err := buildExpr(topEle)
				if err != nil {
					return err
				}
			}
		case RightParenthesis:
			for {
				topEle := toToken(operatorStack.Top())
				operatorStack.Pop()
				if topEle == nil {
					return errors.New(ErrInvalidExpression, fmt.Errorf("no matching '(' for ')' at position %v", val.Index))
				}
				if topEle.Type == LeftParenthesis {
					continue OuterLoop
				}
				err := buildExpr(topEle)
				if err != nil {
					return err
				}
			}
		}

	}

	for operatorStack.Top() != nil {
		topEle := toToken(operatorStack.Top())
		if topEle.Type == LeftParenthesis {
			return errors.New(ErrInvalidExpression, fmt.Errorf("no matching ')' for '(' at position %v", topEle.Index))
		}
		err := buildExpr(topEle)
		if err != nil {
			return err
		}
		operatorStack.Pop()
	}

	*root = *(toNode(operandStack.Top()))
	return nil
}

func operatorPrecedence(op string) int {
	switch op {
	case "^":
		return 4
	case "*", "/":
		return 3
	case "+", "-":
		return 2
	case ">", "<", "==", ">=", "<=":
		return 1
	default:
		return -1
	}
}

func toToken(val interface{}) *Token {
	if val != nil {
		return val.(*Token)
	}
	return nil
}

func toNode(val interface{}) *node {
	if val != nil {
		return val.(*node)
	}
	return nil
}
