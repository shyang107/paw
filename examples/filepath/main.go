package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

var (
	lg    = paw.Logger
	Info  = lg.Info
	Warn  = lg.Warn
	Error = lg.Error
	Debug = lg.Debug
	home  string
)

func init() {
	lg.Level = logrus.DebugLevel
	lg.Info()
	home = homedir()
}

func main() {
	Info()

	// 1. homedir
	fmt.Printf("HomeDir: %q\n", home)

	// 2. clean path
	var (
		paths = []string{
			"/",
			".",
			"./",
			"..",
			"../",
			"~",
			"~/.",
			"~/..",
			"~/../..",
			"~/../../..",
		}
		lgFields = make(map[string]interface{})
	)
	for i, path := range paths {
		key := fmt.Sprintf("path%d", i)
		lgFields[key] = path
	}
	lg.WithFields(lgFields).Debug("clean path")
	for i, path := range paths {
		paths[i] = clPath(path)
	}
	spew.Dump(paths)
	// 3. Glob
	pat := "[a-zA-z0-9]*"
	for _, path := range paths {
		glob(path, pat)
	}
	// 4. ReadDir
	readDir("/dev")
}

func readDir(path string) {
	lg.WithField("path", path).Info()
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		Error("path[" + path + "]: " + err.Error())
	}
	for i, file := range fis {
		fmt.Printf("%5d %q\n", i, file.Name())
	}
}

func glob(path string, pattern string) {
	Info()
	path = filepath.Join(path, pattern)
	fmt.Println("path =", path)
	files, err := filepath.Glob(path)
	if err != nil {
		Error(err)
	}
	for i, file := range files {
		file, _ := filepath.Abs(file)
		fmt.Printf("%5d %s\n", i, file)
	}
}

func homedir() string {
	Info()
	home := os.Getenv("HOME") + "/"
	return home
}

func clPath(path string) string {
	Info()
	hpath := path
	if strings.Contains(hpath, "~") {
		hpath = strings.ReplaceAll(hpath, "~", home)
	}

	cpath := filepath.Clean(hpath)
	apath := cpath
	if !filepath.IsAbs(cpath) {
		apath, err := filepath.Abs(cpath)
		if err != nil {
			apath += " " + err.Error()
		}
	}
	paw.SetLoggerFieldsOrder([]string{"path", "home", "clean", "abs"})
	lg.WithFields(logrus.Fields{
		"path":  path,
		"home":  hpath,
		"clean": cpath,
		"abs":   apath,
	}).Debug()
	return apath
}
