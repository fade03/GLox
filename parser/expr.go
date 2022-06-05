package parser

import (
	"GLox/scanner/token"
)

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func NewBinary(left Expr, operator *token.Token, right Expr) *Binary {
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
	Operator *token.Token
	Right    Expr
}

func NewUnary(operator *token.Token, right Expr) *Unary {
	return &Unary{operator, right}
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

// Variable 也是表达式的一部分
type Variable struct {
	Name *token.Token
}

func NewVariable(name *token.Token) *Variable {
	return &Variable{Name: name}
}

func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}

type Assign struct {
	Name  *token.Token
	Value Expr
}

func NewAssign(name *token.Token, value Expr) *Assign {
	return &Assign{Name: name, Value: value}
}

func (a *Assign) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(a)
}

type Logic struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func NewLogic(left Expr, operator *token.Token, right Expr) *Logic {
	return &Logic{Left: left, Operator: operator, Right: right}
}

func (l *Logic) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogicExpr(l)
}

type Call struct {
	Callee    Expr
	Paren     *token.Token
	Arguments []Expr
}

func NewCall(callee Expr, paren *token.Token, arguments []Expr) *Call {
	return &Call{Callee: callee, Paren: paren, Arguments: arguments}
}

func (c *Call) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitCallExpr(c)
}

type Get struct {
	Object    Expr
	Attribute *token.Token
}

func NewGet(object Expr, attribute *token.Token) *Get {
	return &Get{Object: object, Attribute: attribute}
}

func (g *Get) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGetExpr(g)
}

type Set struct {
	Object    Expr
	Attribute *token.Token
	Value     Expr
}

func NewSet(Object Expr, Attribute *token.Token, Value Expr) *Set {
	return &Set{Object: Object, Attribute: Attribute, Value: Value}
}

func (s *Set) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitSetExpr(s)
}

type This struct {
	Keyword *token.Token
}

func NewThis(keyword *token.Token) *This {
	return &This{Keyword: keyword}
}

func (t *This) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitThisExpr(t)
}
