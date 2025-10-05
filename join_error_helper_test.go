package goutils

import (
	"errors"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestJoinErrorHelper(t *testing.T) {
	singleError := errors.New("single error")
	tests := []struct {
		name     string
		errors   []error
		expected error
	}{
		{
			name:     "no errors appended, expect nil",
			errors:   []error{},
			expected: nil,
		},
		{
			name:     "one error appended, expect err_single",
			errors:   []error{singleError},
			expected: singleError,
		},
		{
			name:     "one valid error appended, expect err_single",
			errors:   []error{nil, singleError, nil},
			expected: singleError,
		},
		{
			name:     "multiple errors appended, expect JoinError",
			errors:   []error{errors.New("0"), errors.New("1")},
			expected: errors.Join(errors.New("0"), errors.New("1")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jeh := &JoinErrorHelper{}
			jeh.Append(tt.errors...)
			asErr := jeh.AsError()
			assert2.Equal(t, tt.expected, asErr)
		})
	}
}
