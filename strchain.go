package paw

import "strings"

// StrChain contains all tools which can be chained.
type StrChain struct {
	v string
	// TBError error
}

// NewStrChain return a instance of `StrChain` and return `*StrChain`
func (t *StrChain) NewStrChain(s string) *StrChain {
	t = &StrChain{s}
	return t
}

// String return `StrChain.sCollection.s`
func (t *StrChain) String() string {
	return t.v
}

// SetText set `StrChain.s` to `txt`
func (t *StrChain) SetText(txt string) *StrChain {
	t.v = txt
	return t
}

// GetText return `StrChain.sCollection.s`
func (t *StrChain) GetText() string {
	return t.v
}

// Len will return the lenth of t.v (would be the sizes of []bytes)
func (t *StrChain) Len() int {
	return len(t.v)
}

// Bytes will convert the string t.v to []byte
//
// Example:
// 	b := StrChain{"ABC€"}
// 	fmt.Println(b.Bytes()) // [65 66 67 226 130 172]
func (t *StrChain) Bytes() []byte {
	return []byte(t.v)
}

// Runes will convert the string t.v to []rune
//
// Example:
// 	r := StrChain{"ABC€"}
// 	fmt.Println(r.Runes())        	// [65 66 67 8364]
// 	fmt.Printf("%U\n", r.Rune()) 	// [U+0041 U+0042 U+0043 U+20AC]
func (t *StrChain) Runes() []rune {
	return []rune(t.v)
}

// GetAbbrString get a abbreviation of `StrChain.sCollection.s` and save to `StrChain.sCollection.s`
//
// 	`maxlen`: maimium length of the abbreviation
// 	`conSymbole`: tailing symbol of the abbreviation
func (t *StrChain) GetAbbrString(maxlen int, contSymbol string) *StrChain {
	t.v = GetAbbrString(t.v, maxlen, contSymbol)
	return t
}

// CountPlaceHolder return `nHan` and `nASCII`
//
// 	`nHan`: number of occupied space in screen for han-character
// 	`nASCII`: number of occupied space in screen for ASCII-character
func (t *StrChain) CountPlaceHolder() (nHan int, nASCII int) {
	return CountPlaceHolder(t.v)
}

// HasChineseChar return true for that `str` include chinese character
//
// Example:
// 	HasChineseChar("abc 中文") return true
// 	HasChineseChar("abccefgh") return false
func (t *StrChain) HasChineseChar() bool {
	return HasChineseChar(t.v)
}

// NumberBanner return numbers' string with length of `StrChain.sCollection.s`
//
// Example:
// 	StrChain.sCollection.s = "Text中文 Collection"
// 	nh, na := CountPlaceHolder（"Text中文 Collection"）
// 	--> nh=4, na=15 --> length = nh + na = 19
// 	NumberBanner() return "12345678901"
func (t *StrChain) NumberBanner() *StrChain {
	h, a := t.CountPlaceHolder()
	t.v = NumberBanner(h + a)
	return t
}

// Reverse packs `Reverse(s string)` based on `rune`
// 	set `StrChain.S` to the result
func (t *StrChain) Reverse() *StrChain {
	t.v = Reverse(t.v)
	return t
}

// ReverseByte packs `ReverseByte(s string)` based on `byte`
// 	set `StrChain.S` to the result
func (t *StrChain) ReverseByte() *StrChain {
	t.v = ReverseByte(t.v)
	return t
}

// HasPrefix return `HasPrefix(t.v, prefix)`
func (t *StrChain) HasPrefix(prefix string) bool {
	return strings.HasPrefix(t.v, prefix)
}

// HasSuffix return `HasSuffix(t.v, Suffix)`
func (t *StrChain) HasSuffix(suffix string) bool {
	return strings.HasSuffix(t.v, suffix)
}

// Contains return `Contains(t.v, substr)`
func (t *StrChain) Contains(substr string) bool {
	return strings.Contains(t.v, substr)
}

// ContainsAny return `ContainsAny(t.v, chars)`
func (t *StrChain) ContainsAny(chars string) bool {
	return strings.ContainsAny(t.v, chars)
}

// Fields return Fields(t.v)
func (t *StrChain) Fields() []string {
	return strings.Fields(t.v)
}

// FieldsFunc return FieldsFunc(t.v, f)
func (t *StrChain) FieldsFunc(f func(rune) bool) []string {
	return strings.FieldsFunc(t.v, f)
}

// ContainsAny return ContainsRune(t.v, r) bool
func (t *StrChain) ContainsRune(r rune) bool {
	return strings.ContainsRune(t.v, r)
}

// EqualFold return EqualFold(t.v,, t) bool
func (t *StrChain) EqualFold(s string) bool {
	return strings.EqualFold(t.v, s)
}

// Index return Index(t.v, substr) int
func (t *StrChain) Index(substr string) int {
	return strings.Index(t.v, substr)
}

// IndexAny return IndexAny(t.v, chars) int
func (t *StrChain) IndexAny(chars string) int {
	return strings.IndexAny(t.v, chars)
}

// IndexByte return IndexByte(t.v, c) int
func (t *StrChain) IndexByte(c byte) int {
	return strings.IndexByte(t.v, c)
}

// IndexFunc return IndexFunc(t.v, f) int
func (t *StrChain) IndexFunc(f func(rune) bool) int {
	return strings.IndexFunc(t.v, f)
}

// IndexRune return IndexRune(t.v, r) int
func (t *StrChain) IndexRune(r rune) int {
	return strings.IndexRune(t.v, r)
}

// LastIndex return LastIndex(t.v, substr) int
func (t *StrChain) LastIndex(substr string) int {
	return strings.LastIndex(t.v, substr)
}

// LastIndexAny return LastIndexAny(t.v, chars) int
func (t *StrChain) LastIndexAny(chars string) int {
	return strings.LastIndexAny(t.v, chars)
}

// LastIndexByte return LastIndexByte(t.v, c) int
func (t *StrChain) LastIndexByte(c byte) int {
	return strings.LastIndexByte(t.v, c)
}

// LastIndexFunc return LastIndexFunc(t.v, f) int
func (t *StrChain) LastIndexFunc(f func(rune) bool) int {
	return strings.LastIndexFunc(t.v, f)
}

// Split return Split(t.v, sep) []string
func (t *StrChain) Split(sep string) []string {
	return strings.Split(t.v, sep)
}

// Split return Split(t.v, sep) []string
func (t *StrChain) SplitN(sep string, n int) []string {
	return strings.SplitN(t.v, sep, n)
}

// Split return Split(t.v, sep) []string
func (t *StrChain) SplitAfter(sep string) []string {
	return strings.SplitAfter(t.v, sep)
}

// Split return Split(t.v, sep) []string
func (t *StrChain) SplitAfterN(sep string, n int) []string {
	return strings.SplitAfterN(t.v, sep, n)
}

// Trim packs `Trim(s, cutset)`
// 	set `StrChain.S` to the result
func (t *StrChain) Trim(cutset string) *StrChain {
	t.v = strings.Trim(t.v, cutset)
	return t
}

// TrimFunc packs `TrimFunc(s string, f func(rune) bool)`
// 	set `StrChain.S` to the result
func (t *StrChain) TrimFunc(f func(rune) bool) *StrChain {
	t.v = strings.TrimFunc(t.v, f)
	return t
}

// TrimLeft packs `TrimLeft(s, cutset string)`
// 	set `StrChain.S` to the result
func (t *StrChain) TrimLeft(cutset string) *StrChain {
	t.v = strings.TrimLeft(t.v, cutset)
	return t
}

// TrimLeftFunc packs `TrimLeftFunc(s string, f func(rune) bool)`
// 	set `StrChain.S` to the result
func (t *StrChain) TrimLeftFunc(f func(rune) bool) *StrChain {
	t.v = strings.TrimLeftFunc(t.v, f)
	return t
}

// TrimPrefix packs `TrimPrefix(s, prefix string)`
// 	set `StrChain.S` to the result
func (t *StrChain) TrimPrefix(prefix string) *StrChain {
	t.v = strings.TrimPrefix(t.v, prefix)
	return t
}

// TrimRight packs `TrimRight(s, cutset string)`
// 	set `StrChain.S` to the result
func (t *StrChain) TrimRight(cutset string) *StrChain {
	t.v = strings.TrimRight(t.v, cutset)
	return t
}

// TrimRightFunc packs `TrimRightFunc(s string, f func(rune) bool)`
// 	set `StrChain.S` to the result
func (t *StrChain) TrimRightFunc(f func(rune) bool) *StrChain {
	t.v = strings.TrimRightFunc(t.v, f)
	return t
}

// TrimSpace packs `TrimSpace(s string)`
// 	set `StrChain.S` to the result
func (t *StrChain) TrimSpace() *StrChain {
	t.v = strings.TrimSpace(t.v)
	return t
}

// TrimSuffix packs `TrimSuffix(s, suffix string)`
// 	set `StrChain.S` to the result
func (t *StrChain) TrimSuffix(suffix string) *StrChain {
	t.v = strings.TrimSuffix(t.v, suffix)
	return t
}

// ToUpper packs `ToUpper(s string)`
// 	set `StrChain.S` to the result
func (t *StrChain) ToUpper() *StrChain {
	t.v = strings.ToUpper(t.v)
	return t
}

// ToTitle packs `ToTitle(s string)`
// 	set `StrChain.S` to the result
func (t *StrChain) ToTitle() *StrChain {
	t.v = strings.ToTitle(t.v)
	return t
}

// ToLower packs ` ToLower(s string)`
// 	set `StrChain.S` to the result
func (t *StrChain) ToLower() *StrChain {
	t.v = strings.ToLower(t.v)
	return t
}

// Title returns packs `Title(s string)`
// 	set `StrChain.S` to the result
func (t *StrChain) Title() *StrChain {
	t.v = strings.Title(t.v)
	return t
}

// Map packs `Map(mapping func(rune) rune, s string)`
// 	set `StrChain.S` to the result
func (t *StrChain) Map(mapping func(rune) rune) *StrChain {
	t.v = strings.Map(mapping, t.v)
	return t
}

// Repeat packs `Repeat(s string, count int)`
// 	set `StrChain.S` to the result
func (t *StrChain) Repeat(count int) *StrChain {
	t.v = strings.Repeat(t.v, count)
	return t
}

// Replace packs `Replace(s, old, new string, n int)`
// 	set `StrChain.S` to the result
func (t *StrChain) Replace(old, new string, n int) *StrChain {
	t.v = strings.Replace(t.v, old, new, n)
	return t
}

// ReplaceAll packs `ReplaceAll(s, old, new string)`
// 	set `StrChain.S` to the result
func (t *StrChain) ReplaceAll(old, new string) *StrChain {
	t.v = strings.ReplaceAll(t.v, old, new)
	return t
}

// GbkToUtf8String packs `GbkToUtf8String(s string)`
func (t *StrChain) GbkToUtf8String() (*StrChain, error) {
	s, e := GbkToUtf8String(t.v)
	if e != nil {
		t.v = ""
		return t, e
	}
	t.v = s
	return t, nil
}

// Utf8ToGbkString packs `Utf8ToGbkString(s string)`
func (t *StrChain) Utf8ToGbkString() (*StrChain, error) {
	s, e := Utf8ToGbkString(t.v)
	if e != nil {
		t.v = ""
		return t, e
	}
	t.v = s
	return t, nil
}

// Big5ToUtf8String packs `Big5ToUtf8String(s string)`
func (t *StrChain) Big5ToUtf8String() (*StrChain, error) {
	s, e := Big5ToUtf8String(t.v)
	if e != nil {
		t.v = ""
		return t, e
	}
	t.v = s
	return t, nil
}

// Utf8ToBig5String packs `Utf8ToBig5String(s string)`
func (t *StrChain) Utf8ToBig5String() (*StrChain, error) {
	s, e := Utf8ToBig5String(t.v)
	if e != nil {
		t.v = ""
		return t, e
	}
	t.v = s
	return t, nil
}

// IsEqualString packs `IsEqualString(a, b string, ignoreCase bool) bool`
func (t *StrChain) IsEqualString(b string, ignoreCase bool) bool {
	return IsEqualString(t.v, b, ignoreCase)
}

// TrimBOM packs `TrimBOM(line string) string`
func (t *StrChain) TrimBOM() *StrChain {
	t.v = TrimBOM(t.v)
	return t
}

// TrimFrontEndSpaceLine packs `TrimFrontEndSpaceLine(content string) string`
func (t *StrChain) TrimFrontEndSpaceLine() *StrChain {
	t.v = TrimFrontEndSpaceLine(t.v)
	return t
}

// The following is adopted from github.com/mattn/go-runewidth

// FillLeft return string filled in left by spaces in w cells
func (t *StrChain) FillLeft(w int) *StrChain {
	t.v = FillLeft(t.v, w)
	return t
}

// FillRight return string filled in left by spaces in w cells
func (t *StrChain) FillRight(w int) *StrChain {
	t.v = FillRight(t.v, w)
	return t
}

// StringWidth will return width as you can see (the numbers of placeholders on terminal)
func (t *StrChain) StringWidth() int {
	return StringWidth(t.v)
}

// // RuneWidth returns the number of cells in r. See http://www.unicode.org/reports/tr11/
// func (t *StrChain) RuneWidth(r rune) int {
// 	return RuneWidth(r)
// }

// Truncate return string truncated with w cells
func (t *StrChain) Truncate(w int, tail string) *StrChain {
	t.v = Truncate(t.v, w, tail)
	return t
}

// Wrap return string wrapped with w cells
func (t *StrChain) Wrap(w int) *StrChain {
	t.v = Wrap(t.v, w)
	return t
}
