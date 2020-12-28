package main

import (
	"os"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/filetree"
)

func exPrintDir(root string) {
	// opt := filetree.NewPrintDirOption()
	opt := &filetree.PrintDirOption{
		Depth: -1,
		// OutOpt : filetree.PDListView,
		// OutOpt : filetree.PDTreeView,
		// OutOpt: filetree.PDListView | filetree.PDTreeView,
		OutOpt: filetree.PDLevelView,
		// OutOpt : filetree.PDTable,
		// Ignore: filetree.DefaultIgnoreFn,
	}
	err := filetree.PrintDir(os.Stdout, root, opt)
	if err != nil {
		paw.Logger.Error(err)
	}
}
