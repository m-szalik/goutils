package goutils

import (
	"strconv"
	"strings"
)

func BoolToStr(b bool, trueVal, falseVal string) string {
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
