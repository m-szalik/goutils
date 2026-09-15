package collector

import (
	"sync"
)

// Stack is a LIFO collection of pointers.
type Stack[T any] interface {
	// Push appends elements to the top of the stack.
	Push(elements ...*T)
	// Pop removes and returns the top element, or nil when empty.
	Pop() *T
	// AsSlice returns stack contents from bottom to top.
	AsSlice() []*T
	// Get returns element at index, or nil when out of bounds.
	Get(index int) *T
	// Length returns current stack size.
	Length() int
}

type stack[T interface{}] struct {
	lock sync.Mutex
	data []*T
}

func (s *stack[T]) Push(elements ...*T) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = append(s.data, elements...)
}

func (s *stack[T]) Pop() *T {
	s.lock.Lock()
	defer s.lock.Unlock()
	length := len(s.data)
	if length == 0 {
		return nil
	}
	result := s.data[length-1]
	s.data = s.data[:length-1]
	return result
}

func (s *stack[T]) AsSlice() []*T {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.data[:]
}

func (s *stack[T]) Get(index int) *T {
	s.lock.Lock()
	defer s.lock.Unlock()
	length := len(s.data)
	if index < 0 || index >= length {
		return nil
	}
	return s.data[index]
}

func (s *stack[T]) Length() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.data)
}

// NewStack creates an empty stack.
func NewStack[T interface{}]() Stack[T] {
	return &stack[T]{
		lock: sync.Mutex{},
		data: make([]*T, 0),
	}
}
