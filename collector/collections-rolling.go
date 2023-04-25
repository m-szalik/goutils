package collector

import "sync"

type rollingCollection[T comparable] struct {
	lock  sync.Mutex
	count int
	data  []*T
}

func (r *rollingCollection[T]) removeIndex(index int) {
	for i := index + 1; i < len(r.data); i++ {
		r.data[i-1] = r.data[i]
	}
	r.count--
}

func (r *rollingCollection[T]) Remove(removeMe T) int {
	r.lock.Lock()
	defer r.lock.Unlock()
	to := len(r.data)
	removals := 0
	for i := 0; i < to; i++ {
		ptr := r.data[i]
		if *ptr == removeMe {
			r.removeIndex(i)
			removals++
			to--
		}
	}
	return removals
}

func (r *rollingCollection[T]) Add(value T) {
	r.lock.Lock()
	r.lock.Unlock()
	if r.count >= cap(r.data) {
		r.removeIndex(0)
	}
	r.data[r.count] = &value
	r.count++
}

func (r *rollingCollection[T]) Length() int {
	return r.count
}

func (r *rollingCollection[T]) Get(index int) *T {
	if index < 0 || index >= r.count {
		return nil
	}
	return r.data[index]
}

func (r *rollingCollection[T]) GetRange() []*T {
	return r.data[0:r.count]
}

func NewRollingCollection[T comparable](maxElements int) Collection[T] {
	return &rollingCollection[T]{
		lock:  sync.Mutex{},
		count: 0,
		data:  make([]*T, maxElements),
	}
}
