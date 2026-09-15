package goutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testStruct struct {
	SArray []string
	Num    int
	unexp  int
}

type iterationResult struct {
	Path  string
	Value any
}

func TestIterateDeep(t *testing.T) {
	var result []iterationResult
	callback := func(path string, depth int, kind reflect.Kind, element interface{}) bool {
		fmt.Printf("-> %s %d -> %s : %v\n", path, depth, kind, element)
		result = append(result, iterationResult{
			Path:  path,
			Value: element,
		})
		return true
	}
	tests := []struct {
		startElement interface{}
		expected     []iterationResult
	}{
		{startElement: "A", expected: []iterationResult{{".", "A"}}},
		{startElement: []string{"a", "b", "c"}, expected: []iterationResult{{".[0]", "a"}, {".[1]", "b"}, {".[2]", "c"}}},
		{startElement: &testStruct{[]string{"ala", "ma"}, 12, -11}, expected: []iterationResult{{".SArray[0]", "ala"}, {".SArray[1]", "ma"}, {".Num", 12}}},
		{startElement: map[string]interface{}{"A": 1, "B": 1.2}, expected: []iterationResult{{".A", 1}, {".B", 1.2}}},
		{startElement: struct {
			P *int
			V int
		}{}, expected: []iterationResult{{".P", nil}, {".V", 0}}},
		{startElement: map[string]any{"b": map[string]any{"c": []any{1, nil}}, "a": 2}, expected: []iterationResult{{".a", 2}, {".b.c[0]", 1}, {".b.c[1]", nil}}},
		{startElement: map[int]string{7: "x"}, expected: []iterationResult{{".7", "x"}}},
		{startElement: nil, expected: []iterationResult{{".", nil}}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d_iterateDeep_over-%v", i, tt.startElement), func(t *testing.T) {
			result = make([]iterationResult, 0)
			IterateDeep(tt.startElement, callback)
			assert.Equal(t, tt.expected, result)
		})
	}
}
