package intr

import (
	"fmt"
	"math"

	"github.com/ikugo-dev/loxogonta/internal/ast"
	"github.com/ikugo-dev/loxogonta/internal/errors"
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

var storage environment = createEnvironment()

func Interpret(statements []ast.Statement) any {
	addNativeFunctions()
	var value any
	for _, statement := range statements {
		value = evalStmt(statement)
	}
	return value
}

func evalStmt(statement ast.Statement) any {
	switch stmt := statement.(type) {
	case *ast.PrintStmt:
		fmt.Println(evalExpr(stmt.Expr))
	case *ast.ExpressionStmt:
		return evalExpr(stmt.Expr)
	case *ast.VarStmt:
		var value any = nil
		if stmt.Initializer != nil {
			value = evalExpr(stmt.Initializer)
		}
		storage.put(stmt.Name.Lexeme, value)
	case *ast.BlockStmt:
		var value any = nil
		oldStorage := storage
		storage = createEnvironmentWithParent(oldStorage)
		for _, statement := range stmt.Statements {
			value = evalStmt(statement)
		}
		storage = oldStorage
		return value
	case *ast.IfStmt:
		if isTruthy(evalExpr(stmt.Condition)) {
			return evalStmt(stmt.ThenBranch)
		} else {
			return evalStmt(stmt.ElseBranch)
		}
	case *ast.WhileStmt:
		for isTruthy(evalExpr(stmt.Condition)) {
			evalStmt(stmt.Body)
		}
	case *ast.FunctionStmt:
		function := &loxFunction{declaration: stmt, closure: &storage}
		storage.put(stmt.Name.Lexeme, function)
	}
	return nil
}

func evalExpr(expression ast.Expression) any {
	switch expr := expression.(type) {
	case *ast.LiteralExpr:
		return expr.Value
	case *ast.GroupingExpr:
		return evalExpr(expr.Expression)
	case *ast.UnaryExpr:
		right := evalExpr(expr.Right)
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
	case *ast.BinaryExpr:
		left := evalExpr(expr.Left)
		right := evalExpr(expr.Right)
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
		case tok.TokenType_Percentage:
			if !areNumbers(left, right) {
				errors.ReportRuntime(0, "%", "Operands must be numbers")
				return nil
			}
			return math.Mod(left.(float64), right.(float64))
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
		case tok.TokenType_EqualEqual:
			return isEqual(left, right)
		case tok.TokenType_BangEqual:
			return !isEqual(left, right)
		}
	case *ast.VariableExpr:
		return storage.get(expr.Name)
	case *ast.AssignExpr:
		value := evalExpr(expr.Value)
		storage.assign(expr.Name, value)
		return value
	case *ast.LogicalExpr:
		leftValue := evalExpr(expr.Left)
		if expr.Operator.TokenType == tok.TokenType_Or && isTruthy(leftValue) {
			return leftValue
		}
		if expr.Operator.TokenType == tok.TokenType_And && !isTruthy(leftValue) {
			return leftValue
		}
		return evalExpr(expr.Right)
	case *ast.CallExpr:
		callee := evalExpr(expr.Callee)
		var arguments []any
		for _, argument := range expr.Arguments {
			arguments = append(arguments, evalExpr(argument))
		}
		callable, ok := callee.(loxCallable)
		if !ok {
			errors.ReportToken(expr.Parenthesis, "Can only call functions.")
		}
		if len(arguments) != callable.arity() {
			errors.ReportToken(expr.Parenthesis, fmt.Sprintf("Expected %d arguments, got %d.", callable.arity(), len(arguments)))
		}
		return callable.call(arguments)
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
