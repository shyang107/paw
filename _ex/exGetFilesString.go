package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/_junk"
)

func exGetFilesString() {
	paw.Logger.Info("")
	sourceFolder := "../"
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
		Fields:    []string{"No.", "File"},
		LenFields: []int{5, 72},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft},
		Padding:   "# ",
	}
	tp.Prepare(os.Stdout)
	tp.SetBeforeMessage(hsb.String())
	tp.PrintSart()

	files, err := _junk.GetFilesFuncString("../", isRecursive,
		func(f _junk.File) bool {
			return !(len(f.FileName) == 0 || strings.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
		})
	if err != nil {
		paw.Logger.Error(err)
	}

	for i, f := range files {
		path := strings.TrimPrefix(f, sourceFolder)
		j := i + 1
		tp.PrintRow(j, path)
		if j%5 == 0 {
			tp.PrintMiddleSepLine()
		}
	}
	tp.PrintEnd()
}
