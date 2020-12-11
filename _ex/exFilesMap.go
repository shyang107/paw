package main

import (
	"os"
	"path/filepath"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/_junk"
	"github.com/shyang107/paw/cast"
)

func exFilesMap(sourceFolder string) {
	head := "GetFilesFuncString:\n"
	// sourceFolder := "../"
	// sourceFolder, _ := homedir.Expand("~/Downloads/md/")
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
	// regexPattern := `\.git|\$RECYCLE\.BIN|desktop\.ini|funk|afero`
	// regexPattern := `\.git|\$RECYCLE\.BIN|desktop\.ini`
	// head += "	- regexPattern:	`" + regexPattern + "`"
	head += "	- regexPattern:	`" + _junk.ExcludePattern + "`"

	// re := regexp.MustCompile(regexPattern)

	fm := _junk.NewFilesMap()
	// exclude := func(f paw.File) bool {
	// 	return (len(f.FileName) == 0 || paw.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
	// }
	// exclude := func(f paw.File) bool {
	// 	return (len(f.FileName) == 0 || paw.HasPrefix(f.FileName, prefix) || paw.REUsuallyExclude.MatchString(f.FullPath))
	// }
	// fm.FindFilesFunc(sourceFolder, isRecursive, exclude)
	fm.FindFiles(sourceFolder, isRecursive)
	// fm.OrderedByFolder()
	fm.OrderedAll()
	pad := "# "
	// fm.Fprint(os.Stdout, paw.OPlainTextMode, head, pad)
	// fm.Fprint(os.Stdout, paw.OTableFormatMode, head, pad)
	fm.Fprint(os.Stdout, _junk.OTreeMode, head, pad)
	// fm.FprintText(os.Stdout, head, pad)
	// fmt.Println(fm)
	// fmt.Println(fm.Text(head, pad))
	// fmt.Println(fm.Table(head, pad))
	// fmt.Println(fm.Tree(head, pad))
}
