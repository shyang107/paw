package filetree

import (
	"errors"
	"io"
	"path/filepath"
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
		OutOpt: PDList,
		Ignore: DefaultIgnoreFn,
	}
}

type PrintDirType int

const (
	// PDList is the option of list view using in PrintDir
	PDList PrintDirType = 1 << iota // 1 << 0 which is 00000001
	// PDTree is the option of tree view using in PrintDir
	PDTree // 1 << 1 which is 00000010
	// PDLevel is the option of level view using in PrintDir
	PDLevel // 1 << 2 which is 00000100
	// PDTable is the option of table view using in PrintDir
	PDTable
)

// PrintDir will find files using codintion `ignore` func
func PrintDir(w io.Writer, path string, opt *PrintDirOption) error {
	root, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	fl := NewFileList(root)
	fl.SetWriters(w)
	fl.FindFiles(opt.Depth, opt.Ignore)

	switch opt.OutOpt {
	case PDList:
		fl.ToList("")
	case PDTree:
		fl.ToTree("")
	case PDList | PDTree:
		fl.ToListTree("")
	case PDLevel:
		fl.ToText("")
	case PDTable:
		fl.ToTable("")
	default:
		return errors.New("No this option of PrintDir")
	}
	return nil
}
