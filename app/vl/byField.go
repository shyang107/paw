package main

import (
	"strings"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	// -------------------------------------------
	// ByField (sort)
	fg_isSortNo = &cli.BoolFlag{
		Name:        "notsort",
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
		Name:        "sortby",
		Aliases:     []string{"f"},
		Value:       "",
		Usage:       "which single `field` to sort by. (case insensitive,field: inode, links, blocks, size, mtime (ot modified), atime (or accessed), ctime (or created), name, lname (lower name, default); «field»[r|R]: reverse sort)",
		Destination: &opt.sortByField,
	}
	fg_isSortByName = &cli.BoolFlag{
		Name:        "byname",
		Aliases:     []string{"bn"},
		Value:       false,
		Usage:       "sort by name in increasing order (single key)",
		Destination: &opt.isSortByName,
	}
	fg_isSortByINode = &cli.BoolFlag{
		Name:        "byinode",
		Aliases:     []string{"bi"},
		Value:       false,
		Usage:       "sort by inode in increasing order (single key)",
		Destination: &opt.isSortByINode,
	}
	fg_isSortBySize = &cli.BoolFlag{
		Name:        "bysize",
		Aliases:     []string{"bz"},
		Value:       false,
		Usage:       "sort by size in increasing order (single key)",
		Destination: &opt.isSortBySize,
	}
	fg_isSortByHDLinks = &cli.BoolFlag{
		Name:        "bylinks",
		Aliases:     []string{"bl"},
		Value:       false,
		Usage:       "sort by links in increasing order (single key)",
		Destination: &opt.isSortByHDLinks,
	}
	fg_isSortByBlocks = &cli.BoolFlag{
		Name:        "byblocks",
		Aliases:     []string{"bb"},
		Value:       false,
		Usage:       "sort by blocks in increasing order (single key)",
		Destination: &opt.isSortByBlocks,
	}

	fg_isSortByUser = &cli.BoolFlag{
		Name:        "byuser",
		Aliases:     []string{"su"},
		Value:       false,
		Usage:       "sort by user in increasing order (single key)",
		Destination: &opt.isSortByUser,
	}
	fg_isSortByGroup = &cli.BoolFlag{
		Name:        "bygroup",
		Aliases:     []string{"bg"},
		Value:       false,
		Usage:       "sort by group in increasing order (single key)",
		Destination: &opt.isSortByGroup,
	}
	fg_isSortByMTime = &cli.BoolFlag{
		Name:        "bymtime",
		Aliases:     []string{"bm"},
		Value:       false,
		Usage:       "sort by modified time in increasing order (single key)",
		Destination: &opt.isSortByMTime,
	}
	fg_isSortByATime = &cli.BoolFlag{
		Name:        "byatime",
		Aliases:     []string{"ba"},
		Value:       false,
		Usage:       "sort by accessed time in increasing order (single key)",
		Destination: &opt.isSortByATime,
	}
	fg_isSortByCTime = &cli.BoolFlag{
		Name:        "byctime",
		Aliases:     []string{"bc"},
		Value:       false,
		Usage:       "sort by created time in increasing order (single key)",
		Destination: &opt.isSortByCTime,
	}
	fg_isSortByMd5 = &cli.BoolFlag{
		Name:        "bymd5",
		Aliases:     []string{"b5"},
		Value:       false,
		Usage:       "sort by md5 string in increasing order (single key)",
		Destination: &opt.isSortByMd5,
	}

	cmd_ByField = &cli.Command{
		Name:    "sort",
		Aliases: []string{"S"},
		Usage:   "sort by fields directly with default option",
		Flags: []cli.Flag{
			// ByField (sort)
			fg_isSortNo, fg_isSortReverse, fg_sortByField, fg_isSortByName,
			fg_isSortByINode, fg_isSortBySize, fg_isSortByHDLinks, fg_isSortByBlocks,
			fg_isSortByUser, fg_isSortByGroup,
			fg_isSortByMTime, fg_isSortByATime, fg_isSortByCTime,
			fg_isSortByMd5,
		},
		Subcommands: []*cli.Command{
			{
				Name:    "field",
				Aliases: []string{"f"},
				Usage:   "which single `field` to sort by. (case insensitive,field: inode, links, blocks, size, mtime (ot modified), atime (or accessed), ctime (or created), name, lname (lower name, default); «field»[r|R]: reverse sort)",
				Action: func(c *cli.Context) error {
					opt.sortByField = c.Args().First()
					return appAction(c)
				},
			},
			{
				Name:    "reverse",
				Aliases: []string{"r"},
				Usage:   "sort in decreasing order, default sort by lower name",
				Action: func(c *cli.Context) error {
					opt.isSortReverse = true
					return appAction(c)
				},
			},
			{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "sort by name in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortByName = true
					return appAction(c)
				},
			},
			{
				Name:    "inode",
				Aliases: []string{"i"},
				Usage:   "sort by inode in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortByINode = true
					return appAction(c)
				},
			},
			{
				Name:    "links",
				Aliases: []string{"l"},
				Usage:   "sort by links in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortByINode = true
					return appAction(c)
				},
			},
			{
				Name:    "size",
				Aliases: []string{"z"},
				Usage:   "sort by size in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortBySize = true
					return appAction(c)
				},
			},
			{
				Name:    "blocks",
				Aliases: []string{"b"},
				Usage:   "sort by blocks in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortByBlocks = true
					return appAction(c)
				},
			},
			{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "sort by user in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortByUser = true
					return appAction(c)
				},
			},
			{
				Name:    "group",
				Aliases: []string{"g"},
				Usage:   "sort by group in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortByINode = true
					return appAction(c)
				},
			},
			{
				Name:    "mtime",
				Aliases: []string{"m"},
				Usage:   "sort by modified time in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortByMTime = true
					return appAction(c)
				},
			},
			{
				Name:    "atime",
				Aliases: []string{"a"},
				Usage:   "sort by accessed time in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortByATime = true
					return appAction(c)
				},
			},
			{
				Name:    "ctime",
				Aliases: []string{"c"},
				Usage:   "sort by created time in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortByCTime = true
					return appAction(c)
				},
			},
			{
				Name:    "md5",
				Aliases: []string{"5"},
				Usage:   "sort by md5 string in increasing order (single key)",
				Flags: []cli.Flag{
					fg_isSortReverse,
				},
				Action: func(c *cli.Context) error {
					opt.isSortByMd5 = true
					return appAction(c)
				},
			},
		},
		Action: appAction,
	}
)

func (opt *option) checkByField() {
	lg.Debug(paw.Caller(1))
	var (
		sflag string
		ok    bool
	)
	if opt.isSortNo {
		opt.byField = vfs.SortByNone
		goto END
	}

	// opt.byField = vfs.SortByLowerName
	sflag = strings.ToLower(opt.sortByField)
	if len(sflag) == 0 {
		sflag = "lname"
	}
	if opt.isSortByINode {
		sflag = "inode"
	}
	if opt.isSortByHDLinks {
		sflag = "links"
	}
	if opt.isSortBySize {
		sflag = "size"
	}
	if opt.isSortByBlocks {
		sflag = "blocks"
	}
	if opt.isSortByMTime {
		sflag = "mtime"
	}
	if opt.isSortByATime {
		sflag = "atime"
	}
	if opt.isSortByCTime {
		sflag = "ctime"
	}
	if opt.isSortByMd5 {
		sflag = "md5"
	}
	if opt.isSortByName {
		sflag = "name"
	}
	if opt.isSortReverse {
		sflag += "r"
	}

	lg.WithFields(logrus.Fields{
		"sflag":   sflag,
		"SortKey": vfs.SortShortNameKeys[sflag],
	}).Trace()

	if opt.byField, ok = vfs.SortShortNameKeys[sflag]; !ok {
		opt.byField = vfs.SortByLowerName
	}

END:
	info(paw.NewValuePair("Sort", opt.byField))
}
