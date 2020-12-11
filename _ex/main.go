package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/go-homedir"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/_junk"
	"github.com/shyang107/paw/cast"
	"github.com/shyang107/paw/filetree"
	"github.com/shyang107/paw/funk"
	"github.com/shyang107/paw/godirwalk"
	"github.com/shyang107/paw/treeprint"
	"github.com/sirupsen/logrus"
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
	// root := os.Args[1]
	// root := "../"
	// root, _ := homedir.Expand("~")
	root, _ := homedir.Expand("~/Downloads")
	// root, _ := homedir.Expand("~/Downloads/0")
	// root := "/Users/shyang/go/src/rover/opcc/"
	// exWalk(root)
	// exFilesMap(root)
	exPathMap(root)
	// exColor()
}

func exPathMap1(root string) {
	root, _ = filepath.Abs(root)
	root = strings.TrimSuffix(root, "/")

	pm := _junk.NewPathMap()
	isRecursive := true

	ignoreHidden := true
	ignoreCondition := false
	targetType := []string{"", "", ""}
	ignoreFile := []string{"index.js"}
	ignorePath := []string{".git"}
	ignoreType := []string{".gitignore", ".exe", ".go"}
	pm.SetCondition(ignoreHidden, ignoreCondition, targetType, ignoreFile, ignorePath, ignoreType)

	pm.FindFiles(root, isRecursive)
	// spew.Dump(pm.GetDirs())
	w := os.Stdout
	// pm.Fprint(w, paw.OPlainTextMode, "", "# ")
	// pm.Fprint(w, paw.OTableFormatMode, "", "")
	pm.Fprint(w, _junk.OTreeMode, "", "# ")
	// spew.Dump(pm)
	// fmt.Println(pm.GetFilesString())
	// spew.Dump(pm.GetCondition())
	// for i, f := range pm.GetPaths() {
	// 	fmt.Println(i+1, f)
	// }
}

func exPathMap2(root string) {
	root, _ = filepath.Abs(root)
	// root = strings.TrimSuffix(root, "/")
	i := 1
	godirwalk.Walk(root, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if b, err := de.IsDirOrSymlinkToDir(); b == true && err == nil {
				if strings.HasPrefix(de.Name(), ".") {
					// return filepath.SkipDir
					return godirwalk.SkipThis
				}
			} else {
				if strings.HasPrefix(de.Name(), ".") {
					return nil
				}
				// fmt.Printf("%d. %s %s %s\n", i, de.ModeType(), osPathname, ext)
				str, _ := filetree.FileLSColorStr(osPathname, de.Name())
				// fmt.Printf("%d. %v %s\n", i, de.ModeType(), str)
				dir := filepath.Dir(osPathname)
				path := filepath.Join(dir, str)
				fi, _ := os.Lstat(osPathname)
				fmt.Printf("%d. %v %s\n", i, fi.Mode(), path)
				i++
			}
			return nil
		},
		// Unsorted: false,
	})
}
func exPathMap(root string) {
	root, _ = filepath.Abs(root)
	// root = "/Users/shyang"
	fns, _ := godirwalk.ReadDirnames(root, nil)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "cannot get list of directory children")
	// }
	sort.Strings(fns)

	for _, f := range fns {
		path := filepath.Join(root, f)
		base, err := filetree.FileLSColorStr(path, f)
		if err != nil {
			paw.Logger.Errorln(err)
		}
		fi, _ := os.Lstat(path)
		fmt.Printf("%s %s\n", fi.Mode(), base)
	}

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
	// outputText(root, dirs, folder)
	// outputTree(root, dirs, folder)
	outputTable(root, dirs, folder)
}

func outputTree(root string, dirs []string, folder map[string][]string) {
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

func outputTable(root string, dirs []string, folder map[string][]string) {
	nd, nf := 0, 0

	w := os.Stdout

	tf := paw.NewTableFormat()
	tf.Fields = []string{"No.", "Files"}
	tf.LenFields = []int{4, 76}
	tf.Aligns = []paw.Align{paw.AlignRight, paw.AlignLeft}
	tf.Prepare(w)

	// tf.SetBeforeMessage(msg)
	tf.PrintSart()
	for i, dir := range dirs {
		level := len(strings.Split(dir, "/")) - 1
		dfiles := len(folder[dir])
		nd++
		subhead := fmt.Sprintf("Depth: %d, .%s ( %d files)", level, dir, dfiles)
		switch {
		case len(dir) == 0:
			level = 0
			nd--
			tf.PrintRow("", fmt.Sprintf("Depth: %d, %s ( %d files)", level, root, dfiles))
		case len(folder[dir]) == 0:
			tf.PrintRow("", subhead)
			goto MID
		default:
			tf.PrintRow("", subhead)
		}
		nf += len(folder[dir])
		level++
		for j, f := range folder[dir] {
			tf.PrintRow(cast.ToString(j+1), f)

		}
	MID:
		// tf.PrintRow("", fmt.Sprintf("Sum: %d files.", len(folder[dir])))
		if i < len(dirs)-1 {
			tf.PrintMiddleSepLine()
		}
	}
	tf.SetAfterMessage(fmt.Sprintf("\n%d directories, %d files\n", nd, nf))
	tf.PrintEnd()
}

func outputText(root string, dirs []string, folder map[string][]string) {
	top := strings.Repeat("=", 80)
	mid := strings.Repeat("-", 80)
	buttom := top
	nd, nf := 0, 0
	w := os.Stdout
	fprintWithLevel(w, 0, top)

	for i, dir := range dirs {
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
		level++
		for j, f := range folder[dir] {
			fprintWithLevel(w, level, fmt.Sprintf("%2d %s", j+1, f))
		}
	MID:
		if i < len(dirs)-1 {
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
