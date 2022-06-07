package parser

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
	VisitVariableExpr(expr *Variable) (interface{}, error)
	VisitAssignExpr(expr *Assign) (interface{}, error)
	VisitLogicExpr(expr *Logic) (interface{}, error)
	VisitCallExpr(expr *Call) (interface{}, error)
	VisitGetExpr(expr *Get) (interface{}, error)
	VisitSetExpr(expr *Set) (interface{}, error)
	VisitThisExpr(expr *This) (interface{}, error)
	VisitSuperExpr(expr *Super) (interface{}, error)
}

// StmtVisitor 中定义的方法相当于直接执行语句，所以不会有返回值
type StmtVisitor interface {
	VisitExprStmt(stmt *ExprStmt) error
	VisitPrintStmt(stmt *PrintStmt) error
	VisitVarDeclStmt(stmt *VarDeclStmt) error
	VisitBlockStmt(stmt *BlockStmt) error
	VisitIfStmt(stmt *IfStmt) error
	VisitWhileStmt(stmt *WhileStmt) error
	VisitFuncDeclStmt(stmt *FuncDeclStmt) error
	VisitReturnStmt(stmt *ReturnStmt) error
	VisitClassDeclStmt(stmt *ClassDeclStmt) error
}
