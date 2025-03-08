package goutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSafeSend(t *testing.T) {
	t.Run("send to open channel", func(t *testing.T) {
		ch := make(chan string, 1)
		defer close(ch)
		err := SafeSend(ch, "hello")
		assert.NoError(t, err)
		str := <-ch
		assert.Equal(t, "hello", str)
	})

	t.Run("send to closed channel", func(t *testing.T) {
		ch := make(chan string, 1)
		close(ch)
		err := SafeSend(ch, "hello")
		assert.ErrorIs(t, err, SafeSendChannelClosedError)
	})
}
