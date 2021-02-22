package filetree

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
	// "github.com/shyang107/paw/treeprint"
)

// FileMap stores directory map to `map[{{ sub-path }}]{{ *File }}`
type FileMap map[string][]*File

func (f FileMap) Count() int {
	return len(f)
}

func (f FileMap) CountDF(dir string) (nd, nf int, size uint64) {
	if len(dir) > 0 {
		if _, ok := f[dir]; !ok {
			return 0, 0, 0
		}
		for _, v := range f[dir] {
			if v.IsDir() {
				nd++
			} else {
				nf++
				size += v.Size
			}
		}
		return nd, nf, size
	}
	for _, files := range f {
		for _, v := range files[:] {
			if v.IsDir() {
				nd++
			} else {
				nf++
				size += v.Size
			}
		}
	}
	return nd, nf, size
}

// FileList stores the list information of File
type FileList struct {
	root      string   // root directory
	store     FileMap  // all files in `root` directory
	dirs      []string // keys of `store`
	depth     int
	totalSize uint64
	git       *GitStatus
	sb        *strings.Builder
	writer    io.Writer
	// writers   []io.Writer
	IsSort    bool // increasing order of Lower(path)
	filesBy   FilesBy
	dirsBy    DirsBy
	IsGrouped bool // grouping files and directories separetly
	ignore    IgnoreFunc
	errors    []*flError
	mux       sync.RWMutex // 互斥鎖
}
type flError struct {
	path     string
	dir      string
	basename string
	err      error
}

func newFileListError(path string, err error, root string) *flError {
	dir := filepath.Dir(path)
	adir := strings.Replace(dir, root, ".", 1)
	basename := filepath.Base(path)
	// paw.Logger.WithFields(logrus.Fields{
	// 	"path": path,
	// 	"dir":  dir,
	// 	"adir": adir,
	// 	"name": basename,
	// }).Debug()
	return &flError{
		path:     path,
		dir:      adir,
		basename: basename,
		err:      err,
	}
}

// func (f flErrors)
// NewFileList will return the instance of `FileList`
func NewFileList(root string) *FileList {
	// if len(root) == 0 {
	// 	return &FileList{}
	// }
	fl := &FileList{
		root:      root,
		store:     make(map[string][]*File),
		dirs:      []string{},
		git:       &GitStatus{NoGit: true},
		sb:        new(strings.Builder),
		IsSort:    true,
		filesBy:   nil,
		dirsBy:    nil,
		IsGrouped: false,
		ignore:    nil,
	}
	fl.SetWriters(fl.sb)
	return fl
}

// String ...
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f *FileList) String() string {
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
	// f.ResetStringBuilder()
	f.sb.Reset()
	// f.writers = []io.Writer{}
	f.SetWriters(f.sb)
}

// // ResetStringBuilder will reset the buffer of FileList
// func (f *FileList) ResetStringBuilder() {
// 	f.stringBuilder.Reset()
// }

// StringBuilder will return the *strings.Builder buffer of FileList
func (f *FileList) StringBuilder() *strings.Builder {
	return f.sb
}

// Dump will dump buffer of FileList to a string
func (f *FileList) Dump() string {
	return f.sb.String()
}

// Writer will return the field writer of FileList
func (f *FileList) Writer() io.Writer {
	if f.writer == nil {
		f.sb.Reset()
		f.SetWriters(f.sb)
	}
	return f.writer
}

// ConfigGit sets up the git status of FileList. Git status of directory just reflects the git status of all sub-files.
//
// Piority:
// 	1. If dir is marked as GitIgnored, then marks all subfiles with GitIgnored.
// 	2. If any of subfiles of dir (including root) has any change of git status, set git status to dir.
// 		2.1 Git status of sub-files are same, the git status of dir is same.
// 		2.2 Git status of sub-files are mutilple status, the git status of dir is GitChanged.
func (f *FileList) ConfigGit() {
	paw.Logger.Info()

	f.git = NewGitStatus(f.root)
	gs := f.git.GetStatus()
	if f.git.NoGit || len(f.store) < 1 || gs == nil {
		return
	}

	// 1. check: if dir is GitIgnored, then marks all subfiles with GitIgnored.
	for _, dir := range f.dirs {
		if len(f.store[dir]) < 2 {
			continue
		}
		fm := f.store[dir]
		fd := fm[0]
		drp := fd.RelPath
		isMarkIgnored := false
		isUntracked := false
		if xy, ok := gs[drp]; ok {
			if isXY(xy, GitIgnored) {
				isMarkIgnored = true
			}
			if isXY(xy, GitUntracked) {
				isUntracked = true
			}
		}
		if isMarkIgnored || isUntracked {
			for _, file := range fm[1:] {
				rp := file.RelPath
				// 1. continue...
				gs[rp] = &GitFileStatus{
					Staging:  gs[drp].Staging,
					Worktree: gs[drp].Worktree,
					Extra:    file.BaseName,
				}
			}
		}
	}
	// 2. if any of subfiles of dir (including root) has any change of git status, set GitChanged to dir
	for _, dir := range f.dirs {
		files := f.GetFiles(dir)
		if files == nil || len(files) < 2 {
			continue
		}
		fd := files[0]
		xs, ys := f.getSubXYs(files[1:], gs)
		if len(xs) > 0 || len(ys) > 0 {
			rp := fd.RelPath
			gs[rp] = &GitFileStatus{
				Staging:  getSC(xs),
				Worktree: getSC(ys),
				Extra:    fd.BaseName + "/",
			}
		}
	}

	f.git.Dump("ConfigGit: modified")
}

// func addRootGitStatus(gs *GitStatus, repPath string) {
// 	if len(gs.status) > 0 { // has git status
// 		xs, ys := getXYs(gs.status)
// 		_, root := filepath.Split(repPath)
// 		root += PathSeparator
// 		gs.status[root] = &GitFileStatus{
// 			Staging:  getSC(xs),
// 			Worktree: getSC(ys),
// 			Extra:    root,
// 		}
// 	}
// }

func (f *FileList) getSubXYs(files []*File, status GStatus) (xs, ys []GitStatusCode) {
	xs = []GitStatusCode{}
	ys = []GitStatusCode{}
	for _, file := range files {
		rp := file.RelPath
		if !file.IsDir() {
			if xy, ok := status[rp]; ok {
				if xy.Staging != GitUnmodified {
					xs = append(xs, xy.Staging)
				}
				if xy.Worktree != GitUnmodified {
					ys = append(ys, xy.Worktree)
				}
			}
		} else {
			sxs, sys := f.getSubXYs(f.store[file.Dir+"/"+file.BaseName], status)
			xs = append(xs, sxs...)
			ys = append(ys, sys...)
		}
	}
	return xs, ys
}

func getXYs(status GStatus) (xs, ys []GitStatusCode) {
	for _, xy := range status {
		if xy.Staging != GitUnmodified {
			xs = append(xs, xy.Staging)
		}
		if xy.Worktree != GitUnmodified {
			ys = append(ys, xy.Worktree)
		}
	}
	return xs, ys
}

func isXY(xy *GitFileStatus, gcode GitStatusCode) bool {
	return xy.Staging == gcode ||
		xy.Worktree == gcode
}

// ConfigGit sets up the git status of FileList
func (f *FileList) GetGitStatus() *GitStatus {
	return f.git
}

// GitXY will retun git XY status of FileList
func (f *FileList) GitXY(relpath string) string {
	return f.git.XY(relpath)
}

// GitXY will retun git colorful XY status of FileList
func (f *FileList) GitXYc(relpath string) string {
	return f.git.XYc(relpath)
}

// SetIgnoreFunc set ignore function to FieldList.ignore
func (f *FileList) SetIgnoreFunc(ignore IgnoreFunc) {
	f.ignore = ignore
}

// Root set path to FieldList.root
func (f *FileList) SetRoot(path string) {
	f.root = path
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

// GetFiles will retun subfiles of dir of FileList
func (f *FileList) GetFiles(dir string) []*File {
	if files, ok := f.store[dir]; ok {
		return files[:]
	}
	return nil
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
		nd, nf, s := f.store.CountDF(dir)
		ndirs += nd - 1
		nfiles += nf
		size += s
	}
	f.totalSize = size
	return ndirs, nfiles, f.totalSize
}

// TotalSummary will return information about dir.
//
// 	Example:
// 	2 directories, 2 files, size ≈ 0b.
func (f *FileList) TotalSummary(wdstty int) string {
	ndirs, nfiles, sumsize := f.NTotalDirsAndFile()
	if wdstty <= 0 {
		wdstty = sttyWidth - 2
	}
	return totalSummary("", ndirs, nfiles, sumsize, wdstty)
}

// DirSummary will return information about dir.
//
// 	Example:
// 	2 directories, 2 files, size ≈ 0b.
func (f *FileList) DirSummary(dir string, wdstty int) string {
	ndirs, nfiles, sumsize := f.NSubDirsAndFiles(dir)
	return dirSummary("", ndirs, nfiles, sumsize, wdstty)
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
	nsDirs, nsFiles, sumsize = f.store.CountDF(dir)
	nsDirs--
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
	// return getDirInfo(f, file)

	if !file.IsDir() {
		return "", 0
	}

	nd, nf, _ := f.store.CountDF(file.Dir)
	nd--
	di := fmt.Sprintf("%d dirs", nd)
	fi := fmt.Sprintf("%d files", nf)
	wdinf = len(di) + len(fi) + 4
	cdinf = fmt.Sprintf("[%v, %v]", cdirp.Sprint(di), cdirp.Sprint(fi))
	return cdinf, wdinf
}

// // GetHead4Meta will return a colorful string of head line for meta information of File
// func (f *FileList) GetHead4Meta(pad, username, groupname string, git GitStatus) (chead string, width int) {
// 	return getColorizedHead(pad, username, groupname, git)
// }

// AddFile will add file into the FileList
func (f *FileList) AddFile(file *File) {
	var dir = file.Dir
	if _, ok := f.store[dir]; !ok {
		f.store[dir] = []*File{}
		f.dirs = append(f.dirs, dir)
	}
	f.store[dir] = append(f.store[dir], file)
	if file.IsFile() {
		f.totalSize += file.Size
	}

	var (
		predir string
		dirs   []string
	)
	if dir == RootMark {
		predir = dir
	} else {
		dirs = file.DirSlice() //strings.Split(dir, PathSeparator)
		ndirs := len(dirs)
		predir = strings.Join(dirs[:ndirs-1], PathSeparator)
	}
	var pdir = f.store[predir][0]
	file.SetUpDir(pdir)

	pdir = f.store[file.Dir][0]
	if file.IsDir() && file.Path != f.root {
		dir = file.GetSubDir()
		if _, ok := f.store[dir]; !ok {
			f.store[dir] = []*File{}
			f.dirs = append(f.dirs, dir)
		}
		dfile, _ := NewFileRelTo(file.Path, f.root)
		dfile.Dir = dir
		dfile.SetUpDir(pdir)
		f.store[dir] = append(f.store[dir], dfile)
		// paw.Logger.WithFields(logrus.Fields{
		// 	"dir":   dir,
		// 	"ddir":  dfile.Dir,
		// 	"pdir":  pdir.Dir,
		// 	"ppath": pdir.Path,
		// 	"path":  dfile.Path,
		// }).Info("AddFile: dfile updir")
	}
}

func (f *FileList) addFilePD(file *File) {
	var dir = file.Dir
	if _, ok := f.store[dir]; !ok {
		f.store[dir] = []*File{}
		f.dirs = append(f.dirs, dir)
		// f.totalSize += file.Size
	}
	f.store[dir] = append(f.store[dir], file)
	if file.IsFile() {
		f.totalSize += file.Size
	}
}

// AddError will add file into the FileList.errors
func (f *FileList) AddError(path string, err error) {
	f.errors = append(f.errors, newFileListError(path, err, f.root))
}

// SetMd5 will set isMd5 to hasMd5, then will get md5 of File
func (f *FileList) SetMd5(isMd5 bool) {
	hasMd5 = isMd5
}

// GetErrorString get the error string in `dir` during find files
func (f *FileList) GetErrorString(dir string) string {
	if len(f.errors) == 0 {
		return ""
	}
	sb := new(strings.Builder)
	for _, e := range f.errors {
		if e.dir == dir {
			sb.WriteString(cerror.Sprint(e.err.Error()))
			sb.WriteRune('\n')
		}
	}
	if len(sb.String()) == 0 {
		return ""
	}
	return sb.String()
}

// GetAllErrorString get the all error string in `dir` during find files
func (f *FileList) GetAllErrorString() string {
	if len(f.errors) == 0 {
		return ""
	}
	sb := new(strings.Builder)
	for _, e := range f.errors {
		sb.WriteString(cerror.Sprint(e.err.Error()))
		sb.WriteRune('\n')
	}
	if len(sb.String()) == 0 {
		return ""
	}
	return sb.String()
}

// FprintErrs prints out error string in `dirent` during find files
func (f *FileList) FprintErrs(w io.Writer, dirent, pad string) {
	errmsg := f.GetErrorString(dirent)
	if len(errmsg) > 0 {
		if len(pad) > 0 {
			errmsg = paw.PaddingString(errmsg, pad)
		}
		fmt.Fprint(w, errmsg)
	}
}

// FprintAllErrs prints out all error string during find files
func (f *FileList) FprintAllErrs(w io.Writer, pad string) {
	errmsg := f.GetAllErrorString()
	if len(errmsg) > 0 {
		if len(pad) == 0 {
			errmsg = paw.PaddingString(errmsg, pad)
		}
		fmt.Fprint(w, errmsg)
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
var SkipThis = errors.New("skip this node")

// IgnoreFn is the type of the function called for each file or directory
// visited by FindFiles. The f argument contains the File argument to FindFiles.
//
// If there was a problem walking to the file or directory named by path, the
// incoming error will describe the problem and the function can decide how
// to handle that error (and FindFiles will not descend into that directory).
// In the case of an error, the info argument will be nil. If an error is
// returned, processing stops. The sole exception is when the function returns
// the special value SkipThis or ErrSkipFile. If the function returns
// SkipThis when invoked on a directory,
// FindFiles skips the directory's contents entirely. If the function returns
// SkipThis when invoked on a non-directory file, FindFiles skips the remaining
// files in the containing directory.
// If the returned error is SkipFile when inviked on a file, FindFiles will
// skip the file.
type IgnoreFunc func(f *File, err error) error

// DefaultIgnoreFn is default IgnoreFn using in FindFiles
//
// _, file := filepath.Split(f.Path)
// 	Skip file: prefix "." of file
var DefaultIgnoreFn IgnoreFunc = func(f *File, err error) error {
	if err != nil {
		return SkipThis
	}
	if strings.HasPrefix(f.BaseName, ".") {
		return SkipThis
	}
	return nil
}

var (
	wg = sync.WaitGroup{}
	// limit goroute number
	sem = make(chan int, 12)
)

func wgosReadDir(f *FileList, dirPath string) {
	sem <- 1
	defer func() {
		<-sem
	}()
	defer wg.Done()

	des, err := os.ReadDir(dirPath) // openDir.Readdirnames(-1)
	if err != nil {
		paw.Logger.Error(err)
		matchs := repath.FindAllString(err.Error(), -1)
		path := dirPath
		if len(matchs) > 0 {
			path = matchs[0]
		}
		f.mux.Lock()
		f.AddError(path, err)
		f.mux.Unlock()
		// return
	}
	if len(des) > 0 {
		wg.Add(1)
		go wgosrdHandleFiles(f, dirPath, des)
	}

	return
}

func wgosrdHandleFiles(f *FileList, dirPath string, des []fs.DirEntry) {
	sem <- 1
	defer func() {
		<-sem
	}()
	defer wg.Done()

	nf := len(des)
	if nf == 0 {
		return
	}

	if nf == 1 {
		skip := false
		name := des[0].Name()
		path := filepath.Join(dirPath, name)
		file, err := NewFileRelTo(path, f.root)
		if err != nil {
			paw.Logger.Error(err)
			f.mux.Lock()
			f.AddError(path, err)
			f.mux.Unlock()
			return
		}

		if err := f.ignore(file, nil); !skip && err == SkipThis {
			if pdOpt.isTrace {
				f.mux.Lock()
				f.AddError(file.Path, fmt.Errorf("%s: %q", err, file.Path))
				f.mux.Unlock()
			}
			skip = true
		}

		idepth := len(file.DirSlice()) - 1
		if !skip && f.depth > 0 {
			if idepth > f.depth {
				if pdOpt.isTrace {
					f.mux.Lock()
					f.AddError(path, fmt.Errorf("The depth (%d) of path > configured depth (%d)  ", f.depth, idepth))
					f.mux.Unlock()
				}
				skip = true
			}
		}

		if !skip {
			f.mux.Lock()
			f.AddFile(file)
			f.mux.Unlock()
			if file.IsDir() {
				wg.Add(1)
				go wgosReadDir(f, path)
			}
		}
	} else {
		wg.Add(2)
		go wgosrdHandleFiles(f, dirPath, des[:nf/2])
		go wgosrdHandleFiles(f, dirPath, des[nf/2:])
	}
}

var repath = regexp.MustCompile(`(/.+/[^:]+?)`)

func osReadDir(f *FileList, dirPath string) {
	des, err := os.ReadDir(dirPath) //odir.ReadDir(-1) //odir.Readdirnames(-1)
	if err != nil {
		paw.Logger.Error(err)
		matchs := repath.FindAllString(err.Error(), -1)
		path := dirPath
		if len(matchs) > 0 {
			path = matchs[0]
		}
		f.AddError(path, err)
		// return
	}
	if len(des) > 0 {
		osrdHandleFiles(f, dirPath, des)
	}

	return
}

func osrdHandleFiles(f *FileList, dirPath string, des []fs.DirEntry) {
	nf := len(des)
	if nf == 0 {
		return
	}
	for _, de := range des {
		skip := false
		path := filepath.Join(dirPath, de.Name())
		file, err := NewFileRelTo(path, f.root)
		if err != nil {
			paw.Logger.Error(err)
			f.AddError(path, err)
			continue
		}

		if err := f.ignore(file, nil); err == SkipThis {
			if pdOpt.isTrace {
				f.AddError(file.Path, fmt.Errorf("%s: %q", err, file.Path))
			}
			skip = true
		}

		idepth := len(file.DirSlice()) - 1
		if !skip && f.depth > 0 {
			if idepth > f.depth {
				if pdOpt.isTrace {
					f.AddError(file.Path, fmt.Errorf("The depth (%d) of path > configured depth (%d)  ", f.depth, idepth))
				}
				skip = true
			}
		}

		if !skip {
			f.AddFile(file)
			if de.IsDir() {
				osReadDir(f, path)
			}
		}
	}
}

// FindFiles will find files using codintion `ignore` func
// 	depth : depth of subfolders
// 		depth < 0 : walk through all directories of {root directory}
// 		depth == 0 : {root directory}/
// 		depth > 0 : {root directory}/{level 1 directory}/.../{{ level n directory }}/
// 	ignore: IgnoreFn func(f *File, err error) error
// 		ignoring condition of files or directory
// 		ignore == nil, using DefaultIgnoreFn
func (f *FileList) FindFiles(depth int) error {
	paw.Logger.Infof("root: %q", f.root)

	if f.ignore == nil {
		f.ignore = DefaultIgnoreFn
	}
	// f.gitstatus, _ = GetShortGitStatus(f.root)
	f.depth = depth
	file, err := NewFileRelTo(f.root, f.root)
	if err != nil {
		paw.Logger.Error(err)
		f.AddError(f.root, err)
		return err
	}
	f.AddFile(file)
	if file.IsLink() {
		f.root = file.LinkPath()
		f.git = NewGitStatus(f.root)
		f.depth = depth
	}

	if hasMd5 {
		paw.Logger.Trace("finding files starts... (goroutine)" + paw.Caller(1))
		wg.Add(1)
		go wgosReadDir(f, f.root)
		wg.Wait()
	} else {
		paw.Logger.Trace("finding files starts..." + paw.Caller(1))
		osReadDir(f, f.root)
		// fpWalkDir(f)
		// gdWalk(f)
	}

	if f.IsSort {
		f.Sort()
	}
	return nil
}

var (
	DefaultFilesBy FilesBy = byName

	// DefaultFilesBy FilesBy = func(fi *File, fj *File) bool {
	// 	if fi.IsDir() && fj.IsFile() {
	// 		return true
	// 	} else if fi.IsFile() && fj.IsDir() {
	// 		return false
	// 	}
	// 	return strings.ToLower(fi.Path) < strings.ToLower(fj.Path)
	// }

	DefaultDirsBy DirsBy = func(di string, dj string) bool {
		return strings.ToLower(di) < strings.ToLower(dj)
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
	paw.Logger.Info("sort...")

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

// Sort0 will sort FileList by sorter of dirsBy and filesBy. (for FileList.Map()[dir][0:])
//
// Default:
// 	Dirs: ToLower(a[i]) < ToLower(a[j])
// 	Map[dir][]*file: ToLower(a[i].Path) < ToLower(a[j].Path)
func (f *FileList) Sort0() {
	f.SortBy0(f.dirsBy, f.filesBy)
}

// SortBy0 will sort FileList using sorters `dirsBy` and `filesBy`. (for FileList.Map()[dir][0:] )
func (f *FileList) SortBy0(dirsBy DirsBy, filesBy FilesBy) {
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
						f.filesBy.Sort(f.store[dir][:])
					} else {
						sfiles := []*File{}
						sdirs := []*File{}
						for _, v := range f.store[dir][:] {
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

func (f *FileList) dumpAll(loglevel logrus.Level) {
	if loglevel != logrus.TraceLevel {
		return
	}
	for _, dir := range f.Dirs() {
		fm := f.Map()[dir]
		level := len(fm[0].DirSlice()) - 1
		sp := paw.Spaces(level * 2)
		fmt.Printf("%sG%d  dir: %q\n", sp, level, paw.Truncate(dir, 60, "..."))
		for j, f := range fm[:] {
			var (
				pname = "x"
				pdir  *File
				// pdirName = RootMark
				fdir, fname = "x", "x"
				rpath       = f.RelPath
			)
			if f.GetUpDir() != nil {
				pdir = f.GetUpDir()
				pname = cdip.Sprint(paw.Truncate(pdir.Dir, 25, "..."))
				fdir = cdip.Sprint(paw.Truncate(f.Dir, 25, "..."))
				fname = f.LSColor().Sprint(paw.FillRight(paw.Truncate(f.Name(), 15, "..."), 15))
			}
			fmt.Printf("%s  %2d dir: \"%v\" pdir: \"%v\" name: \"%v\" rpath: \"%v\"", sp, j, fdir, pname, fname, rpath)
			// fmt.Printf("  %v Excutable: %s owner: %s group: %s others: %s any: %s all: %s", f.PermissionC(),
			// 	bmark(f.IsExecutable()), bmark(f.IsExecOwner()), bmark(f.IsExecGroup()), bmark(f.IsExecOther()), bmark(f.IsExecAny()), bmark(f.IsExecAll()))
			fmt.Print("\n")
		}
	}
}

// DoView will print out FileList according to `out`
func (fl *FileList) DoView(view PDViewFlag, pad string) error {
	switch view {
	case PListView:
		fl.ToListView(pad)
	case PListExtendView:
		fl.ToListExtendView(pad)
	case PTreeView:
		fl.ToTreeView(pad)
	case PTreeExtendView:
		fl.ToTreeExtendView(pad)
	case PListTreeView:
		fl.ToListTreeView(pad)
	case PListTreeExtendView:
		fl.ToListTreeExtendView(pad)
	case PLevelView:
		fl.ToLevelView(pad, false)
	case PLevelExtendView:
		fl.ToLevelView(pad, true)
	case PTableView:
		fl.ToTableView(pad, false)
	case PTableExtendView:
		fl.ToTableView(pad, true)
	case PClassifyView:
		fl.ToClassifyView(pad)
	default:
		return errors.New("No this view option of PrintDir")
	}
	return nil
}
