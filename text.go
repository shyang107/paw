package paw

import (
	"strings"
	"unicode"
	"unicode/utf8"
	// "github.com/gookit/color"
)

// TextCollection contains string `Text`
type TextCollection struct {
	Text string
}

// TextBuilder contains all tools which can be chained.
type TextBuilder struct {
	TextCollection
}

// Build return a instance of `TextCollection`
func (tb *TextBuilder) Build() TextCollection {
	return TextCollection{Text: tb.Text}
}

// SetText set `TextBuilder.TextCollection.Text` to `txt`
func (tb *TextBuilder) SetText(txt string) *TextBuilder {
	tb.Text = txt
	return tb
}

// GetText return `TextBuilder.TextCollection.Text`
func (tb *TextBuilder) GetText() string {
	return tb.TextCollection.Text
}

// String return `TextBuilder.TextCollection.Text`
func (tb *TextBuilder) String() string {
	return tb.TextCollection.Text
}

// GetAbbrString get a abbreviation of `TextBuilder.TextCollection.Text` and save to `TextBuilder.TextCollection.Text`
//
// 	`maxlen`: maimium length of the abbreviation
// 	`conSymbole`: tailing symbol of the abbreviation
func (tb *TextBuilder) GetAbbrString(maxlen int, contSymbol string) *TextBuilder {
	tb.Text = GetAbbrString(tb.TextCollection.Text, maxlen, contSymbol)
	return tb
}

// CountPlaceHolder return `nHan` and `nASCII`
//
// 	`nHan`: number of occupied space in screen for han-character
// 	`nASCII`: number of occupied space in screen for ASCII-character
func (tb *TextBuilder) CountPlaceHolder() (nHan int, nASCII int) {
	return CountPlaceHolder(tb.TextCollection.Text)
}

// HasChineseChar return true for that `str` include chinese character
//
// Example:
// 	HasChineseChar("abc 中文") return true
// 	HasChineseChar("abccefgh") return false
func (tb *TextBuilder) HasChineseChar() bool {
	return HasChineseChar(tb.TextCollection.Text)
}

// NumberBanner return numbers' string with length of `TextBuilder.TextCollection.Text`
//
// Example:
// 	TextBuilder.TextCollection.Text = "Text中文 Collection"
// 	nh, na := CountPlaceHolder（"Text中文 Collection"）
// 	--> nh=4, na=15 --> length = nh + na = 19
// 	NumberBanner() return "12345678901"
func (tb *TextBuilder) NumberBanner() *TextBuilder {
	h, a := tb.CountPlaceHolder()
	tb.TextCollection.Text = NumberBanner(h + a)
	return tb
}

// HasPrefix return `strings.HasPrefix(tb.TextCollection.Text, prefix)`
func (tb *TextBuilder) HasPrefix(prefix string) bool {
	return strings.HasPrefix(tb.TextCollection.Text, prefix)
}

// HasSuffix return `strings.HasSuffix(tb.TextCollection.Text, Suffix)`
func (tb *TextBuilder) HasSuffix(suffix string) bool {
	return strings.HasSuffix(tb.TextCollection.Text, suffix)
}

// Contains return `strings.Contains(tb.TextCollection.Text, substr)`
func (tb *TextBuilder) Contains(substr string) bool {
	return strings.Contains(tb.TextCollection.Text, substr)
}

// Trim packs `Trim(s, cutset)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Trim(cutset string) *TextBuilder {
	tb.TextCollection.Text = strings.Trim(tb.TextCollection.Text, cutset)
	return tb
}

// TrimFunc packs `TrimFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimFunc(f func(rune) bool) *TextBuilder {
	tb.TextCollection.Text = strings.TrimFunc(tb.TextCollection.Text, f)
	return tb
}

// TrimLeft packs `TrimLeft(s, cutset string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimLeft(cutset string) *TextBuilder {
	tb.TextCollection.Text = strings.TrimLeft(tb.TextCollection.Text, cutset)
	return tb
}

// TrimLeftFunc packs `TrimLeftFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimLeftFunc(f func(rune) bool) *TextBuilder {
	tb.TextCollection.Text = strings.TrimLeftFunc(tb.TextCollection.Text, f)
	return tb
}

// TrimPrefix packs `TrimPrefix(s, prefix string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimPrefix(s, prefix string) *TextBuilder {
	tb.TextCollection.Text = strings.TrimPrefix(tb.TextCollection.Text, prefix)
	return tb
}

// TrimRight packs `TrimRight(s, cutset string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimRight(s, cutset string) *TextBuilder {
	tb.TextCollection.Text = strings.TrimRight(tb.TextCollection.Text, cutset)
	return tb
}

// TrimRightFunc packs `TrimRightFunc(s string, f func(rune) bool)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimRightFunc(s string, f func(rune) bool) *TextBuilder {
	tb.TextCollection.Text = strings.TrimRightFunc(tb.TextCollection.Text, f)
	return tb
}

// TrimSpace packs `TrimSpace(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimSpace() *TextBuilder {
	tb.TextCollection.Text = strings.TrimSpace(tb.TextCollection.Text)
	return tb
}

// TrimSuffix packs `TrimSuffix(s, suffix string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) TrimSuffix(suffix string) *TextBuilder {
	tb.TextCollection.Text = strings.TrimSuffix(tb.TextCollection.Text, suffix)
	return tb
}

// ToUpper packs `ToUpper(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) ToUpper() *TextBuilder {
	tb.TextCollection.Text = strings.ToUpper(tb.TextCollection.Text)
	return tb
}

// ToTitle packs `ToTitle(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) ToTitle() *TextBuilder {
	tb.TextCollection.Text = strings.ToUpper(tb.TextCollection.Text)
	return tb
}

// ToLower packs ` ToLower(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) ToLower() *TextBuilder {
	tb.TextCollection.Text = strings.ToLower(tb.TextCollection.Text)
	return tb
}

// Title returns packs `Title(s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Title() *TextBuilder {
	tb.TextCollection.Text = strings.Title(tb.TextCollection.Text)
	return tb
}

// Map packs `Map(mapping func(rune) rune, s string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Map(mapping func(rune) rune) *TextBuilder {
	tb.TextCollection.Text = strings.Map(mapping, tb.TextCollection.Text)
	return tb
}

// Repeat packs `Repeat(s string, count int)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Repeat(count int) *TextBuilder {
	tb.TextCollection.Text = strings.Repeat(tb.TextCollection.Text, count)
	return tb
}

// Replace packs `Replace(s, old, new string, n int)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) Replace(old, new string, n int) *TextBuilder {
	tb.TextCollection.Text = strings.Replace(tb.TextCollection.Text, old, new, n)
	return tb
}

// ReplaceAll packs `ReplaceAll(s, old, new string)`
// 	set `TextCollection.Text` to the result
func (tb *TextBuilder) ReplaceAll(old, new string) *TextBuilder {
	tb.TextCollection.Text = strings.ReplaceAll(tb.TextCollection.Text, old, new)
	return tb
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
// rot13 := func(r rune) rune {
// 	switch {
// 	case r >= 'A' && r <= 'Z':
// 		return 'A' + (r-'A'+13)%26
// 	case r >= 'a' && r <= 'z':
// 		return 'a' + (r-'a'+13)%26
// 	}
// 	return r
// }
// fmt.Println(strings.Map(rot13, "'Twas brillig and the slithy gopher..."))
// out:
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
