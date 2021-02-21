package filetree

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"time"

	"github.com/fatih/color"
	"github.com/karrick/godirwalk"
	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
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
	corp                  = paw.Corp  // orphan file
	cNop                  = paw.CNop  // serial number
	cinp                  = paw.Cinp  // inode
	cpms                  = paw.Cpms  // permission
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
	cmd5p                 = paw.Cmd5p // md5
	cxap                  = paw.Cxap  // extended attributes
	cxbp                  = paw.Cxbp  // extended attributes
	cdashp                = paw.Cdashp
	cnop                  = paw.CNop    // no this file kind
	cbdp                  = paw.Cbdp    // device
	ccdp                  = paw.Ccdp    // CharDevice
	cpip                  = paw.Cpip    // named pipe
	csop                  = paw.Csop    // socket
	cexp                  = paw.Cexp    // execution
	clnp                  = paw.Clnp    // symlink
	cpmpt                 = paw.Cpmpt   // prompt
	cpmptSn               = paw.CpmptSn // number in prompt
	cpmptSu               = paw.CpmptSu // unit in prompt
	currentuser, _        = user.Current()
	urname                = currentuser.Username
	usergp, _             = user.LookupGroupId(currentuser.Gid)
	gpname                = usergp.Name
	curname               = cuup.Sprint(urname)
	cgpname               = cgup.Sprint(gpname)
	ctrace                = paw.Ctrace
	cdebug                = paw.Cdebug
	cinfo                 = paw.Cinfo
	cwarn                 = paw.Cwarn
	cerror                = paw.Cerror
	cfatal                = paw.Cfatal
	cpanic                = paw.Cpanic
	sttyHeight, sttyWidth = paw.GetTerminalSize()
)

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
		meta   = fds.MetaValuesStringC()
		wmeta  = fds.MetaValuesStringWidth()
		spmeta = paw.Spaces(wmeta)
		name   = file.BaseNameToLink()
		wname  = paw.StringWidth(name)
		// wdstty = wdsttylimit - 1
		width = fds.Get(PFieldName).Width // wdstty - wpad - wmeta
	)
	if wname <= width {
		fmt.Fprintln(sb, pad+meta, file.NameC())
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
	if len(root) == 0 || !strings.HasPrefix(dir, root) {
		return dir
	}
	rdir = strings.Replace(dir, root, RootMark, 1)
	return rdir
}

// FileLSColorString will return the color string of `s` according `fullpath` (xxx.yyy)
func FileLSColorString(fullpath, s string) (string, error) {
	file, err := NewFile(fullpath)
	if err != nil {
		return "", err
	}
	return file.LSColor().Sprint(s), nil
}

func GetFileLSColor(file *File) *color.Color {

	if file.IsDir() { // os.ModeDir
		return cdip
	}

	if file.IsLink() { // os.ModeSymlink
		lpath := file.LinkPath()
		if _, err := os.Lstat(lpath); os.IsNotExist(err) {
			return paw.NewLSColor("or")
		}
		return clnp
	}

	if file.IsCharDev() { // os.ModeDevice | os.ModeCharDevice
		return ccdp
	}

	if file.IsDev() { //
		return cbdp
	}

	if file.IsFIFO() { //os.ModeNamedPipe
		return cpip
	}
	if file.IsSocket() { //os.ModeSocket
		return csop
	}

	if file.IsExecutable() && !file.IsDir() {
		return cexp
	}

	if file.IsFile() { // 0
		if att, ok := paw.LSColors[file.BaseName]; ok {
			return color.New(att...)
		}
		if att, ok := paw.LSColors[file.Ext]; ok {
			return color.New(att...)
		}
		for re, att := range paw.ReExtLSColors {
			if re.MatchString(file.BaseName) {
				return color.New(att...)
			}
		}
		return cfip
	}
	return cnop
}

// GetColorizedPath will return a colorful string of {{ dir }}/{{ name }}
func GetColorizedPath(path string, root string) string {
	if path == PathSeparator {
		return cdip.Sprint(path)
	}

	file, err := NewFileRelTo(path, root)
	if err != nil {
		dir, name := filepath.Split(path)
		dir = PathRel(dir, root)
		// c := file.LSColor() //GetFileLSColor(file)
		return cdirp.Sprint(dir) + corp.Sprint(name)
	}

	cname := file.BaseNameToLinkC()
	if file.Dir == "/" {
		return "/" + cname
	} else {
		cdir := cdirp.Sprint(file.Dir)
		return cdir + "/" + cname
	}
}

func pmptColorizedPath(path string, root string) string {
	if path == PathSeparator {
		return cpmpt.Sprint(cdip.Sprint(path))
	}

	file, err := NewFileRelTo(path, root)
	if err != nil {
		dir, name := filepath.Split(path)
		dir = PathRel(dir, root)
		return cpmpt.Sprint(cdirp.Sprint(dir)) + cpmpt.Sprint(cfip.Sprint(name))
	}

	cname := cpmpt.Sprint(file.BaseNameC())
	if file.IsLink() {
		lfile, err := NewFileRelTo(file.LinkPath(), "")
		if err != nil {
			cname += cpmpt.Sprint(cdashp.Sprint(" -> ")) + pmptColorizedPath(lfile.Path, "")
		} else {
			cname += cpmpt.Sprint(cdashp.Sprint(" -> "))
			dir, name := filepath.Split(file.LinkPath())
			c := GetFileLSColor(lfile)
			cname += cpmpt.Sprint(cdirp.Sprint(dir))
			cname += cpmpt.Sprint(c.Sprint(name))
		}
	}
	if file.Dir == "/" {
		return cpmpt.Sprint("/") + cname
	} else {
		cdir := cpmpt.Sprint(cdirp.Sprint(file.Dir))
		return cdir + cpmpt.Sprint("/") + cname
	}
}

func getColorizedRootHead(root string, size uint64, wdstty int) string {
	var (
		ss  = ByteSize(size)
		nss = len(ss)
		sn  = ss[:nss-1] // fmt.Sprintf("%s", ss[:nss-1])
		su  = strings.ToLower(ss[nss-1:])
	)
	if pdOpt != nil && pdOpt.File != nil {
		if pdOpt.File.IsLink() {
			root = pdOpt.File.Path
		}
	}
	chead := cpmpt.Sprint("Root directory: ")
	chead += pmptColorizedPath(root, "")
	chead += cpmpt.Sprint(", size ≈ ")
	chead += cpmptSn.Sprint(sn) + cpmptSu.Sprint(su)
	chead += cpmpt.Sprint(".")

	chead += cpmpt.Sprint(paw.Spaces(wdstty + 1 - paw.StringWidth(paw.StripANSI(chead))))

	// chead := fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, GetColorizedPath(root, ""), GetColorizedSize(size))
	return chead
}

// func getDirInfo(fl *FileList, file *File) (cdinf string, wdinf int) {
// 	files := fl.GetFiles(file.Dir) //fl.Map()[file.Dir]
// 	if !file.IsDir() || files == nil {
// 		return "", 0
// 	}

// 	nd, nf := fl.Map().CountDF(file.Dir)
// 	nd--
// 	// nd, nf := 0, 0
// 	// for _, f := range files[1:] {
// 	// 	if f.IsDir() {
// 	// 		nd++
// 	// 	} else {
// 	// 		nf++
// 	// 	}
// 	// }
// 	di := fmt.Sprintf("%v dirs", nd)
// 	fi := fmt.Sprintf("%v files", nf)
// 	wdinf = len(di) + len(fi) + 4
// 	cdinf = fmt.Sprintf("[%s, %s]", cdirp.Sprint(di), cdirp.Sprint(fi))
// 	return cdinf, wdinf
// }

func dirSummary(pad string, ndirs int, nfiles int, sumsize uint64, wdstty int) string {
	var (
		ss  = ByteSize(sumsize)
		nss = len(ss)
		sn  = fmt.Sprintf("%s", ss[:nss-1])
		su  = strings.ToLower(ss[nss-1:])
	)
	cndirs := cpmptSn.Sprint(ndirs)
	cnfiles := cpmptSn.Sprint(nfiles)
	csumsize := cpmptSn.Sprint(sn) + cpmptSu.Sprint(su)
	msg := pad +
		cndirs +
		cpmpt.Sprint(" directories, ") +
		cnfiles +
		cpmpt.Sprint(" files, size ≈ ") +
		csumsize +
		cpmpt.Sprint(". ")
	nmsg := paw.StringWidth(paw.StripANSI(msg))
	msg += cpmpt.Sprint(paw.Spaces(wdstty + 1 - nmsg))
	// msg := fmt.Sprintf("%s%v directories; %v files, size ≈ %v.\n", pad, cndirs, cnfiles, csumsize)
	return msg
}

func printDirSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize uint64, wdstty int) {
	fmt.Fprintln(w, dirSummary(pad, ndirs, nfiles, sumsize, wdstty))
}

func totalSummary(pad string, ndirs int, nfiles int, sumsize uint64, wdstty int) string {
	var (
		ss  = ByteSize(sumsize)
		nss = len(ss)
		sn  = fmt.Sprintf("%s", ss[:nss-1])
		su  = strings.ToLower(ss[nss-1:])
	)
	cndirs := cpmptSn.Sprint(ndirs)
	cnfiles := cpmptSn.Sprint(nfiles)
	csumsize := cpmptSn.Sprint(sn) + cpmptSu.Sprint(su)
	summary := pad +
		cpmpt.Sprint("Accumulated ") +
		cndirs +
		cpmpt.Sprint(" directories, ") +
		cnfiles +
		cpmpt.Sprint(" files, total size ≈ ") +
		csumsize +
		cpmpt.Sprint(".")
	nsummary := paw.StringWidth(paw.StripANSI(summary))
	summary += cpmpt.Sprint(paw.Spaces(wdstty + 1 - nsummary))
	// fmt.Sprintf("%sAccumulated %v directories, %v files, total size ≈ %v.\n", pad, cndirs, cnfiles, csumsize)
	return summary
}

func printTotalSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize uint64, wdstty int) {

	fmt.Fprintln(w, totalSummary(pad, ndirs, nfiles, sumsize, wdstty))
}

func printBanner(w io.Writer, pad string, mark string, length int) {
	banner := cdashp.Sprintf("%s%s\n", pad, strings.Repeat(mark, length))
	fmt.Fprint(w, banner)
}

// GetColorizedPermission will return a colorful string of mode
// The length of placeholder in terminal is 10.
func GetColorizedPermission(sperm string) string {
	ns := len(sperm)
	cxmark := cdashp.Sprint(string(sperm[ns-1]))
	perm := sperm[ns-10 : ns-1]
	abbr := string(sperm[:ns-10])
	cabbr := ""
	for _, a := range abbr {
		s := string(a)
		cs := "-"
		switch s {
		case "d": // d: is a directory
			cs = "di"
		case "a": // a: append-only
			cs = "ca"
		// case "l": // l: exclusive use
		// case "T": // T: temporary file; Plan 9 only
		case "L": // L: symbolic link
			cs = "ln"
		case "D": // D: device file
			cs = "bd"
		case "c": // c: Unix character device, when ModeDevice is set
			cs = "cd"
		case "p": // p: named pipe (FIFO)
			cs = "pi"
		case "S": // S: Unix domain socket
			cs = "so"
		case "u": // u: setuid
			cs = "su"
		case "g": // g: setgid
			cs = "sg"
		case "t": // t: sticky
			cs = "st"
		case "?": // ?: non-regular file; nothing else is known about this file
			cs = "-"
		case "-":
			s = "."
			cs = "-"
		default:
			cs = "no"
		}
		cabbr += paw.NewEXAColor(cs).Sprint(s)
	}
	c := ""
	// fmt.Println(len(s))
	for i := 0; i < len(perm); i++ {
		s := string(perm[i])
		cs := s
		if cs != "-" {
			switch i {
			case 0, 1, 2:
				cs = "u" + s
			case 3, 4, 5:
				cs = "g" + s
			case 6, 7, 8:
				cs = "t" + s
			}
		}
		// if i == 0 && cs == "-" {
		// 	s = "."
		// }
		// c += color.New(EXAColors[cs]...).Add(color.Bold).Sprint(s)
		c += paw.NewEXAColor(cs).Sprint(s)
	}

	return cabbr + c + cxmark
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

// var cpmap = map[rune]*color.Color{
// 	'L': clnp,
// 	'l': clnp,
// 	'd': cdip,
// 	'r': paw.NewEXAColor("ur"),
// 	'w': paw.NewEXAColor("uw"),
// 	'x': paw.NewEXAColor("ux"),
// 	'-': cdashp,                //color.New(color.Concealed),
// 	'.': cdashp,                //color.New(color.Concealed),
// 	' ': cdashp,                //color.New(color.Concealed), //unmodified
// 	'M': paw.NewEXAColor("gm"), //color.New(EXAColors["gm"]...), //modified
// 	'A': paw.NewEXAColor("ga"), //color.New(EXAColors["ga"]...), //added
// 	'D': paw.NewEXAColor("gd"), //color.New(EXAColors["gd"]...), //deleted
// 	'R': paw.NewEXAColor("gv"), //color.New(EXAColors["gv"]...), //renamed
// 	'C': paw.NewEXAColor("gt"), //color.New(EXAColors["gt"]...), //copied
// 	'U': paw.NewEXAColor("gt"), //color.New(EXAColors["gt"]...), //updated but unmerged
// 	'?': paw.NewEXAColor("gm"), //color.New(EXAColors["gm"]...), //untracked
// 	'N': paw.NewEXAColor("ga"), //color.New(EXAColors["ga"]...), //untracked
// 	'!': cdashp,                //color.New(EXAColors["-"]...),  //ignored
// }

// var ckxy = map[rune]rune{
// 	'M': 'M', //modified
// 	'A': 'A', //added
// 	// 'A': 'N',
// 	'D': 'D', //deleted
// 	'R': 'R', //renamed
// 	'C': 'C', //copied
// 	'U': 'U', //updated but unmerged
// 	// '?': 'N', //untracked
// 	'?': '?', //untracked
// 	' ': '-',
// 	'!': 'I', //ignore
// 	// '!': '!', //ignore
// }

// // getColorizedGitStatus will return a colorful string of shrot status of git.
// // The length of placeholder in terminal is 3.
// func getColorizedGitStatus(git *GitStatus, file *File) string {
// 	x, y := '-', '-'

// 	xy, ok := git.FilesStatus[file.Path]
// 	if ok {
// 		x, y = xy.Split()
// 		x, y = ckxy[x], ckxy[y]
// 	}

// 	if file.IsDir() {
// 		// gits := getGitSlice(git, file)
// 		for k, v := range git.FilesStatus {
// 			if strings.HasPrefix(k, file.Path) {
// 				vx, vy := v.Split()
// 				cx, cy := ckxy[vx], ckxy[vy]
// 				if cx != '-' { //&& x != 'N' {
// 					x = cx
// 				}
// 				if cy != '-' { //&& y != 'N' {
// 					y = cy
// 				}
// 			}
// 		}
// 	}

// 	var sx, sy string
// 	if x == 'N' && y == 'N' {
// 		sx, sy = "-", "N"
// 	} else {
// 		sx, sy = string(x), string(y)
// 	}

// 	return " " + cpmap[x].Sprint(sx) + cpmap[y].Sprint(sy)
// }

// // getGitStatus will return a  string of shrot status of git.
// // The length of placeholder in terminal is 3.
// func getGitStatus(git GitStatus, file *File) string {
// 	x, y := '-', '-'

// 	xy, ok := git.FilesStatus[file.Path]
// 	if ok {
// 		x, y = xy.Split()
// 		x, y = ckxy[x], ckxy[y]
// 	}

// 	if file.IsDir() {
// 		// gits := getGitSlice(git, file)
// 		for k, v := range git.FilesStatus {
// 			if strings.HasPrefix(k, file.Path) {
// 				vx, vy := v.Split()
// 				cx, cy := ckxy[vx], ckxy[vy]
// 				if cx != '-' && x != 'N' {
// 					x = cx
// 				}
// 				if cy != '-' && y != 'N' {
// 					y = cy
// 				}
// 			}
// 		}
// 	}

// 	var sx, sy string
// 	if x == 'N' && y == 'N' {
// 		sx, sy = "-", "N"
// 	} else {
// 		sx, sy = string(x), string(y)
// 	}

// 	return " " + sx + sy
// }

func isEnded(levelsEnded []int, level int) bool {
	for _, l := range levelsEnded {
		if l == level {
			return true
		}
	}
	return false
}

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
			// paw.Error.Printf("")
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

func showlogrus() {
	if pdOpt.isTrace {
		paw.Logger.Trace("trace")
		paw.Logger.Debug("debug")
		paw.Logger.Info("info")
		paw.Logger.Warn("warn")
		paw.Logger.Error("error")
	}
}

func gpReadDirnames(f *FileList, dirPath string) error {
	//  maybe BUGS
	files, err := godirwalk.ReadDirnames(f.root, nil)
	if err != nil {
		return errors.New(f.root + ": " + err.Error())
	}

	file, err := NewFileRelTo(f.root, f.root)
	if err != nil {
		return err
	}

	f.AddFile(file)

	var pdir = file
	for _, name := range files {
		path := filepath.Join(f.root, name)
		file, err := NewFileRelTo(path, f.root)
		if err != nil {
			// paw.Logger.Error(err)
			paw.Error.Printf("accesing path %q, %v\n", path, err)
			// return err
			continue
		}
		if err := f.ignore(file, nil); err != SkipThis {
			// continue
			file.SetUpDir(pdir)
			f.AddFile(file)
		}
	}
	return nil
}

func gpScanDir(f *FileList, dirPath string) error {
	//  maybe BUGS
	dirScan, err := godirwalk.NewScanner(f.root)
	if err != nil {
		return fmt.Errorf("cannot scan directory: %s", err)
	}
	// var pdir = file
	for dirScan.Scan() {
		de, err := dirScan.Dirent()
		if err != nil {
			if pdOpt.isTrace {
				flerr := newFileListError(filepath.Join(f.root, de.Name()), err, f.root)
				paw.Logger.WithFields(logrus.Fields{
					"path": flerr.path,
					// "dir":      flerr.dir,
					// "basename": flerr.basename,
					// "err":      flerr.err,
				}).Error(err)
			}
			f.AddError(filepath.Join(f.root, de.Name()), err)
			continue
		}
		path := filepath.Join(f.root, de.Name())
		file, err := NewFileRelTo(path, f.root)
		if err != nil {
			if pdOpt.isTrace {
				flerr := newFileListError(path, err, f.root)
				paw.Logger.WithFields(logrus.Fields{
					"path": flerr.path,
					// "dir":      flerr.dir,
					// "basename": flerr.basename,
					// "err":      flerr.err,
				}).Error(err)
			}
			f.AddError(path, err)
			continue
		}
		if err := f.ignore(file, nil); err != SkipThis {
			f.AddFile(file)
		}
	}
	return dirScan.Err()
}

func fpWalk(f *FileList) error {
	err := filepath.Walk(f.root, func(path string, info os.FileInfo, err error) error {
		skip := false
		file, errf := NewFileRelTo(path, f.root)
		if errf != nil {
			if pdOpt.isTrace {
				paw.Logger.Error(errf)
			}
			f.AddError(path, errf)
			return nil
		}
		idepth := len(file.DirSlice()) - 1
		if f.depth > 0 {
			if idepth > f.depth {
				skip = true
			}
		}
		if err1 := f.ignore(file, errf); err1 == SkipThis {
			skip = true
			if file.IsDir() {
				return filepath.SkipDir
			}
		}
		if !skip {
			f.AddFile(file)
		}
		return nil
	})
	if err != nil {
		return errors.New(f.root + ": " + err.Error())
	}
	return nil
}

func gdWalk(f *FileList) error {
	err := godirwalk.Walk(f.root, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			file, errf := NewFileRelTo(path, f.root)
			if errf != nil {
				if pdOpt.isTrace {
					paw.Logger.Error(errf)
				}
				return godirwalk.SkipThis
			}
			skip := false
			idepth := len(strings.Split(strings.Replace(path, f.root, ".", 1), PathSeparator)) - 1
			if f.depth > 0 {
				if idepth > f.depth {
					skip = true
				}
			}
			if err1 := f.ignore(file, errf); err1 == SkipThis {
				skip = true
			}
			if !skip {
				f.AddFile(file)
			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			// paw.Logger.Errorf("ERROR: %s\n", err)
			// paw.Error.Printf("ERROR: %s\n", err)
			// if pdOpt.isTrace {
			// 	paw.Logger.WithField("path", osPathname).Error(err)
			// }
			// For the purposes of this example, a simple SkipNode will suffice, although in reality perhaps additional logic might be called for.
			return godirwalk.SkipNode
		},
		FollowSymbolicLinks: true,
		AllowNonDirectory:   false,
		Unsorted:            true, // set true for faster yet non-deterministic enumeration (see godoc)
	})
	if err != nil {
		return errors.New(f.root + ": " + err.Error())
	}
	return nil
}

func osReaddirnames(f *FileList, dirPath string) {
	openDir, err := os.Open(dirPath)
	if err != nil {
		if pdOpt.isTrace {
			paw.Logger.Error(err)
		}
		f.AddError(dirPath, err)
		return
	}
	defer openDir.Close()

	files, err := openDir.Readdirnames(-1)
	if err != nil {
		if pdOpt.isTrace {
			paw.Logger.Error(err)
		}
		f.AddError(dirPath, err)
		// return
	}
	if len(files) > 0 {
		handleFiles(f, dirPath, files)
	}

	return
}

func handleFiles(f *FileList, dirPath string, files []string) {
	nf := len(files)
	if nf == 0 {
		return
	}
	for _, name := range files {
		skip := false
		path := filepath.Join(dirPath, name)
		file, err := NewFileRelTo(path, f.root)
		if err != nil {
			if pdOpt.isTrace {
				paw.Logger.Error(err)
			}
			f.AddError(path, err)
			continue
		}
		if err := f.ignore(file, nil); err == SkipThis {
			skip = true
		}
		idepth := len(file.DirSlice()) - 1
		if f.depth > 0 {
			if idepth > f.depth {
				skip = true
			}
		}
		if !skip {
			f.AddFile(file)
			if f.depth != 0 {
				if file.IsDir() {
					if skip {
						return
					}
					osReaddirnames(f, path)
				}
			}
		}
	}
}

func wgosReaddirnames(f *FileList, dirPath string) {
	sem <- 1
	defer func() {
		<-sem
	}()
	defer wg.Done()

	openDir, err := os.Open(dirPath)
	if err != nil {
		if pdOpt.isTrace {
			paw.Logger.Error(err)
		}
		f.mux.Lock()
		f.AddError(dirPath, err)
		f.mux.Unlock()
		return
	}
	defer openDir.Close()

	files, err := openDir.Readdirnames(-1)
	if err != nil {
		if pdOpt.isTrace {
			paw.Logger.Error(err)
		}
		f.mux.Lock()
		f.AddError(dirPath, err)
		f.mux.Unlock()
		return
	}
	if len(files) > 0 {
		wg.Add(1)
		go wghandleFiles(f, dirPath, files)
	}

	return
}

func wghandleFiles(f *FileList, dirPath string, files []string) {
	sem <- 1
	defer func() {
		<-sem
	}()
	defer wg.Done()

	nf := len(files)
	if nf == 0 {
		return
	}
	// for _, name := range files {
	if nf == 1 {
		skip := false
		name := files[0]
		path := filepath.Join(dirPath, name)
		file, err := NewFileRelTo(path, f.root)
		if err != nil {
			if pdOpt.isTrace {
				paw.Logger.Error(err)
			}
			f.mux.Lock()
			f.AddError(path, err)
			f.mux.Unlock()
			// continue
		}
		if err := f.ignore(file, nil); err == SkipThis {
			skip = true
		}
		idepth := len(file.DirSlice()) - 1
		if f.depth > 0 {
			if idepth > f.depth {
				skip = true
			}
		}
		if !skip {
			f.mux.Lock()
			f.AddFile(file)
			f.mux.Unlock()
			if f.depth != 0 {
				if file.IsDir() {
					if skip {
						return
					}
					wg.Add(1)
					go wgosReaddirnames(f, path)
				}
			}
		}
	} else {
		wg.Add(2)
		go wghandleFiles(f, dirPath, files[:nf/2])
		go wghandleFiles(f, dirPath, files[nf/2:])
	}
	// }
}

func fpWalkDir(f *FileList) {
	root := f.root
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		skip := false
		file, errf := NewFileRelTo(path, root)
		if errf != nil {
			paw.Logger.Error(errf)
			f.AddError(path, errf)
			// return nil
			skip = true
		}
		idepth := len(file.DirSlice()) - 1
		if !skip && f.depth > 0 {
			if idepth > f.depth {
				skip = true
			}
		}
		if err1 := f.ignore(file, errf); !skip && err1 == SkipThis {
			skip = true
			if d.IsDir() {
				return filepath.SkipDir
			}
		}
		if !skip {
			f.AddFile(file)
		}
		return nil
	})

	if err != nil {
		f.AddError(root, err)
		paw.Logger.Error(err)
		paw.Error.Fatal(err)
	}
}
