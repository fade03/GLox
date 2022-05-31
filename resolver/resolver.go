package resolver

import (
	"GLox/interpreter"
	"GLox/parser"
	"GLox/scanner"
)

// Parser -> Resolver -> Interpreter

type Resolver struct {
	interpreter *interpreter.Interpreter
	scopes      *Stack
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter, scopes: NewStack()}
}

func (r *Resolver) ResolveStmt(statements ...parser.Stmt) {
	for _, statement := range statements {
		statement.Accept(r)
	}
}

func (r *Resolver) resolveExpr(expr parser.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr parser.Expr, token *scanner.Token) {
	// 从栈顶向栈底搜索
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		if _, ok := r.scopes.items[i].(Scope)[token.Lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.Size()-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(stmt *parser.FuncStmt) {
	r.beginScope()
	for _, param := range stmt.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStmt(stmt.Body.Stmts...)
	r.endScope()
}
