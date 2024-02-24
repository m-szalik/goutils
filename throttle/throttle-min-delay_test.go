package throttle

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type event struct {
	when    time.Time
	payload any
}
type events struct {
	events []event
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
				testEvents.events = append(testEvents.events, event{
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
		assert.Equal(t, 2, len(testEvents.events))
		if len(testEvents.events) > 1 {
			assert.Equal(t, "one", testEvents.events[0].payload)
			assert.Equal(t, "three", testEvents.events[1].payload)
		}
	})

	t.Run("without_events", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()
		th := NewMinDelayThrottler[any](ctx, 1*time.Second)
		testEvents := testThrottler(ctx, th)
		time.Sleep(2001 * time.Millisecond)
		assert.Equal(t, 0, len(testEvents.events))
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
		assert.Equal(t, 1, len(testEvents.events))
		if len(testEvents.events) > 0 {
			assert.Equal(t, "two", testEvents.events[0].payload)
		}
	})

}
