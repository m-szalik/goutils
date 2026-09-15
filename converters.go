package goutils

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var parseValueIntRegEx = regexp.MustCompile(`^-?\d+$`)
var parseValueFloatRegEx = regexp.MustCompile(`^-?\d+\.\d+$`)

// BoolToStr returns trueVal when b is true, otherwise falseVal.
func BoolToStr(b bool, trueVal, falseVal string) string {
	return BoolTo(b, trueVal, falseVal)
}

// BoolTo returns trueVal when b is true, otherwise falseVal.
func BoolTo[T interface{}](b bool, trueVal, falseVal T) T {
	if b {
		return trueVal
	} else {
		return falseVal
	}
}

// HexToInt converts a hexadecimal string into int.
// Both "0x" and "0X" prefixes are accepted.
func HexToInt(hex string) (int, error) {
	sign, digits := "", hex
	if strings.HasPrefix(digits, "-") || strings.HasPrefix(digits, "+") {
		sign, digits = digits[:1], digits[1:]
	}
	if strings.HasPrefix(digits, "0x") || strings.HasPrefix(digits, "0X") {
		digits = digits[2:]
	}
	i, err := strconv.ParseInt(sign+digits, 16, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// ParseBool parses string values into booleans.
// True values: "true", "1", "on". False values: "false", "0", "off".
func ParseBool(str string) (bool, error) {
	s := strings.ToLower(strings.TrimSpace(str))
	switch s {
	case "true", "1", "on":
		return true, nil
	case "false", "0", "off":
		return false, nil
	default:
		return false, fmt.Errorf("unable to parse '%s' as boolean", str)
	}
}

// ParseValue converts a string to bool, nil, int64, float64, or leaves it as
// string when conversion is not possible.
func ParseValue(str string) interface{} {
	if b, err := ParseBool(str); err == nil {
		return b
	}
	s := strings.ToLower(strings.TrimSpace(str))
	switch s {
	case "null", "nil":
		return nil
	default:
		var r interface{}
		var err error
		matchAny := false
		if parseValueIntRegEx.MatchString(s) {
			r, err = strconv.ParseInt(s, 10, 64)
			matchAny = true
		}
		if parseValueFloatRegEx.MatchString(s) {
			r, err = strconv.ParseFloat(s, 64)
			matchAny = true
		}
		if matchAny && err == nil {
			return r
		}
	}
	return str
}

// AsFloat64 converts input into float64.
// Supported types: float32, float64, int, int32, int64, string, []byte, and
// pointers to these types.
func AsFloat64(input any) (float64, error) {
	if input == nil {
		return 0, fmt.Errorf("cannot convert nil to float64")
	}
	var i any
	rv := reflect.ValueOf(input)
	if rv.Kind() == reflect.Ptr {
		i = rv.Elem().Interface()
	} else {
		i = input
	}
	switch v := i.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case []byte:
		f, err := strconv.ParseFloat(strings.TrimSpace(string(v)), 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert '%s' to float64 - %w", v, err)
		}
		return f, nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert '%s' to float64 - %w", v, err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("unable to covert type %T to float64", i)
	}
}

// RoundFloat rounds val to the given decimal precision.
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func envConvertion[T string | int | bool | float32 | float64 | time.Duration](value string) (*T, error) {
	var zero T
	switch any(zero).(type) {
	case string:
		v := any(value).(T)
		return &v, nil
	case int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		v := any(intVal).(T)
		return &v, nil
	case bool:
		b, err := ParseBool(value)
		if err != nil {
			return nil, err
		}
		v := any(b).(T)
		return &v, nil
	case float32:
		f64, err := AsFloat64(value)
		if err != nil {
			return nil, err
		}
		v := any(float32(f64)).(T)
		return &v, nil
	case float64:
		f64, err := AsFloat64(value)
		if err != nil {
			return nil, err
		}
		v := any(f64).(T)
		return &v, nil
	case time.Duration:
		dur, err := time.ParseDuration(value)
		if err != nil {
			return nil, fmt.Errorf("cannot parse duration '%s':: %w", value, err)
		}
		v := any(dur).(T)
		return &v, nil
	default:
		// Shouldn’t happen given the constraint, but keep a safe fallback
		return nil, fmt.Errorf("unsupported type %T has been passed", zero)
	}
}
