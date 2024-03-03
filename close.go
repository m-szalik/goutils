package goutils

import "io"

// CloseQuietly close quietly. Ignore an error.
func CloseQuietly(closer io.Closer) {
	_ = closer.Close()
}
