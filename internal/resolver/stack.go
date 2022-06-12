package resolver

import "fmt"

type Stack struct {
	items []interface{}
}

func NewStack() *Stack {
	return &Stack{items: []interface{}{}}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (item interface{}) {
	if len(s.items) == 0 {
		return nil
	}

	item = s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]

	return item
}

func (s *Stack) Peek() interface{} {
	return s.items[len(s.items)-1]
}

func (s *Stack) Get(index int) interface{} {
	return s.items[index]
}

func (s *Stack) Size() int {
	return len(s.items)
}

func (s *Stack) isEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) String() string {
	return fmt.Sprintf("%v\n", s.items)
}
