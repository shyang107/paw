package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/urfave/cli"
)

const (
	version     = "0.0.1"
	releaseDate = "2021-3-6"
)

var (
	app         = cli.NewApp()
	programName string
	lg          = paw.Logger
)

func _runFirst() {
	lg.SetLevel(logrus.WarnLevel)
	programName, err := os.Executable()
	if err != nil {
		programName = os.Args[0]
	}
	programName = filepath.Base(programName)

	paw.GologInit(os.Stdout, os.Stderr, os.Stderr, false)

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print only the version",
	}

	app.Name = "vl"
	app.Usage = "list directory (excluding hidden items) in color view."
	app.UsageText = app.Name + " command [command options] [arguments...]"
	app.Version = version
	// app.Compiled = time.Now()
	app.Compiled = cast.ToTime(releaseDate)
	app.Authors = []*cli.Author{
		{
			Name:  "Shuhhua Yang",
			Email: "shyang107@gmail.com",
		},
	}
	app.ArgsUsage = "[path]"

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s @ %v\n",
			c.App.Name,
			paw.NewEXAColor("sb").Sprint(app.Name+c.App.Version),
			paw.NewEXAColor("da").Sprint(c.App.Compiled.Format("Jan 2, 2006")))
	}

	// app.EnableBashCompletion = true

	app.UseShortOptionHandling = true

	app.Commands = []*cli.Command{
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "print only the version",
			Action: func(c *cli.Context) error {
				cli.ShowVersion(c)
				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		// verbose
		fg_isInfo, fg_isDebug, fg_isTrace,
		//  ViewType
		fg_isViewList, fg_isViewLevel, fg_isViewListTree, fg_isViewTree, fg_isViewTable, fg_isViewClassify,
		fg_isViewX, fg_isViewGroup, fg_isViewGroupR,
		fg_isViewNoDirs, fg_isViewNoFiles,
		// Depth
		fg_Depth, fg_isDepthRecurse,
		// ByField (sort)
		fg_isSortNo, fg_isSortReverse, fg_sortByField, fg_isSortByName, fg_isSortBySize, fg_isSortByMTime,
		// SkipConds
		fg_isNoSkip, fg_reIncludePattern, fg_reExcludePattern,
		fg_withNoPrefix, fg_withNoSufix, fg_psDelimiter,
		// ViewFields
		fg_hasAll, fg_hasAllNoGit, fg_hasAllNoMd5, fg_hasAllNoGitMd5,
		fg_hasBasicPSUGN,
		fg_hasINode,
		fg_hasPermission,
		fg_hasHDLinks, fg_hasSize, fg_hasBlocks,
		fg_hasUser, fg_hasGroup,
		fg_hasMTime, fg_hasATime, fg_hasCTime,
		fg_hasGit, fg_hasMd5,
	}

	app.Action = appAction
}

func main() {
	_runFirst()
	// start := time.Now()

	err := app.Run(os.Args)
	if err != nil {
		fatal("run '%s' failed, error:%v", app.Name, err)
	}

	// elapsedTime := time.Since(start)
	// fmt.Println()
	// fmt.Println("Total time for excution:", elapsedTime.String())
}

func info(f string, args ...interface{}) {
	if lg.IsLevelEnabled(logrus.InfoLevel) {
		paw.Info.Printf(programName + ": " + fmt.Sprintf(f, args...) + "\n")
	}
	// fmt.Fprintf(os.Stderr, programName+": "+fmt.Sprintf(f, args...)+"\n")
}

func stderr(f string, args ...interface{}) {
	paw.Error.Printf(programName + ": " + fmt.Sprintf(f, args...) + "\n")
	// fmt.Fprintf(os.Stderr, programName+": "+fmt.Sprintf(f, args...)+"\n")
}

func fatal(f string, args ...interface{}) {
	stderr(f, args...)
	os.Exit(1)
}

func warning(f string, args ...interface{}) {
	if lg.IsLevelEnabled(logrus.WarnLevel) {
		paw.Warning.Printf(programName + ": " + fmt.Sprintf(f, args...) + "\n")
		// stderr(f, args...)
	}
}
