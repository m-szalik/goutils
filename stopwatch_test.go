package goutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_stopWatch(t *testing.T) {
	sw := NewStopWatch()
	str0 := fmt.Sprint(sw)
	sw.Start()
	str1 := fmt.Sprint(sw)
	sw.Stop()
	str2 := fmt.Sprint(sw)
	fmt.Println(str0, str1, str2)
	assert.Equal(t, "0s", str0)
	assert.True(t, strings.HasSuffix(str1, "+"))
	assert.True(t, !strings.HasSuffix(str2, "+"))
}
