package parser

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitVariableExpr(expr *Variable) interface{}
	VisitAssignExpr(expr *Assign) interface{}
	VisitLogicExpr(expr *Logic) interface{}
	VisitCallExpr(expr *Call) interface{}
	VisitGetExpr(expr *Get) interface{}
	VisitSetExpr(expr *Set) interface{}
}

// StmtVisitor 中定义的方法相当于直接执行语句，所以不会有返回值
type StmtVisitor interface {
	VisitExprStmt(stmt *ExprStmt)
	VisitPrintStmt(stmt *PrintStmt)
	VisitVarDeclStmt(stmt *VarDeclStmt)
	VisitBlockStmt(stmt *BlockStmt)
	VisitIfStmt(stmt *IfStmt)
	VisitWhileStmt(stmt *WhileStmt)
	VisitFuncDeclStmt(stmt *FuncDeclStmt)
	VisitReturnStmt(stmt *ReturnStmt)
	VisitClassDeclStmt(stmt *ClassDeclStmt)
}
