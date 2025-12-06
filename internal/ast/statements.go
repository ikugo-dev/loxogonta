package ast

type Statement interface {
	foo()
}

type PrintStmt struct {
	Expr Expression
}
type ExpressionStmt struct {
	Expr Expression
}

func (s *PrintStmt) foo()      {}
func (s *ExpressionStmt) foo() {}
