package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/shyang107/paw/cast"
	"github.com/shyang107/paw/treeprint"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/filetree"
)

func exFileTree(root string) {
	// readdir(root)
	// walk(root)
	// constructFile(root)
	readDirs(root)

}
func readDirs(root string) {
	root, err := filepath.Abs(root)
	if err != nil {
		paw.Logger.Fatal(err)
	}

	fl := filetree.NewFileList(root)
	ignore := func(f *filetree.File, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() && strings.HasPrefix(f.BaseName, ".") {
			return filetree.SkipDir
		}
		if strings.HasPrefix(f.BaseName, ".") {
			return filetree.SkipFile
		}
		if strings.Contains(f.Dir, ".git") {
			return filetree.SkipFile
		}
		return nil
	}

	fl.FindFiles(-1, ignore)
	// spew.Dump(fl.Dirs())
	// fmt.Println(fl.ToTreeString("# "))
	fmt.Println(fl.ToTableString("# "))
}

var pad string
var nd, nf int

func fileTree(path, root string, tree treeprint.Tree) (treeprint.Tree, error) {

	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return tree, filepath.SkipDir
	}
	// nd, nf = 0, 0
	for _, fi := range fis {
		path := filepath.Join(path, fi.Name())
		f := filetree.ConstructFileRelTo(path, root)
		isIgnore, _ := filepath.Match(".*", f.BaseName)
		if isIgnore {
			continue
		}

		cdir := f.LSColorString(f.Dir)
		cfile := f.LSColorString(f.BaseName)
		npad := len(f.DirSlice())
		pad := strings.Repeat("   ", npad-1)
		if fi.IsDir() {
			nd++
			fmt.Println(pad, "D", nd, root, cdir, cfile)
			// one := tree.FindByValue(f)
			// if one == nil {
			// 	one = tree.AddBranch(f)
			// }
			tree, _ = fileTree(path, root, tree)
			continue
		}
		nf++
		fmt.Println(pad, "F", nf, root, cdir, cfile)
		tree.AddNode(f)
	}

	return tree, nil
}

func walk(root string) {
	root, err := filepath.Abs(root)
	if err != nil {
		paw.Logger.Fatal(err)
	}

	tree := treeprint.New()
	tree.SetValue(".")
	// f := filetree.ConstructFile(root)
	// one := tree
	// pre := tree
	nd, nf := 0, 0
	visit := func(path string, info os.FileInfo, err error) error {
		f := filetree.ConstructFileRelTo(path, root)
		ss := f.DirSlice()
		nss := len(ss)
		pad := strings.Repeat("   ", nss-1)
		isIgnore, _ := filepath.Match(root+"/.git*", path)
		// fmt.Println(pad, f.LSColorString(cast.ToString(nd)+". "+f.Path))
		if f.IsDir() {
			if isIgnore {
				return filepath.SkipDir
			}
			nd++
			// one = pre.AddBranch(f)
			fmt.Printf("%s %d %v\n", pad+f.LSColorString(cast.ToString(nd)+". "+f.Dir), nss, ss)
			return nil
		}
		isIgnore, _ = filepath.Match(".*", f.BaseName)
		if isIgnore {
			return nil
		}
		nf++
		fmt.Printf("%s %d %v\n", pad+"   "+f.LSColorString(cast.ToString(nf)+". "+f.Dir+"  "+f.BaseName), nss, ss)
		// one.AddNode(f)
		return nil
	}

	err = filepath.Walk(root, visit)
	if err != nil {
		paw.Logger.Fatal(err)
	}
	tree.SetMetaValue(fmt.Sprintf("%s (%d dirs, %d files)", root, nd-1, nf))

	fmt.Println(string(tree.Bytes()))

}

func readdir(root string) {
	root, err := filepath.Abs(root)
	if err != nil {
		paw.Logger.Fatal(err)
	}

	fis, err := ioutil.ReadDir(root)
	if err != nil {
		paw.Logger.Fatal(err)
	}

	tree := treeprint.New()
	idir := 0
	ifile := 0
	for _, fi := range fis {
		path := filepath.Join(root, fi.Name())
		f := filetree.ConstructFile(path)
		fmt.Println(f.LSColorString(f.Path))
		if fi.IsDir() {
			idir++
			tree.AddMetaBranch(idir, f)
		} else {
			ifile++
			tree.AddMetaNode(ifile, f)
		}
	}
	tree.SetMetaValue(fmt.Sprintf("%q (%d dirs., %d files)", root, idir, ifile))
	fmt.Println(string(tree.Bytes()))

}

func constructFile(root string) {
	root, err := filepath.Abs(root)
	if err != nil {
		paw.Logger.Fatal(err)
	}
	p1 := filetree.ConstructFile(root)
	p2 := filetree.ConstructFile(root)
	fmt.Println(reflect.DeepEqual(p1, p2))

	dir, basename := filepath.Split(root)
	fmt.Printf("dir: %q %q\n", dir, filepath.Dir(root))
	fmt.Printf("basename: %q %q\n", basename, filepath.Base(root))
	stat, _ := os.Stat(root)
	p3 := &filetree.File{
		Path:     root,
		Dir:      filepath.Dir(root),
		BaseName: filepath.Base(root),
		File:     strings.TrimSuffix(filepath.Base(root), filepath.Ext(root)),
		Ext:      filepath.Ext(root),
		Stat:     stat,
	}
	fmt.Println(reflect.DeepEqual(p1, p3))
	spew.Dump(dir)
	spew.Dump(filepath.Dir(root))
	spew.Dump(*p1)
	path := filepath.Join(root, "_ex")
	p4 := filetree.ConstructFileRelTo(path, root)
	spew.Dump(*p4)
}
