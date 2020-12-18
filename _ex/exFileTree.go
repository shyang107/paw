package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/filetree"
	"github.com/shyang107/paw/godirwalk"
)

func exFileTree(root string) {
	// readdir(root)
	// walk(root)
	// constructFile(root)
	readDirs(root)
	// scan(root)
}

func readDirs(root string) {
	root, err := filepath.Abs(root)
	if err != nil {
		paw.Logger.Fatal(err)
	}

	fl := filetree.NewFileList(root)
	// ignore := func(f *filetree.File, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if f.IsDir() && strings.HasPrefix(f.BaseName, ".") {
	// 		return filetree.SkipDir
	// 	}
	// 	if strings.HasPrefix(f.BaseName, ".") {
	// 		return filetree.SkipFile
	// 	}
	// 	// if strings.Contains(f.Dir, ".git") {
	// 	// 	return filetree.SkipFile
	// 	// }
	// 	return nil
	// }
	// fl.FindFiles(-1, ignore)
	fl.FindFiles(-1, nil)
	// spew.Dump(fl.Dirs())
	// fmt.Println(fl.ToTreeString("# "))
	// fmt.Println(fl.ToTableString("# "))
	// fmt.Println(fl.ToTextString("# "))
	fmt.Println(fl.ToListString(""))
	// fmt.Println(fl)
	// spew.Dump(fl.Map())
}

func scan(pathname string) {

	scanner, err := godirwalk.NewScanner(pathname)
	if err != nil {
		fatal("cannot scan directory: %s", err)
	}

	for scanner.Scan() {
		dirent, err := scanner.Dirent()
		if err != nil {
			warning("cannot get dirent: %s", err)
			continue
		}
		name := dirent.Name()
		if name == "break" {
			break
		}
		if name == "continue" {
			continue
		}
		stat, _ := os.Stat(filepath.Join(pathname, name))
		fmt.Printf("%v %v %v\n", dirent.ModeType(), stat.Mode(), name)

	}
	if err := scanner.Err(); err != nil {
		fatal("cannot scan directory: %s", err)
	}
}

func stderr(f string, args ...interface{}) {
	paw.Logger.Error(fmt.Sprintf(f, args...))
}

func fatal(f string, args ...interface{}) {
	stderr(f, args...)
	os.Exit(1)
}

func warning(f string, args ...interface{}) {
	stderr(f, args...)
}
