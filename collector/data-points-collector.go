package collector

import (
	"fmt"
	"github.com/m-szalik/goutils"
	"math"
	"sync"
	"time"
)

type DataPointsCollector interface {
	Collect(value float64)

	// Fork - create copy
	Fork() DataPointsCollector

	// Avg - calculate average value
	Avg(end time.Time, backDuration time.Duration) float64

	Max(end time.Time, backDuration time.Duration) float64

	Min(end time.Time, backDuration time.Duration) float64

	// GetDataPointN - use n=0 for last point
	GetDataPointN(n int) (float64, *time.Time, error)

	GetDataPointsBetween(start, end time.Time) []float64
}

type dataPointEntry struct {
	t     time.Time
	value float64
}

type dataPointsCollectorImpl struct {
	timeProvider goutils.TimeProvider
	lock         sync.Mutex
	data         []*dataPointEntry
	maxSamples   int
	index        int
}

func (a *dataPointsCollectorImpl) Fork() DataPointsCollector {
	a.lock.Lock()
	defer a.lock.Unlock()
	newData := make([]*dataPointEntry, len(a.data))
	for i, dp := range a.data {
		if dp == nil {
			continue
		}
		newData[i] = &dataPointEntry{
			t:     dp.t,
			value: dp.value,
		}
	}
	return &dataPointsCollectorImpl{
		timeProvider: a.timeProvider,
		lock:         sync.Mutex{},
		maxSamples:   a.maxSamples,
		index:        a.index,
		data:         newData,
	}
}

func (a *dataPointsCollectorImpl) GetDataPointN(n int) (float64, *time.Time, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if n >= a.index {
		return 0, nil, fmt.Errorf("out of range")
	}
	i := a.index - 1 - n
	return a.data[i].value, &a.data[i].t, nil
}

func (a *dataPointsCollectorImpl) GetDataPointsBetween(start, end time.Time) []float64 {
	out := make([]float64, 0)
	for _, dp := range a.data {
		if dp == nil || dp.t.Before(start) {
			continue
		}
		if dp.t.After(end) {
			break
		}
		out = append(out, dp.value)
	}
	return out
}

func (a *dataPointsCollectorImpl) Collect(value float64) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.index > a.maxSamples-2 {
		for i := 0; i < a.maxSamples-1; i++ {
			a.data[i] = a.data[i+1]
		}
	} else {
		defer func() {
			a.index++
		}()
	}
	a.data[a.index] = &dataPointEntry{
		t:     a.timeProvider.Now(),
		value: value,
	}
}

func (a *dataPointsCollectorImpl) Avg(end time.Time, backDuration time.Duration) float64 {
	var f float64
	var count int
	a.filter(end, backDuration, func(i int, dp *dataPointEntry) {
		f += dp.value
		count++
	})
	return f / float64(count)
}

func (a *dataPointsCollectorImpl) Max(end time.Time, backDuration time.Duration) float64 {
	var f float64
	if a.filter(end, backDuration, func(i int, dp *dataPointEntry) {
		if i == 0 {
			f = dp.value
		} else {
			f = math.Max(f, dp.value)
		}
	}) {
		return f
	}
	return math.NaN()
}

func (a *dataPointsCollectorImpl) Min(end time.Time, backDuration time.Duration) float64 {
	var f float64
	if a.filter(end, backDuration, func(i int, dp *dataPointEntry) {
		if i == 0 {
			f = dp.value
		} else {
			f = math.Min(f, dp.value)
		}
	}) {
		return f
	}
	return math.NaN()
}

func (a *dataPointsCollectorImpl) filter(end time.Time, backDuration time.Duration, f func(i int, dp *dataPointEntry)) bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	start := end.Add(-backDuration)
	i := 0
	for _, entry := range a.data {
		if entry == nil || entry.t.After(end) {
			break
		}
		if entry.t.Before(start) {
			continue
		}
		f(i, entry)
		i++
	}
	return i > 0
}

func newDataPointsCollectorInternal(maxSamples int, tp goutils.TimeProvider) DataPointsCollector {
	return &dataPointsCollectorImpl{
		timeProvider: tp,
		lock:         sync.Mutex{},
		data:         make([]*dataPointEntry, maxSamples),
		maxSamples:   maxSamples,
		index:        0,
	}
}

func NewDataPointsCollector(maxSamples int) DataPointsCollector {
	return newDataPointsCollectorInternal(maxSamples, goutils.SystemTimeProvider())
}
