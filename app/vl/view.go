package main

import (
	"os"

	"github.com/shyang107/paw/vfs"
)

func (opt *option) View() error {
	lg.Debug()

	// if opt.vopt == nil {
	// 	opt.vopt = vfs.NewVFSOption()
	// }

	lg.Debug(opt.vopt)
	fs := vfs.NewVFS(opt.rootPath, opt.vopt)
	fs.BuildFS()
	fs.View(os.Stdout)

	return nil
}
