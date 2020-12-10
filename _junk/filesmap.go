package paw

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/shyang107/paw/treeprint"
)

// FilesMap store files ordered by folder, use folder name as the key
type FilesMap struct {
	root         string // if root == nil or "", no root directory
	nDirectories int    // number of subfolders icluding empty-files folders
	nFiles       int    // number of files
	// store files in specific folder into store["folder"]
	store map[string][]File
	// store uniq folder into keys
	keys []string
}

// NewFilesMap is constructor of `FilesMap`
func NewFilesMap() *FilesMap {
	return &FilesMap{
		root:  "",
		store: map[string][]File{},
		keys:  []string{},
	}
}

// NewFilesMapFrom is constructor of `FilesMap` from `files`
func NewFilesMapFrom(files []File) *FilesMap {
	o := &FilesMap{
		root:  "",
		store: map[string][]File{},
		keys:  []string{},
	}
	for _, f := range files {
		if f.ShortFolder == "." {
			o.root = f.Folder
		}
		o.SetOne(f.ShortFolder, f)
	}
	return o
}

func (m FilesMap) String() string {
	buf := new(bytes.Buffer)
	m.Fprint(buf, OPlainTextMode, "", "")
	return TrimFrontEndSpaceLine(string(buf.Bytes()))
}

// Keys will return the keys of key-value pairs
//
// Example:
// 	m := NewFilesMap()
// 	for _, key := range m.Keys() {
// 		value, _:= m.Get(key)
// 		fmt.Println(key, value)
// 	}
func (m *FilesMap) Keys() []string {
	return m.keys
}

// GetRoot will return the root directory of `FilesMap`
func (m *FilesMap) GetRoot() string {
	return m.root
}

// SetRoot will store root directory of `FilesMap`
func (m *FilesMap) SetRoot(root string) {
	m.root = root
}

// // NDirectories will return the number of subfolders of root of `FilesMap`
// func (m *FilesMap) NDirectories() int {
// 	return m.ndir
// }

// NDirectories will return the number of subfolders of root of `FilesMap`
func (m *FilesMap) NDirectories() int {
	// nSub, _ := m.calFiles()
	// return nSub
	return m.nDirectories
}

// NFiles will return the number of files of `FilesMap`
func (m *FilesMap) NFiles() int {
	// _, nFiles := m.calFiles()
	// return nFiles
	return m.nFiles
}

var rootHasFile = false

func (m *FilesMap) calFiles() (nSub, nFiles int) {
	nSub = 0
	for _, files := range m.store {
		nSub++
		nFiles += len(files)
		for _, f := range files {
			if strings.EqualFold(f.File, zeroRootFiles) {
				rootHasFile = false
				nFiles--
				break
			}
			if strings.EqualFold(m.root, f.Folder) {
				rootHasFile = true
				break
			}
		}
	}
	// nSub--
	nSub = len(m.keys) - 1
	return nSub, nFiles
}

// Getter ans Setter

// Get will return the value (`[]File`) associated with the key (`folder`).
// If the key does not exist, the second return value will be `false`.
// 	Here, `key` is the name of specific folder.
func (m *FilesMap) Get(key string) ([]File, bool) {
	val, exist := m.store[key]
	return val, exist
}

// GetAll will return the folder-[]file pairs
func (m *FilesMap) GetAll() map[string][]File {
	return m.store
}

// GetFiles will return a copy of all files in store
func (m *FilesMap) GetFiles() []File {
	files := []File{}
	for _, v := range m.store {
		files = append(files, v...)
	}
	return files
}

// Set will store a key-value (folder-files) pair.
// If the key already exists, it will overwrite the existing key-value pair.
func (m *FilesMap) Set(key string, val []File) {
	if _, exist := m.store[key]; !exist {
		m.keys = append(m.keys, key)
	}
	m.store[key] = val
	m.nDirectories++
	m.nFiles += len(val)
	for _, f := range val {
		if f.IsDir() {
			m.nFiles--
		}
	}
}

// SetOne will store a key-value (folder-one_file) pair.
// If the key already exists, it will overwrite the existing key-value pair.
func (m *FilesMap) SetOne(key string, val File) {
	if _, exist := m.store[key]; !exist {
		m.keys = append(m.keys, key)
	}
	for _, file := range m.store[key] {
		if reflect.DeepEqual(file, val) {
			return
		}
	}
	m.store[key] = append(m.store[key], val)
	if !val.IsDir() {
		m.nFiles++
	}
}

// Delete

// Delete will remove the key and its associated vale.
func (m *FilesMap) Delete(key string) {
	m.nDirectories--
	n := len(m.store[key])
	m.nFiles -= n
	delete(m.store, key)
	// find key in slice
	idx := -1
	for i, val := range m.keys {
		if val == key {
			idx = i
			break
		}
	}
	if idx != -1 {
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
	}
}

// Iterator is used to loop through the stored key-value pairs.
// The returned anonymous function returns the index, key and value
//
// Example:
// 	m := NewFilesMap()
//
// 	m.Set("folder1", []File{file11, file12})
// 	m.Set("folder2", []File{file21, file22})
// 	m.Set("folder3", []File{file31, file32})
// 	m.Set("folder4", []File{file41, file42})
// 	m.Delete("folder3")
// 	m.Delete("folder8")
//
// 	iterator := m.Iterator()
//
// 	for {
// 		i, k, v := iterator()
// 		if i == nil {
// 			break
// 		}
// 		fmt.Println(*i, *k, v)
// 	}
func (m *FilesMap) Iterator() func() (*int, *string, []File) {
	var keys = m.keys

	j := 0

	return func() (_ *int, _ *string, _ []File) {
		if j > len(keys)-1 {
			return
		}

		row := keys[j]
		j++

		return &[]int{j - 1}[0], &row, m.store[row]

	}
}

// OrderedByFolder sort keys in increasing order (code number)
func (m *FilesMap) OrderedByFolder() {
	sort.Strings(m.keys)
}

// OrderedByFolderReverse sort keys in decreasing order (code number)
func (m *FilesMap) OrderedByFolderReverse() {
	sort.Sort(sort.Reverse(sort.StringSlice(m.keys)))
}

// OrderedAll sort keys and vals ([]File) in increasing order (code number)
func (m *FilesMap) OrderedAll() {
	m.OrderedByFolder()

	// byFolder := func(f1, f2 *File) bool {
	// 	return f1.Folder < f2.Folder
	// }
	byFileName := func(f1, f2 *File) bool {
		return f1.FileName < f2.FileName
	}
	// OrderedBy(byFolder, byFileName).Sort(m.Files)
	for _, files := range m.store {
		OrderedBy(byFileName).Sort(files)
	}
}

// OrderedAllReverse sort keys and vals ([]File) in decreasing order (code number)
func (m *FilesMap) OrderedAllReverse() {
	m.OrderedByFolder()

	// byFolder := func(f1, f2 *File) bool {
	// 	return f1.Folder > f2.Folder
	// }
	byFileName := func(f1, f2 *File) bool {
		return f1.FileName > f2.FileName
	}
	// OrderedBy(byFolder, byFileName).Sort(m.Files)
	for _, files := range m.store {
		OrderedBy(byFileName).Sort(files)
	}
}

// FindFiles will find files
// 	isRecursive:
// 		false to find files only in `root` directory
//		true  to find recursive files including subfolders
func (m *FilesMap) FindFiles(root string, isRecursive bool) {
	// m.FindFilesFunc(root, isRecursive, func(f File) bool {
	// 	return (len(f.FileName) == 0 || HasPrefix(f.FileName, ".") || REUsuallyExclude.MatchString(f.FullPath))
	// })
	m.FindFilesFunc(root, isRecursive, func(f File) bool {
		return false
	})
}

// FindRegularFiles will find files
// 	isRecursive:
// 		false to find files only in `root` directory
//		true  to find recursive files including subfolders
// 	use `len(f.FileName) == 0 || HasPrefix(f.FileName, ".") || REUsuallyExclude.MatchString(f.FullPath)` to exclude files
func (m *FilesMap) FindRegularFiles(root string, isRecursive bool) {
	m.FindFilesFunc(root, isRecursive, func(f File) bool {
		return (len(f.FileName) == 0 || HasPrefix(f.FileName, ".") || REUsuallyExclude.MatchString(f.FullPath))
	})
}

// FindFilesFunc will find files using codintion `exclude` func
// 	isRecursive:
// 		false to find files only in `root` directory
//		true  to find recursive files including subfolders
func (m *FilesMap) FindFilesFunc(
	root string, isRecursive bool, exclude func(file File) bool) {

	files, err := findFilesMapFunc(root, isRecursive, exclude)
	if err != nil {
		Logger.Error(err)
	}
	m.SetRoot(root)
	// nf := 0
	for _, f := range files {
		if strings.EqualFold(root, f.Folder) && !rootHasFile {
			rootHasFile = true
		}
		m.SetOne(f.ShortFolder, f)
		// nf++
	}
	// m.nFiles = nf
	m.nDirectories = len(m.keys)
	if _, ok := m.Get("./"); ok {
		m.nDirectories--
	}
}

// findFilesMapFunc will find files using codintion `exclude` func
// 	isRecursive:
// 		false to get []File in `folder`
// 		true  to get []File in `folder` and all `subfolders`
// 	exclude(file) return true to exclude
func findFilesMapFunc(folder string, isRecursive bool, exclude func(file File) bool) ([]File, error) {
	var (
		files []File
	)
	if isRecursive {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			file, err := filepath.Abs(path)
			if err != nil {
				return nil
			}
			if folder == path {
				return nil
			}
			if !info.IsDir() && info.Mode().IsRegular() {
				f := ConstructFile(file, folder)
				if !exclude(f) {
					files = append(files, f)
				}
			}
			// f := ConstructFile(file, folder)
			// if !exclude(f) {
			// 	files = append(files, f)
			// }
			return nil
		})
		// files = checkRootFiles(files, folder)
		return files, err
	}

	f, err := os.Open(folder)
	defer f.Close()

	if err != nil {
		return files, err
	}

	if fileinfo, err := f.Readdir(-1); err == nil {
		for _, file := range fileinfo {
			if !file.IsDir() && file.Mode().IsRegular() {
				folder, err := filepath.Abs(folder)
				if err != nil {
					return files, err
				}
				f := ConstructFile(folder+"\\"+file.Name(), folder)
				if !exclude(f) {
					files = append(files, f)
				}
			}
		}
	}

	// files = checkRootFiles(files, folder)
	return files, err
}

// Fprint filelist with `head`
func (m *FilesMap) Fprint(w io.Writer, mode OutputMode, head, pad string) {
	switch mode {
	case OTreeMode:
		m.FprintTree(w, head, pad)
	case OTableFormatMode:
		tf := &TableFormat{
			Fields:    []string{"No.", "Sorted Files"},
			LenFields: []int{5, 75},
			Aligns:    []Align{AlignRight, AlignLeft},
			Padding:   pad,
		}
		m.FprintTable(w, tf, head)
	default: // OPlainTextMode
		m.FprintText(w, head, pad)
	}
}

// Text print out FileList in plain text mode
func (m *FilesMap) Text(head, pad string) string {
	buf := new(bytes.Buffer)
	m.Fprint(buf, OPlainTextMode, head, pad)
	return string(buf.Bytes())
}

// FprintText print out FileList in plain text mode
func (m *FilesMap) FprintText(w io.Writer, head, pad string) {
	fmt.Fprintln(w, PaddingString(head, pad))
	fmt.Fprintln(w, pad)
	j := 0
	for _, k := range m.keys {
		for _, f := range m.store[k] {
			if strings.EqualFold(f.File, zeroRootFiles) {
				continue
			}
			j++
			fmt.Fprintf(w, "%s%4d %s\n", pad, j, f.FullPath)
		}
	}
	fmt.Fprintln(w, pad)
	fmt.Fprintf(w, "%s%d directories, %d Files\n", pad, m.NDirectories(), j)
	fmt.Fprintln(w, pad)
}

// Table will return a string of `FilesMap` in table mode with `head`
func (m *FilesMap) Table(head, pad string) string {
	buf := new(bytes.Buffer)
	m.Fprint(buf, OTableFormatMode, head, pad)
	return TrimFrontEndSpaceLine(string(buf.Bytes()))
}

// FprintTable print `FilesMap` in table mode with `head`
func (m *FilesMap) FprintTable(w io.Writer, tp *TableFormat, head string) {
	tp.Prepare(w)
	tp.SetBeforeMessage(head)
	tp.PrintSart()
	nSubFolders := m.NDirectories()
	nFiles := m.NFiles()
	j := 0
	for i, k := range m.keys {
		if strings.EqualFold(m.store[k][0].File, zeroRootFiles) {
			// nFiles--
			continue
		}
		if strings.EqualFold(k, "./") {
			tp.PrintRow("", fmt.Sprintf("[%d]. root (%q)", i, k))
		} else {
			tp.PrintRow("", fmt.Sprintf("[%d]. subfolder: %q", i, k))
		}

		for _, f := range m.store[k] {
			j++
			tp.PrintRow(j, f.File)
		}
		tp.PrintRow("", fmt.Sprintf("Sum: %d files.", j))
		j = 0
		if i == len(m.keys)-1 {
			break
		}
		tp.PrintMiddleSepLine()

	}

	tp.SetAfterMessage(fmt.Sprintf("%d directories, %d files\n", nSubFolders, nFiles))

	tp.PrintEnd()
}

func trimPath(path string) string {
	mpath := TrimPrefix(path, "./")
	mpath = TrimSuffix(mpath, "/")
	return mpath
}

// Tree will return a string of `FilesMap` in tree mode
func (m *FilesMap) Tree(head, pad string) string {
	buf := new(bytes.Buffer)
	m.FprintTree(buf, head, pad)
	return TrimFrontEndSpaceLine(string(buf.Bytes()))
}

// FprintTree print out `FilesMap` in tree mode
func (m *FilesMap) FprintTree(w io.Writer, head, pad string) {
	fmt.Fprintln(w, PaddingString(head, pad))
	fmt.Fprintln(w, pad)

	// root, rootPath := ".", m.root

	// nSubFolders := m.NDirectories()
	tree := treeprint.New()

	// Logger.Info("rootHasFile:", rootHasFile)
	nFiles, nDir := 0, 0
	for _, k := range m.keys {
		nDir++
		files, _ := m.Get(k)
		// nf := len(files)
		trimfd := trimPath(k)
		ss := strings.Split(trimfd, "/")
		// ns := len(ss)
		// if nf == 1 && files[0].IsDir() {
		// 	nFiles--
		// }
		fmt.Printf("nDir: %2d k: %q;\n\ttrimfd: %q\n", nDir, k, trimfd)
		spew.Dump(ss)
		for _, file := range files {
			if !file.IsDir() {
				nFiles++
			}
			fmt.Printf("\t%2d IsDir: %v %q\n", nFiles, file.IsDir(), file.File)
		}
	}
	// files := m.store[k]
	// 	if ns == 1 {
	// 		if len(ss[0]) == 0 {
	// 			if len(m.store[k]) == 1 && strings.EqualFold(m.store[k][0].File, zeroRootFiles) {
	// 				m.Delete(k)
	// 			}
	// 			tree.SetMetaValue(fmt.Sprintf("%d (%d directories, %d files)", len(m.store[k])-1, m.NDirectories(), m.NFiles()))
	// 			// tree.SetValue(root)
	// 			tree.SetValue(fmt.Sprintf("%s\n» root: %s", root, rootPath))
	// 			for _, v := range m.store[k] {
	// 				if isRootPath(v, rootPath) {
	// 					continue
	// 				}
	// 				tree.AddNode(v.File)
	// 			}
	// 		} else {
	// 			one := tree.AddMetaBranch(cast.ToString(len(m.store[k])-1), ss[0])
	// 			for _, v := range m.store[k] {
	// 				if isRootPath(v, rootPath) {
	// 					continue
	// 				}
	// 				one.AddNode(v.File)
	// 			}
	// 		}
	// 		continue
	// 	}

	// 	treend := make([]treeprint.Tree, ns)
	// 	treend[0] = tree.FindByValue(ss[0])
	// 	if treend[0] == nil {
	// 		// fmt.Println(tree.String())
	// 		// Logger.Debugf("ns = %d; ss[]: %v", ns, ss)
	// 		treend[0] = tree.AddMetaBranch(cast.ToString(len(m.store[k])-1), ss[0])
	// 	}
	// 	for i := 1; i < ns; i++ {
	// 		treend[i] = treend[0].FindByValue(ss[i])
	// 		if treend[i] == nil {
	// 			treend[i] = treend[i-1].AddMetaBranch(cast.ToString(len(m.store[k])-1), ss[i])
	// 		}
	// 		for _, v := range m.store[k] {
	// 			if isRootPath(v, rootPath) || v.IsDir() {
	// 				continue
	// 			}
	// 			treend[i].AddNode(v.File)
	// 		}
	// 	}
	// }

	fmt.Fprintln(w, PaddingString(tree.String(), pad))
	fmt.Fprintf(w, "%s%d directories, %d files\n", pad, m.NDirectories(), m.NFiles())
	keys := make([]string, len(m.keys))
	copy(keys, m.keys)
	sort.Strings(keys)
	spew.Dump(keys)
	// for _, k := range m.keys {
	// 	for _, f := range m.store[k] {
	// 		fmt.Printf("k:%q; f:%q\n", k, f.File)
	// 	}
	// }
}

// // FprintTree print out `FilesMap` in tree mode
// func (m *FilesMap) FprintTree(w io.Writer, head, pad string) {
// 	fmt.Fprintln(w, PaddingString(head, pad))
// 	fmt.Fprintln(w, pad)

// 	root, rootPath := ".", m.root

// 	// nSubFolders := m.NDirectories()
// 	// nFiles := m.NFiles()
// 	tree := treeprint.New()

// 	// Logger.Info("rootHasFile:", rootHasFile)
// 	for _, k := range m.keys {
// 		trimfd := trimPath(k)
// 		ss := strings.Split(trimfd, "/")
// 		ns := len(ss)
// 		// files := m.store[k]
// 		if ns == 1 {
// 			if len(ss[0]) == 0 {
// 				if len(m.store[k]) == 1 && strings.EqualFold(m.store[k][0].File, zeroRootFiles) {
// 					m.Delete(k)
// 				}
// 				tree.SetMetaValue(fmt.Sprintf("%d (%d directories, %d files)", len(m.store[k])-1, m.NDirectories(), m.NFiles()))
// 				// tree.SetValue(root)
// 				tree.SetValue(fmt.Sprintf("%s\n» root: %s", root, rootPath))
// 				for _, v := range m.store[k] {
// 					if isRootPath(v, rootPath) {
// 						continue
// 					}
// 					tree.AddNode(v.File)
// 				}
// 			} else {
// 				one := tree.AddMetaBranch(cast.ToString(len(m.store[k])-1), ss[0])
// 				for _, v := range m.store[k] {
// 					if isRootPath(v, rootPath) {
// 						continue
// 					}
// 					one.AddNode(v.File)
// 				}
// 			}
// 			continue
// 		}

// 		treend := make([]treeprint.Tree, ns)
// 		treend[0] = tree.FindByValue(ss[0])
// 		if treend[0] == nil {
// 			// fmt.Println(tree.String())
// 			// Logger.Debugf("ns = %d; ss[]: %v", ns, ss)
// 			treend[0] = tree.AddMetaBranch(cast.ToString(len(m.store[k])-1), ss[0])
// 		}
// 		for i := 1; i < ns; i++ {
// 			treend[i] = treend[0].FindByValue(ss[i])
// 			if treend[i] == nil {
// 				treend[i] = treend[i-1].AddMetaBranch(cast.ToString(len(m.store[k])-1), ss[i])
// 			}
// 			for _, v := range m.store[k] {
// 				if isRootPath(v, rootPath) || v.IsDir() {
// 					continue
// 				}
// 				treend[i].AddNode(v.File)
// 			}
// 		}
// 	}

// 	fmt.Fprintln(w, PaddingString(tree.String(), pad))
// 	fmt.Fprintf(w, "%s%d directories, %d files\n", pad, m.NDirectories(), m.NFiles())
// 	// spew.Dump(m.keys)
// 	// for _, k := range m.keys {
// 	// 	for _, f := range m.store[k] {
// 	// 		fmt.Printf("k:%q; f:%q\n", k, f.File)
// 	// 	}
// 	// }
// }

func isRootPath(file File, rootPath string) bool {
	return strings.EqualFold(file.FullPath, rootPath)
}
