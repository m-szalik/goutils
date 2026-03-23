package goutils

import (
	"context"
	"fmt"
)

// IterateOver calls callback for each element until processing finishes, callback
// returns an error, or ctx is canceled.
func IterateOver[T any](ctx context.Context, elements []T, callback func(index int, element T) error) error {
	for index, e := range elements {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
			err := callback(index, e)
			if err != nil {
				return fmt.Errorf("error processing element %d: %w", index, err)
			}
		}
	}
	return nil
}
