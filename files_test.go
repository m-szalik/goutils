package goutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	tests := []struct {
		filePath string
		want     bool
	}{
		{"notExistingFile.txt", false},
		{"README.md", true},
		{".github/workflows/go.yml", true},
		{".github/workflows", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Checking file %s expecting %t", tt.filePath, tt.want), func(t *testing.T) {
			assert.Equalf(t, tt.want, FileExists(tt.filePath), "FileExists(%v)", tt.filePath)
		})
	}
}
