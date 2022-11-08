package goutils

import (
	"fmt"
	assert2 "github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_ParseValue(t *testing.T) {
	tests := []struct {
		arg  string
		want interface{}
	}{
		{
			arg:  "true",
			want: true,
		},
		{
			arg:  "nil",
			want: nil,
		},
		{
			arg:  "null ",
			want: nil,
		},
		{
			arg:  " ",
			want: " ",
		},
		{
			arg:  "512",
			want: int64(512),
		},
		{
			arg:  "-512",
			want: int64(-512),
		},
		{
			arg:  " some text ",
			want: " some text ",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("convertValue_%s", tt.arg), func(t *testing.T) {
			got := ParseValue(tt.arg)
			if got != tt.want {
				t.Errorf("convertValue() got = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("convertValue_floats", func(t *testing.T) {
		for _, f := range []float64{-19.01, 193.834} {
			t.Run(fmt.Sprintf("convert_%f", f), func(t *testing.T) {
				str := fmt.Sprint(f)
				value := ParseValue(str)
				assert2.NotNil(t, value)
				if math.Abs(value.(float64)-f) > 0.00001 {
					assert2.Equal(t, value, f)
				}
			})
		}
	})
}
