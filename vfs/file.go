package vfs

import (
	"fmt"
	"io/fs"

	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/cast"
)

// file 代表一個文件
type File struct {
	path    string // full path = filepath.Join(root, relpath, name)
	relpath string
	name    string // basename
	info    FileInfo
	xattrs  []string
	git     *GitStatus
	//
	linkPath string
	isLink   bool
}

func NewFile(path, root string, git *GitStatus) (*File, error) {
	apath, err := filepath.Abs(path)
	if err != nil {
		// paw.Logger.Error(err)
		return nil, &fs.PathError{
			Op:   "NewFile",
			Path: path,
			Err:  err,
		}
	}
	info, err := os.Lstat(apath)
	if err != nil {
		// paw.Logger.Error(err)
		return nil, &fs.PathError{
			Op:   "NewFile",
			Path: path,
			Err:  err,
		}
	}

	var link string
	isLink := false
	if info.Mode()&os.ModeSymlink != 0 {
		info, _ = os.Stat(apath)
		isLink = true
		link = getPathFromLink(apath)
		if !filepath.IsAbs(link) { // get absolute path of link
			dir := filepath.Dir(apath)
			link = filepath.Join(dir, link)
		}
	}

	if info.IsDir() && !isLink {
		err := fmt.Errorf("%q is a directory.", path)
		// paw.Logger.Error(err)
		return nil, &fs.PathError{
			Op:   "NewFile",
			Path: path,
			Err:  err,
		}
	}
	// dir, _ := filepath.Split(apath)
	// git := NewGitStatus(dir)
	relpath := "."
	if len(root) > 0 {
		relpath, _ = filepath.Rel(root, apath)
	}
	name := filepath.Base(apath)
	xattrs, _ := GetXattr(apath)
	return &File{
		path:     apath,
		relpath:  relpath,
		name:     name,
		info:     info,
		xattrs:   xattrs,
		git:      git,
		isLink:   isLink,
		linkPath: link,
	}, nil
}

// 實現 fs.FileInfo 接口
// A FileInfo describes a file and is returned by Stat.
// type FileInfo interface:
//     Name() string       // base name of the file
//     Size() int64        // length in bytes for regular files; system-dependent for others
//     Mode() FileMode     // file mode bits
//     ModTime() time.Time // modification time
//     IsDir() bool        // abbreviation for Mode().IsDir()
//     Sys() interface{}   // underlying data source (can return nil)
//---------------------------------------------------------------------
// 文件也是某個目錄下的檔案目錄項，因此需要實現 fs.DirEntry 接口
// A DirEntry is an entry read from a directory (using the ReadDir function or a ReadDirFile's ReadDir method).
// type DirEntry interface :
// 	Name() string // base name of the file
// 	IsDir() bool // abbreviation for Mode().IsDir()
// 	Type() FileMode // file mode bits
// 	Info() (FileInfo, error) // Info returns the FileInfo for the file or subdirectory described by the entry.
//---------------------------------------------------------------------
// Both interfaces fs.FileInfo and  fs.DirEntry
//     Name() string       // fs.FileInfo & fs.DirEntry
//     Size() int64        // fs.FileInfo
//     Mode() FileMode     // fs.FileInfo
//     ModTime() time.Time // fs.FileInfo
//     IsDir() bool        // fs.FileInfo & fs.DirEntry
//     Sys() interface{}   // fs.FileInfo
// 	   Type() FileMode     // fs.DirEntry; = Mode()
//     Info() (FileInfo, error) // fs.DirEntry
//---------------------------------------------------------------------

//---------------------------------------------------------------------
// 實現 fs.FileInfo & fs.DirEntry 接口：

// Name is base name of the file,  returns the name of the file (or subdirectory) described by the entry.
// This name is only the final element of the path (the base name), not the entire path.
// For example, Name would return "hello.go" not "/home/gopher/hello.go".
func (f *File) Name() string {
	return f.name
}

// Size returns length in bytes for regular files; system-dependent for others
func (f *File) Size() int64 {
	// return int64(f.content.Len())
	return f.info.Size()
}

// Mode returns file mode bits
func (f *File) Mode() FileMode {
	return f.info.Mode()
}

// ModTime returns modification time
func (f *File) ModTime() time.Time {
	return f.info.ModTime()
}

// IsDir is abbreviation for Mode().IsDir()
// IsDir reports whether the entry describes a directory.
func (f *File) IsDir() bool {
	// return f.Mode().IsDir()
	return false
}

// Sys returns underlying data source (can return nil)
func (f *File) Sys() interface{} {
	return f.info.Sys()
}

// Type returns the type bits for the entry.
// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
func (f *File) Type() FileMode {
	return f.Mode()
}

// Info returns the FileInfo for the file or subdirectory described by the entry.
// The returned FileInfo may be from the time of the original directory read
// or from the time of the call to Info. If the file has been removed or renamed
// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
// If the entry denotes a symbolic link, Info reports the information about the link itself,
// not the link's target.
func (f *File) Info() (FileInfo, error) {
	return f.info, nil
}

//---------------------------------------------------------------------
// 實現 Extendeder 接口：

// Xattibutes get the extended attributes of File
// 	implements the interface of Extended
func (f *File) Xattibutes() []string {
	return f.xattrs
}

//---------------------------------------------------------------------
// 實現 Fielder 接口：

// Path get the full-path of File
// 	implements the interface of DirEntryX
func (f *File) Path() string {
	return f.path
}

// RelPath get the relative path of File with respect to some basepath (indicated in creating new intance of File)
// 	implements the interface of DirEntryX
func (f *File) RelPath() string {
	return f.relpath
	// relpath, _ := filepath.Rel(basepath, f.path)
	// return relpath
}

// RelDir get dir part of File.RelPath()
func (f *File) RelDir() string {
	return filepath.Dir(f.RelPath())
}

// LSColor will return LS_COLORS color of File
// 	implements the interface of DirEntryX
func (f *File) LSColor() *Color {
	return GetDexLSColor(f)
}

// NameToLink return colorized name & symlink
func (f *File) NameToLink() string {
	if f.IsLink() {
		return f.name + " -> " + f.LinkPath()
	}
	return f.name
}

// LinkPath report far-end path of a symbolic link.
func (f *File) LinkPath() string {
	return f.linkPath
	// if f.IsLink() {
	// 	// alink, err := filepath.EvalSymlinks(f.Path)
	// 	alink, err := os.Readlink(f.path)
	// 	if err != nil {
	// 		return err.Error()
	// 	}
	// 	return alink
	// }
	// return ""
}

// INode will return the inode number of File
func (f *File) INode() uint64 {
	if stat, ok := f.info.Sys().(*syscall.Stat_t); ok {
		return stat.Ino
	}
	return 0
	// sys := f.Stat.Sys()
	// inode := reflect.ValueOf(sys).Elem().FieldByName("Ino").Uint()
	// return inode
}

// HDLinks will return the number of hard links of File
func (f *File) HDLinks() uint64 {
	if stat, ok := f.info.Sys().(*syscall.Stat_t); ok {
		return uint64(stat.Nlink)
	}
	return 0
}

// Blocks will return number of file system blocks of File
func (f *File) Blocks() uint64 {
	if stat, ok := f.info.Sys().(*syscall.Stat_t); ok {
		return uint64(stat.Blocks)
	}
	return 0
}

// Uid returns user id of File
func (f *File) Uid() uint32 {
	if stat, ok := f.info.Sys().(*syscall.Stat_t); ok {
		return (stat.Uid)
	}
	return uint32(os.Getuid())
}

// User returns user (owner) name of File
func (f *File) User() string {
	u, err := user.LookupId(cast.ToString(f.Uid()))
	if err != nil {
		return err.Error()
	}
	return u.Username
}

// Gid returns group id of File
func (f *File) Gid() uint32 {
	if stat, ok := f.info.Sys().(*syscall.Stat_t); ok {
		return (stat.Gid)
	}
	return uint32(os.Getgid())
}

// Group returns group (owner) name of File
func (f *File) Group() string {
	g, err := user.LookupGroupId(cast.ToString(f.Gid()))
	if err != nil {
		return err.Error()
	}
	return g.Name
}

// Dev will return dev id of File
func (f *File) Dev() uint64 {
	if stat, ok := f.info.Sys().(*syscall.Stat_t); ok {
		return uint64(stat.Rdev)
	}
	return 0
}

// DevNumber returns device number of a Darwin device number.
func (f *File) DevNumber() (uint32, uint32) {
	major, minor := paw.DevNumber(f.Dev())
	return major, minor
}

// DevNumberS returns device number of a Darwin device number.
func (f *File) DevNumberS() string {
	major, minor := paw.DevNumber(f.Dev())
	dev := fmt.Sprintf("%v,%v", major, minor)
	return dev
}

// AccessedTime reports the last access time of File.
func (f *File) AccessedTime() time.Time {
	statT := f.info.Sys().(*syscall.Stat_t)
	return timespecToTime(statT.Atimespec)
}

// CreatedTime reports the create time of file.
func (f *File) CreatedTime() time.Time {
	statT := f.info.Sys().(*syscall.Stat_t)
	return timespecToTime(statT.Birthtimespec)
}

// ModifiedTime reports the modify time of file.
func (f *File) ModifiedTime() time.Time {
	return f.ModTime()
}

// Md5 returns md5 codes of File
func (f *File) Md5() string {
	if !f.info.Mode().IsRegular() {
		return "-"
	} else {
		return paw.GenMd5(f.Path())
	}
}

func (f *File) Git() *GitStatus {
	return f.git
}

func (f *File) XY() string {
	return f.git.XY(f.RelPath())
}

// Field returns the specified value of File according to ViewField
func (f *File) Field(field ViewField) string {
	switch field {
	case ViewFieldNo:
		return cast.ToString(field.Value())
	case ViewFieldINode:
		return cast.ToString(f.INode())
	case ViewFieldPermissions:
		return permissionS(f)
	case ViewFieldLinks:
		return cast.ToString(f.HDLinks())
	case ViewFieldSize:
		return f.SizeS()
	case ViewFieldBlocks:
		if f.Blocks() == 0 {
			return "-"
		}
		return cast.ToString(f.Blocks())
	case ViewFieldUser:
		return f.User()
	case ViewFieldGroup:
		return f.Group()
	case ViewFieldModified:
		return dateS(f.ModifiedTime())
	case ViewFieldCreated:
		return dateS(f.CreatedTime())
	case ViewFieldAccessed:
		return dateS(f.AccessedTime())
	case ViewFieldGit:
		return f.XY()
	case ViewFieldMd5:
		return f.Md5()
	case ViewFieldName:
		return f.NameToLink() //f.Name()
	default:
		return ""
	}
}

// FieldC returns the specified colorful value of File according to ViewField
func (f *File) FieldC(fd ViewField) string {
	switch fd {
	case ViewFieldNo:
		return alNoC(f)
		// return paw.Cfip.Sprint(fd.AlignedS(fd.Value()))
	case ViewFieldPermissions:
		return alPermissionC(f)
	case ViewFieldSize:
		return alSizeC(f)
	case ViewFieldBlocks:
		return alBlockC(f)
	case ViewFieldUser: //"User",
		return alUserC(f)
	case ViewFieldGroup: //"Group",
		return alGroupC(f)
	case ViewFieldGit:
		return alXYC(f)
	case ViewFieldName:
		return alNameC(f)
	default:
		return alFieldC(f, fd)
	}
}

func (f *File) widthOfSize() (width, wmajor, wminor int) {
	if f.IsCharDev() || f.IsDev() {
		major, minor := f.DevNumber()
		wmajor = len(cast.ToString(major))
		wminor = len(cast.ToString(minor))
		return wmajor + wminor + 1, wmajor, wminor
	}
	return len(f.Field(ViewFieldSize)), 0, 0
}

// WidthOf returns width of string of field
func (f *File) WidthOf(field ViewField) int {
	var w int
	switch field {
	case ViewFieldSize:
		w, _, _ = f.widthOfSize()
		// case PFieldGit:
		// 	w = 3
	case ViewFieldMd5:
		w = len(f.Md5())
	case ViewFieldName:
		w = 0
	default:
		w = paw.StringWidth(f.Field(field))
	}
	return w
}

//---------------------------------------------------------------------
// 實現 ISer 接口：

// IsLink() report whether File describes a symbolic link.
func (f *File) IsLink() bool {
	return f.isLink
	// return f.info.Mode()&os.ModeSymlink != 0
}

// IsFile reports whether File describes a file.
func (f *File) IsFile() bool {
	// return f.Mode().IsRegular()
	return true
}

// IsCharDev() report whether File describes a Unix character device, when ModeDevice is set.
func (f *File) IsCharDev() bool {
	return f.info.Mode()&os.ModeCharDevice != 0
}

// IsDev() report whether File describes a device file.
func (f *File) IsDev() bool {
	return f.info.Mode()&os.ModeDevice != 0
}

// IsFIFO() report whether File describes a named pipe.
func (f *File) IsFIFO() bool {
	return f.info.Mode()&os.ModeNamedPipe != 0
}

// IsSocket() report whether File describes a socket.
func (f *File) IsSocket() bool {
	return f.info.Mode()&os.ModeSocket != 0
}

// IsTemporary() report whether File describes a temporary file; Plan 9 only.
func (f *File) IsTemporary() bool {
	return f.info.Mode()&os.ModeTemporary != 0
}

// IsExecOwner is to tell if the file is executable by its owner, use bitmask 0100:
func (f *File) IsExecOwner() bool {
	mode := f.info.Mode()
	return mode&0100 != 0
}

// IsExecGroup is to tell if the file is executable by the group, use bitmask 0010:
func (f *File) IsExecGroup() bool {
	mode := f.info.Mode()
	return mode&0010 != 0
}

// IsExecOther is to tell if the file is executable by others, use bitmask 0001:
func (f *File) IsExecOther() bool {
	mode := f.info.Mode()
	return mode&0001 != 0
}

// IsExecAny is to tell if the file is executable by any of its owner, the group and others, use bitmask 0111:
func (f *File) IsExecAny() bool {
	mode := f.info.Mode()
	return mode&0111 != 0
}

//IsExecAll is to tell if the file is executable by any of its owner, the group and others, again use bitmask 0111 but check if the result equals to 0111:
func (f *File) IsExecAll() bool {
	mode := f.info.Mode()
	return mode&0111 == 0111
}

// IsExecutable is to tell if the file isexecutable.
func (f *File) IsExecutable() bool {
	// return f.IsExecOwner() || f.IsExecGroup() || f.IsExecOther()
	return f.IsExecAny()
}

// ====================================================================

// IsRegularFile reports whether File describes a regular file.
func (f *File) IsRegularFile() bool {
	return f.Mode().IsRegular()
}

func (f *File) SizeS() string {
	return _sizeSC(f, false)
}

func (f *File) SizeC() string {
	return _sizeSC(f, true)
}
