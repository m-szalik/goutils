package goutils

import (
	"fmt"
	"reflect"
	"sort"
)

// IteratorCallback is called by IterateDeep for every leaf value.
// Return false to stop traversal early.
type IteratorCallback func(path string, depth int, kind reflect.Kind, element interface{}) bool

// IterateDeep walks nested structs, maps, slices, arrays, pointers and
// interfaces and calls callback for each non-container value. Pointer and
// interface values are dereferenced before callback execution; nil pointers
// and nil interfaces are reported with a nil element. Map entries are visited
// in the order of their formatted keys.
func IterateDeep(element interface{}, callback IteratorCallback) {
	iterateDeep(".", 0, reflect.ValueOf(element), callback)
}

func iterateDeep(path string, depth int, val reflect.Value, callback IteratorCallback) bool {
	kind := val.Kind()
	switch kind {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return callback(path, depth, kind, nil)
		}
		return iterateDeep(path, depth, val.Elem(), callback)
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
		type entry struct {
			name string
			key  reflect.Value
		}
		entries := make([]entry, 0, val.Len())
		for _, mapKey := range val.MapKeys() {
			entries = append(entries, entry{name: fmt.Sprint(mapKey), key: mapKey})
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })
		for _, e := range entries {
			if !iterateDeep(fmt.Sprintf("%s.%s", path, e.name), depth+1, val.MapIndex(e.key), callback) {
				return false
			}
		}
		return true
	case reflect.Invalid:
		return callback(path, depth, kind, nil)
	default:
		return callback(path, depth, kind, val.Interface())
	}
}
