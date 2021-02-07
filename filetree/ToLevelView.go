package filetree

import (
	"fmt"
	"strings"

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
		w                  = f.StringBuilder()
		dirs               = f.Dirs()
		fm                 = f.Map()
		fNDirs, fNFiles, _ = f.NTotalDirsAndFile()
		nItems             = fNDirs + fNFiles
		git                = f.GetGitStatus()
		ndirs, nfiles      = 0, 0
		wdidx              = len(fmt.Sprint(fNDirs))
		wdjdx              = paw.MaxInts(wdidx, len(fmt.Sprint(fNFiles)))
		wNo                = paw.MaxInt(wdidx, wdjdx) + 1
		wdstty             = sttyWidth - 2 - paw.StringWidth(pad)
		roothead           = getColorizedRootHead(f.root, f.TotalSize(), wdstty)
		fds                = NewFieldSliceFrom(pfieldKeys, git)
	)

	w.Reset()

	fdNo := &Field{
		Key:        PFieldNone,
		Name:       "No",
		Width:      wNo,
		Align:      paw.AlignLeft,
		HeadColor:  chdp,
		ValueColor: cdashp,
	}

	fds.Insert(0, fdNo)
	fds.ModifyWidth(f, wdstty)

	fmt.Fprintln(w, roothead)
	printBanner(w, "", "=", wdstty)

	if fNDirs == 0 && fNFiles == 0 {
		goto END
	}

	for i, dir := range dirs {
		ppad := ""
		if len(fm[dir]) > 1 {
			if !strings.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					// level := len(fm[dir][0].DirSlice()) - 1
					thisDir := fm[dir][0]
					level := len(thisDir.GetUpDir().DirSlice())
					slevel := cNop.Sprintf("L%d: ", level)
					ppad += paw.Spaces(3 * level)
					cistr := slevel + cdip.Sprintf("G%-[1]*[2]d", wdidx, i) + " "
					pipad := ppad + cistr
					fmt.Fprint(w, thisDir.DirNameWrapC(pipad, wdstty))
				}
			}
		} else {
			continue
		}

		fds.ModifyWidth(f, wdstty-paw.StringWidth(ppad))

		fds.PrintHeadRow(w, ppad)

		for _, file := range fm[dir][1:] {
			fds.SetValues(file, git)
			cjstr := ""
			if file.IsDir() {
				ndirs++
				cjstr = cdip.Sprintf("D%-[1]*[2]d", wdjdx, ndirs)
			} else {
				nfiles++
				cjstr = cfip.Sprintf("F%-[1]*[2]d", wdjdx, nfiles)
			}
			fdNo.SetValueC(cjstr)

			fds.PrintRow(w, ppad)

			if isExtended && len(file.XAttributes) > 0 {
				fds.PrintRowXattr(w, ppad, file.XAttributes, "")
			}
		}
		if f.depth != 0 {
			fmt.Fprintln(w, ppad+f.DirSummary(dir))

			if ndirs+nfiles < nItems {
				printBanner(w, "", "-", wdstty)
			}
		}
	}

	printBanner(w, "", "=", wdstty)

END:
	fmt.Fprint(w, f.TotalSummary())

	str := paw.PaddingString(w.String(), pad)
	fmt.Fprintln(f.Writer(), str)

	return str
}
