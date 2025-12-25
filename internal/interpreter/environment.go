package intr

import (
	"fmt"

	"github.com/ikugo-dev/loxogonta/internal/ast"
	"github.com/ikugo-dev/loxogonta/internal/errors"
	"github.com/ikugo-dev/loxogonta/internal/resolver"
	"github.com/ikugo-dev/loxogonta/internal/tokens"
)

type environment struct {
	globals map[string]any
	parent  *environment
	locals  map[ast.Expression]int
}

func createEnvironmentWithResolution(statements []ast.Statement) *environment {
	return &environment{globals: make(map[string]any), parent: nil, locals: rslv.Resolve(statements)}
}
func createEnvironment() *environment {
	return &environment{globals: make(map[string]any), parent: nil, locals: make(map[ast.Expression]int)}
}
func createEnvironmentWithParent(e *environment) *environment {
	return &environment{globals: make(map[string]any), parent: e, locals: e.locals}
}

func (e *environment) put(name string, value any) {
	e.globals[name] = value
}

func (e *environment) assign(expr ast.Expression, value any) {
	assignExpr, ok := expr.(*ast.AssignExpr)
	if !ok {
		errors.ReportRuntime(-1, "During assign expression",
			"Could not validate expression as an assignment expression")
		return
	}
	distance, exists := e.locals[expr]
	if exists {
		e.assignAt(assignExpr.Name.Lexeme, distance, value)
		return
	}
	e.assignStatic(assignExpr.Name, value)
}

func (e *environment) assignAt(name string, distance int, value any) {
	e.ancestor(distance).globals[name] = value
}

func (e *environment) assignStatic(token tok.Token, value any) {
	if _, exists := e.globals[token.Lexeme]; exists {
		e.globals[token.Lexeme] = value
		return
	}
	if e.parent != nil {
		e.parent.assignStatic(token, value)
		return
	}
	errors.ReportRuntime(token.Line, "variable assignment", "Undefined variable "+token.Lexeme)
}

// func (e *environment) getDynamic(token tok.Token) any {
// 	value, exists := e.globals[token.Lexeme]
// 	if exists {
// 		return value
// 	}
// 	if e.parent != nil {
// 		return e.parent.getDynamic(token)
// 	}
// 	errors.ReportRuntime(token.Line, "variable reading", "Undefined variable "+token.Lexeme)
// 	return nil
// }

func (e *environment) lookUp(expr ast.Expression, name string) any {
	distance, exists := e.locals[expr]
	fmt.Printf("distance: %v\n", distance)
	if exists {
		return e.getAt(name, distance)
	}
	fmt.Printf("globals: %v\n", e.globals[name])
	return e.globals[name]
}

func (e *environment) getAt(name string, distance int) any {
	return e.ancestor(distance).globals[name]
}

func (e *environment) ancestor(distance int) *environment {
	storage := e
	for range distance {
		storage = storage.parent
	}
	return storage
}
