package main

import (
	"fmt"

	"github.com/shyang107/paw"
)

func exFileLineCount() {
	paw.Logger.Info("testFileLineCount")
	lc, err := paw.FileLineCount("../README.md")
	if err != nil {
		paw.Logger.Error(err)
	}
	fmt.Println("FileLineCount:", lc)

}
