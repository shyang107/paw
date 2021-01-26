package filetree

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/shyang107/paw"
)

// PrintDirOption is the option of PrintDir
//
// Fields:
// 	Depth:
// 		Depth < 0 : print all files and directories recursively of argument path of PrintDir.
// 		Depth = 0 : print files and directories only in argument path of PrintDir.
// 		Depth > 0 : print files and directories recursively under depth of directory in argument path of PrintDir.
// OutOpt: the view-option of PrintDir
// Call
type PrintDirOption struct {
	Depth     int
	OutOpt    PDViewFlag
	FieldFlag PDFieldFlag
	SortOpt   *PDSortOption
	FiltOpt   *PDFilterOption
	Ignore    IgnoreFunc
}

func NewPrintDirOption() *PrintDirOption {
	return &PrintDirOption{
		Depth:     0,
		OutOpt:    PListView,
		FieldFlag: PFieldModified,
		// SortOpt:
		// FiltOpt:,
		Ignore: DefaultIgnoreFn,
	}
}

var pdOpt *PrintDirOption

type PDViewFlag int

const (
	// PListView is the option of list view using in PrintDir
	PListView PDViewFlag = 1 << iota // 1 << 0 which is 00000001
	// PListExtendView is the option of list view icluding extend attributes using in PrintDir
	PListExtendView
	// PTreeView is the option of tree view using in PrintDir
	PTreeView
	// PTreeExtendView is the option of tree view icluding extend atrribute using in PrintDir
	PTreeExtendView
	// PLevelView is the option of level view using in PrintDir
	PLevelView
	// PLevelExtendView is the option of level view icluding extend attributes using in PrintDir
	PLevelExtendView
	// PTableView is the option of table view using in PrintDir
	PTableView
	// PTableView is the option of table view icluding extend attributes using in PrintDir
	PTableExtendView
	// PClassifyView display type indicator by file names (like as `exa -F` or `exa --classify`) in PrintDir
	PClassifyView
	// PListTreeView is the option of combining list & tree view using in PrintDir
	PListTreeView = PListView | PTreeView
	// PListTreeExtendView is the option of combining list & tree view including extend attribute using in PrintDir
	PListTreeExtendView = PListView | PTreeExtendView
)

var pdview PDViewFlag

// PDSortOption defines sorting way view of PrintDir
//
// Defaut:
//  increasing sort by lower name of path
type PDSortOption struct {
	IsSort  bool
	SortWay PDSortFlag
}

type PDSortFlag int

const (
	PDSort PDSortFlag = 1 << iota
	PDSortReverse
	pdSortKeyName
	pdSortKeyMTime
	pdSortKeySize
	PDSortByName         = PDSort | pdSortKeyName
	PDSortByMtime        = PDSort | pdSortKeyMTime
	PDSortBySize         = PDSort | pdSortKeySize
	PDSortByReverseName  = PDSortByName | PDSortReverse
	PDSortByReverseMtime = PDSortByMtime | PDSortReverse
	PDSortByReverseSize  = PDSortBySize | PDSortReverse
)

type PDFiltFlag int

const (
	PDFiltNoEmptyDir = 1 << iota
	PDFiltJustDirs
	PDFiltJustFiles
	PDFiltJustDirsButNoEmpty     = PDFiltNoEmptyDir | PDFiltJustDirs
	PDFiltJustFilesButNoEmptyDir = PDFiltJustFiles
)

type PDFilterOption struct {
	IsFilt  bool
	FiltWay PDFiltFlag
}

// PrintDir will find files using codintion `ignore` func
func PrintDir(w io.Writer, path string, isGrouped bool, opt *PrintDirOption, pad string) (error, *FileList) {

	root, err := filepath.Abs(path)
	if err != nil {
		return err, nil
	}

	pdOpt = opt
	checkPrintDirOption(pdOpt)

	checkFieldFlag(pdOpt)

	sortOpt := checkSortOpt(pdOpt.SortOpt)

	err = checkAndPrintFile(w, path, pad)
	if err != nil {
		if err == errBreak {
			return nil, nil
		}
		return err, nil
	}

	setIgnoreFn(pdOpt)

	pdview = pdOpt.OutOpt

	fl := setFileList(w, root, isGrouped, sortOpt)
	if !fl.IsSort {
		goto FIND
	}

	checkPDSortOption(fl, pdOpt, sortOpt)

FIND:
	fl.FindFiles(pdOpt.Depth, pdOpt.Ignore)

	cehckAndFiltPrintDirFiltOpt(fl, pdOpt.FiltOpt)

	err = switchFileListView(fl, pdOpt.OutOpt, pad)
	if err != nil {
		return err, nil
	}

	return nil, fl
}

func switchFileListView(fl *FileList, outOpt PDViewFlag, pad string) error {
	switch outOpt {
	case PListView:
		fl.ToListView(pad)
	case PListExtendView:
		fl.ToListExtendView(pad)
	case PTreeView:
		fl.ToTreeView(pad)
	case PTreeExtendView:
		fl.ToTreeExtendView(pad)
	case PListTreeView:
		fl.ToListTreeView(pad)
	case PListTreeExtendView:
		fl.ToListTreeExtendView(pad)
	case PLevelView:
		fl.ToLevelView(pad, false)
	case PLevelExtendView:
		fl.ToLevelView(pad, true)
	case PTableView:
		fl.ToTableView(pad, false)
	case PTableExtendView:
		fl.ToTableView(pad, true)
	case PClassifyView:
		fl.ToClassifyView(pad)
	default:
		return errors.New("No this option of PrintDir")
	}
	return nil
}

func cehckAndFiltPrintDirFiltOpt(fl *FileList, filtOpt *PDFilterOption) {
	if filtOpt != nil && filtOpt.IsFilt {
		switch filtOpt.FiltWay {
		case PDFiltNoEmptyDir: // no empty dir
			flf := NewFileListFilter(fl, []Filter{FiltEmptyDirs})
			flf.Filt()
		case PDFiltJustDirs: // no files
			flf := NewFileListFilter(fl, []Filter{FiltJustDirs})
			flf.Filt()
		case PDFiltJustFiles: // PDFiltJustFilesButNoEmptyDir // no dirs
			flf := NewFileListFilter(fl, []Filter{FiltJustFiles})
			flf.Filt()
		case PDFiltJustDirsButNoEmpty: // no file and no empty dir
			flf := NewFileListFilter(fl, []Filter{FiltEmptyDirs, FiltJustDirs})
			flf.Filt()
		}
	}
}

func checkPDSortOption(fl *FileList, opt *PrintDirOption, sortOpt *PDSortOption) {

	if opt.OutOpt != PTreeView || opt.OutOpt != PListTreeView {
		if sortOpt.IsSort {
			switch sortOpt.SortWay {
			case PDSortByMtime:
				// paw.Info.Println("PDSortByMtime")
				fl.SetFilesSorter(func(fi, fj *File) bool {
					return fi.ModifiedTime().Before(fj.ModifiedTime())
				})
			case PDSortBySize:
				// paw.Info.Println("PDSortBySize")
				fl.SetFilesSorter(func(fi, fj *File) bool {
					if fl.IsGrouped {
						if fi.IsDir() && fj.IsDir() {
							return paw.ToLower(fi.Path) < paw.ToLower(fj.Path)
						}
					}
					return fi.Size < fj.Size
				})
			case PDSortByReverseName:
				// paw.Info.Println("PDSortByReverseName")
				fl.SetFilesSorter(func(fi, fj *File) bool {
					return paw.ToLower(fi.Path) > paw.ToLower(fj.Path)
				})
			case PDSortByReverseMtime:
				// paw.Info.Println("PDSortByReverseMtime")
				fl.SetFilesSorter(func(fi, fj *File) bool {
					return fi.ModifiedTime().After(fj.ModifiedTime())
				})
			case PDSortByReverseSize:
				// paw.Info.Println("PDSortByReverseSize")
				fl.SetFilesSorter(func(fi, fj *File) bool {
					if fl.IsGrouped {
						if fi.IsDir() && fj.IsDir() {
							return paw.ToLower(fi.Path) > paw.ToLower(fj.Path)
						}
					}
					return fi.Size > fj.Size
				})
			default: //case PDSortByName :
				// paw.Info.Println("PDSortByName")
				fl.SetFilesSorter(func(fi, fj *File) bool {
					return paw.ToLower(fi.Path) < paw.ToLower(fj.Path)
				})
			}
		}
	}
}

func setFileList(w io.Writer, root string, isGrouped bool, sortOpt *PDSortOption) *FileList {
	fl := NewFileList(root)
	// fl.IsSort = false
	fl.ResetWriters()
	if w != nil {
		fl.SetWriters(w)
	}
	fl.IsGrouped = isGrouped
	fl.IsSort = sortOpt.IsSort

	return fl
}

func setIgnoreFn(opt *PrintDirOption) {
	if opt.Ignore == nil {
		opt.Ignore = DefaultIgnoreFn
	}
}

func checkFieldFlag(opt *PrintDirOption) {
	if opt.FieldFlag&PFieldINode != 0 {
		pfieldKeys = append(pfieldKeys, PFieldINode)
	}

	pfieldKeys = append(pfieldKeys, PFieldPermissions)

	if opt.FieldFlag&PFieldLinks != 0 {
		pfieldKeys = append(pfieldKeys, PFieldLinks)
	}

	pfieldKeys = append(pfieldKeys, PFieldSize)

	if opt.FieldFlag&PFieldBlocks != 0 {
		pfieldKeys = append(pfieldKeys, PFieldBlocks)
	}

	pfieldKeys = append(pfieldKeys, PFieldUser)
	pfieldKeys = append(pfieldKeys, PFieldGroup)

	if opt.FieldFlag&PFieldModified != 0 {
		pfieldKeys = append(pfieldKeys, PFieldModified)
	}
	if opt.FieldFlag&PFieldCreated != 0 {
		pfieldKeys = append(pfieldKeys, PFieldCreated)
	}
	if opt.FieldFlag&PFieldAccessed != 0 {
		pfieldKeys = append(pfieldKeys, PFieldAccessed)
	}

	if opt.FieldFlag&PFieldGit != 0 {
		pfieldKeys = append(pfieldKeys, PFieldGit)
	}
	// pfieldKeys = append(pfieldKeys, PFieldGit)
	pfieldKeys = append(pfieldKeys, PFieldName)

	for _, k := range pfieldKeys {
		pfields = append(pfields, pfieldsMap[k])
		pfieldWidths = append(pfieldWidths, pfieldWidthsMap[k])
	}
}

func checkPrintDirOption(opt *PrintDirOption) {
	if opt == nil {
		opt = NewPrintDirOption()
		// opt = &PrintDirOption{
		// 	Depth:  0,
		// 	OutOpt: PListView,
		// 	// OutOpt: PListExtendView,
		// 	// OutOpt: PTreeView,
		// 	// OutOpt: PListTreeView,
		// 	// OutOpt: PLevelView,
		// 	// OutOpt: PTableView,
		// 	// OutOpt: PClassifyView,
		// 	FieldFlag: PFieldModified,
		// 	Ignore:    DefaultIgnoreFn,
		// }
	}
}

func checkSortOpt(sortOpt *PDSortOption) *PDSortOption {
	if sortOpt == nil {
		return &PDSortOption{
			IsSort:  true,
			SortWay: PDSortByName,
		}
	}
	return sortOpt
}

var errBreak = errors.New("return nil")

func checkAndPrintFile(w io.Writer, path string, pad string) error {
	file, err := NewFile(path)
	if err != nil {
		return err
	}
	if !file.IsDir() {
		fmt.Fprintf(w, "%sDirectory: %v \n", pad, GetColorizedDirName(file.Dir, ""))
		git, _ := GetShortGitStatus(file.Dir)
		fds := NewFieldSliceFrom(pfieldKeysDefualt, git)
		fl := NewFileList(file.Dir)
		fds.ModifyWidth(fl, sttyWidth-2)
		fds.SetValues(file, git)
		fmt.Fprintln(w, fds.ColorHeadsString())
		// fmt.Fprint(w, rowWrapFileName(file, fds, pad, sttyWidth-2))
		fds.PrintRow(w, pad)
		return errBreak
	}
	return nil
}
