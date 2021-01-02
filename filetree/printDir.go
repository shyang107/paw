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
	Depth  int
	OutOpt PrintDirType
	Ignore IgnoreFunc
}

func NewPrintDirOption() *PrintDirOption {
	return &PrintDirOption{
		Depth:  0,
		OutOpt: PListView,
		Ignore: DefaultIgnoreFn,
	}
}

type PrintDirType int

const (
	// PListView is the option of list view using in PrintDir
	PListView PrintDirType = 1 << iota // 1 << 0 which is 00000001
	// PListExtendView is the option of list view icluding extend attributes using in PrintDir
	PListExtendView
	// PTreeView is the option of tree view using in PrintDir
	PTreeView
	// PLevelView is the option of level view using in PrintDir
	PLevelView
	// PTableView is the option of table view using in PrintDir
	PTableView
	// PClassifyView display type indicator by file names (like as `exa -F` or `exa --classify`) in PrintDir
	PClassifyView
	// PListTreeView is the option of combining list & tree view using in PrintDir
	PListTreeView = PListView | PTreeView
)

var pdview PrintDirType

// PrintDirSort defines sorting way view of PrintDir
//
// Defaut:
//  increasing sort by lower name of path
type PDirSortOption struct {
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

// PrintDir will find files using codintion `ignore` func
func PrintDir(w io.Writer, path string, isGrouped bool, opt *PrintDirOption, s *PDirSortOption, pad string) error {
	root, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if s == nil {
		s = &PDirSortOption{
			IsSort:  true,
			SortWay: PDSortByName,
		}
	}

	if opt == nil {
		opt = &PrintDirOption{
			Depth:  0,
			OutOpt: PListView,
			// OutOpt: PListExtendView,
			// OutOpt: PTreeView,
			// OutOpt: PListTreeView,
			// OutOpt: PLevelView,
			// OutOpt: PTableView,
			// OutOpt: PClassifyView,
			Ignore: DefaultIgnoreFn,
		}
	}
	file, err := NewFile(path)
	if err != nil {
		return err
	}
	if file.IsRegular() {
		git, _ := GetShortStatus(file.Dir)
		chead := getColorizedHead("", urname, gpname, git)
		fmt.Fprintf(w, "%sDirectory: %v \n", pad, getDirName(file.Dir, ""))
		fmt.Fprintln(w, chead)
		meta, _ := file.ColorMeta(git)
		fmt.Fprintf(w, "%s%s%s\n", pad, meta, file.ColorBaseName())
		return nil
	}

	if opt.Ignore == nil {
		opt.Ignore = DefaultIgnoreFn
	}

	pdview = opt.OutOpt

	fl := NewFileList(root)
	// fl.IsSort = false
	fl.SetWriters(w)

	fl.IsGrouped = isGrouped

	fl.IsSort = s.IsSort
	if !fl.IsSort {
		goto FIND
	}
	if opt.OutOpt != PTreeView || opt.OutOpt != PListTreeView {
		if s.IsSort {
			switch s.SortWay {
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
FIND:
	fl.FindFiles(opt.Depth, opt.Ignore)

	switch opt.OutOpt {
	case PListView:
		fl.ToListView(pad)
	case PListExtendView:
		fl.ToListExtendView(pad)
	case PTreeView:
		fl.ToTreeView(pad)
	case PListTreeView:
		fl.ToListTreeView(pad)
	case PLevelView:
		fl.ToLevelView(pad)
	case PTableView:
		fl.ToTableView(pad)
	case PClassifyView:
		fl.ToClassifyView(pad)
	default:
		return errors.New("No this option of PrintDir")
	}
	return nil
}
