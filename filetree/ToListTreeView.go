package filetree

import (
	"fmt"
	"io"

	"github.com/shyang107/paw"
)

// ToListTreeViewBytes will return the []byte of `ToListViewTree(pad)` in list+tree form (like as `exa -T(--tree)`)
func (f *FileList) ToListTreeViewBytes(pad string) []byte {
	return []byte(f.ToListTreeView(pad))
}

// ToListTreeView will return the string of FileList in list+tree form (like as `exa -T(--tree)`)
func (f *FileList) ToListTreeView(pad string) string {
	pdview = PListTreeView
	return toListTreeView(f, pad, false)
}

// ToListTreeExtendViewBytes will return the string of `ToListViewTree(pad)` in list+tree form (like as `exa -T(--tree)`)
func (f *FileList) ToListTreeExtendViewBytes(pad string) []byte {
	return []byte(f.ToListTreeExtendView(pad))
}

// ToListTreeExtendView will return the string of FileList in list+tree form (like as `exa -T(--tree)`)
func (f *FileList) ToListTreeExtendView(pad string) string {
	pdview = PListTreeView
	return toListTreeView(f, pad, true)
}

func toListTreeView(f *FileList, pad string, isExtended bool) string {
	var (
		buf = f.StringBuilder()
		w   = f.Writer()
		fm  = f.store
		git = f.GetGitStatus()
		// chead, wmeta = f.GetHead4Meta(pad, urname, gpname, git)
		fds   = NewFieldSliceFrom(pfieldKeys, git)
		chead = fds.ColorHeadsString()
		wmeta = fds.MetaHeadsStringWidth() + paw.StringWidth(pad)
	)
	// wmeta -= pfieldWidthsMap[PFieldName]

	buf.Reset()

	files := fm[RootMark]
	file := files[0]
	nfiles := len(files)

	// print root file
	meta := pad
	switch pdview {
	case PListTreeView:
		chead, wmeta = modifyTreeHead(fds, f, pad)
		printListln(w, chead)
		fds.SetValues(file, git)
		meta += fds.ColorMetaValuesString()
		// tmeta, _ := file.ColorMeta(f.GetGitStatus())
		// meta += tmeta
	case PTreeView:
		dinf, _ := f.DirInfo(file)
		meta += dinf + " "
		wmeta = paw.StringWidth(pad) + 1
	}

	name := fmt.Sprintf("%v (%v)", file.LSColorString("."), file.ColorDirName(""))
	printListln(w, meta+name)

	// print files in the root dir
	level := 0
	var levelsEnded []int
	for i := 1; i < nfiles; i++ {
		file = files[i]
		edge := EdgeTypeMid
		if i == nfiles-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}

		// printLTFile(w, level, levelsEnded, edge, f, file, git, pad, isExtended, wmeta)
		printLTFile(w, level, levelsEnded, edge, f, file, fds, git, pad, isExtended, wmeta)

		if file.IsDir() && len(fm[file.Dir]) > 1 {
			// printLTDir(w, level+1, levelsEnded, edge, f, file, git, pad, isExtended, wmeta)
			printLTDir(w, level+1, levelsEnded, edge, f, file, fds, git, pad, isExtended, wmeta)
		}
	}

	// print end message
	fmt.Fprintln(w)
	printTotalSummary(w, pad, f.NDirs(), f.NFiles(), f.totalSize)

	return buf.String()
}

func printLTFile(wr io.Writer, level int, levelsEnded []int,
	edge EdgeType, fl *FileList, file *File, fds *FieldSlice, git GitStatus, pad string, isExtended bool, wmeta int) {

	fds.SetValues(file, git)
	sb := paw.NewStringBuilder()
	meta := pad
	padMeta := ""
	if pdview == PListTreeView {
		// tmeta, _ := file.ColorMeta(git)
		meta += fds.ColorMetaValuesString()
		// meta += tmeta
		wmeta = fds.MetaValuesStringWidth() + paw.StringWidth(pad)
		// wmeta += lenmeta
		padMeta = paw.Spaces(wmeta)
	} else {
		padMeta = paw.Spaces(wmeta - 1)
	}

	fmt.Fprintf(sb, "%s", meta)

	aMeta := ""
	for i := 0; i < level; i++ {
		if isEnded(levelsEnded, i) {
			fmt.Fprintf(sb, "%v", paw.Spaces(IndentSize+1))
			aMeta += paw.Spaces(IndentSize + 1)
			wmeta += IndentSize + 1
			continue
		}
		cedge := cdashp.Sprint(EdgeTypeLink)
		fmt.Fprintf(sb, "%s%s", cedge, SpaceIndentSize)
		aMeta += fmt.Sprintf("%s%s", cedge, SpaceIndentSize)
		wmeta += edgeWidth[EdgeTypeLink] + IndentSize
	}
	padMeta += aMeta

	fmt.Fprint(sb, wrapFileString(fl, file, edge, padMeta, wmeta))

	if isExtended && len(file.XAttributes) > 0 {
		fmt.Fprint(sb, xattrLTString(file, pad, edge, padMeta, wmeta))
	}
	fmt.Fprint(wr, sb.String())
}

func wrapFileString(fl *FileList, file *File, edge EdgeType, padMeta string, wmeta int) string {
	var (
		sb           = paw.NewStringBuilder()
		cdinf, ndinf = "", 0
		name         = file.BaseName
		wname        = paw.StringWidth(name)
		width        = sttyWidth - wmeta - edgeWidth[edge] - 4
	)

	if file.IsDir() && fl.depth == -1 {
		cdinf, ndinf = fl.DirInfo(file)
		cdinf += " "
		ndinf++
		// width -= ndinf
	}
	if wname > width {
		cedge := cdashp.Sprint(edge)
		// printListln(sb, cedge, cdinf+file.LSColorString(name[:width]))
		nb := len(paw.Truncate(name, width-ndinf, ""))
		printListln(sb, cedge, cdinf+file.LSColorString(name[:nb]))
		switch edge {
		case EdgeTypeMid:
			cedge = padMeta + cdashp.Sprint(EdgeTypeLink) + SpaceIndentSize
		case EdgeTypeEnd:
			cedge = padMeta + SpaceIndentSize
		}
		if paw.StringWidth(name[nb:]) <= width {
			printListln(sb, cedge, file.LSColorString(name[nb:]))
		} else {
			names := paw.WrapToSlice(name[nb:], width)
			for _, v := range names {
				if edge == EdgeTypeMid {
					printListln(sb, cedge+file.LSColorString(v))
				} else {
					printListln(sb, cedge, file.LSColorString(v))
				}
			}
		}
	} else {
		cedge := cdashp.Sprint(edge)
		cname := cdinf + file.ColorBaseName()
		printListln(sb, cedge, cname)
	}
	return sb.String()
}

func xattrLTString(file *File, pad string, edge EdgeType, padx string, wmeta int) string {
	var (
		sb = paw.NewStringBuilder()
		nx = len(file.XAttributes)
	)

	for i := 0; i < nx; i++ {
		var (
			xattr = file.XAttributes[i]
			wdx   = paw.StringWidth(xattr)
			wdm   = wmeta
		)

		cedge := ""
		switch edge {
		case EdgeTypeMid:
			cedge = padx + cdashp.Sprint(EdgeTypeLink) + SpaceIndentSize
			wdm += edgeWidth[EdgeTypeLink] + IndentSize
		case EdgeTypeEnd:
			cedge = padx + paw.Spaces(IndentSize+1)
			wdm += IndentSize + 1
		}
		width := sttyWidth - 2 - wdm - 2
		if wdx <= width {
			printListln(sb, pad+cedge+cdashp.Sprint("@"), cxp.Sprint(xattr))
		} else {
			x1 := paw.Truncate(xattr, width-2, "")
			b := len(x1)
			printListln(sb, pad+cedge+cdashp.Sprint("@"), cxp.Sprint(xattr[:b]))
			xs := paw.WrapToSlice(xattr[b:], width)
			for _, v := range xs {
				printListln(sb, pad+cedge+" ", cxp.Sprint(v))
			}
		}
	}
	return sb.String()
}

func printLTDir(wr io.Writer, level int, levelsEnded []int,
	edge EdgeType, fl *FileList, file *File, fds *FieldSlice, git GitStatus, pad string, isExtended bool, wmeta int) {
	fm := fl.Map()
	files := fm[file.Dir]
	nfiles := len(files)

	for i := 1; i < nfiles; i++ {
		file = files[i]
		edge := EdgeTypeMid
		if i == nfiles-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}

		printLTFile(wr, level, levelsEnded, edge, fl, file, fds, git, pad, isExtended, wmeta)

		if file.IsDir() && len(fm[file.Dir]) > 1 {
			printLTDir(wr, level+1, levelsEnded, edge, fl, file, fds, git, pad, isExtended, wmeta)
		}
	}
}
