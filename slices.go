package goutils

// AllMatch - check if all elements in the slice match the condition
func AllMatch[T any](input []T, condition func(element T) bool) bool {
	for _, element := range input {
		if !condition(element) {
			return false
		}
	}
	return true
}

// AnyMatch - check if any element in the slice match the condition
func AnyMatch[T any](input []T, condition func(element T) bool) bool {
	for _, element := range input {
		if condition(element) {
			return true
		}
	}
	return false
}

// CountMatch - count elements that matches the condition
func CountMatch[T any](input []T, condition func(element T) bool) int {
	counter := 0
	for _, element := range input {
		if condition(element) {
			counter++
		}
	}
	return counter
}

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

func DistrictValues[T comparable](input []T) []T {
	if input == nil || len(input) < 2 {
		return input
	}
	found := make(map[T]bool)
	ret := make([]T, 0)
	for _, item := range input {
		_, ok := found[item]
		if !ok {
			found[item] = true
			ret = append(ret, item)
		}
	}
	return ret
}

// SliceAllMatch - check if all elements in the slice match the condition
func SliceAllMatch[T any](input []T, condition func(element T) bool) bool {
	for _, element := range input {
		if !condition(element) {
			return false
		}
	}
	return true
}

// SliceAnyMatch - check if any element in the slice match the condition
func SliceAnyMatch[T any](input []T, condition func(element T) bool) bool {
	for _, element := range input {
		if condition(element) {
			return true
		}
	}
	return false
}

// SliceCountMatch - count elements that matches the condition
func SliceCountMatch[T any](input []T, condition func(element T) bool) int {
	counter := 0
	for _, element := range input {
		if condition(element) {
			counter++
		}
	}
	return counter
}
