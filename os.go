package goutils

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", 0)

func ExitNow(code int, message string, messageArgs ...interface{}) {
	logger.Fatalf(message, messageArgs...)
	os.Exit(code)
}


func ExitOnError(err error, code int, message string, messageArgs ...interface{}) {
	if err != nil {
		ExitNow(code, message, messageArgs...)
	}
}

