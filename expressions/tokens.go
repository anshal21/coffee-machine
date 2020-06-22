package expressions

// TokenType is a type to represent different token type allowed in an expression
type TokenType int

// TokenValue is a type to hold value of a token
type TokenValue interface{}

// Set of allowed tokens in an expression
const (
	None TokenType = iota
	Variable
	String
	Number
	Bool
	Operator
	LeftParenthesis
	RightParenthesis
	KeyWord
	Eol
	Unknown
)

func (t TokenType) String() string {
	switch t {
	case None:
		return "None"
	case Variable:
		return "Variable"
	case String:
		return "String"
	case Number:
		return "Number"
	case Bool:
		return "Bool"
	case Operator:
		return "Operator"
	case LeftParenthesis:
		return "LeftParenthesis"
	case RightParenthesis:
		return "RightParenthesis"
	case Eol:
		return "Eol"
	default:
		return "Unknown"
	}
}

// Token represents some token in the input expression
type Token struct {
	Type  TokenType
	Value TokenValue
	Index int
}
