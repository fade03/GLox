package resolver

import (
	"fmt"
	"testing"
)

func TestStack_Push(t *testing.T) {
	var stack = NewStack()
	stack.Push(make(map[string]bool))
	top := stack.Peek().(map[string]bool)

	fmt.Println(top["a"])
	fmt.Println(top["b"])
}
