package errors

import (
	"fmt"

	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

var HadError = false
var HadRuntimeError = false
var HadParseError = false

func formatError(line int, where, message string) string {
	return fmt.Sprintf("[line %d] Error at %s: %s", line, where, message)
}

func Report(line int, where, message string) {
	HadError = true
	fmt.Println(formatError(line, where, message))
}

func ReportRuntime(line int, where, message string) {
	HadRuntimeError = true
	fmt.Println(formatError(line, where, message))
}

func ReportToken(token tok.Token, message string) {
	HadParseError = true
	var where string
	if token.TokenType == tok.TokenType_Eof {
		where = "at end"
	} else {
		where = "at '" + token.Lexeme + "'"
	}
	Report(token.Line, where, message)
	panic(formatError(token.Line, where, message))
}
