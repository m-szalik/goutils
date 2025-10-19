package goutils

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", 0)

// ExitNow exit now
// It supports `/dev/termination-log' file.
func ExitNow(code int, message string, messageArgs ...interface{}) {
	fullMessage := fmt.Sprintf(message, messageArgs...)
	logger.Print(fullMessage)
	terminationFilename := Env("TERMINATION_MESSAGE_PATH", "/dev/termination-log")
	if terminationFilename != "" {
		_, err := os.Stat(terminationFilename)
		if err != nil {
			if !os.IsNotExist(err) {
				logger.Printf("cannot create termination log file '%s'::%s", terminationFilename, err)
			}
		} else {
			f, err := os.OpenFile(terminationFilename, os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				func() {
					defer CloseQuietly(f)
					_, _ = f.WriteString(fullMessage)
				}()
			}
		}
	}
	os.Exit(code)
}

// ExitOnError exit the program if error occur
func ExitOnError(err error, code int) {
	if err != nil {
		ExitNow(code, "error:: %s", err)
	}
}

// ExitOnErrorf exit the program if error occur
func ExitOnErrorf(err error, code int, message string, messageArgs ...interface{}) {
	if err != nil {
		msg := fmt.Sprintf(message, messageArgs...) + ":: " + err.Error()
		ExitNow(code, msg)
	}
}
