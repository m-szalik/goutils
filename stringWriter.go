package goutils

import (
	"fmt"
	"io"
)

type stringWriter struct {
	str string
}

func (s *stringWriter) Write(p []byte) (n int, err error) {
	s.str = fmt.Sprintf("%s%s", s.str, string(p))
	return len(p), nil
}

func (s *stringWriter) String() string {
	return s.str
}

type StringWriter interface {
	io.Writer
	fmt.Stringer
}

func NewStringWriter() StringWriter {
	return &stringWriter{}
}
