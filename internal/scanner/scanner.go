package scn

import (
	"strconv"

	"github.com/ikugo-dev/loxogonta/internal/errors"
	"github.com/ikugo-dev/loxogonta/internal/tokens"
)

var source string
var tokens []tok.Token
var start int
var current int
var line int

func resetData(newSource string) {
	source = newSource
	tokens = []tok.Token{}
	start = 0
	current = 0
	line = 1
}

func ScanSource(source string) []tok.Token {
	resetData(source)
	for !isAtEnd() {
		start = current // We are at the beginning of the next lexeme.
		scanToken()
	}
	tokens = append(tokens, tok.Token{TokenType: tok.TokenType_Eof, Lexeme: "", Literal: 1, Line: line}) // QoL
	return tokens
}

func advance() rune {
	current++
	return rune(source[current-1])
}

func peek() rune {
	if isAtEnd() {
		return ' '
	}
	return rune(source[current])
}

func peekNext() rune {
	current++
	value := peek()
	current--
	return value
}

func isAtEnd() bool {
	return current >= len(source)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func scanNumber() {
	for isDigit(peek()) && !isAtEnd() {
		advance()
	}
	if peek() == '.' && isDigit(peekNext()) { // Look for a fractional part.
		advance() // Consume the "."
		for isDigit(peek()) {
			advance()
		}
	}
	for isDigit(peek()) {
		advance()
	}
	value, err := strconv.ParseFloat(source[start:current], 64)
	if err != nil {
		errors.Report(line, "", "Invalid number.")
		return
	}
	addTokenWithLiteral(tok.TokenType_Number, value)
}

func scanString() {
	for peek() != '"' && !isAtEnd() {
		if peek() == '\n' { // Multi-line string support
			line++
		}
		advance()
	}
	if isAtEnd() {
		errors.Report(line, "", "Unterminated string.")
		return
	}
	// The closing ".
	advance()
	value := source[start+1 : current-1]
	addTokenWithLiteral(tok.TokenType_String, value)
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func identifier() {
	for isAlphaNumeric(peek()) {
		advance()
	}
	text := source[start:current]
	tokenType := keywords[text]
	if tokenType == tok.TokenType_InvalidToken { // If it isnt a keyword
		tokenType = tok.TokenType_Identifier
	}
	addToken(tokenType)
}

// TOKEN FUNCTIONS

func addToken(tokenType tok.TokenType) {
	addTokenWithLiteral(tokenType, nil)
}

func addTokenWithLiteral(tokenType tok.TokenType, literal any) {
	text := source[start:current]
	tokens = append(tokens, tok.Token{TokenType: tokenType, Lexeme: text, Literal: literal, Line: line})
}

func match(expected rune) bool {
	if isAtEnd() {
		return false
	}
	if rune(source[current]) != expected {
		return false
	}
	current++
	return true
}

func matchAddToken(expected rune, a, b tok.TokenType) {
	if match(expected) {
		addToken(a)
	} else {
		addToken(b)
	}
}
