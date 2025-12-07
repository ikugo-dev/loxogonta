package prs

import (
	"github.com/ikugo-dev/loxogonta/internal/ast"
	"github.com/ikugo-dev/loxogonta/internal/errors"
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

// program        → declaration* EOF ;
// declaration    → varDecl | statement ;
// varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;

// statement      → exprStmt | printStmt ;
// printStmt      → "print" expression ";" ;
// exprStmt       → expression ";" ;

// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER;

func program() []ast.Statement {
	var statements []ast.Statement
	for !isAtEnd() {
		statements = append(statements, declaration())
	}
	return statements
}
func declaration() ast.Statement {
	defer func() {
		if err := recover(); err != nil {
			if errors.HadParseError {
				synchronize()
			} else {
				panic(err)
			}
		}
	}()
	if match(tok.TokenType_Var) {
		return varDecl()
	}
	return statement()
}
func varDecl() ast.Statement {
	name := consume(tok.TokenType_Identifier, "Expect variable name.")
	var initializer ast.Expression = nil
	if match(tok.TokenType_Equal) {
		initializer = expression()
	}
	consume(tok.TokenType_Semicolon, "Expect ';' after variable declaration.")
	return &ast.VarStmt{Name: name, Initializer: initializer}
}
func statement() ast.Statement {
	if match(tok.TokenType_Print) {
		return printStmt()
	}
	return expressionStmt()
}
func printStmt() ast.Statement {
	value := expression()
	consume(tok.TokenType_Semicolon, "Expect ';' after value.")
	return &ast.PrintStmt{Expr: value}
}
func expressionStmt() ast.Statement {
	value := expression()
	consume(tok.TokenType_Semicolon, "Expect ';' after value.")
	return &ast.ExpressionStmt{Expr: value}
}

func expression() ast.Expression {
	return equality()
}
func equality() ast.Expression {
	lExpr := comparison()
	for match(tok.TokenType_BangEqual, tok.TokenType_EqualEqual) {
		operator := previous()
		rExpr := comparison()
		lExpr = &ast.Binary{Left: lExpr, Operator: operator, Right: rExpr}
	}
	return lExpr
}
func comparison() ast.Expression {
	lExpr := term()
	for match(tok.TokenType_Greater,
		tok.TokenType_GreaterEqual,
		tok.TokenType_Less,
		tok.TokenType_LessEqual) {

		operator := previous()
		rExpr := term()
		lExpr = &ast.Binary{Left: lExpr, Operator: operator, Right: rExpr}
	}
	return lExpr
}
func term() ast.Expression {
	lExpr := factor()
	for match(tok.TokenType_Plus, tok.TokenType_Minus) {
		operator := previous()
		rExpr := factor()
		lExpr = &ast.Binary{Left: lExpr, Operator: operator, Right: rExpr}
	}
	return lExpr
}
func factor() ast.Expression {
	lExpr := unary()
	for match(tok.TokenType_Slash, tok.TokenType_Star) {
		operator := previous()
		rExpr := unary()
		lExpr = &ast.Binary{Left: lExpr, Operator: operator, Right: rExpr}
	}
	return lExpr
}
func unary() ast.Expression {
	if match(tok.TokenType_Bang, tok.TokenType_Minus) {
		operator := previous()
		rExpr := unary()
		return &ast.Unary{Operator: operator, Right: rExpr}
	}
	return primary()
}
func primary() ast.Expression {
	literalTokenTypes := []tok.TokenType{
		tok.TokenType_Number,
		tok.TokenType_String,
		tok.TokenType_True,
		tok.TokenType_False,
		tok.TokenType_Nil,
	}
	for _, tokenType := range literalTokenTypes {
		if match(tokenType) {
			return &ast.Literal{Value: previous().Literal}
		}
	}
	if match(tok.TokenType_LeftParen) {
		expr := expression()
		consume(tok.TokenType_RightParen, "Expect ')' after expression.")
		return &ast.Grouping{Expression: expr}
	}
	if match(tok.TokenType_Identifier) {
		return &ast.Variable{Name: previous()}
	}
	errors.ReportToken(peek(), "Expect expression.")
	return nil
}
