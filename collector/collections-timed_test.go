package collector

import (
	"github.com/m-szalik/goutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_timedCollectionAdd(t *testing.T) {
	type addPhases struct {
		itemsToAdd []int
		addedTime  time.Duration
	}
	type testCase struct {
		name     string
		phases   []addPhases
		expected []int
	}
	tests := []testCase{
		{name: "over maxElements", phases: []addPhases{{itemsToAdd: []int{0, 1, 2, 3, 4, 5, 6}, addedTime: 0}}, expected: []int{2, 3, 4, 5, 6}},
		{name: "all expired", phases: []addPhases{{itemsToAdd: []int{0, 1, 2, 3, 4, 5, 6}, addedTime: 0}, {itemsToAdd: []int{}, addedTime: 5 * time.Second}}, expected: []int{}},
		{name: "some expired", phases: []addPhases{{itemsToAdd: []int{1, 2, 3, 4, 5}, addedTime: 0}, {itemsToAdd: []int{99}, addedTime: 5 * time.Second}, {itemsToAdd: []int{98}, addedTime: 1 * time.Second}}, expected: []int{99, 98}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := goutils.NewMockTimeProvider()
			col := newTimedCollectionWithTimeProvider[int](5, 3*time.Second, tp)
			for _, phase := range tt.phases {
				tp.Add(phase.addedTime)
				for _, elem := range phase.itemsToAdd {
					col.Add(elem)
				}
			}
			assert.Equal(t, tt.expected, convert(col.AsSlice()))
		})
	}
}
