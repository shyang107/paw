package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/shyang107/paw"
	// "github.com/shyang107/paw/funk"
	"github.com/thoas/go-funk"

	"github.com/keakon/golog"
)

var lg = golog.NewStderrLogger()

func main() {
	// testLineCount()
	// testFileLineCount()
	// rehttp()
	testGetAbbrString()
	exTableFormat()
}

func exTableFormat() {
	// t := paw.NewTableFormat()
	t := &paw.TableFormat{
		Fields:    []string{"No.", "Field 1", "Field 2", "Field 3", "Field 4", "Field 5"},
		LenFields: []int{3, 7, 8, 10, 11, 15},
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
	nr := 6
	for i := 1; i < nr; i++ {
		row[0] = strconv.Itoa(i)
		for j := 1; j < len(t.Fields); j++ {
			row[j] = funk.RandomString(funk.RandomInt(3, 15), []rune("abcdefg中文huaijklmnopq1230456790"))
		}
		t.PrintRow(row)
	}
	t.SetAfterMessage("Total " + strconv.Itoa(nr) + " records")
	t.PrintEnd()
}
func ds(len int) string {
	return paw.NumberBanner(len)
}

func testGetAbbrString() {
	str := "測試中 ab 文測試中文，裝使 ab Cde 中文中文 aaaaa 中文"
	fmt.Println(str)
	hc, ac := paw.CountPlaceHolder(str)
	fmt.Println(ds(hc+ac), hc+ac, "hc:", hc, "ac:", ac)
	maxlen := 28
	abbr := paw.GetAbbrString(str, maxlen, "...")
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
