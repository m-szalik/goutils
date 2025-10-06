package goutils

import "errors"

// JoinErrorHelper - helps to return JoinError
type JoinErrorHelper struct {
	errs []error
}

// Append append error
func (jeh *JoinErrorHelper) Append(errs ...error) {
	if errs == nil {
		return
	}
	for _, err := range errs {
		if err != nil {
			jeh.errs = append(jeh.errs, err)
		}
	}
}

// AsError - return nil if no error, single error if one error was appended, JoinError otherwise
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

// NewJoinErrorHelper - create new JoinErrorHelper
func NewJoinErrorHelper(errs ...error) *JoinErrorHelper {
	if errs == nil {
		errs = make([]error, 0)
	}
	return &JoinErrorHelper{
		errs: errs,
	}
}
