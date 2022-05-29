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

func (r *Resolver) resolveMultiStmt(statements []parser.Stmt) {
	for _, statement := range statements {
		r.resolveSingleStmt(statement)
	}
}

func (r *Resolver) resolveSingleStmt(statement parser.Stmt) {
	statement.Accept(r)
}

func (r *Resolver) resolveExpr(expr parser.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveLocal(expr parser.Expr, token *scanner.Token) {
	for i, item := range r.scopes.items {
		if _, ok := item.(Scope)[token.Lexeme]; ok {
			// todo
			print(i)
			//r.interpreter.resolve(expr, r.scopes.Size()-1-i)
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
	r.resolveMultiStmt(stmt.Body.Stmts)
	r.endScope()
}
