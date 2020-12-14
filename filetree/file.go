package filetree

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// filetree is tree structure of files
//
//  every thing is file, even directory just a special file!
//

const (
	// PathSeparator is OS-specific path separator ('/')
	PathSeparator = string(os.PathSeparator)
	// PathListSeparator is OS-specific path list separator (':')
	PathListSeparator = string(os.PathListSeparator)
)

// File will store information of a file
//
// Fields:
// 	`Path` is an absolute representation of path. If the path is not absolute it will be joined with the current working directory to turn it into an absolute path. The absolute path name for a given file is not guaranteed to be unique.
// 	`Dir` is all but the last element of `Path`, typically the directory of path. After dropping the final element, and clean on the path and trailing slashes are removed. If the path is empty, Dir returns ".". If the path consists entirely of separators, Dir returns a single separator. The returned path does not end in a separator unless it is the root directory.
// 	`BaseName` is the last element of path. Trailing path separators are removed before extracting the last element. If the path is empty, Base returns ".". If the path consists entirely of separators, Base returns a single separator.
// 	`File` is the part of triming the suffix `Ext` of `File`
// 	`Ext` is the file name extension used by `Path`. The extension is the suffix beginning at the final dot in the final element of path; it is empty if there is no dot.
// 	`Stat` is `os.Stat(Path)` but ignoring error.
type File struct {
	Path     string
	Dir      string
	BaseName string
	File     string
	Ext      string
	Stat     os.FileInfo
}

// ConstructFile will the pointer of instance of `File`, and is a constructor of `File`.
func ConstructFile(path string) *File {
	// path = strings.TrimSuffix(path, "/")
	var err error
	if strings.HasPrefix(path, "~") {
		path, err = homedir.Expand(path)
	} else {
		path, err = filepath.Abs(path)
	}
	if err != nil {
		return nil
	}
	stat, err := os.Lstat(path)
	if err != nil {
		return nil
	}
	dir := filepath.Dir(path)
	basename := filepath.Base(path)
	// dir, basename := filepath.Split(path)
	ext := filepath.Ext(path)
	file := strings.TrimSuffix(basename, ext)
	return &File{
		Path:     path,
		Dir:      dir,
		BaseName: basename,
		File:     file,
		Ext:      ext,
		Stat:     stat,
	}
}

const (
	// RootMark = "."
	RootMark = "."
)

// ConstructFileRelTo will the pointer of instance of `File`, and is a constructor of `File`, but `File.Dir` is sub-directory of `root`
// 	If `path` == `root`, then
// 		f.Dir = "."
func ConstructFileRelTo(path, root string) *File {
	path, _ = filepath.Abs(path)
	root, _ = filepath.Abs(root)
	f := ConstructFile(path)
	if f.IsDir() {
		f.Dir = strings.Replace(f.Path, root, ".", 1)
	} else {
		f.Dir = strings.Replace(f.Dir, root, ".", 1)
	}
	// if f.Dir == "." {
	// 	f.Dir = "./"
	// 	// f.BaseName = "."
	// 	// f.File = ""
	// }
	return f
}

func (f File) String() string {
	// return f.Path
	if NoColor {
		return f.BaseName
	}

	cvalue, _ := FileLSColorStr(f.Path, f.BaseName)
	return cvalue
}

// LSColorString will return a color string using LS_COLORS according to `f.Path` of file
func (f *File) LSColorString(s string) string {
	str, _ := FileLSColorStr(f.Path, s)
	return str
}

// IsDir reports whether `f` describes a directory. That is, it tests for the ModeDir bit being set in `f`.
func (f *File) IsDir() bool {
	return f.Stat.IsDir()
}

// IsRegular reports whether `f` describes a regular file. That is, it tests that no mode type bits are set.
func (f *File) IsRegular() bool {
	return f.Stat.Mode().IsRegular()
}

// PathSlice will split `f.Path` following Spearator, seperating it into a string slice.
func (f *File) PathSlice() []string {
	return strings.Split(f.Path, PathSeparator)
}

// DirSlice will split `f.Dir` following Spearator, seperating it into a string slice.
func (f *File) DirSlice() []string {
	return strings.Split(f.Dir, PathSeparator)
}

// // SubPath will return sub-directory based on `root`.
// //  1. If `root` is equal to `f.Dir`, then return "./".
// // 	2. If `root` is not the prefix of `f.Dir`, then return `f.Dir`.
// func (f *File) SubPath(root string) string {
// 	sub := strings.TrimPrefix(f.Dir, root)
// 	if len(sub) == 0 {
// 		return "./"
// 	}
// 	return sub
// }

// // SubPathSlice will split `f.SubPath(root)` following Spearator, seperating it into a string slice.
// // 	The the last element will be empty string, and remove.
// func (f *File) SubPathSlice(root string) []string {
// 	ss := strings.Split(f.SubPath(root), PathSeparator)
// 	ns := len(ss)
// 	return ss[:ns-1]
// }
