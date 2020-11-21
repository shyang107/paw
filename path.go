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

// // GetGlobFilesList 獲得目標檔列表
// func GetGlobFilesList(folder string, pattern string) ([]string, error) {
// 	if len(folder) == 0 {
// 		folder = "" + "."
// 	}
// 	pattern = "" + folder + "/" + pattern
// 	fls, err := filepath.Glob(pattern)
// 	return fls, err
// }

// // GetFolderFileInfo gets and returns the `FileInfo` list from the specific `folder`
// func GetFolderFileInfo(folder string) []os.FileInfo {
// 	var files []os.FileInfo
// 	f, err := os.Open(folder)
// 	if err != nil {
// 		Logger.WithFields(logrus.Fields{
// 			"folder": folder,
// 		}).Fatal(err)
// 	}
// 	defer f.Close()
// 	if fileInfos, err := f.Readdir(-1); err == nil {
// 		for _, fi := range fileInfos {
// 			if !fi.IsDir() {
// 				files = append(files, fi)
// 			}
// 		}
// 	}
// 	return files
// }

// // GetFolderFileString gets and returns the file string list from the specific `folder`
// func GetFolderFileString(folder string) []string {
// 	var files []string
// 	fileInfos := GetFolderFileInfo(folder)
// 	for _, fi := range fileInfos {
// 		files = append(files, fi.Name())
// 	}
// 	return files
// }

// // GetAllSubfolderFileInfo gets and returns all `FileInfo` list in `root` folder and its all subfolders
// func GetAllSubfolderFileInfo(root string) []os.FileInfo {
// 	var files []os.FileInfo
// 	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		if !info.IsDir() {
// 			files = append(files, info)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		Logger.WithFields(logrus.Fields{
// 			"folder": root,
// 		}).Fatal(err)
// 	}
// 	return files
// }

// // GetAllSubfolderString gets and returns all files string list in `root` folder and its all subfolders
// func GetAllSubfolderString(root string) []string {
// 	var files []string
// 	for _, fi := range GetAllSubfolderFileInfo(root) {
// 		files = append(files, fi.Name())
// 	}
// 	return files
// }

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
// 	isInclude(file) return true to include
func GetFilesFunc(folder string, isRecursive bool, isInclude func(file File) bool) ([]File, error) {
	var files []File

	if isRecursive {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			file, err := filepath.Abs(path)
			if err != nil {
				return nil
			}

			if !info.IsDir() {
				f := ConstructFile(file)
				if isInclude(f) {
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
				if isInclude(f) {
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
// 	isInclude(file) return true to include
func GetFilesFuncString(folder string, isRecursive bool, isInclude func(file File) bool) ([]string, error) {
	var files []string
	flist, err := GetFilesFunc(folder, isRecursive, isInclude)
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
	return strings.TrimPrefix(file.Folder, sourceFolder)
}

// byFolder is used in sort with key `Folder`
type byFolder []File

func (f byFolder) Len() int           { return len(f) }
func (f byFolder) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f byFolder) Less(i, j int) bool { return f[i].Folder < f[j].Folder }

type sortByFile []File

func (a sortByFile) Len() int           { return len(a) }
func (a sortByFile) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByFile) Less(i, j int) bool { return a[i].File < a[j].File }

type sortByString []string

func (a sortByString) Len() int           { return len(a) }
func (a sortByString) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByString) Less(i, j int) bool { return a[i] < a[j] }

// GrouppingFiles is groupping `files`, first sorted by fullpath then sorted by file name
func GrouppingFiles(files []File) []File {
	tfiles := files
	sort.Sort(byFolder(tfiles))
	fd := make(map[string][]File)
	fdm := make(map[string]int)
	for _, f := range files {
		fdm[f.Folder] = 1
		fd[f.Folder] = append(fd[f.Folder], f)
	}
	for _, d := range fd {
		sort.Sort(sortByFile(d))
	}
	var fds []string
	for k := range fdm {
		fds = append(fds, k)
	}
	sort.Sort(sortByString(fds))
	var sfiles []File
	for _, s := range fds {
		sfiles = append(sfiles, fd[s]...)
	}
	return sfiles
}
