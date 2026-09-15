package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setAsSlice(t *testing.T) {
	s := NewSet[int]()
	assert.Equal(t, 3, s.Add(1, 2, 3, 3))
	assert.Equal(t, 3, s.Length())
	assert.ElementsMatch(t, []int{1, 2, 3}, convert(s.AsSlice()))
}
