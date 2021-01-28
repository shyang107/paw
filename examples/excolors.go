package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cast"
)

var (
	typeDesc = map[string]string{
		"di": "directory",
		"fi": "file",
		"ln": "symbolic link",
		"pi": "fifo file",
		"so": "socket file",
		"bd": "block (buffered) special file",
		"cd": "character (unbuffered) special file",
		"or": "symbolic link pointing to a non-existent file (orphan)",
		"mi": "non-existent file pointed to by a symbolic link (visible when you type ls -l)",
		"ex": "file which is executable (ie. has 'x' set in permissions)",
	}
	colors = make(map[string]string)
	exts   = []string{}
	// NoColor ...
	NoColor = os.Getenv("TERM") == "dumb" || !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))
)

func init() {
	getcolors()
}

func exColor() {
	// for _, c := range exts {
	// 	str := colorstr(colors[c], c+" : "+colors[c])
	// 	fmt.Println(str)
	// }

	fmt.Println("LSColors = map[string][]color.Attribute {")
	for _, c := range exts {
		// fmt.Println(c, colors[c])
		key := strings.TrimPrefix(c, "*")
		val := strings.ReplaceAll(colors[c], ";", ", ")
		fmt.Printf("\t%q :  []color.Attribute{ %s },\n", key, val)
	}
	fmt.Println("}")
}
func colorstr(code, s string) string {
	att := []color.Attribute{}
	for _, a := range strings.Split(code, ";") {
		att = append(att, color.Attribute(cast.ToInt(a)))
	}
	cs := color.New(att...)
	return cs.Sprint(s)
}

func fileColorStr(ext, s string) string {
	switch {
	case NoColor:
		return s
	default:
		if _, ok := colors[ext]; !ok {
			return s
		}
		return colorstr(colors[ext], s)
	}
}

func getcolors() {
	colorenv := os.Getenv("LS_COLORS")
	args := strings.Split(colorenv, ":")

	// colors := make(map[string]string)
	// ctypes := make(map[string]string)
	// exts := []string{}
	for _, a := range args {
		// fmt.Printf("%v\t", a)
		kv := strings.Split(a, "=")

		// fmt.Printf("%v\n", kv)
		if len(kv) == 2 {
			colors[kv[0]] = kv[1]
			exts = append(exts, kv[0])
		}
	}
	// sort.Strings(exts)
}
