package main

import (
	"os"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/filetree"
)

func exPrintDir(root string) {
	opt := filetree.NewPrintDirOption()
	opt.Depth = -1
	opt.OutOpt = filetree.PDList
	// opt.OutOpt = filetree.PDTree
	// opt.OutOpt = filetree.PDList | filetree.PDTree
	// opt.OutOpt = filetree.PDLevel
	// opt.OutOpt = filetree.PDTable
	err := filetree.PrintDir(os.Stdout, root, opt)
	if err != nil {
		paw.Logger.Error(err)
	}
}
