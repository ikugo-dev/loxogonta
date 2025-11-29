package parser

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

func (e *Literal) foo()  {}
func (e *Grouping) foo() {}
func (e *Unary) foo()    {}
func (e *Binary) foo()   {}
