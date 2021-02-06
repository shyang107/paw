package filetree

import (
	"fmt"

	"github.com/shyang107/paw"
)

// ToListViewBytes will return the []byte of FileList in list form (like as `exa --header --long --time-style=iso --group --git`)
func (f *FileList) ToListViewBytes(pad string) []byte {
	return []byte(f.ToListView(pad))
}

// ToListView will return the string of FileList in list form (like as `exa --header --long --time-style=iso --group --git`)
func (f *FileList) ToListView(pad string) string {
	return toListView(f, pad, false)
}

// ToListExtendViewBytes will return the []byte of FileList in extend list form (like as `exa --header --long --time-style=iso --group --git -@`)
func (f *FileList) ToListExtendViewBytes(pad string) []byte {
	return []byte(f.ToListExtendView(pad))
}

// ToListExtendView will return the string of FileList in extend list form (like as `exa --header --long --time-style=iso --group --git --@`)
func (f *FileList) ToListExtendView(pad string) string {
	return toListView(f, pad, true)
}

// toListView will return the []byte of FileList in list form (like as `exa --header --long --time-style=iso --group --git`)
func toListView(f *FileList, pad string, isExtended bool) string {
	var (
		w                  = f.stringBuilder
		dirs               = f.Dirs()
		fm                 = f.Map()
		git                = f.GetGitStatus()
		fds                = NewFieldSliceFrom(pfieldKeys, git)
		fNDirs, fNFiles, _ = f.NTotalDirsAndFile()
		nItems             = fNDirs + fNFiles
		ndirs, nfiles      = 0, 0
		wdstty             = sttyWidth - 2 - paw.StringWidth(pad)
		roothead           = getColorizedRootHead(f.root, f.TotalSize(), sttyWidth)
	)

	w.Reset()

	fds.ModifyWidth(f, wdstty)

	fmt.Fprintln(w, roothead)

	if fNDirs == 0 && fNFiles == 0 {
		goto END
	}

	printBanner(w, "", "=", wdstty)
	for _, dir := range dirs {
		if len(fm[dir]) <= 1 {
			continue
		}

		if dir != RootMark {
			if f.depth != 0 {
				fmt.Fprint(w, fm[dir][0].DirNameWrapC("", wdstty))
			}
		}

		fds.PrintHeadRow(w, "")

		for _, file := range fm[dir][1:] {
			if file.IsDir() {
				ndirs++
			} else {
				nfiles++
			}

			fds.SetValues(file, git)
			fds.PrintRow(w, "")

			if isExtended && len(file.XAttributes) > 0 {
				fds.PrintRowXattr(w, "", file.XAttributes, "")
			}
		}

		if f.depth != 0 {
			fmt.Fprintln(w, f.DirSummary(dir))
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
