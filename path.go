package paw

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	// log "github.com/sirupsen/logrus"
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

// File : path information
//
// 	Fields:
// 	  `FullPath`: The full path including the folder
// 	  `Folder`: The folder of the file
// 	  `File`: The file name including extension (basename)
// 	  `FileName`: The file name excluding extension
// 	  `Ext`: Extension of the file
type File struct {
	FullPath string // The full path including the folder
	Folder   string // The folder of the file
	File     string // The file name including extension (basename)
	FileName string // The file name excluding extension
	Ext      string // Extension of the file
}

// ConstructFile construct `paw.File` from string
//
// Example:
// 	path := "/aaa/bbb/ccc/example.xxx"
// 	path => File{
// 		FullPath: "/aaa/bbb/ccc/example.xxx",
// 		File:     "example.xxx",
// 		Folder:   "/aaa/bbb/ccc/",
// 		FileName: "example",
// 		Ext:      ".xxx",
// 	}
func ConstructFile(path string) File {
	base := filepath.Base(path)
	ext := filepath.Ext(path)

	return File{
		FullPath: path,
		File:     base,
		Folder:   strings.TrimSuffix(path, base),
		FileName: strings.TrimSuffix(base, ext),
		Ext:      ext,
	}
}

// HasFile : Check if file exists in the current directory
func HasFile(filename string) bool {
	if info, err := os.Stat(filename); os.IsExist(err) {
		return !info.IsDir()
	}
	return false
}

// GetFiles :
// 	isRecursive:
// 		false to return []File in `folder`
//		true  to return []File in `folder` and all `subfolders`
func GetFiles(folder string, isRecursive bool) ([]File, error) {
	return GetFilesFunc(folder, isRecursive, func(f File) bool {
		return true
	})
}

// GetFilesString :
// 	isRecursive:
// 		false to return []File in `folder`
//		true  to return []File in `folder` and all `subfolders`
func GetFilesString(folder string, isRecursive bool) ([]string, error) {
	return GetFilesFuncString(folder, isRecursive, func(f File) bool {
		return true
	})
}

// GetFilesFunc :
// 	isRecursive:
// 		false to get []File in `folder`
// 		true  to get []File in `folder` and all `subfolders`
// 	filter(file) return true to exclude
func GetFilesFunc(folder string, isRecursive bool, filter func(file File) bool) ([]File, error) {
	var files []File
	if isRecursive {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			file, err := filepath.Abs(path)
			if err != nil {
				return nil
			}

			if !info.IsDir() {
				f := ConstructFile(file)
				if !filter(f) {
					files = append(files, f)
				}
			}
			return nil
		})

		return files, err
	}

	f, err := os.Open(folder)
	defer f.Close()

	if err != nil {
		return files, err
	}

	if fileinfo, err := f.Readdir(-1); err == nil {
		for _, file := range fileinfo {
			if !file.IsDir() {
				folder, err := filepath.Abs(folder)
				if err != nil {
					return files, err
				}
				f := ConstructFile(folder + "\\" + file.Name())
				if !filter(f) {
					files = append(files, f)
				}
			}
		}
	} else {
		return files, err
	}

	return files, nil
}

// GetFilesFuncString :
// 	isRecursive:
// 		false to get []File in `folder`
// 		true  to get []File in `folder` and all `subfolders`
// 	filter(file) return true to exclude
func GetFilesFuncString(folder string, isRecursive bool, filter func(file File) bool) ([]string, error) {
	var files []string
	flist, err := GetFilesFunc(folder, isRecursive, filter)
	if err != nil {
		return nil, err
	}
	for _, f := range flist {
		files = append(files, f.FullPath)
	}
	return files, nil
}

// GetNewFilePath change folder in path
//
// Example:
// 	path := "/aaa/bbb/ccc/example.xxx"
// 	path => File{
// 		FullPath: "/aaa/bbb/ccc/example.xxx",
// 		File:     "example.xxx",
// 		Folder:   "/aaa/bbb/ccc/",
// 		FileName: "example",
// 		Ext:      ".xxx",
// 	}
// 	sourceFolder := "/aaa/bbb/"
// 	targetFolder := "ddd/"
// 	return "ddd/ccc/example.xxx"
func GetNewFilePath(file File, sourceFolder, targetFolder string) (string, error) {
	if file.FullPath == "" {
		return "", fmt.Errorf("%s", "Original file is not valid.")
	}
	subfolder := GetSubfolder(file, sourceFolder)
	return targetFolder + subfolder + file.File, nil
}

// GetSubfolder remove `sourceFolder` of path and return the remainder of  subfolder
//
// Example:
// 	path := "/aaa/bbb/ccc/example.xxx"
// 	path => File{
// 		FullPath: "/aaa/bbb/ccc/example.xxx",
// 		File:     "example.xxx",
// 		Folder:   "/aaa/bbb/ccc/",
// 		FileName: "example",
// 		Ext:      ".xxx",
// 	}
// 	sourceFolder := "/aaa/bbb/"
// 	return "ccc/"
func GetSubfolder(file File, sourceFolder string) string {
	return TrimPrefix(file.Folder, sourceFolder)
}

// GrouppingFiles is groupping `files`, first sorted by fullpath then sorted by file name
func GrouppingFiles(files []File) {
	// assemble by folder
	gps := make(map[string][]File)
	gpnames := make(map[string]int)
	var fdnames []string
	for _, f := range files {
		if _, ok := gpnames[f.Folder]; !ok {
			fdnames = append(fdnames, f.Folder)
			gpnames[f.Folder] = 1
		}
		gps[f.Folder] = append(gps[f.Folder], f)
	}
	// sort folder
	sort.Strings(fdnames)
	// sort file in folder
	for _, g := range gps {
		sort.SliceStable(g, func(i, j int) bool {
			return g[i].File < g[j].File
		})
	}
	i := 0
	for _, folder := range fdnames {
		copy(files[i:], gps[folder])
		i += len(gps[folder])
	}
}
