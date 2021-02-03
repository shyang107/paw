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
	chdp                  = paw.Chdp  // head
	cdirp                 = paw.Cdirp // pre-dir part of path
	cdip                  = paw.Cdip  // directory
	cfip                  = paw.Cfip  // file
	cNop                  = paw.CNop  // serial number
	cinp                  = paw.Cinp  // inode
	cpmp                  = paw.Cpmp  // permission
	csnp                  = paw.Csnp  // size number
	csup                  = paw.Csup  // size unit
	cuup                  = paw.Cuup  // user
	cgup                  = paw.Cgup  // group
	cunp                  = paw.Cunp  // user is not you
	cgnp                  = paw.Cgnp  // group without you
	clkp                  = paw.Clkp  // symlink
	cbkp                  = paw.Cbkp  // blocks
	cdap                  = paw.Cdap  // date
	cgitp                 = paw.Cgitp // git
	cxap                  = paw.Cxap  // extended attributes
	cxbp                  = paw.Cxbp  // extended attributes
	cdashp                = paw.Cdashp
	cnop                  = paw.CNop // no this file kind
	cbdp                  = paw.Cbdp // device
	ccdp                  = paw.Ccdp // CharDevice
	cpip                  = paw.Cpip // named pipe
	csop                  = paw.Csop // socket
	cexp                  = paw.Cexp // execution
	clnp                  = paw.Clnp // symlink
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

// // GetColorizedDirName will return a colorful string of {{ dir }}/{{ name }}
// func GetColorizedDirName(path string, root string) string {
// 	return GetColorizedPath(path, root)

// 	if path == PathSeparator {
// 		return cdip.Sprint(path)
// 	}

// 	file, err := NewFile(path)
// 	if err != nil {
// 		dir, name := filepath.Split(path)
// 		if len(root) > 0 {
// 			dir = strings.Replace(dir, root, RootMark, 1)
// 		}
// 		name = cdirp.Sprint(dir) + cdip.Sprint(name)
// 		return name
// 	}
// 	name := file.ColorBaseName()
// 	if file.IsDir() {
// 		dir, _ := filepath.Split(file.Path)
// 		if len(root) > 0 {
// 			// dir = strings.TrimPrefix(dir, root)
// 			dir = strings.Replace(dir, root, RootMark, 1)
// 		}
// 		name = cdirp.Sprint(dir) + name
// 	}
// 	if file.IsLink() {
// 		name += cdashp.Sprint(" -> ") + file.ColorLinkPath()
// 	}
// 	return name
// }

// GetColorizedPath will return a colorful string of {{ dir }}/{{ name }}
func GetColorizedPath(path string, root string) string {
	if path == PathSeparator {
		return cdip.Sprint(path)
	}

	file, err := NewFileRelTo(path, root)
	if err != nil {
		dir, name := filepath.Split(path)
		dir = PathRel(dir, root)
		return cdirp.Sprint(dir) + cfip.Sprint(name)
	}

	cdir := cdirp.Sprint(file.Dir)
	cname := file.ColorBaseNameToLink()
	return cdir + "/" + cname
}

func getColorizedRootHead(pad, root string, size uint64) string {
	head := fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, GetColorizedPath(root, ""), GetColorizedSize(size))
	return head
}

// func getDirAndName(path string, root string) (dir, name string) {
// 	file, err := NewFile(path)
// 	if err != nil {
// 		dir, name = filepath.Split(path)
// 		if len(root) > 0 {
// 			dir = paw.Replace(dir, root, RootMark, 1)
// 		}
// 		return dir, name
// 	}
// 	name = file.BaseName
// 	if file.IsDir() {
// 		dir, _ = filepath.Split(file.Path)
// 		if len(root) > 0 {
// 			dir = paw.Replace(dir, root, RootMark, 1)
// 		}
// 	}
// 	if file.IsLink() {
// 		return dir + name, file.LinkPath()
// 	}
// 	return dir, name
// }

func getDirInfo(fl *FileList, file *File) (cdinf string, wdinf int) {
	if !file.IsDir() {
		return "", 0
	}
	nd, nf := 0, 0
	files := fl.Map()[file.Dir]
	for _, f := range files {
		if f.IsDir() {
			nd++
		} else {
			nf++
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

func dirSummary(pad string, ndirs int, nfiles int, sumsize uint64) string {
	var (
		cndirs   = csnp.Sprint(ndirs)
		cnfiles  = csnp.Sprint(nfiles)
		csumsize = GetColorizedSize(sumsize)
	)
	msg := pad +
		cndirs +
		cdashp.Sprint(" directories, ") +
		cnfiles +
		cdashp.Sprint(" files, size ≈ ") +
		csumsize +
		cdashp.Sprint(".")
	// msg := fmt.Sprintf("%s%v directories; %v files, size ≈ %v.\n", pad, cndirs, cnfiles, csumsize)
	return msg
}

func printDirSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize uint64) {
	fmt.Fprintln(w, dirSummary(pad, ndirs, nfiles, sumsize))
}

func totalSummary(pad string, ndirs int, nfiles int, sumsize uint64) string {
	var (
		cndirs   = csnp.Sprint(ndirs)
		cnfiles  = csnp.Sprint(nfiles)
		csumsize = GetColorizedSize(sumsize)
	)
	summary := pad +
		cdashp.Sprint("Accumulated ") +
		cndirs +
		cdashp.Sprint(" directories, ") +
		cnfiles +
		cdashp.Sprint(" files, total size ≈ ") +
		csumsize +
		cdashp.Sprint(".")
	// fmt.Sprintf("%sAccumulated %v directories, %v files, total size ≈ %v.\n", pad, cndirs, cnfiles, csumsize)
	return summary
}

func printTotalSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize uint64) {

	fmt.Fprintln(w, totalSummary(pad, ndirs, nfiles, sumsize))
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
			if strings.HasPrefix(k, file.Path) {
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
			if strings.HasPrefix(k, file.Path) {
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
func GetColorizePermission(sperm string) string {
	// sperm := fmt.Sprintf("%v", mode)
	// func GetColorizePermission(mode os.FileMode) string {
	// if strings.HasPrefix(sperm, "Dc") {
	// 	sperm = strings.Replace(sperm, "Dc", "c", 1)
	// }
	// if strings.HasPrefix(sperm, "D") {
	// 	sperm = strings.Replace(sperm, "D", "b", 1)
	// }
	c := ""
	// fmt.Println(len(s))
	for i := 0; i < len(sperm); i++ {
		s := string(sperm[i])
		cs := s
		if cs != "-" {
			switch i {
			case 0:
				switch s {
				case "b": // "D":
					cs = "bd"
				case "c":
					cs = "cd"
				case "d":
					cs = "di"
				case "l": // "L":
					cs = "ln"
				case "p":
					cs = "pi"
				case "S":
					cs = "so"
				}
			case 1, 2, 3:
				// switch s {
				// case "c":
				// 	cs = "cd"
				// default:
				// 	cs = "u" + s
				// }
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
		c += paw.NewEXAColor(cs).Sprint(s)
	}

	return c
}

var cpmap = map[rune]*color.Color{
	'L': clnp,
	'l': clnp,
	'd': cdip,
	'r': paw.NewEXAColor("ur"),
	'w': paw.NewEXAColor("uw"),
	'x': paw.NewEXAColor("ux"),
	'-': cdashp,                //color.New(color.Concealed),
	'.': cdashp,                //color.New(color.Concealed),
	' ': cdashp,                //color.New(color.Concealed), //unmodified
	'M': paw.NewEXAColor("gm"), //color.New(EXAColors["gm"]...), //modified
	'A': paw.NewEXAColor("ga"), //color.New(EXAColors["ga"]...), //added
	'D': paw.NewEXAColor("gd"), //color.New(EXAColors["gd"]...), //deleted
	'R': paw.NewEXAColor("gv"), //color.New(EXAColors["gv"]...), //renamed
	'C': paw.NewEXAColor("gt"), //color.New(EXAColors["gt"]...), //copied
	'U': paw.NewEXAColor("gt"), //color.New(EXAColors["gt"]...), //updated but unmerged
	'?': paw.NewEXAColor("gm"), //color.New(EXAColors["gm"]...), //untracked
	'N': paw.NewEXAColor("ga"), //color.New(EXAColors["ga"]...), //untracked
	'!': cdashp,                //color.New(EXAColors["-"]...),  //ignored
}

// GetColorizedSize will return a humman-readable and colorful string of size.
// The length of placeholder in terminal is 6.
func GetColorizedSize(size uint64) (csize string) {
	ss := ByteSize(size)
	nss := len(ss)
	sn := fmt.Sprintf("%s", ss[:nss-1])
	su := strings.ToLower(ss[nss-1:])
	csize = csnp.Sprint(sn) + csup.Sprint(su)
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
	cdate = strings.Join(dates, "")
	cdate = cdate[:len(cdate)-1]
	wd--
	return cdate, wd
}

func printBanner(w io.Writer, pad string, mark string, length int) {
	banner := fmt.Sprintf("%s%s", pad, strings.Repeat(mark, length))
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

func GetFileLSColor(file *File) *color.Color {

	if file.IsDir() { // os.ModeDir
		return cdip
	}

	if file.IsLink() { // os.ModeSymlink
		return clkp
	}

	if file.IsDev() { //
		return cbdp
	}

	if file.IsChardev() { // os.ModeDevice | os.ModeCharDevice
		return ccdp
	}
	if file.IsFiFo() { //os.ModeNamedPipe
		return cpip
	}
	if file.IsSocket() { //os.ModeSocket
		return csop
	}
	if file.IsFile() { // 0
		if _, ok := paw.LSColors[file.BaseName]; ok {
			return color.New(paw.LSColors[file.BaseName]...)
		}
		if _, ok := paw.LSColors[file.Ext]; ok {
			return color.New(paw.LSColors[file.Ext]...)
		}
		for re, att := range paw.ReExtLSColors {
			if re.MatchString(file.BaseName) {
				return color.New(att...)
			}
		}

		mode := file.Stat.Mode()
		sperm := fmt.Sprintf("%v", mode)
		if mode.IsRegular() && !mode.IsDir() && strings.Contains(sperm, "x") {
			return cexp
		}

		return cfip
	}
	return cnop
}

// FileLSColorString will return the color string of `s` according `fullpath` (xxx.yyy)
func FileLSColorString(fullpath, s string) (string, error) {

	file, err := NewFile(fullpath)
	if err != nil {
		return "", err
	}
	return file.LSColor().Sprint(s), nil
}

func rowWrapDirName(dirName, pad string, wdlimit int) string {
	var (
		w         = new(strings.Builder)
		wpad      = paw.StringWidth(paw.StripANSI(pad))
		sppad     = paw.Spaces(wpad)
		dir, name = filepath.Split(dirName)
		wdir      = paw.StringWidth(dir)
		wname     = paw.StringWidth(name)
		wlen      = wpad + wdir + wname
		width     = wdlimit - wpad
	)
	if wlen <= wdlimit { // one line dir
		cname := cdirp.Sprint(dir) + cdip.Sprint(name)
		fmt.Fprintf(w, "%s%s\n", pad, cname)
	} else { // two more lines dir
		if wpad+wdir <= wdlimit { // 1. condisde dir
			// 1st line dir
			fmt.Fprintf(w, "%s%s", pad, cdirp.Sprint(dir))
			wred := width - wdir
			if wname <= wred {
				// 1st line cont. +name
				fmt.Fprintf(w, "%s\n", cdip.Sprint(name))
			} else {
				// 1st line cont. +name
				name0 := paw.Truncate(name, wred, "")
				fmt.Fprintf(w, "%s\n", cdip.Sprint(name0))
				// 2nd line start, only name
				p0 := len(name0)
				name1 := name[p0:]
				wname1 := paw.StringWidth(name1)
				if wname1 <= width {
					// 2nd line, end name
					fmt.Fprintf(w, "%s%s\n", sppad, cdip.Sprint(name1))
				} else {
					// 2nd and more lines, name only
					names := paw.WrapToSlice(name1, width)
					for _, v := range names {
						fmt.Fprintf(w, "%s%s\n", sppad, cdip.Sprint(v))
					}
				}
			}
		} else { // wpad+wdir > wdlimit
			// 1. two and more lines, dir only
			dirs := paw.WrapToSlice(dir, width)
			fmt.Fprintf(w, "%s%s\n", pad, cdirp.Sprint(dirs[0]))
			for i := 1; i < len(dirs)-1; i++ {
				fmt.Fprintf(w, "%s%s\n", sppad, cdirp.Sprint(dirs[i]))
			}
			// 2. last line of dir
			fmt.Fprintf(w, "%s%s", sppad, cdirp.Sprint(dirs[len(dirs)-1]))
			wred := width - paw.StringWidth(dirs[len(dirs)-1])
			if wname <= wred {
				// 3.1 step 2 cont., +name and end
				fmt.Fprintf(w, "%s\n", cdip.Sprint(name))
			} else {
				// 3.2 step 2 cont., +name and more lines
				name0 := paw.Truncate(name, wred, "")
				fmt.Fprintf(w, "%s\n", cdip.Sprint(name0))
				// 4. two and more lines, name
				p0 := len(name0)
				name1 := name[p0:]
				wname1 := paw.StringWidth(name1)
				if wname1 <= width {
					fmt.Fprintf(w, "%s%s\n", sppad, cdip.Sprint(name1))
				} else {
					names := paw.WrapToSlice(name1, width)
					for _, v := range names {
						fmt.Fprintf(w, "%s%s\n", sppad, cdip.Sprint(v))
					}
				}
			}
		}
	}
	return w.String()
}

func rowWrapFileName(file *File, fds *FieldSlice, pad string, wdsttylimit int) string {
	var (
		sb = new(strings.Builder)
		// wpad   = paw.StringWidth(pad)
		meta   = fds.ColorMetaValuesString()
		wmeta  = fds.MetaValuesStringWidth()
		spmeta = paw.Spaces(wmeta)
		name   = file.BaseNameToLink()
		wname  = paw.StringWidth(name)
		// wdstty = wdsttylimit - 1
		width = fds.Get(PFieldName).Width // wdstty - wpad - wmeta
	)
	if wname <= width {
		fmt.Fprintln(sb, pad+meta, file.ColorName())
	} else { // wrap file name
		if err := paw.CheckIndexInString(name, width, "Name"); err != nil {
			paw.Error.Fatal(err, " (may be too many fields)")
		}
		if !file.IsLink() {
			names := paw.WrapToSlice(name, width)
			fmt.Fprintln(sb, pad+meta, file.LSColor().Sprint(names[0]))
			for i := 1; i < len(names); i++ {
				fmt.Fprintln(sb, pad+spmeta, file.LSColor().Sprint(names[i]))
			}
		} else {
			cname := file.LSColor().Sprint(file.BaseName)
			wbname := paw.StringWidth(file.BaseName)
			carrow := cdashp.Sprint(" -> ")
			wbname += 4
			fmt.Fprintln(sb, pad+meta, cname+carrow)
			dir, name := filepath.Split(file.LinkPath())
			wd, wn := paw.StringWidth(dir), paw.StringWidth(name)

			if wd+wn <= width {
				fmt.Fprintln(sb, pad+spmeta, cdirp.Sprint(dir)+cdip.Sprint(name))
			} else {
				if wd <= width {
					clink := cdirp.Sprint(dir) + cdip.Sprint(name[:width-wd])
					fmt.Fprintln(sb, pad+spmeta, clink)
					names := paw.WrapToSlice(name[width-wd:], width)
					for _, v := range names {
						clink = cdip.Sprint(v)
						fmt.Fprintln(sb, pad+spmeta, clink)
					}
				} else { // wd > width
					dirs := paw.WrapToSlice(dir, width)
					nd := len(dirs)
					var clink string
					for i := 0; i < nd-1; i++ {
						clink = cdirp.Sprint(dirs[i])
						fmt.Fprintln(sb, pad+spmeta, clink)
					}
					clink = cdirp.Sprint(dirs[nd-1])
					wdLast := paw.StringWidth(dirs[nd-1])
					if wn <= width-wdLast {
						clink += cdip.Sprint(name)
						fmt.Fprintln(sb, pad+spmeta, clink)
					} else { // wn > wd-width
						clink += cdip.Sprint(name[:width-wdLast])
						fmt.Fprintln(sb, pad+spmeta, clink)
						rname := name[width-wdLast:]
						wr := paw.StringWidth(rname)
						if wr <= width {
							clink = cdip.Sprint(rname)
							fmt.Fprintln(sb, pad+spmeta, clink)
						} else { // wr > width
							names := paw.WrapToSlice(rname, width)
							for _, v := range names {
								clink = cdip.Sprint(v)
								fmt.Fprintln(sb, pad+spmeta, clink)
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
		nx   = len(file.XAttributes)
		sb   = new(strings.Builder)
		edge = EdgeTypeMid
	)

	for i := 0; i < nx; i++ {
		var (
			xattr = file.XAttributes[i]
			wdx   = paw.StringWidth(xattr)
			wdm   = wmeta
		)
		if i == nx-1 {
			edge = EdgeTypeEnd
		}
		var padx = fmt.Sprintf("%s %s ", pad, cdashp.Sprint(edge))
		wdm += edgeWidth[edge] + 2
		width := wdsttylimit - wdm

		if wdx <= width {
			fmt.Fprintln(sb, padx+cxap.Sprint(xattr))
		} else {
			// var wde = wdsttylimit - wdm
			if err := paw.CheckIndexInString(xattr, width, "xattr"); err != nil {
				paw.Error.Fatal(err, " (may be too many fields)")
			}
			x1 := paw.Truncate(xattr, width, "")
			b := len(x1)
			fmt.Fprintln(sb, padx+cxap.Sprint(x1))
			switch edge {
			case EdgeTypeMid:
				padx = fmt.Sprintf("%s %s ", pad, cdashp.Sprint(EdgeTypeLink)+SpaceIndentSize)
			case EdgeTypeEnd:
				padx = fmt.Sprintf("%s %s ", pad, paw.Spaces(edgeWidth[edge]))
			}

			if len(xattr[b:]) <= width {
				fmt.Fprintln(sb, padx+cxap.Sprint(xattr[b:]))
			} else {
				xattrs := paw.WrapToSlice(xattr[b:], width)
				for _, v := range xattrs {
					fmt.Fprintln(sb, padx+cxap.Sprint(v))
				}
			}
		}
	}
	return sb.String()
}

func PathRel(dir, root string) (rdir string) {
	if !strings.HasPrefix(dir, root) {
		return dir
	}
	rdir = strings.Replace(dir, root, RootMark, 1)
	return rdir
}

// func getMaxFileSizeWidth(files []*File) int {
// 	var (
// 		wdsize = 0
// 	)
// 	for _, f := range files {
// 		var size = ByteSize(f.Size)
// 		if wdsize < len(size) {
// 			wdsize = len(size)
// 		}
// 	}
// 	return wdsize
// }

// func modifyHead(fds *FieldSlice, files []*File, pad string, wdstty int) (chead string, wdmeta int) {

// 	wdsize := getMaxFileSizeWidth(files)
// 	fds.Get(PFieldSize).Width = paw.MaxInt(wdsize, fds.Get(PFieldSize).Width)
// 	chead = fds.ColorHeadsString()
// 	wdmeta = fds.MetaHeadsStringWidth() + paw.StringWidth(pad)
// 	fds.Get(PFieldName).Width = wdstty - wdmeta - 1

// 	return chead, wdmeta
// }

// func modifyFDSTreeHead(fds *FieldSlice, fl *FileList) (chead string, wdmeta int) {
// 	wdsize := 0
// 	for _, dir := range fl.Dirs() {
// 		files := fl.Map()[dir][1:]
// 		wd := getMaxFileSizeWidth(files)
// 		if wdsize < wd {
// 			wdsize = wd
// 		}
// 	}
// 	fds.Get(PFieldSize).Width = paw.MaxInt(wdsize, fds.Get(PFieldSize).Width)
// 	chead = fds.ColorHeadsString()
// 	wdmeta = fds.MetaHeadsStringWidth()
// 	return chead, wdmeta
// }

// func modifyFDSWidth(fds *FieldSlice, fl *FileList, sttyLimit int) {
// 	wdsize := 0
// 	for _, dir := range fl.Dirs() {
// 		files := fl.Map()[dir][1:]
// 		wd := getMaxFileSizeWidth(files)
// 		if wdsize < wd {
// 			wdsize = wd
// 		}
// 	}
// 	fds.Get(PFieldSize).Width = paw.MaxInt(wdsize, fds.Get(PFieldSize).Width)

// 	wdmeta := fds.MetaHeadsStringWidth() + 1
// 	fds.Get(PFieldName).Width = sttyLimit - wdmeta
// }
