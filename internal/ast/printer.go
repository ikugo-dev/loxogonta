package ast

import "fmt"

func ToString(e Expression) string {
	switch expr := e.(type) {

	case *LiteralExpr:
		return fmt.Sprintf("%v", expr.Value)
	case *GroupingExpr:
		return "(" + ToString(expr.Expression) + ")"
	case *UnaryExpr:
		return "(" +
			expr.Operator.Lexeme + " " +
			ToString(expr.Right) +
			")"
	case *BinaryExpr:
		return "(" +
			expr.Operator.Lexeme + " " +
			ToString(expr.Left) + " " +
			ToString(expr.Right) +
			")"

	default:
		panic("Unexpected expression")
	}
}
