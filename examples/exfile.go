package main

import (
	"fmt"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/shyang107/paw/filetree"
)

func exFile() {
	// newFile()
	fileSomeMethods()
}

func fileSomeMethods() {
	path := []string{
		"/Users/shyang/go/src/github.com/shyang107/paw/filetree",
		"/Users/shyang/go/src/github.com/shyang107/paw/filetree/",
		"/Users/shyang/go/src/github.com/shyang107/paw/filetree/file.go",
	}
	for i := 0; i < len(path); i++ {
		p := path[i]
		fmt.Println(i+1, p)
		f := filetree.NewFile(p)
		cstr := f.LSColorString(f.BaseName)
		s := filepath.Join(f.Dir, cstr)
		fmt.Println("  ", s)
		fmt.Println("  ", f.LSColorString(f.Path))
		fmt.Println("  Path:", f.Path)
		fmt.Println("  Dir:", f.Dir)
		fmt.Println("  BaseName:", f.BaseName)
		fmt.Println("  File:", f.File)
		fmt.Println("  Ext:", f.Ext)
		fmt.Println("  IsDir:", f.IsDir())
		fmt.Println("  IsRegular:", f.IsRegular())
		pathslice := f.PathSlice()
		fmt.Printf("  PathSlice: %d %v\n", len(pathslice), pathslice)
	}
}

func newFile() {
	path := []string{
		"/Users/shyang/go/src/github.com/shyang107/paw/filetree",
		"/Users/shyang/go/src/github.com/shyang107/paw/filetree/",
		"/Users/shyang/go/src/github.com/shyang107/paw/filetree/file.go",
	}
	for i := 0; i < len(path); i++ {
		p := path[i]
		fmt.Println(p)
		f := filetree.NewFile(p)
		spew.Dump(*f)
	}
}
