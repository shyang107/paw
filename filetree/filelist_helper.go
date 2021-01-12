package filetree

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"time"

	"github.com/fatih/color"
	"github.com/karrick/godirwalk"
	"github.com/shyang107/paw"

	// "github.com/shyang107/paw/treeprint"
	"github.com/xlab/treeprint"
)

//
// ToListTree
//
type EdgeType string

const (
	EdgeTypeLink EdgeType = "│"   //treeprint.EdgeTypeLink
	EdgeTypeMid  EdgeType = "├──" //treeprint.EdgeTypeMid
	EdgeTypeEnd  EdgeType = "└──" //treeprint.EdgeTypeEnd
	IndentSize            = 3     //treeprint.IndentSize
	dateLayout            = "Jan 02, 2006"
	timeLayout            = "01-02-06 15:04"
)

var (
	edgeWidth map[EdgeType]int = map[EdgeType]int{
		EdgeTypeLink: 1,
		EdgeTypeMid:  3,
		EdgeTypeEnd:  3,
	}
	SpaceIndentSize       = paw.Spaces(IndentSize)
	cdashp                = NewEXAColor("-")
	cxp                   = NewEXAColor("xattr")
	chdp                  = NewEXAColor("hd")
	cdirp                 = NewEXAColor("dir")
	lsdip                 = NewLSColor("di")
	currentuser, _        = user.Current()
	urname                = currentuser.Username
	usergp, _             = user.LookupGroupId(currentuser.Gid)
	gpname                = usergp.Name
	curname, cgpname      = getColorizedUGName(urname, gpname)
	sttyHeight, sttyWidth = paw.GetTerminalSize()
)

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
	return getColorDirName(path, root)
}

func getColorDirName(path string, root string) string {
	file, err := NewFile(path)
	if err != nil {
		dir, name := filepath.Split(path)
		if len(root) > 0 {
			dir = paw.Replace(dir, root, RootMark, 1)
		}
		name = cdirp.Sprint(dir) + lsdip.Sprint(name)
		return name
	}
	name := file.LSColorString(file.BaseName)
	if file.IsDir() {
		dir, _ := filepath.Split(file.Path)
		if len(root) > 0 {
			// dir = strings.TrimPrefix(dir, root)
			dir = paw.Replace(dir, root, RootMark, 1)
		}
		name = cdirp.Sprint(dir) + name
	}
	if file.IsLink() {
		name += cdashp.Sprint(" -> ") + file.ColorLinkPath()
	}
	return name
}

func getDirAndName(path string, root string) (dir, name string) {
	file, err := NewFile(path)
	if err != nil {
		dir, name = filepath.Split(path)
		if len(root) > 0 {
			dir = paw.Replace(dir, root, RootMark, 1)
		}
		return dir, name
	}
	name = file.BaseName
	if file.IsDir() {
		dir, _ = filepath.Split(file.Path)
		if len(root) > 0 {
			// dir = strings.TrimPrefix(dir, root)
			dir = paw.Replace(dir, root, RootMark, 1)
		}
	}
	if file.IsLink() {
		return dir + name, file.LinkPath()
	}
	return dir, name
}

func getDirInfo(fl *FileList, file *File) (cdinf string, ndinf int) {
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
	ndinf = len(di) + len(fi) + 4
	cdinf = fmt.Sprintf("[%s, %s]", cdirp.Sprint(di), cdirp.Sprint(fi))
	return cdinf, ndinf
	// return "[" + cl.Sprint(di+", "+fi) + "]"
}

func printLTList(w io.Writer, pad string, parameters ...string) {
	sb := new(strings.Builder)
	sb.Grow(len(parameters))
	for _, p := range parameters {
		sb.WriteString(fmt.Sprintf("%v ", p))
	}
	fmt.Fprintf(w, "%v%v", pad, sb.String())
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
	// msg := KindLSColorString("-", fmt.Sprintf("%s%v directories; %v files, size ≈ %v.\n", pad, ndirs, nfiles, ByteSize(sumsize)))
	msg := fmt.Sprintf("%s%v directories; %v files, size ≈ %v.\n", pad, ndirs, nfiles, ByteSize(sumsize))
	fmt.Fprintf(w, cdashp.Sprint(msg))
}

func printTotalSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize uint64) {
	// fmt.Fprintf(w, "%s\n%sAccumulated %v directories, %v files, total size ≈ %v.\n", pad, pad, ndirs, nfiles, ByteSize(sumsize))
	summary := fmt.Sprintf("%sAccumulated %v directories, %v files, total size ≈ %v.\n", pad, ndirs, nfiles, ByteSize(sumsize))
	fmt.Fprintf(w, cdashp.Sprint(summary))
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
		c += NewEXAColor(cs).Sprint(s)
	}

	return c
}

var cpmap = map[rune]*color.Color{
	'L': NewEXAColor("ln"),
	'l': NewEXAColor("ln"),
	'd': NewEXAColor("di"),
	'r': NewEXAColor("ur"),
	'w': NewEXAColor("uw"),
	'x': NewEXAColor("ux"),
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
	cn := NewEXAColor("sn")
	cu := NewEXAColor("sb")
	csize = cn.Sprint(sn) + cu.Sprint(su)
	return csize
}

func getColorizedUGName(urname, gpname string) (curname, cgpname string) {
	cu := NewEXAColor("uu")
	cg := NewEXAColor("gu")
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
	return NewEXAColor("da").Sprint(modTime.Format(timeLayout))
}

func getColorizedHead(pad, username, groupname string, git GitStatus) string {
	width := paw.MaxInt(4, len(username))
	huser := fmt.Sprintf("%[2]*[1]s", "User", width)
	width = paw.MaxInt(5, len(groupname))
	hgroup := fmt.Sprintf("%[2]*[1]s", "Group", width)

	ssize := fmt.Sprintf("%6s", "Size")
	head := ""
	if git.NoGit {
		head = fmt.Sprintf("%s%s %s %s %s %14s %s", pad, chdp.Sprint("Permissions"), chdp.Sprint(ssize), chdp.Sprint(huser), chdp.Sprint(hgroup), chdp.Sprint(" Data Modified"), chdp.Sprint("Name"))
	} else {
		head = fmt.Sprintf("%s%s %s %s %s %14s %s %s", pad, chdp.Sprint("Permissions"), chdp.Sprint(ssize), chdp.Sprint(huser), chdp.Sprint(hgroup), chdp.Sprint(" Data Modified"), chdp.Sprint("Git"), chdp.Sprint("Name"))
	}
	return head
}

func printBanner(w io.Writer, pad string, mark string, length int) {
	banner := fmt.Sprintf("%s%s\n", pad, paw.Repeat(mark, length))
	fmt.Fprintf(w, cdashp.Sprint(banner))
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
