package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
		// if !fs.ValidPath(arg) {
		// 	fatal(&fs.PathError{
		// 		Op:   "checkArgs",
		// 		Path: arg,
		// 		Err:  fs.ErrInvalid,
		// 	})
		// }
		path, err := filepath.Abs(arg)
		if err != nil {
			paw.Error.Println(&fs.PathError{
				Op:   "checkArgs",
				Path: arg,
				Err:  err,
			})
		}
		fi, err := os.Stat(path)
		if err != nil {
			paw.Error.Println(&fs.PathError{
				Op:   "checkArgs",
				Path: arg,
				Err:  err,
			})
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
			// if !fs.ValidPath(arg) {
			// 	warning(&fs.PathError{
			// 		Op:   "checkArgs",
			// 		Path: arg,
			// 		Err:  fs.ErrInvalid,
			// 	})
			// 	continue
			// }
			path, err := filepath.Abs(arg)
			if err != nil {
				paw.Error.Println(&fs.PathError{
					Op:   "checkArgs",
					Path: arg,
					Err:  err,
				})
				viewPaths_errors = append(viewPaths_errors, err)
				continue
			}
			opt.paths = append(opt.paths, path)
			lg.WithField("path", path).Trace()
		}
		if len(opt.paths) == 0 {
			fatal(&fs.PathError{
				Op:   "checkArgs",
				Path: strings.Join(c.Args().Slice(), ";"),
				Err:  fs.ErrInvalid,
			})
			// fatalf("there is no valid paths: %v", c.Args().Slice())
		}
		lg.WithField("paths", opt.paths).Trace()
	}
}
