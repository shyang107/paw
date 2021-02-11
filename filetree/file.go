package filetree

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
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
	Info        os.FileInfo
	Size        uint64
	XAttributes []string
	// User        string
	// Group       string
	cp    *color.Color
	UpDir *File
}

// NewFile will the pointer of instance of `File`, and is a constructor of `File`.
func NewFile(path string) (*File, error) {
	var (
		info                     os.FileInfo
		err                      error
		dir, basename, ext, file string
		size                     uint64
		xattrs                   []string
	)
	info, err = os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get stat: %s", err)
	}
	dir = filepath.Dir(path)
	basename = filepath.Base(path)
	ext = filepath.Ext(path)
	file = strings.TrimSuffix(basename, ext)
	size = uint64(info.Size())
	xattrs, err = getXattr(path)
	// if err != nil && pdOpt.isTrace {
	// 	paw.Logger.Warn(err)
	// }

	f := &File{
		Path:        path,
		Dir:         dir,
		BaseName:    basename,
		File:        file,
		Ext:         ext,
		Info:        info,
		Size:        size,
		XAttributes: xattrs,
		UpDir:       nil,
	}
	f.cp = GetFileLSColor(f)
	return f, nil
}

func getXattr(path string) ([]string, error) {
	// paw.Logger.WithField("path", path).Info("income")
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
	// paw.Logger.WithField("path", path).Info("input")
	f, err := NewFile(path)
	if err != nil {
		return nil, err
	}
	if len(root) == 0 {
		return f, nil
	}
	if path == root {
		f.Dir = RootMark
	} else {
		f.Dir = PathRel(f.Dir, root)
	}
	return f, nil
}

func (f File) String() string {
	return f.Name()
}

// LSColor will return LS_COLORS color of File
func (f *File) LSColor() *color.Color {
	return f.cp
}

// LSColorstring will return a color string using LS_COLORS according to `f.Path` of file
func (f *File) LSColorstring(s string) string {
	return f.LSColor().Sprint(s)
}

// GetUpDir return the directory file which the File is belong to.
func (f *File) GetUpDir() *File {
	return f.UpDir
}

// SetUpDir return set the directory file to which the File is belong to.
func (f *File) SetUpDir(up *File) *File {
	f.UpDir = up
	return f
}

// Name return File.BaseNameToLink()
func (f File) Name() string {
	return f.BaseNameToLink()
}

func (f File) NameC() string {
	return f.BaseNameToLinkC()
}

// BaseNameC will return a colorful string of BaseName using LS_COLORS like as exa
func (f *File) BaseNameC() string {
	return f.LSColor().Sprint(f.BaseName)
}

// BaseNameToLink return colorized name & symlink
func (f *File) BaseNameToLink() string {
	if f.IsLink() {
		return f.BaseName + " -> " + f.LinkPath()
	}
	return f.BaseName
}

// BaseNameToLinkC return colorized name & symlink
func (f *File) BaseNameToLinkC() string {
	if f.IsLink() {
		return f.BaseNameC() + cdashp.Sprint(" -> ") + f.LinkPathC()
	}
	return f.BaseNameC()
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

// LinkPathC return colorized far-end path string of a symbolic link.
func (f *File) LinkPathC() string {
	return GetColorizedPath(f.LinkPath(), "")
}

// PathSlice will split `f.Path` following Spearator, seperating it into a string slice.
func (f *File) PathSlice() []string {
	return strings.Split(f.Path, PathSeparator)
}

// DirSlice will split `f.Dir` following Spearator, seperating it into a string slice.
func (f *File) DirSlice() []string {
	return strings.Split(f.Dir, PathSeparator)
}

// DirNameC will return a colorful string of {{dir of Path}}+{{name of path }} for human-reading like as exa
func (f *File) DirNameC() string {
	// return GetColorizedPath(f.Path, "")
	dir, _ := filepath.Split(f.Path)
	dir, name := filepath.Split(dir)
	return cdirp.Sprint(dir) + cdip.Sprint(name)
}

// DirNameShortC will return a colorful string of {{dir of Path}}+{{name of path }} (replace root with '.') for human-reading like as exa
func (f *File) DirNameShortC(root string) string {
	// if f.Path == root {
	// 	return cdip.Sprint(".")
	// }
	// return GetColorizedPath(f.Path, root)
	dir, name := filepath.Split(f.Dir)
	dir = strings.Replace(dir, root, ".", 1)
	return cdirp.Sprint(dir) + cdip.Sprint(name)
}

// DirNameWrapC will return a colorful wrapped string according to width adn seprating with '\n'. If width <= 0, use sttyWidth
func (f *File) DirNameWrapC(pad string, width int) string {
	return rowWrapDirName(f.Dir, pad, width)
}

// INode will return the inode number of File
func (f *File) INode() uint64 {
	inode := uint64(0)
	if sys := f.Info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			inode = stat.Ino
		}
	}
	return inode
	// sys := f.Stat.Sys()
	// inode := reflect.ValueOf(sys).Elem().FieldByName("Ino").Uint()
	// return inode
}

// INodeC will return the colorful string of inode number of File
func (f *File) INodeC() string {
	return cinp.Sprint(f.INode())
}

// Permission will return a string of Info.Mode() like as exa.
// The length of placeholder in terminal is 11.
func (f *File) Permission() string {
	sperm := f.Info.Mode().String() //fmt.Sprint(f.Stat.Mode())

	// if strings.HasPrefix(sperm, "Dc") {
	// 	sperm = strings.Replace(sperm, "Dc", "c", 1)
	// }
	// if strings.HasPrefix(sperm, "D") {
	// 	sperm = strings.Replace(sperm, "D", "b", 1)
	// }
	// if strings.HasPrefix(sperm, "L") {
	// 	sperm = strings.Replace(sperm, "L", "l", 1)
	// }

	if f.XAttributes == nil {
		sperm += "?"
	} else {
		if len(f.XAttributes) > 0 {
			sperm += "@"
		} else {
			sperm += " "
		}
	}
	return sperm
}

// PermissionC will return a colorful string of Stat.Mode() like as exa.
// The length of placeholder in terminal is 11.
func (f *File) PermissionC() string {
	// sperm := f.Permission()
	// cxmark := " "
	// if strings.HasSuffix(sperm, "@") || strings.HasSuffix(sperm, "?") || strings.HasSuffix(sperm, " ") {
	// 	cxmark = cdashp.Sprint(string(sperm[len(sperm)-1]))
	// 	sperm = sperm[:len(sperm)-1]
	// }
	// permission := GetColorizedPermission(sperm) + cxmark
	// return permission
	return GetColorizedPermission(f.Permission())
}

// NLinks will return the number of hard links of File
func (f *File) NLinks() uint64 {
	nlink := uint64(0)
	if sys := f.Info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			nlink = uint64(stat.Nlink)
		}
	}
	return nlink
}

// NLinksC will return the colorful string of number of hard links of File
func (f *File) NLinksC() string {
	return clkp.Sprint(f.NLinks())
}

// // Size will return size of `File`
// func (f *File) Size() uint64 {
// 	return f.Size
// }

// ByteSize will retun total size of File in byte-format as human read
func (f *File) ByteSize() string {
	return ByteSize(f.Size)
}

// SizeC will return a colorful string of Size for human-reading like as exa.
// The length of placeholder in terminal is 6.
func (f *File) SizeC() string {
	return GetColorizedSize(f.Size)
}

// Blocks will return number of file system blocks of File
func (f *File) Blocks() uint64 {
	blocks := uint64(0)
	if sys := f.Info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			blocks = uint64(stat.Blocks)
		}
	}
	return blocks
}

// BlocksC will return a colorful string of numbe of file system blocks of File
func (f *File) BlocksC() string {
	if f.IsDir() {
		return cdap.Sprint("-")
	}
	return cbkp.Sprint(f.Blocks())
}

// Uid returns user id of File
func (f *File) Uid() uint32 {
	id := uint32(0)
	if sys := f.Info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			id = (stat.Uid)
		}
	}
	return id
}

// User returns user (owner) name of File
func (f *File) User() string {
	u, err := user.LookupId(fmt.Sprint(f.Uid()))
	if err != nil {
		return err.Error()
	}
	return u.Username
}

// UserC returns colorful user (owner) name of File
func (f *File) UserC() string {
	ur := f.User()
	var c *color.Color
	if ur != urname {
		c = cunp
	} else {
		c = cuup
	}
	return c.Sprint(ur)
}

// Gid returns group id of File
func (f *File) Gid() uint32 {
	id := uint32(0)
	if sys := f.Info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			id = (stat.Gid)
		}
	}
	return id
}

// Group returns group (owner) name of File
func (f *File) Group() string {
	g, err := user.LookupGroupId(fmt.Sprint(f.Gid()))
	if err != nil {
		return err.Error()
	}
	return g.Name
}

// GroupC returns colorful group (owner) name of File
func (f *File) GroupC() string {
	gp := f.Group()
	var c *color.Color
	if gp != gpname {
		c = cgnp
	} else {
		c = cgup
	}
	return c.Sprint(gp)
}

// Dev will return dev id of File
func (f *File) Dev() uint64 {
	dev := uint64(0)
	if sys := f.Info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			dev = uint64(stat.Rdev)
		}
	}
	return dev
	// dev := reflect.ValueOf(f.Stat.Sys()).Elem().FieldByName("dev").Uint()
	// return dev
}

// DevNumber returns device number of a Darwin device number.
func (f *File) DevNumber() (uint32, uint32) {
	major, minor := paw.DevNumber(f.Dev())
	return major, minor
}

// DevNumberString returns device number of a Darwin device number.
func (f *File) DevNumberString() string {
	major, minor := paw.DevNumber(f.Dev())
	dev := fmt.Sprintf("%v,%v", major, minor)
	return dev
}

// DevNumberStringC returns device number of a Darwin device number.
func (f *File) DevNumberStringC() string {
	major, minor := paw.DevNumber(f.Dev())
	dev := csnp.Sprint(major) + cdap.Sprint(",") + csnp.Sprint(minor)
	return dev
}

// AccessedTime reports the last access time of File.
func (f *File) AccessedTime() time.Time {
	statT := f.Info.Sys().(*syscall.Stat_t)
	return timespecToTime(statT.Atimespec)
}

// CreatedTime reports the create time of file.
func (f *File) CreatedTime() time.Time {
	statT := f.Info.Sys().(*syscall.Stat_t)
	return timespecToTime(statT.Birthtimespec)
}

// ModifiedTime reports the modify time of file.
func (f *File) ModifiedTime() time.Time {
	// statT := f.Stat.Sys().(*syscall.Stat_t)
	// return timespecToTime(statT.Mtimespec)
	return f.Info.ModTime()
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

// ModifiedTimeC will return a colorful string of Stat.ModTime() like as exa.
// The length of placeholder in terminal is 14.
func (f *File) ModifiedTimeC() string {
	return GetColorizedTime(f.ModifiedTime()) //+ sp
}

// AccessedTimeC will return a colorful string of File.AccessTime() like as exa.
// The length of placeholder in terminal is 14.
func (f *File) AccessedTimeC() string {
	return GetColorizedTime(f.AccessedTime()) //+ sp
}

// CreatedTimeC will return a colorful string of File.CreateTime() like as exa.
// The length of placeholder in terminal is 14.
func (f *File) CreatedTimeC() string {
	return GetColorizedTime(f.CreatedTime()) //+ sp
}

// GitStatus will return a string of git status like as exa.
// The length of placeholder in terminal is 3.
func (f *File) GitStatus(git GitStatus) string {
	return getGitStatus(git, f)
}

// GitStatusC will return a colorful string of git status like as exa.
// The length of placeholder in terminal is 3.
func (f *File) GitStatusC(git GitStatus) string {
	return getColorizedGitStatus(git, f)
}

// IsDir reports whether `f` describes a directory. That is, it tests for the ModeDir bit being set in `f`.
func (f *File) IsDir() bool {
	return f.Info.IsDir()
}

// IsLink() report whether File describes a symbolic link.
func (f *File) IsLink() bool {
	// return nodeTypeFromFileInfo(f.Info) == kindSymlink
	return f.Info.Mode()&os.ModeSymlink != 0
}

// IsFile reports whether File describes a regular file.
func (f *File) IsFile() bool {
	// if !f.IsDir() && !f.IsLink() {
	// 	return true
	// }
	// return false
	// return nodeTypeFromFileInfo(f.Info) == kindFile //"file"
	return f.Info.Mode().IsRegular()
}

// IsCharDev() report whether File describes a Unix character device, when ModeDevice is set.
func (f *File) IsCharDev() bool {
	// return nodeTypeFromFileInfo(f.Info) == kindChardev
	return f.Info.Mode()&os.ModeCharDevice != 0
}

// IsDev() report whether File describes a device file.
func (f *File) IsDev() bool {
	// return nodeTypeFromFileInfo(f.Info) == kindDev
	return f.Info.Mode()&os.ModeDevice != 0
}

// IsFIFO() report whether File describes a named pipe.
func (f *File) IsFIFO() bool {
	// return nodeTypeFromFileInfo(f.Info) == kindFIFO
	return f.Info.Mode()&os.ModeNamedPipe != 0
}

// IsSocket() report whether File describes a socket.
func (f *File) IsSocket() bool {
	// return nodeTypeFromFileInfo(f.Info) == kindSocket
	return f.Info.Mode()&os.ModeSocket != 0
}

// IsTemporary() report whether File describes a temporary file; Plan 9 only.
func (f *File) IsTemporary() bool {
	// return nodeTypeFromFileInfo(f.Info) == kindSocket
	return f.Info.Mode()&os.ModeTemporary != 0
}

// IsExecOwner is to tell if the file is executable by its owner, use bitmask 0100:
func (f *File) IsExecOwner() bool {
	mode := f.Info.Mode()
	return mode&0100 != 0
}

// IsExecGroup is to tell if the file is executable by the group, use bitmask 0010:
func (f *File) IsExecGroup() bool {
	mode := f.Info.Mode()
	return mode&0010 != 0
}

// IsExecOther is to tell if the file is executable by others, use bitmask 0001:
func (f *File) IsExecOther() bool {
	mode := f.Info.Mode()
	return mode&0001 != 0
}

// IsExecAny is to tell if the file is executable by any of its owner, the group and others, use bitmask 0111:
func (f *File) IsExecAny() bool {
	mode := f.Info.Mode()
	return mode&0111 != 0
}

//IsExecAll is to tell if the file is executable by any of its owner, the group and others, again use bitmask 0111 but check if the result equals to 0111:
func (f *File) IsExecAll() bool {
	mode := f.Info.Mode()
	return mode&0111 == 0111
}

// IsExecutable is to tell if the file isexecutable.
func (f *File) IsExecutable() bool {
	// return f.IsExecOwner() || f.IsExecGroup() || f.IsExecOther()
	return f.IsExecAny()
}

// IsNotIdentify() report whether File describes a not-identify.
func (f *File) IsNotIdentify() bool {
	return nodeTypeFromFileInfo(f.Info) == kindNotIdentify
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
	switch nodeTypeFromFileInfo(f.Info) {
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

// Meta will return a string of meta information of File (including Permission, Size, User, Group, Modified, Git and Name of File) and its' length.
func (f *File) Meta(git GitStatus) (string, int) {

	if len(pfieldKeys) == 0 {
		pfieldKeys = pfieldKeysDefualt
	}

	fds := NewFieldSliceFrom(pfieldKeys, git)
	fds.SetValues(f, git)
	return fds.MetaValuesString(), fds.MetaHeadsStringWidth()
}

// MetaC will return a colorful string of meta information of File (including Permission, Size, User, Group, Data Modified, Git and Name of File) and its' length.
func (f *File) MetaC(git GitStatus) (string, int) {

	if len(pfieldKeys) == 0 {
		pfieldKeys = pfieldKeysDefualt
	}

	fds := NewFieldSliceFrom(pfieldKeys, git)
	fds.SetValues(f, git)
	return fds.MetaValuesStringC(), fds.MetaHeadsStringWidth()
}

func (f *File) subDir() string {
	if f.IsDir() {
		return f.Dir + "/" + f.BaseName
	}
	return f.Dir
}

func (f *File) widthOfSize() (width, wmajor, wminor int) {
	if f.IsCharDev() || f.IsDev() {
		major, minor := f.DevNumber()
		wmajor = len(fmt.Sprint(major))
		wminor = len(fmt.Sprint(minor))
		// width = wmajor + wminor + 1
		return wmajor + wminor + 1, wmajor, wminor
	} else if f.IsDir() {
		return 1, 0, 0
	} else {
		return len(f.ByteSize()), 0, 0
	}
}

// WidthOf returns width of string of field
func (f *File) WidthOf(field PDFieldFlag) int {
	var w int
	switch field {
	case PFieldINode:
		w = len(fmt.Sprint(f.INode()))
	case PFieldPermissions:
		w = len(f.Permission())
	case PFieldLinks:
		w = len(fmt.Sprint(f.NLinks()))
	case PFieldSize:
		w, _, _ = f.widthOfSize()
	case PFieldBlocks:
		w = len(fmt.Sprint(f.Blocks()))
	case PFieldUser:
		w = paw.StringWidth(f.User())
	case PFieldGroup:
		w = paw.StringWidth(f.Group())
	case PFieldModified:
		w = len(DateString(f.ModifiedTime()))
	case PFieldCreated:
		w = len(DateString(f.CreatedTime()))
	case PFieldAccessed:
		w = len(DateString(f.AccessedTime()))
	// case PFieldGit:
	// 	w = 3
	default: // name
		w = 0
	}
	return w
}
