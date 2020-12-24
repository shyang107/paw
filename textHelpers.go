package paw

// TextTools is the collections of tools of text
type TextTools interface {
	GetAbbrString(maxlen int, contSymbol string) TextTools
	Map(mapping func(rune) rune) TextTools
	NumberBanner() TextTools
	Repeat(count int) TextTools
	Replace(old, new string, n int) TextTools
	ReplaceAll(old, new string) TextTools
	Reverse() TextTools
	SetText(txt string) TextTools
	Title() TextTools
	ToLower() TextTools
	ToTitle() TextTools
	ToUpper() TextTools
	Trim(cutset string) TextTools
	TrimFunc(f func(rune) bool) TextTools
	TrimLeft(cutset string) TextTools
	TrimLeftFunc(f func(rune) bool) TextTools
	TrimPrefix(s, prefix string) TextTools
	TrimRight(s, cutset string) TextTools
	TrimRightFunc(s string, f func(rune) bool) TextTools
	TrimSpace() TextTools
	TrimSuffix(suffix string) TextTools
	TrimBOM() TextTools
	TrimFrontEndSpaceLine() TextTools

	Big5ToUtf8String() (TextTools, error)
	Utf8ToBig5String() (TextTools, error)
	GbkToUtf8String() (TextTools, error)
	Utf8ToGbkString() (TextTools, error)

	CountPlaceHolder() (nHan int, nASCII int)
	Contains(substr string) bool
	ContainsAny(chars string) bool
	ContainsRune(r rune) bool
	EqualFold(t string) bool
	Fields() []string
	FieldsFunc(f func(rune) bool) []string
	Index(substr string) int
	IndexAny(chars string) int
	IndexByte(c byte) int
	IndexFunc(f func(rune) bool) int
	IndexRune(r rune) int
	LastIndex(substr string) int
	LastIndexAny(chars string) int
	LastIndexByte(c byte) int
	LastIndexFunc(f func(rune) bool) int
	Split(sep string) []string
	SplitN(sep string, n int) []string
	SplitAfter(sep string) []string
	SplitAfterN(sep string, n int) []string

	GetText() string
	HasChineseChar() bool
	HasPrefix(prefix string) bool
	HasSuffix(suffix string) bool
	IsEqualString(b string, ignoreCase bool) bool
}

// TextBuilder contains all tools which can be chained.
type TextBuilder struct {
	Text    string
	TBError error
}

// NewTextBuilder return a instance of `TextBuilder` and return `TextTools`
func (tb *TextBuilder) NewTextBuilder(s string) TextTools {
	tb = &TextBuilder{Text: s}
	return tb
}

// String return `TextBuilder.TextCollection.Text`
func (tb *TextBuilder) String() string {
	return tb.Text
}

// SetText set `TextBuilder.Text` to `txt`
func (tb *TextBuilder) SetText(txt string) TextTools {
	tb.Text = txt
	return tb
}

// GetText return `TextBuilder.TextCollection.Text`
func (tb *TextBuilder) GetText() string {
	return tb.Text
}

// GetAbbrString get a abbreviation of `TextBuilder.TextCollection.Text` and save to `TextBuilder.TextCollection.Text`
//
// 	`maxlen`: maimium length of the abbreviation
// 	`conSymbole`: tailing symbol of the abbreviation
func (tb *TextBuilder) GetAbbrString(maxlen int, contSymbol string) TextTools {
	tb.Text = GetAbbrString(tb.Text, maxlen, contSymbol)
	return tb
}

// CountPlaceHolder return `nHan` and `nASCII`
//
// 	`nHan`: number of occupied space in screen for han-character
// 	`nASCII`: number of occupied space in screen for ASCII-character
func (tb *TextBuilder) CountPlaceHolder() (nHan int, nASCII int) {
	return CountPlaceHolder(tb.Text)
}

// HasChineseChar return true for that `str` include chinese character
//
// Example:
// 	HasChineseChar("abc 中文") return true
// 	HasChineseChar("abccefgh") return false
func (tb *TextBuilder) HasChineseChar() bool {
	return HasChineseChar(tb.Text)
}

// NumberBanner return numbers' string with length of `TextBuilder.TextCollection.Text`
//
// Example:
// 	TextBuilder.TextCollection.Text = "Text中文 Collection"
// 	nh, na := CountPlaceHolder（"Text中文 Collection"）
// 	--> nh=4, na=15 --> length = nh + na = 19
// 	NumberBanner() return "12345678901"
func (tb *TextBuilder) NumberBanner() TextTools {
	h, a := tb.CountPlaceHolder()
	tb.Text = NumberBanner(h + a)
	return tb
}

// Reverse packs `Reverse(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Reverse() TextTools {
	tb.Text = Reverse(tb.Text)
	return tb
}

// HasPrefix return `HasPrefix(tb.Text, prefix)`
func (tb *TextBuilder) HasPrefix(prefix string) bool {
	return HasPrefix(tb.Text, prefix)
}

// HasSuffix return `HasSuffix(tb.Text, Suffix)`
func (tb *TextBuilder) HasSuffix(suffix string) bool {
	return HasSuffix(tb.Text, suffix)
}

// Contains return `Contains(tb.Text, substr)`
func (tb *TextBuilder) Contains(substr string) bool {
	return Contains(tb.Text, substr)
}

// ContainsAny return `ContainsAny(tb.Text, chars)`
func (tb *TextBuilder) ContainsAny(chars string) bool {
	return ContainsAny(tb.Text, chars)
}

// Fields return Fields(tb.Text)
func (tb *TextBuilder) Fields() []string {
	return Fields(tb.Text)
}

// FieldsFunc return FieldsFunc(tb.Text, f)
func (tb *TextBuilder) FieldsFunc(f func(rune) bool) []string {
	return FieldsFunc(tb.Text, f)
}

// ContainsAny return ContainsRune(tb.Text, r) bool
func (tb *TextBuilder) ContainsRune(r rune) bool {
	return ContainsRune(tb.Text, r)
}

// EqualFold return EqualFold(tb.Text,, t) bool
func (tb *TextBuilder) EqualFold(t string) bool {
	return EqualFold(tb.Text, t)
}

// Index return Index(tb.Text, substr) int
func (tb *TextBuilder) Index(substr string) int {
	return Index(tb.Text, substr)
}

// IndexAny return IndexAny(tb.Text, chars) int
func (tb *TextBuilder) IndexAny(chars string) int {
	return IndexAny(tb.Text, chars)
}

// IndexByte return IndexByte(tb.Text, c) int
func (tb *TextBuilder) IndexByte(c byte) int {
	return IndexByte(tb.Text, c)
}

// IndexFunc return IndexFunc(tb.Text, f) int
func (tb *TextBuilder) IndexFunc(f func(rune) bool) int {
	return IndexFunc(tb.Text, f)
}

// IndexRune return IndexRune(tb.Text, r) int
func (tb *TextBuilder) IndexRune(r rune) int {
	return IndexRune(tb.Text, r)
}

// LastIndex return LastIndex(tb.Text, substr) int
func (tb *TextBuilder) LastIndex(substr string) int {
	return LastIndex(tb.Text, substr)
}

// LastIndexAny return LastIndexAny(tb.Text, chars) int
func (tb *TextBuilder) LastIndexAny(chars string) int {
	return LastIndexAny(tb.Text, chars)
}

// LastIndexByte return LastIndexByte(tb.Text, c) int
func (tb *TextBuilder) LastIndexByte(c byte) int {
	return LastIndexByte(tb.Text, c)
}

// LastIndexFunc return LastIndexFunc(tb.Text, f) int
func (tb *TextBuilder) LastIndexFunc(f func(rune) bool) int {
	return LastIndexFunc(tb.Text, f)
}

// Split return Split(tb.Text, sep) []string
func (tb *TextBuilder) Split(sep string) []string {
	return Split(tb.Text, sep)
}

// Split return Split(tb.Text, sep) []string
func (tb *TextBuilder) SplitN(sep string, n int) []string {
	return SplitN(tb.Text, sep, n)
}

// Split return Split(tb.Text, sep) []string
func (tb *TextBuilder) SplitAfter(sep string) []string {
	return SplitAfter(tb.Text, sep)
}

// Split return Split(tb.Text, sep) []string
func (tb *TextBuilder) SplitAfterN(sep string, n int) []string {
	return SplitAfterN(tb.Text, sep, n)
}

// Trim packs `Trim(s, cutset)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Trim(cutset string) TextTools {
	tb.Text = Trim(tb.Text, cutset)
	return tb
}

// TrimFunc packs `TrimFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimFunc(f func(rune) bool) TextTools {
	tb.Text = TrimFunc(tb.Text, f)
	return tb
}

// TrimLeft packs `TrimLeft(s, cutset string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimLeft(cutset string) TextTools {
	tb.Text = TrimLeft(tb.Text, cutset)
	return tb
}

// TrimLeftFunc packs `TrimLeftFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimLeftFunc(f func(rune) bool) TextTools {
	tb.Text = TrimLeftFunc(tb.Text, f)
	return tb
}

// TrimPrefix packs `TrimPrefix(s, prefix string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimPrefix(s, prefix string) TextTools {
	tb.Text = TrimPrefix(tb.Text, prefix)
	return tb
}

// TrimRight packs `TrimRight(s, cutset string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimRight(s, cutset string) TextTools {
	tb.Text = TrimRight(tb.Text, cutset)
	return tb
}

// TrimRightFunc packs `TrimRightFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimRightFunc(s string, f func(rune) bool) TextTools {
	tb.Text = TrimRightFunc(tb.Text, f)
	return tb
}

// TrimSpace packs `TrimSpace(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimSpace() TextTools {
	tb.Text = TrimSpace(tb.Text)
	return tb
}

// TrimSuffix packs `TrimSuffix(s, suffix string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimSuffix(suffix string) TextTools {
	tb.Text = TrimSuffix(tb.Text, suffix)
	return tb
}

// ToUpper packs `ToUpper(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) ToUpper() TextTools {
	tb.Text = ToUpper(tb.Text)
	return tb
}

// ToTitle packs `ToTitle(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) ToTitle() TextTools {
	tb.Text = ToUpper(tb.Text)
	return tb
}

// ToLower packs ` ToLower(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) ToLower() TextTools {
	tb.Text = ToLower(tb.Text)
	return tb
}

// Title returns packs `Title(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Title() TextTools {
	tb.Text = Title(tb.Text)
	return tb
}

// Map packs `Map(mapping func(rune) rune, s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Map(mapping func(rune) rune) TextTools {
	tb.Text = Map(mapping, tb.Text)
	return tb
}

// Repeat packs `Repeat(s string, count int)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Repeat(count int) TextTools {
	tb.Text = Repeat(tb.Text, count)
	return tb
}

// Replace packs `Replace(s, old, new string, n int)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Replace(old, new string, n int) TextTools {
	tb.Text = Replace(tb.Text, old, new, n)
	return tb
}

// ReplaceAll packs `ReplaceAll(s, old, new string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) ReplaceAll(old, new string) TextTools {
	tb.Text = ReplaceAll(tb.Text, old, new)
	return tb
}

// GbkToUtf8String packs `GbkToUtf8String(s string)`
func (tb *TextBuilder) GbkToUtf8String() (TextTools, error) {
	s, e := GbkToUtf8String(tb.Text)
	if e != nil {
		tb.Text = ""
		return tb, e
	}
	tb.Text = s
	return tb, nil
}

// Utf8ToGbkString packs `Utf8ToGbkString(s string)`
func (tb *TextBuilder) Utf8ToGbkString() (TextTools, error) {
	s, e := Utf8ToGbkString(tb.Text)
	if e != nil {
		tb.Text = ""
		return tb, e
	}
	tb.Text = s
	return tb, nil
}

// Big5ToUtf8String packs `Big5ToUtf8String(s string)`
func (tb *TextBuilder) Big5ToUtf8String() (TextTools, error) {
	s, e := Big5ToUtf8String(tb.Text)
	if e != nil {
		tb.Text = ""
		return tb, e
	}
	tb.Text = s
	return tb, nil
}

// Utf8ToBig5String packs `Utf8ToBig5String(s string)`
func (tb *TextBuilder) Utf8ToBig5String() (TextTools, error) {
	s, e := Utf8ToBig5String(tb.Text)
	if e != nil {
		tb.Text = ""
		return tb, e
	}
	tb.Text = s
	return tb, nil
}

// IsEqualString packs `IsEqualString(a, b string, ignoreCase bool) bool`
func (tb *TextBuilder) IsEqualString(b string, ignoreCase bool) bool {
	return IsEqualString(tb.Text, b, ignoreCase)
}

// TrimBOM packs `TrimBOM(line string) string`
func (tb *TextBuilder) TrimBOM() TextTools {
	tb.Text = TrimBOM(tb.Text)
	return tb
}

// TrimFrontEndSpaceLine packs `TrimFrontEndSpaceLine(content string) string`
func (tb *TextBuilder) TrimFrontEndSpaceLine() TextTools {
	tb.Text = TrimFrontEndSpaceLine(tb.Text)
	return tb
}
