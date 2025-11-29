package scanner

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func (t Token) ToString() string {
	return fmt.Sprintf("%s %s %v", t.tokenType.toString(), t.lexeme, t.literal)
}
