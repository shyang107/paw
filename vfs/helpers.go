package vfs

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/fatih/color"
	"github.com/pkg/xattr"
	"github.com/shyang107/paw"
)

type EdgeType string

const (
	EdgeTypeLink      EdgeType = "│"   //treeprint.EdgeTypeLink
	EdgeTypeMid       EdgeType = "├──" //treeprint.EdgeTypeMid
	EdgeTypeEnd       EdgeType = "└──" //treeprint.EdgeTypeEnd
	IndentSize                 = 3     //treeprint.IndentSize
	dateLayout                 = "Jan 02, 2006"
	timeThisLayout             = "01-02 15:04"
	timeBeforeLayout           = "2006-01-02"
	PathSeparator              = string(os.PathSeparator)
	PathListSeparator          = string(os.PathListSeparator)
	XattrSymbol                = paw.XAttrSymbol
)

var (
	xattrsp                    = paw.Spaces(paw.StringWidth(XattrSymbol))
	hasMd5                     = false
	edgeWidth map[EdgeType]int = map[EdgeType]int{
		EdgeTypeLink: 1,
		EdgeTypeMid:  3,
		EdgeTypeEnd:  3,
	}
	currentuser, _  = user.Current()
	urname          = currentuser.Username
	usergp, _       = user.LookupGroupId(currentuser.Gid)
	gpname          = usergp.Name
	curname         = cuup.Sprint(urname)
	cgpname         = cgup.Sprint(gpname)
	now             = time.Now()
	thisYear        = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	SpaceIndentSize = paw.Spaces(IndentSize)
	chdp            = paw.Chdp  // head
	cdirp           = paw.Cdirp // pre-dir part of path
	cdip            = paw.Cdip  // directory
	cfip            = paw.Cfip  // file
	corp            = paw.Corp  // orphan file
	cNop            = paw.CNop  // serial number
	cinp            = paw.Cinp  // inode
	cpms            = paw.Cpms  // permission
	csnp            = paw.Csnp  // size number
	csup            = paw.Csup  // size unit
	cuup            = paw.Cuup  // user
	cgup            = paw.Cgup  // group
	cunp            = paw.Cunp  // user is not you
	cgnp            = paw.Cgnp  // group without you
	clkp            = paw.Clkp  // symlink
	cbkp            = paw.Cbkp  // blocks
	cdap            = paw.Cdap  // date
	cgitp           = paw.Cgitp // git
	cmd5p           = paw.Cmd5p // md5
	cxap            = paw.Cxap  // extended attributes
	cxbp            = paw.Cxbp  // extended attributes
	cdashp          = paw.Cdashp
	cnop            = paw.CNop    // no this file kind
	cbdp            = paw.Cbdp    // device
	ccdp            = paw.Ccdp    // CharDevice
	cpip            = paw.Cpip    // named pipe
	csop            = paw.Csop    // socket
	cexp            = paw.Cexp    // execution
	clnp            = paw.Clnp    // symlink
	cpmpt           = paw.Cpmpt   // prompt
	cpmptSn         = paw.CpmptSn // number in prompt
	cpmptSu         = paw.CpmptSu // unit in prompt

	ctrace                = paw.Ctrace
	cdebug                = paw.Cdebug
	cinfo                 = paw.Cinfo
	cwarn                 = paw.Cwarn
	cerror                = paw.Cerror
	cfatal                = paw.Cfatal
	cpanic                = paw.Cpanic
	sttyHeight, sttyWidth = paw.GetTerminalSize()
)

// ===

func nameC(de DirEntryX) string {
	return de.LSColor().Sprint(de.Name())
	// if !de.IsDir() {
	// 	return f.LSColor().Sprint(f.Name())
	// }

	// d, _ := de.(*Dir)
	// return cdip.Sprint(d.Name())
}

func linkC(de DirEntryX) string {
	if !de.IsDir() && de.IsLink() {
		link := de.LinkPath()
		_, err := os.Stat(link)
		if err != nil {
			dir, name := filepath.Split(link)
			return cdirp.Sprint(dir+"/") + corp.Sprint(name)
		}
		dir, name := filepath.Split(link)
		return cdirp.Sprint(dir) + paw.FileLSColor(link).Sprint(name)
	}
	return ""
	// f, isFile := de.(*File)
	// if isFile && f.IsLink() {
	// 	link := f.LinkPath()
	// 	_, err := os.Stat(link)
	// 	if err != nil {
	// 		dir, name := filepath.Split(link)
	// 		return cdirp.Sprint(dir+"/") + corp.Sprint(name)
	// 	}
	// 	dir, name := filepath.Split(link)
	// 	return cdirp.Sprint(dir) + paw.FileLSColor(link).Sprint(name)

	// }
	// return ""
}

func nameToLinkC(de DirEntryX) string {
	if !de.IsDir() {
		if de.IsLink() {
			return nameC(de) + cdashp.Sprint(" -> ") + linkC(de)
		} else {
			return nameC(de)
		}
	}
	return cdip.Sprint(de.Name())
	// if !de.IsDir() {
	// 	f := de.(*File)
	// 	if f.IsLink() {
	// 		return nameC(f) + cdashp.Sprint(" -> ") + linkC(f)
	// 	} else {
	// 		return nameC(f)
	// 	}
	// }
	// return cdip.Sprint(de.Name())
}

func iNodeC(de DirEntryX) string {
	return cinp.Sprint(de.INode())
}

func getXattr(path string) ([]string, error) {
	// paw.Logger.WithField("path", path).Info("income")
	xattrs, err := xattr.List(path)
	if err != nil {
		return xattrs, err
	}
	if len(xattrs) > 0 {
		for i, x := range xattrs {
			x, _ := xattr.Get(path, x)
			xattrs[i] = fmt.Sprintf("%s (len %d)", xattrs[i], len(x))
		}
	}
	return xattrs, nil
}

func permissionS(de DirEntryX) string {
	mode := de.Type()
	sperm := mode.String()
	if de.Xattibutes() == nil {
		sperm += "?"
	} else {
		if len(de.Xattibutes()) > 0 {
			sperm += "@"
		} else {
			sperm += " "
		}
	}
	return sperm
}

func permissionC(de DirEntryX) string {
	sperm := permissionS(de)
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

func sizeS(de DirEntryX) string {
	if !de.IsDir() && de.Mode().IsRegular() {
		return bytefmt.ByteSize(uint64(de.Size()))
		// return paw.FillLeft(bytefmt.ByteSize(uint64(f.Size())), 6)
	} else {
		return "-"
	}
	// f, isFile := de.(*File)
	// if isFile && f.Mode().IsRegular() {
	// 	return bytefmt.ByteSize(uint64(f.Size()))
	// 	// return paw.FillLeft(bytefmt.ByteSize(uint64(f.Size())), 6)
	// } else {
	// 	return "-"
	// }
}

func sizeC(de DirEntryX) (csize string) {
	ss := sizeS(de)
	if ss == "-" {
		return cdashp.Sprint(ss)
	}
	nss := len(ss)
	sn := fmt.Sprintf("%s", ss[:nss-1])
	su := strings.ToLower(ss[nss-1:])
	csize = csnp.Sprint(sn) + csup.Sprint(su)
	return csize
}

func sizeCaligned(de DirEntryX) (csize string) {
	var (
		ss  = sizeS(de)
		nss = len(ss)
	)
	if ss == "-" {
		csize = cdashp.Sprint("-")
	} else {
		sn := fmt.Sprintf("%s", ss[:nss-1])
		su := strings.ToLower(ss[nss-1:])
		csize = csnp.Sprint(sn) + csup.Sprint(su)
	}
	var (
		width = paw.MaxInt(nss, ViewFieldSize.Width())
		sp    = paw.Spaces(width - nss)
	)
	return sp + csize
}

func blocksCaligned(de DirEntryX) (cb string) {
	var (
		ss  = "-"
		nss = 1
	)
	cb = cdashp.Sprint(ss)
	var (
		width = paw.MaxInt(nss, ViewFieldBlocks.Width())
		sp    = paw.Spaces(width - nss)
	)
	return sp + cb
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func dateS(date time.Time) (sdate string) {
	sdate = date.Format(timeThisLayout)
	if date.Before(thisYear) {
		sdate = date.Format(timeBeforeLayout)
	}
	return sdate
	// return paw.FillLeft(sdate, 11)
}

func deLSColor(de DirEntryX) *color.Color {
	if de.IsDir() {
		return cdip
	}

	if de.IsLink() { // os.ModeSymlink
		_, err := os.Readlink(de.Path())
		if err != nil {
			return paw.NewLSColor("or")
		}
		return clnp
	}

	if de.IsCharDev() { // os.ModeDevice | os.ModeCharDevice
		return ccdp
	}

	if de.IsDev() { //
		return cbdp
	}

	if de.IsFIFO() { //os.ModeNamedPipe
		return cpip
	}
	if de.IsSocket() { //os.ModeSocket
		return csop
	}

	if de.IsExecutable() && !de.IsDir() {
		return cexp
	}

	name := de.Name()
	if att, ok := paw.LSColors[name]; ok {
		return color.New(att...)
	}
	ext := filepath.Ext(name)
	if att, ok := paw.LSColors[ext]; ok {
		return color.New(att...)
	}

	for re, att := range paw.ReExtLSColors {
		if re.MatchString(name) {
			return color.New(att...)
		}
	}

	return cfip
	// _, isDir := de.(*Dir)
	// if isDir {
	// 	return cdip
	// }

	// file, _ := de.(*File)

	// if file.IsLink() { // os.ModeSymlink
	// 	_, err := os.Readlink(file.Path())
	// 	if err != nil {
	// 		return paw.NewLSColor("or")
	// 	}
	// 	return clnp
	// }

	// if file.IsCharDev() { // os.ModeDevice | os.ModeCharDevice
	// 	return ccdp
	// }

	// if file.IsDev() { //
	// 	return cbdp
	// }

	// if file.IsFIFO() { //os.ModeNamedPipe
	// 	return cpip
	// }
	// if file.IsSocket() { //os.ModeSocket
	// 	return csop
	// }

	// if file.IsExecutable() && !file.IsDir() {
	// 	return cexp
	// }

	// name := file.Name()
	// if att, ok := paw.LSColors[name]; ok {
	// 	return color.New(att...)
	// }
	// ext := filepath.Ext(name)
	// if att, ok := paw.LSColors[ext]; ok {
	// 	return color.New(att...)
	// }

	// for re, att := range paw.ReExtLSColors {
	// 	if re.MatchString(name) {
	// 		return color.New(att...)
	// 	}
	// }

	// return cfip
}

func aligned(field ViewField, value interface{}) string {
	var (
		align = field.Align()
		s     = fmt.Sprintf("%v", value)
		wd    = paw.StringWidth(paw.StripANSI(s))
		width = paw.MaxInt(wd, field.Width())
		sp    = paw.Spaces(width - wd)
	)

	if field&ViewFieldName == ViewFieldName {
		return s
	}

	switch align {
	case paw.AlignLeft:
		return s + sp
	default:
		return sp + s
	}
}

func checkFieldsHasGit(fields []ViewField, isNoGit bool) []ViewField {
	fds := []ViewField{}
	for _, fd := range fields {
		if fd&ViewFieldGit == ViewFieldGit && isNoGit {
			continue
		}
		fds = append(fds, fd)
	}
	return fds
}

func modFieldWidths(v *VFS, fields []ViewField) {
	rd := v.rootDir
	childWidths(rd, fields)
}

func childWidths(d *Dir, fields []ViewField) {
	ds, _ := d.ReadDir(-1)
	d.resetIdx()
	for _, c := range ds {
		f, isFile := c.(*File)
		if isFile {
			for _, fd := range fields {
				if fd&ViewFieldSize == ViewFieldSize {
					if f.IsCharDev() || f.IsDev() {
						fmajor := viewFieldWidths[_ViewFieldMajor]
						fminor := viewFieldWidths[_ViewFieldMinor]
						major, minor := f.DevNumber()
						wdmajor := len(fmt.Sprint(major))
						wdminor := len(fmt.Sprint(minor))
						viewFieldWidths[_ViewFieldMajor] = paw.MaxInt(fmajor, wdmajor)
						viewFieldWidths[_ViewFieldMinor] = paw.MaxInt(fminor, wdminor)
						wdsize := viewFieldWidths[_ViewFieldMajor] + viewFieldWidths[_ViewFieldMinor] + 1
						wd := viewFieldWidths[fd]
						viewFieldWidths[fd] = paw.MaxInt(wd, wdsize)
					}
				}
				wd := f.WidthOf(fd)
				dwd := fd.Width()
				width := paw.MaxInt(dwd, wd)
				viewFieldWidths[fd] = width

			}
		} else {
			d := c.(*Dir)
			for _, fd := range fields {
				wd := d.WidthOf(fd)
				dwd := fd.Width()
				width := paw.MaxInt(dwd, wd)
				viewFieldWidths[fd] = width
			}
			childWidths(d, fields)
		}
	}
}

func dirSummary(pad string, ndirs int, nfiles int, sumsize int64, wdstty int) string {
	var (
		ss  = bytefmt.ByteSize(uint64(sumsize))
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

func fprintDirSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize int64, wdstty int) {
	fmt.Fprintln(w, dirSummary(pad, ndirs, nfiles, sumsize, wdstty))
}

func totalSummary(pad string, ndirs int, nfiles int, sumsize int64, wdstty int) string {
	var (
		ss  = bytefmt.ByteSize(uint64(sumsize))
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

func fprintTotalSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize int64, wdstty int) {

	fmt.Fprintln(w, totalSummary(pad, ndirs, nfiles, sumsize, wdstty))
}

func fprintBanner(w io.Writer, pad string, mark string, length int) {
	banner := cdashp.Sprintf("%s%s\n", pad, strings.Repeat(mark, length))
	fmt.Fprint(w, banner)
}
