package parser

import (
	"GLox/scanner"
	"GLox/utils"
)

// declaration -> varDecl | funcDecl | classDecl | statement
func (p *Parser) declaration() Stmt {
	//defer func() {
	//	if err := recover(); err != nil {
	//		// 当解释器出现错误的时候，进行同步，让解释器跳转到下一个语句或者声明的开头
	//		p.synchronize()
	//	}
	//}()
	// 如果可以匹配到 var 关键字，则为varDecl
	if p.match(scanner.VAR) {
		return p.varDecl()
	}

	// funcDecl -> "fun" function
	// 同 varDecl, 也可以看做是 statement 的一部分
	if p.match(scanner.FUN) {
		return p.functionDecl("function")
	}

	if p.match(scanner.CLASS) {
		return p.classDecl()
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

// functionDecl -> IDENTIFIER "(" parameters? ")" block
// 将 functionDecl 单独抽离出来，可以在定义方法的时候复用这一条规则
// @param kind: "functionDecl" or "method"
func (p *Parser) functionDecl(kind string) Stmt {
	// 获取函数/方法名
	name := p.consume(scanner.IDENTIFIER, "Expect"+kind+" name.")
	p.consume(scanner.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	var parameters []*scanner.Token
	if !p.check(scanner.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				panic(NewParseError(p.peek(), "Can't have more than 255 parameters."))
			}
			// 获取参数名，Lox是动态类型，没有类型声明
			parameters = append(parameters, p.consume(scanner.IDENTIFIER, "Expect parameter name."))
			if !p.match(scanner.COMMA) {
				break
			}
		}
	}
	// consume掉 ")"
	p.consume(scanner.RIGHT_PAREN, "Expect ')' after parameters.")
	// consume掉 "{"，一个函数体（block）的开始
	p.consume(scanner.LEFT_BRACE, "Expect '{' before "+kind+" body.")

	body := NewBlockStmt(p.block())

	return NewFunctionStmt(name, parameters, body)
}

// classDecl -> "class" IDENTIFIER "{" function* "}" ;
func (p *Parser) classDecl() Stmt {
	name := p.consume(scanner.IDENTIFIER, "Expect class name.")
	p.consume(scanner.LEFT_BRACE, "Expect '{' before class body.")

	var methods []*FuncDeclStmt
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.functionDecl("method").(*FuncDeclStmt))
	}

	p.consume(scanner.RIGHT_BRACE, "Expected '}' after class body.")

	return NewClassDeclStmt(name, methods)
}

// statement -> exprStmt | printStmt | block | ifStmt | whileStmt | forStmt ｜ returnStmt
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

	if p.match(scanner.RETURN) {
		// todo
		return p.returnStmt()
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

// returnStmt -> "return" (expression)? ;
func (p *Parser) returnStmt() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(scanner.SEMICOLON) {
		value = p.expression()
	}
	p.consume(scanner.SEMICOLON, "Expect ';' after return value.")

	return NewReturnStmt(keyword, value)
}

// expression -> assignment
func (p *Parser) expression() Expr {
	return p.assignment()
}

// assignment -> ( call "." )? IDENTIFIER "=" assignment | logicOr
func (p *Parser) assignment() Expr {
	// 赋值表达式 = 号左侧其实是一个"伪表达式"，是一个经过计算可以赋值的"东西"，所以这里要先对左侧进行求值
	// expr的计算结果可能是logicOr或者优先级比LogicOr更高的表达式，主要包括**getter表达式**和**primary**
	expr := p.logicOr()
	if p.match(scanner.EQUAL) {
		equals := p.previous()
		value := p.assignment()
		// 只有左侧计算的结果是Variable的时候，才会进行赋值语句
		if ve, ok := expr.(*Variable); ok {
			// =号左侧表达式是一个Variable
			return NewAssign(ve.Name, value)
		} else if getter, ok := expr.(*Get); ok {
			// =号左侧表达式是一个Getter，则返回Setter表达式
			return NewSet(getter.Object, getter.Attribute, value)
		}

		panic(NewParseError(equals, "Invalid assignment target."))
	}

	// 如果右侧没有初始化表达式，那么相当于是一个logicOr表达式
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

// unary -> ("!" | "-") unary | call
func (p *Parser) unary() Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}

	return p.call()
}

// call -> primary ( "(" arguments? ")" | "." IDENTIFIER )*
// 函数调用的优先级仅次于 primary,
// 函数调用本身也可以是callee，如 funcall()()()，从文法角度上说就是 IDENTIFIER + ( "(" arguments? ")" )*
// 一个 argument 本身就是一个 expression, 所以不需要再重新定义它的文法，只需要在解析函数调用的同时解析函数参数即可,
// 属性访问（foo.bar）和函数调用具有相同的优先级
func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(scanner.LEFT_PAREN) {
			var arguments []Expr
			// 当前Token如果不是 ")"，则说明有参数
			if !p.check(scanner.RIGHT_PAREN) {
				for {
					// 限制最大参数量为255
					if len(arguments) >= 255 {
						panic(NewParseError(p.peek(), "Can't have more than 255 arguments."))
					}
					// 添加参数
					arguments = append(arguments, p.expression())
					// 参数之间要以 "," 隔开
					if !p.match(scanner.COMMA) {
						break
					}
				}
			}
			// consume掉 ")"
			paren := p.consume(scanner.RIGHT_PAREN, "Expect ')' after arguments.")
			// 不断迭代expr
			expr = NewCall(expr, paren, arguments)
		} else if p.match(scanner.DOT) {
			attribute := p.consume(scanner.IDENTIFIER, "Expect attribute name after '.'.")
			// 还是不断迭代expr
			expr = NewGet(expr, attribute)
		} else {
			// 如果 "(" 和 "." 都匹配不到，直接break，说明是一个primary
			break
		}
	}

	return expr
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "return" | "(" expression ")" ｜ IDENTIFIER
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
	panic(NewParseError(p.peek(), "Expect expression."))
}
