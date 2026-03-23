package goutils

import "io"

// CloseQuietly closes closer and ignores any returned error.
func CloseQuietly(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}
