package filetree

import (
	"errors"
	"io"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/karrick/godirwalk"

	"github.com/shyang107/paw"
	// "github.com/shyang107/paw/treeprint"
)

// FileMap stores directory map to `map[{{ sub-path }}]{{ *File }}`
type FileMap map[string][]*File

// FileList stores the list information of File
type FileList struct {
	root          string   // root directory
	store         FileMap  // all files in `root` directory
	dirs          []string // keys of `store`
	depth         int
	totalSize     uint64
	gitstatus     GitStatus
	stringBuilder *strings.Builder
	writer        io.Writer
	// writers   []io.Writer
	IsSort    bool // increasing order of Lower(path)
	filesBy   FilesBy
	dirsBy    DirsBy
	IsGrouped bool // grouping files and directories separetly
	// mux       sync.Mutex // 互斥鎖
}

// NewFileList will return the instance of `FileList`
func NewFileList(root string) *FileList {
	if len(root) == 0 {
		return &FileList{}
	}
	fl := &FileList{
		root:          root,
		store:         make(map[string][]*File),
		dirs:          []string{},
		stringBuilder: paw.NewStringBuilder(),
		IsSort:        true,
		filesBy:       nil,
		dirsBy:        nil,
		IsGrouped:     false,
	}
	fl.SetWriters(fl.stringBuilder)
	return fl
}

// String ...
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f FileList) String() string {
	// fmt.Printf("%#v\n", f.writer)
	// oldwr := f.writer
	// f.SetWriters(paw.NewStringBuilder())
	f.DisableColor()
	str := f.ToLevelView("", false)
	f.EnableColor()
	// f.SetWriters(oldwr)
	return str
}

// SetWriters will set writers... to writer of FileList
func (f *FileList) SetWriters(writers ...io.Writer) {
	// f.ResetBuffer()
	// f.writers = append(f.writers, writers...)
	// f.writers = writers
	// if len(writers) == 0 || writers == nil {
	// 	paw.Info.Println("len(writers) == 0 || writers == nil")
	// 	f.ResetWriters()
	// 	return
	// }
	// if len(writers) == 1 && writers[0] == nil {
	// 	paw.Info.Println("len(writers) == 1 && writers[0] == nil")
	// 	f.ResetWriters()
	// 	return
	// }
	// writers = append(writers, f.stringBuilder)
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
	f.ResetStringBuilder()
	// f.writers = []io.Writer{}
	f.SetWriters(f.stringBuilder)
}

// ResetStringBuilder will reset the buffer of FileList
func (f *FileList) ResetStringBuilder() {
	f.stringBuilder.Reset()
}

// StringBuilder will return the *strings.Builder buffer of FileList
func (f *FileList) StringBuilder() *strings.Builder {
	return f.stringBuilder
}

// Dump will dump buffer of FileList to a string
func (f *FileList) Dump() string {
	return f.stringBuilder.String()
}

// Writer will return the field writer of FileList
func (f *FileList) Writer() io.Writer {
	if f.writer == nil {
		f.ResetStringBuilder()
		f.SetWriters(f.stringBuilder)
	}
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

// Dirs will retun directories of FileList (keys of `FileMap`)
func (f *FileList) Dirs() []string {
	return f.dirs
}

// TotalSize will retun total size of FileList
func (f *FileList) TotalSize() uint64 {
	return f.totalSize
}

// TotalByteSize will retun total size string of FileList in byte-format as human read
func (f *FileList) TotalByteSize() string {
	return ByteSize(f.totalSize)
}

// ColorfulTotalByteSize will retun colorful total size string of FileList in byte-format as human read
func (f *FileList) ColorfulTotalByteSize() string {
	return GetColorizedSize(f.totalSize)
}

// NDirs is the numbers of sub-directories of `root`
func (f *FileList) NDirs() int {
	// return len(f.Dirs()) - 1
	ndirs, _, _ := f.NTotalDirsAndFile()
	return ndirs
}

// NFiles is the numbers of all files
func (f *FileList) NFiles() int {
	_, nfiles, _ := f.NTotalDirsAndFile()
	return nfiles
}

// NItems will return FileList.NDirs() + FileList.NFiles()
func (f *FileList) NItems() int {
	ndirs, nfiles, _ := f.NTotalDirsAndFile()
	return ndirs + nfiles
}

// NTotalDirsAndFile will return NDirs, NFiles and TotalSize of FileList
func (f *FileList) NTotalDirsAndFile() (ndirs, nfiles int, size uint64) {
	for _, dir := range f.dirs {
		for _, file := range f.store[dir][1:] {
			if file.IsDir() {
				ndirs++
			} else {
				nfiles++
				// size += file.Size
			}
		}
	}
	return ndirs, nfiles, f.totalSize
}

// TotalSummary will return information about dir.
//
// 	Example:
// 	2 directories, 2 files, size ≈ 0b.
func (f *FileList) TotalSummary() string {
	ndirs, nfiles, sumsize := f.NTotalDirsAndFile()
	return totalSummary("", ndirs, nfiles, sumsize)
}

// DirSummary will return information about dir.
//
// 	Example:
// 	2 directories, 2 files, size ≈ 0b.
func (f *FileList) DirSummary(dir string) string {
	ndirs, nfiles, sumsize := f.NSubDirsAndFiles(dir)
	return dirSummary("", ndirs, nfiles, sumsize)
}

// NSubDirs will return number of sub-directories of (key) dir
func (f *FileList) NSubDirs(dir string) int {
	nsdirs, _, _ := f.NSubDirsAndFiles(dir)
	return nsdirs
}

// NSubFiles will return number of sub-files of (key) dir
func (f *FileList) NSubFiles(dir string) int {
	_, nsfiles, _ := f.NSubDirsAndFiles(dir)
	return nsfiles
}

// NSubDirsAndFiles will return NSubDirs, NSubFiles and sum of size of FileList
func (f *FileList) NSubDirsAndFiles(dir string) (nsDirs, nsFiles int, sumsize uint64) {
	if files, ok := f.Map()[dir]; ok {
		for _, file := range files[1:] {
			if file.IsDir() {
				nsDirs++
			} else {
				nsFiles++
				sumsize += file.Size
			}
		}
	}
	return nsDirs, nsFiles, sumsize
}

// NSubItems will return number of sub-dirs and sub-files in dir
func (f *FileList) NSubItems(dir string) (ndirs, nfiles int) {
	ndirs, nfiles, _ = f.NSubDirsAndFiles(dir)
	return ndirs, nfiles
}

// SubSize will retun total size of dir of FileList
func (f *FileList) SubSize(dir string) uint64 {
	_, _, size := f.NSubDirsAndFiles(dir)
	return size
}

// SubByteSize will retun total size string of FileList in byte-format as human read
func (f *FileList) SubByteSize(dir string) string {
	return ByteSize(f.SubSize(dir))
}

// DirInfo will return the colorful string of sub-dir ( file.IsDir is true) and the width on console.
func (f *FileList) DirInfo(file *File) (cdinf string, wdinf int) {
	return getDirInfo(f, file)
}

// GetGitStatus will return git short status of `FileList`
func (f *FileList) GetGitStatus() GitStatus {
	return f.gitstatus
}

// // GetHead4Meta will return a colorful string of head line for meta information of File
// func (f *FileList) GetHead4Meta(pad, username, groupname string, git GitStatus) (chead string, width int) {
// 	return getColorizedHead(pad, username, groupname, git)
// }

// AddFile will add file into the file list
func (f *FileList) AddFile(file *File) {
	if _, ok := f.store[file.Dir]; !ok {
		f.store[file.Dir] = []*File{}
		f.dirs = append(f.dirs, file.Dir)
		// f.totalSize += file.Size
	}
	f.store[file.Dir] = append(f.store[file.Dir], file)
	if file.IsDir() {
		dd := paw.Split(file.Dir, PathSeparator)
		pdir := paw.Join(dd[:len(dd)-1], PathSeparator)
		if !paw.EqualFold(pdir, file.Dir) {
			f.store[pdir] = append(f.store[pdir], file)
		}
	} else {
		f.totalSize += file.Size
	}
}

func (f *FileList) DisableColor() {
	paw.SetNoColor()
}

func (f *FileList) EnableColor() {
	paw.DefaultNoColor()
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
	f.gitstatus, _ = GetShortGitStatus(f.root)
	f.depth = depth
	switch depth {
	case 0: //{root directory}/*
		// scratchBuffer := make([]byte, godirwalk.MinimumScratchBufferSize)
		files, err := godirwalk.ReadDirnames(f.root, nil)
		if err != nil {
			return errors.New(f.root + ": " + err.Error())
		}
		if f.IsSort {
			sort.Sort(ByLowerString(files))
		}

		file, err := NewFileRelTo(f.root, f.root)
		if err != nil {
			return err
		}
		f.AddFile(file)

		for _, name := range files {
			path := filepath.Join(f.root, name)
			file, err := NewFileRelTo(path, f.root)
			if err != nil {
				return err
			}
			if err := ignore(file, nil); err == SkipThis {
				continue
			}

			f.AddFile(file)
		}
	default: //walk through all directories of {root directory}
		err := godirwalk.Walk(f.root, &godirwalk.Options{
			Callback: func(path string, de *godirwalk.Dirent) error {
				file, err := NewFileRelTo(path, f.root)
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
			return errors.New(f.root + ": " + err.Error())
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

// SortBy will sort FileList using sorters `dirsBy` and `filesBy`
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

	wg := new(sync.WaitGroup)
	nCPU := runtime.NumCPU()
	nDirs := len(f.dirs)
	// paw.Info.Println("nCPU:", nCPU, "nDirs:", nDirs)
	for i := 0; i < nCPU; i++ {
		from := i * nDirs / nCPU
		to := (i + 1) * nDirs / nCPU
		wg.Add(1)
		// go sortParts(f, from, to, wg)
		go func() {
			defer wg.Done()
			for j := from; j < to; j++ {
				dir := f.dirs[j]
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
		}()
	}
	wg.Wait()
	return
	// for _, dir := range f.dirs {
	// 	if len(f.store[dir]) > 1 {
	// 		if !f.IsGrouped {
	// 			f.filesBy.Sort(f.store[dir][1:])
	// 		} else {
	// 			sfiles := []*File{}
	// 			sdirs := []*File{}
	// 			for _, v := range f.store[dir][1:] {
	// 				if v.IsDir() {
	// 					sdirs = append(sdirs, v)
	// 				} else {
	// 					sfiles = append(sfiles, v)
	// 				}
	// 			}
	// 			f.filesBy.Sort(sdirs)
	// 			f.filesBy.Sort(sfiles)
	// 			copy(f.store[dir][1:], sdirs)
	// 			copy(f.store[dir][len(sdirs)+1:], sfiles)
	// 		}
	// 	}
	// }
}
