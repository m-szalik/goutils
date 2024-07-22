package goutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateEmptyObjectOf(t *testing.T) {
	type myTestStruct struct {
		Name string
	}
	str := "Whoa"
	strArr := []string{}
	tests := []struct {
		arg any
	}{
		{arg: &myTestStruct{}},
		{arg: &str},
		{arg: &strArr},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("CreateEmptyObjectOf(type:%T)", tt.arg), func(t *testing.T) {
			got := CreateEmptyObjectOf(tt.arg)
			assert.Equal(t, fmt.Sprintf("%T", tt.arg), fmt.Sprintf("%T", got))
		})
	}
}
