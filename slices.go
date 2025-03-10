package goutils

// SliceIndexOf find an element in a slice
// return index of the element in a slice or -1 if not found
func SliceIndexOf[T comparable](slice []T, e T) int {
	for i, a := range slice {
		if a == e {
			return i
		}
	}
	return -1
}

// SliceContains check if slice contains the element
func SliceContains[T comparable](slice []T, e T) bool {
	return SliceIndexOf(slice, e) >= 0
}

// SliceRemove remove the element from a slice.
// return []T = new slice, int number of removed elements
func SliceRemove[T comparable](slice []T, e any) ([]T, int) {
	mySlice := slice
	size := len(slice)
	removed := 0
	for i := 0; i < size-removed; i++ {
		if mySlice[i] == e {
			mySlice = append(mySlice[:i], mySlice[i+1:]...)
			removed++
			i--
		}
	}
	if removed > 0 {
		return mySlice, removed
	} else {
		return slice, removed
	}
}

// SliceMap - map slice to slice of different object.
// sliceOfStrings := SliceMap[int, string]([]int{2, 7, -11}, func(i int) string { return fmt.Sprint(i) })
func SliceMap[I any, O any](inputData []I, mapper func(I) O) []O {
	if inputData == nil {
		return nil
	}
	outputData := make([]O, len(inputData))
	for i := 0; i < len(inputData); i++ {
		outputData[i] = mapper(inputData[i])
	}
	return outputData
}

// SlicesEq - check if two slices are equal
func SlicesEq[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
