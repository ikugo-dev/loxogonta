package prs

import "fmt"

func ToString(e expression) string {
	switch expr := e.(type) {

	case *literal:
		return fmt.Sprintf("%v", expr.value)
	case *grouping:
		return "(" + ToString(expr.expression) + ")"
	case *unary:
		return "(" +
			expr.operator.Lexeme + " " +
			ToString(expr.right) +
			")"
	case *binary:
		return "(" +
			expr.operator.Lexeme + " " +
			ToString(expr.left) + " " +
			ToString(expr.right) +
			")"

	default:
		panic("unexpected expr")
	}
}
