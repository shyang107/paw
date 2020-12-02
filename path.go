package paw

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/shyang107/paw/cast"
	"github.com/shyang107/paw/treeprint"
	"github.com/sirupsen/logrus"
	// log "github.com/sirupsen/logrus"
)

// IsFileExist return true that `fileName` exist or false for not exist
func IsFileExist(fileName string) bool {
	fi, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !fi.IsDir()
}

// IsDirExists return true that `dir` is dir or false for not
func IsDirExists(dir string) bool {
	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return fi.IsDir()
}

// IsPathExists return true that `path` is dir or false for not
func IsPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	}
	// path/to/whatever exists
	return true
}

// GetCurrPath get the current path
func GetCurrPath() string {
	// file, _ := exec.LookPath(os.Args[0])
	// path, _ := filepath.Abs(file)
	// index := strings.LastIndex(path, string(os.PathSeparator))
	// ret := path[:index]
	// return ret
	var abPath string
	_, fileName, _, ok := runtime.Caller(1)
	if ok {
		abPath = filepath.Dir(fileName)
	}
	return abPath
}

// GetAppDir get the current app directory
func GetAppDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"_os_Args_0": os.Args[0],
		}).Warn(err)

	}
	// Logger.Debugln(dir)
	return dir
}

// GetDotDir return the absolute path of "."
func GetDotDir() string {
	// w, _ := homedir.Expand(".")
	w, _ := filepath.Abs(".")
	// Logger.Debugln("get dot working dir", w)
	return w
}

// GetHomeDir get the home directory of user
func GetHomeDir() string {
	Log.Info("get home dir")
	home, _ := homedir.Dir()
	return home
}

// MakeAll check path and create like as `make -p path`
func MakeAll(path string) error {
	// check
	if IsPathExists(path) {
		return nil
	}
	err := os.MkdirAll(path, os.ModePerm) // 0755
	if err != nil {
		return err
	}
	// check again
	if !IsPathExists(path) {
		return fmt.Errorf("Makeall: fail to create %q", path)
	}
	return nil
}

// File : path information
//
// 	Fields:
// 	  `FullPath`: The full path including the folder
// 	  `ShortPath` : The short path is `FullPath` without rootfolder part (replace with './')
// 	  `Folder`: The folder of the file
// 	  `ShortFolder`: The folder of the file without rootfolder part
// 	  `File`: The file name including extension (basename)
// 	  `FileName`: The file name excluding extension
// 	  `Ext`: Extension of the file
type File struct {
	FullPath    string // The full path including the folder
	ShortPath   string // The short path is `FullPath` without rootfolder part
	Folder      string // The folder of the file
	ShortFolder string // The folder of the file without rootfolder part
	File        string // The file name including extension (basename)
	FileName    string // The file name excluding extension
	Ext         string // Extension of the file
}

// ConstructFile construct `paw.File` from string
//
// Example:
// 	path := "/aaa/bbb/ccc/example.xxx"
//  root := "/aaa/"
// 	path => File{
// 		FullPath:    "/aaa/bbb/ccc/example.xxx",
// 		ShortPath:   "bbb/ccc/example.xxx"
// 		File:        "example.xxx",
// 		Folder:      "/aaa/bbb/ccc/",
// 		ShortFolder: "bbb/ccc/",
// 		FileName:    "example",
// 		Ext:         ".xxx",
// 	}
func ConstructFile(path string, root string) File {
	base := filepath.Base(path)
	ext := filepath.Ext(path)
	folder := TrimSuffix(path, base)
	shortFolder, _ := filepath.Rel(root, folder)
	if shortFolder != "." {
		shortFolder = "./" + shortFolder
	}
	shortFolder += "/"
	shortPath := shortFolder + base
	return File{
		FullPath:    path,
		ShortPath:   shortPath,
		File:        base,
		Folder:      folder,
		ShortFolder: shortFolder,
		FileName:    TrimSuffix(base, ext),
		Ext:         ext,
	}
}

// HasFile : Check if file exists in the current directory
func HasFile(filename string) bool {
	if info, err := os.Stat(filename); os.IsExist(err) {
		return !info.IsDir()
	}
	return false
}

// GetFiles :
// 	isRecursive:
// 		false to return []File in `folder`
//		true  to return []File in `folder` and all `subfolders`
func GetFiles(folder string, isRecursive bool) ([]File, error) {
	return GetFilesFunc(folder, isRecursive, func(f File) bool {
		return true
	})
}

// GetFilesString :
// 	isRecursive:
// 		false to return []File in `folder`
//		true  to return []File in `folder` and all `subfolders`
func GetFilesString(folder string, isRecursive bool) ([]string, error) {
	return GetFilesFuncString(folder, isRecursive, func(f File) bool {
		return true
	})
}

var (
	// ExcludePattern is pattern used in GetFilesFunc to exclude some files/folders such as `.git|.gitignore|.DS_Store|$RECYCLE.BIN|desktop.ini|_gsdata_`
	ExcludePattern = `\.git|\.gitignore|\.DS_Store|\$RECYCLE\.BIN|desktop\.ini|_gsdata_`
	// REUsuallyExclude is regexp used in GetFilesFunc to exclude some files/folders such as `.git|.gitignore|.DS_Store|$RECYCLE.BIN|desktop.ini|_gsdata_`
	REUsuallyExclude = regexp.MustCompile(ExcludePattern)
)

// GetFilesFunc :
// 	isRecursive:
// 		false to get []File in `folder`
// 		true  to get []File in `folder` and all `subfolders`
// 	exclude(file) return true to exclude
func GetFilesFunc(folder string, isRecursive bool, exclude func(file File) bool) ([]File, error) {
	var files []File
	// defer checkRootFiles(files, folder)

	if isRecursive {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			file, err := filepath.Abs(path)
			if err != nil {
				return nil
			}

			if !info.IsDir() && info.Mode().IsRegular() {
				f := ConstructFile(file, folder)
				if !exclude(f) {
					files = append(files, f)
				}
			}
			return nil
		})
		files = checkRootFiles(files, folder)
		return files, err
	}

	f, err := os.Open(folder)
	defer f.Close()

	if err != nil {
		// return files, err
		goto END
	}

	if fileinfo, err := f.Readdir(-1); err == nil {
		for _, file := range fileinfo {
			if !file.IsDir() && file.Mode().IsRegular() {
				folder, err := filepath.Abs(folder)
				if err != nil {
					// return files, err
					goto END
				}
				f := ConstructFile(folder+"\\"+file.Name(), folder)
				if !exclude(f) {
					files = append(files, f)
				}
			}
		}
	} else {
		// return files, err
		goto END
	}
END:
	files = checkRootFiles(files, folder)
	return files, err
}

const zeroRootFiles = "«zeroRootFiles»"

func checkRootFiles(files []File, root string) []File {
	nrf := 0
	// fmt.Println("root:", root)
	for _, f := range files {
		// fmt.Println("f.Folder:", f.Folder)
		if strings.EqualFold(f.Folder, root) {
			nrf++
		}
	}
	if nrf == 0 {
		files = append(files, ConstructFile(root+zeroRootFiles, root))
	}
	return files
}

// GetFilesFuncString :
// 	isRecursive:
// 		false to get []File in `folder`
// 		true  to get []File in `folder` and all `subfolders`
// 	filter(file) return true to exclude
func GetFilesFuncString(folder string, isRecursive bool, filter func(file File) bool) ([]string, error) {
	var files []string
	flist, err := GetFilesFunc(folder, isRecursive, filter)
	if err != nil {
		return nil, err
	}
	for _, f := range flist {
		files = append(files, f.FullPath)
	}
	return files, nil
}

// GetNewFilePath change folder in path
//
// Example:
// 	path := "/aaa/bbb/ccc/example.xxx"
// 	path => File{
// 		FullPath: "/aaa/bbb/ccc/example.xxx",
// 		File:     "example.xxx",
// 		Folder:   "/aaa/bbb/ccc/",
// 		FileName: "example",
// 		Ext:      ".xxx",
// 	}
// 	sourceFolder := "/aaa/bbb/"
// 	targetFolder := "ddd/"
// 	return "ddd/ccc/example.xxx"
func GetNewFilePath(file File, sourceFolder, targetFolder string) (string, error) {
	if file.FullPath == "" {
		return "", fmt.Errorf("%s", "Original file is not valid.")
	}
	subfolder := GetSubfolder(file, sourceFolder)
	return targetFolder + subfolder + file.File, nil
}

// GetSubfolder remove `sourceFolder` of path and return the remainder of  subfolder
//
// Example:
// 	path := "/aaa/bbb/ccc/example.xxx"
// 	path => File{
// 		FullPath: "/aaa/bbb/ccc/example.xxx",
// 		File:     "example.xxx",
// 		Folder:   "/aaa/bbb/ccc/",
// 		FileName: "example",
// 		Ext:      ".xxx",
// 	}
// 	sourceFolder := "/aaa/bbb/"
// 	return "ccc/"
func GetSubfolder(file File, sourceFolder string) string {
	return TrimPrefix(file.Folder, sourceFolder)
}

// CountSubfolders count subfolders of `files`
func CountSubfolders(files []File) int {
	folders := make(map[string]int)
	for _, f := range files {
		if _, ok := folders[f.ShortFolder]; !ok {
			if !strings.EqualFold(f.ShortFolder, "./") {
				folders[f.ShortFolder] = 1
			}
		}
	}
	return len(folders)
}

// GrouppingFiles is groupping `files`, first sorted by fullpath then sorted by file name
func GrouppingFiles(files []File) {
	fl := &FileList{files}
	fl.OrderedByFolder()
}

// LessFunc implement Less()
type LessFunc func(p1, p2 *File) bool

// FilesSorter implements the Sort interface, sorting the files within.
type FilesSorter struct {
	files []File
	less  []LessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *FilesSorter) Sort(files []File) {
	ms.files = files
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...LessFunc) *FilesSorter {
	return &FilesSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *FilesSorter) Len() int {
	return len(ms.files)
}

// Swap is part of sort.Interface.
func (ms *FilesSorter) Swap(i, j int) {
	ms.files[i], ms.files[j] = ms.files[j], ms.files[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *FilesSorter) Less(i, j int) bool {
	p, q := &ms.files[i], &ms.files[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

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
	// fl.PrintWithTableFormat(tf, "")
	fl.Print(buf, OPlainTextMode, "", "")
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
// 	// OPlainTextMode : FileList output in plain text mode (default, use PrintPlain())
// 	OPlainTextMode OutputMode = iota
// 	// OTableFormatMode : FileList output in TableFormat mode (use PrintWithTableFormat())
// 	OTableFormatMode
// 	// OTreeMode : FileList output in tree mode (use PrintTree())
// 	OTreeMode
// )

// Print filelist with `head`
func (fl FileList) Print(w io.Writer, mode OutputMode, head, pad string) {
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
		tf.Prepare(w)
		fl.PrintWithTableFormat(tf, head)
	default: // OPlainTextMode
		fl.PrintPlain(w, head, pad)
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

// PrintPlain print out FileList in plain text mode
func (fl FileList) PrintPlain(w io.Writer, head, pad string) {
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

// PrintWithTableFormat print files with `TableFormat` and `head`
func (fl FileList) PrintWithTableFormat(tp *TableFormat, head string) {
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
