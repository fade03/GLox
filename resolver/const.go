package resolver

type CallableType int

type ClassType int

const (
	None = iota
	Function
	Initializer
	Method
	InClass 
)
