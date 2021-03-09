package main

import (
	"regexp"
	"strings"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
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
		Aliases:     []string{"dl"},
		Value:       ",",
		Usage:       "set `delimiter` needed int mutli-[prefixs|suffixs]",
		Destination: &opt.psDelimiter,
	}

	cmd_SkipConds = &cli.Command{
		Name:    "skip",
		Aliases: []string{"K"},
		Usage:   "skip some specific (name of) files (dirs) directly with default option",
		Flags: []cli.Flag{
			// SkipConds
			fg_isNoSkip, fg_reIncludePattern, fg_reExcludePattern,
			fg_withNoPrefix, fg_withNoSufix, fg_psDelimiter,
		},
		Subcommands: []*cli.Command{
			{
				Name:    "reinclude",
				Aliases: []string{"ri"},
				Usage:   "use regex to find files (not dirs) with matching `pattern`",
				Action: func(c *cli.Context) error {
					opt.reIncludePattern = c.Args().First()
					return appAction(c)
				},
			},
			{
				Name:    "reexclude",
				Aliases: []string{"rx"},
				Usage:   "use regex to find files (not dirs) without matching `pattern`",
				Action: func(c *cli.Context) error {
					opt.reExcludePattern = c.Args().First()
					return appAction(c)
				},
			},
			{
				Name:    "noprefix",
				Aliases: []string{"np", "nopf"},
				Usage:   "skips name of files (dirs) with `prefix`; mutli-prefixs: prefix1,prefix2,...",
				Action: func(c *cli.Context) error {
					opt.withNoPrefix = c.Args().First()
					return appAction(c)
				},
			},
			{
				Name:    "nosuffix",
				Aliases: []string{"ns", "nosf"},
				Usage:   "skips name of files (dirs) with `suffix`; mutli-suffixs: suffix1,suffix2,...",
				Action: func(c *cli.Context) error {
					opt.withNoPrefix = c.Args().First()
					return appAction(c)
				},
			},
			{
				Name:    "delimiter",
				Aliases: []string{"dl", "ps"},
				Usage:   "set `delimiter` needed int mutli-[prefixs|suffixs]",
				Action: func(c *cli.Context) error {
					opt.psDelimiter = c.Args().First()
					return appAction(c)
				},
			},
		},
		Action: appAction,
	}
)

func (opt *option) checkSkips() {
	lg.Debug(paw.Caller(1))

	opt.skips = vfs.NewSkipConds()

	// All files
	if !opt.isNoSkip {
		opt.skips = opt.skips.Add(vfs.DefaultSkiper)
	}

	// reInclude
	if len(opt.reIncludePattern) > 0 {
		pattern := opt.reIncludePattern
		if lg.IsLevelEnabled(logrus.TraceLevel) {
			re := regexp.MustCompile(pattern)
			lg.WithField("ri.pattern", re.String()).Trace()
		}
		reSkiper := vfs.NewSkipFuncRe("«re-include»", pattern, func(de vfs.DirEntryX, re *regexp.Regexp) bool {
			name := strings.TrimSpace(de.Name())
			if re.MatchString(name) || de.IsDir() {
				return false
			}
			return true
		})
		opt.skips.Add(reSkiper)
	}
	// reExclude
	if len(opt.reExcludePattern) > 0 {
		pattern := opt.reExcludePattern
		if lg.IsLevelEnabled(logrus.TraceLevel) {
			re := regexp.MustCompile(pattern)
			lg.WithField("rx.pattern", re.String()).Trace()
		}
		reSkiper := vfs.NewSkipFuncRe("«re-exclude»", pattern, func(de vfs.DirEntryX, re *regexp.Regexp) bool {
			name := strings.TrimSpace(de.Name())
			if !re.MatchString(name) || de.IsDir() {
				return false
			}
			return true
		})
		opt.skips.Add(reSkiper)
	}

	// prefix
	if len(opt.withNoPrefix) > 0 {
		prefixs := strings.Split(opt.withNoPrefix, opt.psDelimiter)
		for _, prefix := range prefixs {
			opt.skips.AddToSkipPrefix(prefix)
		}
		lg.WithFields(logrus.Fields{
			"preifx":    prefixs,
			"delimiter": opt.psDelimiter,
		}).Trace()
	}
	// suffix
	if len(opt.withNoSufix) > 0 {
		suffixs := strings.Split(opt.withNoSufix, opt.psDelimiter)
		for _, suffix := range suffixs {
			opt.skips.AddToSkipSuffix(suffix)
		}
		lg.WithFields(logrus.Fields{
			"suffix":    suffixs,
			"delimiter": opt.psDelimiter,
		}).Trace()
	}
	info(paw.NewValuePair("Skiper", opt.skips))
	// paw.Logger.WithField("skips", opt.skips).Info()
	// info(paw.MesageFieldAndValueC("Skiper", opt.skips, logrus.InfoLevel, paw.Cnop, nil))
}
