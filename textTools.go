package paw

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shyang107/paw/runewidth"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"

	"golang.org/x/text/transform"
	// "github.com/gookit/color"
)

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
	d, e := io.ReadAll(rd)
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
	d, e := io.ReadAll(rd)
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
	d, e := io.ReadAll(rd)
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
	d, e := io.ReadAll(rd)
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

// -----------------------------------------------------------
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
	for i, v := range r {
		if v == '\n' {
			// sb.WriteString("\n")
			sb.WriteRune('\n')
			if i < len(r)-1 {
				sb.WriteString(pad)
			}
		} else {
			sb.WriteRune(v)
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

// AlignWithWidth will return a constant width string according to align. If width of value as you see is greater than width, then return value
func AlignWithWidth(align Align, value string, width int) string {
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
