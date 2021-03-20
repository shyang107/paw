package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/shyang107/paw"
	"github.com/urfave/cli"
)

func (opt *option) checkArgs(c *cli.Context) {
	lg.Debug()

	switch c.NArg() {
	case 0:
		lg.WithField("arg", c.Args().Get(0)).Trace("no argument" + paw.Caller(1))
		path, err := filepath.Abs(".")
		if err != nil {
			paw.Error.Println(err)
		}
		opt.rootPath = path
		info(paw.NewValuePair("Root", opt.rootPath))
	case 1:
		lg.WithField("arg", c.Args().Get(0)).Trace("no argument" + paw.Caller(1))
		arg := c.Args().Get(0)
		if fs.ValidPath(arg) {
			fatalf("%q is not a valid path!", arg)
		}
		path, err := filepath.Abs(arg)
		if err != nil {
			paw.Error.Println(err)
		}
		fi, err := os.Stat(path)
		if err != nil {
			paw.Error.Println(err)
			os.Exit(1)
		}
		if fi.IsDir() {
			opt.rootPath = path
			info(paw.NewValuePair("Root", opt.rootPath))
		} else {
			if opt.paths == nil {
				opt.paths = make([]string, 0)
			}
			opt.paths = append(opt.paths, path)
			info(paw.NewValuePair("Paths", opt.paths))
		}
	default: // > 1
		lg.WithField("arg", c.Args()).Trace("multi-arguments" + paw.Caller(1))
		if opt.paths == nil {
			opt.paths = make([]string, 0, c.NArg())
		}
		lg.WithField("args", c.Args()).Debug()
		for i := 0; i < c.NArg(); i++ {
			arg := c.Args().Get(i)
			if fs.ValidPath(arg) {
				warningf("%q is not a valid path!", arg)
				continue
			}
			path, err := filepath.Abs(arg)
			if err != nil {
				paw.Error.Println(err)
				viewPaths_errors = append(viewPaths_errors, err)
				continue
			}
			opt.paths = append(opt.paths, path)
			lg.WithField("path", path).Trace()
		}
		if len(opt.paths) == 0 {
			fatalf("there is no valid paths: %v", c.Args().Slice())
		}
		lg.WithField("paths", opt.paths).Trace()
	}
}
