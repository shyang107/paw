package paw

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"

	"golang.org/x/text/transform"
	// "github.com/gookit/color"
)

// // HasPrefix return `strings.HasPrefix(str, prefix)`
// func HasPrefix(str string, prefix string) bool {
// 	return strings.HasPrefix(str, prefix)
// }

// // HasSuffix return `strings.HasSuffix(str, Suffix)`
// func HasSuffix(str string, suffix string) bool {
// 	return strings.HasSuffix(str, suffix)
// }

// // Contains return `strings.Contains(str, substr)`
// func Contains(str string, substr string) bool {
// 	return strings.Contains(str, substr)
// }

// // ContainsAny reports whether any Unicode code points in `chars` are within `s`.
// // 	Encapsulates ContainsAny(s, chars string) bool
// func ContainsAny(s, chars string) bool {
// 	return strings.ContainsAny(s, chars)
// }

// // ContainsRune reports whether the Unicode code point `r` is within `s`.
// // 	Encapsulates strings.ContainsRune(s string, r rune) bool
// func ContainsRune(s string, r rune) bool {
// 	return strings.ContainsRune(s, r)
// }

// // EqualFold reports whether `s` and `t`, interpreted as UTF-8 strings, are equal under Unicode case-folding, which is a more general form of case-insensitivity.
// // 	Encapsulates strings.EqualFold(s, t string) bool
// func EqualFold(s, t string) bool {
// 	return strings.EqualFold(s, t)
// }

// // Fields splits the string s around each instance of one or more consecutive white space characters, as defined by unicode.IsSpace, returning a slice of substrings of s or an empty slice if s contains only white space.
// //
// //	Encapsulates strings.Fields(s string) []string
// //
// // Example
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))
// // 	}
// // 	result: Fields are: ["foo" "bar" "baz"]
// func Fields(s string) []string {
// 	return strings.Fields(s)
// }

// // FieldsFunc splits the string s at each run of Unicode code points c satisfying f(c) and returns an array of slices of s. If all code points in s satisfy f(c) or the string is empty, an empty slice is returned.
// //
// // FieldsFunc makes no guarantees about the order in which it calls f(c) and assumes that f always returns the same value for a given c.
// //
// // Example
// //	package main
// //
// //	import (
// //		"fmt"
// //		"strings"
// //		"unicode"
// //	)
// //
// //	func main() {
// //		f := func(c rune) bool {
// //			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
// //		}
// //		fmt.Printf("Fields are: %q", strings.FieldsFunc("  foo1;bar2,baz3...", f))
// //	}
// //	result: Fields are: ["foo1" "bar2" "baz3"]
// func FieldsFunc(s string, f func(rune) bool) []string {
// 	return strings.FieldsFunc(s, f)
// }

// // Trim returns a slice of the string `s` with all leading and trailing Unicode code points contained in `cutset` removed.
// func Trim(s, cutset string) string {
// 	return strings.Trim(s, cutset)
// }

// // TrimFunc returns a slice of the string `s` with all leading and trailing Unicode code points `c` satisfying `f(c)` removed.
// //
// // Example:
// // 	fmt.Print(strings.TrimFunc("¡¡¡Hello, Gophers!!!", func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r)}))
// // 	out: Hello, Gophers
// func TrimFunc(s string, f func(rune) bool) string {
// 	return strings.TrimFunc(s, f)
// }

// // Index returns the index of the first instance of substr in s, or -1 if substr is not present in s.
// //
// // Example
// //
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Println(strings.Index("chicken", "ken")) // 4
// // 		fmt.Println(strings.Index("chicken", "dmr")) // -1
// // 	}
// func Index(s, substr string) int {
// 	return strings.Index(s, substr)
// }

// // IndexAny returns the index of the first instance of any Unicode code point from chars in s, or -1 if no Unicode code point from chars is present in s.
// //
// // Example
// //
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Println(strings.IndexAny("chicken", "aeiouy")) 	// 2
// // 		fmt.Println(strings.IndexAny("crwth", "aeiouy"))	// -1
// // 	}
// func IndexAny(s, chars string) int {
// 	return strings.IndexAny(s, chars)
// }

// // IndexByte returns the index of the first instance of c in s, or -1 if c is not present in s.
// //
// // Example
// //
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Println(strings.IndexByte("golang", 'g'))	// 0
// // 		fmt.Println(strings.IndexByte("gophers", 'h'))	// 3
// // 		fmt.Println(strings.IndexByte("golang", 'x'))	// -1
// // 	}
// // 	result:
// func IndexByte(s string, c byte) int {
// 	return strings.IndexByte(s, c)
// }

// // IndexFunc returns the index into s of the first Unicode code point satisfying f(c), or -1 if none do.
// //
// // Example
// //
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 		"unicode"
// // 	)
// //
// // 	func main() {
// // 		f := func(c rune) bool {
// // 			return unicode.Is(unicode.Han, c)
// // 		}
// // 		fmt.Println(strings.IndexFunc("Hello, 世界", f))	// 7
// // 		fmt.Println(strings.IndexFunc("Hello, world", f))	// -1
// // 	}
// func IndexFunc(s string, f func(rune) bool) int {
// 	return strings.IndexFunc(s, f)
// }

// // IndexRune returns the index of the first instance of the Unicode code point r, or -1 if rune is not present in s. If r is utf8.RuneError, it returns the first instance of any invalid UTF-8 byte sequence.
// //
// // Example
// //
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Println(strings.IndexRune("chicken", 'k'))	// 4
// // 		fmt.Println(strings.IndexRune("chicken", 'd'))	// -1
// // 	}
// func IndexRune(s string, r rune) int {
// 	return strings.IndexRune(s, r)
// }

// // Join concatenates the elements of its first argument to create a single string. The separator string sep is placed between elements in the resulting string.
// //
// // Example
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		s := []string{"foo", "bar", "baz"}
// // 		fmt.Println(strings.Join(s, ", ")) // foo, bar, baz
// // 	}
// func Join(elems []string, sep string) string {
// 	return strings.Join(elems, sep)
// }

// // LastIndex returns the index of the last instance of substr in s, or -1 if substr is not present in s.
// //
// // Example
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Println(strings.Index("go gopher", "go"))			// 0
// // 		fmt.Println(strings.LastIndex("go gopher", "go"))		// 3
// // 		fmt.Println(strings.LastIndex("go gopher", "rodent"))	// -1
// // 	}
// func LastIndex(s, substr string) int {
// 	return strings.LastIndex(s, substr)
// }

// // LastIndexAny returns the index of the last instance of any Unicode code point from chars in s, or -1 if no Unicode code point from chars is present in s.
// //
// // Example
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Println(strings.LastIndexAny("go gopher", "go"))		// 4
// // 		fmt.Println(strings.LastIndexAny("go gopher", "rodent"))	// 8
// // 		fmt.Println(strings.LastIndexAny("go gopher", "fail"))		// -1
// // 	}
// //
// func LastIndexAny(s, chars string) int {
// 	return strings.LastIndexAny(s, chars)
// }

// // LastIndexByte returns the index of the last instance of c in s, or -1 if c is not present in s.
// //
// // Example
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Println(strings.LastIndexByte("Hello, world", 'l')) // 10
// // 		fmt.Println(strings.LastIndexByte("Hello, world", 'o')) // 8
// // 		fmt.Println(strings.LastIndexByte("Hello, world", 'x')) // -1
// // 	}
// //
// func LastIndexByte(s string, c byte) int {
// 	return strings.LastIndexByte(s, c)
// }

// // LastIndexFunc returns the index into s of the last Unicode code point satisfying f(c), or -1 if none do.
// //
// // Example
// //
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 		"unicode"
// // 	)
// //
// // 	func main() {
// // 		fmt.Println(strings.LastIndexFunc("go 123", unicode.IsNumber))	// 5
// // 		fmt.Println(strings.LastIndexFunc("123 go", unicode.IsNumber))	// 2
// // 		fmt.Println(strings.LastIndexFunc("go", unicode.IsNumber))		// -1
// // 	}
// //
// func LastIndexFunc(s string, f func(rune) bool) int {
// 	return strings.LastIndexFunc(s, f)
// }

// // Split slices s into all substrings separated by sep and returns a slice of the substrings between those separators.
// //
// // If s does not contain sep and sep is not empty, Split returns a slice of length 1 whose only element is s.
// //
// // If sep is empty, Split splits after each UTF-8 sequence. If both s and sep are empty, Split returns an empty slice.
// //
// // It is equivalent to SplitN with a count of -1.
// //
// // Example
// //
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Printf("%q\n", strings.Split("a,b,c", ","))
// // 		// ["a" "b" "c"]
// // 		fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a "))
// // 		// ["" "man " "plan " "canal panama"]
// // 		fmt.Printf("%q\n", strings.Split(" xyz ", ""))
// // 		// [" " "x" "y" "z" " "]
// // 		fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins"))
// // 		// [""]
// // 	}
// func Split(s, sep string) []string {
// 	return strings.Split(s, sep)
// }

// //SplitN slices s into substrings separated by sep and returns a slice of the substrings between those separators.
// //
// //The count determines the number of substrings to return:
// //
// //	n > 0: at most n substrings; the last substring will be the unsplit remainder.
// //	n == 0: the result is nil (zero substrings)
// //	n < 0: all substrings
// //
// // Edge cases for s and sep (for example, empty strings) are handled as described in the documentation for Split.
// //
// // Example
// //
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Printf("%q\n", strings.SplitN("a,b,c", ",", 2)) // ["a" "b,c"]
// // 		z := strings.SplitN("a,b,c", ",", 0)
// // 		fmt.Printf("%q (nil = %v)\n", z, z == nil) //[] (nil = true)
// // 	}
// //
// func SplitN(s, sep string, n int) []string {
// 	return strings.SplitN(s, sep, n)
// }

// // SplitAfter slices s into all substrings after each instance of sep and returns a slice of those substrings.
// //
// // If s does not contain sep and sep is not empty, SplitAfter returns a slice of length 1 whose only element is s.
// //
// // If sep is empty, SplitAfter splits after each UTF-8 sequence. If both s and // sep are empty, SplitAfter returns an empty slice.
// //
// // It is equivalent to SplitAfterN with a count of -1.
// //
// // Example
// //
// // 	package main
// //
// // 	import (
// // 		"fmt"
// // 		"strings"
// // 	)
// //
// // 	func main() {
// // 		fmt.Printf("%q\n", strings.SplitAfter("a,b,c", ",")) // ["a," "b," "c"]
// // 	}
// //
// func SplitAfter(s, sep string) []string {
// 	return strings.SplitAfter(s, sep)
// }

// // SplitAfterN slices s into substrings after each instance of sep and returns a slice of those substrings.
// //
// // The count determines the number of substrings to return:
// //
// // 	n > 0: at most n substrings; the last substring will be the unsplit remainder.
// // 	n == 0: the result is nil (zero substrings)
// // 	n < 0: all substrings
// //
// // Edge cases for s and sep (for example, empty strings) are handled as described in the documentation for SplitAfter.
// //
// // Example
// //
// // package main
// //
// // import (
// // 	"fmt"
// // 	"strings"
// // )
// //
// // func main() {
// // 	fmt.Printf("%q\n", strings.SplitAfterN("a,b,c", ",", 2)) // ["a," "b,c"]
// // }
// //
// func SplitAfterN(s, sep string, n int) []string {
// 	return strings.SplitAfterN(s, sep, n)
// }

// // TrimLeft returns a slice of the string `s` with all leading Unicode code points contained in `cutset` removed.
// //
// // To remove a `prefix`, use `TrimPrefix` instead.
// //
// // Example:
// // 	fmt.Print(strings.TrimLeft("¡¡¡Hello, Gophers!!!", "!¡"))
// // 	out: Hello, Gophers!!!
// func TrimLeft(s, cutset string) string {
// 	return strings.TrimLeft(s, cutset)
// }

// // TrimLeftFunc returns a slice of the string `s` with all leading Unicode code points `c` satisfying `f(c)` removed.
// //
// // Example:
// // 	fmt.Print(strings.TrimLeftFunc("¡¡¡Hello, Gophers!!!", func(r rune) bool {return !unicode.IsLetter(r) && !unicode.IsNumber(r)}))
// // 	out: Hello, Gophers!!!
// func TrimLeftFunc(s string, f func(rune) bool) string {
// 	return strings.TrimLeftFunc(s, f)
// }

// // TrimPrefix returns `s` without the provided leading prefix string. If `s` doesn't start with `prefix`, `s` is returned unchanged.
// //
// // Example:
// // 	var s = "¡¡¡Hello, Gophers!!!"
// // 	s = strings.TrimPrefix(s, "¡¡¡Hello, ")
// // 	s = strings.TrimPrefix(s, "¡¡¡Howdy, ")
// // 	fmt.Print(s)
// // 	out: Gophers!!!
// func TrimPrefix(s, prefix string) string {
// 	return strings.TrimPrefix(s, prefix)
// }

// // TrimRight returns a slice of the string `s`, with all trailing Unicode code points contained in cutset removed.
// //
// // To remove a `suffix`, use `TrimSuffix` instead.
// //
// // Example:
// // 	fmt.Print(strings.TrimRight("¡¡¡Hello, Gophers!!!", "!¡"))
// // out: ¡¡¡Hello, Gophers
// func TrimRight(s, cutset string) string {
// 	return strings.TrimRight(s, cutset)
// }

// // TrimRightFunc returns a slice of the string `s` with all trailing Unicode code points `c` satisfying `f(c)` removed.
// //
// // Example:
// // 	fmt.Print(strings.TrimRightFunc("¡¡¡Hello, Gophers!!!", func(r rune) bool {return !unicode.IsLetter(r) && !unicode.IsNumber(r)}))
// // out: ¡¡¡Hello, Gophers
// func TrimRightFunc(s string, f func(rune) bool) string {
// 	return strings.TrimRightFunc(s, f)
// }

// // TrimSpace returns a slice of the string `s`, with all leading and trailing white space removed, as defined by Unicode.
// //
// // Example:
// // 	fmt.Println(strings.TrimSpace(" \t\n Hello, Gophers \n\t\r\n"))
// // 	out: Hello, Gophers
// func TrimSpace(s string) string {
// 	return strings.TrimSpace(s)
// }

// // TrimSuffix returns `s` without the provided trailing `suffix` string.
// // If `s` doesn't end with suffix, `s` is returned unchanged.
// //
// // Example:
// // 	var s = "¡¡¡Hello, Gophers!!!"
// // 	s = strings.TrimSuffix(s, ", Gophers!!!")
// // 	s = strings.TrimSuffix(s, ", Marmots!!!")
// // 	out: ¡¡¡Hello
// func TrimSuffix(s, suffix string) string {
// 	return strings.TrimSuffix(s, suffix)
// }

// // ---------------

// // ToUpper returns `s` with all Unicode letters mapped to their upper case.
// //
// // Example:
// // 	fmt.Println(strings.ToUpper("Gopher"))
// // 	out: GOPHER
// func ToUpper(s string) string {
// 	return strings.ToUpper(s)
// }

// // ToTitle returns a copy of the string `s` with all Unicode letters mapped to their Unicode title case.
// //
// // Example:
// // 	Compare this example to the Title example.
// // 	fmt.Println(strings.ToTitle("her royal highness"))
// // 	fmt.Println(strings.ToTitle("loud noises"))
// // 	fmt.Println(strings.ToTitle("хлеб"))
// // 	out:
// // 	HER ROYAL HIGHNESS
// // 	LOUD NOISES
// // 	ХЛЕБ
// func ToTitle(s string) string {
// 	return strings.ToTitle(s)
// }

// // ToLower returns s with all Unicode letters mapped to their lower case.
// //
// // Example:
// //	fmt.Println(strings.ToLower("Gopher"))
// // 	out: gopher
// func ToLower(s string) string {
// 	return strings.ToLower(s)
// }

// // Title returns a copy of the string `s` with all Unicode letters that begin words mapped to their Unicode title case.
// //
// // BUG(rsc): The rule Title uses for word boundaries does not handle Unicode punctuation properly.
// //
// // Example:
// // 	Compare this example to the ToTitle example.
// // 	fmt.Println(strings.Title("her royal highness"))
// // 	fmt.Println(strings.Title("loud noises"))
// // 	fmt.Println(strings.Title("хлеб"))
// // 	out:
// // 	Her Royal Highness
// // 	Loud Noises
// // 	Хлеб
// func Title(s string) string {
// 	return strings.Title(s)
// }

// // Map returns a copy of the string `s` with all its characters modified according to the `mapping` function. If `mapping` returns a negative value, the character is dropped from the string with no replacement.
// //
// // Example:
// // 	rot13 := func(r rune) rune {
// // 		switch {
// // 		case r >= 'A' && r <= 'Z':
// // 			return 'A' + (r-'A'+13)%26
// // 		case r >= 'a' && r <= 'z':
// // 			return 'a' + (r-'a'+13)%26
// // 		}
// // 		return r
// // 	}
// // 	fmt.Println(strings.Map(rot13, "'Twas brillig and the slithy gopher..."))
// // 	out:
// // 	'Gjnf oevyyvt naq gur fyvgul tbcure...
// func Map(mapping func(rune) rune, s string) string {
// 	return strings.Map(mapping, s)
// }

// // Repeat returns a new string consisting of count copies of the string `s`.
// //
// // It panics if `count` is negative or if the result of (`len(s) * count`) overflows.
// //
// // Example:
// // 	fmt.Println("ba" + strings.Repeat("na", 2))
// // 	out: banana
// func Repeat(s string, count int) string {
// 	return strings.Repeat(s, count)
// }

// // Replace returns a copy of the string `s` with the first `n` non-overlapping instances of `old` replaced by `new`. If `old` is empty, it matches at the beginning of the string and after each UTF-8 sequence, yielding up to `k+1` replacements for a `k-rune string. If `n < 0`, there is no limit on the number of replacements.
// //
// // Example:
// // 	fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))
// // 	fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1))
// // 	out:
// // 	oinky oinky oink
// // 	moo moo moo
// func Replace(s, old, new string, n int) string {
// 	return strings.Replace(s, old, new, n)
// }

// // ReplaceAll returns a copy of the string `s` with all non-overlapping instances of `old` replaced by `new`. If `old` is empty, it matches at the beginning of the string and after each UTF-8 sequence, yielding up to `k+1` replacements for a `k`-rune string.
// //
// // Example:
// // 	fmt.Println(strings.ReplaceAll("oink oink oink", "oink", "moo"))
// // 	out: moo moo moo
// func ReplaceAll(s, old, new string) string {
// 	return strings.ReplaceAll(s, old, new)
// }

// GetAbbrString return the abbreviation string 'xxx...' of `str`
//
// 	`maxlen`: maimium length of the abbreviation
// 	`conSymbole`: tailing symbol of the abbreviation
func GetAbbrString(str string, maxlen int, conSymbole string) string {
	// hc, ac := CountPlaceHolder(str)
	// lenStr := hc + ac
	lenStr := StringWidth(str)
	if lenStr <= maxlen {
		return str
	}
	if StringWidth(conSymbole) < 1 {
		conSymbole = "..."
	}
	return Truncate(str, maxlen, conSymbole)
}

// CountPlaceHolder return `nHan` and `nASCII`
//
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

// NumberBanner return numbers' string with length `width`
//
// Example:
// 	NumberBanner(11) return "01234567890"
func NumberBanner(width int) string {
	nl := []byte("0123456789")
	sb := strings.Builder{}
	for i := 0; i < width; i++ {
		c := nl[i%10]
		sb.Write([]byte{c})
	}
	return sb.String()
}

// Reverse reverse the string `s` based on `rune`
func Reverse(s string) string {
	return ReverseString(s)
}

// ReverseString reverses a string
func ReverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// ReverseByte reverse the string `s` based on `byte`
func ReverseByte(s string) string {
	sl := []byte(s)
	for i, j := 0, len(sl)-1; i < j; i, j = i+1, j-1 {
		sl[i], sl[j] = sl[j], sl[i]
	}
	return string(sl)
}

// GbkToUtf8 decodes GBK to UTF8
func GbkToUtf8(s []byte) ([]byte, error) {
	rd := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(rd)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// GbkToUtf8String decodes GBK to UTF8
func GbkToUtf8String(s string) (string, error) {
	bs := []byte(s)
	us, e := GbkToUtf8(bs)
	return string(us), e
}

// Utf8ToGbk encodes UTF8 to GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	rd := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(rd)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToGbkString encodes UTF8 to GBK
func Utf8ToGbkString(s string) (string, error) {
	bs := []byte(s)
	us, e := Utf8ToGbk(bs)
	return string(us), e
}

// Big5ToUtf8 decodes Big5 to UTF8
func Big5ToUtf8(s []byte) ([]byte, error) {
	rd := transform.NewReader(bytes.NewReader(s), traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(rd)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Big5ToUtf8String decodes Big5 to UTF8
func Big5ToUtf8String(s string) (string, error) {
	bs := []byte(s)
	us, e := Big5ToUtf8(bs)
	return string(us), e
}

// Utf8ToBig5 encodes UTF8 to Big5
func Utf8ToBig5(s []byte) ([]byte, error) {
	rd := transform.NewReader(bytes.NewReader(s), traditionalchinese.Big5.NewEncoder())
	d, e := ioutil.ReadAll(rd)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToBig5String encodes UTF8 to Big5
func Utf8ToBig5String(s string) (string, error) {
	bs := []byte(s)
	us, e := Utf8ToBig5(bs)
	return string(us), e
}

// IsEqualString compares string `a` and `b`
func IsEqualString(a, b string, ignoreCase bool) bool {
	if ignoreCase {
		return strings.EqualFold(a, b)
	}
	i := strings.Compare(a, b)
	if i == 0 {
		return true
	}
	return false
}

// TrimBOM trim the leading BOM character of a string
func TrimBOM(line string) string {
	l := len(line)
	if l >= 3 {
		if line[0] == 0xef && line[1] == 0xbb && line[2] == 0xbf {
			trimLine := line[3:]
			return trimLine
		}
	}
	return line
}

// TrimFrontEndSpaceLine trim the front and end empty line of `content` and return
func TrimFrontEndSpaceLine(content string) string {
	lines := strings.Split(content, "\n")
	fIdx := -1
	eIdx := -1
	for i := 0; i < len(lines); i++ {
		s := strings.TrimSpace(lines[i])
		if len(s) > 0 {
			fIdx = i
			break
		}
	}
	for i := len(lines) - 1; i >= 0; i-- {
		s := strings.TrimSpace(lines[i])
		if len(s) > 0 {
			eIdx = i + 1
			break
		}
	}
	if fIdx == -1 && eIdx == -1 {
		return content
	}
	lines = append([]string{}, lines[fIdx:eIdx]...)
	return strings.Join(lines, "\n")
}

// // type StringBuilder strings.Builder

// // NewStringBuilder will return `*strings.Builder`
// //
// // A Builder is used to efficiently build a string using Write methods. It minimizes memory copying. The zero value is ready to use. Do not copy a non-zero Builder.
// func NewStringBuilder() *strings.Builder {
// 	return new(strings.Builder)
// }

// // NewStringReader returns a new Reader reading from s. It is similar to bytes.NewBufferString but more efficient and read-only.
// func NewStringReader(s string) *strings.Reader {
// 	return strings.NewReader(s)
// }

// // NewBuffer creates and initializes a new Buffer using buf as its initial contents. The new Buffer takes ownership of buf, and the caller should not use buf after this call. NewBuffer is intended to prepare a Buffer to read existing data. It can also be used to set the initial size of the internal buffer for writing. To do that, buf should have the desired capacity but a length of zero.
// //
// // In most cases, new(Buffer) (or just declaring a Buffer variable) is sufficient to initialize a Buffer.
// func NewBuffer(buf []byte) *bytes.Buffer {
// 	return bytes.NewBuffer(buf)
// }

// // NewReader returns a new Reader whose buffer has the default size.
// func NewBufioReader(s string) *bufio.Reader {
// 	// return bufio.NewReader(NewReader(s))
// 	return bufio.NewReader(NewBuffer([]byte(s)))
// }

// The following is adopted from github.com/mattn/go-runewidth

// FillLeft return string filled in left by spaces in w cells
func FillLeft(s string, w int) string {
	return runewidth.FillLeft(s, w)
}

// FillRight return string filled in left by spaces in w cells
func FillRight(s string, w int) string {
	return runewidth.FillRight(s, w)
}

// FillLeftRight return string filled in left and right by spaces in cells with width w (aligns center)
func FillLeftRight(s string, w int) string {
	ns := StringWidth(s)
	if ns <= w {
		return s
	}
	nr := (w - ns) / 2
	nl := w - ns - nr
	lsp := make([]byte, nr)
	for i := range lsp {
		lsp[i] = ' '
	}
	rsp := make([]byte, nl)
	for i := range rsp {
		rsp[i] = ' '
	}
	return string(lsp) + s + string(rsp)
}

// StringWidth will return width as you can see (the numbers of placeholders on terminal)
func StringWidth(s string) int {
	return runewidth.StringWidth(s)
}

// RuneWidth returns the number of cells in r. See http://www.unicode.org/reports/tr11/
func RuneWidth(r rune) int {
	return runewidth.RuneWidth(r)
}

// Truncate return string truncated with w cells
func Truncate(s string, w int, tail string) string {
	return runewidth.Truncate(s, w, tail)
}

// Wrap return string wrapped with w cells
func Wrap(s string, w int) string {
	return runewidth.Wrap(s, w)
}

// WrapToSlice return string slice wrapped with w cells
func WrapToSlice(s string, w int) []string {
	if StringWidth(s) <= w {
		return []string{s}
	}
	return strings.Split(Wrap(s, w), "\n")
}

// Spaces return a string with lenth w of spaces
func Spaces(w int) string {
	if w <= 0 {
		return ""
	}
	return strings.Repeat(" ", w)
}

// CheckIndex will check index idx whether is in range of slice. If not, return error
func CheckIndex(slice interface{}, idx int, varName string) error {
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("CheckIndex: expected slice type, found %q, variable name is %q", v.Kind().String(), varName)
	}
	count := v.Len()
	if idx < 0 || idx > count-1 {
		return fmt.Errorf("CheckIndex: slice range [%d, %d), idx is %d, variable name is %q", 0, count, idx, varName)
	}
	return nil
}

// CheckIndexInString will check index idx whether is in range of string. If not, return error
func CheckIndexInString(s string, idx int, varName string) error {
	var ns = len(s)
	if idx < 0 || idx > ns-1 {
		return fmt.Errorf("CheckIndexInString: bounds out of range, lenth is %d, index: %d, variable name is %q", ns, idx, varName)
	}
	return nil
}

// PaddingString add pad-prefix in every line of string
func PaddingString(s string, pad string) string {
	if !strings.Contains(s, "\n") {
		return pad + s
	}
	r := []rune(s)
	sb := new(strings.Builder)
	sb.WriteString(pad)
	for _, v := range r {
		if v == '\n' {
			// sb.WriteString("\n")
			sb.WriteRune('\n')
			sb.WriteString(pad)
		} else {
			sb.WriteString(string(v))
		}
	}
	return sb.String()
}

// PaddingBytes add pad-prefix in every line('\n') of []byte
func PaddingBytes(bytes []byte, pad string) []byte {
	b := make([]byte, len(bytes))
	b = append(b, pad...)
	for _, v := range bytes {
		b = append(b, v)
		if v == '\n' {
			b = append(b, pad...)
		}
	}
	return b
}

// StringWithWidth will return a constant width string according to align. If width of value as you see is greater than width, then return value
func StringWithWidth(align Align, value string, width int) string {
	if StringWidth(value) >= width || width <= 0 {
		return value
	}
	var r string
	switch align {
	case AlignLeft:
		r = FillRight(value, width)
	case AlignRight:
		r = FillLeft(value, width)
	case AlignCenter:
		r = FillLeftRight(value, width)
	default:
		r = value
	}
	return r
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var reANSI = regexp.MustCompile(ansi)

func StripANSI(str string) string {
	return reANSI.ReplaceAllString(str, "")
}

// // ForEachString higher order function that processes each line of text by callback function.
// // The last non-empty line of input will be processed even if it has no newline.
// // 	`br` : read from `br` reader
// // 	`callback` : the function used to treatment the each line from `br`
// //
// // modify from "github.com/liuzl/goutil"
// func ForEachString(br *bufio.Reader, callback func(string) error) error {
// 	stop := false
// 	for {
// 		if stop {
// 			break
// 		}
// 		line, err := br.ReadString('\n')
// 		if err == io.EOF {
// 			stop = true
// 		} else if err != nil {
// 			return err
// 		}
// 		line = TrimSuffix(line, "\n")
// 		if line == "" {
// 			if !stop {
// 				if err = callback(line); err != nil {
// 					return err
// 				}
// 			}
// 			continue
// 		}
// 		if err = callback(line); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
