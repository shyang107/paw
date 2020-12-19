package main

import (
	"fmt"

	"github.com/shyang107/paw"
)

func ds(len int) string {
	lg.Info("ds")
	return paw.NumberBanner(len)
}

func exGetAbbrString() {
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
