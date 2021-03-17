package vfs

import (
	"fmt"
	"io"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func (v *VFS) ViewListTree(w io.Writer) {
	VFSViewListTree(w, v)
}

func VFSViewListTree(w io.Writer, v *VFS) {
	paw.Logger.WithFields(logrus.Fields{"View type": v.opt.ViewType}).Debug("view...")

	hasList, hasX := v.hasList_hasX()
	viewListTree(w, v.RootDir(), hasX, hasList)
	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))

}

func viewListTree(w io.Writer, rootdir *Dir, hasX, hasList bool) {
	var (
		vfields = rootdir.opt.ViewFields
		fields  []ViewField
		wdstty  = sttyWidth - 2
		// roothead = GetRootHeadC(rootdir, wdstty)
		// rootpath = PathToLinkC(rootdir, nil)
		rootpath = PathTo(rootdir, &PathToOption{true, nil, PRTPathToLink})
	)
	if hasList {
		fields = vfields.GetModifyWidthsNoGitFields(rootdir)
	} else {
		fields = []ViewField{ViewFieldName}
	}

	// fmt.Fprintf(w, "%v\n", roothead)
	// FprintBanner(w, "", "=", wdstty)

	if hasList {
		// head := vfields.GetHeadFunc(paw.ChoseColorH)
		head := vfields.GetHead(paw.Chdp)
		fmt.Fprintf(w, "%v\n", head)
		fmt.Fprintf(w, "%v", vfields.RowStringXNameC(rootdir))
	}
	cdinf, _ := rootdir.DirInfoC()
	fmt.Fprintf(w, " %v %v\n", cdinf, rootpath)
	// fmt.Fprintf(w, " %v %v\n", cdinf, paw.Cdip.Sprint("."))

	des, _ := rootdir.ReadDirAll()
	// print files in the root dir
	level := 0
	var levelsEnded []int
	for i, de := range des {
		edge := EdgeTypeMid
		if i == len(des)-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}
		if de.IsDir() {
			vltFile(w, level, levelsEnded, edge, de, fields, hasX, hasList, wdstty)
			cur := de.(*Dir)
			vltDir(w, level+1, levelsEnded, edge, cur, fields, hasX, hasList, wdstty)
		} else {
			vltFile(w, level, levelsEnded, edge, de, fields, hasX, hasList, wdstty)
		}
	}

	// print end message
	fmt.Fprintln(w)
	// FprintBanner(w, "", "=", wdstty)
	rootdir.FprintlnSummaryC(w, "", wdstty, true)
	// fmt.Fprintln(w, rootdir.SummaryC("", wdstty, true))
}

func vltFile(w io.Writer, level int, levelsEnded []int, edge EdgeType, de DirEntryX, fields []ViewField, hasX bool, hasList bool, wdstty int) {
	var (
		padMeta = ""
		meta    = ""
		wdmeta  int
		cdinf   = ""
		wdinf   int
		cedge   = ""
	)
	// 1. print all fields except Name
	if hasList {
		meta, wdmeta = GetViewFieldWithoutNameA(fields, de)
		padMeta = paw.Spaces(wdmeta)
	}
	fmt.Fprintf(w, "%s ", meta)

	for i := 0; i < level; i++ {
		if isEnded(levelsEnded, i) {
			fmt.Fprintf(w, "%s ", SpaceIndentSize)
			padMeta += SpaceIndentSize + " "
			wdmeta += IndentSize + 1
			continue
		}
		cedge := paw.Cdashp.Sprint(EdgeTypeLink)
		fmt.Fprintf(w, "%s%s", cedge, SpaceIndentSize)
		padMeta += fmt.Sprintf("%s%s", cedge, SpaceIndentSize)
		wdmeta += edgeWidth[EdgeTypeLink] + IndentSize
	}

	xattrs := de.Xattibutes()
	cname := de.FieldC(ViewFieldName)
	if !hasList && !hasX && len(xattrs) > 0 {
		cname += paw.Cdashp.Sprint("@")
	}
	// 2. print out Name field
	if de.IsDir() {
		cdinf, wdinf = de.(*Dir).DirInfoC()
	}
	if wdinf == 0 {
		fmt.Fprintln(w, paw.Cdashp.Sprint(edge), cname)
	} else {
		fmt.Fprintln(w, paw.Cdashp.Sprint(edge), cdinf, cname)
	}

	// 3. print out extended attributes
	if hasX && len(xattrs) > 0 {
		switch edge {
		case EdgeTypeMid:
			cedge = padMeta + paw.Cdashp.Sprint(EdgeTypeLink) + SpaceIndentSize
		case EdgeTypeEnd:
			cedge = padMeta + paw.Spaces(IndentSize+1)
		}
		for _, x := range xattrs {
			fmt.Fprintf(w, " %s%v%v\n",
				cedge,
				paw.Cxbp.Sprint("@ "),
				paw.Cxap.Sprint(x))
		}
	}
}

func isEnded(levelsEnded []int, level int) bool {
	for _, l := range levelsEnded {
		if l == level {
			return true
		}
	}
	return false
}

func vltDir(w io.Writer, level int, levelsEnded []int, edge EdgeType, cur *Dir, fields []ViewField, hasX bool, hasList bool, wdstty int) {
	des, _ := cur.ReadDirAll()
	if len(des) < 1 {
		return
	}
	for i, de := range des {
		edge := EdgeTypeMid
		if i == len(des)-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}
		if de.IsDir() {
			vltFile(w, level, levelsEnded, edge, de, fields, hasX, hasList, wdstty)
			cur := de.(*Dir)
			vltDir(w, level+1, levelsEnded, edge, cur, fields, hasX, hasList, wdstty)
		} else {
			vltFile(w, level, levelsEnded, edge, de, fields, hasX, hasList, wdstty)
		}
	}
}
