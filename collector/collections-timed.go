package collector

import (
	"sync"
	"time"

	"github.com/m-szalik/goutils"
)

type timedWrapper[T comparable] struct {
	time    time.Time
	element T
}

type timedCollection[T comparable] struct {
	rollingCollection[timedWrapper[T]]
	timeProvider goutils.TimeProvider
	duration     time.Duration
}

func (c *timedCollection[T]) Remove(removeMeElements ...T) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	to := len(c.data)
	removals := 0
	for _, removeMe := range removeMeElements {
		for i := 0; i < to; i++ {
			ptr := c.data[i]
			if ptr.element == removeMe {
				c.removeIndex(i)
				removals++
				to--
			}
		}
	}
	return removals
}

func (c *timedCollection[T]) Add(values ...T) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, value := range values {
		if c.count >= cap(c.data) {
			if c.cleanup() == 0 {
				c.removeIndex(0)
			}
		}
		c.data[c.count] = &timedWrapper[T]{
			time:    c.timeProvider.Now().Add(c.duration),
			element: value,
		}
	}
	c.count++
	return len(values)
}

func (c *timedCollection[T]) Length() int {
	return c.count
}

func (c *timedCollection[T]) Get(index int) *T {
	if index < 0 || index >= c.count {
		return nil
	}
	wrapper := c.data[index]
	return &wrapper.element
}

func (c *timedCollection[T]) AsSlice() []*T {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cleanup()
	data := make([]*T, c.count)
	for i, e := range c.data[0:c.count] {
		data[i] = &e.element
	}
	return data
}

func (c *timedCollection[T]) Contains(element T) bool {
	for _, e := range c.data {
		if e == nil {
			continue
		}
		if e.element == element {
			return true
		}
	}
	return false
}

func (c *timedCollection[T]) cleanup() int {
	now := c.timeProvider.Now()
	count := 0
	for i := 0; i < c.count; {
		if c.data[i].time.Before(now) {
			c.removeIndex(i)
			count++
		} else {
			i++
		}
	}
	return count
}

func newTimedCollectionWithTimeProvider[T comparable](maxElements int, duration time.Duration, timeProvider goutils.TimeProvider) Collection[T] {
	return &timedCollection[T]{
		rollingCollection: rollingCollection[timedWrapper[T]]{
			lock:  sync.Mutex{},
			count: 0,
			data:  make([]*timedWrapper[T], maxElements),
		},
		timeProvider: timeProvider,
		duration:     duration,
	}
}

// NewTimedCollection collection that keeps elements for defined duration only
func NewTimedCollection[T comparable](maxElements int, duration time.Duration) Collection[T] {
	return newTimedCollectionWithTimeProvider[T](maxElements, duration, goutils.SystemTimeProvider())
}
