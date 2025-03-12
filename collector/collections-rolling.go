package collector

import "sync"

type rollingCollection[T comparable] struct {
	lock  sync.Mutex
	count int
	data  []*T
}

func (c *rollingCollection[T]) removeIndex(index int) {
	for i := index + 1; i < len(c.data); i++ {
		c.data[i-1] = c.data[i]
	}
	c.count--
}

func (c *rollingCollection[T]) Remove(removeMeElements ...T) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	to := len(c.data)
	removals := 0
	for _, removeMe := range removeMeElements {
		for i := 0; i < to; i++ {
			ptr := c.data[i]
			if *ptr == removeMe {
				c.removeIndex(i)
				removals++
				to--
			}
		}
	}
	return removals
}

func (c *rollingCollection[T]) Add(values ...T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, v := range values {
		if c.count >= cap(c.data) {
			c.removeIndex(0)
		}
		c.data[c.count] = &v
		c.count++
	}
}

func (c *rollingCollection[T]) Length() int {
	return c.count
}

func (c *rollingCollection[T]) Get(index int) *T {
	if index < 0 || index >= c.count {
		return nil
	}
	return c.data[index]
}

func (c *rollingCollection[T]) AsSlice() []*T {
	return c.data[0:c.count]
}

// NewRollingCollection - collection that keeps maxElements only - the oldest elements are removed automatically
func NewRollingCollection[T comparable](maxElements int) Collection[T] {
	return &rollingCollection[T]{
		lock:  sync.Mutex{},
		count: 0,
		data:  make([]*T, maxElements),
	}
}
