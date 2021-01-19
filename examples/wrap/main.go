package main

import (
	"fmt"

	"github.com/shyang107/paw"
)

func main() {
	str := "測序測序123 用檔111案名測1序測序用1檔案名測序測序用檔案名01234abcde測序用檔案名01234abcde測序用檔案名01234abcde測序用檔案名01234abcde測序用檔案名01234abcde測序用檔案名01234abcde測序用檔案名01234abcde測序用檔案名01234abcde"
	fmt.Printf("%v\n", str)
	ws := paw.StringWidth(str)
	paw.StringWidth(str)
	fmt.Println(paw.NumberBanner(ws))
	fmt.Println("width =", ws)

	w := 25
	fmt.Printf(">> warp to string with width %d \n", w)
	fmt.Println(paw.NumberBanner(w))
	wstr := paw.Wrap(str, w)
	fmt.Println(wstr)

	fmt.Printf(">> truncate to string with width %d \n", w)
	fmt.Println(paw.NumberBanner(w))
	tstr := paw.Truncate(str, w, "")
	fmt.Println(tstr, "\t", paw.StringWidth(tstr))

	fmt.Printf(">> meta test \n")
	fmt.Println(paw.NumberBanner(w))
	meta := "[test]"
	wm := paw.StringWidth(meta)
	w1 := w - wm - 1
	fmt.Println("meta =", meta, "w =", w, "wm =", wm, "w-wm-1 =", w1)
	sl1 := paw.Truncate(str, w1, "")
	fmt.Println(paw.NumberBanner(w))
	w1 = len(sl1)
	fmt.Printf("%s %s\t%d\n", meta, sl1, w1)
	rstr := str[w1:]
	rs := paw.Wrap(rstr, w)
	fmt.Println(rs)

}
func wrap(s string, w int) string {
	if paw.StringWidth(s) <= w {
		return s
	}
	// width := 0
	out := ""
	for _, r := range []rune(s) {
		sr := string(r)
		cw := paw.RuneWidth(r)
		csw := paw.StringWidth(sr)
		fmt.Println(sr, "rw =", cw, "sw =", csw)
	}
	return out
}
