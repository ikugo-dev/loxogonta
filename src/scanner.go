package main

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{TokenType_Eof, "", 1, s.line})
	return s.tokens
}

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
	case '/':
		if s.match('/') { // A comment goes until the end of the line.
			// for (peek() != '\n' && !isAtEnd()) {
			for s.peek() != '\n' {
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
		error(s.line, "Unexpected character.")
	}
}
func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{tokenType, text, literal, s.line})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.source[s.current]) != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) matchAddToken(expected rune, a, b TokenType) {
	if s.match(expected) {
		s.addToken(a)
	} else {
		s.addToken(b)
	}
}

func (s *Scanner) advance() rune {
	s.current++
	return rune(s.source[s.current-1])
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\n'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
