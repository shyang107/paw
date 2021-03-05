package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var appAction cli.ActionFunc = func(c *cli.Context) error {

	if opt.isVerbose {
		// pdopt.EnableTrace(opt.isVerbose)
		lg.SetLevel(logrus.TraceLevel)
	} else {
		lg.SetLevel(logrus.WarnLevel)
	}

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
	opt.setVFSOption(vfsOpt)

	// pattern
	// pdopt.Ignore = getPatternflag(opt).Ignore(opt)

	// pdopt.FieldFlag = getFieldFlag(opt)
	// pdopt.SortOpt = getSortOption(opt)
	// pdopt.FiltOpt = getFiltOption(opt)

	// err, _ := filetree.PrintDir(os.Stdout, opt.path, opt.isGrouped, pdopt, "")
	// if err != nil {
	// 	fatal("get file list from %q failed, error: %v", opt.path, err)
	// }

	return nil
}
