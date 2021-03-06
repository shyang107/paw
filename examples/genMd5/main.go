package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func main() {
	// path, _ := filepath.Abs("/Users/shyang/go/src/github.com/shyang107/paw/path.go")
	path, _ := filepath.Abs(`/Volumes/T7 Touch/Action/tube-x/Pink Paradise Paris - Striptease & Table Dance By Ike - Xvid.mp4`)
	start := time.Now()
	md5 := paw.GenMd5(path)
	elapsedTime := time.Since(start)
	fmt.Println("md5 Total time for excution:", elapsedTime.String())
	start = time.Now()
	md5sh := paw.GenMd5sh(path)
	elapsedTime = time.Since(start)
	fmt.Println("md5sh Total time for excution:", elapsedTime.String())
	paw.Logger.SetLevel(logrus.TraceLevel)
	paw.Logger.WithFields(logrus.Fields{
		"md5":       md5,
		"md5-len":   len(md5),
		"md5sh":     md5sh,
		"md5sh-len": len(md5sh),
		"equal":     md5 == md5sh,
	}).Debug()
}
