package goutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_root(t *testing.T) {
	accuracy := 0.005
	tests := []struct {
		x              float64
		n              int
		accuracy       float64
		expectedResult float64
	}{
		{9, 2, 0.001, 3},
		{27, 3, 0.001, 3},
		{2, 2, 0.001, 1.414213},
		{2, 2, 0.00001, 1.414213},
		{881, 8, 0.001, 2.33},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("root(%.2f,%d,%f) should be equal to %f", tt.x, tt.n, tt.accuracy, tt.expectedResult), func(t *testing.T) {
			result := rootN(tt.x, float64(tt.n), tt.accuracy)
			delta := math.Abs(tt.expectedResult - result)
			fmt.Printf("root(%.2f, %d, %f) ~= %f should be %f -> delta: %f\n", tt.x, tt.n, tt.accuracy, result, tt.expectedResult, delta)
			if delta > accuracy {
				assert.Equalf(t, tt.expectedResult, result, "root(%v, %v, %f)", tt.x, tt.n, tt.accuracy)
			}
		})
	}
}
