package filetree

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/shyang107/paw/treeprint"
)

// FileMap ...
type FileMap map[string][]*File

// FileList ...
type FileList struct {
	root  string   // root directory
	store FileMap  // all files in `root` directory
	dirs  []string // keys of `store`
}

// NewFileList will return the instance of `FileList`
func NewFileList(root string) FileList {
	if len(root) == 0 {
		return FileList{}
	}
	return FileList{
		root:  root,
		store: make(map[string][]*File),
		dirs:  []string{},
	}
}

// Root will return the `root` field (root directory)
func (f *FileList) Root() string {
	return f.root
}

// Map will retun the `FileMap`
func (f *FileList) Map() FileMap {
	return f.store
}

// Dirs will retun keys of `FileMap`
func (f *FileList) Dirs() []string {
	return f.dirs
}

// NDirs is the numbers of sub-directories of `root`
func (f *FileList) NDirs() int {
	return len(f.Dirs()) - 1
}

// NFiles is the numbers of all files
func (f *FileList) NFiles() int {
	var nf int
	for _, dir := range f.Dirs() {
		nf += (len(f.store[dir]) - 1)
	}
	return nf
}

// AddFile will add file into the file list
func (f *FileList) AddFile(file *File) {
	if _, ok := f.store[file.Dir]; !ok {
		f.store[file.Dir] = []*File{}
		f.dirs = append(f.dirs, file.Dir)
	}
	f.store[file.Dir] = append(f.store[file.Dir], file)
}

// SkipFile is used as a return value from IgnoreFn to indicate that
// the regular file named in the call is to be skipped. It is not returned
// as an error by any function.
var SkipFile = errors.New("skip the file")

// SkipDir is used as a return value from WalkFuncs to indicate that
// the directory named in the call is to be skipped. It is not returned
// as an error by any function.
var SkipDir = filepath.SkipDir

// IgnoreFn is the type of the function called for each file or directory
// visited by FindFiles. The f argument contains the File argument to FindFiles.
//
// If there was a problem walking to the file or directory named by path, the
// incoming error will describe the problem and the function can decide how
// to handle that error (and FindFiles will not descend into that directory). In the
// case of an error, the info argument will be nil. If an error is returned,
// processing stops. The sole exception is when the function returns the special
// value SkipDir or SkipFile. If the function returns SkipDir when invoked on a directory,
// FindFiles skips the directory's contents entirely. If the function returns SkipDir
// when invoked on a non-directory file, FindFiles skips the remaining files in the
// containing directory.
// If the returned error is SkipFile when inviked on a file, FindFiles will skip the file.
type IgnoreFn func(f *File, err error) error

// FindFiles will find files using codintion `ignore` func
// 	depth : depth of subfolders
// 		< 0 : walk through all directories of {root directory}
// 		0 : {root directory}/*
// 		1 : {root directory}/{level 1 directory}/*
//		...
// 	ignore IgnoreFn func(f *File, err error) error
// 		ignoring condition of files or directory
func (f *FileList) FindFiles(depth int, ignore IgnoreFn) error {
	root := f.Root()
	switch {
	case depth == 0: //{root directory}/*
		fis, err := ioutil.ReadDir(root)
		if err != nil {
			return errors.New(root + ": " + err.Error())
		}

		for _, fi := range fis {
			file := ConstructFileRelTo(root+PathSeparator+fi.Name(), root)
			err := ignore(file, nil)
			if err == SkipFile {
				continue
			}
			f.AddFile(file)
		}
	default: //walk through all directories of {root directory}
		visit := func(path string, info os.FileInfo, err error) error {
			file := ConstructFileRelTo(path, root)
			idepth := len(file.DirSlice()) - 1
			if depth > 0 {
				if idepth > depth {
					return nil
				}
			}
			err1 := ignore(file, err)
			if err1 == SkipFile {
				return nil
			}
			if err1 == SkipDir {
				return err1
			}
			f.AddFile(file)
			return nil
		}

		err := filepath.Walk(root, visit)
		if err != nil {
			return errors.New(root + ": " + err.Error())
		}
	}
	return nil
}

// ToTreeString will return the string of tree view
func (f *FileList) ToTreeString() string {

	tree := treeprint.New()

	dirs := f.Dirs()
	nd := len(dirs) // including root
	ntf := 0
	var one, pre treeprint.Tree
	fm := f.Map()
	// var branch []treeprint.Tree
	for i, dir := range dirs {
		files := f.Map()[dir]
		nf := len(files) // including dir
		ntf += (nf - 1)
		for _, file := range files {
			if file.IsDir() {
				if i == 0 { // root dir
					tree.SetValue(fmt.Sprintf("%v\n%v", file.LSColorString(file.Path), file.LSColorString(file.Dir)))
					tree.SetMetaValue(fmt.Sprintf("%v dirs., %v files", nd-1, nf-1))
					one = tree
				} else {
					pre = preTree(dir, fm, tree)
					one = pre.AddMetaBranch(nf-1, file)
				}
				continue
			}
			// add file node
			if i == 0 {
			}
			one.AddNode(file)
		}
	}
	buf := new(bytes.Buffer)
	buf.Write(tree.Bytes())
	buf.WriteByte('\n')
	buf.WriteString(fmt.Sprintf("%d directoris, %d files.", f.NDirs(), f.NFiles()))
	return string(buf.Bytes())
}

func preTree(dir string, fm FileMap, tree treeprint.Tree) treeprint.Tree {
	dd := strings.Split(dir, PathSeparator)
	nd := len(dd)
	var pre treeprint.Tree
	if nd == 2 { // ./xx
		pre = tree
	} else { //./xx/...
		pre = tree
		for i := 2; i < nd; i++ {
			predir := strings.Join(dd[:i], PathSeparator)
			f := fm[predir][0] // import dir
			pre = pre.FindByValue(f)
		}
	}
	return pre
}
