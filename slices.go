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
