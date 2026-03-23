package goutils

import "errors"

// JoinErrorHelper accumulates errors and exposes them as a single error value.
type JoinErrorHelper struct {
	errs []error
}

// ErrorsCount returns the number of collected non-nil errors.
func (jeh *JoinErrorHelper) ErrorsCount() int {
	return len(jeh.errs)
}

// Append adds non-nil errors to the helper and returns the helper for chaining.
func (jeh *JoinErrorHelper) Append(errs ...error) *JoinErrorHelper {
	if errs == nil {
		return jeh
	}
	for _, err := range errs {
		if err != nil {
			jeh.errs = append(jeh.errs, err)
		}
	}
	return jeh
}

// AsError returns nil when no errors were collected, the single error when one
// was collected, or errors.Join for multiple errors.
func (jeh *JoinErrorHelper) AsError() error {
	switch len(jeh.errs) {
	case 0:
		return nil
	case 1:
		return jeh.errs[0]
	default:
		return errors.Join(jeh.errs...)
	}
}


func (jeh *JoinErrorHelper) Iterate() <-chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		for _, err := range jeh.errs {
			ch <- err
		}
	}()
	return ch
}

// NewJoinErrorHelper creates a helper optionally initialized with errs.
func NewJoinErrorHelper(errs ...error) *JoinErrorHelper {
	if errs == nil {
		errs = make([]error, 0)
	}
	return &JoinErrorHelper{
		errs: errs,
	}
}
