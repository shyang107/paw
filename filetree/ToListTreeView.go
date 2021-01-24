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
		orgpad = pad
		// buf    = f.StringBuilder()
		// w     = f.Writer()
		w     = paw.NewStringBuilder()
		fm    = f.store
		git   = f.GetGitStatus()
		fds   = NewFieldSliceFrom(pfieldKeys, git)
		chead = fds.ColorHeadsString()
		wmeta = fds.MetaHeadsStringWidth()

		wpad   = paw.StringWidth(orgpad)
		wdstty = sttyWidth - 2 - wpad

		rootName = GetColorizedDirName(f.root, "")
		ctdsize  = GetColorizedSize(f.totalSize)
		head     = fmt.Sprintf("Root directory: %v, size ≈ %v", rootName, ctdsize)
	)

	pad = ""

	// buf.Reset()
	modifyFDSWidth(fds, f, wdstty)

	fmt.Fprintln(w, head)
	printBanner(w, "", "=", wdstty)

	files := fm[RootMark]
	file := files[0]
	nfiles := len(files)

	// print root file
	meta := ""
	dinf, _ := f.DirInfo(file)
	switch pdview {
	case PListTreeView:
		chead, wmeta = modifyFDSTreeHead(fds, f)
		fmt.Fprintln(w, chead)
		fds.SetValues(file, git)
		meta = fds.ColorMetaValuesString()
		wmeta = fds.MetaHeadsStringWidth()
	case PTreeView:
		wmeta = 0
	}

	fmt.Fprintln(w, meta, dinf+" "+file.LSColor().Sprint("."))

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

		fds.SetValues(file, git)
		printLTFile(w, level, levelsEnded, edge, f, file, fds, isExtended, wmeta, wdstty)

		if file.IsDir() && len(fm[file.Dir]) > 1 {
			printLTDir(w, level+1, levelsEnded, edge, f, file, git, fds, isExtended, wmeta, wdstty)
		}
	}

	// print end message
	printBanner(w, "", "=", wdstty)
	printTotalSummary(w, pad, f.NDirs(), f.NFiles(), f.totalSize)

	s := paw.PaddingString(w.String(), orgpad)
	s = paw.TrimSpace(s)
	fmt.Fprintln(f.Writer(), s)

	return s // buf.String()
}

func printLTFile(wr io.Writer, level int, levelsEnded []int,
	edge EdgeType, fl *FileList, file *File, fds *FieldSlice, isExtended bool, wmeta, wdstty int) {

	var (
		sb      = paw.NewStringBuilder()
		padMeta = ""
		meta    = ""
	)

	// fds.SetValues(file, git)
	if pdview == PListTreeView {
		meta = fds.ColorMetaValuesString()
		wmeta = fds.MetaHeadsStringWidth()
		padMeta = paw.Spaces(wmeta + 1)
		// 1. print all fields except Name
		fmt.Fprintf(sb, "%s ", meta)
	} else {
		padMeta = paw.Spaces(wmeta)
	}

	aMeta := ""
	for i := 0; i < level; i++ {
		if isEnded(levelsEnded, i) {
			fmt.Fprintf(sb, "%s", paw.Spaces(IndentSize+1))
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

	// 2. print out Name field
	fmt.Fprint(sb, wrapFileString(fl, file, edge, padMeta, wmeta, wdstty))

	if isExtended && len(file.XAttributes) > 0 {
		// 3. print out extended attributes
		fmt.Fprint(sb, xattrLTString(file, edge, padMeta, wmeta, wdstty))
	}

	fmt.Fprint(wr, sb.String())
}

func wrapFileString(fl *FileList, file *File, edge EdgeType, padMeta string, wmeta, wdstty int) string {
	var (
		sb           = paw.NewStringBuilder()
		cdinf, ndinf = "", 0
		name         = file.BaseName
		wname        = paw.StringWidth(name)
		width        = wdstty - wmeta - edgeWidth[edge] - 2
		spmeta       = paw.Spaces(wmeta)
	)
	if file.IsDir() && fl.depth == -1 {
		cdinf, ndinf = fl.DirInfo(file)
		if ndinf > 0 {
			ndinf++
		}
	}
	if ndinf+wname > width {
		cedge := cdashp.Sprint(edge)
		nb := len(paw.Truncate(name, width-ndinf, ""))
		if ndinf == 0 {
			fmt.Fprintln(sb, cedge, file.LSColor().Sprint(name[:nb]))
		} else {
			fmt.Fprintln(sb, cedge, cdinf, file.LSColor().Sprint(name[:nb]))
		}
		switch edge {
		case EdgeTypeMid:
			cedge = cdashp.Sprint(EdgeTypeLink) + SpaceIndentSize
		case EdgeTypeEnd:
			cedge = SpaceIndentSize
		}

		if paw.StringWidth(name[nb:]) <= width {
			switch pdview {
			case PTreeView:
				fmt.Fprintln(sb, padMeta, cedge+file.LSColor().Sprint(name[nb:]))
			default:
				fmt.Fprintln(sb, padMeta+cedge+file.LSColor().Sprint(name[nb:]))
			}
		} else {
			names := paw.WrapToSlice(name[nb:], width)
			for _, v := range names {
				if edge == EdgeTypeMid {
					fmt.Fprintln(sb, spmeta, cedge+file.LSColor().Sprint(v))
				} else {
					fmt.Fprintln(sb, padMeta, cedge+file.LSColor().Sprint(v))
				}
			}
		}
	} else {
		var (
			cedge = cdashp.Sprint(edge)
			cname string
		)
		if ndinf > 0 {
			cname = cdinf + " " + file.ColorBaseName()
		} else {
			cname = file.ColorBaseName()
		}
		fmt.Fprintln(sb, cedge, cname)
	}
	return sb.String()
}

func xattrLTString(file *File, edge EdgeType, padx string, wmeta, wdstty int) string {
	var (
		sb    = paw.NewStringBuilder()
		nx    = len(file.XAttributes)
		wedge = edgeWidth[edge]
		width = wdstty - wmeta - wedge - 2
	)
	for i := 0; i < nx; i++ {
		var (
			xattr = file.XAttributes[i]
			wdx   = paw.StringWidth(xattr)
			cedge = ""
		)

		switch edge {
		case EdgeTypeMid:
			cedge = padx + cdashp.Sprint(EdgeTypeLink) + SpaceIndentSize
		case EdgeTypeEnd:
			cedge = padx + paw.Spaces(IndentSize+1)
		}
		if wdx <= width {
			fmt.Fprintln(sb, cedge+cxbp.Sprint("@ ")+cxap.Sprint(xattr))
		} else {
			x1 := paw.Truncate(xattr, width-2, "")
			b := len(x1)
			fmt.Fprintln(sb, cedge+cxbp.Sprint("@ ")+cxap.Sprint(x1))
			xs := paw.WrapToSlice(xattr[b:], width)
			for _, v := range xs {
				fmt.Fprintln(sb, cedge+cxbp.Sprint("  ")+cxap.Sprint(v))
			}
		}
	}
	return sb.String()
}

func printLTDir(wr io.Writer, level int, levelsEnded []int,
	edge EdgeType, fl *FileList, file *File, git GitStatus, fds *FieldSlice, isExtended bool, wmeta, wdstty int) {

	var (
		fm     = fl.Map()
		files  = fm[file.Dir]
		nfiles = len(files)
	)

	for i := 1; i < nfiles; i++ {
		file = files[i]
		edge := EdgeTypeMid
		if i == nfiles-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}

		fds.SetValues(file, git)
		printLTFile(wr, level, levelsEnded, edge, fl, file, fds, isExtended, wmeta, wdstty)

		if file.IsDir() && len(fm[file.Dir]) > 1 {
			printLTDir(wr, level+1, levelsEnded, edge, fl, file, git, fds, isExtended, wmeta, wdstty)
		}
	}
}
