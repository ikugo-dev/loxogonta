package intr

import (
	"github.com/ikugo-dev/loxogonta/internal/errors"
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

type environment struct {
	values map[string]any
	parent *environment
}

func createEnvironment() *environment {
	return &environment{values: make(map[string]any), parent: nil}
}
func createEnvironmentWithParent(e *environment) *environment {
	return &environment{values: make(map[string]any), parent: e}
}

func (e *environment) put(name string, value any) {
	e.values[name] = value
}

func (e *environment) assign(token tok.Token, value any) {
	_, exists := e.values[token.Lexeme]
	if exists {
		e.put(token.Lexeme, value)
		return
	}
	if e.parent != nil {
		e.parent.put(token.Lexeme, value)
		return
	}
	errors.ReportRuntime(token.Line, "variable assignment", "Undefined variable "+token.Lexeme)
}

func (e *environment) get(token tok.Token) any {
	value, exists := e.values[token.Lexeme]
	if exists {
		return value
	}
	if e.parent != nil {
		return e.parent.get(token)
	}
	errors.ReportRuntime(token.Line, "variable reading", "Undefined variable "+token.Lexeme)
	return nil
}
