package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/xattr"
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
	err := filetree.PrintDir(os.Stdout, root, false, opt, nil, "> ")
	if err != nil {
		paw.Logger.Error(err)
	}

	paw.Info.Println("exPrintDir")
	paw.Warning.Println("exPrintDir")
	paw.Error.Println("exPrintDir")

	r, _ := homedir.Expand("~")
	fl := filetree.NewFileList(r)
	fl.FindFiles(opt.Depth, opt.Ignore)
	for _, dir := range fl.Dirs() {
		for _, file := range fl.Map()[dir] {
			var list []string
			if list, err = xattr.List(file.Path); err != nil {
				paw.Error.Fatal(err)
			}
			if file.IsDir() {
				fmt.Println(file.ColorDirName(r))
			} else {
				fmt.Println(file.ColorBaseName())
			}
			if len(list) > 0 {
				for _, v := range list {
					fmt.Printf("    %q\n", v)
				}
			}
		}
	}

}
