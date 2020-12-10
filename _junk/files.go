package paw

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// File : path information
//
// 	Fields:
// 	  `FullPath`: The full path including the folder
// 	  `ShortPath` : The short path is `FullPath` without rootfolder part (replace with './')
// 	  `Folder`: The folder of the file
// 	  `ShortFolder`: The folder of the file without rootfolder part
// 	  `File`: The file name including extension (basename)
// 	  `FileName`: The file name excluding extension
// 	  `Ext`: Extension of the file
type File struct {
	FullPath    string // The full path including the folder
	ShortPath   string // The short path is `FullPath` without rootfolder part
	Folder      string // The folder of the file
	ShortFolder string // The folder of the file without rootfolder part
	File        string // The file name including extension (basename)
	FileName    string // The file name excluding extension
	Ext         string // Extension of the file
	Info        os.FileInfo
}

// IsDir will return whether the file is directory
func (f *File) IsDir() bool {
	return f.Info.IsDir()
}

// SetFileInfo will store  FileInfo
func (f *File) SetFileInfo(fi os.FileInfo) {
	f.Info = fi
}

// ConstructFile construct `paw.File` from string
//
// Example:
// 	path := "/aaa/bbb/ccc/example.xxx"
//  root := "/aaa/"
// 	path => File{
// 		FullPath:    "/aaa/bbb/ccc/example.xxx",
// 		ShortPath:   "bbb/ccc/example.xxx"
// 		File:        "example.xxx",
// 		Folder:      "/aaa/bbb/ccc/",
// 		ShortFolder: "bbb/ccc/",
// 		FileName:    "example",
// 		Ext:         ".xxx",
// 	}
func ConstructFile(path string, root string) File {
	fi, err := os.Stat(path)
	if err != nil {
		fi = nil
	}
	base := filepath.Base(path)
	ext := filepath.Ext(path)
	folder := TrimSuffix(path, base)
	// shortFolder, _ := filepath.Rel(root, folder)
	shortFolder := TrimPrefix(path, root)
	shortPath := shortFolder + base
	return File{
		FullPath:    path,
		ShortPath:   shortPath,
		File:        base,
		Folder:      folder,
		ShortFolder: shortFolder,
		FileName:    TrimSuffix(base, ext),
		Ext:         ext,
		Info:        fi,
	}
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

var (
	// ExcludePattern is pattern used in GetFilesFunc to exclude some files/folders such as `.git|.gitignore|.DS_Store|$RECYCLE.BIN|desktop.ini|_gsdata_`
	ExcludePattern = `\.git|\.gitignore|\.DS_Store|\$RECYCLE\.BIN|desktop\.ini|_gsdata_`
	// REUsuallyExclude is regexp used in GetFilesFunc to exclude some files/folders such as `.git|.gitignore|.DS_Store|$RECYCLE.BIN|desktop.ini|_gsdata_`
	REUsuallyExclude = regexp.MustCompile(ExcludePattern)
)

// GetFilesFunc :
// 	isRecursive:
// 		false to get []File in `folder`
// 		true  to get []File in `folder` and all `subfolders`
// 	exclude(file) return true to exclude
func GetFilesFunc(folder string, isRecursive bool, exclude func(file File) bool) ([]File, error) {
	var files []File
	// defer checkRootFiles(files, folder)

	if isRecursive {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			file, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			if !info.IsDir() && info.Mode().IsRegular() {
				f := ConstructFile(file, folder)
				if !exclude(f) {
					files = append(files, f)
				}
			}
			return nil
		})
		files = checkRootFiles(files, folder)
		return files, err
	}

	f, err := os.Open(folder)
	defer f.Close()

	if err != nil {
		// return files, err
		goto END
	}

	if fileinfo, err := f.Readdir(-1); err == nil {
		for _, file := range fileinfo {
			if !file.IsDir() && file.Mode().IsRegular() {
				folder, err := filepath.Abs(folder)
				if err != nil {
					// return files, err
					goto END
				}
				f := ConstructFile(folder+"\\"+file.Name(), folder)
				if !exclude(f) {
					files = append(files, f)
				}
			}
		}
	} else {
		// return files, err
		goto END
	}
END:
	files = checkRootFiles(files, folder)
	return files, err
}

const zeroRootFiles = "«zeroRootFiles»"

func checkRootFiles(files []File, root string) []File {
	nrf := 0
	// fmt.Println("root:", root)
	for _, f := range files {
		// fmt.Println("f.Folder:", f.Folder)
		if strings.EqualFold(f.Folder, root) {
			nrf++
		}
	}
	if nrf == 0 {
		files = append(files, ConstructFile(root+zeroRootFiles, root))
		rootHasFile = false
	}
	return files
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

// CountSubfolders count subfolders of `files`
func CountSubfolders(files []File) int {
	folders := make(map[string]int)
	for _, f := range files {
		if _, ok := folders[f.ShortFolder]; !ok {
			if !strings.EqualFold(f.ShortFolder, "./") {
				folders[f.ShortFolder] = 1
			}
		}
	}
	return len(folders)
}

// GrouppingFiles is groupping `files`, first sorted by fullpath then sorted by file name
func GrouppingFiles(files []File) {
	fl := &FileList{files}
	fl.OrderedByFolder()
}

// OutputMode : FileList output mode
type OutputMode uint

const (
	// OPlainTextMode : FileList output in plain text mode (default, use PrintPlain())
	OPlainTextMode OutputMode = iota
	// OTableFormatMode : FileList output in TableFormat mode (use PrintWithTableFormat())
	OTableFormatMode
	// OTreeMode : FileList output in tree mode (use PrintTree())
	OTreeMode
)
