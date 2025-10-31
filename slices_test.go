package goutils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

var allCapitalCasesMatch = func(s string) bool {
	return strings.ToUpper(s) == s
}

func TestAllMatch(t *testing.T) {
	type testCase[T any] struct {
		name      string
		input     []T
		condition func(element T) bool
		want      bool
	}
	tests := []testCase[string]{
		{
			name:      "all matches",
			input:     []string{"JOHN", "ALEX"},
			condition: allCapitalCasesMatch,
			want:      true,
		},
		{
			name:      "all matches for empty input",
			input:     []string{},
			condition: allCapitalCasesMatch,
			want:      true,
		},
		{
			name:      "not all matches",
			input:     []string{"JOHN", "Rene", "FRANK"},
			condition: allCapitalCasesMatch,
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SliceAllMatch(tt.input, tt.condition), "SliceAllMatch(%v, %T)", tt.input, tt.condition)
		})
	}
}

func TestAnyMatch(t *testing.T) {
	type testCase[T any] struct {
		name      string
		input     []T
		condition func(element T) bool
		want      bool
	}
	tests := []testCase[string]{
		{
			name:      "one matches",
			input:     []string{"Frank", "ALEX", "Eric"},
			condition: allCapitalCasesMatch,
			want:      true,
		},
		{
			name:      "empty input",
			input:     []string{},
			condition: allCapitalCasesMatch,
			want:      false,
		},
		{
			name:      "one all matches",
			input:     []string{"John", "Rene", "FRANK"},
			condition: allCapitalCasesMatch,
			want:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SliceAnyMatch(tt.input, tt.condition), "SliceAnyMatch(%v, %T)", tt.input, tt.condition)
		})
	}
}

func TestCountMatch(t *testing.T) {
	type testCase[T any] struct {
		name      string
		input     []T
		condition func(element T) bool
		want      int
	}
	tests := []testCase[string]{
		{
			name:      "one matches",
			input:     []string{"Frank", "ALEX", "Eric"},
			condition: allCapitalCasesMatch,
			want:      1,
		},
		{
			name:      "empty input",
			input:     []string{},
			condition: allCapitalCasesMatch,
			want:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SliceCountMatch(tt.input, tt.condition), "SliceCountMatch(%v, %T)", tt.input, tt.condition)
		})
	}
}

func TestFindFirst(t *testing.T) {
	isEven := func(i int) bool { return i%2 == 0 }
	// found
	in1 := []int{1, 3, 4, 6}
	ptr := FindFirst(in1, isEven)
	if assert.NotNil(t, ptr, "expected non-nil pointer when a match exists") && ptr != nil {
		assert.Equal(t, 4, *ptr)
	}

	// not found
	in2 := []int{1, 3, 5}
	assert.Nil(t, FindFirst(in2, isEven))

	// empty input
	var in3 []int
	assert.Nil(t, FindFirst(in3, isEven))
}

func TestFilter(t *testing.T) {
	isUpper := func(s string) bool { return strings.ToUpper(s) == s }

	// typical filtering
	in1 := []string{"JOHN", "Alex", "MIKE", "rene"}
	out1 := Filter(in1, isUpper)
	assert.Equal(t, []string{"JOHN", "MIKE"}, out1)

	// empty slice -> empty slice
	in2 := []string{}
	out2 := Filter(in2, isUpper)
	assert.Equal(t, 0, len(out2))
	assert.NotNil(t, out2)

	// nil slice -> empty slice (not nil)
	var in3 []string
	out3 := Filter(in3, isUpper)
	assert.Equal(t, 0, len(out3))
	assert.NotNil(t, out3)
}

func TestAllMatch_Function(t *testing.T) {
	allCaps := func(s string) bool { return strings.ToUpper(s) == s }

	assert.True(t, AllMatch([]string{"JOHN", "ALEX"}, allCaps))
	assert.True(t, AllMatch([]string{}, allCaps)) // vacuously true
	assert.False(t, AllMatch([]string{"JOHN", "Rene", "FRANK"}, allCaps))
}

func TestAnyMatch_Function(t *testing.T) {
	allCaps := func(s string) bool { return strings.ToUpper(s) == s }

	assert.True(t, AnyMatch([]string{"Frank", "ALEX", "Eric"}, allCaps))
	assert.False(t, AnyMatch([]string{}, allCaps))
	assert.True(t, AnyMatch([]string{"John", "Rene", "FRANK"}, allCaps))
	assert.False(t, AnyMatch([]string{"john", "rene"}, allCaps))
}

func TestCountMatch_Function(t *testing.T) {
	ge5 := func(i int) bool { return i >= 5 }

	assert.Equal(t, 3, CountMatch([]int{1, 5, 7, 3, 9}, ge5))
	assert.Equal(t, 0, CountMatch([]int{1, 2, 3}, ge5))
	assert.Equal(t, 0, CountMatch([]int{}, ge5))
}
