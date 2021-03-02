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
	relpaths  []string
	skipConds *SkipConds
	level     int
}

// NewVFSWith 創建一個唯讀文件系統的實例
func NewVFS(root string, level int, sortFunc *ByFunc) *VFS {
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
	if sortFunc == nil {
		sortFunc = &ByLowerNameFunc
	}

	v := &VFS{
		Dir: Dir{
			path:     aroot,
			relpath:  relpath,
			name:     name,
			info:     info,
			git:      git,
			relpaths: []string{relpath},
			children: make(map[string]DirEntryX),
			sortFunc: sortFunc,
		},
		relpaths:  []string{relpath},
		skipConds: NewSkipConds().Add(DefaultSkip),
		level:     level,
	}

	return v
}

func NewVFSWithSortKey(root string, level int, sortKey SortKey) *VFS {
	paw.Logger.Info()

	var sortFunc *ByFunc
	if less, ok := SortFuncMap[sortKey]; ok {
		sortFunc = less
	} else {
		sortFunc = &ByLowerNameFunc
	}
	return NewVFS(root, level, sortFunc)
}

func (v *VFS) RootDir() *Dir {
	return &v.Dir
}
func (v *VFS) RelPaths() []string {
	return v.relpaths
}

func (v *VFS) SetSkipConds(skips ...Skiper) {
	if len(skips) < 1 {
		return
	}
	v.skipConds = NewSkipConds().Add(skips...)

}

func (v *VFS) AddSkipFuncs(skips ...Skiper) {
	if len(skips) < 1 {
		return
	}
	v.skipConds.Add(skips...)
}

func (v *VFS) BuildFS() {
	paw.Logger.Trace("building VFS...")
	cur := &v.Dir
	buildFS(cur, cur.Path(), v.level, v.skipConds)

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

func buildFS(cur *Dir, root string, level int, skip *SkipConds) {
	dpath := cur.Path()
	git := cur.git
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
				sortFunc: cur.sortFunc,
			}
		}

		if skip.Is(child) {
			continue
		}

		cur.children[de.Name()] = child

		if level != 0 && child.IsDir() {
			buildFS(child.(*Dir), root, level, skip)
		}
	}
}

func (v *VFS) DumpFS(w io.Writer) {
	color.NoColor = true
	v.View(w, DefaultViewField, ViewLevel)
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

// func (v *VFS) ReadDir(n int) ([]DirEntryX, error) {
// 	return v.Dir.ReadDir(n)
// }

// // implement DirEntry interface of VFS

// func (v *VFS) Name() string {
// 	return v.Dir.Name()
// }

// func (v *VFS) IsDir() bool {
// 	return v.Dir.IsDir()
// }

// func (v *VFS) Type() fs.FileMode {
// 	return v.Dir.Type()
// }

// func (v *VFS) Info() (fs.FileInfo, error) {
// 	return v.Dir.Info()
// }
