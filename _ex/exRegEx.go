package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/shyang107/paw"
)

const tmpText = "const twoMatch = `test string`;\nconst noMatches = `test ${ variabel }`;\nabcde ${ field1 } and ${ Field2}"

func exRegEx() {
	var re = regexp.MustCompile(`(?m)(\${.*?)(\b\w+\b)(.*?})`)
	tokens := map[string]string{
		"variabel": "[token_variabel]",
		"field1":   "[token_field1]",
		"field2":   "[token_field2]",
	}
	fmt.Println(tmpText)
	matchs := re.FindAllStringSubmatch(tmpText, -1)
	spew.Dump(matchs)
	tb := paw.TextBuilder{}
	result := tmpText
	for _, m := range matchs {
		tb.SetText(m[2]).ToLower()
		result = strings.ReplaceAll(result, m[0], tokens[tb.GetText()])
	}
	fmt.Println(result)
	// fmt.Println(re.ReplaceAllString(str, substitution))
}
