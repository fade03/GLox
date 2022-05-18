package interpreter

import (
	"LoxGo/parser"
)

// Interpreter ExprVisitor 和 StmtVisitor 子类之一，计算表达式的值
type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{environment: NewEnvironment(nil)}
}

// semantic.go

// evaluate 计算表达式的值
func (i *Interpreter) evaluate(expr parser.Expr) interface{} {
	return expr.Accept(i)
}

// execute 执行一个statement
func (i *Interpreter) execute(stmt parser.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) executeBlock(block *parser.BlockStmt, env *Environment) {
	previous := i.environment
	// 如果execute方法出现异常，defer还会正常执行，之前的作用域会正常恢复
	defer func() {
		// block执行完毕后，恢复之前的作用域
		i.environment = previous
	}()

	i.environment = env
	for _, stmt := range block.Stmts {
		// 新的env替换当前的env
		// 解释执行block中的statement
		i.execute(stmt)
	}
}

func (i *Interpreter) Interpret(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		i.execute(stmt)
	}
}
