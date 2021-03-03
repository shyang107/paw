package vfs

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
)

// func init() {
// 	paw.Logger.SetLevel(logrus.TraceLevel)
// }

// VFS 是 fs.FS 的唯讀文件系統實現
type VFS struct {
	Dir
	relpaths []string
	// skipConds *SkipConds
	opt *VFSOption
}

// NewVFSWith 創建一個唯讀文件系統的實例
func NewVFS(root string, opt *VFSOption) *VFS {
	paw.Logger.Info()

	aroot, err := filepath.Abs(root)
	if err != nil {
		return nil
	}
	info, err := os.Lstat(aroot)
	if err != nil {
		return nil
	}

	if !info.IsDir() {
		return nil
	}
	git := NewGitStatus(aroot)
	relpath, _ := filepath.Rel(aroot, aroot)
	name := filepath.Base(aroot)

	paw.Logger.Trace("checking VFSOption...")
	checkOpt(opt)

	v := &VFS{
		Dir: Dir{
			path:     aroot,
			relpath:  relpath,
			name:     name,
			info:     info,
			git:      git,
			relpaths: []string{relpath},
			children: make(map[string]DirEntryX),
			opt:      opt,
		},
		relpaths: []string{relpath},
		opt:      opt,
	}

	return v
}

func checkOpt(opt *VFSOption) {
	if opt == nil {
		opt = NewVFSOption()
	} else {
		if opt.Grouping == 0 {
			opt.Grouping = GroupNone
		}
		if opt.By == nil {
			opt.By = &ByLowerNameFunc
		}
		if opt.Skips == nil {
			opt.Skips = NewSkipConds().Add(DefaultSkip)
		}
		if opt.ViewType == 0 {
			opt.ViewType = ViewList
		}
		if opt.ViewFields == 0 {
			opt.ViewFields = DefaultViewField
		}
	}
}

func (v *VFS) RootDir() *Dir {
	return &v.Dir
}
func (v *VFS) RelPaths() []string {
	return v.relpaths
}

func (v *VFS) Option() *VFSOption {
	return v.opt
}

func (v *VFS) SetOption(opt *VFSOption) {
	v.RootDir().SetOption(opt)
}

func (v *VFS) ViewType() ViewType {
	return v.opt.ViewType
}

func (v *VFS) SetViewType(viewType ViewType) {
	v.RootDir().SetViewType(viewType)
}

func (v *VFS) SetSkipConds(skips ...Skiper) {
	if len(skips) < 1 {
		return
	}
	v.opt.Skips = NewSkipConds().Add(skips...)

}

func (v *VFS) AddSkipFuncs(skips ...Skiper) {
	if len(skips) < 1 {
		return
	}
	v.opt.Skips.Add(skips...)
}

func (v *VFS) BuildFS() {
	paw.Logger.Trace("building VFS...")
	cur := &v.Dir
	buildFS(cur, cur.Path())

	paw.Logger.Trace("building VFS.relpaths...")
	v.createRDirs(&v.Dir)

	paw.Logger.Trace("checking VFS.git: dir...")
	checkChildGitDir(&v.Dir)

	paw.Logger.Trace("checking VFS.git: file...")
	checkChildGitFiles(&v.Dir)

	v.git.Dump("checkChildGit: modified")
}

func (v *VFS) createRDirs(cur *Dir) (relpaths []string) {
	ds, _ := cur.ReadDir(-1)
	cur.ResetIndex()
	relpaths = make([]string, 0) //
	for _, d := range ds {
		next, isDir := d.(*Dir)
		if isDir {
			relpaths = append(relpaths, next.RelPath())
			v.relpaths = append(v.relpaths, next.RelPath())
			nextrelpaths := v.createRDirs(next)
			relpaths = append(relpaths, nextrelpaths...)
		}
	}
	cur.relpaths = append(cur.relpaths, relpaths...)
	if len(cur.relpaths) > 0 {
		sort.Slice(cur.relpaths, func(i, j int) bool {
			return strings.ToLower(cur.relpaths[i]) < strings.ToLower(cur.relpaths[j])
		})
	}
	return relpaths
}

func checkChildGitDir(d *Dir) {
	ds, _ := d.ReadDir(-1)
	d.ResetIndex()
	if len(ds) == 0 {
		return
	}
	for _, child := range ds {
		dd, isDir := child.(*Dir)
		if !isDir {
			continue
		}
		dd.checkGitDir()
		checkChildGitDir(dd)
	}
}
func checkChildGitFiles(d *Dir) {
	ds, _ := d.ReadDir(-1)
	d.ResetIndex()
	if len(ds) == 0 {
		return
	}
	for _, child := range ds {
		dd, isDir := child.(*Dir)
		if !isDir {
			continue
		}
		dd.checkGitFiles()
	}
}

func buildFS(cur *Dir, root string) {
	var (
		dpath = cur.Path()
		git   = cur.git
		skip  = cur.opt.Skips
		level = cur.opt.Depth
	)
	nlevel := len(strings.Split(cur.RelPath(), "/"))
	if level > 0 && nlevel > level {
		return
	}
	des, _ := os.ReadDir(dpath)
	for _, de := range des {
		path := filepath.Join(dpath, de.Name())
		info, err := os.Lstat(path)
		if err != nil {
			if cur.errors == nil {
				cur.errors = []error{}
			}
			cur.errors = append(cur.errors, err)
			// cur.errors = append(cur.errors, &fs.PathError{
			// 	Op:   "os", // "buildFS",
			// 	Path: path,
			// 	Err:  err,
			// })
			continue
		}
		relpath, _ := filepath.Rel(root, path)
		xattrs, _ := getXattr(path)
		var child DirEntryX
		if !de.IsDir() {
			child = &File{
				path:    path,
				relpath: relpath,
				name:    de.Name(),
				info:    info,
				xattrs:  xattrs,
				git:     git,
			}
		} else {
			child = &Dir{
				path:     path,
				relpath:  relpath,
				name:     de.Name(),
				info:     info,
				xattrs:   xattrs,
				git:      git,
				relpaths: make([]string, 0),
				children: make(map[string]DirEntryX),
				opt:      cur.opt,
			}
		}

		if skip.Is(child) {
			continue
		}

		cur.children[de.Name()] = child

		if level != 0 && child.IsDir() {
			buildFS(child.(*Dir), root)
		}
	}
}

func (v *VFS) DumpFS(w io.Writer) {
	color.NoColor = true
	v.View(w)
	color.NoColor = paw.NoColor
}

// getDir 通過一個路徑獲取其 dir 類型實例
func (v *VFS) getDir(path string) (*Dir, error) {
	return v.Dir.getDir(path)
	// if path == "." {
	// 	return &v.Dir, nil
	// }
	// parts := strings.Split(path, "/")

	// cur := &v.Dir
	// for _, part := range parts {
	// 	child := cur.children[part]
	// 	if child == nil {
	// 		return nil, fmt.Errorf("%s is not exists", path)
	// 	}

	// 	childDir, ok := child.(*Dir)
	// 	if !ok {
	// 		return nil, fmt.Errorf("%s is not directory", path)
	// 	}

	// 	cur = childDir
	// }

	// return cur, nil
}
