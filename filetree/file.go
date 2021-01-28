package filetree

import (
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
	XattrSymbol       = paw.XAttrSymbol
)

var (
	xattrsp = paw.Spaces(paw.StringWidth(XattrSymbol))
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
	ext := filepath.Ext(path)
	file := strings.TrimSuffix(basename, ext)
	size := uint64(stat.Size())
	// if stat.IsDir() {
	// 	size, _ = sizes(path)
	// }

	xattrs, err := getXattr(path)
	if err != nil {
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
		XAttributes: xattrs,
	}, nil
}

func getXattr(path string) ([]string, error) {
	xattrs, err := xattr.List(path)
	if err != nil {
		return xattrs, err
	}
	if len(xattrs) > 0 {
		for i, x := range xattrs {
			x, _ := xattr.Get(path, x)
			xattrs[i] = fmt.Sprintf("%s (len %d)", xattrs[i], len(x))
		}
	}
	return xattrs, nil
}

const (
	// RootMark = "."
	RootMark  = "."
	UpDirMark = ".."
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
		f.Dir = strings.Replace(f.Path, root, RootMark, 1)
		return f, nil
	}

	f.Dir = strings.Replace(f.Dir, root, RootMark, 1)
	return f, nil
}

func (f File) String() string {
	return f.Name()
}

func (f File) Name() string {
	return f.BaseNameToLink()
}

func (f File) ColorName() string {
	return f.ColorBaseNameToLink()
}

// ColorBaseName will return a colorful string of BaseName using LS_COLORS like as exa
func (f *File) ColorBaseName() string {
	return f.LSColor().Sprint(f.BaseName)
}

// LinkPath report far-end path of a symbolic link.
func (f *File) LinkPath() string {
	if f.IsLink() {
		alink, err := filepath.EvalSymlinks(f.Path)
		if err != nil {
			return fmt.Errorf("%s Err: %s", alink, err.Error()).Error()
		}
		return alink
	}
	return ""
}

// ColorLinkPath return colorized far-end path string of a symbolic link.
func (f *File) ColorLinkPath() string {
	return GetColorizedDirName(f.LinkPath(), "")
}

// BaseNameToLink return colorized name & symlink
func (f *File) BaseNameToLink() string {
	if f.IsLink() {
		return f.BaseName + " -> " + f.LinkPath()
	}
	return f.BaseName
}

// ColorBaseNameToLink return colorized name & symlink
func (f *File) ColorBaseNameToLink() string {
	if f.IsLink() {
		return f.ColorBaseName() + cdashp.Sprint(" -> ") + f.ColorLinkPath()
	}
	return f.ColorBaseName()
}

// PathSlice will split `f.Path` following Spearator, seperating it into a string slice.
func (f *File) PathSlice() []string {
	return strings.Split(f.Path, PathSeparator)
}

// DirSlice will split `f.Dir` following Spearator, seperating it into a string slice.
func (f *File) DirSlice() []string {
	return strings.Split(f.Dir, PathSeparator)
}

// ColorDirName will return a colorful string of {{dir of Path}}+{{name of path }} for human-reading like as exa
func (f *File) ColorDirName() string {
	return GetColorizedDirName(f.Path, "")
}

// ColorShortDirName will return a colorful string of {{dir of Path}}+{{name of path }} (replace root with '.') for human-reading like as exa
func (f *File) ColorShortDirName(root string) string {
	if f.Path == root {
		return cdip.Sprint(".")
	}
	return GetColorizedDirName(f.Path, root)
}

// ColorWrapDirName will return a colorful wrapped string according to width adn seprating with '\n'. If width <= 0, use sttyWidth
func (f *File) ColorWrapDirName(pad string, width int) string {
	return rowWrapDirName(f.Dir, pad, width)
}

// ByteSize will retun total size of File in byte-format as human read
func (f *File) ByteSize() string {
	return ByteSize(f.Size)
}

// Dev will return dev id of File
func (f *File) Dev() uint64 {
	dev := uint64(0)
	if sys := f.Stat.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			dev = uint64(stat.Dev)
		}
	}
	return dev
	// dev := reflect.ValueOf(f.Stat.Sys()).Elem().FieldByName("dev").Uint()
	// return dev
}

// NLinks will return the number of hard links of File
func (f *File) NLinks() uint64 {
	nlink := uint64(0)
	if sys := f.Stat.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			nlink = uint64(stat.Nlink)
		}
	}
	return nlink
	// nlink := reflect.ValueOf(f.Stat.Sys()).Elem().FieldByName("Nlink").Uint()
	// return nlink
}

// INode will return the inode number of File
func (f *File) INode() uint64 {
	inode := uint64(0)
	if sys := f.Stat.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			inode = stat.Ino
		}
	}
	return inode
	// sys := f.Stat.Sys()
	// inode := reflect.ValueOf(sys).Elem().FieldByName("Ino").Uint()
	// return inode
}

// Blocks will return number of file system blocks of File
func (f *File) Blocks() uint64 {
	blocks := uint64(0)
	if sys := f.Stat.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			blocks = uint64(stat.Blocks)
		}
	}
	return blocks
	// block := reflect.ValueOf(f.Stat.Sys()).Elem().FieldByName("Blocks").Int()
	// return block
}

// // Size will return size of `File`
// func (f *File) Size() uint64 {
// 	return uint64(f.Stat.Size())
// }

// LSColor will return LS_COLORS color of File
func (f *File) LSColor() *color.Color {
	return GetFileLSColor(f)
}

// Permission will return a string of Stat.Mode() like as exa.
// The length of placeholder in terminal is 11.
func (f *File) Permission() string {
	permission := fmt.Sprint(f.Stat.Mode())
	if len(f.XAttributes) > 0 {
		permission += "@"
	} else {
		permission += " "
	}
	return permission
}

// ColorPermission will return a colorful string of Stat.Mode() like as exa.
// The length of placeholder in terminal is 11.
func (f *File) ColorPermission() string {
	permission := GetColorizePermission(f.Stat.Mode())
	if len(f.XAttributes) > 0 {
		permission += cdashp.Sprint("@")
	} else {
		permission += " "
	}
	return permission
}

// GitStatus will return a string of git status like as exa.
// The length of placeholder in terminal is 3.
func (f *File) GitStatus(git GitStatus) string {
	return getGitStatus(git, f)
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

// LSColorstring will return a color string using LS_COLORS according to `f.Path` of file
func (f *File) LSColorstring(s string) string {
	// str, _ := FileLSColorString(f.Path, s)
	// return str
	return f.LSColor().Sprint(s)
}

// IsDir reports whether `f` describes a directory. That is, it tests for the ModeDir bit being set in `f`.
func (f *File) IsDir() bool {
	return f.Stat.IsDir()
}

// // IsRegularFile reports whether `f` describes a regular file. That is, it tests that no mode type bits are set.
// func (f *File) IsRegularFile() bool {
// 	return f.Stat.Mode().IsRegular()
// }

// IsLink() report whether File describes a symbolic link.
func (f *File) IsLink() bool {
	return nodeTypeFromFileInfo(f.Stat) == kindSymlink
}

// IsFile reports whether File describes a file (not directory and symbolic link).
func (f *File) IsFile() bool {
	// if !f.IsDir() && !f.IsLink() {
	// 	return true
	// }
	// return false
	return nodeTypeFromFileInfo(f.Stat) == kindFile //"file"
}

// IsChardev() report whether File describes a chardev.
func (f *File) IsChardev() bool {
	return nodeTypeFromFileInfo(f.Stat) == kindChardev
}

// IsDev() report whether File describes a dev.
func (f *File) IsDev() bool {
	return nodeTypeFromFileInfo(f.Stat) == kindDev
}

// IsFiFo() report whether File describes a named pipe.
func (f *File) IsFiFo() bool {
	return nodeTypeFromFileInfo(f.Stat) == kindFIFO
}

// IsSocket() report whether File describes a socket.
func (f *File) IsSocket() bool {
	return nodeTypeFromFileInfo(f.Stat) == kindSocket
}

// IsNotIdentify() report whether File describes a not-identify.
func (f *File) IsNotIdentify() bool {
	return nodeTypeFromFileInfo(f.Stat) == kindNotIdentify
}

type kindType int

const (
	kindFile kindType = iota
	kindDir
	kindSymlink
	kindChardev
	kindDev
	kindFIFO
	kindSocket
	kindNotIdentify
)

func nodeTypeFromFileInfo(fi os.FileInfo) kindType {
	switch fi.Mode() & (os.ModeType | os.ModeCharDevice) {
	case 0:
		return kindFile //"file"
	case os.ModeDir:
		return kindDir //"dir"
	case os.ModeSymlink:
		return kindSymlink // "symlink"
	case os.ModeDevice | os.ModeCharDevice:
		return kindChardev //"chardev"
	case os.ModeDevice:
		return kindDev //"dev"
	case os.ModeNamedPipe:
		return kindFIFO //"fifo"
	case os.ModeSocket:
		return kindSocket //"socket"
	}

	return kindNotIdentify
}

func (f *File) TypeString() string {
	switch nodeTypeFromFileInfo(f.Stat) {
	case kindFile:
		return "file"
	case kindDir:
		return "dir"
	case kindSymlink:
		return "symlink"
	case kindChardev:
		return "chardev"
	case kindDev:
		return "dev"
	case kindFIFO:
		return "fifo"
	case kindSocket:
		return "socket"
	default: //kindNotIdentify
		return "not identify"
	}
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

// AccessedTime reports the last access time of File.
func (f *File) AccessedTime() time.Time {
	statT := f.Stat.Sys().(*syscall.Stat_t)
	return timespecToTime(statT.Atimespec)
}

// CreatedTime reports the create time of file.
func (f *File) CreatedTime() time.Time {
	statT := f.Stat.Sys().(*syscall.Stat_t)
	return timespecToTime(statT.Birthtimespec)
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

// ColorModifyTime will return a colorful string of Stat.ModTime() like as exa.
// The length of placeholder in terminal is 14.
func (f *File) ColorModifyTime() string {
	return GetColorizedTime(f.ModifiedTime()) //+ sp
}

// ColorAccessTime will return a colorful string of File.AccessTime() like as exa.
// The length of placeholder in terminal is 14.
func (f *File) ColorAccessedTime() string {
	return GetColorizedTime(f.AccessedTime()) //+ sp
}

// ColorCreatedTime will return a colorful string of File.CreateTime() like as exa.
// The length of placeholder in terminal is 14.
func (f *File) ColorCreatedTime() string {
	return GetColorizedTime(f.CreatedTime()) //+ sp
}

// ColorMeta will return a colorful string of meta information of File (including Permission, Size, User, Group, Data Modified, Git and Name of File) and its' length.
func (f *File) ColorMeta(git GitStatus) (string, int) {

	if len(pfieldKeys) == 0 {
		pfieldKeys = pfieldKeysDefualt
	}

	fds := NewFieldSliceFrom(pfieldKeys, git)
	fds.SetValues(f, git)
	return fds.ColorMetaValuesString(), fds.MetaHeadsStringWidth()
}

// func getMeta(file *File, git GitStatus) (string, int) {
// 	sb := paw.NewStringBuilder()
// 	csb := paw.NewStringBuilder()
// 	for _, k := range pfieldKeys {
// 		field, cfield := "", ""
// 		switch k {
// 		case PFieldINode: //"inode",
// 			field = fmt.Sprintf("%[1]*[2]v", pfieldWidthsMap[k], file.INode())
// 			cfield = cinp.Sprint(field)
// 		case PFieldPermissions: //"Permissions",
// 			field = fmt.Sprintf("%[1]*[2]v", pfieldWidthsMap[k], file.INode())
// 			cfield = file.ColorPermission()
// 		case PFieldLinks: //"Links",
// 			field = fmt.Sprintf("%[1]*[2]v", pfieldWidthsMap[k], file.NLinks())
// 			cfield = clkp.Sprint(field)
// 		case PFieldSize: //"Size",
// 			if file.IsDir() {
// 				field = fmt.Sprintf("%[1]*[2]v", pfieldWidthsMap[k], "-")
// 				cfield = cdashp.Sprint(field)
// 			} else {
// 				field = fmt.Sprintf("%[1]*[2]v", pfieldWidthsMap[k], ByteSize(file.Size))
// 				cfield = file.ColorSize()
// 			}
// 		case PFieldBlocks: //"User",
// 			if file.IsDir() {
// 				field = fmt.Sprintf("%[1]*[2]v", pfieldWidthsMap[k], "-")
// 				cfield = cdashp.Sprint(field)
// 			} else {
// 				field = fmt.Sprintf("%[1]*[2]v", pfieldWidthsMap[k], file.Blocks())
// 				cfield = cbkp.Sprint(field)
// 			}
// 		case PFieldUser: //"User",
// 			field = fmt.Sprintf("%[1]*[2]v", pfieldWidthsMap[k], urname)
// 			cfield = cuup.Sprint(field)
// 		case PFieldGroup: //"Group",
// 			field = fmt.Sprintf("%[1]*[2]v", pfieldWidthsMap[k], gpname)
// 			cfield = cgup.Sprint(field)
// 		case PFieldModified: //"Date Modified",
// 			field = fmt.Sprintf("%-[1]*[2]v", pfieldWidthsMap[k], DateString(file.ModifiedTime()))
// 			cfield = cdap.Sprint(field)
// 		case PFieldCreated: //"Date Created",
// 			field = fmt.Sprintf("%-[1]*[2]v", pfieldWidthsMap[k], DateString(file.CreatedTime()))
// 			cfield = cdap.Sprint(field)
// 		case PFieldAccessed: //"Date Accessed",
// 			field = fmt.Sprintf("%-[1]*[2]v", pfieldWidthsMap[k], DateString(file.AccessedTime()))
// 			cfield = cdap.Sprint(field)
// 		case PFieldGit: //"Gid",
// 			if git.NoGit {
// 				continue
// 			} else {
// 				field = fmt.Sprintf("%[1]*[2]v", pfieldWidthsMap[k], file.Stat.Mode())
// 				cfield = file.ColorGitStatus(git)
// 			}
// 			// case PFieldName: //"Name",
// 		}
// 		// field := fmt.Sprintf("%[1]*[2]s", pfieldWidthsMap[k], fieldsMap[k])
// 		fmt.Fprintf(sb, "%s ", field)
// 		fmt.Fprintf(csb, "%s ", cfield)
// 	}
// 	head := sb.String()
// 	head = head[:len(head)-1]
// 	width := paw.StringWidth(head)
// 	chead := csb.String()
// 	chead = chead[:len(chead)-1]
// 	return chead, width

// }
