package goutils

import "fmt"

var ErrSafeSendChannelClosed = fmt.Errorf("channel closed")

func SafeSend[T any](ch chan<- T, value T) (exitErr error) {
	defer func() {
		if r := recover(); r != nil {
			exitErr = ErrSafeSendChannelClosed
		}
	}()
	ch <- value
	return nil
}
