package vfs

import (
	"fmt"

	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/fatih/color"
	"github.com/shyang107/paw"
)

// file 代表一個文件
type File struct {
	path    string // full path = filepath.Join(root, relpath, name)
	relpath string
	name    string // basename
	info    fs.FileInfo
	xattrs  []string
	git     *GitStatus
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
func (f *File) Mode() fs.FileMode {
	return f.info.Mode()
}

// ModTime returns modification time
func (f *File) ModTime() time.Time {
	return f.info.ModTime()
}

// IsDir is abbreviation for Mode().IsDir()
// IsDir reports whether the entry describes a directory.
func (f *File) IsDir() bool {
	return f.Mode().IsDir()
	// return false
}

// Sys returns underlying data source (can return nil)
func (f *File) Sys() interface{} {
	return f.info.Sys()
}

// Type returns the type bits for the entry.
// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
func (f *File) Type() fs.FileMode {
	return f.Mode()
}

// Info returns the FileInfo for the file or subdirectory described by the entry.
// The returned FileInfo may be from the time of the original directory read
// or from the time of the call to Info. If the file has been removed or renamed
// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
// If the entry denotes a symbolic link, Info reports the information about the link itself,
// not the link's target.
func (f *File) Info() (fs.FileInfo, error) {
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
func (f *File) LSColor() *color.Color {
	return deLSColor(f)
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
	if f.IsLink() {
		// alink, err := filepath.EvalSymlinks(f.Path)
		alink, err := os.Readlink(f.path)
		if err != nil {
			return err.Error()
		}
		return alink
	}
	return ""
}

// INode will return the inode number of File
func (f *File) INode() uint64 {
	inode := uint64(0)
	if sys := f.info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			inode = stat.Ino
		}
	}
	return inode
	// sys := f.Stat.Sys()
	// inode := reflect.ValueOf(sys).Elem().FieldByName("Ino").Uint()
	// return inode
}

// HDLinks will return the number of hard links of File
func (f *File) HDLinks() uint64 {
	nlink := uint64(0)
	if sys := f.info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			nlink = uint64(stat.Nlink)
		}
	}
	return nlink
}

// Blocks will return number of file system blocks of File
func (f *File) Blocks() uint64 {
	blocks := uint64(0)
	if sys := f.info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			blocks = uint64(stat.Blocks)
		}
	}
	return blocks
}

// Uid returns user id of File
func (f *File) Uid() uint32 {
	id := uint32(0)
	if sys := f.info.Sys(); sys != nil {
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

// Gid returns group id of File
func (f *File) Gid() uint32 {
	id := uint32(0)
	if sys := f.info.Sys(); sys != nil {
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

// Dev will return dev id of File
func (f *File) Dev() uint64 {
	dev := uint64(0)
	if sys := f.info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			dev = uint64(stat.Rdev)
		}
	}
	return dev
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
		return fmt.Sprint(field.Value())
	case ViewFieldINode:
		return fmt.Sprint(f.INode())
	case ViewFieldPermissions:
		return permissionS(f)
	case ViewFieldLinks:
		return fmt.Sprint(f.HDLinks())
	case ViewFieldSize:
		if f.IsCharDev() || f.IsDev() {
			return f.DevNumberS()
		} else {
			return bytefmt.ByteSize(uint64(f.Size()))
		}
	case ViewFieldBlocks:
		return fmt.Sprint(f.Blocks())
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
func (f *File) FieldC(field ViewField) string {
	value := aligned(field, f.Field(field))
	switch field {
	case ViewFieldNo:
		return aligned(field, cfip.Sprint(field.Value()))
	case ViewFieldPermissions:
		return aligned(field, permissionC(f))
	case ViewFieldSize:
		if f.IsCharDev() || f.IsDev() {
			major, minor := f.DevNumber()
			wdmajor := _ViewFieldMajor.Width()
			wdminor := _ViewFieldMinor.Width()
			csj := csnp.Sprintf("%[1]*[2]v", wdmajor, major)
			csn := csnp.Sprintf("%[1]*[2]v", wdminor, minor)
			cdev := csj + cdirp.Sprint(",") + csn
			wdev := wdmajor + wdminor + 1 //len(paw.StripANSI(cdev))
			if wdev < field.Width() {
				cdev = csj + cdirp.Sprint(",") + paw.Spaces(field.Width()-wdev) + csn
			}
			return cdev
		} else {
			return sizeCaligned(f)
		}
	case ViewFieldUser: //"User",
		furname := f.User()
		var c *color.Color
		if furname != urname {
			c = cunp
		} else {
			c = cuup
		}
		return aligned(field, c.Sprint(furname))
	case ViewFieldGroup: //"Group",
		fgpname := f.Group()
		var c *color.Color
		if fgpname != gpname {
			c = cgnp
		} else {
			c = cgup
		}
		return aligned(field, c.Sprint(fgpname))
	case ViewFieldGit:
		return aligned(field, f.git.XYc(f.RelPath()))
	case ViewFieldName:
		return nameToLinkC(f)
	default:
		return field.Color().Sprint(value)
	}
}

func (f *File) widthOfSize() (width, wmajor, wminor int) {
	if f.IsCharDev() || f.IsDev() {
		major, minor := f.DevNumber()
		wmajor = len(fmt.Sprint(major))
		wminor = len(fmt.Sprint(minor))
		// width = wmajor + wminor + 1
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
	return f.info.Mode()&os.ModeSymlink != 0
}

// IsFile reports whether File describes a regular file.
func (f *File) IsFile() bool {
	return f.Mode().IsRegular()
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
