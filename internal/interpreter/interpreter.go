package interpreter

import (
	parser2 "GLox/internal/parser"
	"GLox/internal/scanner/token"
	"time"
)

// Interpreter ExprVisitor 和 StmtVisitor 子类之一，计算表达式的值
type Interpreter struct {
	environment *Environment
	globals     *Environment // globals 存放的是可以全局使用的native函数
	locals      map[parser2.Expr]int
}

func NewInterpreter() *Interpreter {
	g := NewEnvironment(nil)
	g.defineLiteral("clock", NewLoxCallableImpl(func(interpreter *Interpreter, arguments []interface{}) interface{} {
		return time.Now().Unix()
	}, 0))

	return &Interpreter{
		//environment: NewEnvironment(nil),
		environment: NewEnvironment(g),
		globals:     g,
		locals:      make(map[parser2.Expr]int),
	}
}

// semantic.go

// evaluate 计算表达式的值
func (i *Interpreter) evaluate(expr parser2.Expr) (interface{}, error) {
	return expr.Accept(i)
}

// execute 执行一个statement
func (i *Interpreter) execute(stmt parser2.Stmt) error {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(block *parser2.BlockStmt, env *Environment) error {
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
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) Resolve(expr parser2.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) lookUpVariable(token *token.Token, expr parser2.Expr) (interface{}, error) {
	// 现在本地变量表中查询
	if distance, ok := i.locals[expr]; ok {
		return i.environment.getAt(distance, token.Lexeme), nil
	}

	//return i.globals.lookup(token)
	return i.environment.lookup(token)
}

func (i *Interpreter) Interpret(stmts []parser2.Stmt) error {
	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}
