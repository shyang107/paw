package main

import (
	"github.com/shyang107/paw/vfs"
	"github.com/urfave/cli"
)

type option struct {
	isTrace bool
	isDebug bool
	isInfo  bool
	// VFS
	rootPath string
	paths    []string
	vopt     *vfs.VFSOption
	// vfsOpt = &VFSOption{
	// 	Depth:      0,
	// 	Grouping:   GroupNone,
	// 	ByField:    SortByLowerName,
	// 	Skips:      NewSkipConds().Add(DefaultSkiper),
	// 	ViewFields: DefaultViewField,
	// 	ViewType:   ViewList,
	// }
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
	depth             int
	isDepthRecurse    bool
	isDepthScanAllSub bool
	// ByField (sort)
	byField       vfs.SortKey
	isSortNo      bool
	isSortReverse bool
	sortByField   string
	isSortByName  bool //default name
	isSortBySize  bool
	isSortByMTime bool
	// SkipConds
	skips            *vfs.SkipConds
	isNoSkip         bool
	reIncludePattern string
	reExcludePattern string
	psDelimiter      string
	withNoPrefix     string
	withNoSufix      string
	// ViewFields
	viewFields     vfs.ViewField
	hasAll         bool
	hasAllNoMd5    bool
	hasAllNoGit    bool
	hasAllNoGitMd5 bool
	hasBasicPSUGN  bool
	hasINode       bool
	hasPermission  bool
	hasHDLinks     bool
	hasSize        bool
	hasBlocks      bool
	hasUser        bool
	hasGroup       bool
	hasMTime       bool
	hasATime       bool
	hasCTime       bool
	hasGit         bool
	hasMd5         bool
}

var (
	opt  = new(option)
	hasX = false
	// -------------------------------------------
	// Verbose
	fg_isInfo = &cli.BoolFlag{
		Name:        "info",
		Aliases:     []string{},
		Value:       false,
		Usage:       "info",
		Destination: &opt.isInfo,
	}
	fg_isDebug = &cli.BoolFlag{
		Name:        "debug",
		Aliases:     []string{},
		Value:       false,
		Usage:       "debug mode",
		Destination: &opt.isDebug,
	}
	fg_isTrace = &cli.BoolFlag{
		Name:        "trace",
		Aliases:     []string{},
		Value:       false,
		Usage:       "trace mode",
		Destination: &opt.isTrace,
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
	fg_isDepthScanAllSub = &cli.BoolFlag{
		Name:        "scan-all-sub",
		Aliases:     []string{"S"},
		Value:       false,
		Usage:       "anyway, definitely recurse all sub-directories of root",
		Destination: &opt.isDepthScanAllSub,
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
		Aliases:     []string{"fd"},
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
	// SkipConds
	fg_isNoSkip = &cli.BoolFlag{
		Name:        "all",
		Aliases:     []string{"a"},
		Value:       false,
		Usage:       "show all files including hidden files",
		Destination: &opt.isNoSkip,
	}
	fg_reIncludePattern = &cli.StringFlag{
		Name:        "include",
		Aliases:     []string{"ri"},
		Value:       "",
		Usage:       "use regex to find files (not dirs) with matching `pattern`",
		Destination: &opt.reIncludePattern,
	}
	fg_reExcludePattern = &cli.StringFlag{
		Name:        "exclude",
		Aliases:     []string{"rx"},
		Value:       "",
		Usage:       "use regex to find files (not dirs) without matching `pattern`",
		Destination: &opt.reExcludePattern,
	}
	fg_withNoPrefix = &cli.StringFlag{
		Name:        "no-prefix",
		Aliases:     []string{"np"},
		Value:       "",
		Usage:       "skips name of files (dirs) with `prefix`; mutli-prefixs: prefix1,prefix2,...",
		Destination: &opt.withNoPrefix,
	}
	fg_withNoSufix = &cli.StringFlag{
		Name:        "no-suffix",
		Aliases:     []string{"ns"},
		Value:       "",
		Usage:       "skips name of files (dirs) with `suffix`; mutli-suffixs: suffix1,suffix2,...",
		Destination: &opt.withNoSufix,
	}
	fg_psDelimiter = &cli.StringFlag{
		Name:        "delimiter",
		Aliases:     []string{"psd"},
		Value:       ",",
		Usage:       "set `delimiter` needed int mutli-[prefixs|suffixs]",
		Destination: &opt.psDelimiter,
	}
	// -------------------------------------------
	// ViewFields
	fg_hasAll = &cli.BoolFlag{
		Name:        "allfields",
		Aliases:     []string{"x"},
		Value:       false,
		Usage:       "list each file's all fields",
		Destination: &opt.hasAll,
	}
	fg_hasAllNoGit = &cli.BoolFlag{
		Name:        "xgit",
		Aliases:     []string{"xg"},
		Value:       false,
		Usage:       "list each file's all fields, except git",
		Destination: &opt.hasAllNoGit,
	}
	fg_hasAllNoMd5 = &cli.BoolFlag{
		Name:        "xmd5",
		Aliases:     []string{"x5"},
		Value:       false,
		Usage:       "list each file's all fields, except md5",
		Destination: &opt.hasAllNoMd5,
	}
	fg_hasAllNoGitMd5 = &cli.BoolFlag{
		Name:        "xgitmd5",
		Aliases:     []string{"xg5"},
		Value:       false,
		Usage:       "list each file's all fields, except git and md5",
		Destination: &opt.hasAllNoGitMd5,
	}
	fg_hasBasicPSUGN = &cli.BoolFlag{
		Name:        "basic",
		Aliases:     []string{"6"},
		Value:       false,
		Usage:       "list each file's basic fields: inode, permission, user, group, modified, and name (required field)",
		Destination: &opt.hasBasicPSUGN,
	}
	fg_hasINode = &cli.BoolFlag{
		Name:        "inode",
		Aliases:     []string{"I"},
		Value:       false,
		Usage:       "list each file's inode number",
		Destination: &opt.hasINode,
	}
	fg_hasPermission = &cli.BoolFlag{
		Name:        "permissions",
		Aliases:     []string{"P"},
		Value:       false,
		Usage:       "list each file's permissions",
		Destination: &opt.hasPermission,
	}
	fg_hasHDLinks = &cli.BoolFlag{
		Name:        "links",
		Aliases:     []string{"K"},
		Value:       false,
		Usage:       "list each file's number of hard links",
		Destination: &opt.hasHDLinks,
	}
	fg_hasSize = &cli.BoolFlag{
		Name:        "size",
		Aliases:     []string{"Z"},
		Value:       false,
		Usage:       "list each file's size",
		Destination: &opt.hasSize,
	}
	fg_hasBlocks = &cli.BoolFlag{
		Name:        "blocks",
		Aliases:     []string{"B"},
		Value:       false,
		Usage:       "show number of file system blocks",
		Destination: &opt.hasBlocks,
	}
	fg_hasUser = &cli.BoolFlag{
		Name:        "user",
		Aliases:     []string{"s"},
		Value:       false,
		Usage:       "show user's name",
		Destination: &opt.hasUser,
	}
	fg_hasGroup = &cli.BoolFlag{
		Name:        "group",
		Aliases:     []string{"p"},
		Value:       false,
		Usage:       "show user's group name",
		Destination: &opt.hasGroup,
	}
	fg_hasGit = &cli.BoolFlag{
		Name:        "git",
		Aliases:     []string{"g"},
		Value:       false,
		Usage:       " list each file's Git status, if tracked or ignored",
		Destination: &opt.hasGit,
	}
	fg_hasMd5 = &cli.BoolFlag{
		Name:        "md5",
		Aliases:     []string{"5"},
		Value:       false,
		Usage:       " list each file's md5 field",
		Destination: &opt.hasMd5,
	}

	fg_hasMTime = &cli.BoolFlag{
		Name:        "modified",
		Aliases:     []string{"M"},
		Value:       false,
		Usage:       "use the modified timestamp field",
		Destination: &opt.hasMTime,
	}
	fg_hasATime = &cli.BoolFlag{
		Name:        "accessed",
		Aliases:     []string{"A"},
		Value:       false,
		Usage:       "use the accessed timestamp field",
		Destination: &opt.hasATime,
	}
	fg_hasCTime = &cli.BoolFlag{
		Name:        "created",
		Aliases:     []string{"C"},
		Value:       false,
		Usage:       "use the created timestamp field",
		Destination: &opt.hasCTime,
	}
)
