package scn

import (
	"strconv"

	"github.com/ikugo-dev/loxogonta/internal/errors"
	"github.com/ikugo-dev/loxogonta/internal/tokens"
)

type Scanner struct {
	source  string
	tokens  []tok.Token
	start   int
	current int
	line    int
}

func NewScanner(source string) Scanner {
	return Scanner{source: source, line: 1}
}

func (s *Scanner) ScanTokens() []tok.Token {
	for !s.isAtEnd() {
		s.start = s.current // We are at the beginning of the next lexeme.
		s.scanToken()
	}
	s.tokens = append(s.tokens, tok.Token{TokenType: tok.TokenType_Eof, Lexeme: "", Literal: 1, Line: s.line}) // QoL
	return s.tokens
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
		errors.Report(s.line, "", "Invalid number.")
		return
	}
	s.addTokenWithLiteral(tok.TokenType_Number, value)
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' { // Multi-line string support
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		errors.Report(s.line, "", "Unterminated string.")
		return
	}
	// The closing ".
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(tok.TokenType_String, value)
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokenType := keywords[text]
	if tokenType == tok.TokenType_InvalidToken { // If it isnt a keyword
		tokenType = tok.TokenType_Identifier
	}
	s.addToken(tokenType)
}

// TOKEN FUNCTIONS

func (s *Scanner) addToken(tokenType tok.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType tok.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, tok.Token{TokenType: tokenType, Lexeme: text, Literal: literal, Line: s.line})
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

func (s *Scanner) matchAddToken(expected rune, a, b tok.TokenType) {
	if s.match(expected) {
		s.addToken(a)
	} else {
		s.addToken(b)
	}
}
