package filetree

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/shyang107/paw"
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

func (f FileList) String() string {
	var (
		w    = new(bytes.Buffer)
		dirs = f.Dirs()
		fm   = f.Map()
	)

	j := 0
	for i, dir := range dirs {
		for _, file := range fm[dir] {
			if file.IsDir() {
				if strings.EqualFold(file.Dir, RootMark) {
					fmt.Fprintf(w, "D%v. root (%v)\n", i, file.LSColorString(f.Root()))
				} else {
					fmt.Fprintf(w, "D%v. subfolder: %v\n", i, file.LSColorString(file.Dir))
				}
				continue
			}
			j++
			fmt.Fprintf(w, "  %v. %v\n", j, file)
		}

		if i == len(dirs)-1 {
			break
		}
	}

	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "%d directories, %d Files\n", f.NDirs(), f.NFiles())
	return string(w.Bytes())
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

// ToTreeString will return the string of FileList in tree form
func (f *FileList) ToTreeString(pad string) string {
	return string(f.ToTree(pad))
}

// ToTree will return the []byte of FileList in tree form
func (f *FileList) ToTree(pad string) []byte {

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
	// buf.WriteByte('\n')
	b := make([]byte, len(buf.Bytes()))
	b = append(b, pad...)
	for _, v := range buf.Bytes() {
		b = append(b, v)
		if v == '\n' {
			b = append(b, pad...)
		}
	}
	return b
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

// ToTableString will return the string of FileList in table form
func (f *FileList) ToTableString(pad string) string {
	return string(f.ToTable(pad))
}

// ToTable will return the []byte of FileList in table form
func (f *FileList) ToTable(pad string) []byte {

	var (
		buf    = new(bytes.Buffer)
		nDirs  = f.NDirs()
		nFiles = f.NFiles()
		dirs   = f.Dirs()
		fm     = f.Map()
	)

	tf := &paw.TableFormat{
		Fields:    []string{"No.", "Files"},
		LenFields: []int{5, 80},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft},
		Padding:   pad,
	}
	tf.Prepare(buf)

	head := fmt.Sprintf("Root directory: %q", f.Root())
	tf.SetBeforeMessage(head)

	tf.PrintSart()

	for i, dir := range dirs {
		for j, file := range fm[dir] {
			if file.IsDir() {
				if strings.EqualFold(file.Dir, RootMark) {
					tf.PrintRow("", fmt.Sprintf("[%v]. root (%v)", i, file.LSColorString(f.Root())))
				} else {
					tf.PrintRow("", fmt.Sprintf("[%v]. subfolder: %v", i, file.LSColorString(file.Dir)))
				}
				continue
			}
			tf.PrintRow(j, file)
		}

		tf.PrintRow("", fmt.Sprintf("Sum: %v files.", len(fm[dir])-1))

		if i == len(dirs)-1 {
			break
		}
		tf.PrintMiddleSepLine()
	}
	tf.SetAfterMessage(fmt.Sprintf("\n%v directories, %v files", nDirs, nFiles))

	tf.PrintEnd()

	return buf.Bytes()
}

// ToTextString will return the string of FileList in table form
func (f *FileList) ToTextString(pad string) string {
	return string(f.ToText(pad))
}

// ToText will return the []byte of FileList in table form
func (f *FileList) ToText(pad string) []byte {
	var (
		w     = new(bytes.Buffer)
		dirs  = f.Dirs()
		fm    = f.Map()
		width = 80
	)

	// fmt.Fprintln(w, pad)
	fmt.Fprintf(w, "%sRoot directory: %q\n", pad, f.Root())
	fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("=", width))

	ppad := ""
	for i, dir := range dirs {
		for j, file := range fm[dir] {
			if file.IsDir() {
				if strings.EqualFold(file.Dir, RootMark) {
					fmt.Fprintf(w, "%sD%v. root (%v)\n", pad, i, file.LSColorString(f.Root()))
				} else {
					ppad = strings.Repeat("   ", len(file.DirSlice())-1)
					fmt.Fprintf(w, "%sD%v. subfolder: %v\n", pad+ppad, i, file.LSColorString(file.Dir))
				}
				continue
			}
			fmt.Fprintf(w, "%s  %v. %v\n", pad+ppad, j, file)
		}

		fmt.Fprintf(w, "%s  Sum: %v files.\n", pad+ppad, len(fm[dir])-1)

		if i == len(dirs)-1 {
			break
		}
		fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("-", width))
	}

	fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("=", width))
	fmt.Fprintln(w, pad)
	fmt.Fprintf(w, "%s%d directories, %d Files\n", pad, f.NDirs(), f.NFiles())
	// fmt.Fprintln(w, pad)
	return w.Bytes()
}
