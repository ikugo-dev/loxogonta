package prs

import (
	"github.com/ikugo-dev/loxogonta/internal/errors"
	tok "github.com/ikugo-dev/loxogonta/internal/tokens"
)

// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

func (p *parser) expression() expression {
	return p.equality()
}

func (p *parser) equality() expression {
	lExpr := p.comparison()
	for p.match(tok.TokenType_BangEqual, tok.TokenType_EqualEqual) {
		operator := p.previous()
		rExpr := p.comparison()
		lExpr = &binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *parser) comparison() expression {
	lExpr := p.term()
	for p.match(tok.TokenType_Greater,
		tok.TokenType_GreaterEqual,
		tok.TokenType_Less,
		tok.TokenType_LessEqual) {

		operator := p.previous()
		rExpr := p.term()
		lExpr = &binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *parser) term() expression {
	lExpr := p.factor()
	for p.match(tok.TokenType_Plus, tok.TokenType_Minus) {
		operator := p.previous()
		rExpr := p.factor()
		lExpr = &binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *parser) factor() expression {
	lExpr := p.unary()
	for p.match(tok.TokenType_Slash, tok.TokenType_Star) {
		operator := p.previous()
		rExpr := p.unary()
		lExpr = &binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *parser) unary() expression {
	if p.match(tok.TokenType_Bang, tok.TokenType_Minus) {
		operator := p.previous()
		rExpr := p.unary()
		return &unary{operator, rExpr}
	}
	return p.primary()
}
func (p *parser) primary() expression {
	literalTokenTypes := []tok.TokenType{
		tok.TokenType_Number,
		tok.TokenType_String,
		tok.TokenType_True,
		tok.TokenType_False,
		tok.TokenType_Nil,
	}
	for _, tokenType := range literalTokenTypes {
		if p.match(tokenType) {
			return &literal{p.previous().Literal}
		}
	}
	if p.match(tok.TokenType_LeftParen) {
		expr := p.expression()
		p.consume(tok.TokenType_RightParen, "Expect ')' after expression.")
		return &grouping{expr}
	}
	errors.ReportToken(p.peek(), "Expect expression.")
	return nil
}
