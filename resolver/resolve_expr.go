package resolver

import (
	le "GLox/loxerror"
	"GLox/parser"
)

func (r *Resolver) VisitBinaryExpr(expr *parser.Binary) (interface{}, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)

	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr *parser.Grouping) (interface{}, error) {
	r.resolveExpr(expr.Expression)

	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr *parser.Literal) (interface{}, error) {
	// empty implementation

	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr *parser.Unary) (interface{}, error) {
	// empty implementation

	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr *parser.Variable) (interface{}, error) {
	//if prepared, ok := r.scopes.Peek().(Scope)[expr.Name.Lexeme]; !r.scopes.isEmpty() && ok && !prepared {
	//	lerror.ReportLexError(expr.Name.Line, expr.Name.Lexeme, "Can't read local variable in its own initializer.")
	//}

	if r.scopes.isEmpty() {
		return nil, nil
	}

	if prepared, ok := r.scopes.Peek().(Scope)[expr.Name.Lexeme]; ok && !prepared {
		//panic(le.NewRuntimeError(expr.Name, "Can't read local variable in its own initializer."))
		le.ReportResolveError(expr.Name, "Can't read local variable in its own initializer.")
	}

	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr *parser.Assign) (interface{}, error) {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitLogicExpr(expr *parser.Logic) (interface{}, error) {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)

	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr *parser.Call) (interface{}, error) {
	r.resolveExpr(expr.Callee)
	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}

	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr *parser.Get) (interface{}, error) {
	r.resolveExpr(expr.Object)

	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr *parser.Set) (interface{}, error) {
	r.resolveExpr(expr.Object)
	r.resolveExpr(expr.Value)

	return nil, nil
}

// VisitThisExpr : if "this" does not appear in a method, report an error.
func (r *Resolver) VisitThisExpr(expr *parser.This) (interface{}, error) {
	if !(currentClass == InClass || currentClass == SubClass) {
		//panic(le.NewRuntimeError(expr.Keyword, "Can't use 'this' outside of a class."))
		le.ReportResolveError(expr.Keyword, "Can't use 'this' outside of a class.")
	}

	return nil, nil
}

func (r *Resolver) VisitSuperExpr(expr *parser.Super) (interface{}, error) {
	if currentClass == None {
		le.ReportResolveError(expr.Keyword, "Can't use 'super' outside of a class.")
	} else if currentClass != SubClass {
		le.ReportResolveError(expr.Keyword, "Can't use 'super' in a class without superclass.")
	}

	r.resolveLocal(expr, expr.Keyword)

	return nil, nil
}
