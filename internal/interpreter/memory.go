package intr

import (
	"fmt"

	"github.com/ikugo-dev/loxogonta/internal/ast"
	// "github.com/ikugo-dev/loxogonta/internal/errors"
	// "github.com/ikugo-dev/loxogonta/internal/resolver"
	"github.com/ikugo-dev/loxogonta/internal/tokens"
)

// func assign(expr ast.Expression, value any) {
// 	assignExpr, ok := expr.(*ast.AssignExpr)
// 	if !ok {
// 		errors.ReportRuntime(-1, "During assign expression",
// 			"Could not validate expression as an assignment expression")
// 		return
// 	}
// 	distance, exists := locals[expr]
// 	if exists {
// 		assignAt(assignExpr.Name.Lexeme, distance, value)
// 		return
// 	}
// 	env.assign(assignExpr.Name, value)
// }

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
