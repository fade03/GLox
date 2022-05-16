package parser

type Visitor interface {
	VisitBinary(expr *Binary) interface{}
	VisitGrouping(expr *Grouping) interface{}
	VisitLiteral(expr *Literal) interface{}
	VisitUnary(expr *Unary) interface{}
}
