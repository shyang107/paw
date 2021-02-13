package main

import (
	"fmt"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/filetree"
)

func main() {
	root := "." //os.Args[1]
	// fmt.Printf("%q\n", root)
	opt := &filetree.PrintDirOption{
		Depth:    0,
		ViewFlag: filetree.PListView,
		// ViewFlag: filetree.PListExtendView,
		// ViewFlag: filetree.PTreeView,
		// ViewFlag: filetree.PListTreeView,
		// ViewFlag: filetree.PLevelView,
		// ViewFlag: filetree.PTableView,
		// ViewFlag: filetree.PClassifyView,
		Ignore: filetree.DefaultIgnoreFn,
	}
	// sb := &strings.Builder{}
	// w := io.MultiWriter(os.Stdout, sb)
	err, fl := filetree.PrintDir(nil, root, false, opt, "> ")
	if err != nil {
		paw.Logger.Error(err)
	}
	fmt.Println(fl.Dump())
	fmt.Println(fl)
}
