package intr

import (
	"fmt"

	"github.com/ikugo-dev/loxogonta/internal/ast"
	"github.com/ikugo-dev/loxogonta/internal/errors"
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

func Interpret(statements []ast.Statement) {
	for _, statement := range statements {
		evalStatement(statement)
	}
}

func evalStatement(statement ast.Statement) {
	switch s := statement.(type) {
	case *ast.PrintStmt:
		fmt.Println(" --> ", eval(s.Expr))
	case *ast.ExpressionStmt:
		eval(s.Expr)
	}
}

func eval(e ast.Expression) any {
	switch expr := e.(type) {
	case *ast.Literal:
		return expr.Value
	case *ast.Grouping:
		return eval(expr.Expression)
	case *ast.Unary:
		right := eval(expr.Right)
		switch expr.Operator.TokenType {
		case tok.TokenType_Minus:
			if !areNumbers(right) {
				errors.ReportRuntime(0, "-", "Operand must be number")
				return nil
			}
			return -right.(float64)
		case tok.TokenType_Bang:
			return !isTruthy(right)
		}
	case *ast.Binary:
		left := eval(expr.Left)
		right := eval(expr.Right)
		switch expr.Operator.TokenType {
		case tok.TokenType_Plus:
			if !areNumbers(left, right) && !areStrings(left, right) {
				errors.ReportRuntime(0, "+", "Operands must be numbers or strings")
				return nil
			}
			switch left.(type) {
			case float64:
				return left.(float64) + right.(float64)
			case string:
				return left.(string) + right.(string)
			}
		case tok.TokenType_Minus:
			if !areNumbers(left, right) {
				errors.ReportRuntime(0, "-", "Operands must be numbers")
				return nil
			}
			return left.(float64) - right.(float64)
		case tok.TokenType_Slash:
			if !areNumbers(left, right) {
				errors.ReportRuntime(0, "/", "Operands must be numbers")
				return nil
			}
			return left.(float64) / right.(float64)
		case tok.TokenType_Star:
			if !areNumbers(left, right) {
				errors.ReportRuntime(0, "*", "Operands must be numbers")
				return nil
			}
			return left.(float64) * right.(float64)
		case tok.TokenType_Greater:
			if !areNumbers(left, right) {
				errors.ReportRuntime(0, ">", "Operands must be numbers")
				return nil
			}
			return left.(float64) > right.(float64)
		case tok.TokenType_GreaterEqual:
			if !areNumbers(left, right) {
				errors.ReportRuntime(0, ">=", "Operands must be numbers")
				return nil
			}
			return left.(float64) >= right.(float64)
		case tok.TokenType_Less:
			if !areNumbers(left, right) {
				errors.ReportRuntime(0, "<", "Operands must be numbers")
				return nil
			}
			return left.(float64) < right.(float64)
		case tok.TokenType_LessEqual:
			if !areNumbers(left, right) {
				errors.ReportRuntime(0, "<=", "Operands must be numbers")
				return nil
			}
			return left.(float64) <= right.(float64)
		case tok.TokenType_Equal:
			return isEqual(left, right)
		case tok.TokenType_BangEqual:
			return !isEqual(left, right)
		}
	default:
		panic("Unexpected expression")
	}
	return nil // Unreachable.
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

func areNumbers(args ...any) bool {
	for _, arg := range args {
		if _, ok := arg.(float64); !ok {
			return false
		}
	}
	return true
}

func areStrings(args ...any) bool {
	for _, arg := range args {
		if _, ok := arg.(string); !ok {
			return false
		}
	}
	return true
}
