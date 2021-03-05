package main

import (
	"github.com/shyang107/paw/vfs"
	"github.com/urfave/cli"
)

type option struct {
	isVerbose bool
	// VFS
	rootPath string
	paths    []string
	// ViewType
	viewType       vfs.ViewType
	isViewList     bool
	isViewLevel    bool
	isViewListTree bool
	isViewTree     bool
	isViewTable    bool
	isViewClassify bool
	isViewX        bool
	isViewNoFiles  bool
	isViewNoDirs   bool
	// Depth
	depth          int
	isDepthRecurse bool
}

var (
	opt  = new(option)
	vopt = vfs.NewVFSOption()
	// vopt = &VFSOption{
	// 	Depth:      0,
	// 	Grouping:   GroupNone,
	// 	ByField:    SortByLowerName,
	// 	Skips:      NewSkipConds().Add(DefaultSkiper),
	// 	ViewFields: DefaultViewField,
	// 	ViewType:   ViewList,
	// }
	// -------------------------------------------
	// Verbose
	fg_isVerbose = &cli.BoolFlag{
		Name:        "verbose",
		Aliases:     []string{"V"},
		Value:       false,
		Usage:       "show verbose message",
		Destination: &opt.isVerbose,
	}
	// -------------------------------------------
	// ViewType
	fg_isViewList = &cli.BoolFlag{
		Name:        "list",
		Aliases:     []string{"l"},
		Value:       true,
		Usage:       "print out in list view",
		Destination: &opt.isViewList,
	}
	fg_isViewLevel = &cli.BoolFlag{
		Name:        "level",
		Aliases:     []string{"L"},
		Value:       false,
		Usage:       "print out in the level view",
		Destination: &opt.isViewLevel,
	}
	fg_isViewListTree = &cli.BoolFlag{
		Name:        "listtree",
		Aliases:     []string{"t"},
		Value:       false,
		Usage:       "print out in the view of combining list and tree",
		Destination: &opt.isViewListTree,
	}
	fg_isViewTree = &cli.BoolFlag{
		Name:        "tree",
		Aliases:     []string{"T"},
		Value:       false,
		Usage:       "print out in the tree view",
		Destination: &opt.isViewTree,
	}
	fg_isViewTable = &cli.BoolFlag{
		Name:        "table",
		Aliases:     []string{"b"},
		Value:       false,
		Usage:       "print out in the table view",
		Destination: &opt.isViewTable,
	}
	fg_isViewClassify = &cli.BoolFlag{
		Name:        "classify",
		Aliases:     []string{"f"},
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
	fg_isViewNoDirs = &cli.BoolFlag{
		Name:        "nodirs",
		Aliases:     []string{"D"},
		Value:       false,
		Usage:       "show all files but not directories, has high priority than --just-dirs",
		Destination: &opt.isViewNoDirs,
	}
	fg_isViewNoFiles = &cli.BoolFlag{
		Name:        "nofiles",
		Aliases:     []string{"F"},
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
		Usage:       "print out in the level view",
		Destination: &opt.depth,
	}
	fg_isDepthRecurse = &cli.BoolFlag{
		Name:        "recurse",
		Aliases:     []string{"R"},
		Value:       false,
		Usage:       "recurse into directories (equivalent to --depth=-1)",
		Destination: &opt.isDepthRecurse,
	}
)
