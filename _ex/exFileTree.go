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
	// root, _ = homedir.Expand("~")
	root, err := filepath.Abs(root)
	if err != nil {
		paw.Logger.Fatal(err)
	}

	fl := filetree.NewFileList(root)
	fl.FindFiles(-1, nil)
	// spew.Dump(fl.Dirs())
	// fmt.Println(fl.ToTreeString("# "))
	// fmt.Println(fl.ToTableString("# "))
	// fmt.Println(fl.ToTextString("# "))
	fmt.Println(fl.ToListString("# "))
	// fmt.Println(fl)

}
func lscolors() {
	// spew.Dump(fl.Map())
	// fmt.Println(filetree.KindLSColorString(".sh", "sh"))
	// fmt.Println(filetree.KindLSColorString(".go", "go"))
	// fmt.Println(filetree.KindLSColorString("di", "di"))
	// fmt.Println(filetree.KindLSColorString("fi", "fi"))
	// fmt.Println(filetree.KindLSColorString("ln", "ln"))
	// fmt.Println(filetree.KindLSColorString("pi", "pi"))
	// fmt.Println(filetree.KindLSColorString("so", "so"))
	// fmt.Println(filetree.KindLSColorString("bd", "bd"))
	// fmt.Println(filetree.KindLSColorString("cd", "cd"))
	// fmt.Println(filetree.KindLSColorString("or", "or"))
	// fmt.Println(filetree.KindLSColorString("mi", "mi"))
	// fmt.Println(filetree.KindLSColorString("ex", "ex"))
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
