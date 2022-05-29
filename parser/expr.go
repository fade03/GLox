package parser

import "GLox/scanner"

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator *scanner.Token
	Right    Expr
}

func NewBinary(left Expr, operator *scanner.Token, right Expr) *Binary {
	return &Binary{left, operator, right}
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression}
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) Expr {
	return &Literal{value}
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

type Unary struct {
	Operator *scanner.Token
	Right    Expr
}

func NewUnary(operator *scanner.Token, right Expr) *Unary {
	return &Unary{operator, right}
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

// Variable 也是表达式的一部分
type Variable struct {
	Name *scanner.Token
}

func NewVariable(name *scanner.Token) *Variable {
	return &Variable{Name: name}
}

func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}

type Assign struct {
	Name  *scanner.Token
	Value Expr
}

func NewAssign(name *scanner.Token, value Expr) *Assign {
	return &Assign{Name: name, Value: value}
}

func (a *Assign) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(a)
}

type Logic struct {
	Left     Expr
	Operator *scanner.Token
	Right    Expr
}

func NewLogic(left Expr, operator *scanner.Token, right Expr) *Logic {
	return &Logic{Left: left, Operator: operator, Right: right}
}

func (l *Logic) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogicExpr(l)
}

type Call struct {
	Callee    Expr
	Paren     *scanner.Token
	Arguments []Expr
}

func NewCall(callee Expr, paren *scanner.Token, arguments []Expr) *Call {
	return &Call{Callee: callee, Paren: paren, Arguments: arguments}
}

func (c *Call) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitCallExpr(c)
}
