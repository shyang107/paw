package filetree

import (
	"fmt"
	"io"
	"path/filepath"
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
		// w     = new(bytes.Buffer)
		buf     = f.StringBuilder()
		w       = f.Writer()
		dirs    = f.Dirs()
		fm      = f.Map()
		fNDirs  = f.NDirs()
		fNFiles = f.NFiles()
		// git     = f.GetGitStatus()
		// chead                   = f.GetHead4Meta(pad, urname, gpname, git)
		ntdirs  = 1
		nsdirs  = 0
		ntfiles = 0
		i1      = len(fmt.Sprint(fNDirs))
		j1      = paw.MaxInts(i1, len(fmt.Sprint(fNFiles)))
		wNo     = paw.MaxInt(j1+1, 2)
		j       = 0
		sppad   = paw.Spaces(4)
		// wperm       = 11
		// wsize       = 6
		// wdate       = 14
		bannerWidth = sttyWidth - 2
	)
	buf.Reset()

	chead, wmeta := levelHead(wNo)
	wmeta -= 4
	// spmeta := paw.Spaces(wmeta - 1)
	spx := paw.Spaces(wmeta)

	ctdsize := ByteSize(f.totalSize)

	head := fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, getColorDirName(f.root, ""), KindLSColorString("di", ctdsize))
	fmt.Fprintln(w, head)

	printBanner(w, pad, "=", bannerWidth)

	if fNDirs == 0 && fNFiles == 0 {
		goto END
	}

	for i, dir := range dirs {
		sumsize := uint64(0)
		nfiles := 0
		ndirs := 0
		ppad := pad
		// sntd := ""
		if len(fm[dir]) > 0 {
			if !paw.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					ppad = printLevelWrappedDir(w, dir, ppad, i1, i)
					if len(fm[dir]) > 1 {
						ntdirs++
					}
				}
			}
		}
		fpad := ppad
		if i > 1 {
			fpad += sppad
		}
		if len(fm[dir]) > 1 {
			fmt.Fprintf(w, "%s%s\n", fpad, chead)
		}
		for _, file := range fm[dir][1:] {
			cjstr := ""
			if file.IsDir() {
				ndirs, nsdirs = ndirs+1, nsdirs+1
				jstr := fmt.Sprintf("D%-[1]*[2]d", j1, nsdirs)
				cjstr = cdip.Sprint(jstr)
			} else {
				nfiles, ntfiles, j = nfiles+1, ntfiles+1, j+1
				sumsize += file.Size
				jstr := fmt.Sprintf("F%-[1]*[2]d", j1, ntfiles)
				cjstr = cfip.Sprint(jstr)
			}
			printLevelWrappedFile(w, file, fpad, cjstr, wmeta)

			if isExtended {
				fmt.Fprint(w, xattrEdgeString(file, fpad+spx))
			}
		}
		if f.depth != 0 {
			if len(fm[dir]) > 1 {
				printDirSummary(w, fpad, ndirs, nfiles, sumsize)
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
	printTotalSummary(w, pad, fNDirs, fNFiles, f.totalSize)
	// spew.Dump(f.dirs)
	return buf.String()
}

func printLevelWrappedFile(w io.Writer, file *File, pad, cjstr string, wmeta int) {
	spmeta := paw.Spaces(wmeta - 1)
	cperm := file.ColorPermission()
	cfsize := file.ColorSize()
	ctime := file.ColorModifyTime()
	wpad := len(pad)
	name := file.Name()
	wname := paw.StringWidth(name)
	if wpad+wmeta+wname <= sttyWidth {
		cname := file.ColorName()
		printFileItem(w, pad, cjstr, cperm, cfsize, ctime, cname)
	} else {
		end := sttyWidth - wpad - wmeta - 2
		printFileItem(w, pad, cjstr, cperm, cfsize, ctime, file.LSColorString(name[:end]))
		printFileItem(w, pad, spmeta, file.LSColorString(name[end:]))
	}
}

func printLevelWrappedDir(w io.Writer, dir, ppad string, i1, i int) string {
	istr := fmt.Sprintf("G%-[1]*[2]d", i1, i)
	cistr := cdip.Sprint(istr)
	level := len(paw.Split(dir, PathSeparator)) - 1
	ppad += paw.Spaces(4 * level)
	slevel := fmt.Sprintf("L%d: ", level)
	cistr = slevel + cistr
	dir, name := filepath.Split(dir)
	wppad := len(ppad)
	wistr := len(slevel) + len(istr)
	wpi := wppad + wistr
	wdir := len(dir)
	wname := paw.StringWidth(name)
	if wpi+wdir+wname > sttyWidth-4 {
		sp := paw.Spaces(wistr + 1)
		end := sttyWidth - wpi - 4
		if len(dir) < end {
			nend := end - len(dir)
			printFileItem(w, ppad, cistr, "", cdirp.Sprint(dir)+cdip.Sprint(name[:nend]))
			printFileItem(w, ppad, sp, "", cdip.Sprint(name[nend:]))
		} else {
			printFileItem(w, ppad, cistr, "", cdirp.Sprint(dir[:end]))
			printFileItem(w, ppad, sp, "", cdirp.Sprint(dir[end:])+cdip.Sprint(name))
		}
	} else {
		// cname := GetColorizedDirName(dir, f.root)
		cname := cdirp.Sprint(dir) + cdip.Sprint(name)
		printFileItem(w, ppad, cistr, "", cname)
	}
	return ppad
}

func xattrEdgeString(file *File, pad string) string {
	nx := len(file.XAttributes)
	sb := new(strings.Builder)
	if nx > 0 {
		edge := EdgeTypeMid
		for i := 0; i < nx; i++ {
			if i == nx-1 {
				edge = EdgeTypeEnd
			}
			sb.WriteString(pad)
			sb.WriteString(cdashp.Sprint(edge) + " ")
			sb.WriteString(cxp.Sprint(file.XAttributes[i]))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func levelHead(wNo int) (string, int) {
	sno := fmt.Sprintf("%-[1]*[2]s", wNo, "No")
	ssize := fmt.Sprintf("%6s", "Size")
	stime := fmt.Sprintf("%14s", "Data Modified")
	head := fmt.Sprintf("%s %s %s %s %s", sno, "Permissions", ssize, stime, "Name")
	chead := fmt.Sprintf("%s %s %s %s %s", chdp.Sprint(sno), chdp.Sprint("Permissions"), chdp.Sprint(ssize), chdp.Sprint(stime), chdp.Sprint("Name"))
	return chead, paw.StringWidth(head)
}

func printFileItem(w io.Writer, pad string, parameters ...string) {
	str := ""
	for _, p := range parameters {
		str += fmt.Sprintf("%v ", p)
	}
	str += "\n"
	fmt.Fprintf(w, "%v%v", pad, str)
	// fmt.Fprintf(w, "%s%s %s %s %s %s %s %s\n", pad, cperm, cfsize, curname, cgpname, cmodTime, cgit, name)
}
