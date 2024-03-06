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

// ExitOnErrorWithMessage exit the program if error occur
func ExitOnErrorWithMessage(err error, code int, message string, messageArgs ...interface{}) {
	if err != nil {
		msg := fmt.Sprintf(message, messageArgs) + ": " + err.Error()
		ExitNow(code, msg)
	}
}

// Env retrieves the value of the environment variable named by the key. If not defined then default value is returned.
func Env(name string, def string) string {
	s := os.Getenv(name)
	if s == "" {
		return def
	}
	return s
}

// EnvInt retrieves the value as int of the environment variable named by the key. If not defined then default value is returned.
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
