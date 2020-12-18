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
	"time"

	"github.com/fatih/color"
	"github.com/shyang107/paw/cast"
	"github.com/shyang107/paw/godirwalk"

	"code.cloudfoundry.org/bytefmt"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/treeprint"
)

// FileMap ...
type FileMap map[string][]*File

// FileList ...
type FileList struct {
	root      string   // root directory
	store     FileMap  // all files in `root` directory
	dirs      []string // keys of `store`
	depth     int
	totalSize uint64
	gitstatus GitStatus
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

// String ...
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f FileList) String() string {
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
	for i, dir := range dirs {
		istr := fmt.Sprintf("%[2]*[1]d.", i, i1)
		for _, file := range fm[dir] {
			fsize := file.Size
			sdsize := bytefmt.ByteSize(fsize)
			mode := file.Stat.Mode()
			if file.IsDir() {
				if strings.EqualFold(file.Dir, RootMark) {
					fmt.Fprintf(w, file.LSColorString(fmt.Sprintf("%v %10v %6s root (%v)\n", istr, mode, sdsize, f.Root())))
				} else {
					if f.depth != 0 {
						fmt.Fprintf(w, file.LSColorString(fmt.Sprintf("%v %10v %6s %v\n", istr, mode, sdsize, file.Dir)))
					} else {
						fmt.Fprintf(w, file.LSColorString(fmt.Sprintf("%v %10v %6s %v\n", istr, mode, sdsize, file.BaseName)))
					}
				}
				continue
			}
			j++
			jstr := fmt.Sprintf("%[2]*[1]d", j, j1)
			if f.depth == 0 {
				fmt.Fprintf(w, file.LSColorString(fmt.Sprintf("%v. %10v %6s %v\n", jstr, mode, sdsize, file.BaseName)))
			} else {
				fmt.Fprintf(w, file.LSColorString(fmt.Sprintf("    %v. %10v %6s %v\n", jstr, mode, sdsize, file.BaseName)))
			}
		}

		if i == len(dirs)-1 {
			break
		}
	}

	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "%d directories, %d files, total %v\n", f.NDirs(), f.NFiles(), bytefmt.ByteSize(f.totalSize))
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

// GetGitStatus will return git short status of `FileList`
func (f *FileList) GetGitStatus() GitStatus {
	return f.gitstatus
}

// AddFile will add file into the file list
func (f *FileList) AddFile(file *File) {
	if _, ok := f.store[file.Dir]; !ok {
		f.store[file.Dir] = []*File{}
		f.dirs = append(f.dirs, file.Dir)
	}
	f.store[file.Dir] = append(f.store[file.Dir], file)
	// f.totalSize += file.Size()
}

func (f *FileList) DisableColor() {
	SetNoColor()
}

func (f *FileList) EnableColor() {
	DefaultNoColor()
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
	f.totalSize, _ = sizes(root)
	switch {
	case depth == 0: //{root directory}/*
		// fis, err := ioutil.ReadDir(root)
		// if err != nil {
		// 	return errors.New(root + ": " + err.Error())
		// }

		// for _, fi := range fis {
		// 	file := ConstructFileRelTo(root+PathSeparator+fi.Name(), root)
		// 	err := ignore(file, nil)
		// 	if err == SkipFile || err == SkipDir{
		// 		continue
		// 	}
		// 	f.AddFile(file)
		// }
		scratchBuffer := make([]byte, godirwalk.MinimumScratchBufferSize)
		files, err := godirwalk.ReadDirnames(root, scratchBuffer)
		if err != nil {
			return errors.New(root + ": " + err.Error())
		}
		sort.Strings(files)
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
	default: //walk through all directories of {root directory}
		// visit := func(path string, info os.FileInfo, err error) error {
		// 	file := ConstructFileRelTo(path, root)
		// 	idepth := len(file.DirSlice()) - 1
		// 	if depth > 0 {
		// 		if idepth > depth {
		// 			return nil
		// 		}
		// 	}
		// 	err1 := ignore(file, err)
		// 	if err1 == ErrSkipFile {
		// 		return nil
		// 	}
		// 	if err1 == ErrSkipDir {
		// 		return err1
		// 	}
		// 	f.AddFile(file)
		// 	return nil
		// }

		// err := filepath.Walk(root, visit)
		// if err != nil {
		// 	return errors.New(root + ": " + err.Error())
		// }
		err := godirwalk.Walk(root, &godirwalk.Options{
			Callback: func(path string, de *godirwalk.Dirent) error {
				file := ConstructFileRelTo(path, root)
				idepth := len(file.DirSlice()) - 1
				if depth > 0 {
					if idepth > depth {
						return nil
					}
				}
				err1 := ignore(file, nil)
				if err1 == SkipFile {
					return nil
				}
				if err1 == SkipDir {
					return err1
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
			// Unsorted: true, // set true for faster yet non-deterministic enumeration (see godoc)
		})
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

	for i, dir := range dirs {
		files := f.Map()[dir]
		nf := len(files) // including dir
		ntf += (nf - 1)
		for _, file := range files {
			fsize := file.Size
			sfsize := bytefmt.ByteSize(fsize)
			if file.IsDir() {
				if i == 0 { // root dir
					tree.SetValue(fmt.Sprintf("%v\n%v", file.LSColorString(file.Path), file.LSColorString(file.Dir)))
					tree.SetMetaValue(KindLSColorString("di", fmt.Sprintf("%v dirs., %v files, %v", nd-1, nf-1, sfsize)))
					one = tree
				} else {
					pre = preTree(dir, fm, tree)
					if f.depth != 0 {
						// one = pre.AddMetaBranch(nf-1, file)
						one = pre.AddMetaBranch(KindLSColorString("di", fmt.Sprintf("%d files, %v", nf-1, sfsize)), file)
					} else {
						one = pre.AddBranch(file)
					}
				}
				continue
			}
			// add file node
			one.AddNode(file)
		}
	}
	buf := new(bytes.Buffer)
	buf.Write(tree.Bytes())
	buf.WriteByte('\n')
	buf.WriteString(fmt.Sprintf("%d directoris, %d files, total %v.", f.NDirs(), f.NFiles(), bytefmt.ByteSize(f.totalSize)))
	// buf.WriteByte('\n')

	return paddingTree(pad, buf.Bytes())
}

func paddingTree(pad string, bytes []byte) []byte {
	b := make([]byte, len(bytes))
	b = append(b, pad...)
	for _, v := range bytes {
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
	// fmt.Println(dir, nd)
	if nd == 2 { // ./xx
		pre = tree
	} else { //./xx/...
		pre = tree
		for i := 2; i < nd; i++ {
			predir := strings.Join(dd[:i], PathSeparator)
			// fmt.Println("\t", i, predir)
			f := fm[predir][0] // import dir
			pre = pre.FindByValue(f)
		}
	}
	return pre
}

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
	head := fmt.Sprintf("Root directory: %v, size: %v", KindLSColorString("di", f.Root()), KindLSColorString("di", sdsize))
	tf.SetBeforeMessage(head)

	tf.PrintSart()
	SetNoColor()
	for i, dir := range dirs {
		for j, file := range fm[dir] {
			fsize := file.Size
			sfsize := bytefmt.ByteSize(fsize)
			mode := file.Stat.Mode()
			if file.IsDir() {
				idx := fmt.Sprintf("D%d", i)
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
			tf.PrintRow(j, mode, sfsize, file)
		}
		if f.depth != 0 {
			tf.PrintRow("", "", "", fmt.Sprintf("Sum: %v files.", len(fm[dir])-1))

			if i == len(dirs)-1 {
				break
			}
			tf.PrintMiddleSepLine()
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
	fmt.Fprintf(w, "%sRoot directory: %v, size: %v\n", pad, KindLSColorString("di", f.Root()), KindLSColorString("di", sdsize))
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

	for i, dir := range dirs {
		istr := KindLSColorString("di", fmt.Sprintf("%[2]*[1]d.", i, i1))
		for j, file := range fm[dir] {
			mode := file.Stat.Mode()
			fsize := file.Size
			sfsize := bytefmt.ByteSize(fsize)
			sizefile := file.LSColorString(fmt.Sprintf("%8s %v", sfsize, file))
			if file.IsDir() {
				sdsize := KindLSColorString("di", sfsize)
				if f.depth != 0 {
					if strings.EqualFold(file.Dir, RootMark) {
						fmt.Fprintf(w, "%s%v %10v %6s root (%v)\n", pad, istr, mode, sdsize, file.LSColorString(f.Root()))
					} else {
						ppad = strings.Repeat("    ", len(file.DirSlice())-1)
						fmt.Fprintf(w, "%s%v %10v %6s %v\n", pad+ppad, istr, mode, sdsize, file.LSColorString(file.Dir))
					}
				} else {
					ppad = strings.Repeat("    ", len(file.DirSlice())-1)
					fmt.Fprintf(w, "%s%v %10v %v\n", pad+ppad, istr, mode, sizefile)
				}
				continue
			}
			if f.depth != 0 {
				j1 = len(cast.ToString(len(fm[dir]) - 1))
			}
			jstr := fmt.Sprintf("%[2]*[1]d.", j, j1)
			fmt.Fprintf(w, "%s    %v %10v %v\n", pad+ppad, jstr, mode, sizefile)
		}
		if f.depth != 0 {
			fmt.Fprintf(w, "%s    Sum: %v files.\n", pad+ppad, len(fm[dir])-1)

			if i == len(dirs)-1 {
				break
			}
			fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("-", width))
		}
	}

	fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("=", width))
	fmt.Fprintln(w, pad)
	fmt.Fprintf(w, "%s%d directories, %d files, total %v.\n", pad, f.NDirs(), f.NFiles(), bytefmt.ByteSize(f.totalSize))
	// fmt.Fprintln(w, pad)
	return w.Bytes()
}

func (f *FileList) ToListString(pad string) string {
	return string(f.ToList(pad))
}

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
	head := fmt.Sprintf("Root directory: %v, size: %v", KindLSColorString("di", f.Root()), KindLSColorString("di", ctdsize))
	fmt.Fprintln(w, head)
	fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("=", 80))

	chead := getColorizedHead(pad, urname, gpname)
	fmt.Fprintln(w, chead)
	curname, cgpname := getColorizedUGName(urname, gpname)

	for _, dir := range dirs {
		for _, file := range fm[dir] {
			fsize := file.Size
			sfsize := bytefmt.ByteSize(fsize)
			// tsize += fsize
			cperm := getColorizePermission(file.Stat.Mode())
			cmodTime := getColorizedModTime(file.Stat.ModTime())
			cgit := getColorizedGitStatus(f.GetGitStatus(), file)
			if file.IsDir() {
				// cdsize := getColorizedSize(fsize)
				sdsize := fmt.Sprintf("%6s", bytefmt.ByteSize(fsize))
				cdsize := KindLSColorString("di", sdsize)
				if !strings.EqualFold(file.Dir, RootMark) {
					if f.depth != 0 {
						cdsize := KindLSColorString("di", sfsize)
						fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("-", 80))
						fmt.Fprintf(w, "directory: %v, size: %v\n", KindLSColorString("di", file.Dir), cdsize)
						fmt.Fprintln(w, chead)
					}
					fmt.Fprintf(w, "%s%-11s %s %s %s %14s %s %s\n", pad, cperm, cdsize, curname, cgpname, cmodTime, cgit, file.LSColorString(file.BaseName))
				}
				continue
			}
			cfsize := getColorizedSize(fsize)
			fmt.Fprintf(w, "%s%-11s %s %s %s %14s %s %s\n", pad, cperm, cfsize, curname, cgpname, cmodTime, cgit, file.LSColorString(file.BaseName))
		}
	}
	// fmt.Fprintln(w, pad)
	fmt.Fprintf(w, "%s%s\n", pad, strings.Repeat("=", 80))
	fmt.Fprintf(w, "%s%d directories, %d files, total %v.\n", pad, f.NDirs(), f.NFiles(), bytefmt.ByteSize(f.totalSize))
	return w.Bytes()
}

func getColorizedGitStatus(git GitStatus, file *File) string {
	st := "--"
	if fst, ok := git.FilesStatus[file.Path]; ok {
		fst = strings.Replace(fst, " ", "-", 1)
		fst = strings.Replace(fst, "?", "N", 1)
		fst = strings.Replace(fst, "A", "N", 1)
		return getColorizedTag(fst)
	}

	if file.IsDir() {
		gits := getGitSlice(git, file)
		if len(gits) > 0 {
			// fmt.Println("> ", file.Path, file.Dir, gits)
			return getColorizedTag(getGitTag(gits))
		}
		// else {
		// 	return st
		// }
	}
	return getColorizedTag(st)
}
func getColorizedTag(fst string) string {
	x := rune(fst[0])
	y := rune(fst[1])
	return cpmap[x].Sprint(string(x)) + cpmap[y].Sprint(string(y))
}
func getGitTag(gits []string) string {
	// paw.Logger.Info()
	x := getGitTagChar(rune(gits[0][0]))
	y := getGitTagChar(rune(gits[0][1]))
	for i := 1; i < len(gits); i++ {
		c := getGitTagChar(rune(gits[i][0]))
		if c != '-' && x != 'N' {
			x = c
		}
		c = getGitTagChar(rune(gits[i][1]))
		if c != '-' && y != 'N' {
			y = c
		}
	}
	return string(x) + string(y)
}

func getGitTagChar(c rune) rune {
	if c == '?' || c == 'A' {
		return 'N'
	}
	return c
}

func getGitSlice(git GitStatus, file *File) []string {
	gits := []string{}
	for k, v := range git.FilesStatus {
		if strings.HasPrefix(k, file.Path) {
			v = strings.Replace(v, " ", "-", 1)
			gits = append(gits, v)
		}
	}
	return gits
}

func getColorizePermission(mode os.FileMode) string {
	sperm := fmt.Sprintf("%v", mode)
	cperm := ""
	for _, p := range sperm {
		if p == 'd' {
			cperm += KindLSColorString("di", string(p))
			continue
		}
		cperm += cpmap[p].Sprint(string(p))
	}
	return cperm + " "
}

var cpmap = map[rune]*color.Color{
	// 'd': color.New(color.FgCyan).Add(color.Concealed),
	'r': color.New(color.FgYellow).Add(color.Bold),
	'w': color.New(color.FgRed).Add(color.Bold),
	'x': color.New(color.FgGreen).Add(color.Bold),
	'-': color.New(color.Concealed),
	'.': color.New(color.Concealed),
	' ': color.New(color.Concealed),                   //unmodified
	'M': color.New(color.FgBlue).Add(color.Concealed), //modified
	'A': color.New(color.FgBlue).Add(color.Concealed), //added
	'D': color.New(color.FgRed).Add(color.Concealed),  //deleted
	'R': color.New(color.FgBlue).Add(color.Concealed), //renamed
	'C': color.New(color.FgBlue).Add(color.Concealed), //copied
	'U': color.New(color.FgBlue).Add(color.Concealed), //updated but unmerged
	'?': color.New(color.FgHiGreen).Add(color.Bold),   //untracked
	'N': color.New(color.FgHiGreen).Add(color.Bold),   //untracked
	'!': color.New(color.FgBlue).Add(color.Concealed), //ignored
}

func getColorizedSize(size uint64) (csize string) {
	ssize := fmt.Sprintf("%6s", bytefmt.ByteSize(size))
	c := color.New(color.FgHiGreen).Add(color.Bold)
	csize = c.Sprint(ssize)
	return csize
}

func getColorizedUGName(urname, gpname string) (curname, cgpname string) {
	c := color.New(color.FgHiYellow).Add(color.Bold)
	curname = c.Sprint(urname)
	cgpname = c.Sprint(gpname)
	return curname, cgpname
}

func getColorizedModTime(modTime time.Time) string {
	c := color.New(color.FgBlue).Add(color.Concealed)
	s := c.Sprint(modTime.Format("01-02-06 15:04"))
	return s
}

func getColorizedHead(pad, username, groupname string) string {
	c := color.New(color.Concealed).Add(color.Underline)

	width := intmax(4, len(username))
	huser := fmt.Sprintf("%[2]*[1]s", "User", width)
	width = intmax(5, len(groupname))
	hgroup := fmt.Sprintf("%[2]*[1]s", "Group", width)

	ssize := fmt.Sprintf("%6s", "Size")
	head := fmt.Sprintf("%s%s %s %s %s %14s %s %s", pad, c.Sprint("Permissions"), c.Sprint(ssize), c.Sprint(huser), c.Sprint(hgroup), c.Sprint(" Data Modified"), c.Sprint("Git"), c.Sprint("Name"))
	return head
}

func intmax(i1, i2 int) int {
	if i1 >= i2 {
		return i1
	}
	return i2
}

// func below here, invoked from godirwalk/examples/sizes
//  `sizes()`, `sizesStack`, `newSizesStack()`, `(s *sizesStack) EnterDirectory()`, `(s *sizesStack) LeaveDirectory()`, `(s *sizesStack) Accumulate(i int64)`

func sizes(osDirname string) (uint64, error) {
	var size int64
	sizes := newSizesStack()
	return uint64(size), godirwalk.Walk(osDirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				sizes.EnterDirectory()
				return nil
			}

			st, err := os.Stat(osPathname)
			if err != nil {
				return err
			}

			size = st.Size()
			sizes.Accumulate(size)

			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			paw.Logger.Error(err)
			return godirwalk.SkipNode
		},
		PostChildrenCallback: func(osPathname string, de *godirwalk.Dirent) error {
			size = sizes.LeaveDirectory()
			sizes.Accumulate(size) // add this directory's size to parent directory.
			return nil
		},
	})
}

// sizesStack encapsulates operations on stack of directory sizes, with similar
// but slightly modified LIFO semantics to push and pop on a regular stack.
type sizesStack struct {
	sizes []int64 // stack of sizes
	top   int     // index of top of stack
}

func newSizesStack() *sizesStack {
	// Initialize with dummy value at top of stack to eliminate special cases.
	return &sizesStack{sizes: make([]int64, 1, 32)}
}

func (s *sizesStack) EnterDirectory() {
	s.sizes = append(s.sizes, 0)
	s.top++
}

func (s *sizesStack) LeaveDirectory() (i int64) {
	i, s.sizes = s.sizes[s.top], s.sizes[:s.top]
	s.top--
	return i
}

func (s *sizesStack) Accumulate(i int64) {
	s.sizes[s.top] += i
}
