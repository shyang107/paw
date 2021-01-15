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
		bannerWidth             = sttyWidth - 2
	)
	buf.Reset()

	ctdsize := ByteSize(f.totalSize)

	head := fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, getColorDirName(f.root, ""), cdip.Sprint(ctdsize))
	fmt.Fprintln(w, head)

	if fNDirs == 0 && fNFiles == 0 {
		goto END
	}

	printBanner(w, pad, "=", bannerWidth)

	// if paw.IndexOfString(f.dirs, RootMark) != -1 {
	// 	fmt.Fprintln(w, chead)
	// }
	for i, dir := range dirs {
		sumsize := uint64(0)
		nfiles := 0
		ndirs := 0
		sntd := ""
		if len(fm[dir]) > 0 {
			if !paw.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					// sntd = KindEXAColorString("dir", fmt.Sprintf("D%d:", ntdirs))
					dir, name := filepath.Split(dir)
					cname := cdirp.Sprint(dir) + cdip.Sprint(name)
					fmt.Fprintf(w, "%s%s%v\n", pad, sntd, cname)
				}
			}
		}
		if len(fm[dir]) > 1 {
			ntdirs++
			fmt.Fprintln(w, chead)
		}
		for _, file := range fm[dir][1:] {
			sntf := ""
			if file.IsDir() {
				ndirs++
				nsdirs++
				// sntf = file.LSColorString(fmt.Sprintf("D%d(%d):", ndirs, nsdirs))
			} else {
				nfiles++
				ntfiles++
				sumsize += file.Size
				// sntf = file.LSColorString(fmt.Sprintf("F%d(%d):", nfiles, ntfiles))
			}
			meta, metalength := file.ColorMeta(git)
			nameWidth := sttyWidth - metalength - 2
			name := file.BaseName
			if len(name) <= nameWidth {
				fmt.Fprintf(w, "%s%s%s%s\n", pad, sntf, meta, file.ColorName())
			} else {
				fmt.Fprintf(w, "%s%s%s%s\n", pad, sntf, meta, file.LSColorString(file.BaseName[:nameWidth]))
				fmt.Fprintf(w, "%s%s%s%s\n", pad, sntf, paw.Spaces(metalength), file.LSColorString(file.BaseName[nameWidth:]))
				// names := paw.Split(paw.Wrap(name, nameWidth), "\n")
				// fmt.Fprintf(w, "%s%s%s%s\n", pad, sntf, meta, file.LSColorString(names[0]))
				// sp := paw.Spaces(metalength)
				// for k := 1; k < len(names); k++ {
				// 	fmt.Fprintf(w, "%s%s%s%s\n", pad, sntf, sp, file.LSColorString(names[k]))
				// }
			}

			if isExtended {
				sp := paw.Spaces(metalength + len(sntf))
				fmt.Fprint(w, xattrEdgeString(file, sp))
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
