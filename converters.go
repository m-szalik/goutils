package goutils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var parseValueIntRegEx, _ = regexp.Compile("^-?\\d+$")
var parseValueFloatRegEx, _ = regexp.Compile("^-?\\d+\\.\\d+$")

func BoolToStr(b bool, trueVal, falseVal string) string {
	return BoolTo(b, trueVal, falseVal)
}

func BoolTo[T interface{}](b bool, trueVal, falseVal T) T {
	if b {
		return trueVal
	} else {
		return falseVal
	}
}

func HexToInt(hex string) (int, error) {
	hex = strings.Replace(hex, "0x", "", -1)
	hex = strings.Replace(hex, "0X", "", -1)
	i, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// ParseValue - returns one of int64, flot64, string, bool, nil
func ParseValue(str string) interface{} {
	s := strings.ToLower(strings.Trim(str, " \t"))
	switch s {
	case "true":
		return true
	case "false":
		return false
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
