package expressions

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/anshal21/coffee-machine/lib/errors"
)

var (
	_VariableRegex  *regexp.Regexp
	_DecimalRegex   *regexp.Regexp
	_ValidOperators = map[string]struct{}{
		"<":  {},
		">":  {},
		">=": {},
		"<=": {},
		"==": {},
		"+":  {},
		"-":  {},
		"/":  {},
		"*":  {},
		"^":  {},
		"||": {},
		"&&": {},
	}
)

func init() {
	_VariableRegex, _ = regexp.Compile("^[a-zA-Z_][a-zA-Z_0-9]*$")
	// TODO: this matches leading and trailing 0s need a fix for it
	_DecimalRegex, _ = regexp.Compile("^-?[0-9][0-9]*(.[0-9]+)?$")
}

// Lexer is an interface exposes a Lex function, that tokenizes the given
// input expression based on some grammar
type Lexer interface {
	Lex(expression string) ([]*Token, error)
}

type lexer struct {
}

// NewLexer is a constructor to instantiate a Lexer
func NewLexer() Lexer {
	return &lexer{}
}

func (l *lexer) Lex(expression string) ([]*Token, error) {
	expressionStream := newStream(expression)

	tokens := make([]*Token, 0)

	lexerState := stateTransitions[None]

	for {
		val := expressionStream.GetNext()
		if val == _EndOfStream {
			break
		}
		if isDelimiter(val) {
			continue
		}

		expressionStream.Rewind()
		nextToken, err := getNextToken(expressionStream)
		if err != nil {
			return nil, err
		}

		if _, ok := lexerState.nextValidStates[nextToken.Type]; !ok {
			return nil, errors.New(ErrInvalidExpression,
				fmt.Errorf("invalid predicate syntax %v cannot be followed by a %v at position %v",
					lexerState.currentState, nextToken.Type, nextToken.Index))

		}
		tokens = append(tokens, nextToken)

		lexerState = stateTransitions[nextToken.Type]
	}
	return tokens, nil
}

func getNextToken(s *stream) (*Token, error) {
	index := s.Position()
	c := s.GetNext()
	s.Rewind()
	switch c {
	case '"':
		return scanString(s)
	case '(', ')':
		return scanParenthesis(s)
	default:
		token := getNonStringToken(s)
		if isValidBool(token) {
			b, _ := strconv.ParseBool(token)
			return &Token{
				Type:  Bool,
				Value: b,
				Index: index,
			}, nil
		}
		if isValidVariable(token) {
			return &Token{
				Type:  Variable,
				Value: token,
				Index: index,
			}, nil
		}
		if isValidNumber(token) {
			number, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return nil, err
			}
			return &Token{
				Type:  Number,
				Value: number,
				Index: index,
			}, nil
		}
		if isValidOperator(token) {
			return &Token{
				Type:  Operator,
				Value: token,
				Index: index,
			}, nil
		}
		return nil, errors.New(ErrInvalidExpression, fmt.Errorf("unrecognized token %v at position %v", token, index))
	}

}

func isValidVariable(s string) bool {
	return _VariableRegex.MatchString(s)
}

func isValidNumber(s string) bool {
	return _DecimalRegex.MatchString(s)
}

func isValidOperator(s string) bool {
	_, ok := _ValidOperators[s]
	return ok
}

func isValidBool(s string) bool {
	return s == "true" || s == "false"
}

func isDelimiter(c rune) bool {
	return c == ' '
}

func getNonStringToken(s *stream) string {
	token := make([]rune, 0)
	for {
		val := s.GetNext()
		if isDelimiter(val) || val == _EndOfStream || val == ')' {
			if val != _EndOfStream {
				s.Rewind()
			}
			break
		}
		token = append(token, val)
	}
	return string(token)
}

func scanString(s *stream) (*Token, error) {
	index := s.Position()
	startQuote := s.GetNext()
	var endQuote rune
	token := make([]rune, 0)
	for {
		val := s.GetNext()
		if val == _EndOfStream {
			break
		}
		endQuote = val
		if endQuote == startQuote {
			break
		}
		token = append(token, val)
	}
	if startQuote != endQuote {
		return nil, errors.New(ErrInvalidExpression, fmt.Errorf("badly formatted string %v in the expression at position %v", string(token), index))
	}

	return &Token{
		Type:  String,
		Value: string(token),
		Index: index,
	}, nil
}

func scanParenthesis(s *stream) (*Token, error) {
	index := s.Position()
	tokenVal := make([]rune, 0, 1)
	tokenVal = append(tokenVal, s.GetNext())
	tokenStr := string(tokenVal)
	token := &Token{
		Value: tokenStr,
		Index: index,
	}
	if tokenStr == "(" {
		token.Type = LeftParenthesis
	}
	if tokenStr == ")" {
		token.Type = RightParenthesis
	}
	return token, nil
}
