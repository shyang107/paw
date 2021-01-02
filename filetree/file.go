package filetree

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/xattr"
	"github.com/shyang107/paw"
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
// 	`Size` is size of File
// 	`XAttributes` is extend attributes of File but ignore error
type File struct {
	Path        string
	Dir         string
	BaseName    string
	File        string
	Ext         string
	Stat        os.FileInfo
	Size        uint64
	XAttributes []string
}

// NewFile will the pointer of instance of `File`, and is a constructor of `File`.
func NewFile(path string) (*File, error) {
	// path = strings.TrimSuffix(path, "/")
	var err error
	if strings.HasPrefix(path, "~") {
		path, err = homedir.Expand(path)
	} else {
		path, err = filepath.Abs(path)
	}
	if err != nil {
		return nil, err
	}
	stat, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(path)
	basename := filepath.Base(path)
	// dir, basename := filepath.Split(path)
	ext := filepath.Ext(path)
	file := strings.TrimSuffix(basename, ext)
	var size = uint64(stat.Size())
	// if stat.IsDir() {
	// 	size, _ = sizes(path)
	// }
	var list []string
	if list, err = xattr.List(path); err != nil {
		return nil, err
	}
	return &File{
		Path:        path,
		Dir:         dir,
		BaseName:    basename,
		File:        file,
		Ext:         ext,
		Stat:        stat,
		Size:        size,
		XAttributes: list,
	}, nil
}

const (
	// RootMark = "."
	RootMark = "."
)

// NewFileRelTo will the pointer of instance of `File`, and is a constructor of `File`, but `File.Dir` is sub-directory of `root`
// 	If `path` == `root`, then
// 		f.Dir = "."
func NewFileRelTo(path, root string) (*File, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	f, err := NewFile(path)
	if err != nil {
		return nil, err
	}
	if f.IsDir() {
		if f.Path == root {
			f.Dir = paw.Replace(f.Path, root, ".", 1)
		} else {
			f.Dir = paw.Replace(f.Path, root, "..", 1)
		}
		return f, nil
	}

	if f.Dir == root {
		f.Dir = paw.Replace(f.Dir, root, ".", 1)
	} else {
		f.Dir = paw.Replace(f.Dir, root, "..", 1)
	}
	return f, nil
}

// func (f File) String() string {
// 	// return f.Path
// 	if NoColor {
// 		return f.BaseName
// 	}

// 	cvalue, _ := FileLSColorString(f.Path, f.BaseName)
// 	return cvalue
// }

// LSColorString will return a color string using LS_COLORS according to `f.Path` of file
func (f *File) LSColorString(s string) string {
	str, _ := FileLSColorString(f.Path, s)
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

// IsLink() report whether File describes a symbolic link.
func (f *File) IsLink() bool {
	mode := f.Stat.Mode()
	if mode&os.ModeSymlink != 0 {
		return true
	}
	return false
}

// IsFile reports whether File describes a file (not directory and symbolic link).
func (f *File) IsFile() bool {
	// if !f.IsDir() && !f.IsLink() {
	// 	return true
	// }
	// return false
	return nodeTypeFromFileInfo(f.Stat) == "file"
}

func nodeTypeFromFileInfo(fi os.FileInfo) string {
	switch fi.Mode() & (os.ModeType | os.ModeCharDevice) {
	case 0:
		return "file"
	case os.ModeDir:
		return "dir"
	case os.ModeSymlink:
		return "symlink"
	case os.ModeDevice | os.ModeCharDevice:
		return "chardev"
	case os.ModeDevice:
		return "dev"
	case os.ModeNamedPipe:
		return "fifo"
	case os.ModeSocket:
		return "socket"
	}

	return ""
}

// linux int64
// fileInfo, _ := os.Stat(path)
// stat_t := fileInfo.Sys().(*syscall.Stat_t)
// fmt.Println(stat_t.Atim.Sec)
// fmt.Println(stat_t.Ctim.Sec)
// fmt.Println(stat_t.Mtim.Sec)
//
// darwin int64
// fileInfo, _ := os.Stat(path)
// stat_t := fileInfo.Sys().(*syscall.Stat_t)
// fmt.Println(stat_t.Atimespec.Sec)
// fmt.Println(stat_t.Ctimespec.Sec)
// fmt.Println(stat_t.Mtimespec.Sec)
//
// windows int64
// fileInfo, _ := os.Stat(path)
// wFileSys := fileInfo.Sys().(*syscall.Win32FileAttributeData)
// tNanSeconds := wFileSys.CreationTime.Nanoseconds()  /// 返回的是纳秒
// tSec := tNanSeconds/1e9

// AccessTime reports the last access time of File.
func (f *File) AccessTime() time.Time {
	statT := f.Stat.Sys().(*syscall.Stat_t)
	return timespecToTime(statT.Atimespec)
}

// CreateTime reports the create time of file.
func (f *File) CreateTime() time.Time {
	statT := f.Stat.Sys().(*syscall.Stat_t)
	return timespecToTime(statT.Ctimespec)
}

// ModifiedTime reports the modify time of file.
func (f *File) ModifiedTime() time.Time {
	// statT := f.Stat.Sys().(*syscall.Stat_t)
	// return timespecToTime(statT.Mtimespec)
	return f.Stat.ModTime()
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

// // Size will return size of `File`
// func (f *File) Size() uint64 {
// 	return uint64(f.Stat.Size())
// }

// PathSlice will split `f.Path` following Spearator, seperating it into a string slice.
func (f *File) PathSlice() []string {
	return strings.Split(f.Path, PathSeparator)
}

// DirSlice will split `f.Dir` following Spearator, seperating it into a string slice.
func (f *File) DirSlice() []string {
	return strings.Split(f.Dir, PathSeparator)
}

// ColorBaseName will return a colorful string of BaseName using LS_COLORS like as exa
func (f *File) ColorBaseName() string {
	return getName(f)
}

// func getLTName(file *File) string {
func getName(file *File) string {
	name := file.LSColorString(file.BaseName)
	if file.IsDir() && file.Dir == RootMark {
		dir, _ := filepath.Split(file.Path)
		name = KindEXAColorString("dir", dir) + name
	}
	link := checkAndGetColorLink(file)
	if len(link) > 0 {
		// name += cpmap['l'].Sprint(" -> ") + link
		name += " -> " + link
	}
	return name
}

func checkAndGetColorLink(file *File) (link string) {
	mode := file.Stat.Mode()
	if mode&os.ModeSymlink != 0 {
		alink, err := filepath.EvalSymlinks(file.Path)
		if err != nil {
			link = alink + " ERR: " + err.Error()
		} else {
			// link, _ = FileLSColorString(alink, alink)
			link = GetColorizedDirName(alink, "")
		}
	}
	return link
}

func checkAndGetLink(file *File) (link string) {
	SetNoColor()
	link = checkAndGetColorLink(file)
	DefaultNoColor()
	return link
}

// ColorPermission will return a colorful string of Stat.Mode() like as exa.
// The length of placeholder in terminal is 10.
func (f *File) ColorPermission() string {
	permission := getColorizePermission(f.Stat.Mode())
	if len(f.XAttributes) > 0 {
		permission += NewEXAColor("-").Add(color.Concealed).Sprint("@")
	} else {
		permission += " "
	}
	return permission
}

// ColorModifyTime will return a colorful string of Stat.ModTime() like as exa.
// The length of placeholder in terminal is 14.
func (f *File) ColorModifyTime() string {
	return GetColorizedTime(f.Stat.ModTime())
}

// ColorGitStatus will return a colorful string of git status like as exa.
// The length of placeholder in terminal is 3.
func (f *File) ColorGitStatus(git GitStatus) string {
	return getColorizedGitStatus(git, f)
}

// ColorSize will return a colorful string of Size for human-reading like as exa.
// The length of placeholder in terminal is 6.
func (f *File) ColorSize() string {
	return GetColorizedSize(f.Size)
}

// ColorDirName will return a colorful string of {{dir of Path}}+{{name of path }} for human-reading like as exa
func (f *File) ColorDirName(root string) string {
	return GetColorizedDirName(f.Path, root)
}

// ColorMeta will return a colorful string of meta information of File (including Permission, Size, User, Group, Data Modified, Git and Name of File) and its' length.
func (f *File) ColorMeta(git GitStatus) (string, int) {
	meta, length := getMeta("", f, git)
	return meta, length
}

func getMeta(pad string, file *File, git GitStatus) (string, int) {
	width := 0
	buf := new(bytes.Buffer)
	cperm := file.ColorPermission()
	width += len(fmt.Sprintf("%v", file.Stat.Mode())) + 2
	cmodTime := file.ColorModifyTime()
	width += len(fmt.Sprint(file.ModifiedTime().Format("01-02-06 15:04"))) + 1
	cfsize := file.ColorSize()
	if file.IsDir() {
		cfsize = KindLSColorString("-", fmt.Sprintf("%6s", "-"))
	}
	width += 7
	width += len(urname) + len(gpname) + 1
	if git.NoGit {
		printLTList(buf, pad, cperm, cfsize, curname, cgpname, cmodTime)
	} else {
		cgit := file.ColorGitStatus(git)
		width += 4
		printLTList(buf, pad, cperm, cfsize, curname, cgpname, cmodTime, cgit)
	}
	return string(buf.Bytes()), width + 1
}
