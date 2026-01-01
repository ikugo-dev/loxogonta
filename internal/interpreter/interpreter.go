package intr

import (
	"fmt"
	"maps"
	"math"

	"github.com/ikugo-dev/loxogonta/internal/ast"
	"github.com/ikugo-dev/loxogonta/internal/errors"
	"github.com/ikugo-dev/loxogonta/internal/repl"
	"github.com/ikugo-dev/loxogonta/internal/resolver"
	"github.com/ikugo-dev/loxogonta/internal/tokens"
)

var globals *environment
var storage *environment
var locals map[ast.Expression]int

func startInterpreter(statements []ast.Statement) {
	if globals == nil || locals == nil {
		globals = createEnvironment()
		storage = globals
		addNativeFunctions(globals)
	}
	newLocals := rslv.Resolve(statements)
	if errors.HadParseError { // TODO: change to resolver errors
		return
	}

	if locals == nil {
		locals = make(map[ast.Expression]int)
	}
	maps.Copy(locals, newLocals)
}

func Interpret(statements []ast.Statement) any {
	startInterpreter(statements)

	var value any
	for _, statement := range statements {
		value = evalStmt(storage, statement)
	}
	return value
}

func evalStmt(storage *environment, statement ast.Statement) (value any) {
	switch stmt := statement.(type) {
	case *ast.PrintStmt:
		fmt.Println(evalExpr(storage, stmt.Expr))
	case *ast.ExpressionStmt:
		return evalExpr(storage, stmt.Expr)
	case *ast.VarStmt:
		var value any = nil
		if stmt.Initializer != nil {
			value = evalExpr(storage, stmt.Initializer)
		}
		storage.put(stmt.Name.Lexeme, value)
		repl.AddCompletion(stmt.Name.Lexeme, repl.CompletionType_Variable, fmt.Sprintf("%v", value))
	case *ast.BlockStmt:
		var value any = nil
		oldStorage := storage
		storage = createEnvironmentWithParent(oldStorage)
		for _, statement := range stmt.Statements {
			value = evalStmt(storage, statement)
		}
		storage = oldStorage
		return value
	case *ast.IfStmt:
		if isTruthy(evalExpr(storage, stmt.Condition)) {
			return evalStmt(storage, stmt.ThenBranch)
		} else {
			return evalStmt(storage, stmt.ElseBranch)
		}
	case *ast.WhileStmt:
		for isTruthy(evalExpr(storage, stmt.Condition)) {
			evalStmt(storage, stmt.Body)
		}
	case *ast.FunctionStmt:
		function := &loxFunction{declaration: stmt, closure: storage}
		storage.put(stmt.Name.Lexeme, function)
		var paramList []string
		for _, param := range function.declaration.Params {
			paramList = append(paramList, param.Lexeme)
		}
		repl.AddCompletion(stmt.Name.Lexeme, repl.CompletionType_Function, fmt.Sprintf("%v", paramList))
	case *ast.ReturnStmt:
		var value any
		if stmt.Value != nil {
			value = evalExpr(storage, stmt.Value)
		}
		panic(Return{value})
	case *ast.ClassStmt:
		storage.put(stmt.Name.Lexeme, nil)
		class := createLoxClass(stmt.Name.Lexeme)
		storage.assign(stmt.Name, class)
	}
	return nil
}

func evalExpr(storage *environment, expression ast.Expression) (value any) {
	switch expr := expression.(type) {
	case *ast.LiteralExpr:
		return expr.Value
	case *ast.GroupingExpr:
		return evalExpr(storage, expr.Expression)
	case *ast.UnaryExpr:
		right := evalExpr(storage, expr.Right)
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
		left := evalExpr(storage, expr.Left)
		right := evalExpr(storage, expr.Right)
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
		// return storage.get(expr.Name)
		return lookUpVariable(storage, expr, expr.Name)
	case *ast.AssignExpr:
		value := evalExpr(storage, expr.Value)

		// storage.assign(expr.Name, value)
		distance, exists := locals[expr]
		if exists {
			assignAt(storage, expr.Name, value, distance)
		} else {
			globals.assign(expr.Name, value)
		}
		return value
	case *ast.LogicalExpr:
		leftValue := evalExpr(storage, expr.Left)
		if expr.Operator.TokenType == tok.TokenType_Or && isTruthy(leftValue) {
			return leftValue
		}
		if expr.Operator.TokenType == tok.TokenType_And && !isTruthy(leftValue) {
			return leftValue
		}
		return evalExpr(storage, expr.Right)
	case *ast.CallExpr:
		defer func() {
			if err := recover(); err != nil {
				returnValue, ok := err.(Return)
				if !ok {
					panic(err)
				}
				value = returnValue.value
			}
		}()
		callee := evalExpr(storage, expr.Callee)
		var arguments []any
		for _, argument := range expr.Arguments {
			arguments = append(arguments, evalExpr(storage, argument))
		}
		callable, ok := callee.(loxCallable)
		if !ok {
			errors.ReportToken(expr.Parenthesis, "Can only call functions.")
		}
		if len(arguments) != callable.arity() {
			errors.ReportToken(expr.Parenthesis, fmt.Sprintf("Expected %d arguments, got %d.", callable.arity(), len(arguments)))
		}
		return callable.call(storage, arguments)
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

type Return struct {
	value any
}
