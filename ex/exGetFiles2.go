package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shyang107/paw"
)

func exGetFiles2() {
	paw.Logger.Info("exGetFiles2")
	paw.Logger.Info("GetFiles: folder <- '../', isRecursive <- true")
	sourceFolder := "../"
	fmt.Println("sourceFolder:", sourceFolder)
	files, err := paw.GetFiles(sourceFolder, true)
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
