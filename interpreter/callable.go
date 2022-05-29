package interpreter

import "GLox/parser"

type LoxCallableFunc func(interpreter *Interpreter, arguments []interface{}) interface{}

// LoxCallable 任何可以被调用的对象都要实现这个接口，比如定义的函数、类中的方法。
type LoxCallable interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	Arity() int
}

// ################ Native ###################

type Native struct {
	fn LoxCallableFunc
	n  int
}

func NewLoxCallableImpl(fn LoxCallableFunc, n int) *Native {
	return &Native{fn: fn, n: n}
}

func (n *Native) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return n.fn(interpreter, arguments)
}

func (n *Native) Arity() int {
	return n.n
}

func (n *Native) String() string {
	return "<native fn>"
}

// ############### Function ###################

type LoxFunction struct {
	declaration *parser.FuncStmt
	closure     *Environment // 在函数体内部声明的函数
}

func NewLoxFunction(declaration *parser.FuncStmt, closure *Environment) *LoxFunction {
	return &LoxFunction{declaration: declaration, closure: closure}
}

func (lf *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (result interface{}) {
	// 捕获 return 语句
	defer func() {
		if r, ok := recover().(*Return); ok && r != nil {
			result = r.value
		}
	}()

	// 创建函数自己的作用域，enclosing就是全局的函数 ~globals~
	env := NewEnvironment(lf.closure)
	// 将形参和实参绑定起来
	for i, arg := range arguments {
		env.define(lf.declaration.Params[i], arg)
	}
	interpreter.executeBlock(lf.declaration.Body, env)

	return result
}

func (lf *LoxFunction) Arity() int {
	return len(lf.declaration.Params)
}

func (lf *LoxFunction) String() string {
	return "<fn " + lf.declaration.Name.Lexeme + ">"
}
