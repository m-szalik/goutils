package goutils

import "fmt"

// ErrSafeSendChannelClosed is returned by SafeSend when the target channel is closed.
var ErrSafeSendChannelClosed = fmt.Errorf("channel closed")

// SafeSend writes value to ch and returns ErrSafeSendChannelClosed when ch is closed.
func SafeSend[T any](ch chan<- T, value T) (exitErr error) {
	defer func() {
		if r := recover(); r != nil {
			exitErr = ErrSafeSendChannelClosed
		}
	}()
	ch <- value
	return nil
}
