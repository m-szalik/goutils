package goutils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// CmpError describes the first mismatch found by CmpWalkStructAreEqual.
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

// CmpWalkStructAreEqual recursively compares a and b and returns nil when they
// are equal, or a CmpError describing the first mismatch.
func CmpWalkStructAreEqual(a interface{}, b interface{}) CmpError {
	valA := elemValue(a)
	valB := elemValue(b)
	return cmpValue(valA, valB, "")
}

func cmpValue(a reflect.Value, b reflect.Value, fieldPath string) (eError CmpError) {
	defer func() {
		if r := recover(); r != nil {
			eError = cmpError(fieldPath, reportValue(a), reportValue(b), fmt.Sprint(r))
		}
	}()
	if !a.IsValid() || !b.IsValid() {
		if !a.IsValid() && !b.IsValid() {
			return nil
		}
		return cmpError(fieldPath, reportValue(a), reportValue(b), "A and B must be of the same but one is nil")
	}
	if a.Type() != b.Type() {
		return cmpError(fieldPath, reportValue(a), reportValue(b), fmt.Sprintf("A and B must be of the same type but %s and %s had been given", a.Type(), b.Type()))
	}
	switch a.Kind() {
	case reflect.Pointer, reflect.Interface:
		if a.IsNil() != b.IsNil() {
			return cmpError(fieldPath, reportValue(a), reportValue(b), "A and B must be of the same but one is nil")
		}
		if a.IsNil() || (a.Kind() == reflect.Pointer && a.Pointer() == b.Pointer()) {
			return nil
		}
		return cmpValue(a.Elem(), b.Elem(), fieldPath)
	case reflect.Struct:
		for i := 0; i < a.NumField(); i++ {
			structKey := a.Type().Field(i).Name
			if err := cmpValue(a.Field(i), b.Field(i), fmt.Sprintf("%s.%s", fieldPath, structKey)); err != nil {
				return err
			}
		}
	case reflect.Array, reflect.Slice:
		if a.Len() != b.Len() {
			return cmpError(fieldPath, reportValue(a), reportValue(b), "arrays are not the same size")
		}
		if a.Kind() == reflect.Slice && a.IsNil() != b.IsNil() {
			return cmpError(fieldPath, reportValue(a), reportValue(b), "A and B must be of the same but one is nil")
		}
		for i := 0; i < a.Len(); i++ {
			if err := cmpValue(a.Index(i), b.Index(i), fmt.Sprintf("%s[%d]", fieldPath, i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		if a.Len() != b.Len() {
			return cmpError(fieldPath, reportValue(a), reportValue(b), "maps are not the same size")
		}
		if a.IsNil() != b.IsNil() {
			return cmpError(fieldPath, reportValue(a), reportValue(b), "A and B must be of the same but one is nil")
		}
		for _, key := range a.MapKeys() {
			bValue := b.MapIndex(key)
			if !bValue.IsValid() {
				return cmpError(fieldPath, reportValue(a), reportValue(b), fmt.Sprintf("key %v missing in B", key))
			}
			if err := cmpValue(a.MapIndex(key), bValue, fmt.Sprintf("%s[%v]", fieldPath, key)); err != nil {
				return err
			}
		}
	case reflect.Func, reflect.Chan, reflect.UnsafePointer:
		if a.Pointer() != b.Pointer() {
			return cmpError(fieldPath, reportValue(a), reportValue(b), fmt.Sprintf("%v != %v", a, b))
		}
	default:
		if !leafEqual(a, b) {
			return cmpError(fieldPath, reportValue(a), reportValue(b), fmt.Sprintf("%v != %v", a, b))
		}
	}
	return nil
}

// reportValue returns the value held by v for error reporting. Values obtained
// through unexported fields cannot be converted to interface{} and are
// reported as reflect.Value.
func reportValue(v reflect.Value) any {
	if v.IsValid() && v.CanInterface() {
		return v.Interface()
	}
	return v
}

// leafEqual compares scalar values without calling Interface(), so that
// unexported fields can be compared as well.
func leafEqual(a, b reflect.Value) bool {
	switch a.Kind() {
	case reflect.Bool:
		return a.Bool() == b.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() == b.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return a.Uint() == b.Uint()
	case reflect.Float32, reflect.Float64:
		return a.Float() == b.Float()
	case reflect.Complex64, reflect.Complex128:
		return a.Complex() == b.Complex()
	case reflect.String:
		return a.String() == b.String()
	default:
		return a.Interface() == b.Interface()
	}
}

// AcceptFunc decides whether a value at fieldPath should be copied by
// CopyStruct.
type AcceptFunc func(fieldPath string, srcValue reflect.Value) bool

// CopyStruct copies values from src to dst when acceptFunc returns true.
// dst must be a pointer and both values must have the same underlying type.
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

// CopyStructAll copies all fields from src to dst.
func CopyStructAll(src interface{}, dst interface{}) error {
	return CopyStruct(src, dst, func(fieldPath string, srcValue reflect.Value) bool { return true })
}

// CopyStructSelected copies only fields whose path contains any selectedFilePaths value.
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

// CopyStructAllExcept copies all fields except exact paths listed in
// excludedFilePaths.
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
