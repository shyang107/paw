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
	version     = "0.0.2"
	releaseDate = "2021-03-08"
)

var (
	app         *cli.App
	appName     = "vl"
	lg          = paw.Logger
	releaseTime = cast.ToTime(releaseDate)
	authorName  = "Shuhhua Yang"
	authorEmail = "shyang107@gmail.com"

	cInfoPrefix  = paw.Cinfo.Sprintf("[INFO]")
	cWarnPrefix  = paw.Cwarn.Sprintf("[WARN]")
	cErrorPrefix = paw.Cwarn.Sprintf("[ERRO]")
	// cInfoPrefix  = paw.Cinfo.Sprintf("[%s][INFO]", appName)
	// cWarnPrefix  = paw.Cwarn.Sprintf("[%s][WARN]", appName)
	// cErrorPrefix = paw.Cwarn.Sprintf("[%s][ERRO]", appName)
)

func init() {
	lg.SetLevel(logrus.WarnLevel)
	appName, err := os.Executable()
	if err != nil || len(appName) == 0 {
		appName = os.Args[0]
	}
	appName = filepath.Base(appName)

	paw.GologInit(os.Stdout, os.Stderr, os.Stderr, false)

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print only the version",
	}
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s_%v\n",
			c.App.Name,
			paw.NewEXAColor("sb").Sprint(c.App.Name+c.App.Version),
			paw.NewEXAColor("da").Sprint(c.App.Compiled.Format("Jan 2, 2006")))
	}

	cmd_Version := &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print only the version",
		Action: func(c *cli.Context) error {
			cli.ShowVersion(c)
			return nil
		},
	}

	app = &cli.App{
		Name:      appName,
		Usage:     "list directory (excluding hidden items) in color view.",
		UsageText: appName + " command [command options] [arguments...]",
		ArgsUsage: "[path]",
		Version:   version,
		// Compiled : time.Now(),
		Compiled: releaseTime,
		Authors: []*cli.Author{
			{
				Name:  authorName,
				Email: authorEmail,
			},
		},
		UseShortOptionHandling: true,
		// EnableBashCompletion : true,
		Commands: []*cli.Command{
			cmd_Version,
			//  ViewType
			cmd_ViewType,
			// ByField (sort)
			cmd_ByField,
			// SkipConds
			cmd_SkipConds,
			// ViewFields
			cmd_ViewField,
		},

		Flags: []cli.Flag{
			// verbose
			fg_isInfo, fg_isDebug, fg_isTrace,
			//  ViewType
			fg_isViewList, fg_isViewLevel, fg_isViewListTree, fg_isViewTree, fg_isViewTable, fg_isViewClassify,
			fg_isViewX, fg_isViewGroup, fg_isViewGroupR,
			fg_isViewNoDirs, fg_isViewNoFiles,
			// Depth
			fg_Depth, fg_isDepthRecurse, fg_isDepthScanAllSub,
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
		},
		Action: appAction,
	}

}

func main() {
	// start := time.Now()

	err := app.Run(os.Args)
	if err != nil {
		fatalf("run '%s' failed, error:%v", app.Name, err)
	}

	// elapsedTime := time.Since(start)
	// fmt.Println()
	// fmt.Println("Total time for excution:", elapsedTime.String())
}

func info(args ...interface{}) {
	// paw.Info.Printf(programName + ": " + fmt.Sprintf(f, args...) + "\n")
	if lg.IsLevelEnabled(logrus.InfoLevel) {
		fmt.Fprintf(os.Stderr, "%s %v\n", cInfoPrefix, fmt.Sprint(args...))
	}
}
func infof(f string, args ...interface{}) {
	// paw.Info.Printf(programName + ": " + fmt.Sprintf(f, args...) + "\n")
	if lg.IsLevelEnabled(logrus.InfoLevel) {
		fmt.Fprintf(os.Stderr, "%s %v\n", cInfoPrefix, fmt.Sprintf(f, args...))
	}
}

func stderr(args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %v\n", cErrorPrefix, fmt.Sprint(args...))
}

func stderrf(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %v\n", cErrorPrefix, fmt.Sprintf(f, args...))
}

func fatal(args ...interface{}) {
	stderr(args...)
	os.Exit(1)
}
func fatalf(f string, args ...interface{}) {
	stderrf(f, args...)
	os.Exit(1)
}

func warning(args ...interface{}) {
	if lg.IsLevelEnabled(logrus.WarnLevel) {
		fmt.Fprintf(os.Stderr, "%s %v\n", cWarnPrefix, fmt.Sprint(args...))
	}
}
func warningf(f string, args ...interface{}) {
	if lg.IsLevelEnabled(logrus.WarnLevel) {
		fmt.Fprintf(os.Stderr, "%s %v\n", cWarnPrefix, fmt.Sprintf(f, args...))
	}
}
