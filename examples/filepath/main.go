package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/karrick/godirwalk"
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
	home = myhomedir()
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
			"./..//..",
			"..",
			// "~",
			// "~/.",
			"~/..",
			// "~/../..",
			// "~/../../..",
		}
		lgFields = make(map[string]interface{})
	)
	for i, path := range paths {
		key := fmt.Sprintf("path%d", i)
		lgFields[key] = path
	}
	lg.WithFields(lgFields).Debug("clean path")
	for i, path := range paths {
		paths[i], _ = filepath.Abs(path) //clPath(path)
	}
	spew.Dump(paths)
	// return
	// 3. Glob
	pat := "[a-zA-Z]*"
	for _, path := range paths {
		glob(path, pat)
	}
	// 4. ReadDir
	pat = "^[a-zA-z]+"
	re := regexp.MustCompile(pat)
	path := ".."
	path, _ = filepath.Abs(path)
	greadDir(path, re)
	readDir(path)
}

func greadDir(path string, re *regexp.Regexp) {
	// path, _ = filepath.Abs(path)
	lg.WithFields(logrus.Fields{
		"path": path,
		"re":   re.String(),
	}).Info()

	files, err := godirwalk.ReadDirnames(path, nil)
	if err != nil {
		Error("path[" + path + "]: " + err.Error())
	}
	sort.Strings(files)
	for i, file := range files {
		if !re.MatchString(file) {
			continue
		}
		// file, _ := filepath.Abs(filepath.Join(path, file))
		file := filepath.Join(path, file)
		fmt.Printf("%5d %q\n", i, file)
	}
}

func readDir(path string) {
	lg.WithField("path", path).Info()
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		Error("path[" + path + "]: " + err.Error())
	}
	for i, fi := range fis {
		// path, _ := filepath.Abs(filepath.Join(path, fi.Name()))
		path := filepath.Join(path, fi.Name())
		fmt.Printf("%5d %q\n", i, path)
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
		// file, _ := filepath.Abs(file)
		fmt.Printf("%5d %s\n", i, file)
	}
}

func myhomedir() string {
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
	// cpath := hpath
	cpath := filepath.Clean(hpath)
	// cpath := simplifyPath(hpath)
	// cpaths := strings.Split(cpath, string(os.PathSeparator))
	// wp, _ := filepath.Abs(".")
	// for i, p := range cpaths {
	// 	if p == "." {
	// 		cpaths[i] = wp
	// 	}
	// }
	// cpath = strings.Join(cpaths, "/")
	// cpath, _ = homedir.Expand(cpath)
	if strings.Contains(cpath, "..") {
		cpath = filepath.Join(cpath, ".")
	}
	apath := cpath
	apath, err := filepath.Abs(cpath)
	if err != nil {
		apath += " " + err.Error()
	}
	// if !filepath.IsAbs(cpath) {
	// 	apath, err := filepath.Abs(cpath)
	// 	if err != nil {
	// 		apath += " " + err.Error()
	// 	}
	// }
	paw.SetLoggerFieldsOrder([]string{"path", "home", "clean", "abs"})
	lg.WithFields(logrus.Fields{
		"path":  path,
		"home":  hpath,
		"clean": cpath,
		"abs":   apath,
	}).Debug()
	return apath
}

var reSimpPath = regexp.MustCompile("/+")

func simplifyPath(path string) string {
	dirs := reSimpPath.Split(path, -1)
	res := []string{}
	for _, v := range dirs {
		p := string(v)
		if p == ".." {
			if len(res) > 0 {
				res = res[:len(res)-1]
			}
		} else if p != "." && p != "" {
			res = append(res, p)
		}
	}
	str := strings.Join(res, "/")
	return "/" + str
}
