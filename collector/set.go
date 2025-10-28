package collector

import (
	"sync"
)

type threadSafeSet[T comparable] struct {
	lock sync.Mutex
	data map[T]struct{}
}

func (s *threadSafeSet[T]) AsSlice() []*T {
	s.lock.Lock()
	defer s.lock.Unlock()
	ret := make([]*T, len(s.data))
	i := 0
	for k := range s.data {
		ret[i] = &k
		i++
	}
	return ret
}

func (s *threadSafeSet[T]) Length() int {
	return len(s.data)
}

func (s *threadSafeSet[T]) Remove(elements ...T) int {
	s.lock.Lock()
	defer s.lock.Unlock()
	changed := 0
	for _, e := range elements {
		if s.contains(e) {
			delete(s.data, e)
			changed++
		}
	}
	return changed
}

func (s *threadSafeSet[T]) Add(elements ...T) int {
	s.lock.Lock()
	defer s.lock.Unlock()
	changed := 0
	for _, e := range elements {
		if !s.contains(e) {
			s.data[e] = struct{}{}
			changed++
		}
	}
	return changed
}

func (s *threadSafeSet[T]) contains(element T) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.data[element]
	return ok
}

func NewThreadSafeSet[T comparable]() Collection[T] {
	return &threadSafeSet[T]{
		lock: sync.Mutex{},
		data: make(map[T]struct{}),
	}
}
