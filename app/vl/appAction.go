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

	// ViewType
	opt.checkViewType()

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
