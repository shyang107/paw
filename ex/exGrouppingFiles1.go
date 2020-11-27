package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/shyang107/paw"
)

func exGrouppingFiles1() {
	paw.Logger.Info("")
	// sourceFolder := "../"
	sourceFolder := "/Users/shyang/go/src/rover/opcc/"
	isRecursive := true
	sourceFolder, err := filepath.Abs(sourceFolder)
	if err != nil {
		paw.Logger.Error(err)
	}
	sourceFolder += "/"
	hsb := strings.Builder{}
	hsb.WriteString("GetFilesFuncString:\n")
	hsb.WriteString("  sourceFolder: " + `"../" <- "` + sourceFolder + "\"\n")
	hsb.WriteString("   isRecursive: " + strconv.FormatBool(isRecursive) + "\n")
	prefix := "."
	regexPattern := `\.git`
	re := regexp.MustCompile(regexPattern)
	hsb.WriteString("  Exculde:" + "\n")
	hsb.WriteString(`          prefix: "` + prefix + `"` + "\n")
	hsb.WriteString(`    regexPattern: "` + regexPattern + `"`)

	tp := &paw.TableFormat{
		Fields:    []string{"No.", "File", "Sorted Files"},
		LenFields: []int{5, 40, 40},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft, paw.AlignLeft},
		Padding:   "# ",
	}
	tp.Prepare(os.Stdout)
	tp.SetBeforeMessage(hsb.String())
	tp.PrintSart()

	files, err := paw.GetFilesFunc(sourceFolder, isRecursive,
		func(f paw.File) bool {
			return len(f.FileName) == 0 || strings.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath)
		})
	if err != nil {
		paw.Logger.Error(err)
	}
	sfiles := make([]paw.File, len(files))
	copy(sfiles, files)
	paw.GrouppingFiles(sfiles)

	oFolder := sfiles[0].Folder
	for i, f := range files {
		path := strings.TrimPrefix(f.FullPath, sourceFolder)
		spath := strings.TrimPrefix(sfiles[i].FullPath, sourceFolder)
		j := i + 1
		// if j%5 == 0 {
		if oFolder != sfiles[i].Folder {
			oFolder = sfiles[i].Folder
			tp.PrintMiddleSepLine()
		}
		tp.PrintRow(j, path, spath)
	}
	tp.PrintEnd()
}
