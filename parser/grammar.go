package parser

import (
	"LoxGo/scanner"
)

func (p *Parser) expression() Expr {
	return p.equality()
}

// equality -> comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		// 递归解析二元表达式左边的节点
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

// comparison -> term ( (">" | ">=" | "<" | "<=") term )*
func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

// term -> factor ( ("-" | "+") factor )*
func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

// factor -> unary ( ("*" | "/") unary )*
func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(scanner.STAR, scanner.SLASH) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

// unary -> ("!" | "-") unary | primary
func (p *Parser) unary() Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}

	return p.primary()
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")"
func (p *Parser) primary() Expr {
	if p.match(scanner.TRUE) {
		return NewLiteral(true)
	}
	if p.match(scanner.FALSE) {
		return NewLiteral(false)
	}
	if p.match(scanner.NIL) {
		return NewLiteral(nil)
	}

	if p.match(scanner.NUMBER, scanner.STRING) {
		return NewLiteral(p.previous().Literal)
	}

	if p.match(scanner.LEFT_PAREN) {
		expr := p.expression()
		p.consume(scanner.RIGHT_PAREN, "Expected ')' after expression.")
		return NewGrouping(expr)
	}

	// 如果不匹配任何的primary文法，则panic
	panic(newParseError(p.peek(), "Expect expression."))
}
