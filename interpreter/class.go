package interpreter

type LoxClass struct {
	name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{name: name}
}

func (lc *LoxClass) String() string {
	return "<class " + lc.name + ">"
}