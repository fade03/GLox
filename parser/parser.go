package parser

import (
	"GLox/loxerror"
	"GLox/scanner/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{tokens: tokens}
}

// match 逻辑上是OR的关系，只要匹配到current指向的Token和任意一个传入的Token匹配就会返回true，并且会将current+1
func (p *Parser) match(types ...token.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

// check 判断current指向的Token类型和传入的类型t是否匹配
func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)-1
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() *token.Token {
	if p.isAtEnd() {
		return p.tokens[len(p.tokens)-1]
	}
	p.current++
	return p.previous()
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

// 判断current指向的Token是不是传入的t，如果不是则panic，如果是则返回当前token，然后current+1
func (p *Parser) consume(t token.TokenType, msg string) (*token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	//panic(loxerror.NewParseError(p.peek(), msg))
	return nil, loxerror.NewParseError(p.peek(), msg)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}

// Parse 将一个程序（Token序列）解析成多个Stmt
func (p *Parser) Parse() (stmts []Stmt) {
	// 一个程序由多个declaration + EOF组成: program -> declaration* EOF
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			println(err.Error())
			p.synchronize()
		} else {
			stmts = append(stmts, stmt)
		}
	}

	return stmts
}
