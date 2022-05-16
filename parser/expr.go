package parser

import "LoxGo/scanner"

type Expr interface {
	Accept(visitor Visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator *scanner.Token
	Right    Expr
}

func NewBinary(left Expr, operator *scanner.Token, right Expr) *Binary {
	return &Binary{left, operator, right}
}

func (b *Binary) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression}
}

func (g *Grouping) Accept(visitor Visitor) interface{} {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{value}
}

func (l *Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteral(l)
}

type Unary struct {
	Operator *scanner.Token
	Right    Expr
}

func NewUnary(operator *scanner.Token, right Expr) *Unary {
	return &Unary{operator, right}
}

func (u *Unary) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnary(u)
}
