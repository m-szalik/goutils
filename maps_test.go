package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKeys(t *testing.T) {
	t.Run("nil map", func(t *testing.T) {
		assert.Nil(t, MapKeys(nil))
	})

	t.Run("returns all keys", func(t *testing.T) {
		input := map[string]any{
			"one":   1,
			"two":   true,
			"three": "3",
		}

		got := MapKeys(input)
		assert.ElementsMatch(t, []string{"one", "two", "three"}, got)
	})
}

func TestMapMerge(t *testing.T) {
	t.Run("merges without conflicts", func(t *testing.T) {
		initMap := map[string]int{"a": 1}

		got := MapMerge(
			func(_ string, v0, v1 int) int {
				return v0 + v1
			},
			initMap,
			map[string]int{"b": 2},
			map[string]int{"c": 3},
		)

		assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, got)
		assert.Equal(t, initMap, got)
	})

	t.Run("resolves conflicts using conflictFunc", func(t *testing.T) {
		initMap := map[string]int{"a": 1, "b": 2}
		calls := make([]string, 0)

		got := MapMerge(
			func(key string, v0, v1 int) int {
				calls = append(calls, key)
				return v0 + v1
			},
			initMap,
			map[string]int{"a": 3, "c": 4},
			map[string]int{"a": 5, "b": 6},
			nil,
		)

		assert.Equal(t, map[string]int{"a": 9, "b": 8, "c": 4}, got)
		assert.Equal(t, []string{"a", "a", "b"}, calls)
	})
}

func TestMapMergeOverride(t *testing.T) {
	t.Run("nil initMap", func(t *testing.T) {
		got := MapMergeOverride[string, int](nil, map[string]int{"a": 1})
		assert.Equal(t, map[string]int{"a": 1}, got)
	})

	t.Run("incoming value wins", func(t *testing.T) {
		got := MapMergeOverride(map[string]int{"a": 1, "b": 2}, map[string]int{"a": 3}, map[string]int{"a": 4, "c": 5})
		assert.Equal(t, map[string]int{"a": 4, "b": 2, "c": 5}, got)
	})
}
