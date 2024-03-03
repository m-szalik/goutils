package goutils

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestTimeCounter(t *testing.T) {
	t.Run("time-counter", func(t *testing.T) {
		mockTp := NewMockTimeProvider()
		dc := newTimeCounterInternal(mockTp)
		assert.Equal(t, time.Duration(0), dc.Value())
		dc.Start()
		mockTp.Add(2 * time.Second) // 2s
		dc.Stop()
		v1 := dc.Value()
		assertDuration(t, 2*time.Second, v1)
		mockTp.Add(2 * time.Second) // 2s
		v2 := dc.Value()
		assert.Equal(t, v1, v2)
		dc.Start()
		mockTp.Add(2 * time.Second) // 4s
		dc.Stop()
		v3 := dc.Value()
		assertDuration(t, 4*time.Second, v3)
		dc.Stop()
		dc.Reset()
		mockTp.Add(24 * time.Minute)
		assertDuration(t, 0, dc.Value())
	})

}

func assertDuration(t *testing.T, expected time.Duration, val time.Duration) {
	if math.Abs(float64(expected-val)) > float64(50*time.Millisecond) {
		assert.Equal(t, expected, val)
	}
}
