package main

import (
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/urfave/cli"
)

var (
	// -------------------------------------------
	// ViewType
	fg_isViewList = &cli.BoolFlag{
		Name:        "list",
		Aliases:     []string{"li"},
		Value:       true,
		Usage:       "print out in list view",
		Destination: &opt.isViewList,
	}
	fg_isViewLevel = &cli.BoolFlag{
		Name:        "level",
		Aliases:     []string{"le"},
		Value:       false,
		Usage:       "print out in the level view",
		Destination: &opt.isViewLevel,
	}
	fg_isViewListTree = &cli.BoolFlag{
		Name:        "listtree",
		Aliases:     []string{"lt"},
		Value:       false,
		Usage:       "print out in the view of combining list and tree",
		Destination: &opt.isViewListTree,
	}
	fg_isViewTree = &cli.BoolFlag{
		Name:        "tree",
		Aliases:     []string{"tr"},
		Value:       false,
		Usage:       "print out in the tree view",
		Destination: &opt.isViewTree,
	}
	fg_isViewTable = &cli.BoolFlag{
		Name:        "table",
		Aliases:     []string{"ta"},
		Value:       false,
		Usage:       "print out in the table view",
		Destination: &opt.isViewTable,
	}
	fg_isViewClassify = &cli.BoolFlag{
		Name:        "classify",
		Aliases:     []string{"cl"},
		Value:       false,
		Usage:       "display type indicator by file names",
		Destination: &opt.isViewClassify,
	}
	fg_isViewX = &cli.BoolFlag{
		Name:        "extended",
		Aliases:     []string{"@"},
		Value:       false,
		Usage:       "list each file's extended attributes and sizes",
		Destination: &opt.isViewX,
	}
	fg_isViewGroup = &cli.BoolFlag{
		Name:        "grouped",
		Aliases:     []string{"G"},
		Value:       false,
		Usage:       "group files and directories separately",
		Destination: &opt.isViewGroup,
	}
	fg_isViewGroupR = &cli.BoolFlag{
		Name:        "groupedr",
		Aliases:     []string{"H"},
		Value:       false,
		Usage:       "group files and directories separately",
		Destination: &opt.isViewGroupR,
	}
	fg_isViewNoDirs = &cli.BoolFlag{
		Name:        "nodirs",
		Aliases:     []string{"nd"},
		Value:       false,
		Usage:       "show all files but not directories, has high priority than --just-dirs",
		Destination: &opt.isViewNoDirs,
	}
	fg_isViewNoFiles = &cli.BoolFlag{
		Name:        "nofiles",
		Aliases:     []string{"nf"},
		Value:       false,
		Usage:       "show all dirs but not files",
		Destination: &opt.isViewNoFiles,
	}
	// -------------------------------------------
	// Depth
	fg_Depth = &cli.IntFlag{
		Name:        "depth",
		Aliases:     []string{"d"},
		Value:       0,
		Usage:       "set `value` of depth show the files (dirs) under root",
		Destination: &opt.depth,
	}
	fg_isDepthRecurse = &cli.BoolFlag{
		Name:        "recurse",
		Aliases:     []string{"R"},
		Value:       false,
		Usage:       "recurse into directories (equivalent to --depth=-1)",
		Destination: &opt.isDepthRecurse,
	}
	fg_isDepthScanAllSub = &cli.BoolFlag{
		Name:        "scan-all-sub",
		Aliases:     []string{"S"},
		Value:       false,
		Usage:       "anyway, definitely recurse all sub-directories of root",
		Destination: &opt.isDepthScanAllSub,
	}

	cmd_ViewType = &cli.Command{
		Name:    "view",
		Aliases: []string{"V"},
		Usage:   "set ViewType, can be tailing with sub-commands: sort, skip, field",
		Flags: []cli.Flag{
			//  ViewType
			fg_isViewList, fg_isViewLevel, fg_isViewListTree, fg_isViewTree, fg_isViewTable, fg_isViewClassify,
			fg_isViewX, fg_isViewGroup, fg_isViewGroupR,
			fg_isViewNoDirs, fg_isViewNoFiles,
			// Depth
			fg_Depth, fg_isDepthRecurse, fg_isDepthScanAllSub,
		},
		Subcommands: []*cli.Command{
			{
				Name:    "depth",
				Aliases: []string{"d"},
				Usage:   "set `depth` to show the files (dirs) under root",
				Action: func(c *cli.Context) error {
					opt.depth = cast.ToInt(c.Args().First())
					return appAction(c)
				},
			},
			cmd_ByField,
			cmd_SkipConds,
			cmd_ViewField,
		},
		Action: appAction,
	}
)

func (opt *option) checkViewType() {
	lg.Debug(paw.Caller(1))
	// 1. cehck basic ViewType
	if opt.isViewListTree {
		if opt.depth == 0 {
			opt.depth = -1
		}
		opt.viewType = vfs.ViewListTree
	} else if opt.isViewTree {
		if opt.depth == 0 {
			opt.depth = -1
		}
		opt.viewType = vfs.ViewTree
	} else if opt.isViewTable {
		opt.viewType = vfs.ViewTable
	} else if opt.isViewLevel {
		opt.viewType = vfs.ViewLevel
	} else if opt.isViewClassify {
		opt.viewType = vfs.ViewClassify
	} else if opt.isViewList {
		opt.viewType = vfs.ViewList
	}
	lg.WithField("viewType", opt.viewType).Trace()

	// 2. cehck Extended view
	if opt.isViewX {
		hasX = true
		lg.WithField("isViewX", opt.isViewX).Trace()
		if opt.viewType&vfs.ViewClassify == 0 {
			opt.viewType |= vfs.ViewExtended
		}
		lg.WithField("> viewType", opt.viewType).Trace()
	}

	// check Grouping
	lg.WithFields(logrus.Fields{
		"isViewGroup":  opt.isViewGroup,
		"isViewGroupR": opt.isViewGroupR,
	}).Trace()
	if opt.isViewGroup && !opt.isViewGroupR {
		opt.grouping = vfs.Grouped
	} else if !opt.isViewGroup && opt.isViewGroupR {
		opt.grouping = vfs.GroupedR
	} else {
		opt.grouping = vfs.GroupNone
	}

	lg.WithFields(logrus.Fields{
		"isViewNoDirs":  opt.isViewNoDirs,
		"isViewNoFiles": opt.isViewNoFiles,
	}).Trace()
	if opt.isViewNoDirs && !opt.isViewNoFiles {
		switch opt.viewType {
		case vfs.ViewList, vfs.ViewLevel, vfs.ViewTable, vfs.ViewClassify,
			vfs.ViewListX, vfs.ViewLevelX, vfs.ViewTableX:
			opt.viewType |= vfs.ViewNoDirs
		}
		lg.WithField("> viewType", opt.viewType).Trace()
	}

	if !opt.isViewNoDirs && opt.isViewNoFiles {
		switch opt.viewType {
		case vfs.ViewList, vfs.ViewLevel, vfs.ViewTable, vfs.ViewClassify,
			vfs.ViewListX, vfs.ViewLevelX, vfs.ViewTableX:
			opt.viewType |= vfs.ViewNoFiles
		}
		lg.WithField("> viewType", opt.viewType).Trace()
	}
	// lg.Debugf("viewType: %v [%d]; ViewLevelXNoFiles: %v [%d]", opt.viewType, opt.viewType, vfs.ViewLevelXNoFiles, vfs.ViewLevelXNoFiles)

	// Depth
	lg.WithField("depth", opt.depth).Trace()

	lg.WithField("isDepthRecurse", opt.isDepthRecurse).Trace()
	if opt.isDepthRecurse {
		opt.depth = -1
	}

	lg.WithField("isDepthRecurse", opt.isDepthScanAllSub).Trace()

	// lg.WithFields(logrus.Fields{
	// 	"viewType": opt.viewType,
	// 	"grouping": opt.grouping,
	// 	"depth":    opt.depth,
	// }).Info()
	info(paw.MesageFieldAndValueC("View type:", opt.viewType, logrus.InfoLevel, paw.Cnop, nil))
	info(paw.MesageFieldAndValueC("Groupe", opt.grouping, logrus.InfoLevel, paw.Cnop, nil))
	info(paw.MesageFieldAndValueC("Searching depth", opt.depth, logrus.InfoLevel, paw.Cnop, nil))
}
