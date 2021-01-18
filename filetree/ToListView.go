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
		fds                     = NewFieldSliceFrom(pfieldKeys, git)
		chead                   = fds.ColorHeadsString()
		wdmeta                  = fds.MetaHeadsStringWidth() + paw.StringWidth(pad)
		ntdirs, nsdirs, ntfiles = 1, 0, 0
		fNDirs                  = f.NDirs()
		fNFiles                 = f.NFiles()
		bannerWidth             = sttyWidth - 2 - len(pad)
		ctdsize                 = ByteSize(f.totalSize)
		head                    = fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, getColorDirName(f.root, ""), cdip.Sprint(ctdsize))
	)

	buf.Reset()

	printListln(w, pad+head)

	if fNDirs == 0 && fNFiles == 0 {
		goto END
	}

	printListln(w, pad+paw.Repeat("=", bannerWidth))

	for i, dir := range dirs {
		var (
			sumsize = uint64(0)
			nfiles  = 0
			ndirs   = 0
		)
		if len(fm[dir]) > 0 {
			if !paw.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					fmt.Fprint(w, wrapDir(dir, pad, bannerWidth))
				}
			}
		}
		if len(fm[dir]) > 1 {
			ntdirs++
			chead, wdmeta = modifyHead(fds, fm[dir], pad)
			printListln(w, pad+chead)
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
			fmt.Fprint(w, wrapFileName(file, fds, pad, bannerWidth))

			if isExtended && len(file.XAttributes) > 0 {
				sp := paw.Spaces(wdmeta)
				fmt.Fprint(w, xattrEdgeString(file, sp, wdmeta))
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

func wrapDir(dirName, pad string, wdlimit int) string {
	var (
		w         = paw.NewStringBuilder()
		dir, name = filepath.Split(dirName)
		wdir      = paw.StringWidth(dir)
		wname     = paw.StringWidth(name)
		wlen      = wdir + wname
	)
	if wlen <= wdlimit {
		var cname = cdirp.Sprint(dir) + cdip.Sprint(name)
		printListln(w, pad+cname)
	} else {
		if wdir <= wdlimit {
			printListln(w, pad+cdirp.Sprint(dir))
		} else {
			var dirs = paw.WrapToSlice(dir, wdlimit)
			printListln(w, pad+cdirp.Sprint(dirs[0]))
			for i := 1; i < len(dirs); i++ {
				printListln(w, pad+cdirp.Sprint(dirs[i]))
			}
		}
		if wname <= wdlimit {
			printListln(w, pad+cdip.Sprint(name))
		} else {
			var names = paw.WrapToSlice(name, wdlimit)
			printListln(w, pad+cdip.Sprint(names[0]))
			for i := 1; i < len(names); i++ {
				printListln(w, pad+cdip.Sprint(names[i]))
			}
		}
	}
	return w.String()
}
