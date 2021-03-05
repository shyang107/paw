package main

import (
	"os"
	"strings"

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
		opt.Depth = cast.ToInt(os.Args[2])
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

	skipcond := vfs.NewSkipConds().Add(vfs.DefaultSkiper).Add(vfs.DefaultSkiper)
	// skipcond.AddToSkipSuffix("go")
	// skipcond.AddToSkipNames("pd")
	// skipcond.AddToSkipPrefix("make", "read")

	// vfs.SkipSuffix.Add("go")
	// spew.Dump(vfs.SkipSuffix)

	// skipcond := vfs.NewSkipConds().Add(vfs.DefaultSkip).Add(reSkip)
	vfields := vfs.DefaultViewField | vfs.ViewFieldGit //| vfs.ViewFieldMd5
	vopt := &vfs.VFSOption{
		Depth: opt.Depth,
		// Grouping: vfs.GroupedR, //vfs.GroupNone
		// ByField:  vfs.SortBySizeR,
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
}
