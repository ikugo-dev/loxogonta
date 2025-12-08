package ast

import (
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

type Expression interface {
	foo()
}

type Literal struct {
	Value any
}
type Grouping struct {
	Expression Expression
}
type Unary struct {
	Operator tok.Token
	Right    Expression
}
type Binary struct {
	Left     Expression
	Operator tok.Token
	Right    Expression
}
type Variable struct {
	Name tok.Token
}
type Assign struct {
	Name  tok.Token
	Value Expression
}
type Logical struct {
	Left     Expression
	Operator tok.Token
	Right    Expression
}

func (e *Literal) foo()  {}
func (e *Grouping) foo() {}
func (e *Unary) foo()    {}
func (e *Binary) foo()   {}
func (e *Variable) foo() {}
func (e *Assign) foo()   {}
func (e *Logical) foo()  {}
