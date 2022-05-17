package interpreter

import "LoxGo/scanner"

// Environment 用来管理变量名->值之间的映射
type Environment map[string]interface{}

func (e Environment) define(name *scanner.Token, value interface{}) {
	e[name.Lexeme] = value
}

func (e Environment) lookup(name *scanner.Token) interface{} {
	if value, exist := e[name.Lexeme]; exist {
		return value
	}

	panic(NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'."))
}

func (e Environment) assign(name *scanner.Token, value interface{}) {
	if _, exist := e[name.Lexeme]; !exist {
		panic(NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'."))
	}

	e[name.Lexeme] = value
}
