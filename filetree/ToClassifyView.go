package filetree

import (
	"fmt"
	"io"

	"github.com/shyang107/paw"
)

// ToClassifyView will return the string of FileList to display type indicator by file names (like as `exa -F` or `exa --classify`)
func (f *FileList) ToClassifyViewString(pad string) string {
	return string(f.ToClassifyView(pad))
}

// ToClassifyView will return the string of FileList to display type indicator by file names (like as `exa -F` or `exa --classify`)
func (f *FileList) ToClassifyView(pad string) string {
	var (
		buf      = f.StringBuilder()
		w        = f.Writer()
		dirs     = f.Dirs()
		fm       = f.Map()
		wdstty   = sttyWidth - 2 - paw.StringWidth(pad)
		rootName = GetColorizedDirName(f.root, "")
		ctdsize  = f.ColorfulTotalByteSize()
		head     = fmt.Sprintf("%sRoot directory: %v, size ≈ %v", pad, rootName, ctdsize)
		nitems   = f.NDirs() + f.NFiles()
	)
	buf.Reset()

	nsitems := 0
	for i, dir := range dirs {
		if f.depth != 0 {
			if dir == RootMark {
				fmt.Fprintln(w, head)
				printBanner(w, "", "=", wdstty)
			} else {
				fmt.Fprint(w, fm[dir][0].ColorWrapDirName("", wdstty))
			}
		} else {
			fmt.Fprintln(w, head)
			printBanner(w, "", "=", wdstty)
		}

		var (
			files        = fm[dir][1:]
			lens, sumlen = getFileStringWidths(files)
		)

		nsitems += len(files) - 1
		if len(fm[dir]) == 1 {
			goto AFTER
		}
		if sumlen <= wdstty {
			classifyPrintFiles(w, files)
		} else {
			classifyGridPrintFiles(w, files, lens, sumlen, wdstty)
		}
	AFTER:
		if f.depth == 0 {
			break
		} else {
			if i < len(dirs)-1 && nsitems < nitems {
				// fmt.Fprintln(w)
				printBanner(w, "", "-", wdstty)
			}
		}
	}

	printBanner(w, "", "=", wdstty)
	printTotalSummary(w, "", f.NDirs(), f.NFiles(), f.totalSize)

	b := paw.PaddingString(buf.String(), pad)

	return b
}

func classifyGridPrintFiles(w io.Writer, files []*File, lens []int, sumlen int, wdstty int) {

	widths := getFieldWidths(lens, wdstty)
	nFields := len(widths)

	nfolds := len(files) / nFields
	if nfolds*nFields < len(files) {
		nfolds++
	}
	for i := 0; i < nfolds; i++ {
		for iw := 0; iw < nFields; iw++ {
			il := i*nFields + iw
			if il > len(files)-1 {
				break
			}
			name := cgGetFileString(files[il], widths[iw], wdstty)
			fmt.Fprintf(w, "%s", name)
		}
		fmt.Fprintln(w)
	}
}

func cgGetFileString(file *File, width, wdstty int) string {
	var (
		wname = paw.StringWidth(file.BaseName)
		ns    = width - wname
		cname = file.ColorBaseName()
		tail  = ""
	)
	if ns < 0 {
		cname = file.LSColor().Sprint(paw.Wrap(file.BaseName, wdstty))
		ns = 0
	}

	if file.IsDir() || file.IsLink() || len(file.XAttributes) > 0 {
		ws := 0
		if file.IsDir() {
			tail += "/"
			ws++
		}
		if file.IsLink() {
			tail += cdashp.Sprint(">")
			ws++
		}
		if len(file.XAttributes) > 0 {
			tail += cdashp.Sprint("@")
			ws++
		}
		tail += paw.Spaces(ns - ws)
	} else {
		tail = paw.Spaces(ns)
	}

	return cname + tail
}

func getFieldWidths(wds []int, maxwd int) (widths []int) {

	nFields := 1
	for i := len(wds); i > 0; i-- {
		s := paw.SumInts(wds[:i]...)
		if s < maxwd {
			nFields = i
			break
		}
	}
	widths = modifyWidths(wds, nFields, maxwd)
	// fmt.Println("maxwd =", maxwd, "sum(widths) =", paw.SumInts(widths...), len(widths))
	return widths
}

func modifyWidths(wds []int, nFields, maxwd int) (widths []int) {

	widths = make([]int, nFields)
	copy(widths, wds[:nFields])

	if nFields == 0 {
		return []int{sttyWidth - 2}
	}
	nfolds := len(wds) / nFields
	if nfolds*nFields < len(wds) {
		nfolds++
	}
	for i := 0; i < nfolds; i++ {
		for iw := 0; iw < nFields; iw++ {
			il := i*nFields + iw
			if il > len(wds)-1 {
				break
			}
			widths[iw] = paw.MaxInt(widths[iw], wds[il])
		}
	}
	if paw.SumInts(widths...) > maxwd {
		widths = modifyWidths(wds, nFields-1, maxwd)
	}
	return widths
}

func classifyPrintFiles(w io.Writer, files []*File) {

	for _, file := range files {
		cname := file.ColorBaseName()
		if file.IsDir() {
			cname += "/"
		}
		if file.IsLink() {
			cname += cdashp.Sprint(">")
		}
		if len(file.XAttributes) > 0 {
			cname += cdashp.Sprint("@")
		}
		fmt.Fprintf(w, "%s  ", cname)
	}
	fmt.Fprintln(w)
}

// getFileStringWidths will return []int of StringWidth of File.BaseName and summation fo the slice
func getFileStringWidths(files []*File) (leng []int, sum int) {
	sum = 0
	for _, file := range files {
		lenstr := paw.StringWidth(file.BaseName) + 2
		if file.IsDir() {
			lenstr++
		}
		if file.IsLink() {
			lenstr++
		}
		if len(file.XAttributes) > 0 {
			lenstr++
		}
		leng = append(leng, lenstr)
		sum += lenstr
	}
	return leng, sum
}
