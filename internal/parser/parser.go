package prs

import (
	"slices"

	"github.com/ikugo-dev/loxogonta/internal/ast"
	"github.com/ikugo-dev/loxogonta/internal/errors"
	"github.com/ikugo-dev/loxogonta/internal/tokens"
)

var tokens []tok.Token
var current int

func resetData(newTokens []tok.Token) {
	tokens = newTokens
	current = 0
}

func ParseTokens(tokens []tok.Token) []ast.Statement {
	resetData(tokens)
	return program()
}

func match(tokenTypes ...tok.TokenType) bool {
	if slices.ContainsFunc(tokenTypes, check) {
		advance()
		return true
	}
	return false
}

func check(tokenType tok.TokenType) bool {
	if isAtEnd() {
		return false
	}
	return peek().TokenType == tokenType
}

func advance() tok.Token {
	if !isAtEnd() {
		current++
	}
	return previous()
}

func isAtEnd() bool {
	return peek().TokenType == tok.TokenType_Eof
}

func peek() tok.Token {
	return tokens[current]
}

func previous() tok.Token {
	return tokens[current-1]
}

func consume(tokenType tok.TokenType, message string) tok.Token { // is return needed?
	if check(tokenType) {
		return advance()
	}
	errors.ReportToken(peek(), message)
	return tok.Token{}
}

func synchronize() {
	advance()
	for !isAtEnd() {
		if previous().TokenType == tok.TokenType_Semicolon {
			return
		}
		switch peek().TokenType { // finding end
		case tok.TokenType_Class:
			return
		case tok.TokenType_Fun:
			return
		case tok.TokenType_Var:
			return
		case tok.TokenType_For:
			return
		case tok.TokenType_If:
			return
		case tok.TokenType_While:
			return
		case tok.TokenType_Print:
			return
		case tok.TokenType_Return:
			return
		}
		advance()
	}
}
