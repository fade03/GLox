package parser

import "LoxGo/scanner"

type Stmt interface {
	Accept(visitor StmtVisitor)
}

type ExprStmt struct {
	Expr Expr
}

func NewExprStmt(expr Expr) *ExprStmt {
	return &ExprStmt{expr}
}

func (e *ExprStmt) Accept(visitor StmtVisitor) {
	visitor.VisitExprStmt(e)
}

type FuncStmt struct {
	Name   *scanner.Token
	Params []*scanner.Token
	Body   *BlockStmt
}

func NewFunctionStmt(name *scanner.Token, params []*scanner.Token, body *BlockStmt) *FuncStmt {
	return &FuncStmt{Name: name, Params: params, Body: body}
}

func (f *FuncStmt) Accept(visitor StmtVisitor) {
	visitor.VisitFuncStmt(f)
}

type ReturnStmt struct {
	Keyword *scanner.Token
	Value   Expr
}

func NewReturnStmt(keyword *scanner.Token, value Expr) *ReturnStmt {
	return &ReturnStmt{Keyword: keyword, Value: value}
}

func (r *ReturnStmt) Accept(visit StmtVisitor) {
	visit.VisitReturnStmt(r)
}

type PrintStmt struct {
	Expr Expr
}

func NewPrintStmt(expr Expr) *PrintStmt {
	return &PrintStmt{expr}
}

func (p *PrintStmt) Accept(visitor StmtVisitor) {
	visitor.VisitPrintStmt(p)
}

type VarDeclStmt struct {
	Name        *scanner.Token
	Initializer Expr
}

func NewVarDeclStmt(name *scanner.Token, initializer Expr) *VarDeclStmt {
	return &VarDeclStmt{Name: name, Initializer: initializer}
}

func (v *VarDeclStmt) Accept(visitor StmtVisitor) {
	visitor.VisitVarDeclStmt(v)
}

// BlockStmt 一个Block由多个statement组成，包含变量声明、表达式、print语句等
type BlockStmt struct {
	Stmts []Stmt
}

func NewBlockStmt(stmts []Stmt) *BlockStmt {
	return &BlockStmt{Stmts: stmts}
}

func (b *BlockStmt) Accept(visitor StmtVisitor) {
	visitor.VisitBlockStmt(b)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIfStmt(condition Expr, trueBranch Stmt, elseBranch Stmt) *IfStmt {
	return &IfStmt{Condition: condition, ThenBranch: trueBranch, ElseBranch: elseBranch}
}

func (i *IfStmt) Accept(visitor StmtVisitor) {
	visitor.VisitIfStmt(i)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func NewWhileStmt(condition Expr, body Stmt) *WhileStmt {
	return &WhileStmt{Condition: condition, Body: body}
}

func (w *WhileStmt) Accept(visitor StmtVisitor) {
	visitor.VisitWhileStmt(w)
}
