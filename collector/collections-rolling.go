package collector

import "sync"

type rollingCollection[T comparable] struct {
	lock  sync.Mutex
	count int
	data  []*T
}

func (c *rollingCollection[T]) removeIndex(index int) {
	copy(c.data[index:c.count-1], c.data[index+1:c.count])
	c.count--
	c.data[c.count] = nil
}

func (c *rollingCollection[T]) Remove(removeMeElements ...T) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	removals := 0
	for _, removeMe := range removeMeElements {
		for i := 0; i < c.count; {
			if *c.data[i] == removeMe {
				c.removeIndex(i)
				removals++
			} else {
				i++
			}
		}
	}
	return removals
}

func (c *rollingCollection[T]) Add(values ...T) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	if cap(c.data) == 0 {
		return 0
	}
	added := 0
	for _, v := range values {
		value := v
		if c.count >= cap(c.data) {
			c.removeIndex(0)
		}
		c.data[c.count] = &value
		c.count++
		added++
	}
	return added
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

func (c *rollingCollection[T]) Contains(element T) bool {
	for _, e := range c.data {
		if e == nil {
			continue
		}
		if *e == element {
			return true
		}
	}
	return false
}

func (c *rollingCollection[T]) AsSlice() []*T {
	return c.data[0:c.count]
}

// NewRollingCollection returns a fixed-size collection that keeps at most
// maxElements newest values.
func NewRollingCollection[T comparable](maxElements int) Collection[T] {
	return &rollingCollection[T]{
		lock:  sync.Mutex{},
		count: 0,
		data:  make([]*T, maxElements),
	}
}
