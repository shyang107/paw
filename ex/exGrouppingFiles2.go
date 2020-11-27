package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/cast"
)

func exGrouppingFiles2() {
	paw.Logger.Info("")
	sourceFolder := "../"
	// sourceFolder, _ := homedir.Expand("~/Downloads/")
	// sourceFolder := "/Users/shyang/go/src/rover/opcc/"
	isRecursive := true
	sourceFolder, err := filepath.Abs(sourceFolder)
	if err != nil {
		paw.Logger.Error(err)
	}
	sourceFolder += "/"
	hsb := strings.Builder{}
	hsb.WriteString("GetFilesFuncString:\n")
	hsb.WriteString("- sourceFolder:	\"" + sourceFolder + "\"\n")
	hsb.WriteString("- isRecursive:	" + strconv.FormatBool(isRecursive) + "\n")
	prefix := "."
	regexPattern := `\.git|\$RECYCLE\.BIN|desktop\.ini`
	re := regexp.MustCompile(regexPattern)
	hsb.WriteString("- Excluding conditions:" + "\n")
	hsb.WriteString(`	- prefix:	"` + prefix + `"` + "\n")
	hsb.WriteString(`	- regexPattern:	"` + regexPattern + `"`)

	tp := &paw.TableFormat{
		Fields:    []string{"No.", "Sorted Files"},
		LenFields: []int{5, 100},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft},
		Padding:   "# ",
	}
	tp.Prepare(os.Stdout)
	tp.SetBeforeMessage(hsb.String())
	tp.PrintSart()

	files, err := paw.GetFilesFunc(sourceFolder, isRecursive,
		func(f paw.File) bool {
			return (len(f.FileName) == 0 || paw.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
		})
	if err != nil {
		paw.Logger.Error(err)
	}

	paw.GrouppingFiles(files)

	oFolder := files[0].Folder
	gcount := 1
	j := 0
	for i, f := range files {
		if oFolder != f.Folder {
			oFolder = f.Folder
			tp.PrintRow("", "Sum: "+cast.ToString(j)+" files.")
			tp.PrintMiddleSepLine()
			j = 1
			gcount++
		} else {
			j++
		}
		if j == 1 {
			// folder := paw.TrimPrefix(f.Folder, sourceFolder)
			// tp.PrintRow("", fmt.Sprintf("folder %d: %q", gcount, "./"+folder))
			tp.PrintRow("", fmt.Sprintf("folder %d: %q", gcount, f.ShortFolder))
		}
		tp.PrintRow(j, f.File)
		// path := paw.TrimPrefix(f.FullPath, sourceFolder)
		// tp.PrintRow(j, path)
		if i == len(files)-1 {
			tp.PrintRow("", "Sum: "+cast.ToString(j)+" files.")
		}
	}
	tp.SetAfterMessage("Total: " + cast.ToString(gcount) + " subfolders and " + cast.ToString(len(files)) + " files. ")
	tp.PrintEnd()
}
