package paw

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/shyang107/paw/cast"
	"github.com/shyang107/paw/treeprint"
)

// PathMap store paths of files
type PathMap struct {
	root   string
	folder map[string][]string
	dirs   []string
}

// NewPathMap will return an instance of `PathMap`
func NewPathMap() *PathMap {
	p := &PathMap{
		root:   "",
		folder: make(map[string][]string),
		dirs:   []string{},
	}
	return p
}

func (m PathMap) String() string {
	return m.Text("", "")
}

// SetFolder will store `folder`
func (m *PathMap) SetFolder(folder map[string][]string) {
	if folder == nil {
		return
	}
	m.folder = folder
}

// GetRoot will return `root`
func (m *PathMap) GetRoot() string {
	return m.root
}

// SetRoot will store `root`
func (m *PathMap) SetRoot(root string) {
	if len(root) == 0 {
		return
	}
	m.root = root
}

// GetFolder will return `folder`
func (m *PathMap) GetFolder() map[string][]string {
	return m.folder
}

// SetDirs will store `dirs`
func (m *PathMap) SetDirs(dirs []string) {
	if dirs == nil {
		return
	}
	m.dirs = dirs
}

// GetDirs will return `dirs`
func (m *PathMap) GetDirs() []string {
	return m.dirs
}

// GetInfos will return `[]os.FileInfo` of paths
func (m *PathMap) GetInfos() ([]os.FileInfo, error) {
	fis := []os.FileInfo{}
	path := ""
	for k, files := range m.folder {
		if len(k) == 0 {
			path = m.root
		} else {
			path = m.root + k + "/"
		}
		for _, f := range files {
			path += f
			fi, err := os.Stat(path)
			if err != nil {
				return nil, err
			}
			fis = append(fis, fi)
		}
	}
	return fis, nil
}

// NFiles will return the numbers of files
func (m *PathMap) NFiles() int {
	folder := m.GetFolder()
	n := 0
	for _, v := range folder {
		n += len(v)
	}
	return n
}

// NDirs will return the numbers of sub-directories
func (m *PathMap) NDirs() int {
	return len(m.GetDirs()) - 1
}

// Fprint filelist with `head`
func (m *PathMap) Fprint(w io.Writer, mode OutputMode, head, pad string) {
	switch mode {
	case OTreeMode:
		m.FprintTree(w, head, pad)
	case OTableFormatMode:
		tf := &TableFormat{
			Fields:    []string{"No.", "Files"},
			LenFields: []int{5, 75},
			Aligns:    []Align{AlignRight, AlignLeft},
			Padding:   pad,
		}
		m.FprintTable(w, tf, head)
	default: // OPlainTextMode
		m.FprintText(w, head, pad)
	}
}

// Tree will return a string in tree mode (use `FprintTree`)
func (m *PathMap) Tree(head, pad string) string {
	buf := new(bytes.Buffer)
	m.FprintTree(buf, head, pad)
	return TrimFrontEndSpaceLine(string(buf.Bytes()))
}

// FprintTree print out in tree mode
func (m *PathMap) FprintTree(w io.Writer, head, pad string) {
	fmt.Fprintln(w, PaddingString(head, pad))
	fmt.Fprintln(w, pad)
	foutputTree(w, m.GetRoot(), m.GetDirs(), m.GetFolder())
}

func foutputTree(w io.Writer, root string, dirs []string, folder map[string][]string) {

	nd, nf := 0, 0

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

// Table will return a string in table mode with `head` (use `FprintTable`)
func (m *PathMap) Table(head, pad string) string {
	buf := new(bytes.Buffer)
	m.Fprint(buf, OTableFormatMode, head, pad)
	return TrimFrontEndSpaceLine(string(buf.Bytes()))
}

// FprintTable print out in table mode with `head`
func (m *PathMap) FprintTable(w io.Writer, tf *TableFormat, head string) {
	tf.Prepare(w)
	tf.SetBeforeMessage(head)
	foutputTable(tf, m.GetRoot(), m.GetDirs(), m.GetFolder())
}

func foutputTable(tf *TableFormat, root string, dirs []string, folder map[string][]string) {

	nd, nf := 0, 0

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

// Text return string in plain text mode (use `FprintText`)
func (m *PathMap) Text(head, pad string) string {
	buf := new(bytes.Buffer)
	m.Fprint(buf, OPlainTextMode, head, pad)
	return TrimFrontEndSpaceLine(string(buf.Bytes()))
}

// FprintText print out in plain text mode
func (m *PathMap) FprintText(w io.Writer, head, pad string) {
	fmt.Fprintln(w, PaddingString(head, pad))
	fmt.Fprintln(w, pad)
	buf := new(bytes.Buffer)
	foutputText(buf, m.GetRoot(), m.GetDirs(), m.GetFolder())
	fmt.Fprintln(w, PaddingString(string(buf.Bytes()), pad))
	fmt.Fprintln(w, pad)
}

func foutputText(w io.Writer, root string, dirs []string, folder map[string][]string) {

	top := strings.Repeat("=", 80)
	mid := strings.Repeat("-", 80)
	buttom := top
	nd, nf := 0, 0

	fprintWithLevel(w, 0, top)
	for i, dir := range dirs {
		level := len(strings.Split(dir, "/")) - 1
		nd++
		subhead := fmt.Sprintf("%2d %s", i+1, dir)
		switch {
		case len(dir) == 0:
			level = 0
			nd--
			fprintWithLevel(w, level, fmt.Sprintf("%2d %s", i+1, root))
		case len(folder[dir]) == 0:
			fprintWithLevel(w, level, subhead)
			goto MID
			continue
		default:
			fprintWithLevel(w, level, subhead)
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

func stringWithLevel(level int, row string) string {
	ns := 3
	space := " "
	pad := strings.Repeat(space, ns*level)
	return fmt.Sprintln(pad, row)
}

func fprintWithLevel(w io.Writer, level int, row string) {
	ns := 3
	space := " "
	pad := strings.Repeat(space, ns*level)
	fmt.Fprintln(w, pad, row)
}

// FindFiles ...
func (m *PathMap) FindFiles(root string, isRecursive bool) error {

	root, _ = filepath.Abs(root)
	// root = strings.TrimSuffix(root, "/")
	dirs := []string{}
	folder := make(map[string][]string)
	var err error
	if isRecursive {
		err = walkDir(root, &dirs, &folder)
	} else {
		err = ioReadDir(root, &dirs, &folder)
	}
	if err != nil {
		return err
	}

	m.SetRoot(root)
	m.SetDirs(dirs)
	m.SetFolder(folder)

	return nil
}

func osReadDir(root string, dirs *[]string, folder *map[string][]string) error {

	f, err := os.Open(root)
	if err != nil {
		return err
	}

	fis, err := f.Readdir(-1)
	defer f.Close()
	if err != nil {
		return err
	}

	for _, fi := range fis {
		// TODO ignoring conditions
		if strings.HasPrefix(fi.Name(), ".") { // ignore hidden files
			continue
		}
		(*folder)[""] = append((*folder)[""], fi.Name())
	}

	*dirs = append(*dirs, "")

	return nil
}

func ioReadDir(root string, dirs *[]string, folder *map[string][]string) error {
	// root, _ = filepath.Abs(root)
	// root = strings.TrimSuffix(root, "/")
	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		// TODO ignoring conditions
		if strings.HasPrefix(fi.Name(), ".") { // ignore hidden files
			continue
		}
		(*folder)[""] = append((*folder)[""], fi.Name())
	}

	*dirs = append(*dirs, "")

	return nil
}

func walkDir(root string, dirs *[]string, folder *map[string][]string) error {

	// root, _ = filepath.Abs(root)
	// root = strings.TrimSuffix(root, "/")

	visitFile := func(path string, info os.FileInfo, err error) error {
		// fmt.Println(path)
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}

		apath, _ := filepath.Abs(path)

		base := filepath.Base(apath)
		sub := strings.TrimPrefix(apath, root)

		// if paw.REUsuallyExclude.MatchString(path) || strings.HasPrefix(base, ".") {
		// 	return nil
		// }
		// TODO ignoring coditions
		pl := strings.Split(path, "/")
		for _, p := range pl {
			if strings.HasPrefix(p, ".") { // ignore hidden files
				return nil
			}
		}

		// fmt.Printf("%q %q\n", sub, base)

		if info.IsDir() {
			if _, ok := (*folder)[sub]; !ok {
				(*folder)[sub] = []string{}
				(*dirs) = append(*dirs, sub)
			}
		} else {
			sub = strings.TrimSuffix(sub, base)
			sub = strings.TrimSuffix(sub, "/")
			(*folder)[sub] = append((*folder)[sub], base)
		}
		return nil
	}

	err := filepath.Walk(root, visitFile)
	if err != nil {
		return err
	}
	return nil
}
