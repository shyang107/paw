package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func main() {
	// root := `/Users/shyang/go/src/github.com/shyang107/paw/`
	// root := `/dev`
	var (
		root string
		opt  = vfs.NewVFSOption()
	)
	switch len(os.Args) {
	case 2:
		root = os.Args[1]
	case 3:
		root = os.Args[1]
		opt.Depth = cast.ToInt(os.Args[2])
	case 4:
		root = os.Args[1]
		opt.Depth = (cast.ToInt(os.Args[2]))
		if strings.ToLower(os.Args[3]) == "-v" {
			paw.Logger.SetLevel(logrus.TraceLevel)
		}
	default:
		root = "."
		opt.Depth = 0
	}

	// reSkip := vfs.NewSkipFuncRe("not *.go", `.go$`, func(de vfs.DirEntryX, r *regexp.Regexp) bool {
	// 	name := strings.TrimSpace(de.Name())
	// 	if !r.MatchString(name) || de.IsDir() {
	// 		return false
	// 	}
	// 	return true
	// })

	// reSkip := vfs.NewSkipFuncRe("get *.go", `.go$`, func(de vfs.DirEntryX, r *regexp.Regexp) bool {
	// 	name := strings.TrimSpace(de.Name())
	// 	if r.MatchString(name) || de.IsDir() {
	// 		return false
	// 	}
	// 	return true
	// })

	skipcond := vfs.NewSkipConds().Add(vfs.DefaultSkiper)
	// skipcond.AddToSkipSuffix("go")
	// skipcond.AddToSkipNames("pd")
	// skipcond.AddToSkipPrefix("make", "read")

	// vfs.SkipSuffix.Add("go")
	// spew.Dump(vfs.SkipSuffix)

	// skipcond := vfs.NewSkipConds().Add(vfs.DefaultSkip).Add(reSkip)
	vfields := vfs.DefaultViewField | vfs.ViewFieldGit //| vfs.ViewFieldMd5
	vopt := &vfs.VFSOption{
		Depth:          opt.Depth,
		IsForceRecurse: false,
		// Grouping: vfs.GroupedR, //vfs.GroupNone
		ByField:    vfs.SortByNone,
		Skips:      skipcond,
		ViewFields: vfields,
		// ViewType:   vfs.ViewList,
		// ViewType:   vfs.ViewListX,
		// ViewType: vfs.ViewLevel, //vfs.ViewLevel.NoDirs(),
		// ViewType:   vfs.ViewLevelX,
		// ViewType:   vfs.ViewTable,
		// ViewType:   vfs.ViewTableX,
		ViewType: vfs.ViewListTree,
		// ViewType:   vfs.ViewListTreeX,
		// ViewType: vfs.ViewTree,
		// ViewType:   vfs.ViewTreeX,
		// ViewType: vfs.ViewClassify,
	}

	fs := vfs.NewVFS(root, vopt)
	// fs.AddSkipFuncs(reSkip)
	// fs.AddSkipFuncs(vfs.SkipFile)
	fs.BuildFS()
	// fs.View(os.Stdout)
	viewTypes := []vfs.ViewType{
		// vfs.ViewList,
		// vfs.ViewListNoDirs,
		// vfs.ViewListNoFiles,
		// vfs.ViewListX,
		// vfs.ViewListXNoDirs,
		// vfs.ViewListXNoFiles,
		// vfs.ViewLevel,
		// vfs.ViewLevelNoDirs,
		// vfs.ViewLevelNoFiles,
		// vfs.ViewLevelX,
		// vfs.ViewLevelXNoDirs,
		// vfs.ViewLevelXNoFiles,
		// vfs.ViewTable,
		// vfs.ViewTableNoDirs,
		// vfs.ViewTableNoFiles,
		// vfs.ViewTableX,
		// vfs.ViewTableXNoDirs,
		// vfs.ViewTableXNoFiles,
		vfs.ViewListTree,
		// vfs.ViewListTreeX,
		// vfs.ViewTree,
		// vfs.ViewTreeX,
		// vfs.ViewClassify,
		// vfs.ViewClassifyNoDirs,
		// vfs.ViewClassifyNoFiles,
	}

	for _, v := range viewTypes {
		paw.Logger.Infoln(v)
		fs.SetViewType(v)
		ss := strings.Split(vopt.String(), "\n")
		for _, v := range ss {
			paw.Logger.Debugf(v)
		}
		fs.View(os.Stdout)
		// fmt.Println()
	}
	// fs.SetViewType(vfs.ViewLevel)
	// fs.View(os.Stdout)

}

func test() {
	lg.SetLevel(logrus.InfoLevel)
	root, _ := filepath.Abs("../..")
	lg.WithField("root", root).Info()

	opt := vfs.NewVFSOption()
	opt.Depth = 1
	opt.IsForceRecurse = true
	opt.ViewFields = vfs.DefaultViewFieldAllNoMd5
	opt.ViewType = vfs.ViewClassify

	fs := vfs.NewVFS(root, opt)
	fs.BuildFS()
	rd := fs.RootDir()

	lg.SetLevel(logrus.DebugLevel)
	lg.WithFields(logrus.Fields{
		"Depth":          opt.Depth,
		"IsForceRecurse": opt.IsForceRecurse,
	}).Debug()
	curlevel := len(strings.Split(rd.RelPath(), "/"))
	paw.Logger.WithFields(logrus.Fields{
		"1_name":           rd.Name(),
		"3_RelPath":        rd.RelPath(),
		"4_curlevel":       curlevel,
		"IsRelPathNotScan": opt.IsRelPathNotScan(rd.RelPath()),
		"IsRelPathNotView": opt.IsRelPathNotView(rd.RelPath()),
	}).Debug()
	for _, rp := range rd.RelPaths() {
		curlevel := len(strings.Split(rp, "/"))
		isscan := cast.ToString(!opt.IsRelPathNotScan(rp))
		if !opt.IsRelPathNotScan(rp) {
			isscan = paw.Cwarn.Sprint(!opt.IsRelPathNotScan(rp))
		} else {
			isscan = paw.CEven.Sprint(!opt.IsRelPathNotScan(rp))
		}
		isview := cast.ToString(!opt.IsRelPathNotView(rp))
		if !opt.IsRelPathNotView(rp) {
			isview = paw.Cwarn.Sprint(!opt.IsRelPathNotView(rp))
		} else {
			isview = paw.CEven.Sprint(!opt.IsRelPathNotView(rp))
		}
		fv := paw.MesageFieldAndValue("scan", isscan, logrus.InfoLevel)
		fv += paw.MesageFieldAndValue("view", isview, logrus.InfoLevel)
		infof("%v; level= %d; %q", fv, curlevel, rp)
		if opt.IsRelPathNotView(rp) {
			fmt.Println()
			continue
		}
		fmt.Print(">>>")
		lg.WithFields(logrus.Fields{
			"rp":      rp,
			"1 level": curlevel,
			"3 depth": opt.Depth,
			"2 l>d":   curlevel > opt.Depth,
		}).Debug()
		fmt.Print(">>>")
		paw.Logger.WithFields(logrus.Fields{
			"2_RelPath":        rp,
			"3_curlevel":       curlevel,
			"IsRelPathNotScan": opt.IsRelPathNotScan(rp),
			"IsRelPathNotView": opt.IsRelPathNotView(rp),
		}).Debug()
		fmt.Println()
	}
	spew.Dump(rd.RelPaths())
	fs.View(os.Stdout)
}

var (
	lg           = paw.Logger
	cInfoPrefix  = paw.Cinfo.Sprintf("[INFO]")
	cWarnPrefix  = paw.Cwarn.Sprintf("[WARN]")
	cErrorPrefix = paw.Cwarn.Sprintf("[ERRO]")
)

func info(args ...interface{}) {
	if lg.IsLevelEnabled(logrus.InfoLevel) {
		paw.Info.Print(args...)
	}
}

func infof(f string, args ...interface{}) {
	if lg.IsLevelEnabled(logrus.InfoLevel) {
		paw.Info.Printf(f, args...)
	}
}

func stderr(args ...interface{}) {
	paw.Error.Print(args...)
}

func stderrf(f string, args ...interface{}) {
	paw.Error.Printf(f, args...)
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
		paw.Warning.Print(args...)
	}
}
func warningf(f string, args ...interface{}) {
	if lg.IsLevelEnabled(logrus.WarnLevel) {
		paw.Warning.Printf(f, args...)
	}
}
