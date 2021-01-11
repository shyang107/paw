package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/filetree"
)

func main() {
	root := os.Args[1]
	fmt.Printf("%q\n", root)
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
	sb := &strings.Builder{}
	w := io.MultiWriter(os.Stdout, sb)
	err := filetree.PrintDir(w, root, false, opt, nil, nil, "> ")
	if err != nil {
		paw.Logger.Error(err)
	}
	fmt.Println(sb)
}
