package parser

type Stmt interface {
	Accept(visitor StmtVisitor)
}

type ExprStmt struct {
	Expr Expr
}

func NewExprStmt(expr Expr) *ExprStmt {
	return &ExprStmt{expr}
}

func (e *ExprStmt) Accept(visitor StmtVisitor) {
	visitor.VisitExprStmt(e)
}

type PrintStmt struct {
	Expr Expr
}

func NewPrintStmt(expr Expr) *PrintStmt {
	return &PrintStmt{expr}
}

func (p *PrintStmt) Accept(visitor StmtVisitor) {
	visitor.VisitPrintStmt(p)
}
