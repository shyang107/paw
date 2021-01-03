package filetree

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"

	"github.com/thoas/go-funk"

	"github.com/mattn/go-runewidth"

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
	buf       *bytes.Buffer
	writer    io.Writer
	// writers   []io.Writer
	IsSort    bool // increasing order of Lower(path)
	filesBy   FilesBy
	dirsBy    DirsBy
	IsGrouped bool // grouping files and directories separetly
}

// NewFileList will return the instance of `FileList`
func NewFileList(root string) *FileList {
	if len(root) == 0 {
		return &FileList{}
	}
	fl := &FileList{
		root:      root,
		store:     make(map[string][]*File),
		dirs:      []string{},
		buf:       new(bytes.Buffer),
		IsSort:    true,
		filesBy:   nil,
		dirsBy:    nil,
		IsGrouped: false,
	}
	fl.SetWriters(fl.buf)
	return fl
}

// String ...
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f FileList) String() string {
	oldwr := f.writer
	f.writer = new(bytes.Buffer)
	f.DisableColor()
	str := f.ToLevelViewString("")
	f.EnableColor()
	f.writer = oldwr
	return str
}

// SetWriters will set writers... to writer of FileList
func (f *FileList) SetWriters(writers ...io.Writer) {
	// f.ResetBuffer()
	// f.writers = append(f.writers, writers...)
	// f.writers = writers
	f.writer = io.MultiWriter(writers...)
}

// func (f *FileList) checkWriter() {
// 	if f.writer == nil {
// 		f.ResetBuffer()
// 		f.SetWriters(f.Buffer())
// 	}
// }

// ResetWriters will reset default writers... (Buffer of FileList) to writer of FileList
func (f *FileList) ResetWriters() {
	f.ResetBuffer()
	// f.writers = []io.Writer{}
	f.SetWriters(f.buf)
}

// ResetBuffer will reset the buffer of FileList
func (f *FileList) ResetBuffer() {
	f.buf.Reset()
}

// Buffer will return the field buf of FileList
func (f *FileList) Buffer() *bytes.Buffer {
	return f.buf
}

// Writer will return the field writer of FileList
func (f *FileList) Writer() io.Writer {
	return f.writer
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
	dirs := f.dirs
	fm := f.store
	nf := 0
	for _, dir := range dirs {
		for _, file := range fm[dir][1:] {
			if file.IsFile() {
				nf++
			}
		}
	}
	return nf //- f.NDirs()
}

// NSubDirsAndFiles will return the number of sub-dirs and sub-files in dir
func (f *FileList) NSubDirsAndFiles(dir string) (ndirs, nfiles int) {
	if _, ok := f.store[dir]; !ok {
		return ndirs, nfiles
	}
	return getNDirsFiles(f.store[dir])
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
func (f *FileList) GetHead4Meta(pad, username, groupname string, git GitStatus) string {
	return getColorizedHead(pad, username, groupname, git)
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

// SkipThis is used as a return value indicate that the regular path
// (file or directory) named in the Callback is to be skipped.
// It is not returned as an error by any function.
var SkipThis = errors.New("skip the path")

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
type IgnoreFunc func(f *File, err error) error

// DefaultIgnoreFn is default IgnoreFn using in FindFiles
//
// _, file := filepath.Split(f.Path)
// 	Skip file: prefix "." of file
var DefaultIgnoreFn = func(f *File, err error) error {
	if err != nil {
		return err
	}
	_, file := filepath.Split(f.Path)
	if paw.HasPrefix(file, ".") {
		return SkipThis
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
func (f *FileList) FindFiles(depth int, ignore IgnoreFunc) error {
	if ignore == nil {
		ignore = DefaultIgnoreFn
	}
	f.gitstatus, _ = GetShortStatus(f.root)
	f.depth = depth
	root := f.root
	switch depth {
	case 0: //{root directory}/*
		// scratchBuffer := make([]byte, godirwalk.MinimumScratchBufferSize)
		files, err := godirwalk.ReadDirnames(root, nil)
		if err != nil {
			return errors.New(root + ": " + err.Error())
		}
		if f.IsSort {
			sort.Sort(ByLowerString(files))
		}

		file, err := NewFileRelTo(root, root)
		if err != nil {
			return err
		}
		f.AddFile(file)

		for _, name := range files {
			file, err := NewFileRelTo(root+PathSeparator+name, root)
			if err != nil {
				return err
			}
			if err := ignore(file, nil); err == SkipThis {
				continue
			}

			f.AddFile(file)
		}
	default: //walk through all directories of {root directory}
		err := godirwalk.Walk(root, &godirwalk.Options{
			Callback: func(path string, de *godirwalk.Dirent) error {
				file, err := NewFileRelTo(path, root)
				if err != nil {
					return err
				}
				idepth := len(file.DirSlice()) - 1
				if depth > 0 {
					if idepth > depth {
						return godirwalk.SkipThis
					}
				}
				if err1 := ignore(file, nil); err1 == SkipThis {
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
	}
	if f.IsSort {
		f.Sort()
	}
	return nil
}

var (
	DefaultFilesBy FilesBy = func(fi *File, fj *File) bool {

		if fi.IsDir() && fj.IsFile() {
			return true
		} else if fi.IsFile() && fj.IsDir() {
			return false
		}
		return paw.ToLower(fi.Path) < paw.ToLower(fj.Path)
	}
	DefaultDirsBy DirsBy = func(di string, dj string) bool {
		return paw.ToLower(di) < paw.ToLower(dj)
	}
)

// SetFilesSorter will set sorter of Files of FileList
func (f *FileList) SetFilesSorter(by FilesBy) {
	f.filesBy = by
}

// SetDirsSorter will set sorter of Dirs of FileList
func (f *FileList) SetDirsSorter(by DirsBy) {
	f.dirsBy = by
}

// Sort will sort FileList by sorter of dirsBy and filesBy.
//
// Default:
// 	Dirs: ToLower(a[i]) < ToLower(a[j])
// 	Map[dir][]*file: ToLower(a[i].Path) < ToLower(a[j].Path)
func (f *FileList) Sort() {
	f.SortBy(f.dirsBy, f.filesBy)
}

func (f *FileList) SortBy(dirsBy DirsBy, filesBy FilesBy) {
	f.SetDirsSorter(dirsBy)
	f.SetFilesSorter(filesBy)
	if dirsBy == nil {
		f.SetDirsSorter(DefaultDirsBy)
	}
	if filesBy == nil {
		f.SetFilesSorter(DefaultFilesBy)
	}
	f.dirsBy.Sort(f.dirs)
	for _, dir := range f.dirs {
		if len(f.store[dir]) > 1 {
			if !f.IsGrouped {
				f.filesBy.Sort(f.store[dir][1:])
			} else {
				sfiles := []*File{}
				sdirs := []*File{}
				for _, v := range f.store[dir][1:] {
					if v.IsDir() {
						sdirs = append(sdirs, v)
					} else {
						sfiles = append(sfiles, v)
					}
				}
				f.filesBy.Sort(sdirs)
				f.filesBy.Sort(sfiles)
				copy(f.store[dir][1:], sdirs)
				copy(f.store[dir][len(sdirs)+1:], sfiles)
			}
		}
	}
}

// ToTreeViewString will return the string of FileList in tree form
func (f *FileList) ToTreeViewString(pad string) string {
	return string(f.ToTreeView(pad))
}

// ToTreeView will return the []byte of FileList in tree form
func (f *FileList) ToTreeView(pad string) []byte {
	pdview = PTreeView
	return toListTreeView(f, pad)
}

// ToTableViewString will return the string of FileList in table form
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f *FileList) ToTableViewString(pad string) string {
	return string(f.ToTableView(pad))
}

// ToTableView will return the []byte of FileList in table form
// 	`size` of directory shown in the returned value, is accumulated size of sub contents
func (f *FileList) ToTableView(pad string) []byte {

	var (
		// w      = new(bytes.Buffer)
		buf    = f.buf    //f.Buffer()
		w      = f.writer //f.Writer()
		nDirs  = f.NDirs()
		nFiles = f.NFiles()
		dirs   = f.dirs  //f.Dirs()
		fm     = f.store //f.Map()
	)
	buf.Reset()

	f.DisableColor()

	tf := &paw.TableFormat{
		Fields:    []string{"No.", "Mode", "Size", "Files"},
		LenFields: []int{5, 10, 6, 80},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignRight, paw.AlignRight, paw.AlignLeft},
		Padding:   pad,
	}
	tf.Prepare(w)

	sdsize := ByteSize(f.totalSize)
	head := fmt.Sprintf("Root directory: %v, size ≈ %v", f.root, sdsize)
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
						tf.PrintRow(idx, mode, sfsize, file.ColorDirName(f.root))
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

// ToLevelViewString will return the string of FileList in table form
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f *FileList) ToLevelViewString(pad string) string {
	return string(f.ToLevelView(pad))
}

// ToLevelView will return the []byte of FileList in table form
// 	`size` of directory shown in the returned value, is accumulated size of sub contents
func (f *FileList) ToLevelView(pad string) []byte {
	var (
		// w     = new(bytes.Buffer)
		buf   = f.Buffer()
		w     = f.Writer()
		dirs  = f.Dirs()
		fm    = f.Map()
		width = 80
	)
	buf.Reset()

	sdsize := ByteSize(f.totalSize)
	fmt.Fprintf(w, "%sRoot directory: %v, size ≈ %v\n", pad, getDirName(f.root, ""), KindLSColorString("di", sdsize))
	fmt.Fprintf(w, "%s%s\n", pad, paw.Repeat("=", width))

	i1 := len(cast.ToString(f.NDirs()))
	j1 := max(i1, len(cast.ToString(f.NFiles())))
	j := 0
	for i, dir := range dirs {
		ppad := pad
		istr := KindLSColorString("di", fmt.Sprintf("%[2]*[1]d.", i, i1))
		sumsize := uint64(0)
		nfiles := 0
		ndirs := 0
		for jj, file := range fm[dir] {
			cperm := file.ColorPermission()
			fsize := file.Size
			cfsize := file.ColorSize()
			if file.IsDir() {
				cfsize = NewEXAColor("-").Sprint(fmt.Sprintf("%6s", "-"))
				if jj == 0 {
					ppad += paw.Repeat("    ", len(file.DirSlice())-1)
					name := file.ColorDirName(f.root)
					if f.depth != 0 {
						if paw.EqualFold(file.Dir, RootMark) {
							ppad = pad
							name = file.ColorDirName("")
						}
						printFileItem(w, ppad, istr, cperm, name)
					}
					continue
				}
			}
			if f.depth != 0 {
				j1 = len(cast.ToString(len(fm[dir]) - 1))
			}
			jstr := ""
			if !file.IsDir() {
				sumsize += fsize
				j++
				nfiles++
				jstr = fmt.Sprintf("%[2]*[1]d.", j, j1)
			} else {
				ndirs++
				jstr = KindLSColorString("di", fmt.Sprintf("%[2]*[1]d.", ndirs, j1))
			}
			name := file.ColorBaseName()
			printFileItem(w, ppad+"    ", jstr, cperm, cfsize, name)
		}
		if f.depth != 0 {
			printDirSummary(w, ppad+"    ", ndirs, nfiles, sumsize)

			if i != len(dirs)-1 {
				printBanner(w, pad, "-", width)
			}
		}
	}

	printBanner(w, pad, "=", width)
	printTotalSummary(w, pad, f.NDirs(), f.NFiles(), f.totalSize)

	return buf.Bytes()
}

// ToListView will return the string of FileList in list form (like as `exa --header --long --time-style=iso --group --git`)
func (f *FileList) ToListViewString(pad string) string {
	return string(f.ToListView(pad))
}

// ToListView will return the []byte of FileList in list form (like as `exa --header --long --time-style=iso --group --git`)
func (f *FileList) ToListView(pad string) []byte {
	return toListView(f, pad, false)
}

// ToListExtendView will return the string of FileList in extend list form (like as `exa --header --long --time-style=iso --group --git -@`)
func (f *FileList) ToListExtendViewString(pad string) string {
	return string(f.ToListExtendView(pad))
}

// ToListExtendView will return the []byte of FileList in extend list form (like as `exa --header --long --time-style=iso --group --git --@`)
func (f *FileList) ToListExtendView(pad string) []byte {
	return toListView(f, pad, true)
}

// toListView will return the []byte of FileList in list form (like as `exa --header --long --time-style=iso --group --git`)
func toListView(f *FileList, pad string, isExtended bool) []byte {
	var (
		// w     = new(bytes.Buffer)
		buf                     = f.buf
		w                       = f.writer
		dirs                    = f.dirs
		fm                      = f.store
		git                     = f.GetGitStatus()
		chead                   = f.GetHead4Meta(pad, urname, gpname, git)
		ntdirs, nsdirs, ntfiles = 1, 0, 0
		fNDirs                  = f.NDirs()
		fNFiles                 = f.NFiles()
	)
	buf.Reset()

	ctdsize := ByteSize(f.totalSize)

	head := fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, getDirName(f.root, ""), KindLSColorString("di", ctdsize))
	fmt.Fprintln(w, head)

	if f.NDirs() == 0 && f.NFiles() == 0 {
		goto END
	}

	printBanner(w, pad, "=", 80)

	if funk.IndexOfString(f.dirs, RootMark) != -1 {
		fmt.Fprintln(w, chead)
	}
	for _, dir := range dirs {
		sumsize := uint64(0)
		nfiles := 0
		ndirs := 0
		sntd := ""
		if len(fm[dir]) > 1 {
			if !paw.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					// sntd = KindEXAColorString("dir", fmt.Sprintf("D%d:", ntdirs))
					fmt.Fprintf(w, "%s%s%v\n", pad, sntd, GetColorizedDirName(dir, f.root))
					if len(fm[dir]) > 1 {
						ntdirs++
						fmt.Fprintln(w, chead)
					}
				}
			}
		}
		// fmt.Fprintln(w, chead)
		for _, file := range fm[dir][1:] {
			sntf := ""
			if file.IsDir() {
				ndirs++
				nsdirs++
				// sntf = file.LSColorString(fmt.Sprintf("D%d(%d):", ndirs, nsdirs))
			}
			if file.IsFile() {
				nfiles++
				ntfiles++
				sumsize += file.Size
				// sntf = file.LSColorString(fmt.Sprintf("F%d(%d):", nfiles, ntfiles))
			}
			meta, metalength := file.ColorMeta(git)
			fmt.Fprintf(w, "%s%s%s%s\n", pad, sntf, meta, file.ColorBaseName())

			if isExtended {
				nx := len(file.XAttributes)
				sp := paw.Repeat(" ", metalength+len(sntf))
				if nx > 0 {
					edge := EdgeTypeMid
					for i := 0; i < nx; i++ {
						if i == nx-1 {
							edge = EdgeTypeEnd
						}
						fmt.Fprintf(w, "%s%s%s %s\n", pad, sp, NewEXAColor("-").Sprint(edge), file.XAttributes[i])
					}
				}
			}
		}

		if f.depth != 0 {
			if len(fm[dir]) > 1 {
				printDirSummary(w, pad, ndirs, nfiles, sumsize)
				// if i < len(dirs)-1 && nsdirs < fNDirs && ntfiles < fNFiles {
				// 	printBanner(w, pad, "-", 80)
				// }
				switch {
				case fNFiles == 0:
					if nsdirs < fNDirs {
						printBanner(w, pad, "-", 80)
					}
				default:
					if nsdirs < fNDirs && ntfiles < fNFiles {
						printBanner(w, pad, "-", 80)
					}
				}
			}
		}
	}
	printBanner(w, pad, "=", 80)
END:
	// printTotalSummary(w, pad, f.NDirs(), f.NFiles(), f.totalSize)
	printTotalSummary(w, pad, fNDirs, fNFiles, f.totalSize)

	// spew.Dump(dirs)
	return buf.Bytes()
}

// ToListTreeViewString will return the string of `ToListViewTree(pad)` in list+tree form (like as `exa -T(--tree)`)
func (f *FileList) ToListTreeViewString(pad string) string {
	return string(f.ToListTreeView(pad))
}

// ToListTreeView will return the []byte of FileList in list+tree form (like as `exa -T(--tree)`)
func (f *FileList) ToListTreeView(pad string) []byte {
	pdview = PListTreeView
	return toListTreeView(f, pad)
}

func toListTreeView(f *FileList, pad string) []byte {
	var (
		buf = f.Buffer()
		// w  = new(bytes.Buffer)
		w   = f.Writer()
		fm  = f.store
		git = f.GetGitStatus()
	)

	buf.Reset()

	files := fm[RootMark]
	file := files[0]
	nfiles := len(files)

	// print root file
	meta := pad
	switch pdview {
	case PListTreeView:
		chead := f.GetHead4Meta(pad, urname, gpname, git)
		fmt.Fprintf(w, "%v\n", chead)
		tmeta, _ := file.ColorMeta(f.GetGitStatus())
		meta += tmeta
	case PTreeView:
		meta += f.DirInfo(file) + " "
	}

	name := fmt.Sprintf("%v (%v)", file.LSColorString("."), file.ColorBaseName())
	fmt.Fprintf(w, "%v%v\n", meta, name)

	// print files in the root dir
	level := 0
	var levelsEnded []int
	for i := 1; i < nfiles; i++ {
		file = files[i]
		edge := EdgeTypeMid
		if i == nfiles-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}

		printLTFile(w, level, levelsEnded, edge, f, file, git, pad)

		if file.IsDir() && len(fm[file.Dir]) > 1 {
			printLTDir(w, level+1, levelsEnded, edge, f, file, git, pad)
		}
	}

	// print end message
	printTotalSummary(w, pad, f.NDirs(), f.NFiles(), f.totalSize)

	return buf.Bytes()
}

// ToClassifyView will return the string of FileList to display type indicator by file names (like as `exa -F` or `exa --classify`)
func (f *FileList) ToClassifyViewString(pad string) string {
	return string(f.ToClassifyView(pad))
}

// ToClassifyView will return the []byte of FileList to display type indicator by file names (like as `exa -F` or `exa --classify`)
func (f *FileList) ToClassifyView(pad string) []byte {
	var (
		buf              = f.buf
		w                = f.writer
		dirs             = f.dirs
		fm               = f.store
		_, terminalWidth = getTerminalSize()
	)
	buf.Reset()

	for i, dir := range dirs {
		if f.depth != 0 {
			fmt.Fprintf(w, "%v\n", fm[dir][0].ColorDirName(f.root))
			if len(fm[dir]) == 1 {
				fmt.Fprintln(w)
				continue
			}
		}

		files := fm[dir][1:]
		lens, sumlen := getleng(files)
		if sumlen <= terminalWidth {
			classifyPrintFiles(w, files)
		} else {
			classifyGridPrintFiles(w, files, lens, sumlen, terminalWidth)
		}

		if f.depth == 0 {
			break
		} else {
			// printDirSummary(w,ndirs, nfiles, sumsize)
			if i < len(dirs)-1 {
				fmt.Fprintln(w)
			}
		}
	}

	printTotalSummary(w, "", f.NDirs(), f.NFiles(), f.totalSize)

	b := padding(pad, buf.Bytes())

	return b
}

func classifyGridPrintFiles(w io.Writer, files []*File, lens []int, sumlen int, twidth int) {

	nFields := calNFields(lens, twidth)
	widths := calWidth(lens, nFields, twidth)
	nFields = len(widths)

	nfolds := len(lens) / len(widths)
	for i := 0; i < nfolds; i++ {
		for iw := 0; iw < nFields; iw++ {
			il := i*nFields + iw
			ns := widths[iw] - runewidth.StringWidth(files[il].BaseName)
			tail := ""
			cname := files[il].ColorBaseName()
			if files[il].IsDir() {
				tail = "/" + paw.Repeat(" ", ns-1)
			} else {
				if files[il].IsLink() {
					cname = files[il].LSColorString(files[il].BaseName)
					tail = "@" + paw.Repeat(" ", ns-1)
				} else {
					tail = paw.Repeat(" ", ns)
				}
			}
			fmt.Fprintf(w, "%v", cname+tail)
		}
		fmt.Fprintln(w)
	}

	nw := nfolds * nFields
	if len(lens) > nw {
		for i := nw; i < len(lens); i++ {
			iw := i - nw
			ns := widths[iw] - runewidth.StringWidth(files[i].BaseName)
			tail := ""
			cname := files[i].ColorBaseName()
			if files[i].IsDir() {
				tail = "/" + paw.Repeat(" ", ns-1)
			} else {
				if files[i].IsLink() {
					cname = files[i].LSColorString(files[i].BaseName)
					tail = "@" + paw.Repeat(" ", ns-1)
				}
				tail = paw.Repeat(" ", ns)
			}
			fmt.Fprintf(w, "%v", cname+tail)
		}
		fmt.Fprintln(w)
	}
	// fmt.Fprintln(w)
}

func calWidth(lens []int, nFields, limit int) (widths []int) {
	nfolds := len(lens) / nFields
	widths = make([]int, nFields)
	copy(widths, lens[:nFields])
	for i := 0; i < nfolds; i++ {
		for iw := 0; iw < nFields; iw++ {
			il := i*nFields + iw
			if widths[iw] < lens[il] {
				widths[iw] = lens[il]
			}
		}
	}

	nw := nfolds * nFields
	if len(lens) > nw {
		for i := nw; i < len(lens); i++ {
			iw := i - nw
			if widths[iw] < lens[i] {
				widths[iw] = lens[i]
			}
		}
	}

	sw := sum(widths)
	if sw > limit {
		end := 0
		s := sum(widths)
		for i := len(widths) - 1; i > 0; i-- {
			s -= widths[i]
			if s < limit {
				end = i
				break
			}
		}
		widths = calWidth(lens, end, limit)
	}
	return widths
}

func sum(a []int) int {
	s := 0
	for _, v := range a {
		s += v
	}
	return s
}

func calNFields(lens []int, limit int) int {
	count := len(lens)
	n := 0
	for i := 0; i < len(lens); i++ {
		sum := 0
		for j := i; j < len(lens); j++ {
			sum += lens[j]
			if sum > limit {
				n = j - 1
				break
			}
		}
		if n < count {
			count = n
		}
	}
	return count
}

func classifyPrintFiles(w io.Writer, files []*File) {

	for _, file := range files {
		cname := file.ColorBaseName()
		if file.IsLink() {
			cname = file.LSColorString(file.BaseName) + "@"
		}
		if file.IsDir() {
			fmt.Fprintf(w, "%s/  ", cname)
		} else {
			fmt.Fprintf(w, "%s  ", cname)
		}
	}
	fmt.Fprintln(w)
}

func getleng(files []*File) (leng []int, sum int) {
	sum = 0
	for _, file := range files {
		lenstr := runewidth.StringWidth(file.BaseName) + 2
		// h, a := paw.CountPlaceHolder(file.BaseName)
		// lenstr := h + a + 2
		if file.IsDir() {
			lenstr++
		}
		leng = append(leng, lenstr)
		sum += lenstr
	}
	return leng, sum
}
