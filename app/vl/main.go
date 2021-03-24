package main

import (
	"embed"
	// _ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/cnested"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

const (
// version = "0.0.2"
// releaseDate = "2021-03-08"
)

var (
	lg = paw.Logger

	//go:embed assets/*
	cfs embed.FS

	app *cli.App
	ca  *appConfig

	// version     string
	// appName     string
	// releaseTime time.Time
	// authorName  string
	// authorEmail string

	cInfoPrefix  string
	cWarnPrefix  string
	cErrorPrefix string

	// version     string
	// releaseDate string
	// app         *cli.App
	// appName     = "vl"
	// lg          = paw.Logger
	// releaseTime = cast.ToTime(releaseDate)
	// authorName  = "Shuhhua Yang"
	// authorEmail = "shyang107@gmail.com"

	// cInfoPrefix  = paw.Cinfo.Sprintf("[INFO]")
	// cWarnPrefix  = paw.Cwarn.Sprintf("[WARN]")
	// cErrorPrefix = paw.Cwarn.Sprintf("[ERRO]")
	// cInfoPrefix  = paw.Cinfo.Sprintf("[%s][INFO]", appName)
	// cWarnPrefix  = paw.Cwarn.Sprintf("[%s][WARN]", appName)
	// cErrorPrefix = paw.Cwarn.Sprintf("[%s][ERRO]", appName)
)

type appConfig struct {
	Appname string
	Version string `yaml:"version"`
	// Releasestr string `yaml:"release"`
	Release time.Time
	Author  struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	}
	Log struct {
		Pinfo  string `yaml:"pinfo"`
		Pwarn  string `yaml:"pwarn"`
		Perror string `yaml:"perror"`
	}
}

func readConfig() {
	ca = new(appConfig)
	buf, _ := cfs.ReadFile("assets/config.yaml")
	err := yaml.Unmarshal(buf, ca)
	if err != nil {
		fatal(err.Error())
	}
	appName, err := os.Executable()
	if err != nil || len(appName) == 0 {
		appName = os.Args[0]
	}
	appName = filepath.Base(appName)
	ca.Appname = appName

	// spew.Dump(ca)

	cInfoPrefix = paw.Cinfo.Sprint(ca.Log.Pinfo)
	cWarnPrefix = paw.Cwarn.Sprint(ca.Log.Pwarn)
	cErrorPrefix = paw.Cwarn.Sprint(ca.Log.Perror)
}

func init() {
	lg.SetLevel(logrus.WarnLevel)
	readConfig()

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
		Name:      ca.Appname,
		Usage:     "list directory (excluding hidden items) in color view.",
		UsageText: ca.Appname + " command [command options] [arguments...]",
		ArgsUsage: "[path]",
		Version:   ca.Version, //version,
		// Compiled : time.Now(),
		Compiled: ca.Release,
		Authors: []*cli.Author{
			{
				Name:  ca.Author.Name,
				Email: ca.Author.Email,
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
			fg_isInfo, fg_isDebug, fg_isTrace, fg_isDump,
			//  ViewType
			fg_isViewList, fg_isViewLevel, fg_isViewListTree, fg_isViewTree, fg_isViewTable, fg_isViewClassify,
			fg_isViewX, fg_isViewGroup, fg_isViewGroupR,
			fg_isViewNoDirs, fg_isViewNoFiles,
			// Depth
			fg_Depth, fg_IsFindRecurse, fg_isForceRecurse,
			// ByField (sort)
			fg_isSortNo, fg_isSortReverse, fg_sortByField, fg_isSortByName,
			fg_isSortByINode, fg_isSortBySize, fg_isSortByHDLinks, fg_isSortByBlocks,
			fg_isSortByUser, fg_isSortByGroup,
			fg_isSortByMTime, fg_isSortByATime, fg_isSortByCTime,
			fg_isSortByMd5,
			// SkipConds
			fg_isNoSkip, fg_reIncludePattern, fg_reExcludePattern,
			fg_withNoPrefix, fg_withNoSufix, fg_psDelimiter,
			// ViewFields
			fg_hasAll, fg_hasAllNoGit, fg_hasAllNoMd5, fg_hasAllNoGitMd5,
			fg_hasBasicPSUGMN,
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
	// info()
	// warning()
	// stderr()
	// fatal()
}

var (
	traceLogo = cnested.Logos[logrus.TraceLevel]
	debugLogo = cnested.Logos[logrus.DebugLevel]
	infoLogo  = cnested.Logos[logrus.InfoLevel]
	warnLogo  = cnested.Logos[logrus.WarnLevel]
	errorLogo = cnested.Logos[logrus.ErrorLevel]
	fatalLogo = cnested.Logos[logrus.FatalLevel]
)

func info(args ...interface{}) {
	if lg.IsLevelEnabled(logrus.InfoLevel) {
		fmt.Fprint(os.Stderr, infoLogo, " ")
		paw.Info.Print(args...)
		// fmt.Fprint(os.Stderr, msg)
	}
}

func infof(f string, args ...interface{}) {
	if lg.IsLevelEnabled(logrus.InfoLevel) {
		fmt.Fprint(os.Stderr, infoLogo, " ")
		paw.Info.Printf(f, args...)
	}
}

func stderr(args ...interface{}) {
	fmt.Fprint(os.Stderr, errorLogo, " ")
	paw.Error.Print(args...)
}

func stderrf(f string, args ...interface{}) {
	fmt.Fprint(os.Stderr, errorLogo, " ")
	paw.Error.Printf(f, args...)
}

func fatal(args ...interface{}) {
	// stderr(args...)
	fmt.Fprint(os.Stderr, fatalLogo, " ")
	paw.Error.Print(args...)
	os.Exit(1)
}
func fatalf(f string, args ...interface{}) {
	// stderrf(f, args...)
	fmt.Fprint(os.Stderr, fatalLogo, " ")
	paw.Error.Printf(f, args...)
	os.Exit(1)
}

func warning(args ...interface{}) {
	if lg.IsLevelEnabled(logrus.WarnLevel) {
		fmt.Fprint(os.Stderr, warnLogo, " ")
		paw.Warning.Print(args...)
	}
}
func warningf(f string, args ...interface{}) {
	if lg.IsLevelEnabled(logrus.WarnLevel) {
		fmt.Fprint(os.Stderr, warnLogo, " ")
		paw.Warning.Printf(f, args...)
	}
}
