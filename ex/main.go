package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/shyang107/paw/cast"
	"github.com/shyang107/paw/treeprint"

	"github.com/davecgh/go-spew/spew"

	"github.com/sirupsen/logrus"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/funk"
	// "github.com/thoas/go-funk"
)

var (
	// lg = paw.Glog
	lg  = paw.Logger
	log = paw.Logger
)

func init() {
	lg.SetLevel(logrus.DebugLevel)
}

func main() {
	// exLineCount()
	// exFileLineCount()
	// rehttp()
	// exGetAbbrString()
	// exTableFormat()
	// exStringBuilder()
	// exLoger()
	// exReverse()
	// exPrintTree1()
	// exPrintTree2()
	// exShuffle()
	// exGetCurrPath()
	// var n1 = []int{1, 39, 2, 9, 7, 54, 11}
	// var n2 = []int{1, 39, 2, 9, 7, 54, 11}
	// var n3 = []int{1, 39, 2, 9, 7, 54, 11}
	// var n4 = []int{1, 39, 2, 9, 7, 54, 11}
	// // var n1 = []int{4, 3, 2, 10, 12, 1, 5, 6}
	// // var n2 = []int{4, 3, 2, 10, 12, 1, 5, 6}
	// // size := 20
	// // n1 = paw.GenerateSlice(size)
	// InsertionSort(n1)
	// // n2 = paw.GenerateSlice(size)
	// SelectionSort(n2)
	// // n3 = paw.GenerateSlice(size)
	// exCombSort(n3)
	// // n4 = paw.GenerateSlice(size)
	// exMergeSort(n4)
	// exRegEx()
	// exLogger()
	// exFolder()
	// exGetFiles1()
	// exGetFiles2()
	// exGetFiles3()
	// exGetFilesString()
	// exGrouppingFiles1()
	// exGrouppingFiles2()
	// exGrouppingFiles3()
	// exGrouppingFiles4()
	// exTextTemplate()
	// exRegEx2()
	// sourceFolder := os.Args[1]
	sourceFolder := "../"
	// sourceFolder, _ := homedir.Expand("~/Downloads/0")
	// sourceFolder := "/Users/shyang/go/src/rover/opcc/"
	exWalk(sourceFolder)
	// exFilesMap(sourceFolder)
}

func exWalk(root string) {
	root, _ = filepath.Abs(root)
	root = strings.TrimSuffix(root, "/")
	i := 0
	nf := 0
	nd := 0
	fmt.Printf("%q\n", root)
	folder := make(map[string][]string)
	dirs := []string{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// fmt.Println(path)
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}

		apath, _ := filepath.Abs(path)

		i++
		base := filepath.Base(apath)
		sub := strings.TrimPrefix(apath, root)

		// if paw.REUsuallyExclude.MatchString(path) || strings.HasPrefix(base, ".") {
		// 	return nil
		// }

		pl := strings.Split(path, "/")
		for _, p := range pl {
			if strings.HasPrefix(p, ".") { // ignore hidden files
				return nil
			}
		}

		// fmt.Printf("%q %q\n", sub, base)

		if info.IsDir() {
			nd++
			if _, ok := folder[sub]; !ok {
				folder[sub] = []string{}
			}
		} else {
			sub = strings.TrimSuffix(sub, base)
			sub = strings.TrimSuffix(sub, "/")
			nf++
			folder[sub] = append(folder[sub], base)
		}
		// if info.IsDir() {
		// 	fmt.Printf("%4d %v %q\n", i, info.IsDir(), path)
		// 	fmt.Printf("     %q %q\n", sub, base)
		// 	// fmt.Printf("     %q\n", base)
		// }
		return nil
	})
	dirs = funk.Keys(folder).([]string)
	fmt.Println(nd-1, "directories,", nf, "files.")
	sort.Strings(dirs)
	// spew.Dump(folder)
	spew.Dump(dirs)
	// textoutput(root, dirs, folder)
	treeoutput(root, dirs, folder)

}

func treeoutput(root string, dirs []string, folder map[string][]string) {
	nd, nf := 0, 0
	w := os.Stdout

	tree := treeprint.New()
	for _, dir := range dirs {
		nd++
		ss := strings.Split(strings.TrimPrefix(dir, "/"), "/")
		ns := len(ss)
		level := ns
		// fmt.Printf("ss[%d]: %v\n", ns, ss)
		treend := make([]treeprint.Tree, ns)
		switch {
		case len(dir) == 0: // root
			level = 0
			nd--
			// tree.SetMetaValue(fmt.Sprintf("%d (%d directories, %d files)",
			// 	len(folder[dir]), len(dirs)-1, nFiles(folder)))
			tree.SetMetaValue(fmt.Sprintf("%d", len(folder[dir])))
			tree.SetValue(fmt.Sprintf("%s » root: %q", "./", root))
			treend[0] = tree
		default: // subfolder
			treend[0] = tree.FindByValue(ss[0])
			if treend[0] == nil {
				treend[0] = tree.AddMetaBranch(cast.ToString(len(folder[dir])), ss[0])
			}
			for i := 1; i < ns; i++ {
				treend[i] = treend[i-1].FindByValue(ss[i])
				if treend[i] == nil {
					treend[i] = treend[i-1].AddMetaBranch(cast.ToString(len(folder[dir])), ss[i])
				}
			}
		}
		if len(folder[dir]) == 0 {
			continue
		}
		nf += len(folder[dir])
		level++
		// if treend[ns-1] == nil {
		// 	fmt.Println("root:", root, "dir:", dir)
		// 	os.Exit(1)
		// }
		for _, f := range folder[dir] {
			treend[ns-1].AddNode(f)
		}
	}
	fprintWithLevel(w, 0, tree.String())
	fprintWithLevel(w, 0, fmt.Sprintf("%d directories, %d files.", nd, nf))
}

func nFiles(folder map[string][]string) int {
	n := 0
	for _, v := range folder {
		n += len(v)
	}
	return n
}

func textoutput(root string, dirs []string, folder map[string][]string) {
	top := strings.Repeat("=", 80)
	mid := strings.Repeat("-", 80)
	buttom := top
	nd, nf := 0, 0
	w := os.Stdout
	fprintWithLevel(w, 0, top)

	for i, dir := range dirs {
		// var (
		// 	k = 2
		// )
		level := len(strings.Split(dir, "/")) - 1
		nd++
		switch {
		case len(dir) == 0:
			level = 0
			nd--
			fprintWithLevel(w, level, fmt.Sprintf("%2d %s", i+1, root))
		case len(folder[dir]) == 0:
			fprintWithLevel(w, level, fmt.Sprintf("%2d %s", i+1, dir))
			goto MID
			continue
		default:
			fprintWithLevel(w, level, fmt.Sprintf("%2d %s", i+1, dir))
		}
		nf += len(folder[dir])
		// k := 2
		level++
		for j, f := range folder[dir] {
			// if j > k {
			// 	break
			// }
			fprintWithLevel(w, level, fmt.Sprintf("%2d %s", j+1, f))
		}
	MID:
		if i < len(dirs)-1 {
			// level--
			// fprintWithLevel(w, level, mid[:len(mid)-2*level])
			fprintWithLevel(w, 0, mid)
		}
	}
	fprintWithLevel(w, 0, buttom)
	fprintWithLevel(w, 0, fmt.Sprintf("%d directories, %d files.", nd, nf))
}

func fprintWithLevel(w io.Writer, level int, row string) {
	ns := 2
	space := " "
	pad := strings.Repeat(space, ns*level)
	fmt.Fprintln(w, pad, row)
}
func printWithLevel(level int, row string) {
	ns := 2
	space := " "
	pad := strings.Repeat(space, ns*level)
	fmt.Println(pad, row)
}
