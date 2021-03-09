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
		Name:        "sortby",
		Aliases:     []string{"f"},
		Value:       "",
		Usage:       "which single `field` to sort by. (case insensitive,field: inode, links, blocks, size, mtime (ot modified), atime (or accessed), ctime (or created), name, lname (lower name, default); «field»[r|R]: reverse sort)",
		Destination: &opt.sortByField,
	}
	fg_isSortByName = &cli.BoolFlag{
		Name:        "sort-by-name",
		Aliases:     []string{"n"},
		Value:       false,
		Usage:       "sort by name in increasing order (single key)",
		Destination: &opt.isSortByName,
	}
	fg_isSortBySize = &cli.BoolFlag{
		Name:        "sort-by-size",
		Aliases:     []string{"z"},
		Value:       false,
		Usage:       "sort by size in increasing order (single key)",
		Destination: &opt.isSortBySize,
	}
	fg_isSortByMTime = &cli.BoolFlag{
		Name:        "sort-by-mtime",
		Aliases:     []string{"m"},
		Value:       false,
		Usage:       "sort by modified time in increasing order (single key)",
		Destination: &opt.isSortByMTime,
	}

	cmd_ByField = &cli.Command{
		Name:    "sort",
		Aliases: []string{"S"},
		Usage:   "sort by fields directly with default option",
		Flags: []cli.Flag{
			// ByField (sort)
			fg_isSortNo, fg_isSortReverse, fg_sortByField, fg_isSortByName, fg_isSortBySize, fg_isSortByMTime,
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
	if opt.isSortBySize &&
		!opt.isSortByMTime &&
		!opt.isSortByName {
		sflag = "size"
	} else if !opt.isSortBySize &&
		opt.isSortByMTime &&
		!opt.isSortByName {
		sflag = "mtime"
	} else if !opt.isSortBySize &&
		!opt.isSortByMTime &&
		opt.isSortByName {
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
	// lg.WithFields(logrus.Fields{
	// 	"byField": opt.byField,
	// }).Info()
	// info(paw.MesageFieldAndValueC("Sort", opt.byField, logrus.InfoLevel, paw.Cnop, nil))
}
