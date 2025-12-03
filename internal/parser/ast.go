package prs

import tok "github.com/ikugo-dev/loxogonta/internal/tokens"

// OLD:
// expression     → literal
//                | unary
//                | binary
//                | grouping ;
// literal        → NUMBER | STRING | "true" | "false" | "nil" ;
// grouping       → "(" expression ")" ;
// unary          → ( "-" | "!" ) expression ;
// binary         → expression operator expression ;
// operator       → "==" | "!=" | "<" | "<=" | ">" | ">="
//                | "+"  | "-"  | "*" | "/" ;

// NEW:
// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

type expression interface {
	foo()
}

type literal struct {
	value any
}

type grouping struct {
	expression expression
}
type unary struct {
	operator tok.Token
	right    expression
}
type binary struct {
	left     expression
	operator tok.Token
	right    expression
}

func (e *literal) foo()  {}
func (e *grouping) foo() {}
func (e *unary) foo()    {}
func (e *binary) foo()   {}
