package throttle

import (
	"context"
	"time"
)

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
		ticker := time.NewTicker(ignoreDuration)
		ticker.Stop()
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case inp := <-t.input:
				if inp != lastInput {
					lastInput = inp
					ticker.Reset(ignoreDuration)
				}
			case <-ticker.C:
				if lastInput != lastReported {
					lastReported = lastInput
					t.output <- lastInput
					ticker.Stop()
				}
			}
		}
	}()
	return t
}
