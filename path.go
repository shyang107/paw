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
	Log.Info("get home dir")
	home, _ := homedir.Dir()
	return home
}

// MakeAll check path and create like as `make -p path`
func MakeAll(path string) error {
	// check
	if IsPathExists(path) {
		return nil
	}
	err := os.MkdirAll(path, os.ModePerm) // 0755
	if err != nil {
		return err
	}
	// check again
	if !IsPathExists(path) {
		return fmt.Errorf("Makeall: fail to create %q", path)
	}
	return nil
}
