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
	depth          int
	IsFindRecurse  bool
	isForceRecurse bool
	// ByField (sort)
	byField         vfs.SortKey
	isSortNo        bool
	isSortReverse   bool
	sortByField     string
	isSortByName    bool //default name
	isSortByINode   bool
	isSortBySize    bool
	isSortByBlocks  bool
	isSortByHDLinks bool
	isSortByUser    bool
	isSortByGroup   bool
	isSortByMTime   bool
	isSortByATime   bool
	isSortByCTime   bool
	isSortByMd5     bool
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
	hasBasicPSUGMN bool
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
)
