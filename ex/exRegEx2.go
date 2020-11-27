package main

import (
	"fmt"
	"regexp"
)

func exRegEx2() {
	text := "field1,field2 filed3|field4	field5"
	r := regexp.MustCompile(`,| |\t|\|`)
	result := r.Split(text, -1)
	fmt.Printf("case 1 : %q\n\tpattern: %q\n\t%v\n", text, r.String(), result)
	r = regexp.MustCompile(`(?P<year>\d{4})-(?P<month>\d{1,2})-(?P<day>\d{1,2})`)
	text = "2020-10-20 2020-1-05 2020-11-3"
	result2 := r.ReplaceAllString(text, "${year}/${month}/${day}")
	fmt.Printf("case 2 : %q\n\tpattern: %q\n\t%v\n", text, r.String(), result2)
}
