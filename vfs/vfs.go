package vfs

import (
	"fmt"
	"io"
	"io/fs"
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
	rootDir   *Dir
	relpaths  []string
	skipConds *SkipConds
}

// NewVFSWith 創建一個唯讀文件系統的實例
func NewVFSWith(root string, level int) *VFS {
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
	v := &VFS{
		rootDir: &Dir{
			path:     aroot,
			relpath:  relpath,
			name:     name,
			info:     info,
			git:      git,
			relpaths: []string{relpath},
			children: make(map[string]fs.DirEntry),
		},
		relpaths:  []string{relpath},
		skipConds: NewSkipConds(true),
	}
	v.BuildFS(level)

	cur := v.rootDir
	v.createRDirs(cur)
	// sort.Slice(v.rdirs, func(i, j int) bool {
	// 	return strings.ToLower(v.rdirs[i]) < strings.ToLower(v.rdirs[j])
	// })
	// spew.Dump(v.relpaths)
	// spew.Dump(v.rootDir.relpaths)
	// paw.Logger.WithField("reflect.DeepEqual(v.relpaths, v.rootDir.relpaths)", reflect.DeepEqual(v.relpaths, v.rootDir.relpaths)).Debug()

	// v.rootDir.git.Dump("(vfs) ConfigGit: before")
	checkChildGitDir(v.rootDir)
	checkChildGitFiles(v.rootDir)
	// v.rootDir.checkGitDir()
	// v.rootDir.checkGitFiles()
	v.rootDir.git.Dump("checkChildGit: modified")
	return v
}

func (v *VFS) RootDir() *Dir {
	return v.rootDir
}
func (v *VFS) RelPaths() []string {
	return v.relpaths
}

func (v *VFS) SetSkipConds(skips ...*SkipFunc) {
	if len(skips) < 1 {
		return
	}
	v.skipConds = NewSkipConds(false)
	for _, skip := range skips {
		v.skipConds.Add(skip)
	}
}

func (v *VFS) AddSkipFuncs(skips ...*SkipFunc) {
	if len(skips) < 1 {
		return
	}
	for _, skip := range skips {
		v.skipConds.Add(skip)
	}
}

func (v *VFS) createRDirs(cur *Dir) (relpaths []string) {
	ds, _ := cur.ReadDir(-1)
	cur.resetIdx()
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
	d.resetIdx()
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
	d.resetIdx()
	if len(ds) == 0 {
		return
	}
	for _, child := range ds {
		dd, isDir := child.(*Dir)
		if !isDir {
			continue
		}
		dd.checkGitFiles()
		// checkChildGitFiles(dd)
	}
}

func (v *VFS) BuildFS(level int) {
	cur := v.rootDir
	buildFS(cur, cur.Path(), level, v.skipConds)
}

func buildFS(cur *Dir, root string, level int, skipcond *SkipConds) {
	dpath := cur.Path()
	git := cur.git
	nlevel := len(strings.Split(cur.RelPath(), "/"))
	if level > 0 && nlevel > level {
		return
	}
	des, _ := os.ReadDir(dpath)
	for _, de := range des {
		skip := false
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
		skip = skipcond.Is(de)
		if skip {
			continue
		}
		relpath, _ := filepath.Rel(root, path)
		xattrs, _ := getXattr(path)
		if !de.IsDir() {
			cur.children[de.Name()] = &File{
				path:    path,
				relpath: relpath,
				name:    de.Name(),
				info:    info,
				xattrs:  xattrs,
				git:     git,
			}
		} else {
			childtDir := &Dir{
				path:     path,
				relpath:  relpath,
				name:     de.Name(),
				info:     info,
				xattrs:   xattrs,
				git:      git,
				relpaths: make([]string, 0),
				children: make(map[string]fs.DirEntry),
			}
			cur.children[de.Name()] = childtDir
			if level != 0 {
				buildFS(childtDir, root, level, skipcond)
			}
		}
	}
}

func (v *VFS) DumpFS(w io.Writer) {
	color.NoColor = true
	v.View(w, DefaultViewFields, ViewLevel)
	color.NoColor = paw.NoColor
}

// getDir 通過一個路徑獲取其 dir 類型實例
func (v *VFS) getDir(path string) (*Dir, error) {
	if path == "." {
		return v.rootDir, nil
	}
	parts := strings.Split(path, "/")

	cur := v.rootDir
	for _, part := range parts {
		child := cur.children[part]
		if child == nil {
			return nil, fmt.Errorf("%s is not exists", path)
		}

		childDir, ok := child.(*Dir)
		if !ok {
			return nil, fmt.Errorf("%s is not directory", path)
		}

		cur = childDir
	}

	return cur, nil
}
