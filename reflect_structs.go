package goutils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type CmpError interface {
	error
	FieldPath() string
	A() any
	B() any
}

type cmpErrorImpl struct {
	fieldPath string
	msg       string
	a         any
	b         any
}

func (c cmpErrorImpl) A() any {
	return c.a
}

func (c cmpErrorImpl) B() any {
	return c.b
}

func (c cmpErrorImpl) FieldPath() string {
	return c.fieldPath
}

func (c cmpErrorImpl) Error() string {
	return fmt.Sprintf("field: %s - %s", c.fieldPath, c.msg)
}

func cmpError(field string, a, b any, msg string) CmpError {
	return &cmpErrorImpl{field, msg, a, b}
}

func CmpWalkStructAreEqual(a interface{}, b interface{}) CmpError {
	valA := elemValue(a)
	valB := elemValue(b)
	return cmpValue(valA, valB, "")
}

func cmpValue(a reflect.Value, b reflect.Value, fieldPath string) (eError CmpError) {
	defer func() {
		if r := recover(); r != nil {
			eError = cmpError(fieldPath, a, b, eError.Error())
		}
	}()
	if a.Type() != b.Type() || a.Kind() != b.Kind() {
		return cmpError(fieldPath, a, b, fmt.Sprintf("A and B must be of the same type but %T and %T had been given", a, b))
	}
	if !a.IsValid() || !b.IsValid() {
		return cmpError(fieldPath, a, b, "not valid element")
	}
	switch a.Kind() {
	case reflect.Pointer:
		if a.IsNil() != b.IsNil() {
			return cmpError(fieldPath, a, b, "A and B must be of the same but one is nil")
		}
		return cmpValue(a.Elem(), b.Elem(), fieldPath)
	case reflect.Struct:
		for i := 0; i < a.NumField(); i++ {
			aFieldValue := a.Field(i)
			bFieldValue := b.Field(i)
			structKey := a.Type().Field(i).Name
			if err := cmpValue(aFieldValue, bFieldValue, fmt.Sprintf("%s.%s", fieldPath, structKey)); err != nil {
				return err
			}
		}
	case reflect.Array, reflect.Slice:
		if a.Len() != b.Len() {
			return cmpError(fieldPath, a, b, "arrays are not the same size")
		}
		if a.IsNil() != b.IsNil() {
			return cmpError(fieldPath, a, b, "A and B must be of the same but one is nil")
		}
		for i := 0; i < a.Len(); i++ {
			if err := cmpValue(a.Index(i), b.Index(i), fmt.Sprintf("%s[%d]", fieldPath, i)); err != nil {
				return err
			}
		}
	default:
		if a.Interface() != b.Interface() {
			return cmpError(fieldPath, a.Interface(), b.Interface(), fmt.Sprintf("%v != %v", a.Interface(), b.Interface()))
		}
	}
	return nil
}

// AcceptFunc - function used by CopyStruct function.
type AcceptFunc func(fieldPath string, srcValue reflect.Value) bool

// CopyStruct - copy struct from src to dst. acceptFunc is a function that selects fields to be copied.
func CopyStruct(src interface{}, dst interface{}, acceptFunc AcceptFunc) error {
	if reflect.TypeOf(dst).Kind() != reflect.Pointer {
		return fmt.Errorf("dst must be a pointer")
	}
	source := elemValue(src)
	destination := elemValue(dst)
	if source.Type() != destination.Type() || source.Kind() != destination.Kind() {
		return fmt.Errorf("source and destination must be of the same type but %T and %T had been given", source, destination)
	}
	return copyValue(source, destination, "", acceptFunc)
}

// CopyStructAll - copy all struct fields from src to dst.
func CopyStructAll(src interface{}, dst interface{}) error {
	return CopyStruct(src, dst, func(fieldPath string, srcValue reflect.Value) bool { return true })
}

// CopyStructSelected - copy selected fields from struct src to dst.
func CopyStructSelected(src interface{}, dst interface{}, selectedFilePaths ...string) error {
	return CopyStruct(src, dst, func(fieldPath string, srcValue reflect.Value) bool {
		if fieldPath == "" {
			return true
		}
		for _, s := range selectedFilePaths {
			if strings.Contains(fieldPath, s) {
				return true
			}
		}
		return false
	})
}

// CopyStructAllExcept - copy selected fields from struct src to dst.
func CopyStructAllExcept(src interface{}, dst interface{}, excludedFilePaths ...string) error {
	return CopyStruct(src, dst, func(fieldPath string, srcValue reflect.Value) bool {
		for _, s := range excludedFilePaths {
			if s == fieldPath {
				return false
			}
		}
		return true
	})
}

func elemValue(val interface{}) reflect.Value {
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
