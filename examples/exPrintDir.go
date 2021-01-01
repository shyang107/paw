package main

import (
	"os"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/filetree"
)

func exPrintDir(root string) {
	// opt := filetree.NewPrintDirOption()
	opt := &filetree.PrintDirOption{
		Depth: 0,
		// OutOpt: filetree.PListView,
		// OutOpt: filetree.PTreeView,
		// OutOpt: filetree.PListTreeView,
		// OutOpt: filetree.PLevelView,
		// OutOpt: filetree.PTableView,
		OutOpt: filetree.PClassifyView,
		Ignore: filetree.DefaultIgnoreFn,
	}
	err := filetree.PrintDir(os.Stdout, root, opt, "> ")
	if err != nil {
		paw.Logger.Error(err)
	}

	paw.Info.Println("exPrintDir")
	paw.Warning.Println("exPrintDir")
	paw.Error.Println("exPrintDir")
}
