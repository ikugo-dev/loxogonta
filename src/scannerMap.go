package main

func (s *Scanner) scanToken() {
	var c rune = s.advance()
	switch c {
	case '(':
		s.addToken(TokenType_LeftParen)
	case ')':
		s.addToken(TokenType_RightParen)
	case '{':
		s.addToken(TokenType_LeftBrace)
	case '}':
		s.addToken(TokenType_RightBrace)
	case ',':
		s.addToken(TokenType_Comma)
	case '.':
		s.addToken(TokenType_Dot)
	case '-':
		s.addToken(TokenType_Minus)
	case '+':
		s.addToken(TokenType_Plus)
	case ';':
		s.addToken(TokenType_Semicolon)
	case '*':
		s.addToken(TokenType_Star)
	case '!':
		s.matchAddToken('=', TokenType_BangEqual, TokenType_Bang)
	case '=':
		s.matchAddToken('=', TokenType_EqualEqual, TokenType_Equal)
	case '<':
		s.matchAddToken('=', TokenType_LessEqual, TokenType_Less)
	case '>':
		s.matchAddToken('=', TokenType_GreaterEqual, TokenType_Greater)
	case '"':
		s.scanString()
	case '/':
		if s.match('/') { // A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TokenType_Slash)
		}
	case ' ': // Ignore whitespace.
	case '\r':
	case '\t':
	case '\n':
		s.line++
	default:
		if isDigit(c) {
			s.scanNumber()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			error(s.line, "Unexpected character.")
		}
	}
}

var keywords map[string]TokenType = map[string]TokenType{
	"and":    TokenType_And,
	"class":  TokenType_Class,
	"else":   TokenType_Else,
	"false":  TokenType_False,
	"for":    TokenType_For,
	"fun":    TokenType_Fun,
	"if":     TokenType_If,
	"nil":    TokenType_Nil,
	"or":     TokenType_Or,
	"print":  TokenType_Print,
	"return": TokenType_Return,
	"super":  TokenType_Super,
	"this":   TokenType_This,
	"true":   TokenType_True,
	"var":    TokenType_Var,
	"while":  TokenType_While,
}
