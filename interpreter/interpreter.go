package interpreter

import (
	"LoxGo/parser"
	"LoxGo/scanner"
	"fmt"
)

// Interpreter ExprVisitor 和 StmtVisitor 子类之一，计算表达式的值
type Interpreter struct {
	environment Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{environment: make(map[string]interface{})}
}

func (i *Interpreter) VisitBinaryExpr(expr *parser.Binary) interface{} {
	// (递归)计算左右子表达式的值
	lv, rv := i.evaluate(expr.Left), i.evaluate(expr.Right)
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

func (i *Interpreter) VisitGroupingExpr(expr *parser.Grouping) interface{} {
	// 计算中间部分的expression即可
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *parser.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr *parser.Unary) interface{} {
	// 先计算右侧表达式的值
	rv := i.evaluate(expr.Right)
	switch expr.Operator.Type {
	case scanner.MINUS:
		checkNumberOperands(expr.Operator, rv)
		return -(rv.(float64))
	case scanner.BANG:
		return !isTruth(rv)
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *parser.Variable) interface{} {
	return i.environment.lookup(expr.Name)
}

func (i *Interpreter) VisitAssignExpr(expr *parser.Assign) interface{} {
	// 计算Assign的语法树上的value节点
	value := i.evaluate(expr.Value)
	i.environment.assign(expr.Name, value)

	// 因为赋值也是一个表达式，所以这里返回所求的value
	return value
}

func (i *Interpreter) VisitExprStmt(stmt *parser.ExprStmt) {
	i.evaluate(stmt.Expr)
}

func (i *Interpreter) VisitPrintStmt(stmt *parser.PrintStmt) {
	value := i.evaluate(stmt.Expr)
	// 需要打印计算的值
	fmt.Printf("%v", value)
}

func (i *Interpreter) VisitVarDeclStmt(stmt *parser.VarDeclStmt) {
	var value interface{}
	if stmt.Initializer != nil {
		// 对变量的初始化语句求值
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.define(stmt.Name, value)
}

// evaluate 计算表达式的值
func (i *Interpreter) evaluate(expr parser.Expr) interface{} {
	return expr.Accept(i)
}

// execute 执行一个statement
func (i *Interpreter) execute(stmt parser.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) Interpret(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		i.execute(stmt)
	}
}
