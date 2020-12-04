package paw

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/shyang107/paw/cast"
	"github.com/shyang107/paw/treeprint"
)

// // Files []File
// type Files []File

// FileList struct{ Files }
type FileList struct{ Files []File }

func (fl FileList) String() string {
	// tf := &TableFormat{
	// 	Fields:    []string{"No.", "Sorted Files"},
	// 	LenFields: []int{5, 75},
	// 	Aligns:    []Align{AlignRight, AlignLeft},
	// 	// Padding:   "# ",
	// }
	buf := new(bytes.Buffer)
	// tf.Prepare(buf)
	// fl.PrintTable(tf, "")
	fl.Fprint(buf, OPlainTextMode, "", "")
	return TrimPrefix(string(buf.Bytes()), "\n")
}

// GetFilesFunc get files with codintion `exclude` func
func (fl *FileList) GetFilesFunc(srcFolder string, isRecursive bool, exclude func(file File) bool) {
	files, err := GetFilesFunc(srcFolder, isRecursive, exclude)
	if err != nil {
		Logger.Error(err)
	}
	fl.Files = files
}

// OrderedByFolder organizes files ordered by Folder and then by file name
func (fl *FileList) OrderedByFolder() {
	byFolder := func(f1, f2 *File) bool {
		return f1.Folder < f2.Folder
	}
	byFileName := func(f1, f2 *File) bool {
		return f1.FileName < f2.FileName
	}
	OrderedBy(byFolder, byFileName).Sort(fl.Files)
}

// // OutputMode : FileList output mode
// type OutputMode uint

// const (
// 	// OPlainTextMode : FileList output in plain text mode (default, use PlainText())
// 	OPlainTextMode OutputMode = iota
// 	// OTableFormatMode : FileList output in TableFormat mode (use PrintTable())
// 	OTableFormatMode
// 	// OTreeMode : FileList output in tree mode (use PrintTree())
// 	OTreeMode
// )

// Fprint filelist with `head`
func (fl FileList) Fprint(w io.Writer, mode OutputMode, head, pad string) {
	switch mode {
	case OTreeMode:
		fl.PrintTree(w, head, pad)
	case OTableFormatMode:
		tf := &TableFormat{
			Fields:    []string{"No.", "Sorted Files"},
			LenFields: []int{5, 75},
			Aligns:    []Align{AlignRight, AlignLeft},
			Padding:   pad,
		}
		fl.FprintTable(w, tf, head)
	default: // OPlainTextMode
		fl.FprintText(w, head, pad)
	}
}

// PrintTree print out FileList in tree mode
func (fl FileList) PrintTree(w io.Writer, head, pad string) {
	fmt.Fprintln(w, PaddingString(head, pad))
	root, rootPath := findRoot(fl.Files)

	fdm, fdk := collectFiles(fl.Files)
	nfd := len(fdk)
	nfl := len(fl.Files)

	tree := treeprint.New()
	for _, fd := range fdk {
		trimfd := trimPath(fd)
		ss := strings.Split(trimfd, "/")
		ns := len(ss)
		// fmt.Printf("%v %d %v\n", trimfd, ns, ss)
		if ns == 1 {
			if len(ss[0]) == 0 {
				if len(fdm[fd]) == 1 && strings.EqualFold(fdm[fd][0], zeroRootFiles) {
					nfl--
					delete(fdm, fd)
				}
				tree.SetMetaValue(fmt.Sprintf("%d (%d directories, %d files)", len(fdm[fd]), nfd-1, nfl))
				// tree.SetValue(root)
				tree.SetValue(fmt.Sprintf("%s\n» root: %s", root, rootPath))
				for _, v := range fdm[fd] {
					tree.AddNode(v)
				}
			} else {
				one := tree.AddMetaBranch(cast.ToString(len(fdm[fd])), ss[0])
				for _, v := range fdm[fd] {
					one.AddNode(v)
				}
			}
			continue
		}
		treend := make([]treeprint.Tree, ns)
		treend[0] = tree.FindByValue(ss[0])
		if treend[0] == nil {
			treend[0] = tree.AddMetaBranch(cast.ToString(len(fdm[fd])), ss[0])
		}
		for i := 1; i < ns; i++ {
			treend[i] = treend[0].FindByValue(ss[i])
			if treend[i] == nil {
				treend[i] = treend[i-1].AddMetaBranch(cast.ToString(len(fdm[fd])), ss[i])
				for _, v := range fdm[fd] {
					treend[i].AddNode(v)
				}
			}
		}
	}
	// fmt.Println("nfd =", nfd, "nfl =", nfl)
	fmt.Fprintln(w, PaddingString(tree.String(), pad))
	fmt.Fprintf(w, "%s%d directories, %d files\n", pad, nfd-1, nfl)
}

func collectFiles(files []File) (fdm map[string][]string, fdk []string) {
	fdm = make(map[string][]string)
	fdk = []string{}
	sfd := ""
	for _, f := range files {
		if !strings.EqualFold(sfd, f.ShortFolder) {
			sfd = f.ShortFolder
			fdm[f.ShortFolder] = []string{}
			fdk = append(fdk, f.ShortFolder)
		}
		fdm[f.ShortFolder] = append(fdm[f.ShortFolder], f.File)
	}
	if !sort.StringsAreSorted(fdk) {
		sort.Strings(fdk)
	}
	return fdm, fdk
}

// func trimPath(path string) string {
// 	mpath := TrimPrefix(path, "./")
// 	mpath = TrimSuffix(mpath, "/")
// 	return mpath
// }

func findRoot(files []File) (root, fullpath string) {
	var (
		folder string
	)
	root = files[0].Folder
	fullpath = files[0].Folder
	for _, f := range files {
		if !strings.EqualFold(folder, f.ShortFolder) {
			folder = f.ShortFolder
			if len(root) > len(folder) {
				root = folder
				fullpath = f.Folder
			}
		}
	}
	return root, fullpath
}

// FprintText print out FileList in plain text mode
func (fl FileList) FprintText(w io.Writer, head, pad string) {
	fmt.Fprintln(w, PaddingString(head, pad))
	fmt.Fprintln(w, pad)
	nSubFolders := CountSubfolders(fl.Files)
	nFiles := len(fl.Files)
	count := 1
	for _, f := range fl.Files {
		if f.File == zeroRootFiles {
			nFiles--
			continue
		}
		fmt.Fprintf(w, "%s%5d %s\n", pad, count, f.FullPath)
		count++
	}
	fmt.Fprintln(w, pad)
	fmt.Fprintf(w, "%s%d directories, %d files\n", pad, nSubFolders, nFiles)
}

// FprintTable print files with `TableFormat` and `head`
func (fl FileList) FprintTable(w io.Writer, tp *TableFormat, head string) {
	tp.Prepare(w)
	tp.SetBeforeMessage(head)
	tp.PrintSart()
	nSubFolders := CountSubfolders(fl.Files)
	nFiles := len(fl.Files)
	oFolder := fl.Files[0].Folder
	gcount := 1
	j := 0
	for i, f := range fl.Files {
		if f.File == zeroRootFiles {
			oFolder = fl.Files[i+1].Folder
			nFiles--
			continue
		}
		if oFolder != f.Folder {
			oFolder = f.Folder
			tp.PrintRow("", fmt.Sprintf("Sum: %d files.", j))
			tp.PrintMiddleSepLine()
			j = 1
			gcount++
		} else {
			j++
		}
		if j == 1 {
			if strings.EqualFold(f.ShortFolder, "./") {
				gcount--
				tp.PrintRow("", fmt.Sprintf("[%d]. source folder (%q)", gcount, f.ShortFolder))
			} else {
				tp.PrintRow("", fmt.Sprintf("[%d]. subfolder: %q", gcount, f.ShortFolder))
			}
		}

		tp.PrintRow(j, f.File)

		if i == len(fl.Files)-1 {
			tp.PrintRow("", fmt.Sprintf("Sum: %d files.", j))
		}
	}

	tp.SetAfterMessage(fmt.Sprintf("%d directories, %d files\n", nSubFolders, nFiles))

	tp.PrintEnd()
}
