package _junk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/cast"
	"github.com/shyang107/paw/funk"
	"github.com/shyang107/paw/treeprint"
)

var (
	// TypeDesc ...
	TypeDesc = map[string]string{
		"di": "directory",
		"fi": "file",
		"ln": "symbolic link",
		"pi": "fifo file",
		"so": "socket file",
		"bd": "block (buffered) special file",
		"cd": "character (unbuffered) special file",
		"or": "symbolic link pointing to a non-existent file (orphan)",
		"mi": "non-existent file pointed to by a symbolic link (visible when you type ls -l)",
		"ex": "file which is executable (ie. has 'x' set in permissions)",
	}
	// Colors = make(map[string]string)
	Colors = make(map[string][]color.Attribute)
	// exts    = []string{}

	// NoColor check from the type of terminal and
	// determine output to terminal in color (`true`) or not (`false`)
	NoColor = os.Getenv("TERM") == "dumb" || !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))
)

func init() {
	getcolors()
}

// SetNoColor will set `true` to `NoColor`
func SetNoColor() {
	NoColor = true
}

// ResumNoColor will resume the default value of `NoColor`
func ResumNoColor() {
	NoColor = os.Getenv("TERM") == "dumb" || !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))
}

func getcolors() {
	colorenv := os.Getenv("LS_COLORS")
	args := strings.Split(colorenv, ":")

	// colors := make(map[string]string)
	// ctypes := make(map[string]string)
	// exts := []string{}
	for _, a := range args {
		// fmt.Printf("%v\t", a)
		kv := strings.Split(a, "=")

		// fmt.Printf("%v\n", kv)
		if len(kv) == 2 {
			Colors[kv[0]] = getColorAttribute(kv[1])
			// exts = append(exts, kv[0])
		}
	}
	// sort.Strings(exts)
}

func getColorAttribute(code string) []color.Attribute {
	att := []color.Attribute{}
	for _, a := range strings.Split(code, ";") {
		att = append(att, color.Attribute(cast.ToInt(a)))
	}
	return att
}

func colorstr(att []color.Attribute, s string) string {
	cs := color.New(att...)
	return cs.Sprint(s)
}

// FileColorStr will return the color string of `s` form `fullpath`
func FileColorStr(fullpath, s string) (string, error) {
	ext, err := GetColorExt(fullpath)
	if err != nil {
		return "", err
	}
	switch {
	case NoColor:
		return s, nil
	default:
		if _, ok := Colors[ext]; !ok {
			return s, nil
		}
		return colorstr(Colors[ext], s), nil
	}
}

// GetColorExt will return the color key of extention from `fullpath`
func GetColorExt(fullpath string) (ext string, err error) {
	fi, err := os.Lstat(fullpath)
	if err != nil {
		return "", errors.New("GetColorExt:" + err.Error())
	}

	switch mode := fi.Mode(); {
	case mode.IsDir(): // d: is a directory 資料夾模式
		ext = "di" // di = directory

	case mode&os.ModeSymlink != 0: // L: symbolic link 象徵性的關聯
		ext = "ln" // ln = symbolic link
		// if !HasFile(fullpath) { // or = symbolic link pointing to a non-existent file (orphan)
		// 	ext = "or"
		// } else { // mi = non-existent file pointed to by a symbolic link (visible when you type ls -l)
		// 	ext = "mi"
		// }
	case mode&os.ModeSocket != 0: // S: Unix domain socket Unix 主機 socket
		ext = "so" // so = socket file
	case mode&os.ModeNamedPipe != 0:
		ext = "pi" //pi = fifo file
	// case mode&os.ModeDevice != 0:
	// 	ext = ""

	// bd = block (buffered) special file
	// cd = character (unbuffered) special file

	// ex = file which is executable (ie. has 'x' set in permissions)
	default: // fi = file
		ext = "*" + filepath.Ext(fullpath)
	}
	return ext, nil
}

type pmCondition struct {
	ignoreHidden    bool     // 第一優先， `true` 忽略路徑中以 `.` 開頭的檔案或目錄
	ignoreCondition bool     // false 忽略過濾條件
	targetType      []string // 目標檔案型別
	ignoreFile      []string // 忽略檔案 (檔名，包括副檔名)
	ignorePath      []string // 忽略目錄
	ignoreType      []string // 忽略檔案型別
}

// PathMap store paths of files
type PathMap struct {
	root   string              // 根目錄
	folder map[string][]string // 檔案名稱 (basename, xxxx.xxx) map ，按子目錄儲存，完整路徑為 root +
	dirs   []string            // 子目錄索引，路徑不含根目錄
	cond   pmCondition         // 檔案過濾條件
}

// NewPathMap will return an instance of `PathMap`
func NewPathMap() *PathMap {
	p := &PathMap{
		root:   "",
		folder: make(map[string][]string),
		dirs:   []string{},
		cond: pmCondition{
			ignoreHidden:    true,       // 第一優先， `true` 忽略路徑中以 `.` 開頭的檔案或目錄
			ignoreCondition: true,       // false 忽略過濾條件
			targetType:      []string{}, // 目標檔案型別
			ignoreFile:      []string{}, // 忽略檔案 (檔名，包括副檔名)
			ignorePath:      []string{}, // 忽略目錄
			ignoreType:      []string{}, // 忽略檔案型別
		},
	}
	return p
}

// func (m PathMap) String() string {
// 	return m.Text("", "")
// }

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
	sort.Strings(dirs)
	m.dirs = dirs
}

// GetDirs will return `dirs`
func (m *PathMap) GetDirs() []string {
	return m.dirs
}

// GetFileInfo will return `[]os.FileInfo` of paths
func (m *PathMap) GetFileInfo() ([]os.FileInfo, error) {
	fis := []os.FileInfo{}
	for k, files := range m.folder {
		for _, f := range files {
			path := filepath.Join(m.root, k, f)
			fi, err := os.Stat(path)
			if err != nil {
				return nil, err
			}
			fis = append(fis, fi)
		}
	}
	return fis, nil
}

// GetPaths will return string of all fullpaths
func (m *PathMap) GetPaths() []string {
	fs := []string{}
	for _, dir := range m.dirs {
		for _, name := range m.folder[dir] {
			fullpath := filepath.Join(m.root, dir, name)
			fs = append(fs, fullpath)
		}
	}
	return fs
}

// GetPathsString will return string of all fullpaths
func (m *PathMap) GetPathsString() string {
	buf := new(bytes.Buffer)
	i := 1
	for _, dir := range m.dirs {
		for _, name := range m.folder[dir] {
			fullpath := filepath.Join(m.root, dir, name)
			buf.WriteString(fmt.Sprintf("%d. %s\n", i, fullpath))
			i++
		}
	}
	buf.WriteString(fmt.Sprintf("\n%d directories, %d files.\n", m.NDirs(), m.NFiles()))
	return string(buf.Bytes())
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

// SetCondition store conditions to filter files
func (m *PathMap) SetCondition(ignoreHidden, ignoreCondition bool, targetType, ignoreFile, ignorePath, ignoreType []string) {
	m.cond = pmCondition{
		ignoreHidden:    ignoreHidden,    // 第一優先， `true` 忽略路徑中以 `.` 開頭的檔案或目錄
		ignoreCondition: ignoreCondition, // false 忽略過濾條件
		targetType:      targetType,      // 目標檔案型別
		ignoreFile:      ignoreFile,      // 忽略檔案 (檔名，包括副檔名)
		ignorePath:      ignorePath,      // 忽略目錄
		ignoreType:      ignoreType,      // 忽略檔案型別
	}
}

// GetCondition will return `map[{condition}][]string` of filter-conditions of files
func (m *PathMap) GetCondition() map[string][]string {
	cond := make(map[string][]string)
	cond["targetType"] = m.cond.targetType
	cond["ignoreFile"] = m.cond.ignoreFile
	cond["ignorePath"] = m.cond.ignorePath
	cond["ignoreType"] = m.cond.ignoreType
	return cond
}

// Fprint filelist with `head`
func (m *PathMap) Fprint(w io.Writer, mode OutputMode, head, pad string) {
	switch mode {
	case OTreeMode:
		m.FprintTree(w, head, pad)
	case OTableFormatMode:
		tf := &paw.TableFormat{
			Fields:    []string{"No.", "Files"},
			LenFields: []int{5, 75},
			Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft},
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
	return paw.TrimFrontEndSpaceLine(string(buf.Bytes()))
}

// FprintTree print out in tree mode
func (m *PathMap) FprintTree(w io.Writer, head, pad string) {
	if len(head) > 0 {
		fmt.Fprintln(w, paw.PaddingString(head, pad))
	}
	fmt.Fprintln(w, pad)
	buf := new(bytes.Buffer)
	foutputTreeX2(buf, m.GetRoot(), m.GetDirs(), m.GetFolder())
	fmt.Fprintln(w, paw.PaddingString(string(buf.Bytes()), pad))
}

func foutputTree(w io.Writer, root string, dirs []string, folder map[string][]string) {

	nd, nf := 0, 0

	tree := treeprint.New()
	for _, dir := range dirs {
		nd++
		ss := strings.Split(strings.TrimPrefix(dir, "/"), "/")
		ns := len(ss)
		level := ns
		treend := make([]treeprint.Tree, ns)
		switch {
		case len(dir) == 0: // root
			level = 0
			nd--
			// tree.SetMetaValue(fmt.Sprintf("%d (%d directories, %d files)",
			// 	len(folder[dir]), len(dirs)-1, nFiles(folder)))
			tree.SetMetaValue(fmt.Sprintf("%d", len(folder[dir])))
			tree.SetValue(fmt.Sprintf("%s » root: %s", "./", root))
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
		for _, f := range folder[dir] {
			treend[ns-1].AddNode(f)
		}
	}
	// fprintWithLevel(w, 0, tree.String())
	printTreeColor(w, tree.String())
	fprintWithLevel(w, 0, "")
	fprintWithLevel(w, 0, fmt.Sprintf("%d directories, %d files.", nd, nf))
	// fprintWithLevel(w, 0, "")
}

func printTreeColor(w io.Writer, s string) {
	if NoColor {
		fmt.Fprintln(w, s)
		return
	}
	s = paw.TrimFrontEndSpaceLine(s)
	buf := new(bytes.Buffer)

	lines := strings.Split(s, "\n")
	sep := ":"
	fields := strings.Split(lines[0], sep)
	sep = "./"
	part1 := strings.Split(fields[0], sep)
	buf.WriteString(part1[0])
	cstr := colorstr(Colors["di"], sep)
	buf.WriteString(cstr)
	buf.WriteString(part1[1])

	cstr = colorstr(Colors["di"], fields[1])
	buf.WriteString(":" + cstr + "\n") // print match

	for _, row := range lines[1:] {
		sep := "]  "
		fields := strings.Split(row, sep)
		if len(fields) == 2 {
			buf.WriteString(fields[0] + sep) // print text before match
			cstr := colorstr(Colors["di"], fields[1])
			buf.WriteString(cstr + "\n") // print match
			continue
		}

		sep = "── "
		fields = strings.Split(row, sep)
		if len(fields) == 2 {
			buf.WriteString(fields[0] + sep) // print text before match
			ext := "*" + filepath.Ext(fields[1])
			cstr := colorstr(Colors[ext], fields[1])
			buf.WriteString(cstr + "\n") // print match
			continue
		} else {
			paw.Logger.Errorln("file")
			continue
		}
	}
	buf.WriteTo(w)
}

// EdgeType is edge string for tree view output
type EdgeType string

var (
	// EdgeTypeLink is "│"
	EdgeTypeLink EdgeType = "│"
	// EdgeTypeMid is "├──"
	EdgeTypeMid EdgeType = "├──"
	// EdgeTypeEnd is "└──"
	EdgeTypeEnd EdgeType = "└──"
)

// IndentSize is the number of spaces per tree level.
var IndentSize = 3

func foutputTreeX2(w io.Writer, root string, dirs []string, folder map[string][]string) {
	// fmt.Fprintln(w, root)
	// fmt.Println(root)
	// spew.Dump(dirs)
	// spew.Dump(folder)
	var (
		// rootRelDir = "-" // FileSummary.RelDir
		// level0    = ""  // FileSummary.RelDir
		nd, nf int
		sp     = " "
	)

	levelm := getLevelm(dirs)
	dm := getSubdirsMap(dirs)
	// for _, dir := range dirs {
	// 	if len(dm[dir]) == 0 {
	// 		continue
	// 	}
	// 	nd += len(dm[dir])
	// 	ss := strings.Split(dir, string(os.PathSeparator))
	// 	ns := len(ss) - 1
	// 	fmt.Printf("%s%q %d %v\n", strings.Repeat(sp, 3*ns), dir, len(dm[dir]), dm[dir])
	// }
	// fmt.Printf("%d directories\n", nd)
	// return
	for i, dir := range dirs {
		nd++
		fullpath := filepath.Join(root, dir)
		ds := paw.NewFileSummary(fullpath, root)
		// dirpath := ds.RelDir
		level := levelm[dir]
		meta := getSubdirsMeta(dir, dirs, folder)
		if len(dir) == 0 { // root
			nd--
			croot, _ := FileColorStr(fullpath, "./")
			croot += " "
			cstr, _ := FileColorStr(fullpath, root)
			fmt.Fprintf(w, "%s%d. %s\n", strings.Repeat(sp, IndentSize*level), i, croot+meta+cstr)
		} else { // subfolder
			level--
			// edge := EdgeTypeMid
			edge := getEdge(i, 0, ds, dirs, folder, level, levelm, dm)
			cstr, _ := FileColorStr(fullpath, ds.Name)
			fmt.Fprintf(w, "%s %d. %s\n", edge, i, meta+cstr)
			level++
		}
		continue
		for j, file := range folder[dir] {
			fullpath := filepath.Join(root, dir, file)
			fs := paw.NewFileSummary(fullpath, root)
			edge := EdgeTypeMid
			// edge := getEdge(i, j, dir, dirs, folder, level, levelm)
			cstr, _ := FileColorStr(fullpath, fs.Name)
			fmt.Fprintf(w, "%s%s %d. %s\n", strings.Repeat(sp, IndentSize*level), edge, nf+j+1, cstr)
		}
		nf += len(folder[dir])
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%d directories, %d files.\n", nd, nf)
}

func getEdge(id, jf int, ds *paw.FileSummary, dirs []string, folder map[string][]string, level int, levelm map[string]int, dm map[string][]string) string {

	buf := []byte{} // new(bytes.Buffer)
	sp := " "
	var edge EdgeType
	nd := len(dirs) - 1
	// nf := len(folder[dir])
	edge = EdgeTypeMid
	if jf == 0 { // dir
		if level == 0 {
			if id == nd { // last dir
				edge = EdgeTypeEnd
			}
			buf = append(buf, edge...)
			goto END
		}
		// level > 0
		// buf = append(buf, strings.Repeat(sp, IndentSize)...)

		dir := filepath.Join(ds.RelDir, ds.Name)
		ipre := funk.IndexOfString(dm[ds.RelDir], ds.Name)
		fmt.Printf("%d %s %s %s\n", id, ds.RelDir, ds.Name, dir)
		fmt.Printf("  %d pre: %v\n", ipre, dm[ds.RelDir])

		if ipre < 0 || ipre == len(dm[ds.RelDir])-1 {
			edge = EdgeTypeEnd
		}
		buf = append(buf, EdgeTypeLink...)

		for i := 0; i < level-1; i++ {
			buf = append(buf, strings.Repeat(sp, IndentSize)...)
			if ipre == 0 || len(dm[ds.Name]) == 0 {
				buf = append(buf, sp...)
			} else {
				buf = append(buf, EdgeTypeLink...)
			}
		}
		buf = append(buf, strings.Repeat(sp, IndentSize)...)
		// }
		if len(dm[ds.RelDir]) == 0 {
			edge = EdgeTypeEnd
		}
		buf = append(buf, edge...)
		goto END
	}
	// jf != 0 files
END:
	return string(buf)
}

func getSubdirsMap(dirs []string) map[string][]string {
	dm := make(map[string][]string)
	for _, dir := range dirs {
		dm[dir] = []string{}
		dm[dir] = append(dm[dir], getSubdirsList(dir, dirs)...)
	}
	return dm
}

func getSubdirsList(root string, dirs []string) []string {
	sdir := []string{}
	for _, dir := range dirs {
		if strings.HasPrefix(dir, root) {
			path := strings.TrimPrefix(dir, root)
			ss := strings.Split(path, string(os.PathSeparator))
			if len(ss) < 2 {
				continue
			}
			ss = ss[1:]
			if !funk.ContainsString(sdir, ss[0]) {
				sdir = append(sdir, ss[0])
			}
		}
	}
	return sdir
}

func getSubdirsMeta(root string, dirs []string, folder map[string][]string) string {
	ndd := getNumSubdirs(root, dirs)
	ndf := len(folder[root])
	meta := fmt.Sprintf("[%d dirs., %d files] ", ndd, ndf)
	return meta
}

func getNumSubdirs(root string, dirs []string) int {
	nd := 0
	for _, dir := range dirs {
		if strings.HasPrefix(dir, root) {
			nd++
		}
	}
	return nd - 1
}

func transferToFileSummary(root, dir, file string) *paw.FileSummary {
	if len(file) == 0 {
		fullpath := filepath.Join(root, dir)
		fs := paw.NewFileSummary(fullpath, root)
		return fs
	}
	fullpath := filepath.Join(root, dir, file)
	fs := paw.NewFileSummary(fullpath, root)
	return fs
}

func transferToFileSummaryList(root string, dirs []string, folder map[string][]string) []*paw.FileSummary {
	fs := []*paw.FileSummary{}
	for _, d := range dirs {
		if len(folder[d]) == 0 {
			fullpath := filepath.Join(root, d)
			fs = append(fs, paw.NewFileSummary(fullpath, root))
			continue
		}
		for _, f := range folder[d] {
			fullpath := filepath.Join(root, d, f)
			fs = append(fs, paw.NewFileSummary(fullpath, root))
		}
	}
	return fs
}

func getLevelm(dirs []string) map[string]int {
	levelm := make(map[string]int)
	sep := "/"
	for _, d := range dirs {
		ss := strings.Split(d, sep)
		n := len(ss) - 1
		// fmt.Println(d, n, ss)
		levelm[d] = n
	}
	return levelm
}

func printFileSummary(fullpath, root string) {
	// var (
	// 	fullpath string
	// )
	// fullpath = filepath.Join(root, dirs[0])
	// printFileSummary(fullpath, root)

	// fullpath = filepath.Join(root, dirs[1])
	// printFileSummary(fullpath, root)

	// fullpath = filepath.Join(root, dirs[1], folder[dirs[1]][0])
	// printFileSummary(fullpath, root)

	// fullpath = filepath.Join(root, dirs[2])
	// printFileSummary(fullpath, root)

	// fullpath = filepath.Join(root, dirs[2], folder[dirs[2]][0])
	// printFileSummary(fullpath, root)
	fs := paw.NewFileSummary(fullpath, root)
	fmt.Printf("      root: %q\n", root)
	fmt.Printf("  fullpath: %q\n", fullpath)
	fmt.Printf("   AbsPath: %q\n", fs.AbsPath)
	fmt.Printf("    RelDir: %q\n", fs.RelDir)
	fmt.Printf("      Name: %q\n", fs.Name)
	fmt.Printf("       Ext: %q\n", fs.Ext)
	fmt.Printf("     IsDir: %v\n", fs.IsDir)
	fmt.Println()
}

func getDirsLevelsEnded(dirs []string) (dirsLE map[string][]int) {
	spew.Dump(dirs)
	sep := string(os.PathSeparator)
	dirsLE = make(map[string][]int)
	for _, dir := range dirs {
		var levelsEnded []int
		spdirs := strings.Split(dir, sep)
		// spdirs = spdirs[:]
		level := len(spdirs)
		fmt.Printf("%d %v\n", level, spdirs)
		if level < 0 {
			levelsEnded = append(levelsEnded, level)
		} else {
			for i := 0; i < level; i++ {
				levelsEnded = append(levelsEnded, i)
			}
		}
		fmt.Printf("%d %v\n", level, levelsEnded)

		if _, ok := dirsLE[dir]; !ok {
			dirsLE[dir] = []int{}
			dirsLE[dir] = append(dirsLE[dir], levelsEnded...)
		}
	}
	// spew.Dump(dirsLE)
	return dirsLE
}

func foutputTreeX1(w io.Writer, root string, dirs []string, folder map[string][]string) {

	// top := strings.Repeat("=", 80)
	// mid := strings.Repeat("-", 80)
	// buttom := top
	var (
		buf    = new(bytes.Buffer)
		nd, nf = 0, 0
		// prevDir, nextDir = "", ""
		nextDir    = ""
		hasNextDir bool
	)

	hasNextDir = true
	dirsLevelsEnded := getDirsLevelsEnded(dirs)
	iLastDir := len(dirs) - 1
	for i, dir := range dirs {
		// fmt.Println(filepath.Dir(dir), ":", filepath.Base(dir))
		nd++
		if i < iLastDir {
			nextDir = dirs[i+1]
		} else {
			nextDir = ""
		}
		var levelsEnded []int
		level := 0
		if len(dir) == 0 { // root
			path, _ := FileColorStr(root, root)
			root, _ := FileColorStr(root, "./")
			// buf.WriteString(fmt.Sprintf("%2d. %s", i+1, str))
			meta := fmt.Sprintf("%d", len(folder[dir]))
			val := fmt.Sprintf("%s » [%s]", root, path)
			buf.WriteString(fmt.Sprintf("%d. [%s] %s\n", i, meta, val))
		} else { // subfolder: len(dir) > 0
			// buf.WriteString(fmt.Sprintf("%v\n", strings.Split(dir, "/")))
			level = len(strings.Split(dir, "/")) - 1
			levelsEnded = dirsLevelsEnded[dir]
			fullpath := filepath.Join(root, dir)
			hasNextDir = strings.Contains(nextDir, dir) && len(levelsEnded) == len(dirsLevelsEnded[nextDir])
			edge := EdgeTypeMid
			if i == len(dirs)-1 || !hasNextDir {
				edge = EdgeTypeEnd
			}
			printTreeValues(buf, level, levelsEnded, edge, i, fullpath, folder[dir])
		}

		nf += len(folder[dir])
		// subdir := filepath.Join(root, dir)
		// hasNextDir = strings.Contains(nextDir, dir)
		if len(folder[dir]) > 0 {
			// printTreeFiles(buf, level+1, levelsEnded, hasNextDir, subdir, folder[dir])
			files := folder[dir]
			for i, f := range files {
				edge := EdgeTypeMid
				if i == len(files)-1 {
					// levelsEnded = append(levelsEnded, level)
					edge = EdgeTypeEnd
				}
				if hasNextDir {
					edge = EdgeTypeMid
				}

				fullpath := filepath.Join(root, dir, f)
				printTreeValues(buf, level+1, levelsEnded, edge, i+1, fullpath, nil)
			}
		}
		// prevDir = dir
	}
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("%d directories, %d files.", nd, nf))
	fmt.Fprintln(w, string(buf.Bytes()))
}

func printTreeFiles(wr io.Writer, level int, levelsEnded []int, hasNextDir bool, dir string, files []string) {
	for i, f := range files {
		edge := EdgeTypeMid
		if i == len(files)-1 {
			// levelsEnded = append(levelsEnded, level)
			edge = EdgeTypeEnd
		}
		if hasNextDir {
			edge = EdgeTypeMid
		}

		fullpath := filepath.Join(dir, f)
		printTreeValues(wr, level, levelsEnded, edge, i+1, fullpath, nil)
	}
}

func printTreeValues(wr io.Writer,
	level int, levelsEnded []int, edge EdgeType, idx int, fullpath string, files []string) {

	for i := 0; i < level-1; i++ {
		if isEnded(levelsEnded, i) {
			fmt.Fprint(wr, strings.Repeat(" ", IndentSize+1))
			continue
		}
		fmt.Fprintf(wr, "%s%s", EdgeTypeLink, strings.Repeat(" ", IndentSize))
	}

	val, _ := FileColorStr(fullpath, filepath.Base(fullpath))
	if files == nil || len(files) == 0 {
		fmt.Fprintf(wr, "%s %d. %s\n", edge, idx, val)
		return
	}

	fmt.Fprintf(wr, "%s %d. [%d]  %s\n", edge, idx, len(files), val)
}

func isEnded(levelsEnded []int, level int) bool {
	for _, l := range levelsEnded {
		if l == level {
			return true
		}
	}
	return false
}

// Table will return a string in table mode with `head` (use `FprintTable`)
func (m *PathMap) Table(head, pad string) string {
	buf := new(bytes.Buffer)
	SetNoColor()
	m.Fprint(buf, OTableFormatMode, head, pad)
	ResumNoColor()
	return paw.TrimFrontEndSpaceLine(string(buf.Bytes()))
}

// FprintTable print out in table mode with `head`
func (m *PathMap) FprintTable(w io.Writer, tf *paw.TableFormat, head string) {
	tf.Prepare(w)
	tf.SetBeforeMessage(head)
	foutputTable(tf, m.GetRoot(), m.GetDirs(), m.GetFolder())
}

func foutputTable(tf *paw.TableFormat, root string, dirs []string, folder map[string][]string) {

	nd, nf := 0, 0

	// tf.SetBeforeMessage(msg)
	tf.PrintSart()
	for i, dir := range dirs {
		level := len(strings.Split(dir, "/")) - 1
		dfiles := len(folder[dir])
		nd++
		// subhead := fmt.Sprintf("Depth: %d, .%s ( %d files)", level, dir, dfiles)
		fullpath := filepath.Join(root, dir)
		str, _ := FileColorStr(fullpath, dir)
		subhead := fmt.Sprintf("Depth: %d, .%s ( %d files)", level, str, dfiles)
		switch {
		case len(dir) == 0:
			level = 0
			nd--
			// tf.PrintRow("", fmt.Sprintf("Depth: %d, %s ( %d files)", level, root, dfiles))
			str, _ := FileColorStr(root, root)
			tf.PrintRow("", fmt.Sprintf("Depth: %d, %s ( %d files)", level, str, dfiles))
		case len(folder[dir]) == 0:
			tf.PrintRow("", subhead)
			goto MID
		default:
			tf.PrintRow("", subhead)
		}
		nf += len(folder[dir])
		level++
		for j, f := range folder[dir] {
			// tf.PrintRow(cast.ToString(j+1), f)
			fullpath := filepath.Join(root, dir, f)
			str, _ := FileColorStr(fullpath, f)
			tf.PrintRow(cast.ToString(j+1), str)

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
	SetNoColor()
	m.Fprint(buf, OPlainTextMode, head, pad)
	ResumNoColor()
	return paw.TrimFrontEndSpaceLine(string(buf.Bytes()))
}

// FprintText print out in plain text mode
func (m *PathMap) FprintText(w io.Writer, head, pad string) {
	fmt.Fprintln(w, paw.PaddingString(head, pad))
	fmt.Fprintln(w, pad)
	buf := new(bytes.Buffer)
	foutputText(buf, m.GetRoot(), m.GetDirs(), m.GetFolder())
	fmt.Fprintln(w, paw.PaddingString(string(buf.Bytes()), pad))
	fmt.Fprintln(w, pad)
}

func foutputText(w io.Writer,
	root string, dirs []string, folder map[string][]string) {

	top := strings.Repeat("=", 80)
	mid := strings.Repeat("-", 80)
	buttom := top
	nd, nf := 0, 0

	fprintWithLevel(w, 0, top)
	for i, dir := range dirs {
		level := len(strings.Split(dir, "/")) - 1
		nd++
		// subhead := fmt.Sprintf("%2d. %s", i+1, dir)
		fullpath := filepath.Join(root, dir)
		str, _ := FileColorStr(fullpath, dir)
		subhead := fmt.Sprintf("%2d. %s", i+1, str)
		switch {
		case len(dir) == 0:
			level = 0
			nd--
			// fprintWithLevel(w, level, fmt.Sprintf("%2d. %s", i+1, root))
			str, _ := FileColorStr(root, root)
			fprintWithLevel(w, level, fmt.Sprintf("%2d. %s", i+1, str))
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
			// fprintWithLevel(w, level, fmt.Sprintf("%2d. %s", j+1, f))
			fullpath := filepath.Join(root, dir, f)
			str, _ := FileColorStr(fullpath, f)
			fprintWithLevel(w, level, fmt.Sprintf("%2d. %s", j+1, str))
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

// FindFiles finds diles and fills into `PathMap` w.r.t. some conditions (`m.cond`, see `m.GetCondition()`)
// 	isRecursive:
// 		false to find files only in `root` directory
//		true  to find recursive files including subfolders
func (m *PathMap) FindFiles(root string, isRecursive bool) error {

	root, _ = filepath.Abs(root)
	// root = strings.TrimSuffix(root, "/")
	dirs := []string{}
	folder := make(map[string][]string)
	var err error
	if isRecursive {
		err = walkDir(root, &dirs, &folder, &m.cond)
	} else {
		err = ioReadDir(root, &dirs, &folder, &m.cond)
	}
	if err != nil {
		return err
	}

	m.SetRoot(root)
	m.SetDirs(dirs)
	m.SetFolder(folder)

	return nil
}

func osReadDir(root string, dirs *[]string, folder *map[string][]string, cond *pmCondition) error {

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
		if cond.ignoreHidden { // ignore hidden files
			if strings.HasPrefix(fi.Name(), ".") {
				continue
			}
		}
		if isAddFile(fi.Name(), cond) { // 是否執行過濾條件
			(*folder)[""] = append((*folder)[""], fi.Name())
		}
	}

	*dirs = append(*dirs, "")

	return nil
}

func ioReadDir(root string, dirs *[]string, folder *map[string][]string, cond *pmCondition) error {
	// root, _ = filepath.Abs(root)
	// root = strings.TrimSuffix(root, "/")
	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		if cond.ignoreHidden { // ignore hidden files
			if strings.HasPrefix(fi.Name(), ".") {
				continue
			}
		}
		if isAddFile(fi.Name(), cond) { // 是否執行過濾條件
			(*folder)[""] = append((*folder)[""], fi.Name())
		}
	}

	*dirs = append(*dirs, "")

	return nil
}

func walkDir(root string, dirs *[]string, folder *map[string][]string, cond *pmCondition) error {

	visitFile := func(path string, info os.FileInfo, err error) error {
		// fmt.Println(path)
		if err != nil {
			paw.Logger.Errorln(err) // can't walk here,
			return nil              // but continue walking elsewhere
		}

		// apath, _ := filepath.Abs(path)
		// base := filepath.Base(apath)
		base := info.Name()
		sub := strings.TrimPrefix(path, root)

		if cond.ignoreHidden { // ignore hidden files
			pl := strings.Split(path, "/")
			for _, p := range pl {
				if strings.HasPrefix(p, ".") {
					return nil
				}
			}
		}

		if info.IsDir() { // 子目錄
			if cond.ignoreCondition { // 執行過濾條件
				// 過濾被忽略的資料夾 (資料夾名完全相同)
				if isInArray(&cond.ignorePath, base) {
					return filepath.SkipDir
				}
			}
			if _, ok := (*folder)[sub]; !ok {
				(*folder)[sub] = []string{}
				(*dirs) = append(*dirs, sub)
			}
		} else { // 檔案
			sub = strings.TrimSuffix(sub, "/"+base)
			// sub = strings.TrimSuffix(sub, "/")
			if isAddFile(base, cond) { // 是否執行過濾條件
				(*folder)[sub] = append((*folder)[sub], base)
			}
		}
		return nil
	}

	err := filepath.Walk(root, visitFile)
	if err != nil {
		return err
	}
	return nil
}

func isAddFile(base string, c *pmCondition) bool {
	if !c.ignoreCondition { // 不執行過濾條件
		return true // 加入檔案
	}
	// 執行過濾條件
	// 目標檔案型別被指定
	if !isAllEmpty(&c.targetType) {
		// 屬於目標檔案型別
		if isInSuffix(&c.targetType, base) {
			// 忽略檔案為空 或者 目標檔案中不含有指定忽略檔案
			if isAllEmpty(&c.ignoreFile) || !isInArray(&c.ignoreFile, base) {
				return true // 加入檔案
			}
		}
	} else { // 目標檔案型別為空
		// fmt.Printf("%v %q\n", ignoreType, base)
		// 忽略檔案型別被指定
		if !isAllEmpty(&c.ignoreType) {
			// 不屬於忽略檔案型別
			if !isInSuffix(&c.ignoreType, base) {
				// 忽略檔案為空 或者 目標檔案中不含有指定忽略檔案
				if isAllEmpty(&c.ignoreFile) || !isInArray(&c.ignoreFile, base) {
					return true // 加入檔案
				}
			}
		} else { // 忽略檔案型別為空
			// 忽略檔案為空 或者 目標檔案中不含有指定忽略檔案
			if isAllEmpty(&c.ignoreFile) || !isInArray(&c.ignoreFile, base) {
				return true // 加入檔案
			}
		}
	}

	return false // 不加入檔案
}
