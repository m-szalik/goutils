package goutils

import "io"

// CloseQuietly close quietly. Ignore an error.
func CloseQuietly(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}
