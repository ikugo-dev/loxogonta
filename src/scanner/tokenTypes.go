package scanner

type TokenType int

const (
	TokenType_InvalidToken TokenType = iota // Single-character tokens.
	TokenType_LeftParen
	TokenType_RightParen
	TokenType_LeftBrace
	TokenType_RightBrace
	TokenType_Comma
	TokenType_Dot
	TokenType_Minus
	TokenType_Plus
	TokenType_Semicolon
	TokenType_Slash
	TokenType_Star
	TokenType_Bang // One or two character tokens.
	TokenType_BangEqual
	TokenType_Equal
	TokenType_EqualEqual
	TokenType_Greater
	TokenType_GreaterEqual
	TokenType_Less
	TokenType_LessEqual
	TokenType_Identifier // Literals.
	TokenType_String
	TokenType_Number
	TokenType_And // Keywords.
	TokenType_Class
	TokenType_Else
	TokenType_False
	TokenType_Fun
	TokenType_For
	TokenType_If
	TokenType_Nil
	TokenType_Or
	TokenType_Print
	TokenType_Return
	TokenType_Super
	TokenType_This
	TokenType_True
	TokenType_Var
	TokenType_While
	TokenType_Eof
)

func (t TokenType) toString() string {
	return []string{
		"TokenType_InvalidToken", // Single-character tokens.
		"TokenType_LeftParen",
		"TokenType_RightParen",
		"TokenType_LeftBrace",
		"TokenType_RightBrace",
		"TokenType_Comma",
		"TokenType_Dot",
		"TokenType_Minus",
		"TokenType_Plus",
		"TokenType_Semicolon",
		"TokenType_Slash",
		"TokenType_Star",
		"TokenType_Bang", // One or two character tokens.
		"TokenType_BangEqual",
		"TokenType_Equal",
		"TokenType_EqualEqual",
		"TokenType_Greater",
		"TokenType_GreaterEqual",
		"TokenType_Less",
		"TokenType_LessEqual",
		"TokenType_Identifier", // Literals.
		"TokenType_String",
		"TokenType_Number",
		"TokenType_And", // Keywords.
		"TokenType_Class",
		"TokenType_Else",
		"TokenType_False",
		"TokenType_Fun",
		"TokenType_For",
		"TokenType_If",
		"TokenType_Nil",
		"TokenType_Or",
		"TokenType_Print",
		"TokenType_Return",
		"TokenType_Super",
		"TokenType_This",
		"TokenType_True",
		"TokenType_Var",
		"TokenType_While",
		"TokenType_Eof",
	}[t]
}
