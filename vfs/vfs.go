package vfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
func NewVFS(root string, opt *VFSOption) (*VFS, error) {
	paw.Logger.Debug(root)
	// if !fs.ValidPath(root) {
	// 	err := &fs.PathError{
	// 		Op:   "NewVFS",
	// 		Path: root,
	// 		Err:  fs.ErrInvalid,
	// 	}
	// 	return nil, err
	// }

	info, err := os.Stat(root)
	if err != nil {
		return nil, &fs.PathError{
			Op:   "NewVFS",
			Path: root,
			Err:  err,
		}
	}

	if !info.IsDir() {
		return nil, &fs.PathError{
			Op:   "NewVFS",
			Path: root,
			Err:  fmt.Errorf("%s", "not a directory"),
		}
	}

	git := NewGitStatus(root)
	opt.ViewFields = opt.ViewFields.RemoveGit(git.NoGit)

	relpath, _ := filepath.Rel(root, root)
	// name := filepath.Base(root)

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

	dir, err := NewDir(root, root, git, opt)
	if err != nil {
		return nil, err
	}

	v := &VFS{
		Dir:      *dir,
		relpaths: []string{relpath},
		opt:      opt,
	}

	return v, nil
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
	v.opt = opt
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

	buildVFSwalk(cur, cur.Path())
	// buildVFS(cur, cur.Path(), 0)
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

func buildVFSwalk(cur *Dir, root string) {
	var (
		this = cur
		git  = cur.git
		skip = cur.opt.Skips
		dirs = make(map[string]*Dir)
		ok   bool
	)
	dirs["."] = cur
	rfs := os.DirFS(root)
	err := fs.WalkDir(rfs, ".", func(path string, d fs.DirEntry, err error) error {
		dir := filepath.Dir(path)
		level := len(strings.Split(path, "/"))
		if !cur.opt.IsForceRecurse &&
			cur.opt.Depth > 0 &&
			level > cur.opt.Depth || err == fs.SkipDir {
			return nil
		}
		if err != nil {
			this.AddErrors(&fs.PathError{
				Op:   "WalkDir",
				Path: path,
				Err:  err,
			})
			// paw.Error.Printf("WalkDirFunc[dir %q, path %q]: %v", dir, path, err)
			return nil
		}
		if skip.IsSkip(d) {
			if path != "." && d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		var child DirEntryX
		fpath := filepath.Join(root, path)
		if !d.IsDir() {
			child, err = NewFile(fpath, root, git)
		} else {
			child, err = NewDir(fpath, root, git, cur.opt)
		}
		if err != nil {
			this.AddErrors(&fs.PathError{
				Op:   "buildVFSwalk",
				Path: path,
				Err:  err,
			})
			// paw.Error.Printf("[dir:%q, path %q]: %v", dir, path, err)
			return nil
		}
		if d.IsDir() {
			if _, ok = dirs[path]; !ok {
				dirs[path] = child.(*Dir)
			}
		}
		dir = filepath.Dir(path)
		this = dirs[dir]
		this.children[d.Name()] = child
		return nil
	})
	if err != nil {
		cur.AddErrors(err)
		return
	}
}

func buildVFS(cur *Dir, root string, level int) {
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

	des, err := os.ReadDir(dpath)
	if err != nil {
		cur.AddErrors(&fs.PathError{
			Op:   "ReadDir",
			Path: dpath,
			Err:  err,
		})
		// return
	}
	for _, d := range des {
		if skip.IsSkip(d) {
			continue
		}
		path := filepath.Join(dpath, d.Name())
		// _, err := os.Lstat(path)
		// if err != nil {
		// 	cur.AddErrors(err)
		// 	// cur.errors = append(cur.errors, err)
		// 	// cur.errors = append(cur.errors, &fs.PathError{
		// 	// 	Op:   "os", // "buildFS",
		// 	// 	Path: path,
		// 	// 	Err:  err,
		// 	// })
		// 	continue
		// }
		// relpath, _ := filepath.Rel(root, path)
		// xattrs, _ := GetXattr(path)
		var child DirEntryX
		if !d.IsDir() {
			child, err = NewFile(path, root, git)
		} else {
			child, err = NewDir(path, root, git, cur.opt)
		}
		if err != nil {
			cur.AddErrors(&fs.PathError{
				Op:   "buildVFS",
				Path: path,
				Err:  err,
			})
			continue
		}
		// if skip.IsSkip(child) {
		// 	continue
		// }

		cur.children[d.Name()] = child

		// paw.Logger.WithFields(logrus.Fields{
		// 	"name":  child.Name(),
		// 	"IsDir": child.IsDir(),
		// 	"depth": cur.opt.Depth,
		// }).Trace()
		if cur.opt.IsForceRecurse {
			if child.IsDir() {
				buildVFS(child.(*Dir), root, 0)
			}
		} else {
			if cur.opt.Depth != 0 && child.IsDir() {
				buildVFS(child.(*Dir), root, level+1)
			}
		}
	}
}

var _rps = make(map[string]string)

func (v *VFS) createRDirs(cur *Dir) (relpaths []string) {
	ds, _ := cur.ReadDirAll()
	nd, _, _ := cur.NItems(true)
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

// getDir 通過一個路徑獲取其 dir 類型實例
func (v *VFS) getDir(path string) (*Dir, error) {
	return v.Dir.getDir(path)
}
