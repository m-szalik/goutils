package collector

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
			assert.Equal(t, tt.expected, convert(col.GetRange()))
		})
	}
}

func Test_rollingCollectionRemove(t *testing.T) {
	col := NewRollingCollection[int](5)
	for i := 0; i < 7; i++ {
		col.Add(i)
	}
	assert.Equal(t, []int{2, 3, 4, 5, 6}, convert(col.GetRange()))
	col.Remove(3)
	assert.Equal(t, []int{2, 4, 5, 6}, convert(col.GetRange()))
	col.Add(9)
	assert.Equal(t, []int{2, 4, 5, 6, 9}, convert(col.GetRange()))
}

func convert[T comparable](source []*T) []T {
	conv := make([]T, len(source))
	for i, d := range source {
		conv[i] = *d
	}
	return conv
}
