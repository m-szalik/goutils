package goutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var sliceA = []interface{}{"A", "b", 11, 3.14}

func TestSliceIndexOf(t *testing.T) {
	type args struct {
		slice []interface{}
		e     interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "NotFound", args: args{slice: sliceA, e: "bar"}, want: -1},
		{name: "Zero", args: args{slice: sliceA, e: "A"}, want: 0},
		{name: "3.14", args: args{slice: sliceA, e: 3.14}, want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceIndexOf(tt.args.slice, tt.args.e); got != tt.want {
				t.Errorf("SliceIndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceContains(t *testing.T) {
	type args struct {
		slice []interface{}
		e     interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "NotFound", args: args{slice: sliceA, e: "bar"}, want: false},
		{name: "Zero", args: args{slice: sliceA, e: "A"}, want: true},
		{name: "3.14", args: args{slice: sliceA, e: 3.14}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceContains(tt.args.slice, tt.args.e); got != tt.want {
				t.Errorf("SliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceRemove(t *testing.T) {
	type testCase[T comparable] struct {
		name        string
		slice       []T
		e           any
		want        []T
		wantRemoved int
		wantLen     int
	}
	testsStrings := []testCase[string]{
		{name: "remove multiple strings", slice: []string{"abc", "foo", "bar", "go", "foo"}, e: "foo", want: []string{"abc", "bar", "go"}, wantRemoved: 2, wantLen: 3},
		{name: "not remove multiple strings", slice: []string{"abc", "foo", "bar", "go", "foo"}, e: "xyz", want: []string{"abc", "foo", "bar", "go", "foo"}, wantRemoved: 0, wantLen: 5},
	}
	testsInts := []testCase[int]{
		{name: "remove multiple int", slice: []int{0, 1, 2, 3, 0, 5, 0}, e: 0, want: []int{1, 2, 3, 5}, wantRemoved: 3, wantLen: 4},
	}

	t.Run("strings", func(t *testing.T) {
		for _, tt := range testsStrings {
			t.Run(tt.name, func(t *testing.T) {
				got, got1 := SliceRemove(tt.slice, tt.e)
				assert.Equalf(t, tt.want, got, "SliceRemove.slice(%v, %v)", tt.slice, tt.e)
				assert.Equalf(t, tt.wantRemoved, got1, "SliceRemove.removed(%v, %v)", tt.slice, tt.e)
				assert.Equalf(t, tt.wantLen, len(got), "SliceRemove.len(%v, %v)", tt.slice, tt.e)
			})
		}
	})
	t.Run("ints", func(t *testing.T) {
		for _, tt := range testsInts {
			t.Run(tt.name, func(t *testing.T) {
				got, got1 := SliceRemove(tt.slice, tt.e)
				assert.Equalf(t, tt.want, got, "SliceRemove.slice(%v, %v)", tt.slice, tt.e)
				assert.Equalf(t, tt.wantRemoved, got1, "SliceRemove.removed(%v, %v)", tt.slice, tt.e)
				assert.Equalf(t, tt.wantLen, len(got), "SliceRemove.len(%v, %v)", tt.slice, tt.e)
			})
		}
	})

}
