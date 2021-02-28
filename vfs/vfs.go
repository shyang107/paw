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

func (v *VFS) View(w io.Writer, fields []ViewField, viewType ViewType) {
	switch viewType {
	case ViewList:
		v.ViewList(w, fields, false)
	case ViewListX:
		v.ViewList(w, fields, true)
	// case ViewTree:
	// case ViewTreeX:
	case ViewLevel:
		v.ViewLevel(w, fields, false)
	case ViewLevelX:
		v.ViewLevel(w, fields, true)
		// case ViewTable:
		// case ViewTableX:
		// case ViewListTree:
		// case ViewListTreeX:
		// case ViewClassify:
		// default:
	}
}

func (v *VFS) DumpFS(w io.Writer) {
	color.NoColor = true
	v.View(w, DefaultViewFields, ViewLevel)
	color.NoColor = paw.NoColor
}

// func dumpFS(de *Dir, root string, level, wdidx int, fields []PDFieldFlag, nd, nf *int) {
// 	// head := getPFHeadS(chdp, fields...)
// 	des, _ := de.ReadDir(-1)
// 	de.resetIdx()
// 	if len(des) == 0 {
// 		return
// 	}
// 	pad := paw.Spaces(level * 3)
// 	if len(de.errors) > 0 {
// 		de.FprintErrors(os.Stderr, pad)
// 	}
// 	for _, child := range des {
// 		f, isFile := child.(*File)
// 		if isFile {
// 			(*nf)++
// 			sidx := cfip.Sprintf("F%-[1]*[2]d", wdidx, *nf)
// 			fmt.Printf("%s%s ", pad, sidx)
// 			for _, field := range fields {
// 				fmt.Printf("%v ", f.FieldC(field, nil))
// 			}
// 			fmt.Println()
// 		} else {
// 			(*nd)++
// 			sidx := cdip.Sprintf("D%-[1]*[2]d", wdidx, *nd)
// 			d := child.(*Dir)
// 			fmt.Printf("%s%s ", pad, sidx)
// 			for _, field := range fields {
// 				fmt.Printf("%v ", d.FieldC(field, nil))
// 			}
// 			fmt.Println()
// 			ndd, nff := d.NItems()
// 			if ndd+nff > 0 {
// 				// fmt.Printf("%s%s %v\n", pad, paw.Spaces(2*wdidx+2), head)
// 				dumpFS(d, root, level+1, wdidx, fields, nd, nf)
// 			}
// 		}
// 	}
// }

// Open 實現 fs.FS 的 Open 方法
// An FS provides access to a hierarchical file system.
// The FS interface is the minimum implementation required of the file system. A file system may implement additional interfaces, such as ReadFileFS, to provide additional or optimized functionality.
// type FS interface {
//     // Open opens the named file.
//     //
//     // When Open returns an error, it should be of type *PathError
//     // with the Op field set to "open", the Path field set to name,
//     // and the Err field describing the problem.
//     //
//     // Open should reject attempts to open names that do not satisfy
//     // ValidPath(name), returning a *PathError with Err set to
//     // ErrInvalid or ErrNotExist.
//     Open(name string) (File, error)
// }

// func (v *VFS) Open(name string) (fs.File, error) {
// 	// 1、校驗 name
// 	if !fs.ValidPath(name) {
// 		return nil, &fs.PathError{
// 			Op:   "open",
// 			Path: name,
// 			Err:  fs.ErrInvalid,
// 		}
// 	}

// 	// 2、根目錄處理
// 	if name == "." || name == "" {
// 		// 重置目錄的遍歷
// 		v.rootDir.idx = 0
// 		return v.rootDir, nil
// 	}

// 	// 3、根據 name 在目錄樹中進行查找
// 	cur := v.rootDir
// 	parts := strings.Split(name, "/")
// 	for i, part := range parts {
// 		// 不存在返回錯誤
// 		child := cur.children[part]
// 		if child == nil {
// 			return nil, &fs.PathError{
// 				Op:   "open",
// 				Path: name,
// 				Err:  fs.ErrNotExist,
// 			}
// 		}

// 		// 是否是文件
// 		f, ok := child.(*file)
// 		if ok {
// 			// 文件名是最後一項
// 			if i == len(parts)-1 {
// 				return f, nil
// 			}

// 			return nil, &fs.PathError{
// 				Op:   "open",
// 				Path: name,
// 				Err:  fs.ErrNotExist,
// 			}
// 		}

// 		// 是否是目錄
// 		d, ok := child.(*dir)
// 		if !ok {
// 			return nil, &fs.PathError{
// 				Op:   "open",
// 				Path: name,
// 				Err:  errors.New("not a directory"),
// 			}
// 		}
// 		// 重置，避免遍歷問題
// 		d.idx = 0

// 		cur = d
// 	}

// 	return cur, nil
// }

// =====================================

// // MkdirAll 這不是 `io/fs` 的要求，但一個文件系統目錄樹需要可以構建
// // 這個方法就是用來創建目錄
// func (v *VFS) MkdirAll(path string) error {
// 	if !fs.ValidPath(path) {
// 		return errors.New("Invalid path")
// 	}

// 	if path == "." {
// 		return nil
// 	}

// 	cur := v.rootDir
// 	parts := strings.Split(path, "/")
// 	for _, part := range parts {
// 		child := cur.children[part]
// 		if child == nil {
// 			childDir := &dir{
// 				name:     part,
// 				mTime:    time.Now(),
// 				children: make(map[string]fs.DirEntry),
// 			}
// 			cur.children[part] = childDir
// 			cur = childDir
// 		} else {
// 			childDir, ok := child.(*dir)
// 			if !ok {
// 				return fmt.Errorf("%s is not directory", part)
// 			}

// 			cur = childDir
// 		}
// 	}

// 	return nil
// }

// // WriteFile 也不是 `io/fs` 的要求，和 MkdirAll 類似，文件內容也需要有接口寫入
// func (v *VFS) WriteFile(name, content string) error {
// 	if !fs.ValidPath(name) {
// 		return &fs.PathError{
// 			Op:   "write",
// 			Path: name,
// 			Err:  fs.ErrInvalid,
// 		}
// 	}

// 	var err error
// 	dir := v.rootDir

// 	path := filepath.Dir(name)
// 	if path != "." {
// 		dir, err = v.getDir(path)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	filename := filepath.Base(name)

// 	dir.children[filename] = &file{
// 		name:    filename,
// 		content: bytes.NewBufferString(content),
// 		mTime:   time.Now(),
// 	}

// 	return nil
// }

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
