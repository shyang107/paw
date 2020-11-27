package main

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/cast"
)

func exGrouppingFiles3() {
	paw.Logger.Info("")
	head := "GetFilesFuncString:\n"
	sourceFolder := "../"
	// sourceFolder, _ := homedir.Expand("~/Downloads/")
	// sourceFolder := "/Users/shyang/go/src/rover/opcc/"
	sourceFolder, err := filepath.Abs(sourceFolder)
	if err != nil {
		paw.Logger.Error(err)
	}
	sourceFolder += "/"
	head += "- sourceFolder:	\"" + sourceFolder + "\"\n"

	isRecursive := true
	head += "- isRecursive:	" + cast.ToString(isRecursive) + "\n"

	head += "- Excluding conditions:" + "\n"
	prefix := "."
	head += "	- prefix:	`" + prefix + "`" + "\n"
	regexPattern := `\.git|\$RECYCLE\.BIN|desktop\.ini`
	// regexPattern := `\.git|\$RECYCLE\.BIN|desktop\.ini|funk|afero`
	head += "	- regexPattern:	`" + regexPattern + "`"

	re := regexp.MustCompile(regexPattern)

	fileList := paw.FileList{}
	fileList.GetFilesFunc(sourceFolder, isRecursive,
		func(f paw.File) bool {
			return (len(f.FileName) == 0 || paw.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
		})
	fileList.OrderedByFolder()
	fileList.Print(os.Stdout, head, "# ")
	// fmt.Println(head)
	// fmt.Println(fileList)
}
