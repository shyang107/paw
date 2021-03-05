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
	grouping       vfs.Group
	isViewList     bool
	isViewLevel    bool
	isViewListTree bool
	isViewTree     bool
	isViewTable    bool
	isViewClassify bool
	isViewX        bool
	isViewGroup    bool
	isViewGroupR   bool
	isViewNoFiles  bool
	isViewNoDirs   bool
	// Depth
	depth          int
	isDepthRecurse bool
	// ByField (sort)
	byField       vfs.SortKey
	isSortNo      bool
	isSortReverse bool
	sortByField   string
	isSortByName  bool //default name
	isSortBySize  bool
	isSortByMTime bool
	// Skiper
	skips *vfs.SkipConds
	// Fields
	viewFields    vfs.ViewField
	hasINode      bool
	hasPermission bool
	hasHDLinks    bool
	hasSize       bool
	hasBlocks     bool
	hasUser       bool
	hasGroup      bool
	hasMTime      bool
	hasATime      bool
	hasCTime      bool
	hasGit        bool
	hasMd5        bool
}

var (
	opt    = new(option)
	vfsOpt = new(vfs.VFSOption)
	// vfsOpt = &VFSOption{
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
	// -------------------------------------------
	// ByField (sort)
	fg_isSortNo = &cli.BoolFlag{
		Name:        "no-sort",
		Aliases:     []string{"N"},
		Value:       false,
		Usage:       "not sort by name in increasing order (single key)",
		Destination: &opt.isSortNo,
	}
	fg_isSortReverse = &cli.BoolFlag{
		Name:        "reverse",
		Aliases:     []string{"r"},
		Value:       false,
		Usage:       "sort in decreasing order, default sort by name",
		Destination: &opt.isSortReverse,
	}
	fg_sortByField = &cli.StringFlag{
		Name:        "sort",
		Aliases:     []string{"sf"},
		Value:       "",
		Usage:       "which single `field` to sort by. (case insensitive,field: inode, links, blocks, size, mtime (ot modified), atime (or accessed), ctime (or created), name, lname (lower name, default); «field»[r|R]: reverse sort)",
		Destination: &opt.sortByField,
	}
	fg_isSortByName = &cli.BoolFlag{
		Name:        "sort-by-name",
		Aliases:     []string{"sn"},
		Value:       false,
		Usage:       "sort by name in increasing order (single key)",
		Destination: &opt.isSortByName,
	}
	fg_isSortBySize = &cli.BoolFlag{
		Name:        "sort-by-size",
		Aliases:     []string{"sz"},
		Value:       false,
		Usage:       "sort by size in increasing order (single key)",
		Destination: &opt.isSortBySize,
	}
	fg_isSortByMTime = &cli.BoolFlag{
		Name:        "sort-by-mtime",
		Aliases:     []string{"sm"},
		Value:       false,
		Usage:       "sort by modified time in increasing order (single key)",
		Destination: &opt.isSortByMTime,
	}
	// -------------------------------------------
	// Fields
)
