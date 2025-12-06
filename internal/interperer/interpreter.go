package intr

import (
	prs "github.com/ikugo-dev/loxogonta/internal/parser"
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

func Interpret(e prs.Expression) any {
	switch expr := e.(type) {

	case *prs.Literal:
		return expr.Value
	case *prs.Grouping:
		Interpret(expr.Expression)
	case *prs.Unary:
		right := Interpret(expr.Right)
		switch expr.Operator.TokenType {
		case tok.TokenType_Minus:
			return -right.(float64)
		case tok.TokenType_Bang:
			return !isTruthy(right)
		}
		return nil // Unreachable.
	case *prs.Binary:
		left := Interpret(expr.Left)
		right := Interpret(expr.Right)
		switch expr.Operator.TokenType {
		case tok.TokenType_Plus:
			switch left.(type) {
			case float64:
				return left.(float64) + right.(float64)
			case string:
				return left.(string) + right.(string)
			}
		case tok.TokenType_Minus:
			return left.(float64) - right.(float64)
		case tok.TokenType_Slash:
			return left.(float64) / right.(float64)
		case tok.TokenType_Star:
			return left.(float64) * right.(float64)
		case tok.TokenType_Greater:
			return left.(float64) > right.(float64)
		case tok.TokenType_GreaterEqual:
			return left.(float64) >= right.(float64)
		case tok.TokenType_Less:
			return left.(float64) < right.(float64)
		case tok.TokenType_LessEqual:
			return left.(float64) <= right.(float64)
		case tok.TokenType_Equal:
			return isEqual(left, right)
		case tok.TokenType_BangEqual:
			return !isEqual(left, right)
		}
		return nil // Unreachable.
	default:
		panic("unexpected expr")
	}
	return nil
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if boolean, ok := value.(bool); ok {
		return boolean
	}
	return true
}

func isEqual(left any, right any) bool {
	if left == nil {
		return right == nil
	}
	return left == right
}
