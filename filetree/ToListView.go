package filetree

import (
	"fmt"
	"path/filepath"

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
		chead                   = f.GetHead4Meta(pad, urname, gpname, git)
		ntdirs, nsdirs, ntfiles = 1, 0, 0
		fNDirs                  = f.NDirs()
		fNFiles                 = f.NFiles()
		bannerWidth             = sttyWidth - 2 - len(pad)
	)
	buf.Reset()

	ctdsize := ByteSize(f.totalSize)

	head := fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, getColorDirName(f.root, ""), cdip.Sprint(ctdsize))
	printListln(w, pad+head)

	if fNDirs == 0 && fNFiles == 0 {
		goto END
	}

	printListln(w, pad+paw.Repeat("=", bannerWidth))

	// if paw.IndexOfString(f.dirs, RootMark) != -1 {
	// 	fmt.Fprintln(w, chead)
	// }
	for i, dir := range dirs {
		sumsize := uint64(0)
		nfiles := 0
		ndirs := 0
		if len(fm[dir]) > 0 {
			if !paw.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					dir, name := filepath.Split(dir)
					cname := cdirp.Sprint(dir) + cdip.Sprint(name)
					printListln(w, pad+cname)
				}
			}
		}
		if len(fm[dir]) > 1 {
			ntdirs++
			printListln(w, pad+chead)
		}
		for _, file := range fm[dir][1:] {
			sntf := ""
			if file.IsDir() {
				ndirs++
				nsdirs++
			} else {
				nfiles++
				ntfiles++
				sumsize += file.Size
			}
			meta, metalength := file.ColorMeta(git)
			meta = pad + meta
			metalength += len(pad)
			nameWidth := sttyWidth - metalength - 3
			name := file.BaseName
			if paw.StringWidth(name) <= nameWidth {
				printListln(w, meta, file.ColorName())
			} else {
				printListln(w, meta, file.LSColorString(file.BaseName[:nameWidth]))
				printListln(w, paw.Spaces(metalength), file.LSColorString(file.BaseName[nameWidth:]))
			}

			if isExtended {
				sp := paw.Spaces(metalength + len(sntf) + 1)
				fmt.Fprint(w, xattrEdgeString(file, sp, metalength))
			}
		}

		if f.depth != 0 {
			if len(fm[dir]) > 1 {
				printDirSummary(w, pad, ndirs, nfiles, sumsize)
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
	// printTotalSummary(w, pad, f.NDirs(), f.NFiles(), f.totalSize)
	printTotalSummary(w, pad, fNDirs, fNFiles, f.totalSize)

	// spew.Dump(dirs)
	return buf.String()
}
