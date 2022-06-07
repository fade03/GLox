package interpreter

import (
	le "GLox/loxerror"
	"GLox/parser"
	"GLox/scanner/token"
	"fmt"
)

func (i *Interpreter) VisitBinaryExpr(expr *parser.Binary) (interface{}, error) {
	// (递归)计算左右子表达式的值
	//lv, rv := i.evaluate(expr.Left), i.evaluate(expr.Right)
	lv, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	rv, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.MINUS:
		err = checkNumberOperands(expr.Operator, lv, rv)
		if err != nil {
			return nil, err
		}
		return lv.(float64) - rv.(float64), nil
	case token.STAR:
		err = checkNumberOperands(expr.Operator, lv, rv)
		if err != nil {
			return nil, err
		}
		return lv.(float64) * rv.(float64), nil
	case token.SLASH:
		err = checkNumberOperands(expr.Operator, lv, rv)
		if err != nil {
			return nil, err
		}
		return lv.(float64) / rv.(float64), nil
	// 加法操作可以定义在数字和字符之上
	case token.PLUS:
		return doPlus(expr.Operator, lv, rv)
	case token.GREATER:
		err = checkNumberOperands(expr.Operator, lv, rv)
		if err != nil {
			return nil, err
		}
		return lv.(float64) > rv.(float64), nil
	case token.GREATER_EQUAL:
		err = checkNumberOperands(expr.Operator, lv, rv)
		if err != nil {
			return nil, err
		}
		return lv.(float64) >= rv.(float64), nil
	case token.LESS:
		err = checkNumberOperands(expr.Operator, lv, rv)
		if err != nil {
			return nil, err
		}
		return lv.(float64) < rv.(float64), nil
	case token.LESS_EQUAL:
		err = checkNumberOperands(expr.Operator, lv, rv)
		if err != nil {
			return nil, err
		}
		return lv.(float64) <= rv.(float64), nil
	// == 和 != 运算的结果是bool类型
	case token.BANG_EQUAL:
		err = checkNumberOperands(expr.Operator, lv, rv)
		if err != nil {
			return nil, err
		}
		return !isEqual(lv, rv), nil
	case token.EQUAL_EQUAL:
		err = checkNumberOperands(expr.Operator, lv, rv)
		if err != nil {
			return nil, err
		}
		return isEqual(lv, rv), nil
	}
	return nil, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *parser.Grouping) (interface{}, error) {
	// 计算中间部分的expression即可
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *parser.Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *parser.Unary) (interface{}, error) {
	// 先计算右侧表达式的值，如果有运行时错误，直接返回
	rv, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.MINUS:
		err := checkNumberOperands(expr.Operator, rv)
		if err != nil {
			return nil, err
		}

		return -(rv.(float64)), nil
	case token.BANG:
		return !isTruth(rv), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitVariableExpr(expr *parser.Variable) (interface{}, error) {
	// return i.environment.lookup(expr.Name)
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) VisitAssignExpr(expr *parser.Assign) (interface{}, error) {
	// 计算Assign的语法树上的value节点
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	/*
		i.environment.assign(expr.Name, value)

		// 因为赋值也是一个表达式，所以这里返回所求的value
		return value
	*/
	if distance, ok := i.locals[expr]; ok {
		i.environment.assignAt(distance, expr.Name, value)
	} else {
		err := i.globals.assign(expr.Name, value)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitLogicExpr(expr *parser.Logic) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == token.OR {
		if isTruth(left) {
			return left, nil
		}
	} else {
		if !isTruth(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCallExpr(expr *parser.Call) (interface{}, error) {
	//callee, ok := i.evaluate(expr.Callee).(LoxCallable)
	calleeI, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	callee, ok := calleeI.(LoxCallable)
	if !ok {
		//panic(le.NewRuntimeError(expr.Paren, "Can only call functions and classes."))
		return nil, le.NewRuntimeError(expr.Paren, "Can only call functions and classes.")
	}

	var args []interface{}
	for _, arg := range expr.Arguments {
		value, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}

		args = append(args, value)
	}

	// 判断实参和形参的个数是否相同
	if len(args) != callee.Arity() {
		//panic(le.NewRuntimeError(expr.Paren, fmt.Sprintf("Expect %d arguments buf got %d.", len(args), callee.Arity())))
		return nil, le.NewRuntimeError(expr.Paren, fmt.Sprintf("Expect %d arguments buf got %d.", len(args), callee.Arity()))
	}

	return callee.Call(i, args)
}

func (i *Interpreter) VisitGetExpr(expr *parser.Get) (interface{}, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	// object必须是一个Instance
	instance, ok := object.(*LoxInstance)
	if !ok {
		//panic(le.NewRuntimeError(expr.Attribute, "Only instances have attributes."))
		return nil, le.NewRuntimeError(expr.Attribute, "Only instances have attributes.")
	}

	return instance.Get(expr.Attribute)
}

func (i *Interpreter) VisitSetExpr(expr *parser.Set) (interface{}, error) {
	// 计算等号左侧的表达式，找出要复制的属性
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	instance, ok := object.(*LoxInstance)
	if !ok {
		//panic(le.NewRuntimeError(expr.Attribute, "Only instances have attributes."))
		return nil, le.NewRuntimeError(expr.Attribute, "Only instances have attributes.")
	}

	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	instance.fields[expr.Attribute.Lexeme] = value

	return value, nil
}

func (i *Interpreter) VisitThisExpr(expr *parser.This) (interface{}, error) {
	return i.lookUpVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitSuperExpr(expr *parser.Super) (interface{}, error) {
	//superclass := i.lookUpVariable(expr.Keyword, expr).(*LoxClass)
	superclassI, err := i.lookUpVariable(expr.Keyword, expr)
	if err != nil {
		return nil, err
	}

	superclass := superclassI.(*LoxClass)
	instance := NewLoxInstance(superclass)
	method := superclass.findMethod(expr.Identifier.Lexeme)

	return method.bind(instance), nil
}

// ################### Statement #####################

func (i *Interpreter) VisitExprStmt(stmt *parser.ExprStmt) error {
	_, err := i.evaluate(stmt.Expr)

	return err
}

func (i *Interpreter) VisitFuncDeclStmt(stmt *parser.FuncDeclStmt) error {
	// 结束函数定义的区别在于，会创建一个保存了函数节点引用的新变量
	function := NewLoxFunction(stmt, i.environment, false)
	i.environment.define(stmt.Name, function)

	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *parser.ReturnStmt) (err error) {
	var value interface{}
	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return err
		}
	}

	panic(NewReturn(value))
}

func (i *Interpreter) VisitPrintStmt(stmt *parser.PrintStmt) error {
	value, err := i.evaluate(stmt.Expr)
	if err != nil {
		return err
	}

	// 需要打印计算的值
	fmt.Printf("%v\n", value)

	return nil
}

func (i *Interpreter) VisitVarDeclStmt(stmt *parser.VarDeclStmt) (err error) {
	var value interface{}
	if stmt.Initializer != nil {
		// 对变量的初始化语句求值
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return err
		}
	}
	i.environment.define(stmt.Name, value)

	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *parser.BlockStmt) error {
	// 把当前作用域的env传入下一个block
	return i.executeBlock(stmt, NewEnvironment(i.environment))
}

func (i *Interpreter) VisitIfStmt(stmt *parser.IfStmt) error {
	condition, err := i.evaluate(stmt.Condition)
	if err != nil {
		return err
	}

	if isTruth(condition) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}

	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *parser.WhileStmt) error {
	condition, err := i.evaluate(stmt.Condition)
	if err != nil {
		return err
	}

	for isTruth(condition) {
		return i.execute(stmt.Body)
	}

	return nil
}

func (i *Interpreter) VisitClassDeclStmt(stmt *parser.ClassDeclStmt) error {
	var superclass *LoxClass
	if stmt.Superclass != nil {
		tempV, err := i.evaluate(stmt.Superclass)
		if err != nil {
			return err
		}

		if tempC, ok := tempV.(*LoxClass); !ok {
			//panic(le.NewRuntimeError(stmt.Superclass.Name, "Superclass must be a class."))
			return le.NewRuntimeError(stmt.Superclass.Name, "Superclass must be a class.")
		} else {
			superclass = tempC
		}
	}

	i.environment.define(stmt.Name, nil)
	// "super"的作用域位于methods的上层
	if superclass != nil {
		i.environment = NewEnvironment(i.environment)
		i.environment.defineLiteral("super", superclass)
	}
	// methods
	var methods = make(map[string]*LoxFunction)
	for _, method := range stmt.Methods {
		methods[method.Name.Lexeme] = NewLoxFunction(method, i.environment, method.Name.Lexeme == "init")
	}

	class := NewLoxClass(stmt.Name.Lexeme, superclass, methods)
	if superclass != nil {
		// 切换回原来的scoop
		i.environment = i.environment.enclosing
	}

	return i.environment.assign(stmt.Name, class)
}
