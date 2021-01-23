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
		// w     = new(bytes.Buffer)
		buf                     = f.stringBuilder
		w                       = f.Writer()
		dirs                    = f.Dirs()
		fm                      = f.Map()
		git                     = f.GetGitStatus()
		fds                     = NewFieldSliceFrom(pfieldKeys, git)
		chead                   = fds.ColorHeadsString()
		wdmeta                  = fds.MetaHeadsStringWidth() + paw.StringWidth(pad)
		ntdirs, nsdirs, ntfiles = 1, 0, 0
		fNDirs                  = f.NDirs()
		fNFiles                 = f.NFiles()
		wpad                    = paw.StringWidth(pad)
		bannerWidth             = sttyWidth - 2 - wpad
		rootName                = getColorDirName(f.root, "")
		ctdsize                 = GetColorizedSize(f.totalSize)
		head                    = fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, rootName, ctdsize)
	)

	buf.Reset()
	modifyFDSWidth(fds, f, bannerWidth)

	fmt.Fprintln(w, pad+head)

	if fNDirs == 0 && fNFiles == 0 {
		goto END
	}

	printBanner(w, pad, "=", bannerWidth)

	for i, dir := range dirs {
		var (
			sumsize = uint64(0)
			nfiles  = 0
			ndirs   = 0
		)
		if len(fm[dir]) > 0 {
			if !paw.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					fmt.Fprint(w, rowWrapDirName(dir, pad, wpad, bannerWidth))
				}
			}
		}
		if len(fm[dir]) > 1 {
			ntdirs++
			modifyFDSWidth(fds, f, bannerWidth-wpad)
			chead = fds.ColorHeadsString()
			wdmeta = fds.MetaHeadsStringWidth()
			fmt.Fprintln(w, pad+chead)
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
			// meta, _ := file.ColorMeta(git)
			fds.SetValues(file, git)
			fmt.Fprint(w, rowWrapFileName(file, fds, pad, bannerWidth))

			if isExtended && len(file.XAttributes) > 0 {
				sp := paw.Spaces(wdmeta)
				fmt.Fprint(w, xattrEdgeString(file, sp, wdmeta, bannerWidth))
			}
		}

		if f.depth != 0 {
			if len(fm[dir]) > 1 {
				printDirSummary(w, pad, ndirs, nfiles, sumsize)
				switch {
				case nsdirs < fNDirs && fNFiles == 0:
					printBanner(w, pad, "-", bannerWidth)
				case nsdirs <= fNDirs && ntfiles < fNFiles:
					printBanner(w, pad, "-", bannerWidth)
				default:
					if i < len(f.dirs)-1 {
						printBanner(w, pad, "-", bannerWidth)
					}
				}
			}
		}
	}

	printBanner(w, pad, "=", bannerWidth)

END:
	// printTotalSummary(w, pad, f.NDirs(), f.NFiles(), f.totalSize)
	printTotalSummary(w, pad, fNDirs, fNFiles, f.totalSize)

	// spew.Dump(dirs)
	return buf.String()
}
