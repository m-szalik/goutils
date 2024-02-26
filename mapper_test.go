package goutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_MapSlice(t *testing.T) {
	for _, vSlice := range [][]string{
		{"a", "shd ", "faaz"},
		{"xxkjd ", "!@#", "faz"},
		{"xxkjd ", "a(z@", ""},
		{},
	} {
		t.Run(fmt.Sprintf("MapSlice_%v", vSlice), func(t *testing.T) {
			res := MapSlice[string, string](vSlice, func(s string) string {
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
