package paw

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	// log "github.com/sirupsen/logrus"
)

// HasFile : Check if file exists in the current directory
func HasFile(filename string) bool {
	if info, err := os.Stat(filename); os.IsExist(err) {
		return !info.IsDir()
	}
	return false
}

// IsFileExist reports whether the named file exists as a boolean
func IsFileExist(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

// IsDirExist reports whether the dir exists as a boolean
func IsDirExist(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}
	return false
}

// IsExists reports whether the file or dir exists as a boolean
func IsExist(name string) bool {
	return IsDirExist(name) || IsFileExist(name)
}

// IsPathExist return whether the path exists.
func IsPathExist(path string) bool {
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

// GetAppDir get the current app directory
func GetAppDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"_os_Args_0": os.Args[0],
		}).Warn(err)

	}
	// Logger.Debugln(dir)
	return dir
}

// GetDotDir return the absolute path of "."
func GetDotDir() string {
	// w, _ := homedir.Expand(".")
	w, _ := filepath.Abs(".")
	// Logger.Debugln("get dot working dir", w)
	return w
}

// GetHomeDir get the home directory of user
func GetHomeDir() string {
	// Log.Info("get home dir")
	home, err := homedir.Dir()
	if err != nil {
		Logger.Error(err)
	}
	return home
}

// MakeAll check path and create like as `make -p path`
func MakeAll(path string) error {
	// check
	if IsPathExist(path) {
		return nil
	}
	err := os.MkdirAll(path, os.ModePerm) // 0755
	if err != nil {
		return err
	}
	// check again
	if !IsPathExist(path) {
		return fmt.Errorf("Makeall: fail to create %q", path)
	}
	return nil
}

// GetPathFromLink return path of symblink.
// 	If
//  1. path is not a link or effective path, return ""
// 	2. there is error, return error
func GetPathFromLink(path string) string {
	info, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return ""
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return ""
	} else {
		alink, err := os.Readlink(path)
		if err != nil {
			return err.Error()
		}
		return alink
	}
}
