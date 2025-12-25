package rslv

import (
	"slices"

	"github.com/ikugo-dev/loxogonta/internal/ast"
	"github.com/ikugo-dev/loxogonta/internal/errors"
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

var scopes []map[string]bool
var locals map[ast.Expression]int

func Resolve(statements []ast.Statement) map[ast.Expression]int {
	locals = make(map[ast.Expression]int)
	for _, statement := range statements {
		resolveStmt(statement)
	}
	return locals
}

func resolveStmt(statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.PrintStmt:
		resolveExpr(stmt.Expr)
	case *ast.ExpressionStmt:
		resolveExpr(stmt.Expr)
	case *ast.VarStmt:
		declare(stmt.Name)
		if stmt.Initializer != nil {
			resolveExpr(stmt.Initializer)
		}
		define(stmt.Name)
	case *ast.BlockStmt:
		beginScope()
		for _, statement := range stmt.Statements {
			resolveStmt(statement)
		}
		endScope()
	case *ast.IfStmt:
		resolveExpr(stmt.Condition)
		resolveStmt(stmt.ThenBranch)
		if stmt.ElseBranch != nil {
			resolveStmt(stmt.ElseBranch)
		}
	case *ast.WhileStmt:
		resolveExpr(stmt.Condition)
		resolveStmt(stmt.Body)
	case *ast.FunctionStmt:
		declare(stmt.Name) //XXX do we really need them back to back? not really... right?
		define(stmt.Name)
		resolveFunction(*stmt)
	case *ast.ReturnStmt:
		if stmt.Value != nil {
			resolveExpr(stmt.Value)
		}
	default:
	}
}

func resolveExpr(expression ast.Expression) {
	switch expr := expression.(type) {
	case *ast.LiteralExpr: // NOTHING
	case *ast.GroupingExpr:
		resolveExpr(expr.Expression)
	case *ast.UnaryExpr:
		resolveExpr(expr.Right)
	case *ast.BinaryExpr:
		resolveExpr(expr.Left)
		resolveExpr(expr.Right)
	case *ast.VariableExpr:
		if len(scopes) != 0 && scopes[len(scopes)-1][expr.Name.Lexeme] == false {
			errors.ReportToken(expr.Name, "Can't read local variable in its own initializer.")
		}
		resolveLocal(expr, expr.Name.Lexeme)
	case *ast.AssignExpr:
		resolveExpr(expr)
		resolveLocal(expr, expr.Name.Lexeme)
	case *ast.LogicalExpr:
		resolveExpr(expr.Left)
		resolveExpr(expr.Right)
	case *ast.CallExpr:
		resolveExpr(expr.Callee)
		for _, argument := range expr.Arguments {
			resolveExpr(argument)
		}
	default:
	}
}

func resolveLocal(expr ast.Expression, name string) {
	for i, scope := range slices.Backward(scopes) {
		if _, ok := scope[name]; ok {
			// interpreter.resolve(expr, scopes.size() - 1 - i);
			locals[expr] = len(scope) - 1 - i
			return
		}
	}
}

func resolveFunction(function ast.FunctionStmt) { // TODO: maybe change to pointer?
	beginScope()
	for _, token := range function.Params {
		declare(token) //XXX do we really need them back to back? not really... right?
		define(token)
	}
	for _, statement := range function.Body {
		resolveStmt(statement)
	}
	endScope()
}

func beginScope() {
	scopes = append(scopes, make(map[string]bool))
}
func endScope() {
	scopes = scopes[:len(scopes)-1]
}

func declare(token tok.Token) {
	if len(scopes) == 0 {
		return
	}
	scopes[len(scopes)-1][token.Lexeme] = false // different than default (nil)
}

func define(token tok.Token) {
	if len(scopes) == 0 {
		return
	}
	scopes[len(scopes)-1][token.Lexeme] = true
}
