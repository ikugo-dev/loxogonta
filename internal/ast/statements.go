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

func (s *PrintStmt) foo()      {}
func (s *ExpressionStmt) foo() {}
func (s *VarStmt) foo()        {}
