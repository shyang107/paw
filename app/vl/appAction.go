package main

import (
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var appAction cli.ActionFunc = func(c *cli.Context) error {
	lg.Debug()

	lg.SetLevel(logrus.WarnLevel)

	if opt.isInfo {
		lg.SetLevel(logrus.InfoLevel)
	}

	if opt.isDebug {
		lg.SetLevel(logrus.DebugLevel)
	}

	if opt.isTrace {
		lg.SetLevel(logrus.TraceLevel)
	}

	opt.vopt = vfs.NewVFSOption()

	opt.checkArgs(c)

	// ViewType (during view)
	opt.checkViewType()

	//ByField (sort))
	opt.checkByField()

	// SkipConds (during BuildVFS)
	opt.checkSkips()

	// ViewFields
	opt.checkViewFields()

	// Setuo vfs.VFSOption
	opt.setVFSOption()

	// View
	if len(opt.paths) < 1 {
		err := opt.view()
		if err != nil {
			stderrf("view: %s", err.Error())
		}
	} else {
		err := opt.viewPaths()
		if err != nil {
			stderrf("view: %s", err.Error())
		}
	}

	return nil
}
