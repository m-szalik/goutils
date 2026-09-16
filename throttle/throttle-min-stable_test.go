package throttle

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMinStable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	var last atomic.Int64
	th := NewMinStable[int](ctx, 1*time.Second, 0)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case l := <-th.Output():
				last.Store(int64(l))
			}
		}
	}()

	t.Run("no-change", func(t *testing.T) {
		th.Input() <- 0 // no change
		time.Sleep(2 * time.Second)
		assert.Equal(t, int64(0), last.Load())
	})

	t.Run("no-changes-but-to-be-ignored", func(t *testing.T) {
		th.Input() <- 1 // change
		th.Input() <- 0 // back
		time.Sleep(2 * time.Second)
		assert.Equal(t, int64(0), last.Load())
	})

	t.Run("real-change", func(t *testing.T) {
		th.Input() <- 2 // change
		time.Sleep(2 * time.Second)
		assert.Equal(t, int64(2), last.Load())
	})

}

func TestNewMinStableRapidChanges(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	th := NewMinStable[any](ctx, 100*time.Millisecond, 0)
	testEvents := testThrottler(ctx, th)
	th.Input() <- 1
	time.Sleep(70 * time.Millisecond)
	th.Input() <- 2
	time.Sleep(70 * time.Millisecond)
	th.Input() <- 3
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 0, testEvents.len(), "no value was stable long enough yet")
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 1, testEvents.len())
	assert.Equal(t, 3, testEvents.payload(0))
}
