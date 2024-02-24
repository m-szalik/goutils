package throttle

import (
	"context"
	"time"
)

func NewPeriodicThrottler[E any](ctx context.Context, period time.Duration) Throttler[E] {
	t := &throttler[E]{
		input:  make(chan E),
		output: make(chan E),
	}
	go func() {
		defer close(t.input)
		defer close(t.output)
		tick := time.NewTimer(period)
		defer tick.Stop()
		var lastInput *E = nil
		for {
			select {
			case <-ctx.Done():
				return
			case inp := <-t.input:
				lastInput = &inp
			case <-tick.C:
				if lastInput != nil {
					t.output <- *lastInput
					lastInput = nil
				}
			}
		}
	}()
	return t
}

func NewMinDelayThrottler[E any](ctx context.Context, minDelay time.Duration) Throttler[E] {
	t := &throttler[E]{
		input:  make(chan E),
		output: make(chan E),
	}
	go func() {
		defer close(t.input)
		defer close(t.output)
		nextAllowedPassAt := time.Now()
		for {
			select {
			case <-ctx.Done():
				return
			case inp := <-t.input:
				now := time.Now()
				if !now.Before(nextAllowedPassAt) {
					nextAllowedPassAt = now.Add(minDelay)
					t.output <- inp
				}
			}
		}
	}()
	return t
}
