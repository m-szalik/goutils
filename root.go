package goutils

import "math"

// Root - this function calculate approximation of n-th root of a number.
func Root(x, n float64) float64 {
	return root(x, n, 0.001)
}

// root - accuracyFactor = the lower number the lower error
func root(x, n float64, accuracyFactor float64) float64 {
	if x == 0 {
		return x
	}
	var y float64
	left := float64(0)
	right := math.Max(1, x)
	y = (left + right) / 2
	for (y - left) > accuracyFactor {
		check := math.Pow(y, n)
		if check > x {
			right = y
		} else if check < x {
			left = y
		} else {
			break
		}
		y = (left + right) / 2
	}
	return y
}
