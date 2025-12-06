package prs

import "fmt"

func ToString(e Expression) string {
	switch expr := e.(type) {

	case *Literal:
		return fmt.Sprintf("%v", expr.Value)
	case *Grouping:
		return "(" + ToString(expr.Expression) + ")"
	case *Unary:
		return "(" +
			expr.Operator.Lexeme + " " +
			ToString(expr.Right) +
			")"
	case *Binary:
		return "(" +
			expr.Operator.Lexeme + " " +
			ToString(expr.Left) + " " +
			ToString(expr.Right) +
			")"

	default:
		panic("unexpected expr")
	}
}
