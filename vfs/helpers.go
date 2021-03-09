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

	"github.com/fatih/color"
	"github.com/pkg/xattr"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/bytefmt"
	"github.com/spf13/cast"
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
	currentuser, _        = user.Current()
	urname                = currentuser.Username
	usergp, _             = user.LookupGroupId(currentuser.Gid)
	gpname                = usergp.Name
	curname               = cuup.Sprint(urname)
	cgpname               = cgup.Sprint(gpname)
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
	bgpmpt                = []color.Attribute{48, 5, 236}
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

// GetPathC returns  cdip.Sprint(cname)+cdirp.Sprint(cdir)
func GetPathC(path string) (cdir, cname, cpath string) {
	cdir, cname = filepath.Split(path)
	cname = cdip.Sprint(cname)
	cdir = cdirp.Sprint(cdir)
	return cdir, cname, cpath
}

func nameC(de DirEntryX) string {
	return de.LSColor().Sprint(de.Name())
}

func linkC(de DirEntryX) string {
	if de.IsLink() {
		link := de.LinkPath()
		dir, name := filepath.Split(link)
		_, err := os.Stat(link)
		if err != nil {
			return cdirp.Sprint(dir) + corp.Sprint(name)
		}
		return cdirp.Sprint(dir) + paw.FileLSColor(link).Sprint(name)
	}
	return ""
}

func nameToLinkC(de DirEntryX) string {
	if de.IsLink() {
		return nameC(de) + cdashp.Sprint(" -> ") + linkC(de)
	} else {
		return nameC(de)
	}
}

func PathToLinkC(de DirEntryX, bgc []color.Attribute) string {
	if bgc == nil {
		dir, name := filepath.Split(de.Path())
		if de.IsLink() {
			return cdirp.Sprint(dir) + de.LSColor().Sprint(name) + cdashp.Sprint(" -> ") + linkC(de)
		} else {
			return cdirp.Sprint(dir) + de.LSColor().Sprint(name)
		}
	} else {
		dir, name := filepath.Split(de.Path())
		var (
			ccdirp  = (*cdirp)
			cnamep  = (*de.LSColor())
			ccdashp = (*cdashp)
		)
		ccdirp.Add(bgc...)
		cnamep.Add(bgc...)
		ccdashp.Add(bgc...)
		if de.IsLink() {
			return ccdirp.Sprint(dir) + cnamep.Sprint(name) + ccdashp.Sprint(" -> ") + PathToLinkC(de, bgc)
		} else {
			return ccdirp.Sprint(dir) + cnamep.Sprint(name)
		}
	}
}

func iNodeC(de DirEntryX) string {
	return cinp.Sprint(de.INode())
}

func GetXattr(path string) ([]string, error) {
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
	var s string
	if !de.IsDir() && de.Mode().IsRegular() {
		s = bytefmt.ByteSize(de.Size())
		// s= paw.FillLeft(bytefmt.ByteSize(uint64(f.Size())), 6)
	} else {
		s = "-"
	}
	return strings.ToLower(s)
}

func sizeC(de DirEntryX) (csize string) {
	ss := sizeS(de)
	if ss == "-" {
		return cdashp.Sprint(ss)
	}
	nss := len(ss)
	sn := fmt.Sprintf("%s", ss[:nss-1])
	su := ss[nss-1:]
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
		su := ss[nss-1:]
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

func modFieldWidths(d *Dir, fields []ViewField) {
	childWidths(d, fields)
	hasFieldNo := false
	for _, fd := range fields {
		if !hasFieldNo && fd&ViewFieldNo != 0 {
			hasFieldNo = true
			break
		}
	}
	if hasFieldNo {
		nd, nf, _ := d.NItems()
		wdidx := GetMaxWidthOf(nd, nf)
		ViewFieldNo.SetWidth(wdidx + 1)
	}
	ViewFieldName.SetWidth(GetViewFieldNameWidthOf(fields))
}

func childWidths(d *Dir, fields []ViewField) {
	ds, _ := d.ReadDirAll()
	var (
		wd, dwd int
	)
	for _, de := range ds {
		for _, fd := range fields {
			wd = de.WidthOf(fd)
			dwd = fd.Width()
			if !de.IsDir() && fd&ViewFieldSize == ViewFieldSize {
				if de.IsCharDev() || de.IsDev() {
					fmajor := ViewFieldMajor.Width()
					fminor := ViewFieldMinor.Width()
					major, minor := de.DevNumber()
					wdmajor := len(fmt.Sprint(major))
					wdminor := len(fmt.Sprint(minor))
					ViewFieldMajor.SetWidth(paw.MaxInt(fmajor, wdmajor))
					ViewFieldMinor.SetWidth(paw.MaxInt(fminor, wdminor))
					wd = ViewFieldMajor.Width() +
						ViewFieldMinor.Width() + 1
				}
			}
			width := paw.MaxInt(dwd, wd)
			fd.SetWidth(width)
		}
		if de.IsDir() {
			child := de.(*Dir)
			childWidths(child, fields)
		}
	}

}

func dirSummary(pad string, ndirs int, nfiles int, sumsize int64, wdstty int) string {
	var (
		ss  = bytefmt.ByteSize(sumsize)
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
		ss  = bytefmt.ByteSize(sumsize)
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
func GetRootHeadC(de DirEntryX, wdstty int) string {
	var size int64
	d, isDir := de.(*Dir)
	if isDir {
		size = d.TotalSize()
	}
	var (
		ss  = bytefmt.ByteSize(size)
		nss = len(ss)
		sn  = ss[:nss-1] // fmt.Sprintf("%s", ss[:nss-1])
		su  = strings.ToLower(ss[nss-1:])
	)
	// if pdOpt != nil && pdOpt.File != nil {
	// 	if pdOpt.File.IsLink() {
	// 		root = pdOpt.File.Path
	// 	}
	// }
	// "prompt":   {38, 5, 251, 48, 5, 236}
	chead := cpmpt.Sprint("Root directory: ")
	chead += PathToLinkC(de, bgpmpt)
	chead += cpmpt.Sprint(", size ≈ ")
	chead += cpmptSn.Sprint(sn) + cpmptSu.Sprint(su)
	chead += cpmpt.Sprint(".")

	chead += cpmpt.Sprint(paw.Spaces(wdstty + 1 - paw.StringWidth(paw.StripANSI(chead))))

	// chead := fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, GetColorizedPath(root, ""), GetColorizedSize(size))
	return chead
}

func FprintTotalSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize int64, wdstty int) {

	fmt.Fprintln(w, totalSummary(pad, ndirs, nfiles, sumsize, wdstty))
}

func FprintBanner(w io.Writer, pad string, mark string, length int) {
	banner := cdashp.Sprintf("%s%s\n", pad, strings.Repeat(mark, length))
	fmt.Fprint(w, banner)
}

func FprintXattrs(w io.Writer, wdpad int, xattrs []string) {
	if len(xattrs) < 1 {
		return
	}
	sp := paw.Spaces(wdpad)
	for _, x := range xattrs {
		fmt.Fprintf(w, "%s%v%v\n",
			sp,
			cxbp.Sprint(XattrSymbol),
			cxap.Sprint(x))
	}
}

func GetViewFieldWithoutName(vfields ViewField, de DirEntryX) (meta string, wdmeta int) {
	meta, wdmeta = GetViewFieldWithoutNameA(vfields.Fields(), de)
	return meta, wdmeta
}

func GetViewFieldWithoutNameA(fields []ViewField, de DirEntryX) (meta string, wdmeta int) {
	for _, field := range fields {
		if field&ViewFieldName != 0 {
			continue
		}
		wdmeta += field.Width() + 1
		meta += fmt.Sprintf("%v ", de.FieldC(field))
	}
	return meta, wdmeta
}

func GetViewFieldWidthWithoutName(vfields ViewField) int {
	wds := vfields.Widths()
	wdmeta := paw.SumInts(wds[:len(wds)-1]...) + len(wds) - 1
	return wdmeta
}

func GetViewFieldNameWidth(vfields ViewField) int {
	wdmeta := GetViewFieldWidthWithoutName(vfields)
	return sttyWidth - 2 - wdmeta
}

func GetViewFieldNameWidthOf(fields []ViewField) int {
	wdmeta := 0
	for _, f := range fields {
		if f&ViewFieldName != 0 {
			continue
		}
		wdmeta += f.Width() + 1
	}
	return sttyWidth - 2 - wdmeta
}

func GetMaxWidthOf(a interface{}, b interface{}) int {
	wda := len(cast.ToString(a))
	wdb := len(cast.ToString(b))
	return paw.MaxInt(wda, wdb)
}
