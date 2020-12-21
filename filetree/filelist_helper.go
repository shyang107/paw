package filetree

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/fatih/color"
	"github.com/karrick/godirwalk"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/treeprint"
)

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

//
// ToText
//

func printDirSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize uint64) {
	fmt.Fprintf(w, "%s%v directories, %v files, size: %v.\n", pad, ndirs, nfiles, bytefmt.ByteSize(sumsize))
}

func printFileItem(w io.Writer, pad string, parameters ...string) {
	str := ""
	for _, p := range parameters {
		str += fmt.Sprintf("%v ", p)
	}
	str += "\n"
	fmt.Fprintf(w, "%s%s", pad, str)
	// fmt.Fprintf(w, "%s%s %s %s %s %s %s %s\n", pad, cperm, cfsize, curname, cgpname, cmodTime, cgit, name)
}

func getColorizedGitStatus(git GitStatus, file *File) string {
	st := "--"
	xy, ok := git.FilesStatus[file.Path]

	if ok {
		xy = checkXY(xy)
		st = xy.String()
	}

	if file.IsDir() {
		gits := getGitSlice(git, file)
		if len(gits) > 0 {
			st = getGitTag(gits)
		}
	}
	return getColorizedTag(st)
}

func checkXY(xy XY) XY {
	st := xy.String()
	st = strings.Replace(st, " ", "-", -1)
	st = strings.Replace(st, "??", "-N", -1)
	st = strings.Replace(st, "?", "N", -1)
	st = strings.Replace(st, "A", "N", -1)
	return ToXY(st)
}

func getColorizedTag(fst string) string {
	x := rune(fst[0])
	y := rune(fst[1])
	return " " + cpmap[x].Sprint(string(x)) + cpmap[y].Sprint(string(y))
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
			xy := checkXY(v)
			gits = append(gits, xy.String())
		}
	}
	return gits
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
		c += color.New(EXAColors[cs]...).Add(color.Bold).Sprint(s)
	}

	return c + " "
}

var cpmap = map[rune]*color.Color{
	'L': color.New(LSColors["ln"]...).Add(color.Concealed),
	'l': color.New(LSColors["ln"]...).Add(color.Concealed),
	'd': color.New(LSColors["di"]...).Add(color.Concealed),
	'r': color.New(color.FgYellow).Add(color.Bold),
	'w': color.New(color.FgRed).Add(color.Bold),
	'x': color.New([]color.Attribute{38, 5, 155}...).Add(color.Bold),
	'-': color.New(color.Concealed),
	'.': color.New(color.Concealed),
	' ': color.New(color.Concealed), //unmodified
	// 'M': color.New(color.FgBlue).Add(color.Concealed), //modified
	// 'A': color.New(color.FgBlue).Add(color.Concealed), //added
	// 'D': color.New(color.FgRed).Add(color.Concealed),  //deleted
	// 'R': color.New(color.FgBlue).Add(color.Concealed), //renamed
	// 'C': color.New(color.FgBlue).Add(color.Concealed), //copied
	// 'U': color.New(color.FgBlue).Add(color.Concealed), //updated but unmerged
	// '?': color.New(color.FgHiGreen).Add(color.Bold),   //untracked
	// 'N': color.New(color.FgHiGreen).Add(color.Bold),   //untracked
	// '!': color.New(color.FgBlue).Add(color.Concealed), //ignored
	'M': color.New(EXAColors["gm"]...), //modified
	'A': color.New(EXAColors["ga"]...), //added
	'D': color.New(EXAColors["gd"]...), //deleted
	'R': color.New(EXAColors["gv"]...), //renamed
	'C': color.New(EXAColors["gt"]...), //copied
	'U': color.New(EXAColors["gt"]...), //updated but unmerged
	'?': color.New(EXAColors["gm"]...), //untracked
	'N': color.New(EXAColors["ga"]...), //untracked
	'!': color.New(EXAColors["-"]...),  //ignored
}

func getColorizedSize(size uint64) (csize string) {
	ss := bytefmt.ByteSize(size)
	nss := len(ss)
	sn := fmt.Sprintf("%5s", ss[:nss-1])
	su := strings.ToLower(ss[nss-1:])
	// c := color.New(color.FgHiGreen).Add(color.Bold)
	cn := color.New(EXAColors["sn"]...).Add(color.Bold)
	cu := color.New(EXAColors["sb"]...)
	csize = cn.Sprint(sn) + cu.Sprint(su)
	return csize
}

func getColorizedUGName(urname, gpname string) (curname, cgpname string) {
	cu := color.New(EXAColors["uu"]...).Add(color.Bold)
	cg := color.New(EXAColors["gu"]...).Add(color.Bold)
	curname = cu.Sprint(urname)
	cgpname = cg.Sprint(gpname)
	return curname, cgpname
}

func getColorizedModTime(modTime time.Time) string {
	c := color.New(EXAColors["da"]...)
	s := c.Sprint(modTime.Format("01-02-06 15:04"))
	return s
}

func getColorizedHead(pad, username, groupname string) string {
	c := color.New(EXAColors["hd"]...).Add(color.Underline)

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
