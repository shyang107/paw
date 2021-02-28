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
	fs.View(os.Stdout, vfs.AllViewFieldsNoMd5, vfs.ViewLevel)
}
