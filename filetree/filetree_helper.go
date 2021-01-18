package filetree

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"time"

	"github.com/fatih/color"
	"github.com/karrick/godirwalk"
	"github.com/shyang107/paw"
	// "github.com/shyang107/paw/treeprint"
)

//
// ToListTree
//
type EdgeType string

const (
	EdgeTypeLink     EdgeType = "│"   //treeprint.EdgeTypeLink
	EdgeTypeMid      EdgeType = "├──" //treeprint.EdgeTypeMid
	EdgeTypeEnd      EdgeType = "└──" //treeprint.EdgeTypeEnd
	IndentSize                = 3     //treeprint.IndentSize
	dateLayout                = "Jan 02, 2006"
	timeThisLayout            = "01-02 15:04"
	timeBeforeLayout          = "2006-01-02"
)

var (
	edgeWidth map[EdgeType]int = map[EdgeType]int{
		EdgeTypeLink: 1,
		EdgeTypeMid:  3,
		EdgeTypeEnd:  3,
	}
	now                   = time.Now()
	thisYear              = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	SpaceIndentSize       = paw.Spaces(IndentSize)
	cdashp                = NewEXAColor("-")
	cxp                   = NewEXAColor("xattr")
	chdp                  = NewEXAColor("hd")
	cdirp                 = NewEXAColor("dir")
	lsdip                 = NewLSColor("di")
	cdip                  = NewEXAColor("di")
	cfip                  = NewEXAColor("fi")
	cuup                  = NewEXAColor("uu")
	cgup                  = NewEXAColor("gu")
	cinp                  = NewEXAColor("in")
	clkp                  = NewEXAColor("lk")
	cbkp                  = NewEXAColor("bk")
	cdap                  = NewEXAColor("da")
	currentuser, _        = user.Current()
	urname                = currentuser.Username
	usergp, _             = user.LookupGroupId(currentuser.Gid)
	gpname                = usergp.Name
	curname               = cuup.Sprint(urname)
	cgpname               = cgup.Sprint(gpname)
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
			dir = paw.Replace(dir, root, RootMark, 1)
		}
	}
	if file.IsLink() {
		return dir + name, file.LinkPath()
	}
	return dir, name
}

func getDirInfo(fl *FileList, file *File) (cdinf string, wdinf int) {
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
	wdinf = len(di) + len(fi) + 4
	cdinf = fmt.Sprintf("[%s, %s]", cdirp.Sprint(di), cdirp.Sprint(fi))
	return cdinf, wdinf
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
		x, y = ckxy[x], ckxy[y]
	}

	if file.IsDir() {
		// gits := getGitSlice(git, file)
		for k, v := range git.FilesStatus {
			if paw.HasPrefix(k, file.Path) {
				vx, vy := v.Split()
				cx, cy := ckxy[vx], ckxy[vy]
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

// getGitStatus will return a  string of shrot status of git.
// The length of placeholder in terminal is 3.
func getGitStatus(git GitStatus, file *File) string {
	x, y := '-', '-'

	xy, ok := git.FilesStatus[file.Path]
	if ok {
		x, y = xy.Split()
		x, y = ckxy[x], ckxy[y]
	}

	if file.IsDir() {
		// gits := getGitSlice(git, file)
		for k, v := range git.FilesStatus {
			if paw.HasPrefix(k, file.Path) {
				vx, vy := v.Split()
				cx, cy := ckxy[vx], ckxy[vy]
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

	return " " + sx + sy
}

// GetColorizePermission will return a colorful string of mode
// The length of placeholder in terminal is 10.
func GetColorizePermission(mode os.FileMode) string {
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
func GetColorizedSize(size uint64) (csize string) {
	ss := ByteSize(size)
	nss := len(ss)
	sn := fmt.Sprintf("%5s", ss[:nss-1])
	su := paw.ToLower(ss[nss-1:])
	cn := NewEXAColor("sn")
	cu := NewEXAColor("sb")
	csize = cn.Sprint(sn) + cu.Sprint(su)
	return csize
}

func getColorizedDates(file *File) (cdate string, wd int) {
	dates := make([]string, len(pfieldKeys))
	for _, k := range pfieldKeys {
		cdate := ""
		var date time.Time
		switch k {
		case PFieldModified:
			cdate = file.ColorModifyTime()
			date = file.ModifiedTime()
			// sv = DateString(date)
		case PFieldAccessed:
			cdate = file.ColorAccessedTime()
			date = file.AccessedTime()
			// sv = DateString(date)
		case PFieldCreated:
			cdate = file.ColorCreatedTime()
			date = file.CreatedTime()
			// sv = DateString(date)
		default:
			continue
		}
		fwd := pfieldWidthsMap[k]
		fwdd := len(DateString(date))
		sp := paw.Spaces(fwd - fwdd + 1)
		dates = append(dates, cdate+sp)
		// wd += paw.StringWidth(sv)
		wd += fwd + 1
	}
	cdate = paw.Join(dates, "")
	cdate = cdate[:len(cdate)-1]
	wd--
	return cdate, wd
}

func getColorizedHead(pad, username, groupname string, git GitStatus) (chead string, width int) {
	sb := paw.NewStringBuilder()
	csb := paw.NewStringBuilder()
	sb.WriteString(pad)
	csb.WriteString(pad)
	for i, k := range pfieldKeys {
		switch k {
		// case PFieldINode: //"inode",
		// 	// field = fmt.Sprintf("%-[1]*[2]s", pfieldWidthsMap[k], fieldsMap[k])
		// case PFieldPermissions: //"Permissions",
		// case PFieldLinks: //"Links",
		// case PFieldSize: //"Size",
		// case PFieldUser: //"User",
		// case PFieldGroup: //"Group",
		// case PFieldModified: //"Date Modified",
		// case PFieldCreated: //"Date Created",
		// case PFieldAccessed: //"Date Accessed",
		case PFieldGit: //"Git",
			if git.NoGit {
				pfieldKeys = append(pfieldKeys[:i], pfieldKeys[i+1:]...)
				pfields = append(pfields[:i], pfields[i+1:]...)
				continue
			}
			// case PFieldName: //"Name",
		}
		field := fmt.Sprintf("%[1]*[2]s", pfieldWidthsMap[k], pfieldsMap[k])
		fmt.Fprintf(sb, "%s ", field)
		fmt.Fprintf(csb, "%s ", chdp.Sprint(field))
	}
	head := sb.String()
	head = head[:len(head)-1]
	width = paw.StringWidth(head)
	chead = csb.String()
	chead = chead[:len(chead)-1]
	return chead, width
	// width := paw.MaxInts(4, paw.StringWidth(username))
	// huser := fmt.Sprintf("%[1]*[2]s", width, "User")
	// width = paw.MaxInts(5, paw.StringWidth(groupname))
	// hgroup := fmt.Sprintf("%[1]*[2]s", width, "Group")

	// ssize := fmt.Sprintf("%6s", "Size")

	// dates := make([]string, len(fieldKeys))
	// for _, v := range fields {
	// 	dates = append(dates, chdp.Sprint(v)+" ")
	// }
	// cdate := paw.Join(dates, "")
	// cdate = paw.TrimSpace(cdate)

	// head := ""
	// if git.NoGit {
	// 	// head = fmt.Sprintf("%s%s %s %s %s %14s %s", pad, chdp.Sprint("Permissions"), chdp.Sprint(ssize), chdp.Sprint(huser), chdp.Sprint(hgroup), chdp.Sprint(" Date Modified"), chdp.Sprint("Name"))
	// 	head = fmt.Sprintf("%s%s %s %s %s %s %s", pad, chdp.Sprint("Permissions"), chdp.Sprint(ssize), chdp.Sprint(huser), chdp.Sprint(hgroup), cdate, chdp.Sprint("Name"))
	// } else {
	// 	// head = fmt.Sprintf("%s%s %s %s %s %14s %s %s", pad, chdp.Sprint("Permissions"), chdp.Sprint(ssize), chdp.Sprint(huser), chdp.Sprint(hgroup), chdp.Sprint(" Date Modified"), chdp.Sprint("Git"), chdp.Sprint("Name"))
	// 	head = fmt.Sprintf("%s%s %s %s %s %s %s %s", pad, chdp.Sprint("Permissions"), chdp.Sprint(ssize), chdp.Sprint(huser), chdp.Sprint(hgroup), cdate, chdp.Sprint("Git"), chdp.Sprint("Name"))
	// }
	// return head
}

func printBanner(w io.Writer, pad string, mark string, length int) {
	banner := fmt.Sprintf("%s%s", pad, paw.Repeat(mark, length))
	fmt.Fprintln(w, cdashp.Sprint(banner))
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

func printListln(w io.Writer, items ...interface{}) {
	sb := paw.NewStringBuilder()
	sb.Grow(len(items))
	for i := 0; i < len(items); i++ {
		if i < len(items)-1 {
			fmt.Fprintf(sb, "%v ", items[i])
		} else {
			fmt.Fprintf(sb, "%v", items[i])
		}
	}
	fmt.Fprintln(w, sb.String())
}

func printFileName(w io.Writer, pad string, fds *FieldSlice, fdNo *Field, file *File) {
}

func DateString(date time.Time) (sdate string) {
	sdate = date.Format(timeThisLayout)
	if date.Before(thisYear) {
		sdate = date.Format(timeBeforeLayout)
	}
	return sdate
}

// GetColorizedTime will return a colorful string of time.
// The length of placeholder in terminal is 14.
func GetColorizedTime(date time.Time) string {
	return cdap.Sprint(DateString(date))
}

func wrapFileName(file *File, fds *FieldSlice, pad string, wdsttylimit int) string {
	var (
		sb     = paw.NewStringBuilder()
		wpad   = paw.StringWidth(pad)
		meta   = fds.ColorMetaValuesString()
		wmeta  = fds.MetaValuesStringWidth()
		spmeta = paw.Spaces(wmeta)
		name   = file.BaseNameToLink()
		wname  = paw.StringWidth(name)
		wdstty = wdsttylimit - 1
		limit  = wdstty - wpad - wmeta
	)
	if wname <= limit {
		printListln(sb, pad+meta, file.ColorName())
	} else { // wrap file name
		if err := checkStringRange(name, limit, "Name"); err != nil {
			paw.Error.Fatal(err, " (may be too many fields)")
		}
		if !file.IsLink() {
			names := paw.WrapToSlice(name, limit)
			printListln(sb, pad+meta, file.LSColorString(names[0]))
			for i := 1; i < len(names); i++ {
				printListln(sb, pad+spmeta, file.LSColorString(names[i]))
			}
		} else {
			cname := file.LSColorString(file.BaseName)
			wbname := paw.StringWidth(file.BaseName)
			carrow := cdashp.Sprint(" -> ")
			wbname += 4
			printListln(sb, pad+meta, cname+carrow)
			dir, name := filepath.Split(file.LinkPath())
			wd, wn := paw.StringWidth(dir), paw.StringWidth(name)

			if wd+wn <= limit {
				printListln(sb, pad+spmeta, cdirp.Sprint(dir)+cdip.Sprint(name))
			} else {
				if wd <= limit {
					clink := cdirp.Sprint(dir) + cdip.Sprint(name[:limit-wd])
					printListln(sb, pad+spmeta, clink)
					names := paw.WrapToSlice(name[limit-wd:], limit)
					for _, v := range names {
						clink = cdip.Sprint(v)
						printListln(sb, pad+spmeta, clink)
					}
				} else { // wd > limit
					dirs := paw.WrapToSlice(dir, limit)
					nd := len(dirs)
					var clink string
					for i := 0; i < nd-1; i++ {
						clink = cdirp.Sprint(dirs[i])
						printListln(sb, pad+spmeta, clink)
					}
					clink = cdirp.Sprint(dirs[nd-1])
					wdLast := paw.StringWidth(dirs[nd-1])
					if wn <= limit-wdLast {
						clink += cdip.Sprint(name)
						printListln(sb, pad+spmeta, clink)
					} else { // wn > wd-limit
						clink += cdip.Sprint(name[:limit-wdLast])
						printListln(sb, pad+spmeta, clink)
						rname := name[limit-wdLast:]
						wr := paw.StringWidth(rname)
						if wr <= limit {
							clink = cdip.Sprint(rname)
							printListln(sb, pad+spmeta, clink)
						} else { // wr > limit
							names := paw.WrapToSlice(rname, limit)
							for _, v := range names {
								clink = cdip.Sprint(v)
								printListln(sb, pad+spmeta, clink)
							}
						}
					}
				}
			}
		}
	}

	return sb.String()
}

func xattrEdgeString(file *File, pad string, wmeta int, wdsttylimit int) string {
	var (
		nx = len(file.XAttributes)
		sb = paw.NewStringBuilder()
	)
	if nx > 0 {
		var edge = EdgeTypeMid
		for i := 0; i < nx; i++ {
			var wdm = wmeta
			xattr := file.XAttributes[i]
			if i == nx-1 {
				edge = EdgeTypeEnd
			}
			var padx = fmt.Sprintf("%s %s ", pad, cdashp.Sprint(edge))
			wdm += edgeWidth[edge] + 2
			var wdx = len(xattr)
			if wdm+wdx <= wdsttylimit {
				printListln(sb, padx+cxp.Sprint(xattr))
			} else {
				var wde = wdsttylimit - wdm
				if err := checkStringRange(xattr, wde, "xattr"); err != nil {
					paw.Error.Fatal(err, " (may be too many fields)")
				}
				printListln(sb, padx+cxp.Sprint(xattr[:wde]))
				switch edge {
				case EdgeTypeMid:
					padx = fmt.Sprintf("%s %s ", pad, cdashp.Sprint(EdgeTypeLink)+SpaceIndentSize)
				case EdgeTypeEnd:
					padx = fmt.Sprintf("%s %s ", pad, paw.Spaces(edgeWidth[edge]))
				}
				if err := checkStringRange(xattr, wde, "xattr"); err != nil {
					paw.Error.Fatal(err, "may be too many fields")
				}
				if len(xattr[wde:]) <= wde {
					printListln(sb, padx+cxp.Sprint(xattr[wde:]))
				} else {
					xattrs := paw.Split(paw.Wrap(xattr[wde:], wde), "\n")
					for _, v := range xattrs {
						printListln(sb, padx+cxp.Sprint(v))
					}
				}
			}
		}
	}
	return sb.String()
}

func checkStringRange(s string, idx int, name string) error {
	var ns = len(s)
	if idx < 0 || idx > ns-1 {
		return errors.New(fmt.Sprintf("slice (%s) bounds out of range,len(slice) = %d, index: %d", name, ns, idx))
	}
	return nil
}

func getMaxFileSizeWidth(files []*File) int {
	var (
		wdsize = 0
	)
	for _, f := range files {
		var size = ByteSize(f.Size)
		if wdsize < len(size) {
			wdsize = len(size)
		}
	}
	return wdsize
}

func modifyHead(fds *FieldSlice, files []*File, pad string) (chead string, wdmeta int) {
	wdsize := getMaxFileSizeWidth(files)
	fds.Get(PFieldSize).Width = paw.MaxInt(wdsize, fds.Get(PFieldSize).Width)
	chead = fds.ColorHeadsString()
	wdmeta = fds.MetaHeadsStringWidth() + paw.StringWidth(pad)
	return chead, wdmeta
}

func modifyTreeHead(fds *FieldSlice, fl *FileList, pad string) (chead string, wdmeta int) {
	wdsize := 0
	for _, dir := range fl.Dirs() {
		files := fl.Map()[dir][1:]
		wd := getMaxFileSizeWidth(files)
		if wdsize < wd {
			wdsize = wd
		}
	}
	fds.Get(PFieldSize).Width = paw.MaxInt(wdsize, fds.Get(PFieldSize).Width)
	chead = fds.ColorHeadsString()
	wdmeta = fds.MetaHeadsStringWidth() + paw.StringWidth(pad)
	return chead, wdmeta
}

func modifySizeWidth(fds *FieldSlice, fl *FileList) {
	wdsize := 0
	for _, dir := range fl.Dirs() {
		files := fl.Map()[dir][1:]
		wd := getMaxFileSizeWidth(files)
		if wdsize < wd {
			wdsize = wd
		}
	}
	fds.Get(PFieldSize).Width = paw.MaxInt(wdsize, fds.Get(PFieldSize).Width)
}
