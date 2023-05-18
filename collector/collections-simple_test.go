package collector

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_simpleCollectionRemove(t *testing.T) {
	col := NewSimpleCollection[int]()
	for i := 0; i < 7; i++ {
		col.Add(i)
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6}, convert(col.AsSlice()))
	col.Remove(3)
	assert.Equal(t, []int{0, 1, 2, 4, 5, 6}, convert(col.AsSlice()))
	col.Add(9)
	assert.Equal(t, []int{0, 1, 2, 4, 5, 6, 9}, convert(col.AsSlice()))
	assert.Equal(t, 7, col.Length())
}

func Test_simpleCollectionString(t *testing.T) {
	col := NewSimpleCollection[int]()
	for i := 0; i < 3; i++ {
		col.Add(i)
	}
	assert.Equal(t, "0,1,2", fmt.Sprint(col))
}
