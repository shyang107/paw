package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/funk"
	"github.com/shyang107/paw/log"
	"github.com/shyang107/paw/treeprint"
	// "github.com/shyang107/paw/Log"
	// "github.com/thoas/go-funk"
)

var (
	// lg = paw.Glog
	lg = paw.Logger
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
	exGetCurrPath()
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
	log.Info.Println("飛雪無情的博客:", "http://www.flysnow.org")
	log.Warn.Printf("飛雪無情的微信公眾號：%s\n", "flysnow_org")
	log.Error.Println("歡迎關注留言")

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
	tb.Build(s)
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
		// Padding:   "  ",
		// Sep:       " ",
		// TopChar:    "*",
		// MiddleChar: "-",
		// BottomChar: "^",
	}
	t.Prepare(os.Stdout)
	t.SetBeforeMessage("Table: test")
	t.PrintSart()
	row := make([]string, len(t.Fields))
	nr := 2
	for i := 0; i < nr; i++ {
		row[0] = strconv.Itoa(i + 1)
		for j := 1; j < len(t.Fields); j++ {
			row[j] = funk.RandomString(funk.RandomInt(3, 15), []rune("abcdefg中文huaijklmnopq1230456790"))
		}
		t.PrintRow(row)
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
