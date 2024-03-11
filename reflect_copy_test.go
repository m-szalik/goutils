package goutils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCopyStruct(t *testing.T) {
	keys := make([]string, 0)
	stringValueForPtr := "my string as ptr"
	type otherStruc struct {
		Name string
	}
	type testObject struct {
		Name   string
		Age    int32
		Emails []string
		Ost    otherStruc
		Ptr    *string
		Except int
	}

	src := testObject{
		Name:   "John Doe",
		Age:    44,
		Emails: []string{"john@nowhere.com", "doe@nowhere.com"},
		Ost: otherStruc{
			Name: "John Rambo",
		},
		Ptr:    &stringValueForPtr,
		Except: 99,
	}
	dst := testObject{}
	err := CopyStruct(src, &dst, func(fieldPath string, srcValue reflect.Value) bool {
		t.Logf("TestCopyStruct path '%s' --> %+v", fieldPath, srcValue)
		keys = append(keys, fieldPath)
		if fieldPath == ".Except" {
			return false
		}
		return true
	})
	assert.NoError(t, err)
	assert.Equal(t, src.Name, dst.Name)
	assert.Equal(t, src.Age, dst.Age)
	assert.Equal(t, src.Ost.Name, dst.Ost.Name)
	assert.Equal(t, src.Emails, dst.Emails)
	assert.Equal(t, src.Ptr, dst.Ptr)
	assert.Equal(t, 0, dst.Except)
	assert.Equal(t, []string{"", ".Name", ".Age", ".Emails", ".Emails[0]", ".Emails[1]", ".Ost", ".Ost.Name", ".Ptr", ".Except"}, keys)
}
