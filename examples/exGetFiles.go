package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/_junk"
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

	files, err := _junk.GetFilesFunc("../", true, func(f _junk.File) bool {
		return (len(f.FileName) == 0 || strings.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
	})
	if err != nil {
		paw.Logger.Error(err)
	}

	for i, f := range files {
		newPath, err := _junk.GetNewFilePath(f, sourceFolder, "./")
		if err != nil {
			paw.Logger.Error(err)
		}
		rows := []interface{}{i + 1, newPath}
		tp.PrintRow(rows...)
	}
	tp.PrintEnd()
}
func exGetFiles2() {
	paw.Logger.Info("exGetFiles2")
	paw.Logger.Info("GetFiles: folder <- '../', isRecursive <- true")
	sourceFolder := "../"
	fmt.Println("sourceFolder:", sourceFolder)
	files, err := _junk.GetFiles(sourceFolder, true)
	if err != nil {
		paw.Logger.Error(err)
	}
	for i, f := range files {
		fmt.Printf("%3d. %s\n", i, f.FullPath)
	}
	i := 0
	prefix := "."
	regexPattern := `\.git`
	re := regexp.MustCompile(regexPattern)
	fmt.Println("Exculde:")
	fmt.Printf("\t      prefix: %q\n", prefix)
	fmt.Printf("\tregexPattern: %q\n", regexPattern)
	for _, f := range files {
		if strings.HasPrefix(f.FileName, prefix) {
			continue
		} else if len(f.FileName) == 0 {
			continue
		} else if re.MatchString(f.FullPath) {
			continue
		}
		i++
		fmt.Printf("%3d. %s\n", i, f.FullPath)
	}
}
func exGetFiles1() {
	paw.Logger.Info("exGetFiles1")
	paw.Logger.Info("GetFiles: folder <- '~/', isRecursive <- false")
	homepath, err := homedir.Dir()
	if err != nil {
		paw.Logger.Error(err)
	}
	files, err := _junk.GetFiles(homepath, false)
	if err != nil {
		paw.Logger.Error(err)
	}
	for i, f := range files {
		fmt.Printf("%2d. %s\n", i+1, f.FullPath)
	}
}
