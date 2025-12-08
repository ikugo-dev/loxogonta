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
	Body      Statement
}

func (s *PrintStmt) foo()      {}
func (s *ExpressionStmt) foo() {}
func (s *VarStmt) foo()        {}
func (e *BlockStmt) foo()      {}
func (e *IfStmt) foo()         {}
func (e *WhileStmt) foo()      {}
