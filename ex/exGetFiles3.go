package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shyang107/paw"
)

func exGetFiles3() {
	paw.Logger.Info("exGetFiles3")
	sourceFolder := "../"
	sourceFolder, err := filepath.Abs(sourceFolder)
	if err != nil {
		paw.Logger.Error(err)
	}
	sourceFolder += "/"
	head := "\nGetFilesFunc: folder <- '../', isRecursive <- true\n"
	head += "  sourceFolder: " + sourceFolder + "\n"
	prefix := "."
	regexPattern := `\.git`
	re := regexp.MustCompile(regexPattern)
	head += "  Exculde:" + "\n"
	head += fmt.Sprintf("          prefix: %q\n", prefix)
	head += fmt.Sprintf("    regexPattern: %q", regexPattern)

	tp := &paw.TableFormat{
		Fields:    []string{"No.", "File"},
		LenFields: []int{5, 72},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft},
		Padding:   "  ",
	}
	tp.Prepare(os.Stdout)
	tp.SetBeforeMessage(head)
	tp.PrintSart()

	files, err := paw.GetFilesFunc("../", true, func(f paw.File) bool {
		return (len(f.FileName) == 0 || strings.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
	})
	if err != nil {
		paw.Logger.Error(err)
	}

	for i, f := range files {
		newPath, err := paw.GetNewFilePath(f, sourceFolder, "./")
		if err != nil {
			paw.Logger.Error(err)
		}
		rows := []interface{}{i + 1, newPath}
		tp.PrintRow(rows...)
	}
	tp.PrintEnd()
}
