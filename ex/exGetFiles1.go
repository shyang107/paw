package main

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/shyang107/paw"
)

func exGetFiles1() {
	paw.Logger.Info("exGetFiles1")
	paw.Logger.Info("GetFiles: folder <- '~/', isRecursive <- false")
	homepath, err := homedir.Dir()
	if err != nil {
		paw.Logger.Error(err)
	}
	files, err := paw.GetFiles(homepath, false)
	if err != nil {
		paw.Logger.Error(err)
	}
	for i, f := range files {
		fmt.Printf("%2d. %s\n", i+1, f.FullPath)
	}
}
