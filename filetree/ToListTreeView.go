package filetree

import (
	"fmt"
	"io"
	"strings"

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
	paw.Logger.Debug("ListTreeView")
	var (
		buf      = f.StringBuilder()
		w        = f.Writer()
		fm       = f.store
		git      = f.GetGitStatus()
		fds      = NewFieldSliceFrom(pdOpt.FieldKeys(), git)
		chead    = fds.HeadsStringC()
		wmeta    = fds.MetaHeadsStringWidth()
		wdstty   = sttyWidth - 2 - paw.StringWidth(pad)
		roothead = getColorizedRootHead(f.root, f.TotalSize(), wdstty)
		dfile    *File
	)

	buf.Reset()

	fds.ModifyWidth(f, wdstty)

	fmt.Fprintln(buf, roothead)
	f.FprintAllErrs(buf, "")
	printBanner(buf, "", "=", wdstty)

	files := fm[RootMark]
	nfiles := len(files)

	file := files[0]

	// print root file
	meta := ""
	cdinf, _ := f.DirInfo(file)
	cdinf += " "
	switch pdview {
	case PListTreeView:
		chead = fds.HeadsStringC()
		wmeta = fds.HeadsStringWidth()
		fmt.Fprintln(buf, chead)
		fds.SetValues(file, git)
		meta = fds.MetaValuesStringC()
		fmt.Fprintln(buf, meta, cdinf+file.DirNameShortC(f.Root()))
	case PTreeView:
		wmeta = 0
		fmt.Fprintln(buf, cdinf+file.DirNameShortC(f.Root()))
	}

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
		if file.IsDir() {
			sdir := file.GetSubDir()
			dfile = f.Map()[sdir][0]
			printLTFile(buf, level, levelsEnded, edge, f, dfile, fds, isExtended, wmeta, wdstty)

			if len(fm[file.Dir]) > 1 {
				printLTDir(buf, level+1, levelsEnded, edge, f, dfile, fds, isExtended, wmeta, wdstty)
			}
		} else {
			printLTFile(buf, level, levelsEnded, edge, f, file, fds, isExtended, wmeta, wdstty)
		}
	}

	// print end message
	printBanner(buf, "", "=", wdstty)

	fmt.Fprint(buf, f.TotalSummary(wdstty))

	str := paw.PaddingString(buf.String(), pad)
	fmt.Fprintln(w, str)

	return str
}

func printLTFile(wr io.Writer, level int, levelsEnded []int,
	edge EdgeType, fl *FileList, file *File, fds *FieldSlice, isExtended bool, wmeta, wdstty int) {

	// paw.Logger.WithFields(logrus.Fields{
	// 	"dir":  file.Dir,
	// 	"name": file.BaseNameToLink(),
	// }).Info()

	var (
		// sb      = new(strings.Builder)
		padMeta = ""
		meta    = ""
	)

	// fds.SetValues(file, git)
	if pdview == PListTreeView {
		meta = fds.MetaValuesStringC()
		wmeta = fds.MetaHeadsStringWidth()
		padMeta = paw.Spaces(wmeta + 1)
		// 1. print all fields except Name
		fmt.Fprintf(wr, "%s ", meta)
	} else {
		padMeta = paw.Spaces(wmeta)
	}

	aMeta := ""
	for i := 0; i < level; i++ {
		if isEnded(levelsEnded, i) {
			fmt.Fprintf(wr, "%s", paw.Spaces(IndentSize+1))
			aMeta += paw.Spaces(IndentSize + 1)
			wmeta += IndentSize + 1
			continue
		}
		cedge := cdashp.Sprint(EdgeTypeLink)
		fmt.Fprintf(wr, "%s%s", cedge, SpaceIndentSize)
		aMeta += fmt.Sprintf("%s%s", cedge, SpaceIndentSize)
		wmeta += edgeWidth[EdgeTypeLink] + IndentSize
	}
	padMeta += aMeta

	// 2. print out Name field
	fmt.Fprint(wr, wrapLTFileString(fl, file, edge, padMeta, wmeta, wdstty))

	if isExtended && len(file.XAttributes) > 0 {
		// 3. print out extended attributes
		fmt.Fprint(wr, xattrLTString(file, edge, padMeta, wmeta, wdstty))
	}

	// fmt.Fprint(wr, sb.String())
}

func wrapLTFileString(fl *FileList, file *File, edge EdgeType, padMeta string, wmeta, wdstty int) string {
	var (
		sb           = new(strings.Builder)
		cdinf, ndinf = "", 0
		name         = file.Name() //file.BaseName
		wname        = paw.StringWidth(name)
		width        = wdstty - wmeta - edgeWidth[edge] - 2
		spmeta       = paw.Spaces(wmeta)
	)
	if file.IsDir() {
		cdinf, ndinf = fl.DirInfo(file)
		ndinf++
	}
	// fmt.Fprintln(sb, "cdinf =", cdinf, "ndinf =", ndinf)

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
			cname = cdinf + " " + file.NameC() //file.BaseNameC()
		} else {
			cname = file.NameC() //file.BaseNameC()
		}
		fmt.Fprintln(sb, cedge, cname)
	}
	return sb.String()
}

func xattrLTString(file *File, edge EdgeType, padx string, wmeta, wdstty int) string {
	var (
		sb    = new(strings.Builder)
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
	edge EdgeType, fl *FileList, file *File, fds *FieldSlice, isExtended bool, wmeta, wdstty int) {

	var (
		git    = fl.GetGitStatus()
		fm     = fl.Map()
		files  = fm[file.Dir]
		nfiles = len(files)
		dfile  *File
	)

	for i := 1; i < nfiles; i++ {
		file = files[i]
		edge := EdgeTypeMid
		if i == nfiles-1 {
			edge = EdgeTypeEnd
			levelsEnded = append(levelsEnded, level)
		}

		fds.SetValues(file, git)
		if file.IsDir() {
			sdir := file.Dir + "/" + file.BaseName
			dfile = fl.Map()[sdir][0]
			printLTFile(wr, level, levelsEnded, edge, fl, dfile, fds, isExtended, wmeta, wdstty)

			if len(fm[file.Dir]) > 1 {
				printLTDir(wr, level+1, levelsEnded, edge, fl, dfile, fds, isExtended, wmeta, wdstty)
			}
		} else {
			printLTFile(wr, level, levelsEnded, edge, fl, file, fds, isExtended, wmeta, wdstty)
		}

	}
}
