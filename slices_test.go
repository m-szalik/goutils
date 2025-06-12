package goutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
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

func Test_SliceMap(t *testing.T) {
	for _, vSlice := range [][]string{
		{"a", "shd ", "faaz"},
		{"xxkjd ", "!@#", "faz"},
		{"xxkjd ", "a(z@", ""},
		{},
	} {
		t.Run(fmt.Sprintf("MapSlice_%v", vSlice), func(t *testing.T) {
			res := SliceMap[string, string](vSlice, func(s string) string {
				return strings.ToUpper(s)
			})
			assert.Equal(t, len(vSlice), len(res))
			for i := 0; i < len(vSlice); i++ {
				upper := strings.ToUpper(vSlice[i])
				assert.Equal(t, upper, res[i])
			}
		})
	}
}

func TestSliceMapIntToString(t *testing.T) {
	sliceOfStrings := SliceMap[int, string]([]int{2, 7, -11}, func(i int) string { return fmt.Sprint(i) })
	assert.Equal(t, []string{"2", "7", "-11"}, sliceOfStrings)
}

func TestSlicesEq(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		a    []T
		b    []T
		want bool
	}
	tests := []testCase[string]{
		{name: "not equal by size", a: []string{"a"}, b: []string{"a", "b"}, want: false},
		{name: "not equal by content", a: []string{"a", "b"}, b: []string{"a", "z"}, want: false},
		{name: "equal", a: []string{"a", "b"}, b: []string{"a", "b"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SlicesEq(tt.a, tt.b), "SlicesEq(%v, %v)", tt.a, tt.b)
		})
	}
}

func TestDistrictValues(t *testing.T) {
	type testCase[T comparable] struct {
		name  string
		input []T
		want  []T
	}
	tests := []testCase[string]{
		{name: "nil value", input: nil, want: nil},
		{name: "empty slice", input: []string{}, want: []string{}},
		{name: "all unique element", input: []string{"a", "b"}, want: []string{"a", "b"}},
		{name: "some duplicates element", input: []string{"a", "b", "b", "c", "a"}, want: []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, DistrictValues(tt.input), "DistrictValues(%v)", tt.input)
		})
	}
}
