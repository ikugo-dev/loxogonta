package intr

import (
	"github.com/ikugo-dev/loxogonta/internal/errors"
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

var values map[string]any = make(map[string]any)

func put(name string, value any) {
	values[name] = value
}

func get(token tok.Token) any {
	value, exists := values[token.Lexeme]
	if exists {
		return value
	} else {
		errors.ReportRuntime(token.Line, "interpreter variable reading", "Undefined variable "+token.Lexeme)
		return nil
	}
}
