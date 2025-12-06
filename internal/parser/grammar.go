package prs

import (
	"github.com/ikugo-dev/loxogonta/internal/ast"
	"github.com/ikugo-dev/loxogonta/internal/errors"
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

// program        → statement* EOF ;
// statement      → exprStmt | printStmt ;
// printStmt      → "print" expression ";" ;
// exprStmt       → expression ";" ;

// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;

func (p *parser) program() []ast.Statement {
	var statements []ast.Statement
	for !p.isAtEnd() {
		statements = append(statements, p.statement())
	}
	return statements
}
func (p *parser) statement() ast.Statement {
	if p.match(tok.TokenType_Print) {
		return p.printStmt()
	}
	return p.expressionStmt()
}
func (p *parser) printStmt() ast.Statement {
	value := p.expression()
	p.consume(tok.TokenType_Semicolon, "Expect ';' after value.")
	return &ast.PrintStmt{Expr: value}
}
func (p *parser) expressionStmt() ast.Statement {
	value := p.expression()
	p.consume(tok.TokenType_Semicolon, "Expect ';' after value.")
	return &ast.ExpressionStmt{Expr: value}
}

func (p *parser) expression() ast.Expression {
	return p.equality()
}
func (p *parser) equality() ast.Expression {
	lExpr := p.comparison()
	for p.match(tok.TokenType_BangEqual, tok.TokenType_EqualEqual) {
		operator := p.previous()
		rExpr := p.comparison()
		lExpr = &ast.Binary{Left: lExpr, Operator: operator, Right: rExpr}
	}
	return lExpr
}
func (p *parser) comparison() ast.Expression {
	lExpr := p.term()
	for p.match(tok.TokenType_Greater,
		tok.TokenType_GreaterEqual,
		tok.TokenType_Less,
		tok.TokenType_LessEqual) {

		operator := p.previous()
		rExpr := p.term()
		lExpr = &ast.Binary{Left: lExpr, Operator: operator, Right: rExpr}
	}
	return lExpr
}
func (p *parser) term() ast.Expression {
	lExpr := p.factor()
	for p.match(tok.TokenType_Plus, tok.TokenType_Minus) {
		operator := p.previous()
		rExpr := p.factor()
		lExpr = &ast.Binary{Left: lExpr, Operator: operator, Right: rExpr}
	}
	return lExpr
}
func (p *parser) factor() ast.Expression {
	lExpr := p.unary()
	for p.match(tok.TokenType_Slash, tok.TokenType_Star) {
		operator := p.previous()
		rExpr := p.unary()
		lExpr = &ast.Binary{Left: lExpr, Operator: operator, Right: rExpr}
	}
	return lExpr
}
func (p *parser) unary() ast.Expression {
	if p.match(tok.TokenType_Bang, tok.TokenType_Minus) {
		operator := p.previous()
		rExpr := p.unary()
		return &ast.Unary{Operator: operator, Right: rExpr}
	}
	return p.primary()
}
func (p *parser) primary() ast.Expression {
	literalTokenTypes := []tok.TokenType{
		tok.TokenType_Number,
		tok.TokenType_String,
		tok.TokenType_True,
		tok.TokenType_False,
		tok.TokenType_Nil,
	}
	for _, tokenType := range literalTokenTypes {
		if p.match(tokenType) {
			return &ast.Literal{Value: p.previous().Literal}
		}
	}
	if p.match(tok.TokenType_LeftParen) {
		expr := p.expression()
		p.consume(tok.TokenType_RightParen, "Expect ')' after expression.")
		return &ast.Grouping{Expression: expr}
	}
	errors.ReportToken(p.peek(), "Expect expression.")
	return nil
}
