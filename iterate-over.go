package goutils

import (
	"context"
	"fmt"
)

// IterateOver iterate over slice until the end or until context exited.
func IterateOver[T any](ctx context.Context, elements []T, callback func(index int, element T) error) error {
	if elements != nil {
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
	}
	return nil
}
