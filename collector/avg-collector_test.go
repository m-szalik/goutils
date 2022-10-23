package collector

import (
	"github.com/m-szalik/goutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAvgCollector(t *testing.T) {
	maxSamples := 10
	t.Run("avg-collector", func(t *testing.T) {
		mockTp := goutils.NewMockTimeProvider()
		mockTp.Add(1 * time.Second)
		avgc := newAvgCollectorInternal(maxSamples, mockTp)
		avgc.Collect(4)
		mockTp.Add(1 * time.Second)
		avgc.Collect(2)
		s1End := mockTp.Now()
		assert.Equal(t, 2.0, avgc.Avg(s1End, 500*time.Millisecond))
		assert.Equal(t, 3.0, avgc.Avg(s1End, 4*time.Second))
		for i := 0; i < maxSamples+1; i++ {
			avgc.Collect(5)
		}
		assert.Equal(t, 5.0, avgc.Avg(s1End.Add(1*time.Second), 14*time.Second))
	})
}
