package goutils

import "fmt"

var SafeSendChannelClosedError = fmt.Errorf("channel closed")

func SafeSend[T any](ch chan T, value T) (exitErr error) {
	defer func() {
		if r := recover(); r != nil {
			exitErr = SafeSendChannelClosedError
		}
	}()
	ch <- value
	return nil
}
