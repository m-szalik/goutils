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

// BoolToStr - return string for true or false bool value
func BoolToStr(b bool, trueVal, falseVal string) string {
	return BoolTo(b, trueVal, falseVal)
}

// BoolTo - return T object for true or false bool value
func BoolTo[T interface{}](b bool, trueVal, falseVal T) T {
	if b {
		return trueVal
	} else {
		return falseVal
	}
}

// HexToInt convert hex representation to int
func HexToInt(hex string) (int, error) {
	hex = strings.Replace(hex, "0x", "", -1) //nolint:staticcheck
	hex = strings.Replace(hex, "0X", "", -1) //nolint:staticcheck
	i, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// ParseBool - return bool. True is one of "true", "1", "on", false is one of "false", "0", "off"
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

// ParseValue - converts string to one of int64, flot64, string, bool, nil.
// If impossible to convert, the same string is returned.
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
			r, err = strconv.ParseInt(s, 10, 32)
			matchAny = true
		}
		if parseValueFloatRegEx.MatchString(s) {
			r, err = strconv.ParseFloat(s, 32)
			matchAny = true
		}
		if matchAny {
			if err != nil {
				panic(fmt.Sprintf("error parsing '%s' as number", s))
			}
			return r
		}
	}
	return str
}

// AsFloat64 - convert multiple types (float32,int,int32,int64,string,[]byte) to float64
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

// RoundFloat round float64 number
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
