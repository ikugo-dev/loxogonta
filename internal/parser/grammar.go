package prs

import (
	"github.com/ikugo-dev/loxogonta/internal/ast"
	"github.com/ikugo-dev/loxogonta/internal/errors"
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

// program        → declaration* EOF ;
// declaration    → varDecl | statement ;
// varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;

// statement      → block | printStmt | expressionStmt | ifStmt | whileStmt | forStmt
// block          → "{" declaration* "}" ;
// printStmt      → "print" expression ";" ;
// expressionStmt → expression ";" ;
// ifStmt         → "if" "(" expression ")" statement ( "else" statement )? ;
// whileStmt      → "while" "(" expression ")" statement ;
// forStmt        → "for" "(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ;

// expression     → assignment ;
// assignment     → IDENTIFIER "=" assignment | logic_or ;
// logic_or       → logic_and ( "or" logic_and )* ;
// logic_and      → equality ( "and" equality )* ;
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
	consume(tok.TokenType_Identifier, "Expect variable name.")
	name := previous()
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
	if match(tok.TokenType_LeftBrace) {
		return &ast.BlockStmt{Statements: block()}
	}
	if match(tok.TokenType_If) {
		return ifStmt()
	}
	if match(tok.TokenType_While) {
		return whileStmt()
	}
	if match(tok.TokenType_For) {
		return forStmt()
	}
	return expressionStmt()
}
func block() []ast.Statement {
	var statements []ast.Statement
	for !check(tok.TokenType_RightBrace) && !isAtEnd() {
		statements = append(statements, declaration())
	}
	consume(tok.TokenType_RightBrace, "Expect '}' after block.")
	return statements
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
func ifStmt() ast.Statement {
	consume(tok.TokenType_LeftParen, "Expect '(' after 'if'.")
	condition := expression()
	consume(tok.TokenType_RightParen, "Expect ')' after if condition.")
	thenBranch := statement()
	var elseBranch ast.Expression = nil
	if match(tok.TokenType_Else) {
		elseBranch = statement()
	}
	return &ast.IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}
func whileStmt() ast.Expression {
	consume(tok.TokenType_LeftParen, "Expect '(' after 'while'.")
	condition := expression()
	consume(tok.TokenType_RightParen, "Expect ')' after while condition.")
	body := statement()
	return &ast.WhileStmt{Condition: condition, Body: body}
}
func forStmt() ast.Expression { // desugaring
	consume(tok.TokenType_LeftParen, "Expect '(' after 'for'.")
	var initializer ast.Statement
	if match(tok.TokenType_Semicolon) {
		initializer = nil
	} else if match(tok.TokenType_Var) {
		initializer = varDecl()
	} else {
		initializer = expressionStmt()
	}

	var condition ast.Expression = nil
	if !check(tok.TokenType_Semicolon) {
		condition = expression()
	}
	consume(tok.TokenType_Semicolon, "Expect ';' after loop condition.")

	var increment ast.Expression = nil
	if !check(tok.TokenType_RightParen) {
		increment = expression()
	}
	consume(tok.TokenType_RightParen, "Expect ')' after for clauses.")

	var body ast.Statement = statement()

	if increment != nil {
		body = &ast.BlockStmt{Statements: []ast.Statement{body, &ast.ExpressionStmt{Expr: increment}}}
	}
	if condition == nil {
		condition = &ast.Literal{Value: true}
	}
	body = &ast.WhileStmt{Condition: condition, Body: body}
	if initializer != nil {
		body = &ast.BlockStmt{Statements: []ast.Statement{initializer, body}}
	}

	return body
}
func expression() ast.Expression {
	return assignment()
}
func assignment() ast.Expression {
	lExpr := logicalOr()
	if match(tok.TokenType_Equal) {
		equals := previous()
		rExpr := assignment()                        // assignment is right associative
		if lvalue, ok := lExpr.(*ast.Variable); ok { // turn expression into lvalue
			return &ast.Assign{Name: lvalue.Name, Value: rExpr}
		}
		errors.ReportToken(equals, "Invalid assignment target.")
	}
	return lExpr
}
func logicalOr() ast.Expression {
	lExpr := logicalAnd()
	for match(tok.TokenType_Or) {
		operator := previous()
		rExpr := logicalAnd()
		lExpr = &ast.Logical{Left: lExpr, Operator: operator, Right: rExpr}
	}
	return lExpr
}
func logicalAnd() ast.Expression {
	lExpr := equality()
	for match(tok.TokenType_And) {
		operator := previous()
		rExpr := equality()
		lExpr = &ast.Logical{Left: lExpr, Operator: operator, Right: rExpr}
	}
	return lExpr
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
