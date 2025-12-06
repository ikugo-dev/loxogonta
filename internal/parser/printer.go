package prs

import (
	"fmt"

	"github.com/ikugo-dev/loxogonta/internal/ast"
)

func ToString(e ast.Expression) string {
	switch expr := e.(type) {

	case *ast.Literal:
		return fmt.Sprintf("%v", expr.Value)
	case *ast.Grouping:
		return "(" + ToString(expr.Expression) + ")"
	case *ast.Unary:
		return "(" +
			expr.Operator.Lexeme + " " +
			ToString(expr.Right) +
			")"
	case *ast.Binary:
		return "(" +
			expr.Operator.Lexeme + " " +
			ToString(expr.Left) + " " +
			ToString(expr.Right) +
			")"

	default:
		panic("unexpected expr")
	}
}
