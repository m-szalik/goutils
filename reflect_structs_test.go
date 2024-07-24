package goutils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCopyStructDifferentStructs(t *testing.T) {
	type bType struct {
		Name string
	}
	type aType struct {
		Age int32
	}
	err := CopyStructAll(aType{Age: 32}, &bType{})
	assert.Error(t, err)
}

func TestCopyStructDstNotAPointer(t *testing.T) {
	type structType struct {
		Name string
	}
	err := CopyStructAll(structType{Name: "Rambo"}, structType{})
	assert.Error(t, err)
}

func TestCopyStruct(t *testing.T) {
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

	t.Run("copy with filter", func(t *testing.T) {
		keys := make([]string, 0)
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
	})

	t.Run("copy all except", func(t *testing.T) {
		dst := testObject{}
		err := CopyStructAllExcept(src, &dst, ".Except")
		assert.NoError(t, err)
		assert.Equal(t, src.Name, dst.Name)
		assert.Equal(t, src.Age, dst.Age)
		assert.Equal(t, src.Ost.Name, dst.Ost.Name)
		assert.Equal(t, src.Emails, dst.Emails)
		assert.Equal(t, src.Ptr, dst.Ptr)
		assert.Equal(t, 0, dst.Except)
	})

	t.Run("copy selected", func(t *testing.T) {
		dst := testObject{}
		err := CopyStructSelected(src, &dst, ".Name", ".Age")
		assert.NoError(t, err)
		assert.Equal(t, src.Name, dst.Name)
		assert.Equal(t, src.Age, dst.Age)
		assert.Equal(t, "", dst.Ost.Name)
		assert.Nil(t, dst.Emails)
		assert.Nil(t, dst.Ptr)
		assert.Equal(t, 0, dst.Except)
	})

}
