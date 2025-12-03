package errors

import (
	"fmt"

	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

var HadError = false

func Report(line int, where, message string) string {
	HadError = true
	err := fmt.Sprintf("[line %d] Error %s + %s", line, where, message)
	fmt.Println(err)
	return err
}

type ParseError struct {
	message string
}

func ReportToken(token tok.Token, message string) {
	var where string
	if token.TokenType == tok.TokenType_Eof {
		where = "at end"
	} else {
		where = "at '" + token.Lexeme + "'"
	}
	panic(ParseError{Report(token.Line, where, message)})
}
