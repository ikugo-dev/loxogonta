package parser

import "fmt"

func Print(e Expression) string {
	switch expr := e.(type) {

	case *Literal:
		return fmt.Sprintf("%v", expr.Value)
	case *Grouping:
		return "(" + Print(expr.Expression) + ")"
	case *Unary:
		return "(" +
			expr.Operator.ToString() + " " +
			Print(expr.Right) +
			")"
	case *Binary:
		return "(" +
			expr.Operator.ToString() + " " +
			Print(expr.Left) + " " +
			Print(expr.Right) +
			")"

	default:
		panic("unexpected expr")
	}
}
