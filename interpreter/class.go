package interpreter

// LoxClass implements LoxCallable, specific the constructor method.
type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name: name, superclass: superclass, methods: methods}
}

// Call means "constructor", e.g. class Foo {}; print(Foo());
func (lc *LoxClass) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	instance := NewLoxInstance(lc)
	// init() will be called when an instance is initialized
	if initializer, ok := lc.methods["init"]; ok {
		// 类中的方法首先要经过bind处理，为特殊变量this绑定值
		initializer.bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (lc *LoxClass) Arity() int {
	if initializer, ok := lc.methods["init"]; ok {
		return initializer.Arity()
	}

	return 0
}

func (lc *LoxClass) findMethod(name string) *LoxFunction {
	if method, ok := lc.methods[name]; ok {
		return method
	}

	if lc.superclass != nil {
		return lc.superclass.findMethod(name)
	}

	return nil
}

func (lc *LoxClass) String() string {
	if lc.superclass != nil {
		return "<class " + lc.name + " inherit " + lc.superclass.name + ">"
	}

	return "<class " + lc.name + ">"
}
