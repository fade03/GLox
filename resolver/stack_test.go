package resolver

import (
	"fmt"
	"testing"
)

var stack = NewStack()

func TestStack_Push(t *testing.T) {
	stack.Push(1)
	stack.Push(2)

	fmt.Println(stack.String())
}

func TestStack_Pop(t *testing.T) {
	stack.Pop()
	fmt.Println(stack.String())
}

func TestStack_Peek(t *testing.T) {
	fmt.Println(stack.Peek())
}
