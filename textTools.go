package paw

import (
	"bytes"
	"io/ioutil"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"

	"golang.org/x/text/transform"
	// "github.com/gookit/color"
)

// HasPrefix return `strings.HasPrefix(str, prefix)`
func HasPrefix(str string, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// HasSuffix return `strings.HasSuffix(str, Suffix)`
func HasSuffix(str string, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

// Contains return `strings.Contains(str, substr)`
func Contains(str string, substr string) bool {
	return strings.Contains(str, substr)
}

// Trim returns a slice of the string `s` with all leading and trailing Unicode code points contained in `cutset` removed.
func Trim(s, cutset string) string {
	return strings.Trim(s, cutset)
}

// TrimFunc returns a slice of the string `s` with all leading and trailing Unicode code points `c` satisfying `f(c)` removed.
//
// Example:
// 	fmt.Print(strings.TrimFunc("¡¡¡Hello, Gophers!!!", func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r)}))
// 	out: Hello, Gophers
func TrimFunc(s string, f func(rune) bool) string {
	return strings.TrimFunc(s, f)
}

// TrimLeft returns a slice of the string `s` with all leading Unicode code points contained in `cutset` removed.
//
// To remove a `prefix`, use `TrimPrefix` instead.
//
// Example:
// 	fmt.Print(strings.TrimLeft("¡¡¡Hello, Gophers!!!", "!¡"))
// 	out: Hello, Gophers!!!
func TrimLeft(s, cutset string) string {
	return strings.TrimLeft(s, cutset)
}

// TrimLeftFunc returns a slice of the string `s` with all leading Unicode code points `c` satisfying `f(c)` removed.
//
// Example:
// 	fmt.Print(strings.TrimLeftFunc("¡¡¡Hello, Gophers!!!", func(r rune) bool {return !unicode.IsLetter(r) && !unicode.IsNumber(r)}))
// 	out: Hello, Gophers!!!
func TrimLeftFunc(s string, f func(rune) bool) string {
	return strings.TrimLeftFunc(s, f)
}

// TrimPrefix returns `s` without the provided leading prefix string. If `s` doesn't start with `prefix`, `s` is returned unchanged.
//
// Example:
// 	var s = "¡¡¡Hello, Gophers!!!"
// 	s = strings.TrimPrefix(s, "¡¡¡Hello, ")
// 	s = strings.TrimPrefix(s, "¡¡¡Howdy, ")
// 	fmt.Print(s)
// 	out: Gophers!!!
func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

// TrimRight returns a slice of the string `s`, with all trailing Unicode code points contained in cutset removed.
//
// To remove a `suffix`, use `TrimSuffix` instead.
//
// Example:
// 	fmt.Print(strings.TrimRight("¡¡¡Hello, Gophers!!!", "!¡"))
// out: ¡¡¡Hello, Gophers
func TrimRight(s, cutset string) string {
	return strings.TrimRight(s, cutset)
}

// TrimRightFunc returns a slice of the string `s` with all trailing Unicode code points `c` satisfying `f(c)` removed.
//
// Example:
// 	fmt.Print(strings.TrimRightFunc("¡¡¡Hello, Gophers!!!", func(r rune) bool {return !unicode.IsLetter(r) && !unicode.IsNumber(r)}))
// out: ¡¡¡Hello, Gophers
func TrimRightFunc(s string, f func(rune) bool) string {
	return strings.TrimRightFunc(s, f)
}

// TrimSpace returns a slice of the string `s`, with all leading and trailing white space removed, as defined by Unicode.
//
// Example:
// 	fmt.Println(strings.TrimSpace(" \t\n Hello, Gophers \n\t\r\n"))
// 	out: Hello, Gophers
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// TrimSuffix returns `s` without the provided trailing `suffix` string.
// If `s` doesn't end with suffix, `s` is returned unchanged.
//
// Example:
// 	var s = "¡¡¡Hello, Gophers!!!"
// 	s = strings.TrimSuffix(s, ", Gophers!!!")
// 	s = strings.TrimSuffix(s, ", Marmots!!!")
// 	out: ¡¡¡Hello
func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

// ---------------

// ToUpper returns `s` with all Unicode letters mapped to their upper case.
//
// Example:
// 	fmt.Println(strings.ToUpper("Gopher"))
// 	out: GOPHER
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToTitle returns a copy of the string `s` with all Unicode letters mapped to their Unicode title case.
//
// Example:
// 	Compare this example to the Title example.
// 	fmt.Println(strings.ToTitle("her royal highness"))
// 	fmt.Println(strings.ToTitle("loud noises"))
// 	fmt.Println(strings.ToTitle("хлеб"))
// 	out:
// 	HER ROYAL HIGHNESS
// 	LOUD NOISES
// 	ХЛЕБ
func ToTitle(s string) string {
	return strings.ToUpper(s)
}

// ToLower returns s with all Unicode letters mapped to their lower case.
//
// Example:
//	fmt.Println(strings.ToLower("Gopher"))
// 	out: gopher
func ToLower(s string) string {
	return strings.ToLower(s)
}

// Title returns a copy of the string `s` with all Unicode letters that begin words mapped to their Unicode title case.
//
// BUG(rsc): The rule Title uses for word boundaries does not handle Unicode punctuation properly.
//
// Example:
// 	Compare this example to the ToTitle example.
// 	fmt.Println(strings.Title("her royal highness"))
// 	fmt.Println(strings.Title("loud noises"))
// 	fmt.Println(strings.Title("хлеб"))
// 	out:
// 	Her Royal Highness
// 	Loud Noises
// 	Хлеб
func Title(s string) string {
	return strings.Title(s)
}

// Map returns a copy of the string `s` with all its characters modified according to the `mapping` function. If `mapping` returns a negative value, the character is dropped from the string with no replacement.
//
// Example:
// 	rot13 := func(r rune) rune {
// 		switch {
// 		case r >= 'A' && r <= 'Z':
// 			return 'A' + (r-'A'+13)%26
// 		case r >= 'a' && r <= 'z':
// 			return 'a' + (r-'a'+13)%26
// 		}
// 		return r
// 	}
// 	fmt.Println(strings.Map(rot13, "'Twas brillig and the slithy gopher..."))
// 	out:
// 	'Gjnf oevyyvt naq gur fyvgul tbcure...
func Map(mapping func(rune) rune, s string) string {
	return strings.Map(mapping, s)
}

// Repeat returns a new string consisting of count copies of the string `s`.
//
// It panics if `count` is negative or if the result of (`len(s) * count`) overflows.
//
// Example:
// 	fmt.Println("ba" + strings.Repeat("na", 2))
// 	out: banana
func Repeat(s string, count int) string {
	return strings.Repeat(s, count)
}

// Replace returns a copy of the string `s` with the first `n` non-overlapping instances of `old` replaced by `new`. If `old` is empty, it matches at the beginning of the string and after each UTF-8 sequence, yielding up to `k+1` replacements for a `k-rune string. If `n < 0`, there is no limit on the number of replacements.
//
// Example:
// 	fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))
// 	fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1))
// 	out:
// 	oinky oinky oink
// 	moo moo moo
func Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

// ReplaceAll returns a copy of the string `s` with all non-overlapping instances of `old` replaced by `new`. If `old` is empty, it matches at the beginning of the string and after each UTF-8 sequence, yielding up to `k+1` replacements for a `k`-rune string.
//
// Example:
// 	fmt.Println(strings.ReplaceAll("oink oink oink", "oink", "moo"))
// 	out: moo moo moo
func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// GetAbbrString return the abbreviation string 'xxx...' of `str`
//
// 	`maxlen`: maimium length of the abbreviation
// 	`conSymbole`: tailing symbol of the abbreviation
func GetAbbrString(str string, maxlen int, conSymbole string) string {
	hc, ac := CountPlaceHolder(str)
	lenStr := hc + ac
	if lenStr <= maxlen {
		return str
	}
	if len(conSymbole) < 1 {
		conSymbole = "..."
	}
	limit := maxlen - len(conSymbole)
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
	str = sb.String() + conSymbole
	return str
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

// Reverse reverse the string `s`
func Reverse(s string) string {
	sl := []rune(s)
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
		s := TrimSpace(lines[i])
		if len(s) > 0 {
			fIdx = i
			break
		}
	}
	for i := len(lines) - 1; i >= 0; i-- {
		s := TrimSpace(lines[i])
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
