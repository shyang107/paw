package paw

import (
	"github.com/mattn/go-runewidth"
)

// StrChain chains some function about string
type StrChain string

func (s StrChain) String() string {
	return string(s)
}

// GetAbbrString get a abbreviation of StrChain
//
// 	`maxlen`: maimium length of the abbreviation
// 	`conSymbole`: tailing symbol of the abbreviation
func (s *StrChain) GetAbbrString(maxlen int, contSymbol string) *StrChain {
	*s = StrChain(GetAbbrString(string(*s), maxlen, contSymbol))
	return s
}

// RuneStringWidth will return width as you can see
func (s *StrChain) RuneStringWidth() int {
	return runewidth.StringWidth(string(*s))
}

// CountPlaceHolder return `nHan` and `nASCII`
//
// 	`nHan`: number of occupied space in screen for han-character
// 	`nASCII`: number of occupied space in screen for ASCII-character
func (s *StrChain) CountPlaceHolder() (nHan int, nASCII int) {
	return CountPlaceHolder(string(*s))
}

// HasChineseChar return true for that `str` include chinese character
//
// Example:
// 	HasChineseChar("abc 中文") return true
// 	HasChineseChar("abccefgh") return false
func (s *StrChain) HasChineseChar() bool {
	return HasChineseChar(string(*s))
}

// NumberBanner return numbers' string with length of `TextBuilder.TextCollection.Text`
//
// Example:
// 	TextBuilder.TextCollection.Text = "Text中文 Collection"
// 	nh, na := CountPlaceHolder（"Text中文 Collection"）
// 	--> nh=4, na=15 --> length = nh + na = 19
// 	NumberBanner() return "12345678901"
func (s *StrChain) NumberBanner() *StrChain {
	// h, a := s.CountPlaceHolder()
	// *s = StrChain(NumberBanner(h + a))
	*s = StrChain(NumberBanner(s.RuneStringWidth()))
	return s
}

// Reverse packs `Reverse(s string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) Reverse() *StrChain {
	*s = StrChain(Reverse(string(*s)))
	return s
}

// HasPrefix return `HasPrefix(string(*s), prefix)`
func (s *StrChain) HasPrefix(prefix string) bool {
	return HasPrefix(string(*s), prefix)
}

// HasSuffix return `HasSuffix(string(*s), Suffix)`
func (s *StrChain) HasSuffix(suffix string) bool {
	return HasSuffix(string(*s), suffix)
}

// Contains return `Contains(string(*s), substr)`
func (s *StrChain) Contains(substr string) bool {
	return Contains(string(*s), substr)
}

// ContainsAny return `ContainsAny(string(*s), chars)`
func (s *StrChain) ContainsAny(chars string) bool {
	return ContainsAny(string(*s), chars)
}

// Fields return Fields(string(*s))
func (s *StrChain) Fields() []string {
	return Fields(string(*s))
}

// FieldsFunc return FieldsFunc(string(*s), f)
func (s *StrChain) FieldsFunc(f func(rune) bool) []string {
	return FieldsFunc(string(*s), f)
}

// ContainsAny return ContainsRune(string(*s), r) bool
func (s *StrChain) ContainsRune(r rune) bool {
	return ContainsRune(string(*s), r)
}

// EqualFold return EqualFold(string(*s),, t) bool
func (s *StrChain) EqualFold(t string) bool {
	return EqualFold(string(*s), t)
}

// Index return Index(string(*s), substr) int
func (s *StrChain) Index(substr string) int {
	return Index(string(*s), substr)
}

// IndexAny return IndexAny(string(*s), chars) int
func (s *StrChain) IndexAny(chars string) int {
	return IndexAny(string(*s), chars)
}

// IndexByte return IndexByte(string(*s), c) int
func (s *StrChain) IndexByte(c byte) int {
	return IndexByte(string(*s), c)
}

// IndexFunc return IndexFunc(string(*s), f) int
func (s *StrChain) IndexFunc(f func(rune) bool) int {
	return IndexFunc(string(*s), f)
}

// IndexRune return IndexRune(string(*s), r) int
func (s *StrChain) IndexRune(r rune) int {
	return IndexRune(string(*s), r)
}

// LastIndex return LastIndex(string(*s), substr) int
func (s *StrChain) LastIndex(substr string) int {
	return LastIndex(string(*s), substr)
}

// LastIndexAny return LastIndexAny(string(*s), chars) int
func (s *StrChain) LastIndexAny(chars string) int {
	return LastIndexAny(string(*s), chars)
}

// LastIndexByte return LastIndexByte(string(*s), c) int
func (s *StrChain) LastIndexByte(c byte) int {
	return LastIndexByte(string(*s), c)
}

// LastIndexFunc return LastIndexFunc(string(*s), f) int
func (s *StrChain) LastIndexFunc(f func(rune) bool) int {
	return LastIndexFunc(string(*s), f)
}

// Split return Split(string(*s), sep) []string
func (s *StrChain) Split(sep string) []string {
	return Split(string(*s), sep)
}

// Split return Split(string(*s), sep) []string
func (s *StrChain) SplitN(sep string, n int) []string {
	return SplitN(string(*s), sep, n)
}

// Split return Split(string(*s), sep) []string
func (s *StrChain) SplitAfter(sep string) []string {
	return SplitAfter(string(*s), sep)
}

// Split return Split(string(*s), sep) []string
func (s *StrChain) SplitAfterN(sep string, n int) []string {
	return SplitAfterN(string(*s), sep, n)
}

// Trim packs `Trim(s, cutset)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) Trim(cutset string) *StrChain {
	*s = StrChain(Trim(string(*s), cutset))
	return s
}

// TrimFunc packs `TrimFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) TrimFunc(f func(rune) bool) *StrChain {
	*s = StrChain(TrimFunc(string(*s), f))
	return s
}

// TrimLeft packs `TrimLeft(s, cutset string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) TrimLeft(cutset string) *StrChain {
	*s = StrChain(TrimLeft(string(*s), cutset))
	return s
}

// TrimLeftFunc packs `TrimLeftFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) TrimLeftFunc(f func(rune) bool) *StrChain {
	*s = StrChain(TrimLeftFunc(string(*s), f))
	return s
}

// TrimPrefix packs `TrimPrefix(s, prefix string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) TrimPrefix(prefix string) *StrChain {
	*s = StrChain(TrimPrefix(string(*s), prefix))
	return s
}

// TrimRight packs `TrimRight(s, cutset string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) TrimRight(cutset string) *StrChain {
	*s = StrChain(TrimRight(string(*s), cutset))
	return s
}

// TrimRightFunc packs `TrimRightFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) TrimRightFunc(f func(rune) bool) *StrChain {
	*s = StrChain(TrimRightFunc(string(*s), f))
	return s
}

// TrimSpace packs `TrimSpace(s string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) TrimSpace() *StrChain {
	*s = StrChain(TrimSpace(string(*s)))
	return s
}

// TrimSuffix packs `TrimSuffix(s, suffix string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) TrimSuffix(suffix string) *StrChain {
	*s = StrChain(TrimSuffix(string(*s), suffix))
	return s
}

// ToUpper packs `ToUpper(s string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) ToUpper() *StrChain {
	*s = StrChain(ToUpper(string(*s)))
	return s
}

// ToTitle packs `ToTitle(s string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) ToTitle() *StrChain {
	*s = StrChain(ToUpper(string(*s)))
	return s
}

// ToLower packs ` ToLower(s string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) ToLower() *StrChain {
	*s = StrChain(ToLower(string(*s)))
	return s
}

// Title returns packs `Title(s string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) Title() *StrChain {
	*s = StrChain(Title(string(*s)))
	return s
}

// Map packs `Map(mapping func(rune) rune, s string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) Map(mapping func(rune) rune) *StrChain {
	*s = StrChain(Map(mapping, string(*s)))
	return s
}

// Repeat packs `Repeat(s string, count int)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) Repeat(count int) *StrChain {
	*s = StrChain(Repeat(string(*s), count))
	return s
}

// Replace packs `Replace(s, old, new string, n int)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) Replace(old, new string, n int) *StrChain {
	*s = StrChain(Replace(string(*s), old, new, n))
	return s
}

// ReplaceAll packs `ReplaceAll(s, old, new string)`
// 	set `TextCollection.Text` to the result
func (s *StrChain) ReplaceAll(old, new string) *StrChain {
	*s = StrChain(ReplaceAll(string(*s), old, new))
	return s
}

// GbkToUtf8String packs `GbkToUtf8String(s string)`
func (s *StrChain) GbkToUtf8String() (*StrChain, error) {
	t, e := GbkToUtf8String(string(*s))
	if e != nil {
		*s = StrChain("")
		return s, e
	}
	*s = StrChain(t)
	return s, nil
}

// Utf8ToGbkString packs `Utf8ToGbkString(s string)`
func (s *StrChain) Utf8ToGbkString() (*StrChain, error) {
	t, e := Utf8ToGbkString(string(*s))
	if e != nil {
		*s = StrChain("")
		return s, e
	}
	*s = StrChain(t)
	return s, nil
}

// Big5ToUtf8String packs `Big5ToUtf8String(s string)`
func (s *StrChain) Big5ToUtf8String() (*StrChain, error) {
	t, e := Big5ToUtf8String(string(*s))
	if e != nil {
		*s = StrChain("")
		return s, e
	}
	*s = StrChain(t)
	return s, nil
}

// Utf8ToBig5String packs `Utf8ToBig5String(s string)`
func (s *StrChain) Utf8ToBig5String() (*StrChain, error) {
	t, e := Utf8ToBig5String(string(*s))
	if e != nil {
		*s = StrChain("")
		return s, e
	}
	*s = StrChain(t)
	return s, nil
}

// IsEqualString packs `IsEqualString(a, b string, ignoreCase bool) bool`
func (s *StrChain) IsEqualString(b string, ignoreCase bool) bool {
	return IsEqualString(string(*s), b, ignoreCase)
}

// TrimBOM packs `TrimBOM(line string) string`
func (s *StrChain) TrimBOM() *StrChain {
	*s = StrChain(TrimBOM(string(*s)))
	return s
}

// TrimFrontEndSpaceLine packs `TrimFrontEndSpaceLine(content string) string`
func (s *StrChain) TrimFrontEndSpaceLine() *StrChain {
	*s = StrChain(TrimFrontEndSpaceLine(string(*s)))
	return s
}
