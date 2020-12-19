package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/_junk"
)

func exFolder() {
	paw.Logger.Info("exFolder")
	path := "/aaa/bbb/ccc/example.xxx"
	fmt.Println("                            path:", path)
	file := _junk.ConstructFile(path, "")
	fmt.Println("ConstructFile(path):")
	spew.Dump(file)
	sourceFolder := "/aaa/bbb/"
	fmt.Println("                    sourceFolder:", sourceFolder)
	subfolder := _junk.GetSubfolder(file, sourceFolder)
	fmt.Println("GetSubfolder(file, sourceFolder):", subfolder)
	targetFolder := "ddd/"
	fmt.Println("                    targetFolder:", targetFolder)
	newFolder, _ := _junk.GetNewFilePath(file, sourceFolder, targetFolder)
	fmt.Println("GetNewFilePath(file, sourceFolder, targetFolder):", newFolder)
}
