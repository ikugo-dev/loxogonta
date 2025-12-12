package scn

import (
	"github.com/ikugo-dev/loxogonta/internal/errors"
	t "github.com/ikugo-dev/loxogonta/internal/tokens"
)

func scanToken() {
	var character rune = advance()
	switch character {
	case '(':
		addToken(t.TokenType_LeftParen)
	case ')':
		addToken(t.TokenType_RightParen)
	case '{':
		addToken(t.TokenType_LeftBrace)
	case '}':
		addToken(t.TokenType_RightBrace)
	case ',':
		addToken(t.TokenType_Comma)
	case '.':
		addToken(t.TokenType_Dot)
	case '-':
		addToken(t.TokenType_Minus)
	case '+':
		addToken(t.TokenType_Plus)
	case ';':
		addToken(t.TokenType_Semicolon)
	case '*':
		addToken(t.TokenType_Star)
	case '%':
		addToken(t.TokenType_Percentage)
	case '!':
		matchAddToken('=', t.TokenType_BangEqual, t.TokenType_Bang)
	case '=':
		matchAddToken('=', t.TokenType_EqualEqual, t.TokenType_Equal)
	case '<':
		matchAddToken('=', t.TokenType_LessEqual, t.TokenType_Less)
	case '>':
		matchAddToken('=', t.TokenType_GreaterEqual, t.TokenType_Greater)
	case '"':
		scanString()
	case '/':
		if match('/') { // A comment goes until the end of the line.
			for peek() != '\n' && !isAtEnd() {
				advance()
			}
		} else {
			addToken(t.TokenType_Slash)
		}
	case ' ': // Ignore whitespace.
	case '\r':
	case '\t':
	case '\n':
		line++
	default:
		if isDigit(character) {
			scanNumber()
		} else if isAlpha(character) {
			identifier()
		} else {
			errors.Report(line, "", "Unexpected character.")
		}
	}
}

var keywords map[string]t.TokenType = map[string]t.TokenType{
	"and":    t.TokenType_And,
	"class":  t.TokenType_Class,
	"else":   t.TokenType_Else,
	"false":  t.TokenType_False,
	"for":    t.TokenType_For,
	"fun":    t.TokenType_Fun,
	"if":     t.TokenType_If,
	"nil":    t.TokenType_Nil,
	"or":     t.TokenType_Or,
	"print":  t.TokenType_Print,
	"return": t.TokenType_Return,
	"super":  t.TokenType_Super,
	"this":   t.TokenType_This,
	"true":   t.TokenType_True,
	"var":    t.TokenType_Var,
	"while":  t.TokenType_While,
}
