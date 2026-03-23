package goutils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Env returns the value of name converted to T, or def when the variable is
// missing or cannot be converted.
// Supported types: string, int, bool, float32, float64, time.Duration.
func Env[T string | int | bool | float32 | float64 | time.Duration](name string, def T) T {
	s := os.Getenv(name)
	if s == "" { // not set
		return def
	}
	v, err := envConvertion[T](s)
	if err != nil {
		return def
	}
	return *v
}

// EnvRequired returns the value of name converted to T.
// It panics when the variable is missing or conversion fails.
// Supported types: string, int, bool, float32, float64, time.Duration.
func EnvRequired[T string | int | bool | float32 | float64 | time.Duration](name string) T {
	val := os.Getenv(name)
	if val == "" {
		panic(fmt.Sprintf("enviroment variable %s not defined", name))
	}
	v, err := envConvertion[T](val)
	if err != nil {
		panic(fmt.Sprintf("cannot convert enviroment variable %s value '%s' :: %s", name, val, err))
	}
	return *v
}

// EnvInt returns the value of name as int, or def on missing value or parse
// failure.
// Deprecated: use: [Env].
func EnvInt(name string, def int) int {
	s := os.Getenv(name)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
