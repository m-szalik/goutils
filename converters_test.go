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
		t.Run(fmt.Sprintf("ParseValue_%s", tt.arg), func(t *testing.T) {
			got := ParseValue(tt.arg)
			if got != tt.want {
				t.Errorf("ParseValue() got = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("ParseValue_floats", func(t *testing.T) {
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

func Test_ParseBool(t *testing.T) {
	tests := []struct {
		arg         string
		want        bool
		expectError bool
	}{
		{
			arg:  "true",
			want: true,
		},
		{
			arg:  " 0",
			want: false,
		},
		{
			arg:  "\toff ",
			want: false,
		},
		{
			arg:  "TrUE",
			want: true,
		},
		{
			arg:  " on ",
			want: true,
		},
		{
			arg:         "blabla",
			want:        false,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("ParseBool_%s", tt.arg), func(t *testing.T) {
			got, err := ParseBool(tt.arg)
			assert2.Equal(t, tt.want, got)
			if tt.expectError {
				if err == nil {
					t.Errorf("ParseBool() error expected")
				}
			} else {
				assert2.NoError(t, err)
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

func TestAsFloat64(t *testing.T) {
	tests := []struct {
		arg     interface{}
		want    float64
		wantErr bool
	}{
		{"-17.2", -17.2, false},
		{"17.2", 17.2, false},
		{" 17.2 ", 17.2, false},
		{"17", 17, false},
		{17, 17, false},
		{float32(17), 17, false},
		{int64(17), 17, false},
		{int32(17), 17, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("converting %v of type %T", tt.arg, tt.arg), func(t *testing.T) {
			got, err := AsFloat64(tt.arg)
			if tt.wantErr && err == nil {
				assert2.Fail(t, "error expected")
			}
			assert2.Equalf(t, tt.want, got, "AsFloat64(%v)", tt.arg)
		})
	}
}
func TestRoundFloat64(t *testing.T) {
	tests := []struct {
		arg       float64
		precision uint
		want      float64
	}{
		{1.0 / 3.0, 9, 0.333333333},
		{1.0 / 3.0, 2, 0.33},
		{1.0 / 3.0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("RoundFloat(%f, %d) = %f", tt.arg, tt.precision, tt.want), func(t *testing.T) {
			result := RoundFloat(tt.arg, tt.precision)
			assert2.Equal(t, tt.want, result)
		})
	}
}

func TestHexToInt(t *testing.T) {
	tests := []struct {
		args    string
		want    int
		wantErr assert2.ErrorAssertionFunc
	}{
		{"A", 10, assert2.NoError},
		{"0xA", 10, assert2.NoError},
		{"0X0A", 10, assert2.NoError},
		{"0B", 11, assert2.NoError},
		{"-B", -11, assert2.NoError},
		{"-b", -11, assert2.NoError},
		{"-0xb", -11, assert2.NoError},
		{"invalid", 0, assert2.Error},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("HexToInt(\"%s\") is %d", tt.args, tt.want), func(t *testing.T) {
			got, err := HexToInt(tt.args)
			if !tt.wantErr(t, err, fmt.Sprintf("HexToInt(%v)", tt.args)) {
				return
			}
			assert2.Equalf(t, tt.want, got, "HexToInt(%v)", tt.args)
		})
	}
}
