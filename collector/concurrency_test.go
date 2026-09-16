package collector

import (
	"sync"
	"testing"
	"time"
)

// Exercises readers and writers concurrently; meaningful under -race.
func Test_concurrentAccess(t *testing.T) {
	collections := map[string]Collection[int]{
		"rolling": NewRollingCollection[int](8),
		"simple":  NewSimpleCollection[int](),
		"timed":   NewTimedCollection[int](8, time.Minute),
	}
	for name, col := range collections {
		t.Run(name, func(t *testing.T) {
			var wg sync.WaitGroup
			for i := 0; i < 4; i++ {
				wg.Add(2)
				go func(i int) {
					defer wg.Done()
					for j := 0; j < 50; j++ {
						col.Add(i*100 + j)
						col.Remove(i*100 + j - 1)
					}
				}(i)
				go func() {
					defer wg.Done()
					for j := 0; j < 50; j++ {
						_ = col.Length()
						_ = col.Contains(j)
						_ = col.AsSlice()
						if ic, ok := col.(IndexableCollection[int]); ok {
							_ = ic.Get(0)
						}
					}
				}()
			}
			wg.Wait()
		})
	}
	t.Run("stack", func(t *testing.T) {
		st := NewStack[int]()
		var wg sync.WaitGroup
		for i := 0; i < 4; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				for j := 0; j < 50; j++ {
					v := j
					st.Push(&v)
					st.Pop()
				}
			}()
			go func() {
				defer wg.Done()
				for j := 0; j < 50; j++ {
					_ = st.Length()
					_ = st.AsSlice()
					_ = st.Get(0)
				}
			}()
		}
		wg.Wait()
	})
	t.Run("datapoints", func(t *testing.T) {
		dp := NewDataPointsCollector(8)
		var wg sync.WaitGroup
		for i := 0; i < 4; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				for j := 0; j < 50; j++ {
					dp.Collect(float64(j))
				}
			}()
			go func() {
				defer wg.Done()
				for j := 0; j < 50; j++ {
					now := time.Now()
					_ = dp.GetDataPointsBetween(now.Add(-time.Minute), now)
					_ = dp.Avg(now, time.Minute)
					_, _, _ = dp.GetDataPointN(0)
				}
			}()
		}
		wg.Wait()
	})
}
