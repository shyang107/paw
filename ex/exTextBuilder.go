package main

import (
	"fmt"

	"github.com/shyang107/paw"
)

func exTextBuilder() {
	lg.Info("exTextBuilder")
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
