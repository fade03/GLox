package resolver

import (
	"GLox/interpreter"
	"GLox/parser"
	"GLox/scanner/token"
)

var (
	currentClass    ClassType    = None
	currentCallable CallableType = None
)

// Parser -> Resolver -> Interpreter

// Resolver implement ExprVisitor, StmtVisitor
type Resolver struct {
	interpreter *interpreter.Interpreter
	scopes      *Stack
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter, scopes: NewStack()}
}

func (r *Resolver) ResolveStmt(statements ...parser.Stmt) {
	for _, statement := range statements {
		_ = statement.Accept(r)
	}
}

func (r *Resolver) resolveExpr(expr parser.Expr) {
	_, _ = expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr parser.Expr, token *token.Token) {
	// 从栈顶向栈底搜索
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		if _, ok := r.scopes.items[i].(Scope)[token.Lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.Size()-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(stmt *parser.FuncDeclStmt, ct CallableType) {
	currentCallable = ct
	r.beginScope()
	for _, param := range stmt.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStmt(stmt.Body.Stmts...)
	r.endScope()
}
