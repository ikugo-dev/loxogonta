package main

import (
	"strconv"
)

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
		} else {
			error(s.line, "Unexpected character.")
		}
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
		return ' '
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	s.current++
	value := s.peek()
	s.current--
	return value
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) scanNumber() {
	for isDigit(s.peek()) && !s.isAtEnd() {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) { // Look for a fractional part.
		s.advance() // Consume the "."
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	for isDigit(s.peek()) {
		s.advance()
	}
	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		error(s.line, "Invalid number.")
		return
	}
	s.addTokenWithLiteral(TokenType_Number, value)
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' { // Multi-line string support
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		error(s.line, "Unterminated string.")
		return
	}
	// The closing ".
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(TokenType_String, value)
}
