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
		bannerWidth = sttyWidth - 2 - len(pad)
	)
	buf.Reset()

	chead, wmeta := levelHead(wNo)
	// wmeta -= 2 //len(fields)
	// spmeta := paw.Spaces(wmeta - 1)
	spx := paw.Spaces(wmeta)

	ctdsize := ByteSize(f.totalSize)

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
		ppad := pad
		// sntd := ""
		if len(fm[dir]) > 0 {
			if !paw.EqualFold(dir, RootMark) {
				if f.depth != 0 {
					ppad = printLevelWrappedDir(w, fm[dir][0], ppad, i1, i)
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
			printListln(w, fpad+chead)
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
				fmt.Fprint(w, xattrEdgeString(file, fpad+spx, wmeta+len(fpad)))
			}
		}
		if f.depth != 0 {
			if len(fm[dir]) > 1 {
				printDirSummary(w, fpad, ndirs, nfiles, sumsize)
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

func printLevelWrappedFile(w io.Writer, file *File, pad, cjstr string, wmeta int) {
	spmeta := paw.Spaces(wmeta)
	cperm := file.ColorPermission()
	cfsize := file.ColorSize()
	// ctime := file.ColorModifyTime()
	ctime, _ := getColorizedDates(file)
	wpad := paw.StringWidth(pad)
	name := file.Name()
	wname := paw.StringWidth(name)
	if wpad+wmeta+wname <= sttyWidth {
		cname := file.ColorName()
		printListln(w, pad+cjstr, cperm, cfsize, ctime, cname)
	} else {
		end := sttyWidth - wpad - wmeta - 3
		printListln(w, pad+cjstr, cperm, cfsize, ctime, file.LSColorString(name[:end]))
		printListln(w, pad+spmeta, file.LSColorString(name[end:]))
	}
}

func printLevelWrappedDir(w io.Writer, file *File, ppad string, i1, i int) string {
	level := len(file.DirSlice()) - 1 //len(paw.Split(dir, PathSeparator)) - 1
	ppad += paw.Spaces(4 * level)
	slevel := fmt.Sprintf("L%d: ", level)
	istr := fmt.Sprintf("G%-[1]*[2]d", i1, i)
	cistr := cdip.Sprint(istr)
	cistr = slevel + cistr
	dir, name := filepath.Split(file.Dir)
	wppad := paw.StringWidth(ppad)
	wistr := len(slevel) + len(istr)
	wpi := wppad + wistr
	wdir := len(dir)
	wname := paw.StringWidth(name)
	if wpi+wdir+wname > sttyWidth-4 {
		sp := paw.Spaces(wistr + 1)
		end := sttyWidth - wpi - 4
		if len(dir) < end {
			nend := end - len(dir)
			printListln(w, ppad+cistr, "", cdirp.Sprint(dir)+cdip.Sprint(name[:nend]))
			printListln(w, ppad+sp, "", cdip.Sprint(name[nend:]))
		} else {
			printListln(w, ppad+cistr, "", cdirp.Sprint(dir[:end]))
			printListln(w, ppad+sp, "", cdirp.Sprint(dir[end:])+cdip.Sprint(name))
		}
	} else {
		// cname := GetColorizedDirName(dir, f.root)
		cname := cdirp.Sprint(dir) + cdip.Sprint(name)
		printListln(w, ppad+cistr, "", cname)
	}
	return ppad
}

func xattrEdgeString(file *File, pad string, wmeta int) string {
	nx := len(file.XAttributes)
	sb := new(strings.Builder)
	if nx > 0 {
		edge := EdgeTypeMid
		for i := 0; i < nx; i++ {
			wdm := wmeta
			xattr := file.XAttributes[i]
			if i == nx-1 {
				edge = EdgeTypeEnd
			}
			padx := fmt.Sprintf("%s %s ", pad, cdashp.Sprint(edge))
			wdm += edgeWidth[edge] + 2
			wdx := len(xattr)
			if wdm+wdx <= sttyWidth-2 {
				printListln(sb, padx+cxp.Sprint(xattr))
			} else {
				wde := sttyWidth - 2 - wdm
				printListln(sb, padx+cxp.Sprint(xattr[:wde]))
				switch edge {
				case EdgeTypeMid:
					padx = fmt.Sprintf("%s %s ", pad, cdashp.Sprint(EdgeTypeLink)+SpaceIndentSize)
				case EdgeTypeEnd:
					padx = fmt.Sprintf("%s %s ", pad, paw.Spaces(edgeWidth[edge]))
				}
				if len(xattr[wde:]) <= wde {
					printListln(sb, padx+cxp.Sprint(xattr[wde:]))
				} else {
					xattrs := paw.Split(paw.Wrap(xattr[wde:], wde), "\n")
					for _, v := range xattrs {
						printListln(sb, padx+cxp.Sprint(v))
					}
				}
			}
		}
	}
	return sb.String()
}

func levelHead(wNo int) (chead string, width int) {
	sb := new(strings.Builder)
	csb := new(strings.Builder)
	sno := fmt.Sprintf("%-[1]*[2]s", wNo, "No")
	fmt.Fprintf(sb, "%s ", sno)
	fmt.Fprintf(csb, "%s ", chdp.Sprint(sno))
	for _, k := range fieldKeys {
		// field := ""
		switch k {
		case PFieldINode, PFieldLinks, PFieldUser, PFieldGroup, PFieldGit:
			continue
			// case PFieldPermissions: //"Permissions",
			// case PFieldSize: //"Size",
			// case PFieldModified: //"Date Modified",
			// case PFieldCreated: //"Date Created",
			// case PFieldAccessed: //"Date Accessed",
			// case PFieldName: //"Name",
		default:
			field := fmt.Sprintf("%[1]*[2]s", fieldWidthsMap[k], fieldsMap[k])
			fmt.Fprintf(sb, "%s ", field)
			fmt.Fprintf(csb, "%s ", chdp.Sprint(field))
		}
	}
	head := sb.String()
	head = head[:len(head)-1]
	width = paw.StringWidth(head) - fieldWidthsMap[PFieldName] - 1
	chead = csb.String()
	chead = chead[:len(chead)-1]
	return chead, width

	// ssize := fmt.Sprintf("%6s", "Size")

	// cdate := ""
	// sdate := ""
	// for _, v := range fields {
	// 	cdate += chdp.Sprint(v) + " "
	// 	sdate += v + " "
	// }
	// cdate = cdate[:len(cdate)-1]
	// sdate = cdate[:len(sdate)-1]

	// head := fmt.Sprintf("%s %s %s %s %s", sno, "Permissions", ssize, sdate, "Name")
	// chead := fmt.Sprintf("%s %s %s %s %s", chdp.Sprint(sno), chdp.Sprint("Permissions"), chdp.Sprint(ssize), cdate, chdp.Sprint("Name"))
	// return chead, paw.StringWidth(head) + len(fields) - 5
}
