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
		Depth:  0,
		OutOpt: filetree.PListView,
		// OutOpt: filetree.PListExtendView,
		// OutOpt: filetree.PTreeView,
		// OutOpt: filetree.PListTreeView,
		// OutOpt: filetree.PLevelView,
		// OutOpt: filetree.PTableView,
		// OutOpt: filetree.PClassifyView,
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
