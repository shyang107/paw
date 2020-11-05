package paw

import (
	"os"
)

// IsFileExist return true that `fileName` exist or false for not exist
func IsFileExist(fileName string) bool {
	fi, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !fi.IsDir()
}

// IsDirExists return true that `dir` is dir or false for not
func IsDirExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	}
	// path/to/whatever exists
	return true
}
