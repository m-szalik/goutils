package collector

import (
	"fmt"
	"strings"
	"sync"
)

type simpleCollection[T comparable] struct {
	lock sync.Mutex
	data []*T
}

func (c *simpleCollection[T]) removeIndex(index int) {
	for i := index + 1; i < len(c.data); i++ {
		c.data[i-1] = c.data[i]
	}
}

func (c *simpleCollection[T]) Remove(removeMeElements ...T) int {
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
	if removals > 0 {
		c.data = c.data[0:to]
	}
	return removals
}

func (c *simpleCollection[T]) Add(elements ...T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, elem := range elements {
		c.data = append(c.data, &elem)
	}
}

func (c *simpleCollection[T]) Length() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

func (c *simpleCollection[T]) Get(index int) *T {
	c.lock.Lock()
	defer c.lock.Unlock()
	if index < 0 || index >= len(c.data) {
		return nil
	}
	return c.data[index]
}

func (c *simpleCollection[T]) AsSlice() []*T {
	c.lock.Lock()
	defer c.lock.Unlock()
	ret := make([]*T, len(c.data))
	copy(ret, c.data)
	return ret
}

func (c *simpleCollection[T]) String() string {
	strs := make([]string, len(c.data))
	func() {
		c.lock.Lock()
		defer c.lock.Unlock()
		for i, e := range c.data {
			if e == nil {
				strs[i] = "nil"
			} else {
				strs[i] = fmt.Sprint(*e)
			}
		}
	}()
	return strings.Join(strs, ",")
}

// NewSimpleCollection - collection that keeps all elements, slice that grows when needed
func NewSimpleCollection[T comparable]() Collection[T] {
	return &simpleCollection[T]{
		lock: sync.Mutex{},
		data: make([]*T, 0),
	}
}
