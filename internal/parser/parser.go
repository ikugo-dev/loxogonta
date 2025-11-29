package parser

import (
	"slices"

	"github.com/ikugo-dev/loxogonta/internal/errors"
	"github.com/ikugo-dev/loxogonta/internal/tokens"
)

type Parser struct {
	tokens  []tok.Token
	current int
}

func (p *Parser) Parse() {
}

func (p *Parser) match(tokenTypes ...tok.TokenType) bool {
	if slices.ContainsFunc(tokenTypes, p.check) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) check(tokenType tok.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() tok.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == tok.TokenType_Eof
}

func (p *Parser) peek() tok.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() tok.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType tok.TokenType) {
	if p.check(tokenType) {
		p.advance()
	} else {
		errors.Report(p.current, "", "Expect ')' after expression.")
	}
}
