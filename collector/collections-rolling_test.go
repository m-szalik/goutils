package collector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_rollingCollectionAdd(t *testing.T) {
	type testCase struct {
		name           string
		noOfItemsToAdd int
		expected       []int
	}
	tests := []testCase{
		{name: "below maxElements", noOfItemsToAdd: 3, expected: []int{0, 1, 2}},
		{name: "over maxElements", noOfItemsToAdd: 10, expected: []int{5, 6, 7, 8, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := NewRollingCollection[int](5)
			for i := 0; i < tt.noOfItemsToAdd; i++ {
				col.Add(i)
			}
			assert.Equal(t, tt.expected, convert(col.AsSlice()))
		})
	}
}

func Test_rollingCollectionRemove(t *testing.T) {
	col := NewRollingCollection[int](5)
	for i := 0; i < 7; i++ {
		col.Add(i)
	}
	assert.Equal(t, []int{2, 3, 4, 5, 6}, convert(col.AsSlice()))
	col.Remove(3)
	assert.Equal(t, []int{2, 4, 5, 6}, convert(col.AsSlice()))
	col.Add(9)
	assert.Equal(t, []int{2, 4, 5, 6, 9}, convert(col.AsSlice()))
}

func convert[T comparable](source []*T) []T {
	conv := make([]T, len(source))
	for i, d := range source {
		conv[i] = *d
	}
	return conv
}

func Test_rollingCollectionAddMany(t *testing.T) {
	col := NewRollingCollection[int](5)
	assert.Equal(t, 3, col.Add(1, 2, 3))
	assert.Equal(t, []int{1, 2, 3}, convert(col.AsSlice()))
}

func Test_rollingCollectionRemoveNotFull(t *testing.T) {
	col := NewRollingCollection[int](5)
	col.Add(1, 2, 2, 3)
	assert.Equal(t, 0, col.Remove(9))
	assert.Equal(t, 2, col.Remove(2))
	assert.Equal(t, []int{1, 3}, convert(col.AsSlice()))
	assert.False(t, col.Contains(2))
}

func Test_rollingCollectionRemoveLastWhenFull(t *testing.T) {
	col := NewRollingCollection[int](3)
	col.Add(1, 2, 3)
	assert.Equal(t, 1, col.Remove(3))
	assert.Equal(t, 2, col.Length())
	assert.False(t, col.Contains(3))
	assert.Equal(t, []int{1, 2}, convert(col.AsSlice()))
}

func Test_rollingCollectionZeroCapacity(t *testing.T) {
	col := NewRollingCollection[int](0)
	assert.Equal(t, 0, col.Add(1))
	assert.Equal(t, 0, col.Length())
	timed := NewTimedCollection[int](0, time.Minute)
	assert.Equal(t, 0, timed.Add(1))
	assert.Equal(t, 0, timed.Length())
}

func Test_rollingCollectionAsSliceIsCopy(t *testing.T) {
	col := NewRollingCollection[int](2)
	col.Add(1, 2)
	snapshot := col.AsSlice()
	col.Add(3)
	assert.Equal(t, []int{1, 2}, convert(snapshot))
	assert.Equal(t, []int{2, 3}, convert(col.AsSlice()))
}
