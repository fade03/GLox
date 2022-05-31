package interpreter

import (
	"GLox/parser"
	"GLox/scanner"
	"time"
)

// Interpreter ExprVisitor 和 StmtVisitor 子类之一，计算表达式的值
type Interpreter struct {
	environment *Environment
	globals     *Environment // globals 存放的是可以全局使用的native函数
	locals      map[parser.Expr]int
}

func NewInterpreter() *Interpreter {
	g := NewEnvironment(nil)
	g.defineStr("clock", NewLoxCallableImpl(func(interpreter *Interpreter, arguments []interface{}) interface{} {
		return time.Now().Unix()
	}, 0))

	return &Interpreter{
		//environment: NewEnvironment(nil),
		environment: NewEnvironment(g),
		globals:     g,
		locals:      make(map[parser.Expr]int),
	}
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

func (i *Interpreter) Resolve(expr parser.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) lookUpVariable(token *scanner.Token, expr parser.Expr) interface{} {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.getAt(distance, token.Lexeme)
	}

	//return i.globals.lookup(token)
	return i.environment.lookup(token)
}

func (i *Interpreter) Interpret(stmts []parser.Stmt) {
	for _, stmt := range stmts {
		i.execute(stmt)
	}
}
