package goutils

import "errors"

// JoinErrorHelper - helps to return JoinError
type JoinErrorHelper struct {
	errs []error
}

// Append append error
func (jeh *JoinErrorHelper) Append(err error) {
	if err != nil {
		jeh.errs = append(jeh.errs, err)
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
