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
		// OutOpt: filetree.PListExtendView,
		// OutOpt: filetree.PTreeView,
		// OutOpt: filetree.PListTreeView,
		OutOpt: filetree.PLevelView,
		// OutOpt: filetree.PTableView,
		// OutOpt: filetree.PClassifyView,
		Ignore: filetree.DefaultIgnoreFn,
	}

	err, _ := filetree.PrintDir(os.Stdout, root, false, opt, ">")
	if err != nil {
		paw.Logger.Error(err)
	}

	// paw.Info.Println("exPrintDir")
	// paw.Warning.Println("exPrintDir")
	// paw.Error.Println("exPrintDir")

	// r, _ := homedir.Expand("~")
	// fl := filetree.NewFileList(r)
	// fl.FindFiles(opt.Depth, opt.Ignore)
	// for _, dir := range fl.Dirs() {
	// 	for _, file := range fl.Map()[dir] {
	// 		if file.IsDir() {
	// 			fmt.Println(file.DirNameC(r))
	// 		} else {
	// 			fmt.Println(file.BaseNameC())
	// 		}
	// 		if len(file.XAttributes) > 0 {
	// 			for _, v := range file.XAttributes {
	// 				fmt.Printf("    %q\n", v)
	// 			}
	// 		}
	// 	}
	// }

	// fis, _ := ioutil.ReadDir(r)
	// for _, fi := range fis {
	// 	path := filepath.Join(root, fi.Name())
	// 	var list []string
	// 	if list, err = xattr.List(path); err != nil {
	// 		paw.Error.Print(err)
	// 		continue
	// 	}
	// 	if len(list) > 0 {
	// 		fmt.Println(fi.Name())
	// 		for _, v := range list {
	// 			vp, _ := xattr.Get(path, v)
	// 			fmt.Printf("    %s %d %d\n", v, len(v), len(vp))

	// 		}
	// 	}
	// }

}
