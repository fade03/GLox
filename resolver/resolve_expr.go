package resolver

import (
	"GLox/lerror"
	"GLox/parser"
)

func (r *Resolver) VisitBinaryExpr(expr *parser.Binary) interface{} {
	//TODO implement me
	panic("implement me")
}

func (r *Resolver) VisitGroupingExpr(expr *parser.Grouping) interface{} {
	//TODO implement me
	panic("implement me")
}

func (r *Resolver) VisitLiteralExpr(expr *parser.Literal) interface{} {
	//TODO implement me
	panic("implement me")
}

func (r *Resolver) VisitUnaryExpr(expr *parser.Unary) interface{} {
	//TODO implement me
	panic("implement me")
}

func (r *Resolver) VisitVariableExpr(expr *parser.Variable) interface{} {
	if !r.scopes.isEmpty() && r.scopes.Peek().(Scope)[expr.Name.Lexeme] == false {
		lerror.Report(expr.Name.Line, expr.Name.Lexeme, "Can't read local variable in its own initializer.")
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
	//TODO implement me
	panic("implement me")
}

func (r *Resolver) VisitCallExpr(call *parser.Call) interface{} {
	//TODO implement me
	panic("implement me")
}
