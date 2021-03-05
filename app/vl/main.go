package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shyang107/paw"
	"github.com/urfave/cli"
)

const (
	version     = "0.0.0.1"
	releaseDate = "2021-03-5"
)

var (
	app         = cli.NewApp()
	programName string
	lg          = paw.Logger
)

func _runFirst() {
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
	app.Compiled = time.Now()
	// app.Compiled = cast.ToTime(releaseDate)
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
			paw.NewEXAColor("sb").Sprint("gl"+c.App.Version),
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
		fg_isVerbose,
		//  ViewType
		fg_isViewList, fg_isViewLevel, fg_isViewListTree, fg_isViewTree, fg_isViewTable, fg_isViewClassify,
		fg_isViewX,
		fg_isViewNoDirs, fg_isViewNoFiles,
		// Depth
		fg_Depth, fg_isDepthRecurse,

		// &allFilesFlag, &includePatternFlag, &excludePatternFlag,
		// &isNoEmptyDirsFlag, &isJustDirsFlag, &isJustFilesFlag,
		// &isFieldINodeFlag, &isFieldLinksFlag,
		// // &isFieldPermissionsFlag,
		// // &isFieldSizeFlag,
		// &isFieldBlocksFlag,
		// // &isFieldUserFlag, &isFieldGroupFlag,
		// &isModifiedFlag, &isAccessedFlag, &isCreatedFlag,
		// &isFieldMd5Flag,
		// &isFieldGitFlag,
		// &isExtendedFlag,
		// &isNoSortFlag, &isReverseFlag, &sortByFieldFlag, &isSortByNameFlag, &isSortBySizeFlag, &isSortByMTimeFlag,
		// &isGroupedFlag,
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
	if opt.isVerbose {
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
	if opt.isVerbose {
		paw.Warning.Printf(programName + ": " + fmt.Sprintf(f, args...) + "\n")
		// stderr(f, args...)
	}
}
