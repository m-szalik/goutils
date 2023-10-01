package goutils

func SliceIndexOf(slice []interface{}, e interface{}) int {
	for i, a := range slice {
		if a == e {
			return i
		}
	}
	return -1
}

func SliceContains(slice []interface{}, e interface{}) bool {
	return SliceIndexOf(slice, e) >= 0
}

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
