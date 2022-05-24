package interpreter

type Return struct {
	value interface{}
}

func NewReturn(value interface{}) *Return {
	return &Return{value: value}
}
