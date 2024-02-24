package throttle

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMinStable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	var last = 0
	th := NewMinStable[int](ctx, 1*time.Second, last)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case l := <-th.Output():
				last = l
			}
		}
	}()

	t.Run("no-change", func(t *testing.T) {
		th.Input() <- 0 // no change
		time.Sleep(2 * time.Second)
		assert.Equal(t, 0, last)
	})

	t.Run("no-changes-but-to-be-ignored", func(t *testing.T) {
		th.Input() <- 1 // change
		th.Input() <- 0 // back
		time.Sleep(2 * time.Second)
		assert.Equal(t, 0, last)
	})

	t.Run("real-change", func(t *testing.T) {
		th.Input() <- 2 // change
		time.Sleep(2 * time.Second)
		assert.Equal(t, 2, last)
	})

}
