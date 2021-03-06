package vfs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/cast"
	"github.com/shyang107/paw/cnested"
	"github.com/sirupsen/logrus"
)

// dir 代表一個目錄
type Dir struct {
	File
	// path     string // full path = filepath.Join(root, relpath, name)
	// relpath  string
	// name     string // basename
	// info     FileInfo
	// xattrs   []string
	// git      *GitStatus
	// linkPath string
	// isLink   bool

	// 存放該目錄下的子項，value 可能是 *dir 或 *file
	// map[basename]DirEntryX
	children map[string]DirEntryX
	relpaths []string
	errors   []error
	// ReadDir 遍歷用
	idx int
	opt *VFSOption
	//
}

func NewDir(dirpath, root string, git *GitStatus, opt *VFSOption) (*Dir, error) {
	f, err := _NewFile(dirpath, root, git)
	if err != nil {
		return nil, &fs.PathError{
			Op:   "NewDir",
			Path: dirpath,
			Err:  err,
		}
	}
	if !f.info.IsDir() {
		err := fmt.Errorf("%q is not a directory.", root)
		// paw.Logger.Error(err)
		return nil, &fs.PathError{
			Op:   "NewDir",
			Path: root,
			Err:  err,
		}
	}
	return &Dir{
		File: *f,
		// relpaths: []string{},
		relpaths: []string{f.relpath},
		// errors:   []error{},
		children: make(map[string]DirEntryX),
		opt:      opt,
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
// 需要實現 fs.DirEntry 接口
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
// fs.FileInfo & fs.DirEntry 接口：

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
func (d *Dir) Mode() FileMode {
	return d.info.Mode()
}

// ModTime returns modification time
func (d *Dir) ModTime() time.Time {
	return d.info.ModTime()
}

// IsDir is abbreviation for Mode().IsDir()
// IsDir reports whether the entry describes a directory.
func (d *Dir) IsDir() bool {
	// return d.Mode().IsDir()
	return true
}

// Sys returns underlying data source (can return nil)
func (d *Dir) Sys() interface{} {
	return d.info.Sys()
}

// Type returns the type bits for the entry.
// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
func (d *Dir) Type() FileMode {
	return d.Mode()
}

// Info returns the FileInfo for the file or subdirectory described by the entry.
// The returned FileInfo may be from the time of the original directory read
// or from the time of the call to Info. If the file has been removed or renamed
// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
// If the entry denotes a symbolic link, Info reports the information about the link itself,
// not the link's target.
func (d *Dir) Info() (FileInfo, error) {
	return d.info, nil
}

//---------------------------------------------------------------------
// 實現 Extendeder 接口：

// Xattibutes get the extended attributes of Dir
// 	implements the interface of Extended
func (d *Dir) Xattibutes() []string {
	return d.xattrs
}

//---------------------------------------------------------------------
// 實現 Fielder 接口：

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

// RelDir get dir part of File.RelPath()
func (d *Dir) RelDir() string {
	return filepath.Dir(d.RelPath())
}

// LSColor will return LS_COLORS color of File
// 	implements the interface of DirEntryX
func (d *Dir) LSColor() *Color {
	return GetDexLSColor(d)
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
	return d.linkPath
	// if d.IsLink() {
	// 	// alink, err := filepath.EvalSymlinks(f.Path)
	// 	alink, err := os.Readlink(d.path)
	// 	if err != nil {
	// 		return err.Error()
	// 	}
	// 	return alink
	// }
	// return ""
}

// INode will return the inode number of File
func (d *Dir) INode() uint64 {
	if stat, ok := d.info.Sys().(*syscall.Stat_t); ok {
		return stat.Ino
	}
	return 0
}

// HDLinks will return the number of hard links of File
func (d *Dir) HDLinks() uint64 {
	if stat, ok := d.info.Sys().(*syscall.Stat_t); ok {
		return uint64(stat.Nlink)
	}
	return 0
}

// Blocks will return number of file system blocks of File
func (d *Dir) Blocks() uint64 {
	if stat, ok := d.info.Sys().(*syscall.Stat_t); ok {
		return uint64(stat.Blocks)
	}
	return 0
}

// Uid returns user id of File
func (d *Dir) Uid() uint32 {
	if stat, ok := d.info.Sys().(*syscall.Stat_t); ok {
		return (stat.Uid)
	}
	return uint32(os.Geteuid())
}

// User returns user (owner) name of File
func (d *Dir) User() string {
	u, err := user.LookupId(cast.ToString(d.Uid()))
	if err != nil {
		return err.Error()
	}
	return u.Username
}

// Gid returns group id of File
func (d *Dir) Gid() uint32 {
	if stat, ok := d.info.Sys().(*syscall.Stat_t); ok {
		return (stat.Gid)
	}
	return uint32(os.Getgid())
}

// Group returns group (owner) name of File
func (d *Dir) Group() string {
	g, err := user.LookupGroupId(cast.ToString(d.Gid()))
	if err != nil {
		return err.Error()
	}
	return g.Name
}

// Dev will return dev id of File
func (d *Dir) Dev() uint64 {
	if stat, ok := d.info.Sys().(*syscall.Stat_t); ok {
		return uint64(stat.Rdev)
	}
	return 0
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

func (d *Dir) Git() *GitStatus {
	return d.git
}

func (d *Dir) XY() string {
	return d.git.XY(d.RelPath() + "/")
}

// Field returns the specified value of File according to ViewField
func (d *Dir) Field(field ViewField) string {
	switch field {
	case ViewFieldNo:
		return cast.ToString(field.Value())
	case ViewFieldINode:
		return cast.ToString(d.INode())
	case ViewFieldPermissions:
		return permissionS(d)
	case ViewFieldLinks:
		return cast.ToString(d.HDLinks())
	case ViewFieldSize:
		return "-"
	case ViewFieldBlocks:
		return "-"
	case ViewFieldUser:
		return d.User()
	case ViewFieldGroup:
		return d.Group()
	case ViewFieldModified, ViewFieldCreated, ViewFieldAccessed:
		return deDateS(d, field)
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
func (d *Dir) FieldC(fd ViewField) string {
	switch fd {
	case ViewFieldNo:
		return alNoC(d)
	case ViewFieldPermissions:
		return alPermissionC(d)
	case ViewFieldSize:
		return alSizeC(d)
	case ViewFieldBlocks:
		return alBlockC(d)
	case ViewFieldUser: //"User",
		return alUserC(d)
	case ViewFieldGroup: //"Group",
		return alGroupC(d)
	case ViewFieldGit:
		return alXYC(d)
	case ViewFieldName:
		return alNameC(d)
	default:
		return alFieldC(d, fd)
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
		w = paw.StringWidth(d.Field(field))
	}
	return w
}

//---------------------------------------------------------------------
// 實現 ISer 接口：

// IsLink() report whether File describes a symbolic link.
func (d *Dir) IsLink() bool {
	return d.isLink
	// return d.info.Mode()&os.ModeSymlink != 0
}

// IsFile reports whether File describes a regular file.
func (d *Dir) IsFile() bool {
	// return d.Mode().IsRegular()
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

//---------------------------------------------------------------------

// ReadDir 實現 fs.ReadDirFile 接口，方便遍歷目錄
func (d *Dir) ReadDir(n int) ([]DirEntryX, error) {
	// 1. reading items
	names := make([]string, 0, len(d.children))
	for name := range d.children {
		names = append(names, name)
	}

	totalEntries := len(names)
	if n <= 0 {
		n = totalEntries
	}

	dxs := make([]DirEntryX, 0, n)
	dirs := make([]DirEntryX, 0)
	files := make([]DirEntryX, 0)
	for i := d.idx; i < n && i < totalEntries; i++ {
		child := d.children[names[i]]
		if d.opt.Grouping == GroupNone {
			dxs = append(dxs, child)
		} else { //grouping items
			if child.IsDir() {
				dirs = append(dirs, child)
			} else {
				files = append(files, child)
			}
		}
		d.idx = i
	}

	// 2. sort items
	if d.opt.Grouping == GroupNone {
		d.opt.Sort(dxs)
	} else { //grouping items
		d.opt.Sort(dirs)
		d.opt.Sort(files)
		switch d.opt.Grouping {
		case Grouped:
			dxs = append(dirs, files...)
		case GroupedR:
			dxs = append(files, dirs...)
		}
	}

	return dxs, nil
}

// ====================================================================

func (d *Dir) ReadDirAll() ([]DirEntryX, error) {
	dxs, err := d.ReadDir(-1)
	d.ResetIndex()
	return dxs, err
}

func (d *Dir) ResetIndex() {
	d.idx = 0
}

// func (d *Dir) _ReadDirClose() {
// 	d.ResetIndex()
// }

func (d *Dir) Option() *VFSOption {
	return d.opt
}

func (d *Dir) SetOption(opt *VFSOption) {
	_SetOption(d, opt)
}

func _SetOption(cur *Dir, opt *VFSOption) {
	cur.opt = opt
	for _, dx := range cur.children {
		if dx.IsDir() {
			child := dx.(*Dir)
			_SetOption(child, opt)
		}
	}
}

func (d *Dir) SetViewType(viewType ViewType) {
	d.opt.ViewType = viewType
}

func (d *Dir) SetSortField(sortField SortKey) {
	d.opt.ByField = sortField
}

func (d *Dir) RelPaths() []string {
	return d.relpaths
}

// func (d *Dir) AddRelPath(relpaths ...string) {
// 	d.relpaths = append(d.relpaths, relpaths...)
// }

func (d *Dir) AddErrors(errs ...error) {
	if d.errors == nil {
		d.errors = make([]error, 0, len(errs))
	}
	d.errors = append(d.errors, errs...)
}

func (d *Dir) Errors(pad string, isRecurse bool) string {
	sb := new(strings.Builder)
	d.FprintErrors(sb, pad, isRecurse)
	return sb.String()
}

func (d *Dir) FprintErrors(w io.Writer, pad string, isRecurse bool) {
	errors := make([]error, 0)
	if isRecurse {
		errors = getAllErrors(d)
	} else {
		errors = d.errors
	}
	if len(errors) > 0 {
		for _, err := range errors {
			if paw.CnestedFMT.IsLogo {
				fmt.Fprintf(w, "%s%s %v\n", pad,
					cnested.Logos[logrus.ErrorLevel],
					paw.Cerror.Sprint(err))
			} else {
				fmt.Fprintf(w, "%s%v\n", pad, paw.Cerror.Sprint(err))
			}
		}
	}
}

func getAllErrors(d *Dir) []error {
	errors := make([]error, 0)
	if len(d.errors) > 0 {
		errors = append(errors, d.errors...)
	}
	des, _ := d.ReadDirAll()
	for _, cur := range des {
		if cur.IsDir() {
			cerrs := getAllErrors(cur.(*Dir))
			errors = append(errors, cerrs...)
		}
	}
	return errors
}

// NItems returns numbers of dirs and files of resurse this dir if isRecurse is true; otherwise, returns just this dir.
func (d *Dir) NItems(isRecurse bool) (ndirs, nfiles, nitems int) {
	level := 0
	if d.RelPath() != "." {
		level = len(strings.Split(d.RelPath(), "/"))
	}
	ndirs, nfiles = _NItems(d, level, true)
	return ndirs, nfiles, ndirs + nfiles
}

func _NItems(d *Dir, levle int, isRecurse bool) (ndirs, nfiles int) {
	if d.opt.Depth > 0 && levle > d.opt.Depth {
		return
	}
	dxs, _ := d.ReadDirAll()
	for _, de := range dxs {
		if !de.IsDir() {
			nfiles++
		} else {
			ndirs++
			child := de.(*Dir)
			if isRecurse {
				nd, nf := _NItems(child, levle+1, isRecurse)
				ndirs += nd
				nfiles += nf
			}
		}
	}
	return ndirs, nfiles
}

// DirInfoC will return the colorful string of sub-dir ( file.IsDir is true) and the width on console.
func (d *Dir) DirInfoC() (cdinf string, wdinf int) {
	nd, nf, _ := d.NItems(true)
	cnd := paw.Csnp.Sprint(nd)
	cnf := paw.Csnp.Sprint(nf)
	di := " dirs"
	fi := " files"
	cdi := paw.Cdirp.Sprintf(di)
	cfi := paw.Cdirp.Sprintf(fi)
	wdinf = len(di) + len(fi) + 4
	cdinf = fmt.Sprintf("[%v%v, %v%v]", cnd, cdi, cnf, cfi)
	return cdinf, wdinf
}

func (d *Dir) TotalSize() int64 {
	if d.opt.ViewType&ViewNoFiles != 0 {
		return 0
	}
	level := 0
	if d.RelPath() != "." {
		level = len(strings.Split(d.RelPath(), "/"))
	}
	return calcSize(d, level)
}

func calcSize(cur *Dir, level int) (size int64) {
	if cur.opt.Depth > 0 && level > cur.opt.Depth {
		return size
	}
	dxs, _ := cur.ReadDirAll()
	for _, de := range dxs {
		if de.IsDir() {
			next := de.(*Dir)
			size += calcSize(next, level+1)
		}
		if de.Mode().IsRegular() {
			size += de.Size()
		}
	}
	return size
}

func (d *Dir) SetGit(git *GitStatus) {
	d.git = git
}

func (d *Dir) CheckGitDir() {
	// paw.Logger.Debug(paw.Caller(1))
	// 1. check: if dir is GitIgnored, then marks all subfiles with GitIgnored.
	// if d.git == nil {
	// 	root := filepath.Dir(d.Path())
	// 	d.git = NewGitStatus(root)
	// }
	dxs, _ := d.ReadDirAll()
	if len(dxs) == 0 {
		return
	}
	for _, child := range dxs {
		if child.IsDir() {
			next := child.(*Dir)
			_checkGitDir(next)
			next.CheckGitDir()
		}
	}
}

func _checkGitDir(d *Dir) {
	gs := d.git.GetStatus()
	if d.git.NoGit ||
		len(d.children) < 1 ||
		gs == nil {
		return
	}
	// 1. check: if dir is GitIgnored, then marks all subfiles with GitIgnored.
	isMarkIgnored := false
	isUntracked := false
	var xy GitFileStatus
	rp := d.RelPath() + "/"
	if gxy, ok := gs[rp]; ok {
		if isXY(gxy, GitIgnored) {
			// paw.Logger.WithField("rp", rp).Debug("GitIgnored")
			isMarkIgnored = true
			xy = *gxy
		}
		if isXY(gxy, GitUntracked) {
			// paw.Logger.WithField("rp", rp).Debug("GitUntracked")
			isUntracked = true
			xy = *gxy
		}
	}

	if isMarkIgnored || isUntracked {
		markChildGit(d, &xy)
	}
}

func markChildGit(d *Dir, xy *GitFileStatus) {
	gs := d.git.GetStatus()
	ds, _ := d.ReadDirAll()
	if len(ds) == 0 {
		return
	}
	for _, child := range ds {
		var rp string
		if child.IsDir() {
			rp = child.RelPath() + "/"
		} else {
			rp = child.RelPath()
		}
		gs[rp] = &GitFileStatus{
			Staging:  xy.Staging,
			Worktree: xy.Worktree,
			Extra:    child.Name(),
		}
		if child.IsDir() {
			next := child.(*Dir)
			markChildGit(next, xy)
		}
	}
}

func isXY(xy *GitFileStatus, gcode GitStatusCode) bool {
	return xy.Staging == gcode ||
		xy.Worktree == gcode
}

func (d *Dir) CheckGitFiles() {
	// paw.Logger.Debug(paw.Caller(1))
	gs := d.git.GetStatus()
	if d.git.NoGit ||
		len(d.children) < 1 ||
		gs == nil {
		return
	}
	// 2. if any of subfiles of dir (including root) has any change of git status, set GitChanged to dir
	for _, de := range d.children {
		if de.IsDir() {
			next := de.(*Dir)
			next.setSubDirXY()
		}
	}
	d.setSubDirXY()
}

func (d *Dir) setSubDirXY() {
	gs := d.git.GetStatus()
	xs, ys := d.getSubXYs()
	// paw.Logger.WithFields(logrus.Fields{
	// 	"rp": "" + color.New(color.FgMagenta).Sprint(d.RelPath()) + "",
	// 	"xs": xs,
	// 	"ys": ys,
	// }).Debug(paw.Caller(1))
	if len(xs) > 0 || len(ys) > 0 {
		rp := d.RelPath()
		if rp != "." {
			rp += "/"
		}
		paw.Logger.WithFields(logrus.Fields{
			"rp": "" + color.New(color.FgMagenta).Sprint(rp) + "",
			"xs": xs,
			"ys": ys,
		}).Debug()
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

// func (d *Dir) XYC() string {
// 	return d.git.XYC(d.RelPath() + "/")
// }

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
		child, ok := cur.children[part]
		if !ok {
			return nil, &fs.PathError{
				Op:   "getDir",
				Path: relpath,
				Err:  fs.ErrNotExist,
			}
			// return nil, fmt.Errorf("%s is not exists", relpath)
		}

		childDir, ok := child.(*Dir)
		if !ok {
			return nil, &fs.PathError{
				Op:   "getDir",
				Path: relpath,
				Err:  fmt.Errorf("%s is not directory", relpath),
			}
			// return nil, fmt.Errorf("%s is not directory", relpath)
		}

		cur = childDir
	}

	return cur, nil
}

func (d *Dir) FprintlnSummaryC(w io.Writer, pad string, wdstty int, isRecurse bool) {
	fmt.Fprintln(w, d.SummaryC(pad, wdstty, isRecurse))
}

func (d *Dir) SummaryC(pad string, wdstty int, isRecurse bool) string {
	var (
		ndirs, nfiles, _ = d.NItems(isRecurse)
		tsize            = d.TotalSize()
	)
	if isRecurse {
		return totalSummary(pad, ndirs, nfiles, tsize, sttyWidth-2)
	} else {
		return dirSummary(pad, ndirs, nfiles, tsize, sttyWidth-2)
	}
}

func (d *Dir) FprintlnRelPathC(w io.Writer, pad string, isBg bool) {
	fmt.Fprintf(w, "%s:\n", d.RelPathC(pad, isBg))
}

func (d *Dir) RelPathC(pad string, isBg bool) string {
	var bgc []Attribute
	if isBg {
		bgc = paw.EXAColorAttributes["bgpmpt"]
	}
	rp := PathTo(d, &PathToOption{true, bgc, PRTRelPath})
	return fmt.Sprintf("%s%s", pad, rp)
	// return getRelPath(pad, "", d.RelPath(), isBg)
}
