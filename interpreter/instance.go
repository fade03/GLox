package interpreter

import (
	"GLox/scanner"
)

type LoxInstance struct {
	class *LoxClass
	fields map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class, fields: make(map[string]interface{})}
}

// Get will first look for fields defined in the class, then methods. 
func (ls *LoxInstance) Get(attribute *scanner.Token) interface{} {
	if field, ok := ls.fields[attribute.Lexeme]; ok {
		return field
	}

	if method, ok := ls.class.methods[attribute.Lexeme]; ok {
		return method.bind(ls)
	}

	panic(NewRuntimeError(attribute, "undefined attribute " + attribute.Lexeme + "."))
}

func (ls *LoxInstance) String() string {
	return "<" + ls.class.name + " instance>"
}

