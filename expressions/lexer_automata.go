package expressions

type state struct {
	currentState    TokenType
	nextValidStates map[TokenType]struct{}
}

var stateTransitions = map[TokenType]*state{
	None: &state{
		currentState: None,
		nextValidStates: map[TokenType]struct{}{
			Variable:        {},
			String:          {},
			Number:          {},
			Bool:            {},
			LeftParenthesis: {},
			Eol:             {},
		},
	},
	Variable: &state{
		currentState: Variable,
		nextValidStates: map[TokenType]struct{}{
			Operator:         {},
			Eol:              {},
			RightParenthesis: {},
		},
	},
	String: &state{
		currentState: String,
		nextValidStates: map[TokenType]struct{}{
			Operator:         {},
			Eol:              {},
			RightParenthesis: {},
		},
	},
	Number: &state{
		currentState: Number,
		nextValidStates: map[TokenType]struct{}{
			Operator:         {},
			Eol:              {},
			RightParenthesis: {},
		},
	},
	Operator: &state{
		currentState: Operator,
		nextValidStates: map[TokenType]struct{}{
			Variable:        {},
			String:          {},
			Bool:            {},
			Number:          {},
			LeftParenthesis: {},
		},
	},
	LeftParenthesis: &state{
		currentState: LeftParenthesis,
		nextValidStates: map[TokenType]struct{}{
			Variable:        {},
			String:          {},
			Number:          {},
			Bool:            {},
			LeftParenthesis: {},
		},
	},
	RightParenthesis: &state{
		currentState: RightParenthesis,
		nextValidStates: map[TokenType]struct{}{
			Operator:         {},
			Eol:              {},
			RightParenthesis: {},
		},
	},
}
