package goutils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Env retrieves the value of the environment variable named by the key. If not defined then default value is returned.
// Supported types: string,int,bool,float32,float64,time.Duration
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

// EnvRequired retrieves the value of the environment variable named by the key. If not defined then panic.
// Supported types: string,int,bool,float32,float64,time.Duration
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

// EnvInt retrieves the value as int of the environment variable named by the key. If not defined then default value is returned.
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
