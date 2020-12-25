package filetree

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/karrick/godirwalk"
	"github.com/spf13/cast"

	"github.com/shyang107/paw"
	// "github.com/shyang107/paw/treeprint"
)

// FileMap stores directory map to `map[{{ sub-path }}]{{ *File }}`
type FileMap map[string][]*File

// FileList stores the list information of File
type FileList struct {
	root      string   // root directory
	store     FileMap  // all files in `root` directory
	dirs      []string // keys of `store`
	depth     int
	totalSize uint64
	gitstatus GitStatus
}

// NewFileList will return the instance of `FileList`
func NewFileList(root string) *FileList {
	if len(root) == 0 {
		return &FileList{}
	}
	return &FileList{
		root:  root,
		store: make(map[string][]*File),
		dirs:  []string{},
	}
}

// String ...
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f FileList) String() string {
	f.DisableColor()
	str := f.ToTextString("")
	f.EnableColor()
	return str
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
	dirs := f.Dirs()
	fm := f.Map()
	for _, dir := range dirs {
		for _, file := range fm[dir] {
			if !file.IsDir() {
				nf++
			}
		}
	}
	return nf
}

// NSubDirsAndFiles will return the number of sub-dirs and sub-files in dir
func (f *FileList) NSubDirsAndFiles(dir string) (ndirs, nfiles int) {
	if _, ok := f.Map()[dir]; !ok {
		return ndirs, nfiles
	}
	return getNDirsFiles(f.Map()[dir])
}

// DirInfo will return the colorful string of sub-dir ( file.IsDir is true)
func (f *FileList) DirInfo(file *File) string {
	return getDirInfo(f, file)
}

// GetGitStatus will return git short status of `FileList`
func (f *FileList) GetGitStatus() GitStatus {
	return f.gitstatus
}

// GetHead4Meta will return a colorful string of head line for meta information of File
func (f *FileList) GetHead4Meta(pad, username, groupname string) string {
	return getColorizedHead(pad, username, groupname)
}

// AddFile will add file into the file list
func (f *FileList) AddFile(file *File) {
	if _, ok := f.store[file.Dir]; !ok {
		f.store[file.Dir] = []*File{}
		f.dirs = append(f.dirs, file.Dir)
		f.totalSize += file.Size
	}
	f.store[file.Dir] = append(f.store[file.Dir], file)
	f.totalSize += file.Size
	if file.IsDir() {
		pdir := findPreDir(file.Dir)
		if !paw.EqualFold(pdir, file.Dir) {
			f.store[pdir] = append(f.store[pdir], file)
		}
	}
}

func findPreDir(dir string) string {
	ddirs := paw.Split(dir, PathSeparator)
	if len(ddirs) == 1 {
		ddirs = []string{RootMark}
	}
	if ddirs[0] == ".." && len(ddirs) == 2 {
		ddirs[0] = RootMark
	}
	pdir := filepath.Join(ddirs[:len(ddirs)-1]...)
	// fmt.Println(dir, ddirs, pdir)
	return pdir

}

func (f *FileList) DisableColor() {
	SetNoColor()
}

func (f *FileList) EnableColor() {
	DefaultNoColor()
}

type FileSortByPathP []*File

func (a FileSortByPathP) Len() int           { return len(a) }
func (a FileSortByPathP) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FileSortByPathP) Less(i, j int) bool { return a[i].BaseName < a[j].BaseName }

func (f *FileList) Sort() {
	sort.Strings(f.dirs)

	for _, dir := range f.dirs {
		fm := FileSortByPathP(f.store[dir])
		for _, file := range fm {
			fmt.Println("Before:", file)
		}
		sort.Sort(fm)
		for _, file := range fm {
			fmt.Println("After:", file)
		}
		f.store[dir] = fm
	}
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
// value ErrSkipDir or ErrSkipFile. If the function returns ErrSkipDir when invoked on a directory,
// FindFiles skips the directory's contents entirely. If the function returns ErrSkipDir
// when invoked on a non-directory file, FindFiles skips the remaining files in the
// containing directory.
// If the returned error is SkipFile when inviked on a file, FindFiles will skip the file.
// TODO maybe has a better way
type IgnoreFn func(f *File, err error) error

// DefaultIgnoreFn is default IgnoreFn using in FindFiles
//
// 	Skip file: prefix "." of files
// 	Skip folder: prefix "." of directory
var DefaultIgnoreFn = func(f *File, err error) error {
	if err != nil {
		return err
	}
	if f.IsDir() && paw.HasPrefix(f.BaseName, ".") {
		return SkipDir
	}
	if paw.HasPrefix(f.BaseName, ".") {
		return SkipFile
	}
	return nil
}

// FindFiles will find files using codintion `ignore` func
// 	depth : depth of subfolders
// 		depth < 0 : walk through all directories of {root directory}
// 		depth == 0 : {root directory}/
// 		depth > 0 : {root directory}/{level 1 directory}/.../{{ level n directory }}/
// 	ignore: IgnoreFn func(f *File, err error) error
// 		ignoring condition of files or directory
// 		ignore == nil, using DefaultIgnoreFn
func (f *FileList) FindFiles(depth int, ignore IgnoreFn) error {
	if ignore == nil {
		ignore = DefaultIgnoreFn
	}
	f.gitstatus, _ = GetShortStatus(f.Root())
	f.depth = depth
	root := f.Root()
	switch {
	case depth == 0: //{root directory}/*
		scratchBuffer := make([]byte, godirwalk.MinimumScratchBufferSize)
		files, err := godirwalk.ReadDirnames(root, scratchBuffer)
		if err != nil {
			return errors.New(root + ": " + err.Error())
		}

		sort.Slice(files, func(i, j int) bool {
			return paw.ToLower(files[i]) < paw.ToLower(files[j])
		})
		file := NewFileRelTo(root, root)
		f.AddFile(file)
		for _, name := range files {
			file := NewFileRelTo(root+PathSeparator+name, root)
			err := ignore(file, nil)
			if err == SkipFile || err == SkipDir {
				continue
			}
			f.AddFile(file)
		}
		// f.Sort()
	default: //walk through all directories of {root directory}
		err := godirwalk.Walk(root, &godirwalk.Options{
			Callback: func(path string, de *godirwalk.Dirent) error {
				file := NewFileRelTo(path, root)
				idepth := len(file.DirSlice()) - 1
				if depth > 0 {
					if idepth > depth {
						return godirwalk.SkipThis
					}
				}
				if err1 := ignore(file, nil); err1 == SkipFile || err1 == SkipDir {
					return godirwalk.SkipThis
				}
				f.AddFile(file)
				return nil
			},
			ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
				// fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
				paw.Logger.Errorf("ERROR: %s\n", err)

				// For the purposes of this example, a simple SkipNode will suffice, although in reality perhaps additional logic might be called for.
				return godirwalk.SkipNode
			},
			Unsorted: true, // set true for faster yet non-deterministic enumeration (see godoc)
		})
		if err != nil {
			return errors.New(root + ": " + err.Error())
		}
		// sort
		sort.Slice(f.dirs, func(i, j int) bool {
			return paw.ToLower(f.dirs[i]) < paw.ToLower(f.dirs[j])
		})
		for _, dir := range f.dirs {
			sort.Slice(f.store[dir], func(i, j int) bool {
				return paw.ToLower(f.store[dir][i].Path) < paw.ToLower(f.store[dir][j].Path)
			})
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
	return toListTree(f, pad, false)
}

// // ToTree will return the []byte of FileList in tree form
// func (f *FileList) ToTree(pad string) []byte {

// 	tree := treeprint.New()

// 	dirs := f.Dirs()
// 	// nd := len(dirs) // including root
// 	ntf := 0
// 	var one, pre treeprint.Tree
// 	fm := f.Map()

// 	for i, dir := range dirs {
// 		files := f.Map()[dir]
// 		ndirs, nfiles := getNDirsFiles(files) // excluding the dir
// 		ntf += nfiles
// 		for jj, file := range files {
// 			// fsize := file.Size
// 			// sfsize := ByteSize(fsize)
// 			if jj == 0 && file.IsDir() {
// 				if i == 0 { // root dir
// 					// tree.SetValue(fmt.Sprintf("%v (%v)", file.LSColorString(file.Dir), file.LSColorString(file.Path)))
// 					tree.SetValue(getName(file))
// 					tree.SetMetaValue(KindLSColorString("di", fmt.Sprintf("%d dirs", ndirs)+", "+KindLSColorString("fi", fmt.Sprintf("%d files", nfiles))))
// 					one = tree
// 				} else {
// 					pre = preTree(dir, fm, tree)
// 					if f.depth != 0 {
// 						// one = pre.AddMetaBranch(nf-1, file)
// 						one = pre.AddMetaBranch(KindLSColorString("di", fmt.Sprintf("%d dirs", ndirs)+", "+KindLSColorString("fi", fmt.Sprintf("%d files", nfiles))), file)
// 					} else {
// 						one = pre.AddBranch(file)
// 					}
// 				}
// 				continue
// 			}
// 			// add file node
// 			link := checkAndGetColorLink(file)
// 			if !file.IsDir() {
// 				if len(link) > 0 {
// 					one.AddMetaNode(link, file)
// 				} else {
// 					one.AddNode(file)
// 				}
// 			}
// 		}
// 	}
// 	buf := new(bytes.Buffer)
// 	buf.Write(tree.Bytes())

// 	printTotalSummary(buf, "", f.NDirs(), f.NFiles(), f.totalSize)

// 	return paddingTree(pad, buf.Bytes())
// }

// ToTableString will return the string of FileList in table form
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f *FileList) ToTableString(pad string) string {
	return string(f.ToTable(pad))
}

// ToTable will return the []byte of FileList in table form
// 	`size` of directory shown in the returned value, is accumulated size of sub contents
func (f *FileList) ToTable(pad string) []byte {

	var (
		buf    = new(bytes.Buffer)
		nDirs  = f.NDirs()
		nFiles = f.NFiles()
		dirs   = f.Dirs()
		fm     = f.Map()
	)

	f.DisableColor()

	tf := &paw.TableFormat{
		Fields:    []string{"No.", "Mode", "Size", "Files"},
		LenFields: []int{5, 10, 6, 80},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignRight, paw.AlignRight, paw.AlignLeft},
		Padding:   pad,
	}
	tf.Prepare(buf)

	sdsize := ByteSize(f.totalSize)
	head := fmt.Sprintf("Root directory: %v, size: %v", f.Root(), sdsize)
	tf.SetBeforeMessage(head)

	tf.PrintSart()
	j := 0
	for i, dir := range dirs {
		sumsize := uint64(0)
		ndirs, nfiles := 0, 0
		for jj, file := range fm[dir] {
			fsize := file.Size
			sfsize := ByteSize(fsize)
			mode := file.Stat.Mode()
			if jj == 0 && file.IsDir() {
				idx := fmt.Sprintf("D%d", i)
				sfsize = "-"
				switch f.depth {
				case 0:
					if len(fm[dir]) > 1 && !paw.EqualFold(file.Dir, RootMark) {
						tf.PrintRow(idx, mode, sfsize, file)
					}
				default:
					if paw.EqualFold(file.Dir, RootMark) {
						tf.PrintRow(idx, mode, sfsize, file.ColorDirName(""))
					} else {
						tf.PrintRow(idx, mode, sfsize, file.ColorDirName(f.Root()))
					}
				}
				continue
			}
			jdx := fmt.Sprintf("d%d", ndirs+1)
			name := file.ColorBaseName()
			if !file.IsDir() {
				sumsize += fsize
				j++
				nfiles++
				tf.PrintRow(j, mode, sfsize, name)
			} else {
				ndirs++
				tf.PrintRow(jdx, mode, "-", name)
			}

		}
		if f.depth != 0 {
			// printDirSummary(buf, pad, ndirs, nfiles, sumsize)
			tf.PrintRow("", "", "", fmt.Sprintf("%v directories, %v files, size: %v.", ndirs, nfiles, ByteSize(sumsize)))

			if i != len(dirs)-1 {
				tf.PrintMiddleSepLine()
			}
		}
	}

	tf.SetAfterMessage(fmt.Sprintf("\n%v directories, %v files, total %v.", nDirs, nFiles, ByteSize(f.totalSize)))
	tf.PrintEnd()

	f.EnableColor()

	return buf.Bytes()
}

// ToTextString will return the string of FileList in table form
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f *FileList) ToTextString(pad string) string {
	return string(f.ToText(pad))
}

// ToText will return the []byte of FileList in table form
// 	`size` of directory shown in the returned value, is accumulated size of sub contents
func (f *FileList) ToText(pad string) []byte {
	var (
		w     = new(bytes.Buffer)
		dirs  = f.Dirs()
		fm    = f.Map()
		width = 80
	)

	sdsize := ByteSize(f.totalSize)
	fmt.Fprintf(w, "%sRoot directory: %v, size: %v\n", pad, getDirName(f.Root(), ""), KindLSColorString("di", sdsize))
	fmt.Fprintf(w, "%s%s\n", pad, paw.Repeat("=", width))

	ppad := ""

	i1 := len(cast.ToString(f.NDirs()))
	j1 := intmax(i1, len(cast.ToString(f.NFiles())))
	j := 0
	for i, dir := range dirs {
		istr := KindLSColorString("di", fmt.Sprintf("%[2]*[1]d.", i, i1))
		sumsize := uint64(0)
		nfiles := 0
		ndirs := 0
		files := fm[dir]
		for jj, file := range files {
			cperm := file.ColorPermission()
			fsize := file.Size
			cfsize := file.ColorSize()
			if file.IsDir() {
				cfsize = NewEXAColor("-").Sprint(fmt.Sprintf("%6s", "-"))
			}
			if jj == 0 && file.IsDir() {
				switch f.depth {
				case 0:
					if len(files) > 1 && !paw.EqualFold(file.Dir, RootMark) {
						ppad = paw.Repeat("    ", len(file.DirSlice())-1)
						fmt.Fprintf(w, "%s%v %10v %v %v\n", pad+ppad, istr, cperm, cfsize, file.ColorDirName(f.Root()))
					}
				default: // f.depth !=0
					if paw.EqualFold(file.Dir, RootMark) {
						printFileItem(w, pad, istr, cperm, file.ColorDirName(""))
					} else {
						ppad = paw.Repeat("    ", len(file.DirSlice())-1)
						printFileItem(w, pad+ppad, istr, cperm, file.ColorDirName(f.Root()))
					}
				}
				continue
			}
			sumsize += fsize
			if f.depth != 0 {
				j1 = len(cast.ToString(len(fm[dir]) - 1))
			}
			jstr := ""
			if !file.IsDir() {
				j++
				nfiles++
				jstr = fmt.Sprintf("%[2]*[1]d.", j, j1)
			} else {
				ndirs++
				jstr = KindLSColorString("di", fmt.Sprintf("%[2]*[1]d.", ndirs, j1))
			}
			name := file.ColorBaseName()
			printFileItem(w, pad+ppad+"    ", jstr, cperm, cfsize, name)
		}
		if f.depth != 0 {
			printDirSummary(w, pad+ppad+"    ", ndirs, nfiles, sumsize)

			if i != len(dirs)-1 {
				fmt.Fprintf(w, "%s%s\n", pad, paw.Repeat("-", width))
			}
		}
	}

	fmt.Fprintf(w, "%s%s\n", pad, paw.Repeat("=", width))
	printTotalSummary(w, pad, f.NDirs(), f.NFiles(), f.totalSize)

	return w.Bytes()
}

// ToList will return the string of FileList in list form (like as `exa`)
func (f *FileList) ToListString(pad string) string {
	return string(f.ToList(pad))
}

// ToList will return the []byte of FileList in list form (like as `exa`)
func (f *FileList) ToList(pad string) []byte {
	var (
		w     = new(bytes.Buffer)
		dirs  = f.Dirs()
		fm    = f.Map()
		chead = f.GetHead4Meta(pad, urname, gpname)
	)

	ctdsize := ByteSize(f.totalSize)
	head := fmt.Sprintf("%sRoot directory: %v, size: %v", pad, getDirName(f.Root(), ""), KindLSColorString("di", ctdsize))
	fmt.Fprintln(w, head)
	fmt.Fprintf(w, "%s%s\n", pad, paw.Repeat("=", 80))

	fmt.Fprintln(w, chead)

	git := f.GetGitStatus()
	for i, dir := range dirs {
		sumsize := uint64(0)
		nfiles := 0
		ndirs := 0
		for jj, file := range fm[dir] {
			if file.IsDir() {
				if jj == 0 {
					if !paw.EqualFold(file.Dir, RootMark) {
						if f.depth != 0 {
							fmt.Fprintf(w, "%s%v\n", pad, file.ColorDirName(f.Root()))
							fmt.Fprintln(w, chead)
						}
					}
					continue
				}
				ndirs++
			} else {
				nfiles++
			}
			sumsize += file.Size
			name := file.ColorBaseName()
			fmt.Fprintf(w, "%s%s%s\n", pad, file.ColorMeta(git), name)
		}

		if f.depth != 0 {
			printDirSummary(w, pad, ndirs, nfiles, sumsize)
			if i < len(dirs)-1 {
				fmt.Fprintf(w, "%s%s\n", pad, paw.Repeat("-", 80))
			}
		}
	}

	fmt.Fprintf(w, "%s%s\n", pad, paw.Repeat("=", 80))
	printTotalSummary(w, pad, f.NDirs(), f.NFiles(), f.totalSize)

	return w.Bytes()
}

// ToListTreeString will return the string of `ToListTree(pad)` in list+tree form (like as `exa -T(--tree)`)
func (f *FileList) ToListTreeString(pad string) string {
	return string(f.ToListTree(pad))
}

// ToListTree will return the []byte of FileList in list+tree form (like as `exa -T(--tree)`)
func (f *FileList) ToListTree(pad string) []byte {
	return toListTree(f, pad, true)
}

func toListTree(f *FileList, pad string, isMeta bool) []byte {
	var (
		buf = new(bytes.Buffer)
		fm  = f.Map()
	)

	// print head
	if isMeta {
		chead := f.GetHead4Meta(pad, urname, gpname)
		buf.WriteString(chead)
		buf.WriteByte('\n')
	}

	files := fm[RootMark]
	file := files[0]
	nfiles := len(files)

	// print root file
	meta := pad
	if isMeta {
		meta += file.ColorMeta(f.GetGitStatus())
	}

	name := file.ColorBaseName()

	buf.WriteString(fmt.Sprintf("%v%v", meta, name))
	buf.WriteByte('\n')

	// print files in the root dir
	git := f.GetGitStatus()
	level := 0
	var levelsEnded []int
	for i := 1; i < nfiles; i++ {
		file = files[i]
		edge := EdgeTypeMid
		if i == nfiles-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}

		printLTFile(buf, level, levelsEnded, edge, f, file, git, pad, isMeta)

		if file.IsDir() && len(fm[file.Dir]) > 1 {
			printLTDir(buf, level+1, levelsEnded, edge, f, file, git, pad, isMeta)
		}
	}

	// print end message
	printTotalSummary(buf, pad, f.NDirs(), f.NFiles(), f.totalSize)

	return buf.Bytes()
}
