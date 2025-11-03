package goutils

import "os"

// FileExists check if file exists
func FileExists(filePath string) bool {
	stat, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

// DirExists check if directory exists
func DirExists(filePath string) bool {
	stat, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return stat.IsDir()
}
