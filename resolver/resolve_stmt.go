package resolver

import (
	le "GLox/loxerror"
	"GLox/parser"
	"GLox/utils"
)

func (r *Resolver) VisitExprStmt(stmt *parser.ExprStmt) error {
	r.resolveExpr(stmt.Expr)
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *parser.PrintStmt) error {
	r.resolveExpr(stmt.Expr)
	return nil
}

func (r *Resolver) VisitVarDeclStmt(stmt *parser.VarDeclStmt) error {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitBlockStmt(stmt *parser.BlockStmt) error {
	r.beginScope()
	r.ResolveStmt(stmt.Stmts...)
	r.endScope()
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *parser.IfStmt) error {
	r.resolveExpr(stmt.Condition)
	r.ResolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.ResolveStmt(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *parser.WhileStmt) error {
	r.resolveExpr(stmt.Condition)
	r.ResolveStmt(stmt.Body)
	return nil
}

func (r *Resolver) VisitFuncDeclStmt(stmt *parser.FuncDeclStmt) error {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, Function)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *parser.ReturnStmt) error {
	if stmt.Value != nil {
		if currentCallable == Initializer {
			le.ReportResolveError(stmt.Keyword, "Can't return a value from initializer.")
		}
		r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitClassDeclStmt(stmt *parser.ClassDeclStmt) error {
	var enclosingClass = currentClass
	currentClass = InClass

	// Lox允许将一个类声明为局部变量
	r.declare(stmt.Name)
	r.define(stmt.Name)
	//if stmt.Superclass != nil && stmt.Superclass.Name.Lexeme == stmt.Name.Lexeme {
	//	panic(le.NewRuntimeError(stmt.Superclass.Name, "A class can't inherit from itself."))
	//}
	if stmt.Superclass != nil {
		if stmt.Superclass.Name.Lexeme == stmt.Name.Lexeme {
			le.ReportResolveError(stmt.Superclass.Name, "A class can't inherit from itself.")
		} else {
			currentClass = SubClass
			r.resolveExpr(stmt.Superclass)
			r.beginScope()
			r.scopes.Peek().(Scope)["super"] = true // "super"的作用域位于"this"的上层
		}
	}

	// 处理特殊的变量"this"，为它创建一个单独的作用域，位于类中方法的上层
	r.beginScope()
	r.scopes.Peek().(Scope)["this"] = true
	// resolve类中的方法
	for _, method := range stmt.Methods {
		callableType := utils.Ternary[CallableType](method.Name.Lexeme == "init", Initializer, Method)
		r.resolveFunction(method, callableType)
	}
	r.endScope() // 对应81行开启的"this"的作用域

	if stmt.Superclass != nil {
		r.endScope() // 对应75行开启的"super"的作用域
	}

	currentClass = enclosingClass
	return nil
}
