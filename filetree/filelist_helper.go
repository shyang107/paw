package filetree

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"time"

	"github.com/fatih/color"
	"github.com/karrick/godirwalk"
	"github.com/shyang107/paw"
	"github.com/spf13/cast"

	// "github.com/shyang107/paw/treeprint"
	"github.com/xlab/treeprint"
)

//
// ToListTree
//
type EdgeType string

var (
	EdgeTypeLink     EdgeType = "│"   //treeprint.EdgeTypeLink
	EdgeTypeMid      EdgeType = "├──" //treeprint.EdgeTypeMid
	EdgeTypeEnd      EdgeType = "└──" //treeprint.EdgeTypeEnd
	IndentSize                = 3     //treeprint.IndentSize
	SpaceIndentSize           = paw.Repeat(" ", IndentSize)
	currentuser, _            = user.Current()
	urname                    = currentuser.Username
	usergp, _                 = user.LookupGroupId(currentuser.Gid)
	gpname                    = usergp.Name
	curname, cgpname          = getColorizedUGName(urname, gpname)
)

func printLTFile(wr io.Writer, level int, levelsEnded []int,
	edge EdgeType, fl *FileList, file *File, git GitStatus, pad string) {

	meta := pad
	if pdview == PListTreeView {
		tmeta, _ := file.ColorMeta(git)
		meta += tmeta
	}

	fmt.Fprintf(wr, "%v", meta)
	for i := 0; i < level; i++ {
		if isEnded(levelsEnded, i) {
			fmt.Fprint(wr, paw.Repeat(" ", IndentSize+1))
			continue
		}
		cedge := KindLSColorString("-", string(EdgeTypeLink))
		fmt.Fprintf(wr, "%v%s", cedge, paw.Repeat(" ", IndentSize))
	}

	cedge := KindLSColorString("-", string(edge))
	name := file.ColorBaseName()
	if file.IsDir() && fl.depth == -1 {
		dinf := fl.DirInfo(file)
		name = dinf + " " + name
	}

	fmt.Fprintf(wr, "%v %v\n", cedge, name)
}

func printLTDir(wr io.Writer, level int, levelsEnded []int,
	edge EdgeType, fl *FileList, file *File, git GitStatus, pad string) {
	fm := fl.Map()
	files := fm[file.Dir]
	nfiles := len(files)

	for i := 1; i < nfiles; i++ {
		file = files[i]
		edge := EdgeTypeMid
		if i == nfiles-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}

		printLTFile(wr, level, levelsEnded, edge, fl, file, git, pad)

		if file.IsDir() && len(fm[file.Dir]) > 1 {
			printLTDir(wr, level+1, levelsEnded, edge, fl, file, git, pad)
		}
	}
}

func isLast(file *File, fl *FileList) (isLastPreBranch, hasFiles bool) {
	dir := file.Dir
	hasFiles = len(fl.Map()[dir]) > 1

	ddir := paw.Split(dir, PathSeparator)
	pdir := paw.Join(ddir[:len(ddir)-1], PathSeparator)
	pfiles := fl.Map()[pdir]
	iplast := len(pfiles) - 1
	for i, pfile := range pfiles {
		if pfile.Path == file.Path && i == iplast {
			isLastPreBranch = true
		}
	}
	return isLastPreBranch, hasFiles
}

func isEnded(levelsEnded []int, level int) bool {
	for _, l := range levelsEnded {
		if l == level {
			return true
		}
	}
	return false
}

// GetColorizedDirName will return a colorful string of {{ dir }}/{{ name }}
func GetColorizedDirName(path string, root string) string {
	return getDirName(path, root)
}

func getDirName(path string, root string) string {
	file, err := NewFile(path)
	if err != nil {
		dir, name := filepath.Split(path)
		if len(root) > 0 {
			dir = paw.Replace(dir, root, "..", 1)
		}
		name = KindEXAColorString("dir", dir) + KindLSColorString("di", name)
		return name
	}
	name := file.LSColorString(file.BaseName)
	if file.IsDir() {
		dir, _ := filepath.Split(file.Path)
		if len(root) > 0 {
			// dir = strings.TrimPrefix(dir, root)
			dir = paw.Replace(dir, root, "..", 1)
		}
		name = KindEXAColorString("dir", dir) + name
	}
	link := checkAndGetColorLink(file)
	if len(link) > 0 {
		name += cpmap['l'].Sprint(" -> ") + link
	}
	return name
}

func getDirInfo(fl *FileList, file *File) string {
	nd, nf := 0, 0
	if file.IsDir() {
		files := fl.Map()[file.Dir]
		for _, f := range files {
			if f.IsDir() {
				nd++
			} else {
				nf++
			}
		}
	}
	di := fmt.Sprintf("%v dirs", nd-1)
	fi := fmt.Sprintf("%v files", nf)
	// return "[" + KindEXAColorString("di", di) + ", " + KindEXAColorString("dir", fi) + "]"
	return "[" + KindEXAColorString("dir", di) + ", " + KindEXAColorString("dir", fi) + "]"
	// cl := color.New(EXAColors["dir"]...).Add(color.Underline)
	// return "[" + cl.Sprint(di+", "+fi) + "]"
}

func printLTList(w io.Writer, pad string, parameters ...string) {
	str := ""
	for _, p := range parameters {
		str += fmt.Sprintf("%v ", p)
	}
	fmt.Fprintf(w, "%v%v", pad, str)
}

//
// ToTree
//

func getNDirsFiles(files []*File) (ndirs, nfiles int) {
	for _, file := range files {
		if file.IsDir() {
			ndirs++
		} else {
			nfiles++
		}
	}
	return ndirs - 1, nfiles
}

func padding(pad string, bytes []byte) []byte {
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
	dd := paw.Split(dir, PathSeparator)
	nd := len(dd)
	var pre treeprint.Tree
	// fmt.Println(dir, nd)
	if nd == 2 { // ./xx
		pre = tree
	} else { //./xx/...
		pre = tree
		for i := 2; i < nd; i++ {
			predir := paw.Join(dd[:i], PathSeparator)
			// fmt.Println("\t", i, predir)
			f := fm[predir][0] // import dir
			pre = pre.FindByValue(f)
		}
	}
	return pre
}

//
// Tolist
//

func printDirSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize uint64) {
	msg := KindLSColorString("-", fmt.Sprintf("%s%v directories; %v files, size ≈ %v.\n", pad, ndirs, nfiles, ByteSize(sumsize)))
	fmt.Fprintf(w, msg)
}

func printTotalSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize uint64) {
	fmt.Fprintf(w, "%s\n%sAccumulated %v directories, %v files, total size ≈ %v.\n", pad, pad, ndirs, nfiles, ByteSize(sumsize))
}

func printFileItem(w io.Writer, pad string, parameters ...string) {
	str := ""
	for _, p := range parameters {
		str += fmt.Sprintf("%v ", p)
	}
	str += "\n"
	fmt.Fprintf(w, "%v%v", pad, str)
	// fmt.Fprintf(w, "%s%s %s %s %s %s %s %s\n", pad, cperm, cfsize, curname, cgpname, cmodTime, cgit, name)
}

var ckxy = map[rune]rune{
	'M': 'M', //modified
	// "A": "A", //added
	'A': 'N',
	'D': 'D', //deleted
	'R': 'R', //renamed
	'C': 'C', //copied
	'U': 'U', //updated but unmerged
	'?': 'N', //untracked
	' ': '-',
	'!': 'N', //ignore
}

// getColorizedGitStatus will return a colorful string of shrot status of git.
// The length of placeholder in terminal is 3.
func getColorizedGitStatus(git GitStatus, file *File) string {
	x, y := '-', '-'

	xy, ok := git.FilesStatus[file.Path]

	if ok {
		x, y = xy.Split()
		x = ckxy[x]
		y = ckxy[y]
	}

	if file.IsDir() {
		// gits := getGitSlice(git, file)
		for k, v := range git.FilesStatus {
			if paw.HasPrefix(k, file.Path) {
				vx, vy := v.Split()
				cx := ckxy[vx]
				cy := ckxy[vy]
				if cx != '-' && x != 'N' {
					x = cx
				}
				if cy != '-' && y != 'N' {
					y = cy
				}
			}
		}
	}

	var sx, sy string
	if x == 'N' && y == 'N' {
		sx, sy = "-", "N"
	} else {
		sx, sy = string(x), string(y)
	}

	return " " + cpmap[x].Sprint(sx) + cpmap[y].Sprint(sy)
}

// GetColorizePermission will return a colorful string of mode
// The length of placeholder in terminal is 10.
func GetColorizePermission(mode os.FileMode) string {
	return getColorizePermission(mode)
}
func getColorizePermission(mode os.FileMode) string {
	sperm := fmt.Sprintf("%v", mode)
	c := ""
	// fmt.Println(len(s))
	for i := 0; i < len(sperm); i++ {
		s := string(sperm[i])
		cs := s
		if cs != "-" {
			switch i {
			case 0:
				switch s {
				case "d":
					cs = "di"
				case "L":
					cs = "ln"
				}
			case 1, 2, 3:
				cs = "u" + s
			case 4, 5, 6:
				cs = "g" + s
			case 7, 8, 9:
				cs = "t" + s
			}
		}
		if i == 0 && cs == "-" {
			s = "."
		}
		// c += color.New(EXAColors[cs]...).Add(color.Bold).Sprint(s)
		c += NewEXAColor(cs).Add(color.Bold).Sprint(s)
	}

	return c
}

var cpmap = map[rune]*color.Color{
	'L': NewLSColor("ln").Add(color.Concealed),
	'l': NewLSColor("ln").Add(color.Concealed),
	'd': NewLSColor("di").Add(color.Concealed),
	'r': NewEXAColor("ur").Add(color.Bold),
	'w': NewEXAColor("uw").Add(color.Bold),
	'x': NewEXAColor("ux").Add(color.Bold),
	'-': NewEXAColor("-"),  //color.New(color.Concealed),
	'.': NewEXAColor("."),  //color.New(color.Concealed),
	' ': NewEXAColor(" "),  //color.New(color.Concealed), //unmodified
	'M': NewEXAColor("gm"), //color.New(EXAColors["gm"]...), //modified
	'A': NewEXAColor("ga"), //color.New(EXAColors["ga"]...), //added
	'D': NewEXAColor("gd"), //color.New(EXAColors["gd"]...), //deleted
	'R': NewEXAColor("gv"), //color.New(EXAColors["gv"]...), //renamed
	'C': NewEXAColor("gt"), //color.New(EXAColors["gt"]...), //copied
	'U': NewEXAColor("gt"), //color.New(EXAColors["gt"]...), //updated but unmerged
	'?': NewEXAColor("gm"), //color.New(EXAColors["gm"]...), //untracked
	'N': NewEXAColor("ga"), //color.New(EXAColors["ga"]...), //untracked
	'!': NewEXAColor("-"),  //color.New(EXAColors["-"]...),  //ignored
}

// GetColorizedSize will return a humman-readable and colorful string of size.
// The length of placeholder in terminal is 6.
func GetColorizedSize(size uint64) string {
	return getColorizedSize(size)
}
func getColorizedSize(size uint64) (csize string) {
	ss := ByteSize(size)
	nss := len(ss)
	sn := fmt.Sprintf("%5s", ss[:nss-1])
	su := paw.ToLower(ss[nss-1:])
	cn := NewEXAColor("sn").Add(color.Bold)
	cu := NewEXAColor("sb")
	csize = cn.Sprint(sn) + cu.Sprint(su)
	return csize
}

func getColorizedUGName(urname, gpname string) (curname, cgpname string) {
	cu := NewEXAColor("uu").Add(color.Bold)
	cg := NewEXAColor("gu").Add(color.Bold)
	curname = cu.Sprint(urname)
	cgpname = cg.Sprint(gpname)
	return curname, cgpname
}

// GetColorizedTime will return a colorful string of time.
// The length of placeholder in terminal is 14.
func GetColorizedTime(modTime time.Time) string {
	return getColorizedModTime(modTime)
}
func getColorizedModTime(modTime time.Time) string {
	return NewEXAColor("da").Sprint(modTime.Format("01-02-06 15:04"))
}

func getColorizedHead(pad, username, groupname string, git GitStatus) string {
	c := NewEXAColor("hd").Add(color.Underline)
	width := max(4, len(username))
	huser := fmt.Sprintf("%[2]*[1]s", "User", width)
	width = max(5, len(groupname))
	hgroup := fmt.Sprintf("%[2]*[1]s", "Group", width)

	ssize := fmt.Sprintf("%6s", "Size")
	head := ""
	if git.NoGit {
		head = fmt.Sprintf("%s%s %s %s %s %14s %s", pad, c.Sprint("Permissions"), c.Sprint(ssize), c.Sprint(huser), c.Sprint(hgroup), c.Sprint(" Data Modified"), c.Sprint("Name"))
	} else {
		head = fmt.Sprintf("%s%s %s %s %s %14s %s %s", pad, c.Sprint("Permissions"), c.Sprint(ssize), c.Sprint(huser), c.Sprint(hgroup), c.Sprint(" Data Modified"), c.Sprint("Git"), c.Sprint("Name"))
	}
	return head
}

func max(i1, i2 int) int {
	if i1 >= i2 {
		return i1
	}
	return i2
}

func printBanner(w io.Writer, pad string, mark string, length int) {
	fmt.Fprintf(w, "%s%s\n", pad, paw.Repeat(mark, length))
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

func getTerminalSize() (height, width int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		// log.Fatal(err)
		paw.Error.Println("run stty, err: ", err)
		return 38, 100
	}
	size := paw.Split(paw.TrimSuffix(string(out), "\n"), " ")
	height = cast.ToInt(size[0])
	width = cast.ToInt(size[1])
	return height, width
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
