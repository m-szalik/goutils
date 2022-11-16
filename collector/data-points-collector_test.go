package collector

import (
	"github.com/m-szalik/goutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDataPointsCollector(t *testing.T) {
	maxSamples := 10
	mockTp := goutils.NewMockTimeProvider()
	mockTp.Add(1 * time.Second)
	dpCollector := newDataPointsCollectorInternal(maxSamples, mockTp)
	dpCollector.Collect(4)
	mockTp.Add(1 * time.Second)
	dpCollector.Collect(2)
	s1End := mockTp.Now()

	t.Run("Avg", func(t *testing.T) {
		assert.Equal(t, 2.0, dpCollector.Avg(s1End, 500*time.Millisecond))
		assert.Equal(t, 3.0, dpCollector.Avg(s1End, 4*time.Second))
		dpCollectorFork := dpCollector.Fork()
		for i := 0; i < maxSamples+1; i++ {
			dpCollectorFork.Collect(5)
		}
		assert.Equal(t, 5.0, dpCollectorFork.Avg(s1End.Add(1*time.Second), 14*time.Second))
	})

	t.Run("GetDataPointN", func(t *testing.T) {
		dp, tp, err := dpCollector.GetDataPointN(0)
		assert.NoError(t, err)
		assert.Equal(t, tp, &s1End)
		assert.Equal(t, float64(2), dp)
	})

	t.Run("GetDataPointsBetween", func(t *testing.T) {
		data := dpCollector.GetDataPointsBetween(s1End.Add(-1*time.Millisecond), s1End.Add(1*time.Millisecond))
		assert.Equal(t, []float64{2}, data)
	})

	t.Run("Max", func(t *testing.T) {
		maxValue := dpCollector.Max(s1End.Add(+1*time.Millisecond), 10*time.Minute)
		assert.Equal(t, float64(4), maxValue)
	})

	t.Run("Min", func(t *testing.T) {
		minValue := dpCollector.Min(s1End.Add(+1*time.Millisecond), 10*time.Minute)
		assert.Equal(t, float64(2), minValue)
	})
}
