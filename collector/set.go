package collector

type Set[T any] interface {
	Collection[T]
}

type setImpl[T comparable] struct {
	data map[T]struct{}
}

func (s *setImpl[T]) AsSlice() []*T {
	ret := make([]*T, len(s.data))
	i := 0
	for k := range s.data {
		ret[i] = &k
		i++
	}
	return ret
}

func (s *setImpl[T]) Length() int {
	return len(s.data)
}

func (s *setImpl[T]) Remove(elements ...T) int {
	changed := 0
	for _, e := range elements {
		if s.contains(e) {
			delete(s.data, e)
			changed++
		}
	}
	return changed
}

func (s *setImpl[T]) Add(elements ...T) int {
	changed := 0
	for _, e := range elements {
		if !s.contains(e) {
			s.data[e] = struct{}{}
			changed++
		}
	}
	return changed
}

func (s *setImpl[T]) Contains(element T) bool {
	return s.contains(element)
}

func (s *setImpl[T]) contains(element T) bool {
	_, ok := s.data[element]
	return ok
}

func NewSet[T comparable]() Set[T] {
	return &setImpl[T]{
		data: make(map[T]struct{}),
	}
}
