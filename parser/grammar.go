package parser

import (
	"GLox/loxerror"
	"GLox/scanner/token"
	"GLox/utils"
)

// declaration -> varDecl | funcDecl | classDecl | statement
func (p *Parser) declaration() (Stmt, error) {
	//defer func() {
	//	if err := recover(); err != nil {
	//		// 当解释器出现错误的时候，进行同步，让解释器跳转到下一个语句或者声明的开头
	//		p.synchronize()
	//	}
	//}()
	// 如果可以匹配到 var 关键字，则为varDecl
	if p.match(token.VAR) {
		return p.varDecl()
	}

	// funcDecl -> "fun" function
	// 同 varDecl, 也可以看做是 statement 的一部分
	if p.match(token.FUN) {
		return p.functionDecl("function")
	}

	if p.match(token.CLASS) {
		return p.classDecl()
	}

	return p.statement()
}

// varDecl -> "var" IDENTIFIER ( "=" expression )? ";"
// varDecl 本身也可以看作是statement的一部分
func (p *Parser) varDecl() (Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	// initializer是可选的
	var initializer Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	// 注意最后consume掉一个分号
	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")

	return NewVarDeclStmt(name, initializer), err
}

// functionDecl -> IDENTIFIER "(" parameters? ")" block
// 将 functionDecl 单独抽离出来，可以在定义方法的时候复用这一条规则
// @param kind: "functionDecl" or "method"
func (p *Parser) functionDecl(kind string) (Stmt, error) {
	// 获取函数/方法名
	name, err := p.consume(token.IDENTIFIER, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}

	var parameters []*token.Token
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				panic(loxerror.NewParseError(p.peek(), "Can't have more than 255 parameters."))
			}
			// 获取参数名，Lox是动态类型，没有类型声明
			para, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, para)
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	// consume掉 ")"
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	// consume掉 "{"，一个函数体（block）的开始
	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	stmts, err := p.block()
	if err != nil {
		return nil, err
	}

	body := NewBlockStmt(stmts)

	return NewFunctionStmt(name, parameters, body), nil
}

// classDecl -> "class" IDENTIFIER ( "<" IDENTIFIER )? "{" function* "}" ;
func (p *Parser) classDecl() (Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}

	var superclass *Variable
	if p.match(token.LESS) {
		identifier, err := p.consume(token.IDENTIFIER, "Expect superclass name after '<'.")
		if err != nil {
			return nil, err
		}

		superclass = NewVariable(identifier)
	}

	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	var methods []*FuncDeclStmt
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		method, err := p.functionDecl("method")
		if err != nil {
			return nil, err
		}

		methods = append(methods, method.(*FuncDeclStmt))
	}

	_, err = p.consume(token.RIGHT_BRACE, "Expected '}' after class body.")

	return NewClassDeclStmt(name, superclass, methods), err
}

// statement -> exprStmt | printStmt | block | ifStmt | whileStmt | forStmt ｜ returnStmt
func (p *Parser) statement() (Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStmt()
	}

	if p.match(token.LEFT_BRACE) {
		stmts, err := p.block()

		return NewBlockStmt(stmts), err
	}

	if p.match(token.IF) {
		return p.ifStmt()
	}

	if p.match(token.WHILE) {
		return p.whileStmt()
	}

	if p.match(token.FOR) {
		return p.forStmt()
	}

	if p.match(token.RETURN) {
		return p.returnStmt()
	}

	return p.exprStmt()
}

// exprStmt -> expression ";"
func (p *Parser) exprStmt() (Stmt, error) {
	// 解析expression的值
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	// 如果下一个Token是';'，则consume掉，否则就出现了语法错误
	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")

	return NewExprStmt(value), err
}

// printStmt -> "print" expression ";"
func (p *Parser) printStmt() (Stmt, error) {
	// "print"已经在statement()中consume掉了（用于区分stmt的类型），所以这里不需要再match一遍
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	// consume ';'
	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")

	return NewPrintStmt(value), err
}

// block -> "{" + declaration* + "}"
func (p *Parser) block() (stmts []Stmt, err error) {
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)
	}
	_, err = p.consume(token.RIGHT_BRACE, "Expect '}' after block.")

	return stmts, err
}

// ifStmt -> "if" "(" expression ")" statement ( "else" statement )?
func (p *Parser) ifStmt() (Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after 'if condition'.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return NewIfStmt(condition, thenBranch, elseBranch), nil
}

// whileStmt -> "while" "(" expression ")" statement
func (p *Parser) whileStmt() (Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after 'while condition'.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()

	return NewWhileStmt(condition, body), err
}

// forStmt 由语法糖实现
func (p *Parser) forStmt() (Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer Stmt
	// for循环的初始化部分
	if p.match(token.SEMICOLON) {
		// 如果匹配到 ; 号，说明for循环的初始化被省略
		initializer = nil
	} else if p.match(token.VAR) {
		// 如果匹配到 var ，则为initializer
		initializer, err = p.varDecl()
		if err != nil {
			return nil, err
		}
	} else {
		// 如果没有匹配到var，则一定是一个表达式
		initializer, err = p.exprStmt()
		if err != nil {
			return nil, err
		}
	}
	// for循环的条件表达式
	var condition Expr
	// 如果没有匹配到 ; 号，则说明条件表达式没有被省略
	if !p.check(token.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	// for循环的增量语句
	var increment Expr
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	// for循环的主体
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

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

	return body, nil
}

// returnStmt -> "return" (expression)? ;
func (p *Parser) returnStmt() (Stmt, error) {
	keyword := p.previous()
	var value Expr
	var err error
	if !p.check(token.SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after return value.")

	return NewReturnStmt(keyword, value), err
}

// expression -> assignment
func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

// assignment -> ( call "." )? IDENTIFIER "=" assignment | logicOr
func (p *Parser) assignment() (Expr, error) {
	// 赋值表达式 = 号左侧其实是一个"伪表达式"，是一个经过计算可以赋值的"东西"，所以这里要先对左侧进行求值
	// expr的计算结果可能是logicOr或者优先级比LogicOr更高的表达式，主要包括**getter表达式**和**primary**
	expr, err := p.logicOr()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		// 只有左侧计算的结果是Variable的时候，才会进行赋值语句
		if ve, ok := expr.(*Variable); ok {
			// =号左侧表达式是一个Variable
			return NewAssign(ve.Name, value), nil
		} else if getter, ok := expr.(*Get); ok {
			// =号左侧表达式是一个Getter，则返回Setter表达式
			return NewSet(getter.Object, getter.Attribute, value), nil
		}

		// panic(loxerror.NewParseError(equals, "Invalid assignment target."))
		return nil, loxerror.NewParseError(equals, "Invalid assignment target.")
	}

	// 如果右侧没有初始化表达式，那么相当于是一个logicOr表达式
	return expr, nil
}

// logicOr -> logicAnd ( "or" logicAnd )*
func (p *Parser) logicOr() (Expr, error) {
	expr, err := p.logicAnd()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.logicAnd()
		if err != nil {
			return nil, err
		}

		expr = NewLogic(expr, operator, right)
	}

	return expr, nil
}

// logicAnd -> equality ( "and" equality )*
func (p *Parser) logicAnd() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = NewLogic(expr, operator, right)
	}

	return expr, nil
}

// equality -> comparison ( ( "!=" | "==" ) comparison )*
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		// match如果匹配，则会将current+1，之前的Token已经被consume掉了，所以下一行取的是之前的一个Token
		// 后面同理
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		// 递归解析二元表达式左边的节点
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

// comparison -> term ( (">" | ">=" | "<" | "<=") term )*
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

// term -> factor ( ("-" | "+") factor )*
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

// factor -> unary ( ("*" | "/") unary )*
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.STAR, token.SLASH) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

// unary -> ("!" | "-") unary | call
func (p *Parser) unary() (Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		return NewUnary(operator, right), err
	}

	return p.call()
}

// call -> primary ( "(" arguments? ")" | "." IDENTIFIER )*
// 函数调用的优先级仅次于 primary,
// 函数调用本身也可以是callee，如 funcall()()()，从文法角度上说就是 IDENTIFIER + ( "(" arguments? ")" )*
// 一个 argument 本身就是一个 expression, 所以不需要再重新定义它的文法，只需要在解析函数调用的同时解析函数参数即可,
// 属性访问（foo.bar）和函数调用具有相同的优先级
func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) {
			var arguments []Expr
			// 当前Token如果不是 ")"，则说明有参数
			if !p.check(token.RIGHT_PAREN) {
				for {
					// 限制最大参数量为255
					if len(arguments) >= 255 {
						//panic(loxerror.NewParseError(p.peek(), "Can't have more than 255 arguments."))
						return nil, loxerror.NewParseError(p.peek(), "Can't have more than 255 arguments.")
					}
					// 添加参数
					argument, err := p.expression()
					if err != nil {
						return nil, err
					}

					arguments = append(arguments, argument)
					// 参数之间要以 "," 隔开
					if !p.match(token.COMMA) {
						break
					}
				}
			}
			// consume掉 ")"
			paren, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
			if err != nil {
				return nil, err
			}

			// 不断迭代expr
			expr = NewCall(expr, paren, arguments)
		} else if p.match(token.DOT) {
			attribute, err := p.consume(token.IDENTIFIER, "Expect attribute name after '.'.")
			if err != nil {
				return nil, err
			}

			// 还是不断迭代expr
			expr = NewGet(expr, attribute)
		} else {
			// 如果 "(" 和 "." 都匹配不到，直接break，说明是一个primary
			break
		}
	}

	return expr, nil
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "return" | "(" expression ")" ｜ IDENTIFIER | "this" | super "." IDENTIFIER
// #### "super" isn't allowed to appear alone ###
func (p *Parser) primary() (Expr, error) {
	if p.match(token.TRUE) {
		return NewLiteral(true), nil
	}
	if p.match(token.FALSE) {
		return NewLiteral(false), nil
	}
	if p.match(token.NIL) {
		return NewLiteral(nil), nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return NewLiteral(p.previous().Literal), nil
	}

	if p.match(token.IDENTIFIER) {
		return NewVariable(p.previous()), nil
	}

	if p.match(token.THIS) {
		return NewThis(p.previous()), nil
	}

	if p.match(token.SUPER) {
		keyword := p.previous()
		_, err := p.consume(token.DOT, "Expect '.' after 'super'.")
		if err != nil {
			return nil, err
		}

		identifier, err := p.consume(token.IDENTIFIER, "Expect a identifier after '.'")

		return NewSuper(keyword, identifier), err
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")

		return NewGrouping(expr), err
	}

	//panic(loxerror.NewParseError(p.peek(), "Unknown expression."))
	return nil, loxerror.NewParseError(p.peek(), "Unknown expression.")
}
