package main

import (
	"regexp"
	"strings"
	"testing"
)

func BenchmarkReplace2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		strings.Replace("seafood seafood", "foo", "oof", -1)
	}
}
func BenchmarkReplace5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		strings.Replace("seafood seafood seafood seafood seafood",
			"foo", "oof", -1)
	}
}
func BenchmarkReplace10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		strings.Replace("seafood seafood seafood seafood seafood seafood seafood seafood seafood seafood", "foo", "oof", -1)
	}
}
func BenchmarkRegexp2(b *testing.B) {
	r, _ := regexp.Compile("foo")
	for n := 0; n < b.N; n++ {
		r.ReplaceAllString("seafood seafood", "oof")
	}
}
func BenchmarkRegexp5(b *testing.B) {
	r, _ := regexp.Compile("foo")
	for n := 0; n < b.N; n++ {
		r.ReplaceAllString("seafood seafood seafood seafood seafood", "oof")
	}
}
func BenchmarkRegexp10(b *testing.B) {
	r, _ := regexp.Compile("foo")
	for n := 0; n < b.N; n++ {
		r.ReplaceAllString("seafood seafood seafood seafood seafood seafood seafood seafood seafood seafood", "oof")
	}
}
