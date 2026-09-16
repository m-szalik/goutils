package throttle

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type event struct {
	when    time.Time
	payload any
}
type events struct {
	lock   sync.Mutex
	events []event
}

func (e *events) add(ev event) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.events = append(e.events, ev)
}

func (e *events) len() int {
	e.lock.Lock()
	defer e.lock.Unlock()
	return len(e.events)
}

func (e *events) payload(i int) any {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.events[i].payload
}

func testThrottler(ctx context.Context, th Throttler[any]) *events {
	testEvents := &events{
		events: make([]event, 0),
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case out := <-th.Output():
				testEvents.add(event{
					when:    time.Now(),
					payload: out,
				})
			}
		}
	}()
	return testEvents
}

func TestNewMinDelayThrottler(t *testing.T) {
	t.Run("with_events", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()
		th := NewMinDelayThrottler[any](ctx, 1*time.Second)
		testEvents := testThrottler(ctx, th)
		th.Input() <- "one"
		time.Sleep(500 * time.Millisecond)
		th.Input() <- "two"
		time.Sleep(501 * time.Millisecond)
		th.Input() <- "three"
		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, 2, testEvents.len())
		if testEvents.len() > 1 {
			assert.Equal(t, "one", testEvents.payload(0))
			assert.Equal(t, "three", testEvents.payload(1))
		}
	})

	t.Run("without_events", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()
		th := NewMinDelayThrottler[any](ctx, 1*time.Second)
		testEvents := testThrottler(ctx, th)
		time.Sleep(2001 * time.Millisecond)
		assert.Equal(t, 0, testEvents.len())
	})

}

func TestNewPeriodicThrottler(t *testing.T) {
	t.Run("with_events", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()
		th := NewPeriodicThrottler[any](ctx, 1*time.Second)
		testEvents := testThrottler(ctx, th)
		th.Input() <- "one"
		th.Input() <- "two"
		time.Sleep(1001 * time.Millisecond)
		th.Input() <- "three"
		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, 1, testEvents.len())
		if testEvents.len() > 0 {
			assert.Equal(t, "two", testEvents.payload(0))
		}
	})

}

func TestNewPeriodicThrottlerMultiplePeriods(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	th := NewPeriodicThrottler[any](ctx, 50*time.Millisecond)
	testEvents := testThrottler(ctx, th)
	for i := 1; i <= 3; i++ {
		th.Input() <- i
		time.Sleep(80 * time.Millisecond)
	}
	assert.Equal(t, 3, testEvents.len())
	for i := 0; i < testEvents.len(); i++ {
		assert.Equal(t, i+1, testEvents.payload(i))
	}
}
