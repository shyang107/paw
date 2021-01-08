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
	Text string
	// TBError error
}

// NewTextBuilder return a instance of `TextBuilder` and return `TextTools`
func (t *TextBuilder) NewTextBuilder(s string) TextTools {
	t = &TextBuilder{s}
	return t
}

// String return `TextBuilder.TextCollection.Text`
func (t *TextBuilder) String() string {
	return t.Text
}

// SetText set `TextBuilder.Text` to `txt`
func (t *TextBuilder) SetText(txt string) TextTools {
	t.Text = txt
	return t
}

// GetText return `TextBuilder.TextCollection.Text`
func (t *TextBuilder) GetText() string {
	return t.Text
}

// GetAbbrString get a abbreviation of `TextBuilder.TextCollection.Text` and save to `TextBuilder.TextCollection.Text`
//
// 	`maxlen`: maimium length of the abbreviation
// 	`conSymbole`: tailing symbol of the abbreviation
func (t *TextBuilder) GetAbbrString(maxlen int, contSymbol string) TextTools {
	t.Text = GetAbbrString(t.Text, maxlen, contSymbol)
	return t
}

// CountPlaceHolder return `nHan` and `nASCII`
//
// 	`nHan`: number of occupied space in screen for han-character
// 	`nASCII`: number of occupied space in screen for ASCII-character
func (t *TextBuilder) CountPlaceHolder() (nHan int, nASCII int) {
	return CountPlaceHolder(t.Text)
}

// HasChineseChar return true for that `str` include chinese character
//
// Example:
// 	HasChineseChar("abc 中文") return true
// 	HasChineseChar("abccefgh") return false
func (t *TextBuilder) HasChineseChar() bool {
	return HasChineseChar(t.Text)
}

// NumberBanner return numbers' string with length of `TextBuilder.TextCollection.Text`
//
// Example:
// 	TextBuilder.TextCollection.Text = "Text中文 Collection"
// 	nh, na := CountPlaceHolder（"Text中文 Collection"）
// 	--> nh=4, na=15 --> length = nh + na = 19
// 	NumberBanner() return "12345678901"
func (t *TextBuilder) NumberBanner() TextTools {
	h, a := t.CountPlaceHolder()
	t.Text = NumberBanner(h + a)
	return t
}

// Reverse packs `Reverse(s string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) Reverse() TextTools {
	t.Text = Reverse(t.Text)
	return t
}

// HasPrefix return `HasPrefix(t.Text, prefix)`
func (t *TextBuilder) HasPrefix(prefix string) bool {
	return HasPrefix(t.Text, prefix)
}

// HasSuffix return `HasSuffix(t.Text, Suffix)`
func (t *TextBuilder) HasSuffix(suffix string) bool {
	return HasSuffix(t.Text, suffix)
}

// Contains return `Contains(t.Text, substr)`
func (t *TextBuilder) Contains(substr string) bool {
	return Contains(t.Text, substr)
}

// ContainsAny return `ContainsAny(t.Text, chars)`
func (t *TextBuilder) ContainsAny(chars string) bool {
	return ContainsAny(t.Text, chars)
}

// Fields return Fields(t.Text)
func (t *TextBuilder) Fields() []string {
	return Fields(t.Text)
}

// FieldsFunc return FieldsFunc(t.Text, f)
func (t *TextBuilder) FieldsFunc(f func(rune) bool) []string {
	return FieldsFunc(t.Text, f)
}

// ContainsAny return ContainsRune(t.Text, r) bool
func (t *TextBuilder) ContainsRune(r rune) bool {
	return ContainsRune(t.Text, r)
}

// EqualFold return EqualFold(t.Text,, t) bool
func (t *TextBuilder) EqualFold(s string) bool {
	return EqualFold(t.Text, s)
}

// Index return Index(t.Text, substr) int
func (t *TextBuilder) Index(substr string) int {
	return Index(t.Text, substr)
}

// IndexAny return IndexAny(t.Text, chars) int
func (t *TextBuilder) IndexAny(chars string) int {
	return IndexAny(t.Text, chars)
}

// IndexByte return IndexByte(t.Text, c) int
func (t *TextBuilder) IndexByte(c byte) int {
	return IndexByte(t.Text, c)
}

// IndexFunc return IndexFunc(t.Text, f) int
func (t *TextBuilder) IndexFunc(f func(rune) bool) int {
	return IndexFunc(t.Text, f)
}

// IndexRune return IndexRune(t.Text, r) int
func (t *TextBuilder) IndexRune(r rune) int {
	return IndexRune(t.Text, r)
}

// LastIndex return LastIndex(t.Text, substr) int
func (t *TextBuilder) LastIndex(substr string) int {
	return LastIndex(t.Text, substr)
}

// LastIndexAny return LastIndexAny(t.Text, chars) int
func (t *TextBuilder) LastIndexAny(chars string) int {
	return LastIndexAny(t.Text, chars)
}

// LastIndexByte return LastIndexByte(t.Text, c) int
func (t *TextBuilder) LastIndexByte(c byte) int {
	return LastIndexByte(t.Text, c)
}

// LastIndexFunc return LastIndexFunc(t.Text, f) int
func (t *TextBuilder) LastIndexFunc(f func(rune) bool) int {
	return LastIndexFunc(t.Text, f)
}

// Split return Split(t.Text, sep) []string
func (t *TextBuilder) Split(sep string) []string {
	return Split(t.Text, sep)
}

// Split return Split(t.Text, sep) []string
func (t *TextBuilder) SplitN(sep string, n int) []string {
	return SplitN(t.Text, sep, n)
}

// Split return Split(t.Text, sep) []string
func (t *TextBuilder) SplitAfter(sep string) []string {
	return SplitAfter(t.Text, sep)
}

// Split return Split(t.Text, sep) []string
func (t *TextBuilder) SplitAfterN(sep string, n int) []string {
	return SplitAfterN(t.Text, sep, n)
}

// Trim packs `Trim(s, cutset)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) Trim(cutset string) TextTools {
	t.Text = Trim(t.Text, cutset)
	return t
}

// TrimFunc packs `TrimFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) TrimFunc(f func(rune) bool) TextTools {
	t.Text = TrimFunc(t.Text, f)
	return t
}

// TrimLeft packs `TrimLeft(s, cutset string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) TrimLeft(cutset string) TextTools {
	t.Text = TrimLeft(t.Text, cutset)
	return t
}

// TrimLeftFunc packs `TrimLeftFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) TrimLeftFunc(f func(rune) bool) TextTools {
	t.Text = TrimLeftFunc(t.Text, f)
	return t
}

// TrimPrefix packs `TrimPrefix(s, prefix string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) TrimPrefix(s, prefix string) TextTools {
	t.Text = TrimPrefix(t.Text, prefix)
	return t
}

// TrimRight packs `TrimRight(s, cutset string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) TrimRight(s, cutset string) TextTools {
	t.Text = TrimRight(t.Text, cutset)
	return t
}

// TrimRightFunc packs `TrimRightFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) TrimRightFunc(s string, f func(rune) bool) TextTools {
	t.Text = TrimRightFunc(t.Text, f)
	return t
}

// TrimSpace packs `TrimSpace(s string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) TrimSpace() TextTools {
	t.Text = TrimSpace(t.Text)
	return t
}

// TrimSuffix packs `TrimSuffix(s, suffix string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) TrimSuffix(suffix string) TextTools {
	t.Text = TrimSuffix(t.Text, suffix)
	return t
}

// ToUpper packs `ToUpper(s string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) ToUpper() TextTools {
	t.Text = ToUpper(t.Text)
	return t
}

// ToTitle packs `ToTitle(s string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) ToTitle() TextTools {
	t.Text = ToUpper(t.Text)
	return t
}

// ToLower packs ` ToLower(s string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) ToLower() TextTools {
	t.Text = ToLower(t.Text)
	return t
}

// Title returns packs `Title(s string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) Title() TextTools {
	t.Text = Title(t.Text)
	return t
}

// Map packs `Map(mapping func(rune) rune, s string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) Map(mapping func(rune) rune) TextTools {
	t.Text = Map(mapping, t.Text)
	return t
}

// Repeat packs `Repeat(s string, count int)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) Repeat(count int) TextTools {
	t.Text = Repeat(t.Text, count)
	return t
}

// Replace packs `Replace(s, old, new string, n int)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) Replace(old, new string, n int) TextTools {
	t.Text = Replace(t.Text, old, new, n)
	return t
}

// ReplaceAll packs `ReplaceAll(s, old, new string)`
// 	set `TextCollection.Text` to the result
func (t *TextBuilder) ReplaceAll(old, new string) TextTools {
	t.Text = ReplaceAll(t.Text, old, new)
	return t
}

// GbkToUtf8String packs `GbkToUtf8String(s string)`
func (t *TextBuilder) GbkToUtf8String() (TextTools, error) {
	s, e := GbkToUtf8String(t.Text)
	if e != nil {
		t.Text = ""
		return t, e
	}
	t.Text = s
	return t, nil
}

// Utf8ToGbkString packs `Utf8ToGbkString(s string)`
func (t *TextBuilder) Utf8ToGbkString() (TextTools, error) {
	s, e := Utf8ToGbkString(t.Text)
	if e != nil {
		t.Text = ""
		return t, e
	}
	t.Text = s
	return t, nil
}

// Big5ToUtf8String packs `Big5ToUtf8String(s string)`
func (t *TextBuilder) Big5ToUtf8String() (TextTools, error) {
	s, e := Big5ToUtf8String(t.Text)
	if e != nil {
		t.Text = ""
		return t, e
	}
	t.Text = s
	return t, nil
}

// Utf8ToBig5String packs `Utf8ToBig5String(s string)`
func (t *TextBuilder) Utf8ToBig5String() (TextTools, error) {
	s, e := Utf8ToBig5String(t.Text)
	if e != nil {
		t.Text = ""
		return t, e
	}
	t.Text = s
	return t, nil
}

// IsEqualString packs `IsEqualString(a, b string, ignoreCase bool) bool`
func (t *TextBuilder) IsEqualString(b string, ignoreCase bool) bool {
	return IsEqualString(t.Text, b, ignoreCase)
}

// TrimBOM packs `TrimBOM(line string) string`
func (t *TextBuilder) TrimBOM() TextTools {
	t.Text = TrimBOM(t.Text)
	return t
}

// TrimFrontEndSpaceLine packs `TrimFrontEndSpaceLine(content string) string`
func (t *TextBuilder) TrimFrontEndSpaceLine() TextTools {
	t.Text = TrimFrontEndSpaceLine(t.Text)
	return t
}
