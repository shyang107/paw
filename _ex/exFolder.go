package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/shyang107/paw"
)

func exFolder() {
	paw.Logger.Info("exFolder")
	path := "/aaa/bbb/ccc/example.xxx"
	fmt.Println("                            path:", path)
	file := paw.ConstructFile(path, "")
	fmt.Println("ConstructFile(path):")
	spew.Dump(file)
	sourceFolder := "/aaa/bbb/"
	fmt.Println("                    sourceFolder:", sourceFolder)
	subfolder := paw.GetSubfolder(file, sourceFolder)
	fmt.Println("GetSubfolder(file, sourceFolder):", subfolder)
	targetFolder := "ddd/"
	fmt.Println("                    targetFolder:", targetFolder)
	newFolder, _ := paw.GetNewFilePath(file, sourceFolder, targetFolder)
	fmt.Println("GetNewFilePath(file, sourceFolder, targetFolder):", newFolder)
}
