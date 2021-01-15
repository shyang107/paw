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
	var (
		buf = f.StringBuilder()
		w   = f.Writer()
		fm  = f.store
		git = f.GetGitStatus()
	)

	buf.Reset()

	files := fm[RootMark]
	file := files[0]
	nfiles := len(files)

	// print root file
	meta := pad
	switch pdview {
	case PListTreeView:
		chead := f.GetHead4Meta(pad, urname, gpname, git)
		printListln(w, chead)
		tmeta, _ := file.ColorMeta(f.GetGitStatus())
		meta += tmeta
	case PTreeView:
		dinf, _ := f.DirInfo(file)
		meta += dinf + " "
	}

	name := fmt.Sprintf("%v (%v)", file.LSColorString("."), file.ColorDirName(""))
	// fmt.Fprintf(w, "%v%v\n", meta, name)
	printListln(w, meta, name)

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

		printLTFile(w, level, levelsEnded, edge, f, file, git, pad, isExtended)

		if file.IsDir() && len(fm[file.Dir]) > 1 {
			printLTDir(w, level+1, levelsEnded, edge, f, file, git, pad, isExtended)
		}
	}

	// print end message
	fmt.Fprintln(w)
	printTotalSummary(w, pad, f.NDirs(), f.NFiles(), f.totalSize)

	return buf.String()
}

func printLTFile(wr io.Writer, level int, levelsEnded []int,
	edge EdgeType, fl *FileList, file *File, git GitStatus, pad string, isExtended bool) {

	sb := new(strings.Builder)
	meta := pad
	wmeta := paw.StringWidth(meta)
	if pdview == PListTreeView {
		tmeta, lenmeta := file.ColorMeta(git)
		meta += tmeta
		wmeta += lenmeta
	}
	fmt.Fprintf(sb, "%v ", meta)

	// awmeta := wmeta
	aMeta := ""
	padMeta := paw.Spaces(wmeta)
	for i := 0; i < level; i++ {
		if isEnded(levelsEnded, i) {
			fmt.Fprintf(sb, "%v", paw.Spaces(IndentSize+1))
			aMeta += paw.Spaces(IndentSize + 1)
			wmeta += (IndentSize + 1)
			continue
		}
		cedge := cdashp.Sprint(EdgeTypeLink)
		fmt.Fprintf(sb, "%s%s", cedge, SpaceIndentSize)
		aMeta += fmt.Sprintf("%s%s", cedge, SpaceIndentSize)
		wmeta += (edgeWidth[EdgeTypeLink] + IndentSize)
	}
	padMeta += aMeta

	fmt.Fprint(sb, wrapFileString(fl, file, edge, wmeta, padMeta))

	if isExtended {
		fmt.Fprint(sb, xattrLTString(file, pad, edge, padMeta))
	}
	fmt.Fprint(wr, sb.String())
}

func wrapFileString(fl *FileList, file *File, edge EdgeType, wmeta int, padMeta string) string {
	sb := new(strings.Builder)
	cdinf, ndinf := "", 0
	if file.IsDir() && fl.depth == -1 {
		cdinf, ndinf = fl.DirInfo(file)
		cdinf += " "
		wmeta += ndinf
	}

	name := file.BaseName
	wname := paw.StringWidth(name)
	if wmeta+wname+edgeWidth[edge]+1 >= sttyWidth {
		end := sttyWidth - wmeta - edgeWidth[edge] - 3
		if ndinf != 0 {
			ndinf++
			end--
		}
		cedge := cdashp.Sprint(edge)
		printListln(sb, cedge, cdinf+file.LSColorString(name[:end]))
		switch edge {
		case EdgeTypeMid:
			cedge = padMeta + cdashp.Sprint(EdgeTypeLink) + SpaceIndentSize
		case EdgeTypeEnd:
			cedge = padMeta + SpaceIndentSize
		}
		printListln(sb, "", cedge+paw.Spaces(ndinf)+file.LSColorString(name[end:]))
	} else {
		cedge := cdashp.Sprint(edge)
		cname := cdinf + file.ColorBaseName()
		// fmt.Fprintf(sb, "%v %v\n", cedge, cname)
		printListln(sb, cedge, cname)
	}
	return sb.String()
}

func xattrLTString(file *File, pad string, edge EdgeType, padx string) string {
	nx := len(file.XAttributes)
	sb := new(strings.Builder)
	if nx > 0 {
		// edge := EdgeTypeMid
		for i := 0; i < nx; i++ {
			cedge := ""
			switch edge {
			case EdgeTypeMid:
				cedge = padx + cdashp.Sprint(EdgeTypeLink) + SpaceIndentSize
			case EdgeTypeEnd:
				cedge = padx + paw.Spaces(IndentSize+1)
			}
			// fmt.Fprintf(sb, "%s%s%s %s\n", pad, cedge, cdashp.Sprint("@"), cxp.Sprint(file.XAttributes[i]))
			printListln(sb, "", pad+cedge+cdashp.Sprint("@"), cxp.Sprint(file.XAttributes[i]))
		}
	}
	return sb.String()
}

func printLTDir(wr io.Writer, level int, levelsEnded []int,
	edge EdgeType, fl *FileList, file *File, git GitStatus, pad string, isExtended bool) {
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

		printLTFile(wr, level, levelsEnded, edge, fl, file, git, pad, isExtended)

		if file.IsDir() && len(fm[file.Dir]) > 1 {
			printLTDir(wr, level+1, levelsEnded, edge, fl, file, git, pad, isExtended)
		}
	}
}
