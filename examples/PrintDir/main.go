package main

import (
	"os"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/filetree"
)

func main() {
	root := os.Args[1]
	opt := &filetree.PrintDirOption{
		Depth: 0,
		// OutOpt: filetree.PListView,
		// OutOpt: filetree.PListExtendView,
		// OutOpt: filetree.PTreeView,
		// OutOpt: filetree.PListTreeView,
		OutOpt: filetree.PLevelView,
		// OutOpt: filetree.PTableView,
		// OutOpt: filetree.PClassifyView,
		Ignore: filetree.DefaultIgnoreFn,
	}
	sb :=
	w := io.MultiWriter(os.Stdout, )
	err := filetree.PrintDir(os.Stdout, root, false, opt, nil, nil, ">")
	if err != nil {
		paw.Logger.Error(err)
	}
}
