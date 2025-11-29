package parser

import tok "github.com/ikugo-dev/loxogonta/internal/tokens"

func (p *Parser) expression() Expression {
	return p.equality()
}

func (p *Parser) equality() Expression {
	lExpr := p.comparison()
	for p.match(tok.TokenType_BangEqual, tok.TokenType_EqualEqual) {
		operator := p.previous()
		rExpr := p.comparison()
		lExpr = &Binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *Parser) comparison() Expression {
	lExpr := p.term()
	for p.match(tok.TokenType_Greater,
		tok.TokenType_GreaterEqual,
		tok.TokenType_Less,
		tok.TokenType_LessEqual) {

		operator := p.previous()
		rExpr := p.term()
		lExpr = &Binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *Parser) term() Expression {
	lExpr := p.factor()
	for p.match(tok.TokenType_Plus, tok.TokenType_Minus) {
		operator := p.previous()
		rExpr := p.factor()
		lExpr = &Binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *Parser) factor() Expression {
	lExpr := p.unary()
	for p.match(tok.TokenType_Slash, tok.TokenType_Star) {
		operator := p.previous()
		rExpr := p.unary()
		lExpr = &Binary{lExpr, operator, rExpr}
	}
	return lExpr
}
func (p *Parser) unary() Expression {
	if p.match(tok.TokenType_Bang, tok.TokenType_Minus) {
		operator := p.previous()
		rExpr := p.unary()
		return &Unary{operator, rExpr}
	}
	return p.primary()
}
func (p *Parser) primary() Expression {
	literalTokenTypes := []tok.TokenType{
		tok.TokenType_Number,
		tok.TokenType_String,
		tok.TokenType_True,
		tok.TokenType_False,
		tok.TokenType_Nil,
	}
	for _, tokenType := range literalTokenTypes {
		if p.match(tokenType) {
			return &Literal{p.previous().Literal}
		}
	}
	if p.match(tok.TokenType_LeftParen) {
		expr := p.expression()
		p.consume(tok.TokenType_RightParen)
		return &Grouping{expr}
	}
	// errors.Report(p.current, "Parser", "Failed to parse primary")
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
