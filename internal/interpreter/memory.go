package intr

import (
	"github.com/ikugo-dev/loxogonta/internal/ast"
	"github.com/ikugo-dev/loxogonta/internal/tokens"
)

func lookUpVariable(storage *environment, expr ast.Expression, name tok.Token) any {
	distance, exists := locals[expr]
	if exists {
		return getAt(storage, name.Lexeme, distance)
	}
	return globals.get(name)
}

func getAt(storage *environment, name string, distance int) any {
	return ancestor(storage, distance).values[name]
}

func ancestor(storage *environment, distance int) *environment {
	for range distance {
		storage = storage.parent
	}
	return storage
}

func assignAt(storage *environment, token tok.Token, value any, distance int) {
	ancestor(storage, distance).values[token.Lexeme] = value
}
