package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/shyang107/paw/cast"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/funk"

	"github.com/shyang107/paw/treeprint"
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
	// testLineCount()
	// testFileLineCount()
	// rehttp()
	// testGetAbbrString()
	// exTableFormat()
	// exStringBuilder()
	// exLoger()
	// exReverse()
	// exPrintTree()
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
	exGrouppingFiles2()
}

func exGrouppingFiles2() {
	paw.Logger.Info("")
	// sourceFolder := "../"
	sourceFolder, _ := homedir.Expand("~/Downloads/")
	// sourceFolder := "/Users/shyang/go/src/rover/opcc/"
	isRecursive := true
	sourceFolder, err := filepath.Abs(sourceFolder)
	if err != nil {
		paw.Logger.Error(err)
	}
	sourceFolder += "/"
	hsb := strings.Builder{}
	hsb.WriteString("GetFilesFuncString:\n")
	hsb.WriteString("  sourceFolder: \"" + sourceFolder + "\"\n")
	hsb.WriteString("   isRecursive: " + strconv.FormatBool(isRecursive) + "\n")
	prefix := "."
	regexPattern := `\.git|\$RECYCLE\.BIN`
	re := regexp.MustCompile(regexPattern)
	hsb.WriteString("  Exculde:" + "\n")
	hsb.WriteString(`          prefix: "` + prefix + `"` + "\n")
	hsb.WriteString(`    regexPattern: "` + regexPattern + `"`)

	tp := &paw.TableFormat{
		Fields:    []string{"No.", "Sorted Files"},
		LenFields: []int{5, 80},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft},
		Padding:   "# ",
	}
	tp.Prepare(os.Stdout)
	tp.SetBeforeMessage(hsb.String())
	tp.PrintSart()

	files, err := paw.GetFilesFunc(sourceFolder, isRecursive,
		func(f paw.File) bool {
			return (len(f.FileName) == 0 || paw.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
		})
	if err != nil {
		paw.Logger.Error(err)
	}

	paw.GrouppingFiles(files)

	oFolder := files[0].Folder
	gcount := 0
	for i, f := range files {
		path := paw.TrimPrefix(f.FullPath, sourceFolder)
		if oFolder != f.Folder {
			oFolder = f.Folder
			tp.PrintRow("", "Sum: "+cast.ToString(gcount)+" files.")
			tp.PrintMiddleSepLine()
			gcount = 1
		} else {
			gcount++
		}
		tp.PrintRow(gcount, path)
		if i == len(files)-1 {
			tp.PrintRow("", "Sum: "+cast.ToString(gcount)+" files.")
		}
	}
	tp.SetAfterMessage("Total: " + cast.ToString(len(files)) + " files. ")
	tp.PrintEnd()
}
func exGrouppingFiles1() {
	paw.Logger.Info("")
	// sourceFolder := "../"
	sourceFolder := "/Users/shyang/go/src/rover/opcc/"
	isRecursive := true
	sourceFolder, err := filepath.Abs(sourceFolder)
	if err != nil {
		paw.Logger.Error(err)
	}
	sourceFolder += "/"
	hsb := strings.Builder{}
	hsb.WriteString("GetFilesFuncString:\n")
	hsb.WriteString("  sourceFolder: " + `"../" <- "` + sourceFolder + "\"\n")
	hsb.WriteString("   isRecursive: " + strconv.FormatBool(isRecursive) + "\n")
	prefix := "."
	regexPattern := `\.git`
	re := regexp.MustCompile(regexPattern)
	hsb.WriteString("  Exculde:" + "\n")
	hsb.WriteString(`          prefix: "` + prefix + `"` + "\n")
	hsb.WriteString(`    regexPattern: "` + regexPattern + `"`)

	tp := &paw.TableFormat{
		Fields:    []string{"No.", "File", "Sorted Files"},
		LenFields: []int{5, 40, 40},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft, paw.AlignLeft},
		Padding:   "# ",
	}
	tp.Prepare(os.Stdout)
	tp.SetBeforeMessage(hsb.String())
	tp.PrintSart()

	files, err := paw.GetFilesFunc(sourceFolder, isRecursive,
		func(f paw.File) bool {
			return len(f.FileName) == 0 || strings.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath)
		})
	if err != nil {
		paw.Logger.Error(err)
	}
	sfiles := make([]paw.File, len(files))
	copy(sfiles, files)
	paw.GrouppingFiles(sfiles)

	oFolder := sfiles[0].Folder
	for i, f := range files {
		path := strings.TrimPrefix(f.FullPath, sourceFolder)
		spath := strings.TrimPrefix(sfiles[i].FullPath, sourceFolder)
		j := i + 1
		// if j%5 == 0 {
		if oFolder != sfiles[i].Folder {
			oFolder = sfiles[i].Folder
			tp.PrintMiddleSepLine()
		}
		tp.PrintRow(j, path, spath)
	}
	tp.PrintEnd()
}

func exGetFilesString() {
	paw.Logger.Info("")
	sourceFolder := "../"
	isRecursive := true
	sourceFolder, err := filepath.Abs(sourceFolder)
	if err != nil {
		paw.Logger.Error(err)
	}
	sourceFolder += "/"
	hsb := strings.Builder{}
	hsb.WriteString("GetFilesFuncString:\n")
	hsb.WriteString("  sourceFolder: " + `"../" <- "` + sourceFolder + "\"\n")
	hsb.WriteString("   isRecursive: " + strconv.FormatBool(isRecursive) + "\n")
	prefix := "."
	regexPattern := `\.git`
	re := regexp.MustCompile(regexPattern)
	hsb.WriteString("  Exculde:" + "\n")
	hsb.WriteString(`          prefix: "` + prefix + `"` + "\n")
	hsb.WriteString(`    regexPattern: "` + regexPattern + `"`)

	tp := &paw.TableFormat{
		Fields:    []string{"No.", "File"},
		LenFields: []int{5, 72},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft},
		Padding:   "# ",
	}
	tp.Prepare(os.Stdout)
	tp.SetBeforeMessage(hsb.String())
	tp.PrintSart()

	files, err := paw.GetFilesFuncString("../", isRecursive,
		func(f paw.File) bool {
			return !(len(f.FileName) == 0 || strings.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
		})
	if err != nil {
		paw.Logger.Error(err)
	}

	for i, f := range files {
		path := strings.TrimPrefix(f, sourceFolder)
		j := i + 1
		tp.PrintRow(j, path)
		if j%5 == 0 {
			tp.PrintMiddleSepLine()
		}
	}
	tp.PrintEnd()
}

func exGetFiles3() {
	paw.Logger.Info("exGetFiles3")
	sourceFolder := "../"
	sourceFolder, err := filepath.Abs(sourceFolder)
	if err != nil {
		paw.Logger.Error(err)
	}
	sourceFolder += "/"
	head := "\nGetFilesFunc: folder <- '../', isRecursive <- true\n"
	head += "  sourceFolder: " + sourceFolder + "\n"
	prefix := "."
	regexPattern := `\.git`
	re := regexp.MustCompile(regexPattern)
	head += "  Exculde:" + "\n"
	head += fmt.Sprintf("          prefix: %q\n", prefix)
	head += fmt.Sprintf("    regexPattern: %q", regexPattern)

	tp := &paw.TableFormat{
		Fields:    []string{"No.", "File"},
		LenFields: []int{5, 72},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft},
		Padding:   "  ",
	}
	tp.Prepare(os.Stdout)
	tp.SetBeforeMessage(head)
	tp.PrintSart()

	files, err := paw.GetFilesFunc("../", true, func(f paw.File) bool {
		return (len(f.FileName) == 0 || strings.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
	})
	if err != nil {
		paw.Logger.Error(err)
	}

	for i, f := range files {
		newPath, err := paw.GetNewFilePath(f, sourceFolder, "./")
		if err != nil {
			paw.Logger.Error(err)
		}
		rows := []interface{}{i + 1, newPath}
		tp.PrintRow(rows...)
	}
	tp.PrintEnd()
}

func exGetFiles2() {
	paw.Logger.Info("exGetFiles2")
	paw.Logger.Info("GetFiles: folder <- '../', isRecursive <- true")
	sourceFolder := "../"
	fmt.Println("sourceFolder:", sourceFolder)
	files, err := paw.GetFiles(sourceFolder, true)
	if err != nil {
		paw.Logger.Error(err)
	}
	for i, f := range files {
		fmt.Printf("%3d. %s\n", i, f.FullPath)
	}
	i := 0
	prefix := "."
	regexPattern := `\.git`
	re := regexp.MustCompile(regexPattern)
	fmt.Println("Exculde:")
	fmt.Printf("\t      prefix: %q\n", prefix)
	fmt.Printf("\tregexPattern: %q\n", regexPattern)
	for _, f := range files {
		if strings.HasPrefix(f.FileName, prefix) {
			continue
		} else if len(f.FileName) == 0 {
			continue
		} else if re.MatchString(f.FullPath) {
			continue
		}
		i++
		fmt.Printf("%3d. %s\n", i, f.FullPath)
	}
}
func exGetFiles1() {
	paw.Logger.Info("exGetFiles1")
	paw.Logger.Info("GetFiles: folder <- '~/', isRecursive <- false")
	homepath, err := homedir.Dir()
	if err != nil {
		paw.Logger.Error(err)
	}
	files, err := paw.GetFiles(homepath, false)
	if err != nil {
		paw.Logger.Error(err)
	}
	for i, f := range files {
		fmt.Printf("%2d. %s\n", i+1, f.FullPath)
	}
}

func exFolder() {
	paw.Logger.Info("exFolder")
	path := "/aaa/bbb/ccc/example.xxx"
	fmt.Println("                            path:", path)
	file := paw.ConstructFile(path)
	fmt.Println("ConstructFile(path):")
	spew.Dump(file)
	sourceFolder := "/aaa/bbb/"
	fmt.Println("                    sourceFolder:", sourceFolder)
	subfolder := paw.GetSubfolder(file, sourceFolder)
	fmt.Println("GetSubfolder(file, sourceFolder):", subfolder)
	targetFolder := "ddd/"
	fmt.Println("                    targetFolder:", targetFolder)
	newFolder, _ := paw.GetNewFilePath(file, sourceFolder, targetFolder)
	fmt.Println("GetNewFilePath(file, sourceFolder, targetFolder):", newFolder)
}

func exLogger() {
	paw.Logger.Info("exLogger")
	paw.Logger.Debug("exLogger")
	paw.Logger.Warn("exLogger")
	paw.Logger.Trace("exLogger")
	fmt.Println("  GetDotDir()", paw.GetDotDir())
	fmt.Println("GetCurrPath()", paw.GetCurrPath())
	fmt.Println("  GetAppDir()", paw.GetAppDir())

}

const tmpText = "const twoMatch = `test string`;\nconst noMatches = `test ${ variabel }`;\nabcde ${ field1 } and ${ Field2}"

func exRegEx() {
	var re = regexp.MustCompile(`(?m)(\${.*?)(\b\w+\b)(.*?})`)
	tokens := map[string]string{
		"variabel": "[token_variabel]",
		"field1":   "[token_field1]",
		"field2":   "[token_field2]",
	}
	fmt.Println(tmpText)
	matchs := re.FindAllStringSubmatch(tmpText, -1)
	spew.Dump(matchs)
	tb := paw.TextBuilder{}
	result := tmpText
	for _, m := range matchs {
		tb.SetText(m[2]).ToLower()
		result = strings.ReplaceAll(result, m[0], tokens[tb.GetText()])
	}
	fmt.Println(result)
	// fmt.Println(re.ReplaceAllString(str, substitution))
}

// func exMergeSort(n []int) {
// 	fmt.Println("MergeSort\n", n)
// 	n = paw.MergeSort(n)
// 	fmt.Println(n)
// 	// paw.CombSortFunc(n, 1.8, func(a, b int) bool { return a < b })
// 	// fmt.Println(n)
// }
// func exCombSort(n []int) {
// 	// n := paw.GenerateSlice(8)
// 	fmt.Println("CombSort\n", n)
// 	paw.CombSort(n, 1.8)
// 	fmt.Println(n)
// 	paw.CombSortFunc(n, 1.8, func(a, b int) bool { return a < b })
// 	fmt.Println(n)
// }

// SelectionSort 選擇排序
func SelectionSort(n []int) {
	fmt.Println("SelectionSort\n", n)
	// count := 0
	// for i := 0; i < len(n); i++ {
	// 	minIndex := i
	// 	for j := i + 1; j < len(n); j++ {
	// 		if n[minIndex] > n[j] {
	// 			minIndex = j
	// 		}
	// 	}
	// 	n[i], n[minIndex] = n[minIndex], n[i]
	// 	count++
	// 	fmt.Println(count, n)
	// }
	paw.SelectionSort(n)
	fmt.Println(n)
	paw.SelectionSortFunc(n, func(a, b int) bool { return a < b })
	fmt.Println(n)

}

// InsertionSort 插入排序
func InsertionSort(n []int) {
	fmt.Println("InsertionSort\n", n)
	// count := 0
	// i := 1
	// for i < len(a) {
	// 	j := i
	// 	for j >= 1 && a[j] < a[j-1] {
	// 		a[j-1], a[j] = a[j], a[j-1]
	// 		count++
	// 		fmt.Println(count, a)
	// 		j--
	// 	}
	// 	i++
	// }
	paw.InsertionSort(n)
	fmt.Println(n)
	paw.InsertionSortFunc(n, func(a, b int) bool { return a < b })
	fmt.Println(n)
}
func exGetCurrPath() {
	fmt.Println(paw.GetCurrPath())
}

func exShuffle() {
	s := []rune("abcdefg")
	slice := make([]interface{}, len(s))
	for i, val := range s {
		slice[i] = string(val)
	}
	fmt.Println(slice)
	for i := 0; i < 10; i++ {
		paw.Shuffle(slice)
		fmt.Println(slice)
	}
}
func exPrintTree() {
	data := []treeprint.Org{
		// {"A001", "Dept1", "0 -----th top"},
		{"A001", "Dept1", "0"},
		{"A011", "Dept2", "0"},
		{"A002", "subDept1", "A001"},
		{"A005", "subDept2", "A001"},
		{"A003", "sub_subDept1", "A002"},
		{"A006", "gran_subDept", "A003"},
		{"A004", "sub_subDept2", "A002"},
		{"A012", "subDept1", "A011"},
	}

	treeprint.PrintOrgTree("ORG", data, "0", 3)
}
func exReverse() {
	lg.Info("exReverse")
	s := "Text中文 Collection"
	tb := paw.TextBuilder{}
	tb.SetText(s)
	fmt.Println("           s:", s)
	fmt.Println("     tb.Text:", tb.GetText())
	fmt.Println("tb.Reverse():", tb.Reverse())
	fmt.Println("     tb.Text:", tb.GetText())
	fmt.Println("tb.Reverse():", tb.Reverse())

}
func exLoger() {
	log.Infoln("飛雪無情的博客:", "http://www.flysnow.org")
	log.Warnln("飛雪無情的微信公眾號：%s\n", "flysnow_org")
	log.Errorln("歡迎關注留言")

	lg.Infoln("飛雪無情的博客:", "http://www.flysnow.org")
	lg.Warnln("飛雪無情的博客:", "http://www.flysnow.org")
	lg.Debugln("飛雪無情的博客:", "http://www.flysnow.org")
	lg.Errorln("飛雪無情的博客:", "http://www.flysnow.org")
	// lg.Traceln("飛雪無情的博客:", "http://www.flysnow.org")
	// lg.Fatalln("飛雪無情的博客:", "http://www.flysnow.org")

}
func exStringBuilder() {
	lg.Info("exStringBuilder")
	s := "Text中文 Collection"
	tb := &paw.TextBuilder{}
	tb.NewTextBuilder(s)
	fmt.Println("                         s:", s)
	fmt.Println("              tb.GetText():", tb.GetText())
	fmt.Println("tb.NumberBanner().String():", tb.NumberBanner())
	fmt.Println(`   b.GetAbbrString(8, "»"):`, tb.GetAbbrString(8, "»"))
	fmt.Println(`tb.NumberBanner().String():`, tb.GetAbbrString(8, "»").NumberBanner())
	h, a := tb.CountPlaceHolder()
	fmt.Println("     tb.CountPlaceHolder():", h, a)
	fmt.Println("       tb.HasChineseChar():", tb.HasChineseChar())
	fmt.Println("              tb.GetText():", tb.GetText())
	// out:
	// [I 2020-11-08 12:28:38 main:33] exStringBuilder
	//                          s: Text中文 Collection
	//      tb.SetText(s).Build(): Text中文 Collection
	// tb.NumberBanner().String(): 1234567890123456789
	//    b.GetAbbrString(8, "»"): 123456»
	// tb.NumberBanner().String(): 1234567
	//      tb.CountPlaceHolder(): 0 7
	//        tb.HasChineseChar(): false
	//      tb.SetText(s).Build(): Text中文 Collection
}

func exTableFormat() {
	lg.Info("exTableFormat")
	// t := paw.NewTableFormat()
	t := &paw.TableFormat{
		Fields:    []string{"No.", "Field 1", "Field 2", "Field 3", "Field 4", "Field 5"},
		LenFields: []int{3, 10, 10, 10, 15, 20},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignRight, paw.AlignRight, paw.AlignLeft, paw.AlignRight, paw.AlignCenter},
		Padding:   "# ",
		// Sep:       " ",
		// TopChar:    "*",
		// MiddleChar: "-",
		// BottomChar: "^",
	}
	t.Prepare(os.Stdout)
	t.SetBeforeMessage("Table: test\ntest\ntest")
	t.PrintSart()
	row := make([]interface{}, len(t.Fields))
	nr := 2
	for i := 0; i < nr; i++ {
		row[0] = strconv.Itoa(i + 1)
		for j := 1; j < len(t.Fields); j++ {
			row[j] = funk.RandomString(funk.RandomInt(3, 15), []rune("abcdefg中文huaijklmnopq1230456790"))
		}
		t.PrintRow(row...)
	}
	t.SetAfterMessage("Total " + strconv.Itoa(nr) + " records")
	t.PrintEnd()
}

func ds(len int) string {
	lg.Info("ds")
	return paw.NumberBanner(len)
}

func testGetAbbrString() {
	lg.Info("testGetAbbrString")
	str := "測試中 ab 文測試中文，裝使 ab Cde 中文中文 aaaaa 中文"
	fmt.Println(str)
	hc, ac := paw.CountPlaceHolder(str)
	fmt.Println(ds(hc+ac), hc+ac, "hc:", hc, "ac:", ac)
	maxlen := 28
	abbr := paw.GetAbbrString(str, maxlen, "»")
	fmt.Println(abbr)
	fmt.Println(ds(maxlen), maxlen)

}
func testFileLineCount() {
	lg.Info("testFileLineCount")
	lc, err := paw.FileLineCount("../README.md")
	if err != nil {
		lg.Error(err)
	}
	fmt.Println("FileLineCount:", lc)

}
func testLineCount() {
	lg.Info("testLineCount")
	fr, err := os.Open("../README.md")
	defer func() {
		if err != nil {
			lg.Error(err)
		}
		fr.Close()
	}()

	br := bufio.NewReader(fr)
	lc, err := paw.LineCount(br)
	if err != nil {
		lg.Error(err)
	}
	fmt.Println("LineCount:", lc)

}

func rehttp() {
	lg.Info("rehttp")
	var re = regexp.MustCompile(`(?m)(https|http|ftp):\/\/([\w\.\/\-\#]+)`)
	var str = `- _IndexOfFloat32: https://godoc.org/github.com/thoas/go-funk#IndexOfFloat32
- _IndexOfFloat64: https://godoc.org/github.com/thoas/go-funk#IndexOfFloat64
- _IndexOfInt: https://godoc.org/github.com/thoas/go-funk#IndexOfInt
- _IndexOfInt64: https://godoc.org/github.com/thoas/go-funk#IndexOfInt64
- _IndexOfString: https://godoc.org/github.com/thoas/go-funk#IndexOfString`

	for i, match := range re.FindAllString(str, -1) {
		fmt.Println(match, "found at index", i)
	}
	for i, match := range re.FindAllStringSubmatch(str, -1) {
		fmt.Println(match, "found at index", i)
	}
}
