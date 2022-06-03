package interpreter

// LoxClass implements LoxCallable, specific the constructor method.
// e.g. class Foo {}; print(Foo());
type LoxClass struct {
	name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{name: name}
}

func (lc *LoxClass) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return NewLoxInstance(lc)
}

func (lc *LoxClass) Arity() int {
	return 0
}

func (lc *LoxClass) String() string {
	return "<class " + lc.name + ">"
}