package collector

import (
	"sync"
)

type Stack[T any] interface {
	Push(elements ...*T)
	Pop() *T
	AsSlice() []*T
	Get(index int) *T
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
	if index >= length || length == 0 {
		return nil
	}
	return s.data[index]
}

func (s *stack[T]) Length() int {
	return len(s.data)
}

func NewStack[T interface{}]() Stack[T] {
	return &stack[T]{
		lock: sync.Mutex{},
		data: make([]*T, 0),
	}
}
