package main

import (
	"os"
	"strconv"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/funk"
)

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
