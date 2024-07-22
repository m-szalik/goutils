package goutils

import (
	"fmt"
	"reflect"
)

func CreateEmptyObjectOf(source any) any {
	if reflect.ValueOf(source).Kind() != reflect.Ptr {
		panic(fmt.Sprintf("source must be pointer to a type buf %T found", source))
	}
	val := reflect.ValueOf(source)
	result := reflect.Indirect(reflect.New(val.Type())).Interface()
	return result
}
