package resolver

import (
	"GLox/interpreter"
	"GLox/parser"
)

func (r *Resolver) VisitBinaryExpr(expr *parser.Binary) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)

	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *parser.Grouping) interface{} {
	r.resolveExpr(expr.Expression)

	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *parser.Literal) interface{} {
	// empty implementation

	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *parser.Unary) interface{} {
	// empty implementation

	return nil
}

func (r *Resolver) VisitVariableExpr(expr *parser.Variable) interface{} {
	//if prepared, ok := r.scopes.Peek().(Scope)[expr.Name.Lexeme]; !r.scopes.isEmpty() && ok && !prepared {
	//	lerror.Report(expr.Name.Line, expr.Name.Lexeme, "Can't read local variable in its own initializer.")
	//}

	if r.scopes.isEmpty() {
		return nil
	}

	if prepared, ok := r.scopes.Peek().(Scope)[expr.Name.Lexeme]; ok && !prepared {
		// lerror.Report(expr.Name.Line, expr.Name.Lexeme, "Can't read local variable in its own initializer.")
		panic(interpreter.NewRuntimeError(expr.Name, "Can't read local variable in its own initializer."))
		// return nil
	}

	r.resolveLocal(expr, expr.Name)

	return nil
}

func (r *Resolver) VisitAssignExpr(expr *parser.Assign) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)

	return nil
}

func (r *Resolver) VisitLogicExpr(expr *parser.Logic) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)

	return nil
}

func (r *Resolver) VisitCallExpr(call *parser.Call) interface{} {
	r.resolveExpr(call.Callee)
	for _, arg := range call.Arguments {
		r.resolveExpr(arg)
	}

	return nil
}

func (r *Resolver) VisitGetExpr(expr *parser.Get) interface {} {
	r.resolveExpr(expr.Object)

	return nil 
}

func (r *Resolver) VisitSetExpr(expr *parser.Set) interface{} {
	r.resolveExpr(expr.Object)
	r.resolveExpr(expr.Value)

	return nil 
}
