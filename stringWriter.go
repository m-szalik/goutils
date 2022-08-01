package goutils

import (
	"fmt"
	"io"
)

type stringWriter struct {
	buff []byte
}

func (s *stringWriter) Write(p []byte) (n int, err error) {
	s.buff = append(s.buff, p...)
	return len(p), nil
}

func (s *stringWriter) String() string {
	return string(s.buff)
}

type StringWriter interface {
	io.Writer
	fmt.Stringer
}

func NewStringWriter() StringWriter {
	return &stringWriter{
		buff: make([]byte, 0),
	}
}
