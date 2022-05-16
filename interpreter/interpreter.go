package interpreter

import (
	"LoxGo/parser"
	"LoxGo/scanner"
)

// Interpreter Visitor 子类之一，计算表达式的值
type Interpreter struct {
}

func (i *Interpreter) VisitBinary(expr *parser.Binary) interface{} {
	// (递归)计算左右子表达式的值
	lv, rv := i.Interpret(expr.Left), i.Interpret(expr.Right)
	switch expr.Operator.Type {
	case scanner.MINUS:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) - rv.(float64)
	case scanner.STAR:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) * rv.(float64)
	case scanner.SLASH:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) / rv.(float64)
	// 加法操作可以定义在数字和字符之上
	case scanner.PLUS:
		return doPlus(expr.Operator, lv, rv)
	case scanner.GREATER:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) > rv.(float64)
	case scanner.GREATER_EQUAL:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) >= rv.(float64)
	case scanner.LESS:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) < rv.(float64)
	case scanner.LESS_EQUAL:
		checkNumberOperands(expr.Operator, lv, rv)
		return lv.(float64) <= rv.(float64)
	// != 和 == 运算的结果是bool类型
	case scanner.BANG_EQUAL:
		return !isEqual(lv, rv)
	case scanner.EQUAL_EQUAL:
		return isEqual(lv, rv)
	}
	return nil
}

func (i *Interpreter) VisitGrouping(expr *parser.Grouping) interface{} {
	// 计算中间部分的expression即可
	return i.Interpret(expr.Expression)
}

func (i *Interpreter) VisitLiteral(expr *parser.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnary(expr *parser.Unary) interface{} {
	// 先计算右侧表达式的值
	rv := i.Interpret(expr.Right)
	switch expr.Operator.Type {
	case scanner.MINUS:
		checkNumberOperands(expr.Operator, rv)
		return -(rv.(float64))
	case scanner.BANG:
		return !isTruth(rv)
	}

	return nil
}

func (i *Interpreter) Interpret(expr parser.Expr) interface{} {
	return expr.Accept(i)
}
