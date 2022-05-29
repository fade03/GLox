package resolver

import "GLox/parser"

func (r *Resolver) VisitExprStmt(stmt *parser.ExprStmt) {
	r.resolveExpr(stmt.Expr)
}

func (r *Resolver) VisitPrintStmt(stmt *parser.PrintStmt) {
	r.resolveExpr(stmt.Expr)
}

func (r *Resolver) VisitVarDeclStmt(stmt *parser.VarDeclStmt) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
}

func (r *Resolver) VisitBlockStmt(stmt *parser.BlockStmt) {
	r.beginScope()
	r.resolveMultiStmt(stmt.Stmts)
	r.endScope()
}

func (r *Resolver) VisitIfStmt(stmt *parser.IfStmt) {
	r.resolveExpr(stmt.Condition)
	r.resolveSingleStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveSingleStmt(stmt.ElseBranch)
	}
}

func (r *Resolver) VisitWhileStmt(stmt *parser.WhileStmt) {
	//TODO implement me
	panic("implement me")
}

func (r *Resolver) VisitFuncDeclStmt(stmt *parser.FuncStmt) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt)
}

func (r *Resolver) VisitReturnStmt(stmt *parser.ReturnStmt) {
	r.resolveExpr(stmt.Value)
}
