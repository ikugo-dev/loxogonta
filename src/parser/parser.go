package parser

import (
	"github.com/ikugo-dev/loxogonta/src/errors"
	"github.com/ikugo-dev/loxogonta/src/scanner"
)

type Parser struct {
	Tokens  []scanner.Token
	current int
}

func (p *Parser) Parse() {
}

func (p *Parser) expression() Expression {
	return p.equality()
}

func (p *Parser) equality() Expression {
	lExpr := p.comparison()
	for p.match(scanner.TokenType_BangEqual, scanner.TokenType_EqualEqual) {
		operator := p.previous()
		rExpr := p.comparison()
		lExpr = &Binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *Parser) comparison() Expression {
	lExpr := p.term()
	for p.match(scanner.TokenType_Greater,
		scanner.TokenType_GreaterEqual,
		scanner.TokenType_Less,
		scanner.TokenType_LessEqual) {

		operator := p.previous()
		rExpr := p.term()
		lExpr = &Binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *Parser) term() Expression {
	lExpr := p.factor()
	for p.match(scanner.TokenType_Plus, scanner.TokenType_Minus) {
		operator := p.previous()
		rExpr := p.factor()
		lExpr = &Binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *Parser) factor() Expression {
	lExpr := p.unary()
	for p.match(scanner.TokenType_Slash, scanner.TokenType_Star) {
		operator := p.previous()
		rExpr := p.unary()
		lExpr = &Binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *Parser) unary() Expression {
	if p.match(scanner.TokenType_Bang, scanner.TokenType_Minus) {
		operator := p.previous()
		rExpr := p.unary()
		return &Unary{operator, rExpr}
	}
	return p.primary()
}
func (p *Parser) primary() Expression {
	literalTokenTypes := []scanner.TokenType{
		scanner.TokenType_Number,
		scanner.TokenType_String,
		scanner.TokenType_True,
		scanner.TokenType_False,
		scanner.TokenType_Nil,
	}
	for _, tokenType := range literalTokenTypes {
		if p.match(tokenType) {
			return &Literal{p.previous().GetLiteral()}
		}
	}
	if p.match(scanner.TokenType_LeftParen) {
		expr := p.expression()
		p.consume(scanner.TokenType_RightParen, "Expect ')' after expression.")
		return &Grouping{expr}
	}
	errors.Report(p.current, "Parser", "Failed to parse primary")
	return nil
}

// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

func (p *Parser) match(tokenTypes ...scanner.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}
func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().GetTokenType() == tokenType
}
func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().GetTokenType() == scanner.TokenType_Eof
}
func (p *Parser) peek() scanner.Token {
	return p.Tokens[p.current]
}
func (p *Parser) previous() scanner.Token {
	return p.Tokens[p.current-1]
}
