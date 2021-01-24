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
		orgpad                  = pad
		w                       = f.stringBuilder
		dirs                    = f.Dirs()
		fm                      = f.Map()
		nDirs                   = f.NDirs()
		nFiles                  = f.NFiles()
		nIterms                 = nDirs + nFiles
		git                     = f.GetGitStatus()
		fds                     = NewFieldSliceFrom(pfieldKeys, git)
		wdmeta                  = fds.MetaHeadsStringWidth() + paw.StringWidth(pad)
		ntdirs, nsdirs, ntfiles = 1, 0, 0
		fNDirs                  = f.NDirs()
		fNFiles                 = f.NFiles()
		wpad                    = paw.StringWidth(orgpad)
		bannerWidth             = sttyWidth - 2 - wpad
		rootName                = getColorDirName(f.root, "")
		ctdsize                 = GetColorizedSize(f.totalSize)
		head                    = fmt.Sprintf("Root directory: %v, size ≈ %v", rootName, ctdsize)
	)

	w.Reset()
	modifyFDSWidth(fds, f, bannerWidth-wpad)

	fmt.Fprintln(w, pad+head)

	if fNDirs == 0 && fNFiles == 0 {
		goto END
	}

	printBanner(w, "", "=", bannerWidth)

	for i, dir := range dirs {
		var (
			sumsize = uint64(0)
			nfiles  = 0
			ndirs   = 0
		)
		if len(fm[dir]) > 0 {
			if !paw.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					fmt.Fprint(w, rowWrapDirName(dir, "", wpad, bannerWidth))
				}
			}
		}
		if len(fm[dir]) > 1 {
			ntdirs++
			fds.PrintHeadRow(w, "")
		}
		for _, file := range fm[dir][1:] {
			if file.IsDir() {
				ndirs++
				nsdirs++
			} else {
				nfiles++
				ntfiles++
				sumsize += file.Size
			}

			fds.SetValues(file, git)
			fds.PrintRow(w, "")
			// fmt.Fprint(w, rowWrapFileName(file, fds, pad, bannerWidth))

			if isExtended && len(file.XAttributes) > 0 {
				sp := paw.Spaces(wdmeta)
				fmt.Fprint(w, xattrEdgeString(file, sp, wdmeta, bannerWidth))
			}
		}

		if f.depth != 0 {
			if len(fm[dir]) > 1 {
				printDirSummary(w, "", ndirs, nfiles, sumsize)
			}
			if i < len(f.dirs)-1 && ndirs+nfiles < nIterms {
				printBanner(w, "", "-", bannerWidth)
			}
		}
	}

	printBanner(w, "", "=", bannerWidth)

END:
	printTotalSummary(w, "", fNDirs, fNFiles, f.totalSize)

	// spew.Dump(dirs)
	str := paw.PaddingString(w.String(), orgpad)
	fmt.Fprintln(f.Writer(), str)

	return str
}
