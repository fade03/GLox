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
	r.ResolveStmt(stmt.Stmts...)
	r.endScope()
}

func (r *Resolver) VisitIfStmt(stmt *parser.IfStmt) {
	r.resolveExpr(stmt.Condition)
	r.ResolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.ResolveStmt(stmt.ElseBranch)
	}
}

func (r *Resolver) VisitWhileStmt(stmt *parser.WhileStmt) {
	r.resolveExpr(stmt.Condition)
	r.ResolveStmt(stmt.Body)
}

func (r *Resolver) VisitFuncDeclStmt(stmt *parser.FuncDeclStmt) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt)
}

func (r *Resolver) VisitReturnStmt(stmt *parser.ReturnStmt) {
	r.resolveExpr(stmt.Value)
}

func (r *Resolver) VisitClassDeclStmt(stmt *parser.ClassDeclStmt) {
	// Lox允许将一个类声明为局部变量
	r.declare(stmt.Name)
	r.define(stmt.Name)
}
