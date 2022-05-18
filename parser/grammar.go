package parser

import (
	"LoxGo/scanner"
)

// declaration -> varDecl | statement
func (p *Parser) declaration() Stmt {
	defer func() {
		if err := recover(); err != nil {
			// 当解释器出现错误的时候，进行同步，让解释器跳转到下一个语句或者声明的开头
			p.synchronize()
		}
	}()
	// 如果可以匹配到 var 关键字，则为varDecl
	if p.match(scanner.VAR) {
		return p.varDecl()
	}
	return p.statement()
}

// varDecl -> "var" IDENTIFIER ( "=" expression )? ";"
// varDecl 本身也可以看作是statement的一部分
func (p *Parser) varDecl() Stmt {
	name := p.consume(scanner.IDENTIFIER, "Expected variable name.")
	// initializer是可选的
	var initializer Expr
	if p.match(scanner.EQUAL) {
		initializer = p.expression()
	}
	// 注意最后consume掉一个分号
	p.consume(scanner.SEMICOLON, "Expected ';' after variable declaration.")
	return NewVarDeclStmt(name, initializer)
}

// statement -> exprStmt | printStmt | block
func (p *Parser) statement() Stmt {
	if p.match(scanner.PRINT) {
		return p.printStmt()
	}

	if p.match(scanner.LEFT_BRACE) {
		return NewBlockStmt(p.block())
	}

	return p.exprStmt()
}

// exprStmt -> expression ";"
func (p *Parser) exprStmt() Stmt {
	// 解析expression的值
	value := p.expression()
	// 如果下一个Token是';'，则consume掉，否则就出现了语法错误
	p.consume(scanner.SEMICOLON, "Expected ';' after value.")

	return NewExprStmt(value)
}

// printStmt -> "print" expression ";"
func (p *Parser) printStmt() Stmt {
	// "print"已经在statement()中consume掉了（用于区分stmt的类型），所以这里不需要再match一遍
	value := p.expression()
	// consume ';'
	p.consume(scanner.SEMICOLON, "Expected ';' after value.")

	return NewPrintStmt(value)
}

// block -> "{" + declaration* + "}"
func (p *Parser) block() (stmts []Stmt) {
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	p.consume(scanner.RIGHT_BRACE, "Expect '}' after block.")

	return stmts
}

// expression -> assignment
func (p *Parser) expression() Expr {
	return p.assignment()
}

// assignment -> IDENTIFIER "=" assignment | equality
// 某种意义上来说 assignment -> equality "=" assignment | equality
func (p *Parser) assignment() Expr {
	// 赋值表达式 = 号左侧其实是一个"伪表达式"，是一个经过计算可以赋值的"东西"，所以这里要先对左侧进行求值
	expr := p.equality()
	if p.match(scanner.EQUAL) {
		equals := p.previous()
		value := p.assignment()
		// 只有左侧计算的结果是Variable的时候，才会进行赋值语句
		if ve, ok := expr.(*Variable); ok {
			return NewAssign(ve.Name, value)
		}

		panic(newParseError(equals, "Invalid assignment target."))
	}

	return expr
}

// equality -> comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		// match如果匹配，则会将current+1，之前的Token已经被consume掉了，所以下一行取的是之前的一个Token
		// 后面同理
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

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ｜ IDENTIFIER
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

	if p.match(scanner.IDENTIFIER) {
		return NewVariable(p.previous())
	}

	if p.match(scanner.LEFT_PAREN) {
		expr := p.expression()
		p.consume(scanner.RIGHT_PAREN, "Expected ')' after expression.")
		return NewGrouping(expr)
	}

	// 如果不匹配任何的primary文法，则panic
	panic(newParseError(p.peek(), "Expect expression."))
}
