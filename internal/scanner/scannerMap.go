package scn

import (
	"github.com/ikugo-dev/loxogonta/internal/errors"
	t "github.com/ikugo-dev/loxogonta/internal/tokens"
)

func (s *Scanner) scanToken() {
	var character rune = s.advance()
	switch character {
	case '(':
		s.addToken(t.TokenType_LeftParen)
	case ')':
		s.addToken(t.TokenType_RightParen)
	case '{':
		s.addToken(t.TokenType_LeftBrace)
	case '}':
		s.addToken(t.TokenType_RightBrace)
	case ',':
		s.addToken(t.TokenType_Comma)
	case '.':
		s.addToken(t.TokenType_Dot)
	case '-':
		s.addToken(t.TokenType_Minus)
	case '+':
		s.addToken(t.TokenType_Plus)
	case ';':
		s.addToken(t.TokenType_Semicolon)
	case '*':
		s.addToken(t.TokenType_Star)
	case '!':
		s.matchAddToken('=', t.TokenType_BangEqual, t.TokenType_Bang)
	case '=':
		s.matchAddToken('=', t.TokenType_EqualEqual, t.TokenType_Equal)
	case '<':
		s.matchAddToken('=', t.TokenType_LessEqual, t.TokenType_Less)
	case '>':
		s.matchAddToken('=', t.TokenType_GreaterEqual, t.TokenType_Greater)
	case '"':
		s.scanString()
	case '/':
		if s.match('/') { // A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(t.TokenType_Slash)
		}
	case ' ': // Ignore whitespace.
	case '\r':
	case '\t':
	case '\n':
		s.line++
	default:
		if isDigit(character) {
			s.scanNumber()
		} else if isAlpha(character) {
			s.identifier()
		} else {
			errors.Report(s.line, "", "Unexpected character.")
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
