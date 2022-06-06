package resolver

import (
	le "GLox/loxerror"
	"GLox/parser"
	"GLox/utils"
)

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

	r.resolveFunction(stmt, Function)
}

func (r *Resolver) VisitReturnStmt(stmt *parser.ReturnStmt) {
	if stmt.Value != nil {
		if currentCallable == Initializer {
			panic(le.NewRuntimeError(stmt.Keyword, "Can't return a value from initializer."))
		}
		r.resolveExpr(stmt.Value)
	}
}

func (r *Resolver) VisitClassDeclStmt(stmt *parser.ClassDeclStmt) {
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
			panic(le.NewRuntimeError(stmt.Superclass.Name, "A class can't inherit from itself."))
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
}
