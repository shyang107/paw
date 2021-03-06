package main

import (
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
		lg.WithField("rootPath", opt.rootPath).Info()
	case 1:
		lg.WithField("arg", c.Args().Get(0)).Trace("no argument" + paw.Caller(1))
		path, err := filepath.Abs(c.Args().Get(0))
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
			lg.WithField("rootPath", opt.rootPath).Info()
		} else {
			if opt.paths == nil {
				opt.paths = make([]string, 0)
			}
			opt.paths = append(opt.paths, path)
			lg.WithField("paths", opt.paths).Info()
		}
	default: // > 1
		lg.WithField("arg", c.Args()).Trace("multi-arguments" + paw.Caller(1))
		if opt.paths == nil {
			opt.paths = make([]string, 0)
		}
		for i := 0; i < c.NArg(); i++ {
			lg.WithField("args", c.Args().Get(i)).Info()
			path, err := filepath.Abs(c.Args().Get(i))
			if err != nil {
				paw.Error.Println(err)
				continue
			}
			opt.paths = append(opt.paths, path)
			lg.WithField("paths", path).Info()
		}
	}
}
