package main

import (
	"os"
	"strings"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func main() {
	// root := `/Users/shyang/go/src/github.com/shyang107/paw/`
	// root := `/dev`
	var (
		root     string
		level    int
		loglevel = logrus.WarnLevel
	)
	switch len(os.Args) {
	case 2:
		root = os.Args[1]
	case 3:
		root = os.Args[1]
		level = cast.ToInt(os.Args[2])
	case 4:
		root = os.Args[1]
		level = cast.ToInt(os.Args[2])
		if strings.ToLower(os.Args[3]) == "-v" {
			loglevel = logrus.TraceLevel
		}
	default:
		root = "."
		level = 0
	}
	paw.Logger.SetLevel(loglevel)

	fs := vfs.NewVFSWith(root, level)

	// reSkip := vfs.NewSkipFuncRe("not *.go", `.go$`, func(de vfs.DirEntryX, r *regexp.Regexp) bool {
	// 	name := strings.TrimSpace(de.Name())
	// 	if !r.MatchString(name) || de.IsDir() {
	// 		return false
	// 	}
	// 	return true
	// })
	// fs.AddSkipFuncs(reSkip)
	// fs.AddSkipFuncs(vfs.SkipFile)
	fs.BuildFS()

	vfields := vfs.DefaultViewField //| vfs.ViewFieldMd5
	// fs.View(os.Stdout, vfields, vfs.ViewList)
	// fs.View(os.Stdout, vfields, vfs.ViewListX)
	fs.View(os.Stdout, vfields, vfs.ViewLevel)
	// fs.View(os.Stdout, vfields, vfs.ViewLevelX)
	// fs.View(os.Stdout, vfields, vfs.ViewTable)
	// fs.View(os.Stdout, vfields, vfs.ViewTableX)
	// fs.View(os.Stdout, vfields, vfs.ViewListTree)
	// fs.View(os.Stdout, vfields, vfs.ViewListTreeX)
	// fs.View(os.Stdout, vfields, vfs.ViewTree)
	// fs.View(os.Stdout, vfields, vfs.ViewTreeX)
	// fs.View(os.Stdout, vfields, vfs.ViewClassify)

}
