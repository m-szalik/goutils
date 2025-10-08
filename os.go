package goutils

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var logger = log.New(os.Stderr, "", 0)

func ExitNow(code int, message string, messageArgs ...interface{}) {
	logger.Printf(message, messageArgs...)
	os.Exit(code)
}

// ExitOnError exit the program if error occur
func ExitOnError(err error, code int) {
	if err != nil {
		ExitNow(code, "error: %s", err)
	}
}

// ExitOnErrorf exit the program if error occur
func ExitOnErrorf(err error, code int, message string, messageArgs ...interface{}) {
	if err != nil {
		msg := fmt.Sprintf(message, messageArgs...) + ": " + err.Error()
		ExitNow(code, msg)
	}
}

// Env retrieves the value of the environment variable named by the key. If not defined then default value is returned.
// Supported types: string,int,bool,float32,float64
func Env[T string | int | bool | float32 | float64](name string, def T) T {
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
// Supported types: string,int,bool,float32,float64
func EnvRequired[T string | int | bool | float32 | float64](name string) T {
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

func envConvertion[T string | int | bool | float32 | float64](value string) (*T, error) {
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
	default:
		// Shouldn’t happen given the constraint, but keep a safe fallback
		return nil, fmt.Errorf("unsupported type %T has been passed", zero)
	}
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
