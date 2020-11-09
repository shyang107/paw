package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/funk"
	// "github.com/thoas/go-funk"
)

var (
	lg = paw.Log
)

func init() {
	// paw.SetLogLevel(paw.InfoLevel)
}

func main() {
	// testLineCount()
	// testFileLineCount()
	// rehttp()
	// testGetAbbrString()
	// exTableFormat()
	exStringBuilder()
}

func exStringBuilder() {
	lg.Info("exStringBuilder")
	s := "Text中文 Collection"
	tb := &paw.TextBuilder{}
	tc := tb.SetText(s).Build()
	fmt.Println("                         s:", s)
	fmt.Println("     tb.SetText(s).Build():", tc.Text)
	fmt.Println("tb.NumberBanner().String():", tb.NumberBanner())
	fmt.Println(`   b.GetAbbrString(8, "»"):`, tb.GetAbbrString(8, "»"))
	fmt.Println(`tb.NumberBanner().String():`, tb.GetAbbrString(8, "»").NumberBanner())
	h, a := tb.CountPlaceHolder()
	fmt.Println("     tb.CountPlaceHolder():", h, a)
	fmt.Println("       tb.HasChineseChar():", tb.HasChineseChar())
	fmt.Println("     tb.SetText(s).Build():", tc.Text)
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
