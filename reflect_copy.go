package goutils

import (
	"errors"
	"fmt"
	"reflect"
)

type AcceptFunc func(fieldPath string, srcValue reflect.Value) bool

func CopyStructAll(src interface{}, dst interface{}) error {
	return CopyStruct(src, dst, func(fieldPath string, srcValue reflect.Value) bool { return true })
}

func CopyStruct(src interface{}, dst interface{}, acceptFunc AcceptFunc) error {
	if reflect.TypeOf(dst).Kind() != reflect.Pointer {
		return fmt.Errorf("dst must be a pointer")
	}
	source := dereference(src)
	destination := dereference(dst)
	if source.Type() != destination.Type() {
		return fmt.Errorf("source and destination must be of the same type but %T and %T had been given", source, destination)
	}
	return copyValue(source, destination, "", acceptFunc)
}

func dereference(val interface{}) reflect.Value {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Pointer {
		return v.Elem()
	} else {
		return v
	}
}

func copyValue(source reflect.Value, destination reflect.Value, fieldPath string, acceptFunc AcceptFunc) (eError error) {
	defer func() {
		if r := recover(); r != nil {
			eError = errors.New(fmt.Sprint(r))
		}
	}()
	if !source.IsValid() || source.IsZero() {
		return nil
	}
	if source.Kind() == reflect.Pointer && source.IsNil() {
		return nil
	}
	if !acceptFunc(fieldPath, source) {
		return nil
	}
	switch source.Kind() {
	case reflect.Struct:
		for i := 0; i < source.NumField(); i++ {
			srcFieldValue := source.Field(i)
			dstFieldValue := destination.Field(i)
			structKey := source.Type().Field(i).Name
			newFieldPath := fmt.Sprintf("%s.%s", fieldPath, structKey)
			if err := copyValue(srcFieldValue, dstFieldValue, newFieldPath, acceptFunc); err != nil {
				return err
			}
		}
	case reflect.Array, reflect.Slice:
		// make an array
		sliceInstance := reflect.MakeSlice(source.Type(), source.Len(), source.Cap())
		destination.Set(sliceInstance)
		// copy values
		for i := 0; i < source.Len(); i++ {
			if err := copyValue(source.Index(i), destination.Index(i), fmt.Sprintf("%s[%d]", fieldPath, i), acceptFunc); err != nil {
				return err
			}
		}
	case reflect.String:
		destination.Set(source)
	default:
		destination.Set(source)
	}
	return nil
}
