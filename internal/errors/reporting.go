package errors

import "fmt"

var HadError = false

func Report(line int, where, message string) {
	fmt.Printf("[line %d] Error %s + %s", line, where, message)
	HadError = true
}

func ReportToken(token , String message) {
    if (token.type == TokenType.EOF) {
      report(token.line, " at end", message);
    } else {
      report(token.line, " at '" + token.lexeme + "'", message);
    }
  }
