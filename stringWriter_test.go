package goutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_stringWriter_Write(t *testing.T) {
	t.Run("write huge content", func(t *testing.T) {
		sw := NewStringWriter()
		content := "1234567890abcdefghijklmnoprstuwxyz"
		got1, err1 := sw.Write([]byte(content))
		got2, err2 := sw.Write([]byte(" "))
		got3, err3 := sw.Write([]byte(content))
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
		assert.Equal(t, len(content), got1)
		assert.Equal(t, 1, got2)
		assert.Equal(t, len(content), got3)
		assert.Equal(t, content+" "+content, sw.String())
	})

}
