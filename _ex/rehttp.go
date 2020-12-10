package main

import (
	"fmt"
	"regexp"
)

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
