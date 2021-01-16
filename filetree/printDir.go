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

type PDFieldFlag int

const (
	// PFieldINode uses inode field
	PFieldINode PDFieldFlag = 1 << iota
	// PFieldPermissions uses permission field
	PFieldPermissions
	// PFieldLinks uses hard link field
	PFieldLinks
	// PFieldSize uses size field
	PFieldSize
	// PFieldBlocks uses blocks field
	PFieldBlocks
	// PFieldUser uses user field
	PFieldUser
	// PFieldGroup uses group field
	PFieldGroup
	// PFieldModified uses date modified field
	PFieldModified
	// PFieldAccessed uses date accessed field
	PFieldAccessed
	// PFieldCreated uses date created field
	PFieldCreated
	// PFieldGit uses git field
	PFieldGit
	// PFieldName uses name field
	PFieldName
)

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

var (
	fields    = []string{}
	fieldsMap = map[PDFieldFlag]string{
		PFieldINode:       "inode",
		PFieldPermissions: "Permissions",
		PFieldLinks:       "Links",
		PFieldSize:        "Size",
		PFieldBlocks:      "Blocks",
		PFieldUser:        "User",
		PFieldGroup:       "Group",
		PFieldModified:    "Date Modified",
		PFieldCreated:     "Date Created",
		PFieldAccessed:    "Date Accessed",
		PFieldGit:         "Git",
		PFieldName:        "Name",
	}
	fieldWidthsMap = map[PDFieldFlag]int{
		PFieldINode:       paw.MaxInt(8, len(fieldsMap[PFieldINode])),
		PFieldPermissions: paw.MaxInt(11, len(fieldsMap[PFieldPermissions])),
		PFieldLinks:       paw.MaxInt(2, len(fieldsMap[PFieldLinks])),
		PFieldSize:        paw.MaxInt(6, len(fieldsMap[PFieldSize])),
		PFieldBlocks:      paw.MaxInt(6, len(fieldsMap[PFieldBlocks])),
		PFieldUser:        paw.MaxInt(paw.StringWidth(urname), len(fieldsMap[PFieldUser])),
		PFieldGroup:       paw.MaxInt(paw.StringWidth(gpname), len(fieldsMap[PFieldGroup])),
		PFieldModified:    paw.MaxInt(11, len(fieldsMap[PFieldModified])),
		PFieldCreated:     paw.MaxInt(11, len(fieldsMap[PFieldCreated])),
		PFieldAccessed:    paw.MaxInt(11, len(fieldsMap[PFieldAccessed])),
		PFieldGit:         paw.MaxInt(2, len(fieldsMap[PFieldGit])),
		PFieldName:        paw.MaxInt(4, len(fieldsMap[PFieldName])),
	}
	fieldKeys = []PDFieldFlag{}
)

func checkFieldFlag(opt *PrintDirOption) {
	if opt.FieldFlag&PFieldINode != 0 {
		fieldKeys = append(fieldKeys, PFieldINode)
	}

	fieldKeys = append(fieldKeys, PFieldPermissions)

	if opt.FieldFlag&PFieldLinks != 0 {
		fieldKeys = append(fieldKeys, PFieldLinks)
	}

	fieldKeys = append(fieldKeys, PFieldSize)

	if opt.FieldFlag&PFieldBlocks != 0 {
		fieldKeys = append(fieldKeys, PFieldBlocks)
	}

	fieldKeys = append(fieldKeys, PFieldUser)
	fieldKeys = append(fieldKeys, PFieldGroup)

	if opt.FieldFlag&PFieldModified != 0 {
		fieldKeys = append(fieldKeys, PFieldModified)
	}
	if opt.FieldFlag&PFieldCreated != 0 {
		fieldKeys = append(fieldKeys, PFieldCreated)
	}
	if opt.FieldFlag&PFieldAccessed != 0 {
		fieldKeys = append(fieldKeys, PFieldAccessed)
	}

	if opt.FieldFlag&PFieldGit != 0 {
		fieldKeys = append(fieldKeys, PFieldGit)
	}
	// fieldKeys = append(fieldKeys, PFieldGit)
	fieldKeys = append(fieldKeys, PFieldName)

	for _, k := range fieldKeys {
		fields = append(fields, fieldsMap[k])
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
	if file.IsFile() || file.IsLink() {
		git, _ := GetShortGitStatus(file.Dir)
		chead, _ := getColorizedHead("", urname, gpname, git)
		fmt.Fprintf(w, "%sDirectory: %v \n", pad, getColorDirName(file.Dir, ""))
		fmt.Fprintln(w, chead)
		meta, _ := file.ColorMeta(git)
		fmt.Fprintf(w, "%s%s%s\n", pad, meta, file.ColorName())
		return errBreak
	}
	return nil
}
