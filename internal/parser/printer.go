package parser

import (
	"GLox/utils"
	"bytes"
)

// Printer ExprVisitor 子类之一，以特殊的形式打印出语法树上的节点
type Printer struct {
}

func (p *Printer) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p *Printer) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return p.parenthesize("group", expr.Expression), nil
}

func (p *Printer) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	// 打印value的字面量
	return utils.ToString(expr.Value), nil
}

func (p *Printer) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (p *Printer) VisitVariableExpr(expr *Variable) (interface{}, error) {
	// empty implementation

	return nil, nil
}

func (p *Printer) VisitAssignExpr(expr *Assign) (interface{}, error) {
	// empty implementation

	return nil, nil
}

func (p *Printer) VisitLogicExpr(expr *Logic) (interface{}, error) {
	// empty implementation

	return nil, nil
}

func (p *Printer) VisitCallExpr(expr *Call) (interface{}, error) {
	// empty implementation

	return nil, nil
}

func (p *Printer) VisitGetExpr(expr *Get) (interface{}, error) {
	// empty implementation

	return nil, nil
}

func (p *Printer) VisitSetExpr(expr *Set) (interface{}, error) {
	// empty implementation

	return nil, nil
}

func (p *Printer) VisitThisExpr(expr *This) (interface{}, error) {
	// empty implementation

	return nil, nil
}

func (p *Printer) VisitSuperExpr(expr *Super) (interface{}, error) {
	// empty implementation

	return nil, nil
}

func (p *Printer) parenthesize(name string, exprs ...Expr) string {
	var buffer bytes.Buffer
	buffer.WriteString("(" + name)
	for _, expr := range exprs {
		buffer.WriteString(" ")
		output, _ := expr.Accept(p)
		buffer.WriteString(output.(string))
	}
	buffer.WriteString(")")

	return buffer.String()
}

func (p *Printer) Print(expr Expr) string {
	output, _ := expr.Accept(p)
	return output.(string)
}
