package collector

import (
	"github.com/m-szalik/goutils"
	"sync"
	"time"
)

type AvgCollector interface {
	Collect(value float64)
	Avg(end time.Time, dur time.Duration) float64
}

type avgDataEntry struct {
	t     time.Time
	value float64
}

type avgCollectorImpl struct {
	timeProvider goutils.TimeProvider
	lock         sync.Mutex
	data         []*avgDataEntry
	maxSamples   int
	index        int
}

func (a *avgCollectorImpl) Collect(value float64) {
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
	a.data[a.index] = &avgDataEntry{
		t:     a.timeProvider.Now(),
		value: value,
	}
}

func (a *avgCollectorImpl) Avg(end time.Time, backDuration time.Duration) float64 {
	a.lock.Lock()
	defer a.lock.Unlock()
	start := end.Add(-backDuration)
	var f float64
	var count int
	for _, entry := range a.data {
		if entry == nil || entry.t.After(end) {
			break
		}
		if entry.t.Before(start) {
			continue
		}
		f += entry.value
		count++
	}
	return f / float64(count)
}

func newAvgCollectorInternal(maxSamples int, tp goutils.TimeProvider) AvgCollector {
	return &avgCollectorImpl{
		timeProvider: tp,
		lock:         sync.Mutex{},
		data:         make([]*avgDataEntry, maxSamples),
		maxSamples:   maxSamples,
		index:        0,
	}
}

func NewAvgCollector(maxSamples int) AvgCollector {
	return newAvgCollectorInternal(maxSamples, goutils.SystemTimeProvider())
}
