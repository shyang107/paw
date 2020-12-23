package filetree

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/spf13/cast"

	"code.cloudfoundry.org/bytefmt"

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
	SetNoColor()
	str := f.ToTextString("")
	DefaultNoColor()
	return str

	var (
		w    = new(bytes.Buffer)
		dirs = f.Dirs()
		fm   = f.Map()
	)

	i1 := len(cast.ToString(f.NDirs()))
	j1 := len(cast.ToString(f.NFiles()))
	if f.depth == 0 {
		if i1 < j1 {
			i1 = j1
		} else {
			j1 = i1
		}
	}
	// i1 := len(cast.ToString(len(dirs)))
	j := 0
	fmt.Fprintf(w, "%s\n", strings.Repeat("=", 80))
	for i, dir := range dirs {
		istr := fmt.Sprintf("%[2]*[1]d.", i, i1)
		sumsize := uint64(0)
		nfiles := 0
		for jj, file := range fm[dir] {
			mode := file.Stat.Mode()
			perm := getColorizePermission(mode)
			fsize := file.Size
			sdsize := getColorizedSize(fsize)
			if file.IsDir() {
				sdsize = cpmap['-'].Sprint(fmt.Sprintf("%6s", "--"))
			}
			if jj == 0 && file.IsDir() {
				if strings.EqualFold(file.Dir, RootMark) {
					fmt.Fprintf(w, fmt.Sprintf("%v %10v %v root (%v)\n", istr, perm, sdsize, f.Root()))
				} else {
					if f.depth != 0 {
						fmt.Fprintf(w, fmt.Sprintf("%v %10v %v %v\n", istr, perm, sdsize, file.Dir))
					} else {
						fmt.Fprintf(w, fmt.Sprintf("%v %10v %v %v\n", istr, perm, sdsize, file.BaseName))
					}
				}
				continue
			}
			jstr := ""
			if !file.IsDir() {
				sumsize += fsize
				j++
				nfiles++
				jstr = fmt.Sprintf("%[2]*[1]d.", j, j1)
			} else {
				jstr = fmt.Sprintf("%[2]*[1]s ", "", j1)
			}
			name := file.LSColorString(file.BaseName)
			link := checkAndGetLink(file)
			if len(link) > 0 {
				name += " -> " + link
			}

			if f.depth == 0 {
				fmt.Fprintf(w, fmt.Sprintf("%v %10v %6s %v\n", jstr, perm, sdsize, name))
			} else {
				fmt.Fprintf(w, fmt.Sprintf("    %v %10v %6s %v\n", jstr, perm, sdsize, name))
			}
		}
		if f.depth != 0 {
			fmt.Fprintf(w, "    Sum: %v files, size: %v.\n", nfiles, bytefmt.ByteSize(sumsize))

			if i == len(dirs)-1 {
				break
			}
			fmt.Fprintf(w, "%s\n", strings.Repeat("-", 80))
		}

	}
	fmt.Fprintf(w, "%s\n", strings.Repeat("=", 80))

	printTotalSummary(w, "", f.NDirs(), f.NFiles(), f.totalSize)

	return string(w.Bytes())
}

func checkAndGetLink(file *File) (link string) {
	mode := file.Stat.Mode()
	if mode&os.ModeSymlink != 0 {
		alink, err := filepath.EvalSymlinks(file.Path)
		if err != nil {
			link += alink + " ERR: " + err.Error()
		} else {
			link = alink
		}
	}

	return link
}

func checkAndGetColorLink(file *File) (link string) {
	mode := file.Stat.Mode()
	if mode&os.ModeSymlink != 0 {
		alink, err := filepath.EvalSymlinks(file.Path)
		if err != nil {
			link = alink + " ERR: " + err.Error()
		} else {
			link, _ = FileLSColorString(alink, alink)
		}
	}
	return link
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

// GetGitStatus will return git short status of `FileList`
func (f *FileList) GetGitStatus() GitStatus {
	return f.gitstatus
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
		if !strings.EqualFold(pdir, file.Dir) {
			f.store[pdir] = append(f.store[pdir], file)
		}
	}
}

func findPreDir(dir string) string {
	ddirs := strings.Split(dir, PathSeparator)
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
	if f.IsDir() && strings.HasPrefix(f.BaseName, ".") {
		return SkipDir
	}
	if strings.HasPrefix(f.BaseName, ".") {
		return SkipFile
	}
	return nil
}

// FindFiles will find files using codintion `ignore` func
// 	depth : depth of subfolders
// 		< 0 : walk through all directories of {root directory}
// 		0 : {root directory}/*
// 		1 : {root directory}/{level 1 directory}/*
//		...
// 	`ignore` IgnoreFn func(f *File, err error) error
// 		ignoring condition of files or directory
// 		`ignore` == nil, using `DefaultIgnoreFn`
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
			return strings.ToLower(files[i]) < strings.ToLower(files[j])
		})
		file := ConstructFileRelTo(root, root)
		f.AddFile(file)
		for _, name := range files {
			file := ConstructFileRelTo(root+PathSeparator+name, root)
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
				file := ConstructFileRelTo(path, root)
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
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)

				// For the purposes of this example, a simple SkipNode will suffice,
				// although in reality perhaps additional logic might be called for.
				return godirwalk.SkipNode
			},
			Unsorted: true, // set true for faster yet non-deterministic enumeration (see godoc)
		})
		if err != nil {
			return errors.New(root + ": " + err.Error())
		}
		// sort
		sort.Slice(f.dirs, func(i, j int) bool {
			return strings.ToLower(f.dirs[i]) < strings.ToLower(f.dirs[j])
		})
		for _, dir := range f.dirs {
			sort.Slice(f.store[dir], func(i, j int) bool {
				return strings.ToLower(f.store[dir][i].Path) < strings.ToLower(f.store[dir][j].Path)
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
// 			// sfsize := bytefmt.ByteSize(fsize)
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

	tf := &paw.TableFormat{
		Fields:    []string{"No.", "Mode", "Size", "Files"},
		LenFields: []int{5, 10, 6, 80},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignRight, paw.AlignRight, paw.AlignLeft},
		Padding:   pad,
	}
	tf.Prepare(buf)

	sdsize := bytefmt.ByteSize(f.totalSize)
	head := fmt.Sprintf("Root directory: %v, size: %v", f.Root(), sdsize)
	tf.SetBeforeMessage(head)

	SetNoColor()
	tf.PrintSart()
	j := 0
	for i, dir := range dirs {
		sumsize := uint64(0)
		nfiles := 0
		for jj, file := range fm[dir] {
			fsize := file.Size
			sfsize := bytefmt.ByteSize(fsize)
			mode := file.Stat.Mode()
			if jj == 0 && file.IsDir() {
				idx := fmt.Sprintf("D%d", i)
				sfsize = "-"
				if f.depth != 0 {
					if strings.EqualFold(file.Dir, RootMark) {
						tf.PrintRow(idx, mode, sfsize, file.LSColorString(f.Root()))
					} else {
						tf.PrintRow(idx, mode, sfsize, file.LSColorString(file.Dir))
					}
				} else if i > 0 {
					tf.PrintRow(idx, mode, sfsize, file)

				}
				continue
			}

			name := file.LSColorString(file.BaseName)
			link := checkAndGetLink(file)
			if len(link) > 0 {
				name += cpmap['l'].Sprint(" -> ") + link
			}
			if !file.IsDir() {
				sumsize += fsize
				j++
				nfiles++
				tf.PrintRow(j, mode, sfsize, name)
			} else {
				tf.PrintRow("", mode, "--", name)
			}

		}
		if f.depth != 0 {
			tf.PrintRow("", "", "", fmt.Sprintf("Sum: %v files, size: %v.", nfiles, bytefmt.ByteSize(sumsize)))

			if i != len(dirs)-1 {
				tf.PrintMiddleSepLine()
			}
		}
	}

	tf.SetAfterMessage(fmt.Sprintf("\n%v directories, %v files, total %v.", nDirs, nFiles, bytefmt.ByteSize(f.totalSize)))
	tf.PrintEnd()

	DefaultNoColor()
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

	sdsize := bytefmt.ByteSize(f.totalSize)
	fmt.Fprintf(w, "%sRoot directory: %v, size: %v\n", pad, getDirName(f.Root(), ""), KindLSColorString("di", sdsize))
	fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("=", width))

	ppad := ""

	i1 := len(cast.ToString(f.NDirs()))
	j1 := len(cast.ToString(f.NFiles()))
	if f.depth == 0 {
		if i1 < j1 {
			i1 = j1
		} else {
			j1 = i1
		}
	}
	j := 0
	for i, dir := range dirs {
		istr := KindLSColorString("di", fmt.Sprintf("%[2]*[1]d.", i, i1))
		sumsize := uint64(0)
		nfiles := 0
		ndirs := 0
		for jj, file := range fm[dir] {
			mode := file.Stat.Mode()
			cperm := getColorizePermission(mode)
			fsize := file.Size
			cfsize := getColorizedSize(fsize)
			if file.IsDir() {
				cfsize = cpmap['-'].Sprint(fmt.Sprintf("%6s", "-"))
			}
			if jj == 0 && file.IsDir() {
				if f.depth != 0 {
					if strings.EqualFold(file.Dir, RootMark) {
						// printFileItem(w, pad, istr, cperm, file.LSColorString(file.Dir)+" ("+getDirName(f.Root(), "")+")")
						printFileItem(w, pad, istr, cperm, getDirName(f.Root(), ""))
					} else {
						ppad = strings.Repeat("    ", len(file.DirSlice())-1)
						printFileItem(w, pad+ppad, istr, cperm, getDirName(file.Dir, f.Root()))
					}
				} else {
					ppad = strings.Repeat("    ", len(file.DirSlice())-1)
					fmt.Fprintf(w, "%s%v %10v %v %v\n", pad+ppad, istr, mode, cfsize, file.LSColorString(RootMark))
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
			name := file.LSColorString(file.BaseName)
			link := checkAndGetColorLink(file)
			if len(link) > 0 {
				name += cpmap['l'].Sprint(" -> ") + link
			}
			// fmt.Fprintf(w, "%s    %v %10v %v %v\n", pad+ppad, jstr, cperm, cfsize, name)
			printFileItem(w, pad+ppad+"    ", jstr, cperm, cfsize, name)
		}
		if f.depth != 0 {
			printDirSummary(w, pad+ppad+"    ", ndirs, nfiles, sumsize)

			if i != len(dirs)-1 {
				fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("-", width))
			}
		}
	}

	fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("=", width))
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
		w              = new(bytes.Buffer)
		dirs           = f.Dirs()
		fm             = f.Map()
		currentuser, _ = user.Current()
		urname         = currentuser.Username
		usergp, _      = user.LookupGroupId(currentuser.Gid)
		gpname         = usergp.Name
	)

	ctdsize := bytefmt.ByteSize(f.totalSize)
	head := fmt.Sprintf("%sRoot directory: %v, size: %v", pad, getDirName(f.Root(), ""), KindLSColorString("di", ctdsize))
	fmt.Fprintln(w, head)
	fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("=", 80))

	chead := getColorizedHead(pad, urname, gpname)
	fmt.Fprintln(w, chead)
	curname, cgpname := getColorizedUGName(urname, gpname)

	for i, dir := range dirs {
		sumsize := uint64(0)
		nfiles := 0
		ndirs := 0
		for jj, file := range fm[dir] {
			cperm := getColorizePermission(file.Stat.Mode())
			cmodTime := getColorizedModTime(file.Stat.ModTime())
			cgit := getColorizedGitStatus(f.GetGitStatus(), file)
			fsize := file.Size
			cfsize := getColorizedSize(fsize)
			if file.IsDir() {
				cfsize = cpmap['-'].Sprint(fmt.Sprintf("%6s", "-"))
			}
			if jj == 0 && file.IsDir() {
				reldir := file.Dir
				if !strings.EqualFold(file.Dir, RootMark) {
					if f.depth != 0 {
						fmt.Fprintf(w, "%s%v\n", pad, KindLSColorString("di", reldir))
						fmt.Fprintln(w, chead)
					}
				}
				continue
			}
			if file.IsDir() {
				ndirs++
			} else {
				nfiles++
			}
			sumsize += fsize
			name := file.LSColorString(file.BaseName)
			link := checkAndGetColorLink(file)
			if len(link) > 0 {
				name += cpmap['l'].Sprint(" -> ") + link
			}
			printFileItem(w, pad, cperm, cfsize, curname, cgpname, cmodTime, cgit, name)
		}

		if f.depth != 0 {
			printDirSummary(w, pad, ndirs, nfiles, sumsize)
			if i == len(dirs)-1 {
				break
			}
			fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("-", 80))
		}
	}

	fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("=", 80))
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

	// print heat
	if isMeta {
		chead := getColorizedHead(pad, urname, gpname)
		buf.WriteString(chead)
		buf.WriteByte('\n')
	}

	files := fm[RootMark]
	file := files[0]
	nfiles := len(files)

	// print root file
	meta := pad
	name := getName(file)
	if isMeta {
		meta = getMeta(pad, file, f.GetGitStatus())
	} else {
		meta += "[" + KindLSColorString("di", fmt.Sprintf("%d dirs", f.NDirs())) + ", " + KindLSColorString("fi", fmt.Sprintf("%d files", f.NFiles())) + "] "
	}

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
