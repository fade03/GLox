package parser

import (
	"LoxGo/scanner"
	"LoxGo/utils"
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
	name := p.consume(scanner.IDENTIFIER, "Expect variable name.")
	// initializer是可选的
	var initializer Expr
	if p.match(scanner.EQUAL) {
		initializer = p.expression()
	}
	// 注意最后consume掉一个分号
	p.consume(scanner.SEMICOLON, "Expect ';' after variable declaration.")
	return NewVarDeclStmt(name, initializer)
}

// statement -> exprStmt | printStmt | block | ifStmt | whileStmt | forStmt
func (p *Parser) statement() Stmt {
	if p.match(scanner.PRINT) {
		return p.printStmt()
	}

	if p.match(scanner.LEFT_BRACE) {
		return NewBlockStmt(p.block())
	}

	if p.match(scanner.IF) {
		return p.ifStmt()
	}

	if p.match(scanner.WHILE) {
		return p.whileStmt()
	}

	if p.match(scanner.FOR) {
		return p.forStmt()
	}

	return p.exprStmt()
}

// exprStmt -> expression ";"
func (p *Parser) exprStmt() Stmt {
	// 解析expression的值
	value := p.expression()
	// 如果下一个Token是';'，则consume掉，否则就出现了语法错误
	p.consume(scanner.SEMICOLON, "Expect ';' after value.")

	return NewExprStmt(value)
}

// printStmt -> "print" expression ";"
func (p *Parser) printStmt() Stmt {
	// "print"已经在statement()中consume掉了（用于区分stmt的类型），所以这里不需要再match一遍
	value := p.expression()
	// consume ';'
	p.consume(scanner.SEMICOLON, "Expect ';' after value.")

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

// ifStmt -> "if" "(" expression ")" statement ( "else" statement )?
func (p *Parser) ifStmt() Stmt {
	p.consume(scanner.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(scanner.RIGHT_PAREN, "Expect ')' after 'if condition'.")

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(scanner.ELSE) {
		elseBranch = p.statement()
	}

	return NewIfStmt(condition, thenBranch, elseBranch)
}

// whileStmt -> "while" "(" expression ")" statement
func (p *Parser) whileStmt() Stmt {
	p.consume(scanner.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(scanner.RIGHT_PAREN, "Expect ')' after 'while condition'.")
	body := p.statement()

	return NewWhileStmt(condition, body)
}

// forStmt 由语法糖实现
func (p *Parser) forStmt() Stmt {
	p.consume(scanner.LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer Stmt
	// for循环的初始化部分
	if p.match(scanner.SEMICOLON) {
		// 如果匹配到 ; 号，说明for循环的初始化被省略
		initializer = nil
	} else if p.match(scanner.VAR) {
		// 如果匹配到 var ，则为initializer
		initializer = p.varDecl()
	} else {
		// 如果没有匹配到var，则一定是一个表达式
		initializer = p.exprStmt()
	}
	// for循环的条件表达式
	var condition Expr
	// 如果没有匹配到 ; 号，则说明条件表达式没有被省略
	if !p.check(scanner.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(scanner.SEMICOLON, "Expect ';' after loop condition.")
	// for循环的增量语句
	var increment Expr
	if !p.check(scanner.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(scanner.RIGHT_PAREN, "Expect ')' after for clauses.")
	// for循环的主体
	body := p.statement()
	// 开始语法脱糖
	if increment != nil {
		// 如果增量语句不为nil，则相当于每次while循环body执行完毕之后再执行increment，把两个合成为一个block
		body = NewBlockStmt([]Stmt{body, NewExprStmt(increment)})
	}

	condition = utils.Ternary(condition == nil, NewLiteral(true), condition)
	body = NewWhileStmt(condition, body)

	if initializer != nil {
		body = NewBlockStmt([]Stmt{initializer, body})
	}

	return body
}

// expression -> assignment
func (p *Parser) expression() Expr {
	return p.assignment()
}

// assignment -> IDENTIFIER "=" assignment | logicOr
// 某种意义上来说 assignment -> logicOr "=" assignment | logicOr
func (p *Parser) assignment() Expr {
	// 赋值表达式 = 号左侧其实是一个"伪表达式"，是一个经过计算可以赋值的"东西"，所以这里要先对左侧进行求值
	expr := p.logicOr()
	if p.match(scanner.EQUAL) {
		equals := p.previous()
		value := p.assignment()
		// 只有左侧计算的结果是Variable的时候，才会进行赋值语句
		if ve, ok := expr.(*Variable); ok {
			return NewAssign(ve.Name, value)
		}

		panic(newParseError(equals, "Invalid assignment target."))
	}

	// 如果右侧没有初始化表达式，那么相当于就是一个equality表达式
	return expr
}

// logicOr -> logicAnd ( "or" logicAnd )*
func (p *Parser) logicOr() Expr {
	expr := p.logicAnd()
	for p.match(scanner.OR) {
		operator := p.previous()
		right := p.logicAnd()
		expr = NewLogic(expr, operator, right)
	}

	return expr
}

// logicAnd -> equality ( "and" equality )*
func (p *Parser) logicAnd() Expr {
	expr := p.equality()
	for p.match(scanner.AND) {
		operator := p.previous()
		right := p.equality()
		expr = NewLogic(expr, operator, right)
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
		p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(expr)
	}

	// 如果不匹配任何的primary文法，则panic
	panic(newParseError(p.peek(), "Expect expression."))
}
