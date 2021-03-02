package vfs

import (
	"fmt"
	"io"

	"github.com/shyang107/paw"
)

func (v *VFS) ViewListTree(w io.Writer, fields []ViewField, hasX, onlyTree bool) {
	paw.Logger.Info("[vfs] ViewListTree...")

	cur := v.RootDir()

	if fields == nil {
		fields = DefaultViewFieldSlice
	}
	if onlyTree {
		fields = []ViewField{ViewFieldName}
	} else {
		fields = checkFieldsHasGit(fields, cur.git.NoGit)
		modFieldWidths(v, fields)
		ViewFieldName.SetWidth(GetViewFieldNameWidthOf(fields))
	}

	viewListTree(w, cur, fields, hasX, onlyTree)
	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))

}

func viewListTree(w io.Writer, cur *Dir, fields []ViewField, hasX, onlyTree bool) {
	var (
		wdstty   = sttyWidth - 2
		wdmeta   = 0
		roothead = getRootHeadC(cur, wdstty)
		head     = getPFHeadS(chdp, fields...)
	)

	fmt.Fprintf(w, "%v\n", roothead)
	fprintBanner(w, "", "=", wdstty)

	if hasX {
		for _, fd := range fields {
			if fd&ViewFieldName == ViewFieldName {
				continue
			}
			wdmeta += fd.Width() + 1
		}
	}
	fmt.Fprintf(w, "%v\n", head)

	cdinf, _ := cur.DirInfoC()
	for _, field := range fields {
		if field&ViewFieldName != 0 {
			fmt.Fprintf(w, "%v %v", cdinf, cdip.Sprint("."))
		} else {
			fmt.Fprintf(w, "%v ", cur.FieldC(field))
		}
	}
	fmt.Fprintln(w)
	des, _ := cur.ReadDir(-1)
	cur.ResetIndex()
	// print files in the root dir
	level := 0
	var levelsEnded []int
	for i, de := range des {
		edge := EdgeTypeMid
		if i == len(des)-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}
		if !de.IsDir() {
			vltFile(w, level, levelsEnded, edge, de, fields, hasX, onlyTree, wdstty)
		} else {
			vltFile(w, level, levelsEnded, edge, de, fields, hasX, onlyTree, wdstty)
			cur := de.(*Dir)
			vltDir(w, level+1, levelsEnded, edge, cur, fields, hasX, onlyTree, wdstty)
		}
	}

	// print end message
	fprintBanner(w, "", "=", wdstty)
	tnd, tnf := cur.NItems()
	fprintTotalSummary(w, "", tnd, tnf, cur.TotalSize(), wdstty)
}

func vltFile(w io.Writer, level int, levelsEnded []int, edge EdgeType, de DirEntryX, fields []ViewField, hasX bool, onlyTree bool, wdstty int) {
	var (
		padMeta = ""
		meta    = ""
		wdmeta  int
		cdinf   = ""
		wdinf   int
		cedge   = ""
	)
	if !onlyTree {
		// 1. print all fields except Name
		for _, field := range fields {
			if field&ViewFieldName != 0 {
				continue
			}
			wdmeta += field.Width() + 1
			meta += fmt.Sprintf("%v ", de.FieldC(field))
		}
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
		cedge := cdashp.Sprint(EdgeTypeLink)
		fmt.Fprintf(w, "%s%s", cedge, SpaceIndentSize)
		padMeta += fmt.Sprintf("%s%s", cedge, SpaceIndentSize)
		wdmeta += edgeWidth[EdgeTypeLink] + IndentSize
	}

	xattrs := de.Xattibutes()
	cname := de.FieldC(ViewFieldName)
	if onlyTree && !hasX && len(xattrs) > 0 {
		cname += cdashp.Sprint(" @")
	}
	// 2. print out Name field
	if de.IsDir() {
		cdinf, wdinf = de.(*Dir).DirInfoC()
	}
	if wdinf == 0 {
		fmt.Fprintln(w, cdashp.Sprint(edge), cname)
	} else {
		fmt.Fprintln(w, cdashp.Sprint(edge), cdinf, cname)
	}

	// 3. print out extended attributes
	if hasX && len(xattrs) > 0 {
		switch edge {
		case EdgeTypeMid:
			cedge = padMeta + cdashp.Sprint(EdgeTypeLink) + SpaceIndentSize
		case EdgeTypeEnd:
			cedge = padMeta + paw.Spaces(IndentSize+1)
		}
		for _, x := range xattrs {
			fmt.Fprintf(w, " %s%v%v\n",
				cedge,
				cxbp.Sprint("@ "),
				cxap.Sprint(x))
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

func vltDir(w io.Writer, level int, levelsEnded []int, edge EdgeType, cur *Dir, fields []ViewField, hasX bool, onlyTree bool, wdstty int) {
	des, _ := cur.ReadDir(-1)
	cur.ResetIndex()
	if len(des) < 1 {
		return
	}
	for i, de := range des {
		edge := EdgeTypeMid
		if i == len(des)-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}
		if !de.IsDir() {
			vltFile(w, level, levelsEnded, edge, de, fields, hasX, onlyTree, wdstty)
		} else {
			vltFile(w, level, levelsEnded, edge, de, fields, hasX, onlyTree, wdstty)
			cur := de.(*Dir)
			vltDir(w, level+1, levelsEnded, edge, cur, fields, hasX, onlyTree, wdstty)
		}
	}
}
