package goutils

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldStopWhenContextExited(t *testing.T) {
	tests := []struct {
		name            string
		tasksCount      int
		cancelAfter     int
		expectedError   error
		expectedCounter int
	}{
		{
			name:            "Should finish",
			tasksCount:      2,
			cancelAfter:     10,
			expectedCounter: 2,
		},
		{
			name:            "Should cancel before finish",
			tasksCount:      10,
			cancelAfter:     3,
			expectedCounter: 3,
			expectedError:   context.Canceled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			counter := 0
			var waits []time.Duration
			for i := 0; i < tt.tasksCount; i++ {
				waits = append(waits, 200*time.Millisecond)
			}
			err := IterateOver(ctx, waits, func(index int, element time.Duration) error {
				counter++
				if counter == tt.cancelAfter {
					cancel()
				}
				return nil
			})
			assert.ErrorIs(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedCounter, counter)
		})
	}
}
