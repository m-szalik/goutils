package goutils

import (
	"log"
	"os"
	"strconv"
)

var logger = log.New(os.Stderr, "", 0)

func ExitNow(code int, message string, messageArgs ...interface{}) {
	logger.Printf(message, messageArgs...)
	os.Exit(code)
}

func ExitOnError(err error, code int, message string, messageArgs ...interface{}) {
	if err != nil {
		ExitNow(code, message, messageArgs...)
	}
}

func Env(name string, def string) string {
	s := os.Getenv(name)
	if s == "" {
		return def
	}
	return s
}

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
