package scanner

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func (t Token) ToString() string {
	return fmt.Sprintf("%s", t.lexeme)
}

func (t Token) GetTokenType() TokenType {
	return t.tokenType
}

func (t Token) GetLiteral() any {
	return t.literal
}
