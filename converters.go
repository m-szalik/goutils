package goutils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var parseValueIntRegEx, _ = regexp.Compile("^-?\\d+$")
var parseValueFloatRegEx, _ = regexp.Compile("^-?\\d+\\.\\d+$")

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
	hex = strings.Replace(hex, "0x", "", -1)
	hex = strings.Replace(hex, "0X", "", -1)
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
// If impossible to covert the same string is returned.
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
func AsFloat64(i interface{}) (float64, error) {
	if i == nil {
		return 0, fmt.Errorf("cannot convert nil to float64")
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
