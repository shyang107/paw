package vfs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
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
	paw.Logger.Debug(root)

	aroot, err := filepath.Abs(root)
	if err != nil {
		paw.Error.Fatal(err)
	}
	info, err := os.Stat(aroot)
	if err != nil {
		paw.Error.Fatal(err)
	}

	if !info.IsDir() {
		return nil
	}

	git := NewGitStatus(aroot)
	opt.ViewFields = opt.ViewFields.RemoveGit(git.NoGit)

	relpath, _ := filepath.Rel(aroot, aroot)
	name := filepath.Base(aroot)

	opt.Check()

	paw.Logger.WithFields(logrus.Fields{
		"Depth":          opt.Depth,
		"IsForceRecurse": opt.IsForceRecurse,
		"Grouping":       opt.Grouping,
		"ByField":        opt.ByField,
		"Skips":          opt.Skips,
		"ViewFields":     opt.ViewFields,
		"ViewType":       opt.ViewType,
	}).Debug()

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
	paw.Logger.Debug("building VFS...")
	cur := v.RootDir()

	buildFS(cur, cur.Path(), 0)
	// nd, nf := cur.NItems()
	// paw.Logger.WithFields(logrus.Fields{
	// 	"nd": nd,
	// 	"nf": nf,
	// }).Debug()

	paw.Logger.Debug("building VFS.relpaths...")
	v.createRDirs(&v.Dir)

	paw.Logger.Tracef("checking VFS.git: dir...[%q]", cur.RelPath())
	cur.CheckGitDir()

	paw.Logger.Tracef("checking VFS.git: files...[%q]", cur.RelPath())
	cur.CheckGitFiles()

	v.git.Dump("checkChildGit: modified")
}

func buildFS(cur *Dir, root string, level int) {
	var (
		dpath = cur.Path()
		git   = cur.git
		skip  = cur.opt.Skips
	)
	if !cur.opt.IsForceRecurse &&
		cur.opt.Depth > 0 &&
		level > cur.opt.Depth {
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
		xattrs, _ := GetXattr(path)
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
		if skip.IsSkip(child) {
			continue
		}

		cur.children[de.Name()] = child

		if cur.opt.IsForceRecurse {
			if child.IsDir() {
				buildFS(child.(*Dir), root, 0)
			}
		} else {
			if cur.opt.Depth != 0 && child.IsDir() {
				buildFS(child.(*Dir), root, level+1)
			}
		}
	}
}

func (v *VFS) createRDirs(cur *Dir) (relpaths []string) {
	ds, _ := cur.ReadDirAll()
	nd, _, _ := cur.NItems()
	relpaths = make([]string, 0, nd) //
	for _, d := range ds {
		if d.IsDir() {
			next := d.(*Dir)
			relpaths = append(relpaths, next.RelPath())
			v.relpaths = append(v.relpaths, next.RelPath())
			nextrelpaths := v.createRDirs(next)
			relpaths = append(relpaths, nextrelpaths...)
		}
	}
	cur.relpaths = append(cur.relpaths, relpaths...)
	// if len(cur.relpaths) > 0 {
	// 	sort.Sort(ByLowerString(cur.relpaths))
	// }
	return relpaths
}

func (v *VFS) DumpFS(w io.Writer) {
	color.NoColor = true
	v.View(w)
	color.NoColor = paw.NoColor
}

// getDir 通過一個路徑獲取其 dir 類型實例
func (v *VFS) getDir(path string) (*Dir, error) {
	return v.Dir.getDir(path)
}
