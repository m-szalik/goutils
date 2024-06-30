package goutils

import (
	"fmt"
	"reflect"
)

// IteratorCallback
// Parameters:
//
//	path position as string;
//	depth element depth in the structure;
//	kind - kind of the element;
//	element - item tp iterate over;
//
// return false if stop iterating
type IteratorCallback func(path string, depth int, kind reflect.Kind, element interface{}) bool

// IterateDeep - iterate over the structure element. For each element callback function is called.
// All pointers are passed into the callback as value.
func IterateDeep(element interface{}, callback IteratorCallback) {
	iterateDeep(".", 0, reflect.ValueOf(element), callback)
}

func iterateDeep(path string, depth int, val reflect.Value, callback IteratorCallback) bool {
	kind := val.Kind()
	switch kind {
	case reflect.Pointer:
		uptrValue := val.Elem()
		return iterateDeep(path, depth, uptrValue, callback)
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			elementValue := val.Index(i)
			if !iterateDeep(fmt.Sprintf("%s[%d]", path, i), depth+1, elementValue, callback) {
				return false
			}
		}
		return true
	case reflect.Struct:
		if path == "." {
			path = ""
		}
		for i := 0; i < val.NumField(); i++ {
			ft := val.Type().Field(i)
			if !ft.IsExported() {
				continue
			}
			fieldName := ft.Name
			fieldValue := val.FieldByName(fieldName)
			if !iterateDeep(fmt.Sprintf("%s.%s", path, fieldName), depth+1, fieldValue, callback) {
				return false
			}
		}
		return true
	case reflect.Map:
		if path == "." {
			path = ""
		}
		for _, mapKey := range val.MapKeys() {
			mapValue := val.MapIndex(mapKey)
			if !iterateDeep(fmt.Sprintf("%s.%s", path, mapKey), depth+1, mapValue, callback) {
				return false
			}
		}
		return true
	default:
		element := val.Interface()
		return callback(path, depth, kind, element)
	}
}
