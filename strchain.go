package paw

// StrChain contains all tools which can be chained.
type StrChain struct {
	S string
	// TBError error
}

// NewStrChain return a instance of `StrChain` and return `*StrChain`
func (t *StrChain) NewStrChain(s string) *StrChain {
	t = &StrChain{s}
	return t
}

// String return `StrChain.sCollection.s`
func (t *StrChain) String() string {
	return t.S
}

// SetText set `StrChain.s` to `txt`
func (t *StrChain) SetText(txt string) *StrChain {
	t.S = txt
	return t
}

// GetText return `StrChain.sCollection.s`
func (t *StrChain) GetText() string {
	return t.S
}

// Len will return the lenth of t.S (would be the sizes of []bytes)
func (t *StrChain) Len() int {
	return len(t.S)
}

// Bytes will convert the string t.S to []byte
//
// Example:
// 	b := StrChain{"ABC€"}
// 	fmt.Println(b.Bytes()) // [65 66 67 226 130 172]
func (t *StrChain) Bytes() []byte {
	return []byte(t.S)
}

// Runes will convert the string t.S to []rune
//
// Example:
// 	r := StrChain{"ABC€"}
// 	fmt.Println(r.Runes())        	// [65 66 67 8364]
// 	fmt.Printf("%U\n", r.Rune()) 	// [U+0041 U+0042 U+0043 U+20AC]
func (t *StrChain) Runes() []rune {
	return []rune(t.S)
}

// GetAbbrString get a abbreviation of `StrChain.sCollection.s` and save to `StrChain.sCollection.s`
//
// 	`maxlen`: maimium length of the abbreviation
// 	`conSymbole`: tailing symbol of the abbreviation
func (t *StrChain) GetAbbrString(maxlen int, contSymbol string) *StrChain {
	t.S = GetAbbrString(t.S, maxlen, contSymbol)
	return t
}

// CountPlaceHolder return `nHan` and `nASCII`
//
// 	`nHan`: number of occupied space in screen for han-character
// 	`nASCII`: number of occupied space in screen for ASCII-character
func (t *StrChain) CountPlaceHolder() (nHan int, nASCII int) {
	return CountPlaceHolder(t.S)
}

// HasChineseChar return true for that `str` include chinese character
//
// Example:
// 	HasChineseChar("abc 中文") return true
// 	HasChineseChar("abccefgh") return false
func (t *StrChain) HasChineseChar() bool {
	return HasChineseChar(t.S)
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
	t.S = NumberBanner(h + a)
	return t
}

// Reverse packs `Reverse(s string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) Reverse() *StrChain {
	t.S = Reverse(t.S)
	return t
}

// HasPrefix return `HasPrefix(t.S, prefix)`
func (t *StrChain) HasPrefix(prefix string) bool {
	return HasPrefix(t.S, prefix)
}

// HasSuffix return `HasSuffix(t.S, Suffix)`
func (t *StrChain) HasSuffix(suffix string) bool {
	return HasSuffix(t.S, suffix)
}

// Contains return `Contains(t.S, substr)`
func (t *StrChain) Contains(substr string) bool {
	return Contains(t.S, substr)
}

// ContainsAny return `ContainsAny(t.S, chars)`
func (t *StrChain) ContainsAny(chars string) bool {
	return ContainsAny(t.S, chars)
}

// Fields return Fields(t.S)
func (t *StrChain) Fields() []string {
	return Fields(t.S)
}

// FieldsFunc return FieldsFunc(t.S, f)
func (t *StrChain) FieldsFunc(f func(rune) bool) []string {
	return FieldsFunc(t.S, f)
}

// ContainsAny return ContainsRune(t.S, r) bool
func (t *StrChain) ContainsRune(r rune) bool {
	return ContainsRune(t.S, r)
}

// EqualFold return EqualFold(t.S,, t) bool
func (t *StrChain) EqualFold(s string) bool {
	return EqualFold(t.S, s)
}

// Index return Index(t.S, substr) int
func (t *StrChain) Index(substr string) int {
	return Index(t.S, substr)
}

// IndexAny return IndexAny(t.S, chars) int
func (t *StrChain) IndexAny(chars string) int {
	return IndexAny(t.S, chars)
}

// IndexByte return IndexByte(t.S, c) int
func (t *StrChain) IndexByte(c byte) int {
	return IndexByte(t.S, c)
}

// IndexFunc return IndexFunc(t.S, f) int
func (t *StrChain) IndexFunc(f func(rune) bool) int {
	return IndexFunc(t.S, f)
}

// IndexRune return IndexRune(t.S, r) int
func (t *StrChain) IndexRune(r rune) int {
	return IndexRune(t.S, r)
}

// LastIndex return LastIndex(t.S, substr) int
func (t *StrChain) LastIndex(substr string) int {
	return LastIndex(t.S, substr)
}

// LastIndexAny return LastIndexAny(t.S, chars) int
func (t *StrChain) LastIndexAny(chars string) int {
	return LastIndexAny(t.S, chars)
}

// LastIndexByte return LastIndexByte(t.S, c) int
func (t *StrChain) LastIndexByte(c byte) int {
	return LastIndexByte(t.S, c)
}

// LastIndexFunc return LastIndexFunc(t.S, f) int
func (t *StrChain) LastIndexFunc(f func(rune) bool) int {
	return LastIndexFunc(t.S, f)
}

// Split return Split(t.S, sep) []string
func (t *StrChain) Split(sep string) []string {
	return Split(t.S, sep)
}

// Split return Split(t.S, sep) []string
func (t *StrChain) SplitN(sep string, n int) []string {
	return SplitN(t.S, sep, n)
}

// Split return Split(t.S, sep) []string
func (t *StrChain) SplitAfter(sep string) []string {
	return SplitAfter(t.S, sep)
}

// Split return Split(t.S, sep) []string
func (t *StrChain) SplitAfterN(sep string, n int) []string {
	return SplitAfterN(t.S, sep, n)
}

// Trim packs `Trim(s, cutset)`
// 	set `TextCollection.s` to the result
func (t *StrChain) Trim(cutset string) *StrChain {
	t.S = Trim(t.S, cutset)
	return t
}

// TrimFunc packs `TrimFunc(s string, f func(rune) bool)`
// 	set `TextCollection.s` to the result
func (t *StrChain) TrimFunc(f func(rune) bool) *StrChain {
	t.S = TrimFunc(t.S, f)
	return t
}

// TrimLeft packs `TrimLeft(s, cutset string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) TrimLeft(cutset string) *StrChain {
	t.S = TrimLeft(t.S, cutset)
	return t
}

// TrimLeftFunc packs `TrimLeftFunc(s string, f func(rune) bool)`
// 	set `TextCollection.s` to the result
func (t *StrChain) TrimLeftFunc(f func(rune) bool) *StrChain {
	t.S = TrimLeftFunc(t.S, f)
	return t
}

// TrimPrefix packs `TrimPrefix(s, prefix string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) TrimPrefix(prefix string) *StrChain {
	t.S = TrimPrefix(t.S, prefix)
	return t
}

// TrimRight packs `TrimRight(s, cutset string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) TrimRight(cutset string) *StrChain {
	t.S = TrimRight(t.S, cutset)
	return t
}

// TrimRightFunc packs `TrimRightFunc(s string, f func(rune) bool)`
// 	set `TextCollection.s` to the result
func (t *StrChain) TrimRightFunc(f func(rune) bool) *StrChain {
	t.S = TrimRightFunc(t.S, f)
	return t
}

// TrimSpace packs `TrimSpace(s string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) TrimSpace() *StrChain {
	t.S = TrimSpace(t.S)
	return t
}

// TrimSuffix packs `TrimSuffix(s, suffix string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) TrimSuffix(suffix string) *StrChain {
	t.S = TrimSuffix(t.S, suffix)
	return t
}

// ToUpper packs `ToUpper(s string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) ToUpper() *StrChain {
	t.S = ToUpper(t.S)
	return t
}

// ToTitle packs `ToTitle(s string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) ToTitle() *StrChain {
	t.S = ToTitle(t.S)
	return t
}

// ToLower packs ` ToLower(s string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) ToLower() *StrChain {
	t.S = ToLower(t.S)
	return t
}

// Title returns packs `Title(s string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) Title() *StrChain {
	t.S = Title(t.S)
	return t
}

// Map packs `Map(mapping func(rune) rune, s string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) Map(mapping func(rune) rune) *StrChain {
	t.S = Map(mapping, t.S)
	return t
}

// Repeat packs `Repeat(s string, count int)`
// 	set `TextCollection.s` to the result
func (t *StrChain) Repeat(count int) *StrChain {
	t.S = Repeat(t.S, count)
	return t
}

// Replace packs `Replace(s, old, new string, n int)`
// 	set `TextCollection.s` to the result
func (t *StrChain) Replace(old, new string, n int) *StrChain {
	t.S = Replace(t.S, old, new, n)
	return t
}

// ReplaceAll packs `ReplaceAll(s, old, new string)`
// 	set `TextCollection.s` to the result
func (t *StrChain) ReplaceAll(old, new string) *StrChain {
	t.S = ReplaceAll(t.S, old, new)
	return t
}

// GbkToUtf8String packs `GbkToUtf8String(s string)`
func (t *StrChain) GbkToUtf8String() (*StrChain, error) {
	s, e := GbkToUtf8String(t.S)
	if e != nil {
		t.S = ""
		return t, e
	}
	t.S = s
	return t, nil
}

// Utf8ToGbkString packs `Utf8ToGbkString(s string)`
func (t *StrChain) Utf8ToGbkString() (*StrChain, error) {
	s, e := Utf8ToGbkString(t.S)
	if e != nil {
		t.S = ""
		return t, e
	}
	t.S = s
	return t, nil
}

// Big5ToUtf8String packs `Big5ToUtf8String(s string)`
func (t *StrChain) Big5ToUtf8String() (*StrChain, error) {
	s, e := Big5ToUtf8String(t.S)
	if e != nil {
		t.S = ""
		return t, e
	}
	t.S = s
	return t, nil
}

// Utf8ToBig5String packs `Utf8ToBig5String(s string)`
func (t *StrChain) Utf8ToBig5String() (*StrChain, error) {
	s, e := Utf8ToBig5String(t.S)
	if e != nil {
		t.S = ""
		return t, e
	}
	t.S = s
	return t, nil
}

// IsEqualString packs `IsEqualString(a, b string, ignoreCase bool) bool`
func (t *StrChain) IsEqualString(b string, ignoreCase bool) bool {
	return IsEqualString(t.S, b, ignoreCase)
}

// TrimBOM packs `TrimBOM(line string) string`
func (t *StrChain) TrimBOM() *StrChain {
	t.S = TrimBOM(t.S)
	return t
}

// TrimFrontEndSpaceLine packs `TrimFrontEndSpaceLine(content string) string`
func (t *StrChain) TrimFrontEndSpaceLine() *StrChain {
	t.S = TrimFrontEndSpaceLine(t.S)
	return t
}

// The following is adopted from github.com/mattn/go-runewidth

// FillLeft return string filled in left by spaces in w cells
func (t *StrChain) FillLeft(w int) *StrChain {
	t.S = FillLeft(t.S, w)
	return t
}

// FillRight return string filled in left by spaces in w cells
func (t *StrChain) FillRight(w int) *StrChain {
	t.S = FillRight(t.S, w)
	return t
}

// StringWidth will return width as you can see (the numbers of placeholders on terminal)
func (t *StrChain) StringWidth() int {
	return StringWidth(t.S)
}

// // RuneWidth returns the number of cells in r. See http://www.unicode.org/reports/tr11/
// func (t *StrChain) RuneWidth(r rune) int {
// 	return RuneWidth(r)
// }

// Truncate return string truncated with w cells
func (t *StrChain) Truncate(w int, tail string) *StrChain {
	t.S = Truncate(t.S, w, tail)
	return t
}

// Wrap return string wrapped with w cells
func (t *StrChain) Wrap(w int) *StrChain {
	t.S = Wrap(t.S, w)
	return t
}
