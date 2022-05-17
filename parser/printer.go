package parser

import (
	"LoxGo/utils"
	"bytes"
)

// Printer ExprVisitor 子类之一，以特殊的形式打印出语法树上的节点
type Printer struct {
}

func (p *Printer) VisitBinaryExpr(expr *Binary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *Printer) VisitGroupingExpr(expr *Grouping) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p *Printer) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	// 打印value的字面量
	return utils.ToString(expr.Value)
}

func (p *Printer) VisitUnaryExpr(expr *Unary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *Printer) VisitVariableExpr(expr *Variable) interface{} {
	// TODO implement me
	panic("implement me")
}

func (p *Printer) VisitAssignExpr(expr *Assign) interface{} {
	//TODO implement me
	panic("implement me")
}

func (p *Printer) parenthesize(name string, exprs ...Expr) string {
	var buffer bytes.Buffer
	buffer.WriteString("(" + name)
	for _, expr := range exprs {
		buffer.WriteString(" ")
		buffer.WriteString(expr.Accept(p).(string))
	}
	buffer.WriteString(")")

	return buffer.String()
}

func (p *Printer) Print(expr Expr) string {
	return expr.Accept(p).(string)
}
