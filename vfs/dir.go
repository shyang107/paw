package vfs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

// dir 代表一個目錄
type Dir struct {
	// full path = filepath.Join(root, reldir, name)
	path    string
	relpath string
	// basename
	name     string
	info     fs.FileInfo
	xattrs   []string
	git      *GitStatus
	relpaths []string

	// 存放該目錄下的子項，value 可能是 *dir 或 *file
	// map[basename]fs.DirEntry
	children map[string]fs.DirEntry

	errors []error

	// ReadDir 遍歷用
	idx int
}

// 實現 fs.FileInfo 接口
// A FileInfo describes a file and is returned by Stat.
// type FileInfo interface {
//     Name() string       // base name of the file
//     Size() int64        // length in bytes for regular files; system-dependent for others
//     Mode() FileMode     // file mode bits
//     ModTime() time.Time // modification time
//     IsDir() bool        // abbreviation for Mode().IsDir()
//     Sys() interface{}   // underlying data source (can return nil)
// }
// 文件也是某個目錄下的目錄項，因此需要實現 fs.DirEntry 接口
// A DirEntry is an entry read from a directory (using the ReadDir function or a ReadDirFile's ReadDir method).
// type DirEntry interface {
// 	Name() string // Name returns the name of the file (or subdirectory) described by the entry.
// 	// This name is only the final element of the path (the base name), not the entire path.
// 	// For example, Name would return "hello.go" not "/home/gopher/hello.go".

// 	IsDir() bool // IsDir reports whether the entry describes a directory.

// 	Type() FileMode // Type returns the type bits for the entry.
// 	// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.

// 	Info() (FileInfo, error) // Info returns the FileInfo for the file or subdirectory described by the entry.
// 	// The returned FileInfo may be from the time of the original directory read
// 	// or from the time of the call to Info. If the file has been removed or renamed
// 	// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
// 	// If the entry denotes a symbolic link, Info reports the information about the link itself,
// 	// not the link's target.
// }

// Both interfaces fs.FileInfo and  fs.DirEntry
//     Name() string       // base name of the file
//     Size() int64        // length in bytes for regular files; system-dependent for others
//     Mode() FileMode     // file mode bits
// 	Type() FileMode // Type returns the type bits for the entry.
//     ModTime() time.Time // modification time
//     IsDir() bool        // abbreviation for Mode().IsDir()
//     Sys() interface{}   // underlying data source (can return
// 	Info() (FileInfo, error) // Info returns the FileInfo for the file or subdirectory described by the entry.

// Name is base name of the file,  returns the name of the file (or subdirectory) described by the entry.
// This name is only the final element of the path (the base name), not the entire path.
// For example, Name would return "hello.go" not "/home/gopher/hello.go".
func (d *Dir) Name() string {
	return d.name
}

// Size returns length in bytes for regular files; system-dependent for others
func (d *Dir) Size() int64 {
	return d.info.Size()
}

// Mode returns file mode bits
func (d *Dir) Mode() fs.FileMode {
	return d.info.Mode()
}

// ModTime returns modification time
func (d *Dir) ModTime() time.Time {
	return d.info.ModTime()
}

// IsDir is abbreviation for Mode().IsDir()
// IsDir reports whether the entry describes a directory.
func (d *Dir) IsDir() bool {
	return true
}

// Sys returns underlying data source (can return nil)
func (d *Dir) Sys() interface{} {
	return d.info.Sys()
}

// Type returns the type bits for the entry.
// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
func (d *Dir) Type() fs.FileMode {
	return d.Mode()
}

// Info returns the FileInfo for the file or subdirectory described by the entry.
// The returned FileInfo may be from the time of the original directory read
// or from the time of the call to Info. If the file has been removed or renamed
// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
// If the entry denotes a symbolic link, Info reports the information about the link itself,
// not the link's target.
func (d *Dir) Info() (fs.FileInfo, error) {
	return d.info, nil
}

// ReadDir 實現 fs.ReadDirFile 接口，方便遍歷目錄
func (d *Dir) ReadDir(n int) ([]fs.DirEntry, error) {
	names := make([]string, 0, len(d.children))
	for name := range d.children {
		names = append(names, name)
	}

	sort.Slice(names, func(i, j int) bool {
		return strings.ToLower(names[i]) < strings.ToLower(names[j])
	})

	totalEntry := len(names)
	if n <= 0 {
		n = totalEntry
	}

	dirEntries := make([]fs.DirEntry, 0, n)
	for i := d.idx; i < n && i < totalEntry; i++ {
		name := names[i]
		child := d.children[name]

		f, isFile := child.(*File)
		if isFile {
			dirEntries = append(dirEntries, f)
		} else {
			dirEntry := child.(*Dir)
			dirEntries = append(dirEntries, dirEntry)
		}

		d.idx = i
	}

	return dirEntries, nil
}

// Xattibutes get the extended attributes of Dir
// 	implements the interface of Extended
func (d *Dir) Xattibutes() []string {
	return d.xattrs
}

// Path get the full-path of Dir
// 	implements the interface of DirEntryX
func (d *Dir) Path() string {
	return d.path
}

// Path get the relative path of Dir with respect to some basepath (indicated in creating new intance of Dir)
// 	implements the interface of DirEntryX
func (d *Dir) RelPath() string {
	return d.relpath
	// relpath, _ := filepath.Rel(basepath, d.path)
	// return relpath
}

// NameToLink return colorized name & symlink
func (d *Dir) NameToLink() string {
	if d.IsLink() {
		return d.name + " -> " + d.LinkPath()
	}
	return d.name
}

// LinkPath report far-end path of a symbolic link.
func (d *Dir) LinkPath() string {
	if d.IsLink() {
		// alink, err := filepath.EvalSymlinks(f.Path)
		alink, err := os.Readlink(d.path)
		if err != nil {
			return err.Error()
		}
		return alink
	}
	return ""
}

// LSColor will return LS_COLORS color of File
// 	implements the interface of DirEntryX
func (d *Dir) LSColor() *color.Color {
	return deLSColor(d)
}

// INode will return the inode number of File
func (d *Dir) INode() uint64 {
	inode := uint64(0)
	if sys := d.info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			inode = stat.Ino
		}
	}
	return inode
	// sys := d.Stat.Sys()
	// inode := reflect.ValueOf(sys).Elem().FieldByName("Ino").Uint()
	// return inode
}

// HDLinks will return the number of hard links of File
func (d *Dir) HDLinks() uint64 {
	nlink := uint64(0)
	if sys := d.info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			nlink = uint64(stat.Nlink)
		}
	}
	return nlink
}

// Blocks will return number of file system blocks of File
func (d *Dir) Blocks() uint64 {
	blocks := uint64(0)
	if sys := d.info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			blocks = uint64(stat.Blocks)
		}
	}
	return blocks
}

// Uid returns user id of File
func (d *Dir) Uid() uint32 {
	id := uint32(0)
	if sys := d.info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			id = (stat.Uid)
		}
	}
	return id
}

// User returns user (owner) name of File
func (d *Dir) User() string {
	u, err := user.LookupId(fmt.Sprint(d.Uid()))
	if err != nil {
		return err.Error()
	}
	return u.Username
}

// Gid returns group id of File
func (d *Dir) Gid() uint32 {
	id := uint32(0)
	if sys := d.info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			id = (stat.Gid)
		}
	}
	return id
}

// Group returns group (owner) name of File
func (d *Dir) Group() string {
	g, err := user.LookupGroupId(fmt.Sprint(d.Gid()))
	if err != nil {
		return err.Error()
	}
	return g.Name
}

// Dev will return dev id of File
func (d *Dir) Dev() uint64 {
	dev := uint64(0)
	if sys := d.info.Sys(); sys != nil {
		if stat, ok := sys.(*syscall.Stat_t); ok {
			dev = uint64(stat.Rdev)
		}
	}
	return dev
}

// DevNumber returns device number of a Darwin device number.
func (d *Dir) DevNumber() (uint32, uint32) {
	major, minor := paw.DevNumber(d.Dev())
	return major, minor
}

// DevNumberS returns device number of a Darwin device number.
func (d *Dir) DevNumberS() string {
	major, minor := paw.DevNumber(d.Dev())
	dev := fmt.Sprintf("%v,%v", major, minor)
	return dev
}

// AccessedTime reports the last access time of File.
func (d *Dir) AccessedTime() time.Time {
	statT := d.info.Sys().(*syscall.Stat_t)
	return timespecToTime(statT.Atimespec)
}

// CreatedTime reports the create time of file.
func (d *Dir) CreatedTime() time.Time {
	statT := d.info.Sys().(*syscall.Stat_t)
	return timespecToTime(statT.Birthtimespec)
}

// ModifiedTime reports the modify time of file.
func (d *Dir) ModifiedTime() time.Time {
	return d.ModTime()
}

// Md5 returns md5 codes of Dir
func (d *Dir) Md5() string {
	return "-"
}

// IsLink() report whether File describes a symbolic link.
func (d *Dir) IsLink() bool {
	return d.info.Mode()&os.ModeSymlink != 0
}

// IsFile reports whether File describes a regular file.
func (d *Dir) IsFile() bool {
	return false
}

// IsCharDev() report whether File describes a Unix character device, when ModeDevice is set.
func (d *Dir) IsCharDev() bool {
	return false
}

// IsDev() report whether File describes a device file.
func (d *Dir) IsDev() bool {
	return false
}

// IsFIFO() report whether File describes a named pipe.
func (d *Dir) IsFIFO() bool {
	return false
}

// IsSocket() report whether File describes a socket.
func (d *Dir) IsSocket() bool {
	return false
}

// IsTemporary() report whether File describes a temporary file; Plan 9 only.
func (d *Dir) IsTemporary() bool {
	return d.info.Mode()&os.ModeTemporary != 0
}

// IsExecOwner is to tell if the file is executable by its owner, use bitmask 0100:
func (d *Dir) IsExecOwner() bool {
	mode := d.info.Mode()
	return mode&0100 != 0
}

// IsExecGroup is to tell if the file is executable by the group, use bitmask 0010:
func (d *Dir) IsExecGroup() bool {
	mode := d.info.Mode()
	return mode&0010 != 0
}

// IsExecOther is to tell if the file is executable by others, use bitmask 0001:
func (d *Dir) IsExecOther() bool {
	mode := d.info.Mode()
	return mode&0001 != 0
}

// IsExecAny is to tell if the file is executable by any of its owner, the group and others, use bitmask 0111:
func (d *Dir) IsExecAny() bool {
	mode := d.info.Mode()
	return mode&0111 != 0
}

//IsExecAll is to tell if the file is executable by any of its owner, the group and others, again use bitmask 0111 but check if the result equals to 0111:
func (d *Dir) IsExecAll() bool {
	mode := d.info.Mode()
	return mode&0111 == 0111
}

// IsExecutable is to tell if the file isexecutable.
func (d *Dir) IsExecutable() bool {
	// return d.IsExecOwner() || d.IsExecGroup() || d.IsExecOther()
	return d.IsExecAny()
}

// =====================================

func (d *Dir) checkGitDir() {

	gs := d.git.GetStatus()
	if d.git.NoGit || len(d.children) < 1 || gs == nil || len(d.children) < 1 {
		return
	}
	// 1. check: if dir is GitIgnored, then marks all subfiles with GitIgnored.
	isMarkIgnored := false
	isUntracked := false
	var xy GitFileStatus
	rp := d.RelPath() + "/"
	if gxy, ok := gs[rp]; ok {
		if isXY(gxy, GitIgnored) {
			paw.Logger.WithField("rp", rp).Trace("GitIgnored")
			isMarkIgnored = true
			xy = *gxy
		}
		if isXY(gxy, GitUntracked) {
			paw.Logger.WithField("rp", rp).Trace("GitUntracked")
			isUntracked = true
			xy = *gxy
		}
	}

	if isMarkIgnored || isUntracked {
		markChildGit(d, &xy)
	}
}

func (d *Dir) checkGitFiles() {
	gs := d.git.GetStatus()
	if d.git.NoGit || len(d.children) < 1 || gs == nil || len(d.children) < 1 {
		return
	}
	// 2. if any of subfiles of dir (including root) has any change of git status, set GitChanged to dir
	for _, e := range d.children {
		next, isDir := e.(*Dir)
		if isDir {
			next.setSubDirXY()
		}
	}
	d.setSubDirXY()
}

func (d *Dir) setSubDirXY() {
	gs := d.git.GetStatus()
	xs, ys := d.getSubXYs()
	if len(xs) > 0 || len(ys) > 0 {
		rp := d.RelPath() + "/"
		paw.Logger.WithFields(logrus.Fields{
			"rp": "" + color.New(color.FgMagenta).Sprint(rp) + "",
			"xs": xs,
			"ys": ys,
		}).Trace()
		gs[rp] = &GitFileStatus{
			Staging:  getSC(xs),
			Worktree: getSC(ys),
			Extra:    d.Name() + "/",
		}
	}
}

func (d *Dir) getSubXYs() (xs, ys []GitStatusCode) {
	gs := d.git.GetStatus()

	xs = make([]GitStatusCode, 0)
	ys = make([]GitStatusCode, 0)
	for _, e := range d.children {
		f, isFiele := e.(*File)
		if isFiele {
			rp := f.RelPath()
			if xy, ok := gs[rp]; ok {
				if xy.Staging != GitUnmodified {
					xs = append(xs, xy.Staging)
				}
				if xy.Worktree != GitUnmodified {
					ys = append(ys, xy.Worktree)
				}
			}
		} else {
			d := e.(*Dir)
			sxs, sys := d.getSubXYs()
			xs = append(xs, sxs...)
			ys = append(ys, sys...)
		}
	}
	return xs, ys
}

func markChildGit(d *Dir, xy *GitFileStatus) {
	gs := d.git.GetStatus()
	ds, _ := d.ReadDir(-1)
	d.resetIdx()
	if len(ds) == 0 {
		return
	}
	for _, child := range ds {
		f, isFile := child.(*File)
		if isFile {
			rp := f.RelPath()
			gs[rp] = &GitFileStatus{
				Staging:  xy.Staging,
				Worktree: xy.Worktree,
				Extra:    f.Name(),
			}
		} else {
			dd := child.(*Dir)
			rp := dd.RelPath() + "/"
			gs[rp] = &GitFileStatus{
				Staging:  xy.Staging,
				Worktree: xy.Worktree,
				Extra:    dd.Name(),
			}
			// paw.Logger.WithFields(logrus.Fields{
			// 	"rp": rp,
			// 	"xy": gs[rp].Staging.String() + gs[rp].Worktree.String(),
			// }).Debug()
			dd.checkGitDir()
		}
	}
}

func isXY(xy *GitFileStatus, gcode GitStatusCode) bool {
	return xy.Staging == gcode ||
		xy.Worktree == gcode
}

func (d *Dir) Git() *GitStatus {
	return d.git
}

func (d *Dir) XY() string {
	return d.git.XY(d.RelPath() + "/")
}

func (d *Dir) FprintErrors(w io.Writer, pad string) {
	if d.errors != nil && len(d.errors) > 0 {
		for _, err := range d.errors {
			fmt.Fprintf(w, "%s%v\n", pad, cerror.Sprint(err))
		}
	}
}

func (d *Dir) resetIdx() {
	d.idx = 0
}

func (d *Dir) NItems() (ndirs, nfiles int) {
	for _, entry := range d.children {
		_, isFile := entry.(*File)
		if isFile {
			nfiles++
		} else {
			ndirs++
			dd := entry.(*Dir)
			nd, nf := dd.NItems()
			ndirs += nd
			nfiles += nf
		}
	}
	return ndirs, nfiles
}

func (d *Dir) RelPaths() []string {
	return d.relpaths
}

func (d *Dir) RelDir() string {
	return filepath.Dir(d.RelPath())
}

// // NameToLink return colorized name & symlink
// func (d *Dir) NameToLink() string {
// 	if d.IsLink() {
// 		return d.name + " -> " + d.LinkPath()
// 	}
// 	return d.name
// }

// // LinkPath report far-end path of a symbolic link.
// func (d *Dir) LinkPath() string {
// 	if d.IsLink() {
// 		// alink, err := filepath.EvalSymlinks(d.Path)
// 		alink, err := os.Readlink(d.path)
// 		if err != nil {
// 			return err.Error()
// 		}
// 		return alink
// 	}
// 	return ""
// }

// Field returns the specified value of File according to ViewField
func (d *Dir) Field(field ViewField, git *GitStatus) string {
	switch field {
	case ViewFieldNo:
		return fmt.Sprint(field.Value())
	case ViewFieldINode:
		return fmt.Sprint(d.INode())
	case ViewFieldPermissions:
		return permissionS(d)
	case ViewFieldLinks:
		return fmt.Sprint(d.HDLinks())
	case ViewFieldSize, ViewFieldBlocks:
		// return bytefmt.ByteSize(uint64(d.Size()))
		return "-"
	// case ViewFieldBlocks:
	// 	return "-"
	case ViewFieldUser:
		return d.User()
	case ViewFieldGroup:
		return d.Group()
	case ViewFieldModified:
		return dateS(d.ModifiedTime())
	case ViewFieldCreated:
		return dateS(d.CreatedTime())
	case ViewFieldAccessed:
		return dateS(d.AccessedTime())
	case ViewFieldGit:
		return d.XY()
	case ViewFieldMd5:
		return d.Md5()
	case ViewFieldName:
		return d.Name()
	default:
		return ""
	}
}

// FieldC returns the specified colorful value of File according to ViewField
func (d *Dir) FieldC(field ViewField, git *GitStatus) string {
	value := aligned(field, d.Field(field, git))
	switch field {
	case ViewFieldNo:
		return aligned(field, cdip.Sprint(field.Value()))
	case ViewFieldPermissions:
		return aligned(field, permissionC(d))
	case ViewFieldSize:
		return sizeCaligned(d)
	case ViewFieldBlocks:
		return blocksCaligned(d)
	case ViewFieldUser: //"User",
		furname := d.User()
		var c *color.Color
		if furname != urname {
			c = cunp
		} else {
			c = cuup
		}
		return aligned(field, c.Sprint(furname))
	case ViewFieldGroup: //"Group",
		fgpname := d.Group()
		var c *color.Color
		if fgpname != gpname {
			c = cgnp
		} else {
			c = cgup
		}
		return aligned(field, c.Sprint(fgpname))
	case ViewFieldGit:
		return aligned(field, d.git.XYc(d.RelPath()+"/"))
	case ViewFieldName:
		return cdip.Sprint(d.Name())
	default:
		return field.Color().Sprint(value)
	}
}

func (d *Dir) widthOfSize() (width, wmajor, wminor int) {
	return 1, 0, 0
}

// WidthOf returns width of string of field
func (d *Dir) WidthOf(field ViewField) int {
	var w int
	switch field {
	case ViewFieldSize, ViewFieldBlocks:
		w = 1
		// case PFieldGit:
		// 	w = 3
	case ViewFieldMd5:
		w = len(d.Md5())
	case ViewFieldName:
		w = 0
	default:
		w = paw.StringWidth(d.Field(field, nil))
	}
	return w
}

// getDir 通過一個路徑獲取其 dir 類型實例
func (d *Dir) getDir(relpath string) (*Dir, error) {
	if relpath == "." {
		return d, nil
	}
	if strings.HasPrefix(relpath, d.RelPath()+"/") {
		relpath = strings.TrimPrefix(relpath, d.RelPath()+"/")
	}
	parts := strings.Split(relpath, "/")

	cur := d
	for _, part := range parts {
		child := cur.children[part]
		if child == nil {
			return nil, fmt.Errorf("%s is not exists", relpath)
		}

		childDir, ok := child.(*Dir)
		if !ok {
			return nil, fmt.Errorf("%s is not directory", relpath)
		}

		cur = childDir
	}

	return cur, nil
}