package goutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
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

func Test_stopWatchWithMock(t *testing.T) {
	tp := NewMockTimeProvider()
	sw := NewStopWatchWithTimeProvider(tp)
	strBefore := fmt.Sprint(sw)
	sw.Start()
	tp.Add(3*time.Second + 500*time.Millisecond)
	sw.Stop()
	strEnd := fmt.Sprint(sw)
	assert.Equal(t, "0s", strBefore)
	assert.Equal(t, "3.5s", strEnd)
}
