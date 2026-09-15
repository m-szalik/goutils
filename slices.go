package goutils

// FindFirst returns the first element matching condition, or nil when none
// match.
func FindFirst[T any](input []T, condition func(element T) bool) *T {
	for _, element := range input {
		if condition(element) {
			return &element
		}
	}
	return nil
}

// Filter returns a new slice containing elements that match condition.
func Filter[T any](input []T, condition func(element T) bool) []T {
	ret := make([]T, 0)
	for _, element := range input {
		if condition(element) {
			ret = append(ret, element)
		}
	}
	return ret
}

// AllMatch reports whether all elements match condition.
func AllMatch[T any](input []T, condition func(element T) bool) bool {
	for _, element := range input {
		if !condition(element) {
			return false
		}
	}
	return true
}

// AnyMatch reports whether at least one element matches condition.
func AnyMatch[T any](input []T, condition func(element T) bool) bool {
	for _, element := range input {
		if condition(element) {
			return true
		}
	}
	return false
}

// CountMatch returns the number of elements matching condition.
func CountMatch[T any](input []T, condition func(element T) bool) int {
	counter := 0
	for _, element := range input {
		if condition(element) {
			counter++
		}
	}
	return counter
}

// SliceIndexOf returns the index of e in slice, or -1 when not found.
func SliceIndexOf[T comparable](slice []T, e T) int {
	for i, a := range slice {
		if a == e {
			return i
		}
	}
	return -1
}

// SliceContains reports whether slice contains e.
func SliceContains[T comparable](slice []T, e T) bool {
	return SliceIndexOf(slice, e) >= 0
}

// SliceRemove removes all occurrences of e from slice.
// It returns the resulting slice and the number of removed elements.
// The input slice is left unchanged.
func SliceRemove[T comparable](slice []T, e any) ([]T, int) {
	result := make([]T, 0, len(slice))
	removed := 0
	for _, v := range slice {
		if v == e {
			removed++
		} else {
			result = append(result, v)
		}
	}
	if removed == 0 {
		return slice, 0
	}
	return result, removed
}

// SliceMap maps each element in inputData to a new value.
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

// SlicesEq reports whether slices a and b are equal.
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

// DistrictValues returns unique values from input, preserving first
// occurrence order.
func DistrictValues[T comparable](input []T) []T {
	if len(input) < 2 {
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

// SliceAllMatch reports whether all elements match condition.
func SliceAllMatch[T any](input []T, condition func(element T) bool) bool {
	for _, element := range input {
		if !condition(element) {
			return false
		}
	}
	return true
}

// SliceAnyMatch reports whether at least one element matches condition.
func SliceAnyMatch[T any](input []T, condition func(element T) bool) bool {
	for _, element := range input {
		if condition(element) {
			return true
		}
	}
	return false
}

// SliceCountMatch returns the number of elements matching condition.
func SliceCountMatch[T any](input []T, condition func(element T) bool) int {
	counter := 0
	for _, element := range input {
		if condition(element) {
			counter++
		}
	}
	return counter
}
