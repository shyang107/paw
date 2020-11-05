package paw

import (
	"strings"
	"unicode"
	"unicode/utf8"
	// "github.com/gookit/color"
)

// var Color = color.Green

// GetAbbrString return a abbreviation string 'xxx...' of `str` with
// maximum length `maxlen`
func GetAbbrString(str string, maxlen int) string {
	hc, ac := CountPlaceHolder(str)
	lenStr := hc + ac
	if lenStr < maxlen {
		return str
	}

	ignore := "..."
	limit := maxlen - len(ignore)
	c := 0
	sb := strings.Builder{}
	for _, ch := range str {
		rl := utf8.RuneLen(ch)
		if rl == 3 {
			c += 2
		} else {
			c++
		}
		if c <= limit {
			sb.WriteRune(ch)
		} else {
			break
		}
	}
	if c < limit {
		for i := 0; i < limit-c; i++ {
			sb.WriteString(" ")
		}
	}
	str = sb.String() + ignore
	return str
}

// CountPlaceHolder return `nHan` and `nASCII`
//    `nHan`: number of occupied space in terminal for han-character
//    `nASCII`: number of occupied space in terminal for ASCII-character
func CountPlaceHolder(str string) (nHan int, nASCII int) {
	nHan, nASCII = 0, 0
	for _, ch := range str {
		rl := utf8.RuneLen(ch)
		if rl == 3 {
			nHan += 2
		} else {
			nASCII++
		}
	}
	return nHan, nASCII
}

// HasChineseChar return true for that `str` include chinese character
func HasChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}
