package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/karrick/godirwalk"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/filetree"
	// "github.com/shyang107/paw/3rd-party/filetree"
	// "github.com/shyang107/paw/filetree"
	// "github.com/xlab/treeprint"
)

func exFileTree(root string) {
	// readdir(root)
	// walk(root)
	// constructFile(root)
	readDirs(root)
	// scan(root)
	// xafero(root)
}

func readDirs(root string) {
	// root, _ = homedir.Expand("~")
	root, err := filepath.Abs(root)
	if err != nil {
		paw.Logger.Error(err)
	}

	fl := filetree.NewFileList(root)
	fl.SetIgnoreFunc(nil)
	fl.FindFiles(-1)

	// spew.Dump(fl.Dirs())
	// fmt.Println(fl.ToTreeViewString("# "))
	// fmt.Println(fl.ToTableViewString("# "))
	// fmt.Println(fl.ToLevelViewString("# "))
	// fmt.Println(fl.ToListViewString("# "))
	fmt.Println(fl.ToListTreeView("# "))
	// fl.SetWriters(os.Stdout)
	// fl.ToListTreeViewString("# ")
	// fmt.Println("out:\n", out)
	// fmt.Println(fl)
	// listfl(fl)
}

var appFs = afero.NewMemMapFs()

func xafero(root string) {
	root, err := filepath.Abs(root)
	if err != nil {
		paw.Logger.Fatal(err)
	}
	// re := regexp.MustCompile(`^[^.].+$`)
	// re := regexp.MustCompile(`^[.].+$`)

	// fmt.Println(re.String(), `.git`, !re.MatchString(`.git`))
	a := afero.Afero{
		Fs: afero.NewOsFs(),
		// Fs: afero.NewRegexpFs(afero.NewOsFs(), re),
	}

	fis, err := a.ReadDir(root)
	if err != nil {
		paw.Logger.Fatal(err)
	}
	git, _ := filetree.GetShortGitStatus(root)
	for i, fi := range fis {
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		file, _ := filetree.NewFile(root + "/" + fi.Name())
		meta, _ := file.MetaC(git)
		fmt.Printf("%2d %s %v\n", i+1, meta, file)
	}
}

func listfl(fl *filetree.FileList) {
	dirs := fl.Dirs()
	fm := fl.Map()
	for i, dir := range dirs {
		fmt.Println(i, dir)
		for j, file := range fm[dir] {
			fmt.Println("    ", j+1, file)
		}
	}
}

func LSColors() {
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
