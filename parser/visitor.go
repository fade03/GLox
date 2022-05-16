package parser

type ExprVisitor interface {
	VisitBinary(expr *Binary) interface{}
	VisitGrouping(expr *Grouping) interface{}
	VisitLiteral(expr *Literal) interface{}
	VisitUnary(expr *Unary) interface{}
}

// StmtVisitor 中定义的方法相当于直接执行语句，所以不会有返回值
type StmtVisitor interface {
	VisitExprStmt(stmt *ExprStmt)
	VisitPrintStmt(stmt *PrintStmt)
}
