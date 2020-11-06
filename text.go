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
func GetAbbrString(str string, maxlen int, ignore string) string {
	hc, ac := CountPlaceHolder(str)
	lenStr := hc + ac
	if lenStr <= maxlen {
		return str
	}
	if len(ignore) < 1 {
		ignore = "..."
	}
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
	hc, ac = CountPlaceHolder(sb.String())
	c = hc + ac
	if c < limit {
		for i := 0; i < limit-c; i++ {
			sb.WriteString(" ")
		}
	}
	str = sb.String() + ignore
	return str
}

// CountPlaceHolder return `nHan` and `nASCII`
// 	`nHan`: number of occupied space in screen for han-character
// 	`nASCII`: number of occupied space in screen for ASCII-character
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
//
// Example:
// 	HasChineseChar("abc 中文") return true
// 	HasChineseChar("abccefgh") return false
func HasChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

// NumberBanner return numbers' string with length `len`
//
// Example:
// 	NumberBanner(11) return "12345678901"
func NumberBanner(len int) string {
	nl := []byte("1234567890")
	sb := strings.Builder{}
	for i := 0; i < len; i++ {
		c := nl[i%10]
		sb.Write([]byte{c})
	}
	return sb.String()
}
