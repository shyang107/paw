package filetree

import (
	"fmt"

	"github.com/shyang107/paw"
)

// ToLevelViewBytes will return the []byte of FileList in table form
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f *FileList) ToLevelViewBytes(pad string) []byte {
	return []byte(f.ToLevelView(pad, false))
}

// ToLevelExtendViewString will return the string involving extend attribute of FileList in level form
// 	`size` of directory shown in the string, is accumulated size of sub contents
func (f *FileList) ToLevelExtendViewBytes(pad string) []byte {
	return []byte(f.ToLevelView(pad, true))
}

// ToLevelView will return the string of FileList in level form
// 	`size` of directory shown in the returned value, is accumulated size of sub contents
// 	If `isExtended` is true to involve extend attribute
func (f *FileList) ToLevelView(pad string, isExtended bool) string {
	var (
		// w     = new(bytes.Buffer)
		buf     = f.StringBuilder()
		w       = f.Writer()
		dirs    = f.Dirs()
		fm      = f.Map()
		fNDirs  = f.NDirs()
		fNFiles = f.NFiles()
		git     = f.GetGitStatus()
		ntdirs  = 1
		nsdirs  = 0
		ntfiles = 0
		i1      = len(fmt.Sprint(fNDirs))
		j1      = paw.MaxInts(i1, len(fmt.Sprint(fNFiles)))
		wNo     = paw.MaxInt(j1+1, 2)
		j       = 0
		// nopad   = paw.Spaces(4)
		// wperm       = 11
		// wsize       = 6
		// wdate       = 14
		bannerWidth = sttyWidth - 2 - paw.StringWidth(pad)
		fds         = NewFieldSliceFrom(pfieldKeys, git)
		chead       = fds.ColorHeadsString()
		wdmeta      = fds.MetaHeadsStringWidth()
		spmeta      = paw.Spaces(wdmeta)
		// spx     = paw.Spaces(wdmeta)
		ctdsize = ByteSize(f.totalSize)
	)
	buf.Reset()

	fdNo := &Field{
		Key:        PFieldNone,
		Name:       "No",
		Width:      wNo,
		Align:      paw.AlignLeft,
		HeadColor:  chdp,
		ValueColor: cdashp,
	}

	fds.Insert(0, fdNo)

	head := fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, getColorDirName(f.root, ""), KindLSColorString("di", ctdsize))
	// fmt.Fprintln(w, head)
	printListln(w, pad+head)
	printListln(w, pad+paw.Repeat("=", bannerWidth))
	// printBanner(w, pad, "=", bannerWidth)

	if fNDirs == 0 && fNFiles == 0 {
		goto END
	}

	for i, dir := range dirs {
		sumsize := uint64(0)
		nfiles := 0
		ndirs := 0
		// level := len(fm[dir][0].DirSlice()) - 1
		ppad := pad //+ paw.Spaces(4*level)
		// sntd := ""
		if len(fm[dir]) > 0 {
			if !paw.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					level := len(fm[dir][0].DirSlice()) - 1
					slevel := fmt.Sprintf("L%d: ", level)

					ppad += paw.Spaces(4 * level)
					wppad := paw.StringWidth(ppad)

					istr := fmt.Sprintf("G%-[1]*[2]d", i1, i)
					cistr := slevel + cdip.Sprint(istr) + " "
					wistr := paw.StringWidth(slevel) + paw.StringWidth(istr) + 1

					pipad := ppad + cistr
					wpipad := wppad + wistr
					fmt.Fprint(w, rowWrapDirName(dir, pipad, wpipad, bannerWidth))
					// ppad = printLevelWrappedDir(w, fm[dir][0], ppad, i1, i)
					if len(fm[dir]) > 1 {
						ntdirs++
					}
				}
			}
		}

		if len(fm[dir]) > 1 {
			chead, wdmeta = modifyHead(fds, fm[dir], ppad)
			printListln(w, ppad+chead)
		}
		for _, file := range fm[dir][1:] {
			jstr, cjstr := "", ""
			if file.IsDir() {
				ndirs, nsdirs = ndirs+1, nsdirs+1
				jstr = fmt.Sprintf("D%-[1]*[2]d", j1, nsdirs)
				cjstr = cdip.Sprint(jstr)
			} else {
				nfiles, ntfiles, j = nfiles+1, ntfiles+1, j+1
				sumsize += file.Size
				jstr = fmt.Sprintf("F%-[1]*[2]d", j1, ntfiles)
				cjstr = cfip.Sprint(jstr)
			}
			fdNo.SetValue(jstr)
			fdNo.SetColorfulValue(cjstr)
			fds.SetValues(file, git)
			fmt.Fprint(w, rowWrapFileName(file, fds, ppad, bannerWidth))

			if isExtended && len(file.XAttributes) > 0 {
				spmeta = paw.Spaces(wdmeta)
				fmt.Fprint(w, xattrEdgeString(file, ppad+spmeta, wdmeta+len(ppad), bannerWidth))
			}
		}
		if f.depth != 0 {
			if len(fm[dir]) > 1 {
				printDirSummary(w, ppad, ndirs, nfiles, sumsize)
				switch {
				case nsdirs < fNDirs && fNFiles == 0:
					printListln(w, pad+paw.Repeat("-", bannerWidth))
				case nsdirs <= fNDirs && ntfiles < fNFiles:
					printListln(w, pad+paw.Repeat("-", bannerWidth))
				default:
					if i < len(f.dirs)-1 {
						printListln(w, pad+paw.Repeat("-", bannerWidth))
					}
				}
			}
		}
	}

	printListln(w, pad+paw.Repeat("=", bannerWidth))

END:
	printTotalSummary(w, pad, fNDirs, fNFiles, f.totalSize)
	// spew.Dump(f.dirs)
	return buf.String()
}
