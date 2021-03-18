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
	"github.com/shyang107/paw/cast"
)

type Color = color.Color
type Attribute = color.Attribute

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
	curname               = paw.Cuup.Sprint(urname)
	cgpname               = paw.Cgup.Sprint(gpname)
	now                   = time.Now()
	thisYear              = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	SpaceIndentSize       = paw.Spaces(IndentSize)
	sttyHeight, sttyWidth = paw.GetTerminalSize()
)

// ===
type PathReturnType int

const (
	PRTName PathReturnType = iota + 1
	PRTNameToLink
	PRTLink
	PRTPath
	PRTPathToLink
	PRTRelPath
	PRTRelPathToLink
)

type PathToOption struct {
	IsColor    bool
	Bgc        []Attribute
	PathReturn PathReturnType
}

func PathTo(de DirEntryX, opt *PathToOption) (p string) {
	switch opt.PathReturn {
	case PRTName:
		if opt.IsColor {
			p = nameCbg(de, opt.Bgc)
		} else {
			p = de.Name()
		}
	case PRTNameToLink:
		if opt.IsColor {
			p = nameToLinkCbg(de, opt.Bgc)
		} else {
			p = de.NameToLink()
		}
	case PRTLink:
		if opt.IsColor {
			p = linkCbg(de, opt.Bgc)
		} else {
			p = de.LinkPath()
		}
	case PRTPath:
		_, _, p = pathC(de.Path(), opt.Bgc)
		if !opt.IsColor {
			p = paw.StripANSI(p)
		}
	case PRTPathToLink:
		p = pathToLinkC(de, opt.Bgc)
		if !opt.IsColor {
			p = paw.StripANSI(p)
		}
	case PRTRelPath:
		p = relPathC(de.RelPath(), opt.Bgc)
		if !opt.IsColor {
			p = paw.StripANSI(p)
		}
	case PRTRelPathToLink:
		if opt.IsColor {
			_, _, p1 := pathC(de.RelPath(), opt.Bgc)
			c := paw.CloneColor(paw.Cdashp)
			if opt.Bgc != nil {
				c = c.Add(opt.Bgc...)
			}
			p = p1 + c.Sprint(" -> ") + linkCbg(de, opt.Bgc)
		} else {
			p = de.RelPath() + " -> " + de.LinkPath()
		}
	}

	return p
}

func choiceSprintf(opt *PathToOption) (dsprintf, nsprintf func(string, ...interface{}) string) {
	var (
		cdir  = paw.CloneColor(paw.Cdirp)
		cname = paw.CloneColor(paw.Cdip)
	)
	if opt.IsColor {
		if opt.Bgc != nil {
			cdir.Add(opt.Bgc...)
			cname.Add(opt.Bgc...)
		}
		dsprintf = cdir.Sprintf
		nsprintf = cname.Sprintf
	} else {
		dsprintf = fmt.Sprintf
		nsprintf = fmt.Sprintf
	}
	return dsprintf, nsprintf
}

// pathC return color string of path
func pathC(path string, bgc []Attribute) (cdir, cname, cpath string) {
	var (
		cdirp = paw.CloneColor(paw.Cdirp)
		cdip  = paw.CloneColor(paw.Cdip)
	)

	if bgc != nil {
		cdirp = cdirp.Add(bgc...)
		cdip = cdip.Add(bgc...)
	}
	cdir, cname = filepath.Split(path)
	if len(cname) > 0 {
		cdir = cdirp.Sprint(cdir)
		cname = cdip.Sprint(cname)
	} else {
		cdir = cdip.Sprint(cdir)
	}
	cpath = cdir + cname
	return cdir, cname, cpath
}

func nameC(de DirEntryX) string {
	return de.LSColor().Sprint(de.Name())
}
func nameCbg(de DirEntryX, bgc []Attribute) string {
	c := paw.CloneColor(de.LSColor())
	if bgc != nil {
		c = c.Add(bgc...)
	}
	return c.Sprint(de.Name())
}

func linkC(de DirEntryX) string {
	if de.IsLink() {
		alink := de.LinkPath()
		dir := filepath.Dir(de.Path())
		link := alink
		if filepath.IsAbs(link) { // get rel path from absolute path
			link, _ = filepath.Rel(dir, alink)
		}
		dir, name := filepath.Split(link)
		if _, err := os.Stat(alink); os.IsNotExist(err) {
			fmt.Println(err)
			return paw.Cdirp.Sprint(dir) + paw.Corp.Sprint(name)
		}
		return paw.Cdirp.Sprint(dir) + paw.FileLSColor(alink).Sprint(name)
	}
	return ""
}

func linkCbg(de DirEntryX, bgc []Attribute) string {
	var (
		cdirp = paw.CloneColor(paw.Cdirp)
		cdip  = paw.CloneColor(paw.Cdip)
		corp  = paw.CloneColor(paw.Corp)
		c     *Color
	)

	if bgc != nil {
		cdirp = cdirp.Add(bgc...)
		cdip = cdip.Add(bgc...)
		corp = corp.Add(bgc...)
	}
	if de.IsLink() {
		alink := de.LinkPath()
		dir := filepath.Dir(de.Path())
		link := alink
		if filepath.IsAbs(link) { // get rel path from absolute path
			link, _ = filepath.Rel(dir, alink)
		}
		// else {
		// 	alink = filepath.Join(dir, alink)
		// }
		dir, name := filepath.Split(link)
		if _, err := os.Stat(alink); os.IsNotExist(err) {
			fmt.Println(err)
			return cdirp.Sprint(dir) + corp.Sprint(name)
		}
		c = paw.FileLSColor(alink)
		if bgc != nil {
			c = c.Add(bgc...)
		}
		return cdirp.Sprint(dir) + c.Sprint(name)
	}
	return ""
}

func nameToLinkC(de DirEntryX) string {
	if de.IsLink() {
		return nameC(de) + paw.Cdashp.Sprint(" -> ") + linkC(de)
	} else {
		return nameC(de)
	}
}
func nameToLinkCbg(de DirEntryX, bgc []Attribute) string {
	c := paw.CloneColor(paw.Cdashp)
	if bgc != nil {
		c = c.Add(bgc...)
	}
	if de.IsLink() {
		return nameCbg(de, bgc) + c.Sprint(" -> ") + linkCbg(de, bgc)
	} else {
		return nameCbg(de, bgc)
	}
}

func getPathFromLink(path string) string {
	alink, err := os.Readlink(path)
	if err != nil {
		return err.Error()
	}
	return alink
}

func pathToLinkC(de DirEntryX, bgc []Attribute) (cpath string) {
	var (
		cdirp  = paw.CloneColor(paw.Cdirp)
		cnamep = paw.CloneColor(de.LSColor())
		cdashp = paw.CloneColor(paw.Cdashp)
		// clnamep *color.Color
	)
	if bgc != nil {
		cdirp = cdirp.Add(bgc...)
		cnamep = cnamep.Add(bgc...)
		cdashp = cdashp.Add(bgc...)
	}

	dir, name := filepath.Split(de.Path())
	cpath = cdirp.Sprint(dir) + cnamep.Sprint(name)
	if de.IsLink() {
		_, _, lpath := pathC(de.LinkPath(), bgc)
		cpath += cdashp.Sprint(" -> ") + lpath
	}
	return cpath
}

func iNodeC(de DirEntryX) string {
	return paw.Cinp.Sprint(de.INode())
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
	cxmark := paw.Cdashp.Sprint(string(sperm[ns-1]))
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
		return paw.Cdashp.Sprint(ss)
	}
	nss := len(ss)
	sn := fmt.Sprintf("%s", ss[:nss-1])
	su := ss[nss-1:]
	csize = paw.Csnp.Sprint(sn) + paw.Csup.Sprint(su)
	return csize
}

// func sizeCaligned(de DirEntryX) (csize string) {
// 	var (
// 		ss  = sizeS(de)
// 		nss = len(ss)
// 	)
// 	if ss == "-" {
// 		csize = paw.Cdashp.Sprint("-")
// 	} else {
// 		sn := fmt.Sprintf("%s", ss[:nss-1])
// 		su := ss[nss-1:]
// 		csize = paw.Csnp.Sprint(sn) + paw.Csup.Sprint(su)
// 	}
// 	var (
// 		width = paw.MaxInt(nss, ViewFieldSize.Width())
// 		sp    = paw.Spaces(width - nss)
// 	)
// 	return sp + csize
// }

// func blocksCaligned(de DirEntryX) (cb string) {
// 	var (
// 		ss  = "-"
// 		nss = 1
// 	)
// 	cb = paw.Cdashp.Sprint(ss)
// 	var (
// 		width = paw.MaxInt(nss, ViewFieldBlocks.Width())
// 		sp    = paw.Spaces(width - nss)
// 	)
// 	return sp + cb
// }

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

func GetDexLSColor(de DirEntryX) *Color {
	if de.IsDir() {
		return paw.Cdip
	}

	if de.IsLink() { // os.ModeSymlink
		_, err := os.Readlink(de.Path())
		if err != nil {
			return paw.NewLSColor("or")
		}
		return paw.Clnp
	}

	if de.IsCharDev() { // os.ModeDevice | os.ModeCharDevice
		return paw.Ccdp
	}

	if de.IsDev() { //
		return paw.Cbdp
	}

	if de.IsFIFO() { //os.ModeNamedPipe
		return paw.Cpip
	}
	if de.IsSocket() { //os.ModeSocket
		return paw.Csop
	}

	if de.IsExecutable() && !de.IsDir() {
		return paw.Cexp
	}

	name := de.Name()
	if att, ok := paw.LSColorAttributes[name]; ok {
		return color.New(att...)
	}
	ext := filepath.Ext(name)
	if att, ok := paw.LSColorAttributes[ext]; ok {
		return color.New(att...)
	}
	file := strings.TrimSuffix(name, ext)
	if att, ok := paw.LSColorAttributes[file]; ok {
		return color.New(att...)
	}
	for re, att := range paw.ReExtLSColors {
		if re.MatchString(name) {
			return color.New(att...)
		}
	}

	return paw.Cfip
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

func dirSummary(pad string, ndirs int, nfiles int, sumsize int64, wdstty int) string {
	var (
		ss  = bytefmt.ByteSize(sumsize)
		nss = len(ss)
		sn  = fmt.Sprintf("%s", ss[:nss-1])
		su  = strings.ToLower(ss[nss-1:])
	)
	cndirs := paw.CpmptSn.Sprint(ndirs)
	cnfiles := paw.CpmptSn.Sprint(nfiles)
	csumsize := paw.CpmptSn.Sprint(sn) + paw.CpmptSu.Sprint(su)
	msg := pad +
		cndirs +
		paw.Cpmpt.Sprint(" directories, ") +
		cnfiles +
		paw.Cpmpt.Sprint(" files, size ≈ ") +
		csumsize +
		paw.Cpmpt.Sprint(". ")
	nmsg := paw.StringWidth(paw.StripANSI(msg))
	msg += paw.Cpmpt.Sprint(paw.Spaces(wdstty + 1 - nmsg))
	return msg
}

func FprintDirSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize int64, wdstty int) {
	fmt.Fprintln(w, dirSummary(pad, ndirs, nfiles, sumsize, wdstty))
}

func FprintDirSummaryNoColor(w io.Writer, pad string, ndirs int, nfiles int, sumsize int64, wdstty int) {
	s := dirSummary(pad, ndirs, nfiles, sumsize, wdstty)
	s = paw.StripANSI(s)
	fmt.Fprintln(w, s)
}

func totalSummary(pad string, ndirs int, nfiles int, sumsize int64, wdstty int) string {
	var (
		ss  = bytefmt.ByteSize(sumsize)
		nss = len(ss)
		sn  = fmt.Sprintf("%s", ss[:nss-1])
		su  = strings.ToLower(ss[nss-1:])
	)
	cndirs := paw.CpmptSn.Sprint(ndirs)
	cnfiles := paw.CpmptSn.Sprint(nfiles)
	csumsize := paw.CpmptSn.Sprint(sn) + paw.CpmptSu.Sprint(su)
	summary := pad +
		paw.Cpmpt.Sprint("Accumulated ") +
		cndirs +
		paw.Cpmpt.Sprint(" directories, ") +
		cnfiles +
		paw.Cpmpt.Sprint(" files, total size ≈ ") +
		csumsize +
		paw.Cpmpt.Sprint(".")
	nsummary := paw.StringWidth(paw.StripANSI(summary))
	summary += paw.Cpmpt.Sprint(paw.Spaces(wdstty + 1 - nsummary))
	return summary
}
func GetRootHeadC(de DirEntryX, wdstty int) string {
	var (
		size  int64
		csize string
	)

	if de.IsDir() {
		size = de.(*Dir).TotalSize()
	}
	if size > 0 {
		ss := bytefmt.ByteSize(size)
		nss := len(ss)
		sn := ss[:nss-1] // fmt.Sprintf("%s", ss[:nss-1])
		su := strings.ToLower(ss[nss-1:])
		csize = paw.CpmptSn.Sprint(sn) + paw.CpmptSu.Sprint(su)
	} else {
		csize = paw.CpmptDashp.Sprint("-")
	}
	chead := paw.Cpmpt.Sprint("Root directory: ")
	chead += PathTo(de, &PathToOption{true, paw.EXAColorAttributes["bgpmpt"], PRTPathToLink})
	chead += paw.Cpmpt.Sprint(", size ≈ ")
	chead += csize
	chead += paw.Cpmpt.Sprint(".")
	chead += paw.Cpmpt.Sprint(paw.Spaces(wdstty + 1 - paw.StringWidth(paw.StripANSI(chead))))
	return chead
}

func FprintRelPath(w io.Writer, pad, slevel, cidx, rp string, isBg bool) {
	fmt.Fprintln(w, getRelPath(pad, slevel+" "+cidx, rp, isBg))
}

func getRelPath(pad, slevel, rp string, isBg bool) string {
	var bgc []Attribute
	var (
		cdirp   = paw.CloneColor(paw.Cdirp)
		cdip    = paw.CloneColor(paw.Cdip)
		clevelp = paw.CloneColor(paw.Cfield)
	)

	if isBg {
		bgc = paw.EXAColorAttributes["bgpmpt"]
		cdirp = cdirp.Add(bgc...)
		cdip = cdip.Add(bgc...)
		clevelp = clevelp.Add(bgc...)
	}
	cdir, cname, cpath := pathC(rp, bgc)
	cpath = cdirp.Sprint("./") + cdir + cname
	clevel := clevelp.Sprintf("%s", slevel)
	return fmt.Sprintf("%s%s%v", pad, clevel, cpath)
}
func relPathC(rp string, bgc []Attribute) string {
	var (
		cdirp   = paw.CloneColor(paw.Cdirp)
		cdip    = paw.CloneColor(paw.Cdip)
		clevelp = paw.CloneColor(paw.Cfield)
	)

	if bgc != nil {
		// bgc = paw.EXAColors["bgpmpt"]
		cdirp = cdirp.Add(bgc...)
		cdip = cdip.Add(bgc...)
		clevelp = clevelp.Add(bgc...)
	}
	cdir, cname, cpath := pathC(rp, bgc)
	cpath = cdirp.Sprint("./") + cdir + cname
	return fmt.Sprintf("%v", cpath)
}

func FprintTotalSummary(w io.Writer, pad string, ndirs int, nfiles int, sumsize int64, wdstty int) {
	fmt.Fprintln(w, totalSummary(pad, ndirs, nfiles, sumsize, wdstty))
}

func FprintTotalSummaryNoColor(w io.Writer, pad string, ndirs int, nfiles int, sumsize int64, wdstty int) {
	s := totalSummary(pad, ndirs, nfiles, sumsize, wdstty)
	s = paw.StripANSI(s)
	fmt.Fprintln(w, s)
}

func FprintBanner(w io.Writer, pad string, mark string, length int) {
	banner := paw.Cdashp.Sprintf("%s%s\n", pad, strings.Repeat(mark, length))
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
			paw.Cxbp.Sprint(XattrSymbol),
			paw.Cxap.Sprint(x))
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
	wdmeta := paw.SumIntA(wds[:len(wds)-1]...) + len(wds) - 1
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

func isSkipViewItem(de DirEntryX, isViewNoDirs, isViewNoFiles bool, nitems, curnd, curnf *int, size *int64) bool {
	if de.IsDir() {
		if isViewNoDirs {
			(*nitems)--
			return true
		}
		(*curnd)++
	} else {
		if isViewNoFiles {
			(*nitems)--
			return true
		}
		(*curnf)++
		if de.Mode().IsRegular() {
			(*size) += de.Size()
		}
	}
	return false
}
