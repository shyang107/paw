package paw

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return fi.IsDir()
}

// IsPathExists return true that `path` is dir or false for not
func IsPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	}
	// path/to/whatever exists
	return true
}

// GetCurrPath get the current path
func GetCurrPath() string {
	// file, _ := exec.LookPath(os.Args[0])
	// path, _ := filepath.Abs(file)
	// index := strings.LastIndex(path, string(os.PathSeparator))
	// ret := path[:index]
	// return ret
	var abPath string
	_, fileName, _, ok := runtime.Caller(1)
	if ok {
		abPath = filepath.Dir(fileName)
	}
	return abPath
}

// MakeAll check path and create like as `make -p path`
func MakeAll(path string) error {
	// check
	if IsPathExists(path) {
		return nil
	}
	err := os.MkdirAll(path, 0711) // 0755
	if err != nil {
		return err
	}
	// check again
	if !IsPathExists(path) {
		return fmt.Errorf("Makeall: fail to create %q", path)
	}
	return nil
}
