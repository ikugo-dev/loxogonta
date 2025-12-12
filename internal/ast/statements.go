package ast

import tok "github.com/ikugo-dev/loxogonta/internal/tokens"

type Statement interface {
	foo()
}

type PrintStmt struct {
	Expr Expression
}
type ExpressionStmt struct {
	Expr Expression
}
type VarStmt struct {
	Name        tok.Token
	Initializer Expression
}
type BlockStmt struct {
	Statements []Statement
}
type IfStmt struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}
type WhileStmt struct {
	Condition Expression
	Body      Statement // block statement containing one or more arguments
	// but also doesnt have to be block / doesnt need { ... }
}
type FunctionStmt struct {
	Name   tok.Token
	Params []tok.Token
	Body   []Statement // avoids block statement because syntax forces { ... }
	// so why isnt it define as `Body BlockStmt`? Because were already going to
	// implement our own scope, so it would be redundant
}

func (s *PrintStmt) foo()      {}
func (s *ExpressionStmt) foo() {}
func (s *VarStmt) foo()        {}
func (e *BlockStmt) foo()      {}
func (e *IfStmt) foo()         {}
func (e *WhileStmt) foo()      {}
func (e *FunctionStmt) foo()   {}
