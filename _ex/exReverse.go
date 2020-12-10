package main

import (
	"fmt"

	"github.com/shyang107/paw"
)

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
