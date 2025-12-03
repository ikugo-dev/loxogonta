package prs

import (
	"slices"

	"github.com/ikugo-dev/loxogonta/internal/errors"
	"github.com/ikugo-dev/loxogonta/internal/tokens"
)

type parser struct {
	tokens  []tok.Token
	current int
}

func NewParser(tokens []tok.Token) parser {
	return parser{tokens, 0}
}

func (p *parser) Parse() expression {
	return p.declaration()
}

func (p *parser) declaration() expression {
	defer func() {
		err := recover()
		if err != nil {
			if _, ok := err.(errors.ParseError); ok {
				// this is for later
				// p.synchronize()
			} else {
				panic(err)
			}
		}
	}()
	return p.expression()
}

func (p *parser) match(tokenTypes ...tok.TokenType) bool {
	if slices.ContainsFunc(tokenTypes, p.check) {
		p.advance()
		return true
	}
	return false
}

func (p *parser) check(tokenType tok.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *parser) advance() tok.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) isAtEnd() bool {
	return p.peek().TokenType == tok.TokenType_Eof
}

func (p *parser) peek() tok.Token {
	return p.tokens[p.current]
}

func (p *parser) previous() tok.Token {
	return p.tokens[p.current-1]
}

func (p *parser) consume(tokenType tok.TokenType, message string) tok.Token { // is return needed?
	if p.check(tokenType) {
		return p.advance()
	}
	errors.ReportToken(p.peek(), message)
	return tok.Token{}
}

func (p *parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().TokenType == tok.TokenType_Semicolon {
			return
		}
		switch p.peek().TokenType { // finding end
		case tok.TokenType_Class:
			return
		case tok.TokenType_Fun:
			return
		case tok.TokenType_Var:
			return
		case tok.TokenType_For:
			return
		case tok.TokenType_If:
			return
		case tok.TokenType_While:
			return
		case tok.TokenType_Print:
			return
		case tok.TokenType_Return:
			return
		}
		p.advance()
	}
}
