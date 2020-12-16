package main

import (
	"fmt"
	"path/filepath"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/filetree"
)

func exFileTree(root string) {
	// readdir(root)
	// walk(root)
	// constructFile(root)
	readDirs(root)
}

func readDirs(root string) {
	root, err := filepath.Abs(root)
	if err != nil {
		paw.Logger.Fatal(err)
	}

	fl := filetree.NewFileList(root)
	// ignore := func(f *filetree.File, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if f.IsDir() && strings.HasPrefix(f.BaseName, ".") {
	// 		return filetree.SkipDir
	// 	}
	// 	if strings.HasPrefix(f.BaseName, ".") {
	// 		return filetree.SkipFile
	// 	}
	// 	// if strings.Contains(f.Dir, ".git") {
	// 	// 	return filetree.SkipFile
	// 	// }
	// 	return nil
	// }
	// fl.FindFiles(-1, ignore)
	fl.FindFiles(-1, nil)
	// spew.Dump(fl.Dirs())
	fmt.Println(fl.ToTreeString("# "))
	// fmt.Println(fl.ToTableString("# "))
	// fmt.Println(fl.ToTextString("# "))
	// fmt.Println(fl)
	// spew.Dump(fl.Map())
}
