package main

import (
	"fmt"

	"github.com/shyang107/paw"
)

func exLogger() {
	paw.Logger.Info("exLogger")
	paw.Logger.Debug("exLogger")
	paw.Logger.Warn("exLogger")
	paw.Logger.Trace("exLogger")
	fmt.Println("  GetDotDir()", paw.GetDotDir())
	fmt.Println("GetCurrPath()", paw.GetCurrPath())
	fmt.Println("  GetAppDir()", paw.GetAppDir())

}
