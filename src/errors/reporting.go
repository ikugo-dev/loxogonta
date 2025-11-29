package errors

import "fmt"

var HadError = false

func Report(line int, where, message string) {
	fmt.Printf("[line %d] Error %s + %s", line, where, message)
	HadError = true
}
