package dfs

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
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

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func GetDexLSColor(de DirEntryX) *color.Color {
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

func getLinkPath(path string) string {
	alink, err := os.Readlink(path)
	if err != nil {
		return err.Error()
	}
	return alink
}
