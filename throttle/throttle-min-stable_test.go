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
