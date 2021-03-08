package main

import (
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	// -------------------------------------------
	// ViewFields
	fg_hasAll = &cli.BoolFlag{
		Name:        "allfields",
		Aliases:     []string{"af"},
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

	cmd_ViewField = &cli.Command{
		Name:    "field",
		Aliases: []string{"F"},
		Usage:   "set what fields to show directly with default option",
		Flags: []cli.Flag{
			// ViewFields
			fg_hasAll, fg_hasAllNoGit, fg_hasAllNoMd5, fg_hasAllNoGitMd5,
			fg_hasBasicPSUGN,
			fg_hasINode,
			fg_hasPermission,
			fg_hasHDLinks, fg_hasSize, fg_hasBlocks,
			fg_hasUser, fg_hasGroup,
			fg_hasMTime, fg_hasATime, fg_hasCTime,
			fg_hasGit, fg_hasMd5,
		},
		Subcommands: []*cli.Command{
			{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "show all fields",
				Action: func(c *cli.Context) error {
					opt.hasAll = true
					return appAction(c)
				},
			},
			{
				Name:    "allxgit",
				Aliases: []string{"xg"},
				Usage:   "show all fields, except git",
				Action: func(c *cli.Context) error {
					opt.hasAllNoGit = true
					return appAction(c)
				},
			},
			{
				Name:    "allxmd5",
				Aliases: []string{"x5"},
				Usage:   "show all fields, except md5",
				Action: func(c *cli.Context) error {
					opt.hasAllNoMd5 = true
					return appAction(c)
				},
			},
			{
				Name:    "allxgitmd5",
				Aliases: []string{"xg5"},
				Usage:   "show all fields, except git and md5",
				Action: func(c *cli.Context) error {
					opt.hasAllNoGitMd5 = true
					return appAction(c)
				},
			},
			{
				Name:    "basic",
				Aliases: []string{"6"},
				Usage:   "show basic fields: inode, permission, user, group, modified, and name (required field)",
				Action: func(c *cli.Context) error {
					opt.hasBasicPSUGN = true
					return appAction(c)
				},
			},
		},
		Action: appAction,
	}
)

func (opt *option) checkViewFields() {
	lg.Debug(paw.Caller(1))

	var (
		viewFields vfs.ViewField
		isOk       bool = false
		hasBasic        = opt.hasBasicPSUGN
	)

	if opt.hasAll {
		opt.viewFields = vfs.DefaultViewFieldAll
		goto END
	}
	if opt.hasAllNoMd5 {
		opt.viewFields = vfs.DefaultViewFieldAllNoMd5
		goto END
	}
	if opt.hasAllNoGit {
		opt.viewFields = vfs.DefaultViewFieldAllNoGit
		goto END
	}
	if opt.hasAllNoGitMd5 {
		opt.viewFields = vfs.DefaultViewFieldAllNoGitMd5
		goto END
	}

	if opt.hasINode {
		isOk = true
		viewFields |= vfs.ViewFieldINode
	}
	if hasBasic {
		isOk = true
		viewFields |= vfs.ViewFieldPermissions
	} else {
		if opt.hasPermission {
			isOk = true
			viewFields |= vfs.ViewFieldPermissions
		}
	}
	if opt.hasHDLinks {
		isOk = true
		viewFields |= vfs.ViewFieldLinks
	}
	if hasBasic {
		isOk = true
		viewFields |= vfs.ViewFieldSize
	} else {
		if opt.hasSize {
			isOk = true
			viewFields |= vfs.ViewFieldSize
		}
	}
	if opt.hasBlocks {
		isOk = true
		viewFields |= vfs.ViewFieldBlocks
	}
	if hasBasic {
		isOk = true
		viewFields |= vfs.ViewFieldUser
		viewFields |= vfs.ViewFieldGroup
		viewFields |= vfs.ViewFieldModified
	} else {
		if opt.hasUser {
			isOk = true
			viewFields |= vfs.ViewFieldUser
		}
		if opt.hasGroup {
			isOk = true
			viewFields |= vfs.ViewFieldGroup
		}
		if opt.hasMTime {
			isOk = true
			viewFields |= vfs.ViewFieldModified
		}
	}
	if opt.hasCTime {
		isOk = true
		viewFields |= vfs.ViewFieldCreated
	}
	if opt.hasATime {
		isOk = true
		viewFields |= vfs.ViewFieldAccessed
	}
	if opt.hasMd5 {
		isOk = true
		viewFields |= vfs.ViewFieldMd5
	}
	if opt.hasGit {
		isOk = true
		viewFields |= vfs.ViewFieldGit
	}

	viewFields |= vfs.ViewFieldName
	lg.WithFields(logrus.Fields{
		"isOk":       isOk,
		"viewFields": viewFields,
	}).Debug()

	if isOk {
		opt.viewFields = viewFields
		// if viewFields == vfs.ViewFieldName {
		// 	opt.viewFields = vfs.DefaultViewField
		// } else {
		// 	opt.viewFields = viewFields
		// }
	} else {
		opt.viewFields = vfs.DefaultViewField
	}
END:
	// lg.WithField("viewFields", opt.viewFields).Info()
	info(paw.MesageFieldAndValueC("View fields", opt.viewFields, logrus.InfoLevel, paw.Cnop, nil))

}
