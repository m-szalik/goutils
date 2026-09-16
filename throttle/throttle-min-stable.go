package throttle

import (
	"context"
	"time"
)

// NewMinStable emits value changes only after the input stays unchanged for
// ignoreDuration.
// initValue is the initial internal state and is not emitted automatically.
func NewMinStable[E comparable](ctx context.Context, ignoreDuration time.Duration, initValue E) Throttler[E] {
	t := &throttler[E]{
		input:  make(chan E),
		output: make(chan E),
	}
	go func() {
		defer close(t.input)
		defer close(t.output)
		var lastReported = initValue
		var lastInput = initValue
		var timer *time.Timer
		var timerC <-chan time.Time
		defer func() {
			if timer != nil {
				timer.Stop()
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case inp := <-t.input:
				if inp != lastInput {
					lastInput = inp
					if timer != nil {
						timer.Stop()
					}
					timer = time.NewTimer(ignoreDuration)
					timerC = timer.C
				}
			case <-timerC:
				timerC = nil
				if lastInput != lastReported {
					lastReported = lastInput
					t.output <- lastInput
				}
			}
		}
	}()
	return t
}
