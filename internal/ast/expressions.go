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

type LiteralExpr struct {
	Value any
}
type GroupingExpr struct {
	Expression Expression
}
type UnaryExpr struct {
	Operator tok.Token
	Right    Expression
}
type BinaryExpr struct {
	Left     Expression
	Operator tok.Token
	Right    Expression
}
type VariableExpr struct {
	Name tok.Token
}
type AssignExpr struct {
	Name  tok.Token
	Value Expression
}
type LogicalExpr struct {
	Left     Expression
	Operator tok.Token
	Right    Expression
}
type CallExpr struct {
	Callee      Expression
	Parenthesis tok.Token
	Arguments   []Expression
}

func (e *LiteralExpr) foo()  {}
func (e *GroupingExpr) foo() {}
func (e *UnaryExpr) foo()    {}
func (e *BinaryExpr) foo()   {}
func (e *VariableExpr) foo() {}
func (e *AssignExpr) foo()   {}
func (e *LogicalExpr) foo()  {}
func (e *CallExpr) foo()     {}
