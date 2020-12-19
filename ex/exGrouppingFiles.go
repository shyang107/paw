package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/_junk"
	"github.com/spf13/cast"
)

func exGrouppingFiles4() {
	paw.Logger.Info("")
	head := "GetFilesFuncString:\n"
	// sourceFolder := "../"
	// sourceFolder, _ := homedir.Expand("~/Downloads/")
	sourceFolder := "/Users/shyang/go/src/rover/opcc/"
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

	fileList := _junk.FileList{}
	// exclude := func(f paw.File) bool {
	// 	return (len(f.FileName) == 0 || paw.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
	// }
	exclude := func(f _junk.File) bool {
		return (len(f.FileName) == 0 || paw.HasPrefix(f.FileName, prefix) || _junk.REUsuallyExclude.MatchString(f.FullPath))
	}
	fileList.GetFilesFunc(sourceFolder, isRecursive, exclude)
	fileList.OrderedByFolder()
	// fileList.Print(os.Stdout, paw.OPlainTextMode, head, "# ")
	// fileList.Print(os.Stdout, paw.OTableFormatMode, head, "# ")
	fileList.Fprint(os.Stdout, _junk.OTreeMode, head, "# ")
	// fmt.Println(fileList)
}

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

	fileList := _junk.FileList{}
	fileList.GetFilesFunc(sourceFolder, isRecursive,
		func(f _junk.File) bool {
			return (len(f.FileName) == 0 || paw.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
		})
	fileList.OrderedByFolder()
	fileList.Fprint(os.Stdout, _junk.OTableFormatMode, head, "# ")
	// fmt.Println(head)
	// fmt.Println(fileList)
}

func exGrouppingFiles2() {
	paw.Logger.Info("")
	sourceFolder := "../"
	// sourceFolder, _ := homedir.Expand("~/Downloads/")
	// sourceFolder := "/Users/shyang/go/src/rover/opcc/"
	isRecursive := true
	sourceFolder, err := filepath.Abs(sourceFolder)
	if err != nil {
		paw.Logger.Error(err)
	}
	sourceFolder += "/"
	hsb := strings.Builder{}
	hsb.WriteString("GetFilesFuncString:\n")
	hsb.WriteString("- sourceFolder:	\"" + sourceFolder + "\"\n")
	hsb.WriteString("- isRecursive:	" + strconv.FormatBool(isRecursive) + "\n")
	prefix := "."
	regexPattern := `\.git|\$RECYCLE\.BIN|desktop\.ini`
	re := regexp.MustCompile(regexPattern)
	hsb.WriteString("- Excluding conditions:" + "\n")
	hsb.WriteString(`	- prefix:	"` + prefix + `"` + "\n")
	hsb.WriteString(`	- regexPattern:	"` + regexPattern + `"`)

	tp := &paw.TableFormat{
		Fields:    []string{"No.", "Sorted Files"},
		LenFields: []int{5, 100},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft},
		Padding:   "# ",
	}
	tp.Prepare(os.Stdout)
	tp.SetBeforeMessage(hsb.String())
	tp.PrintSart()

	files, err := _junk.GetFilesFunc(sourceFolder, isRecursive,
		func(f _junk.File) bool {
			return (len(f.FileName) == 0 || paw.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath))
		})
	if err != nil {
		paw.Logger.Error(err)
	}

	_junk.GrouppingFiles(files)

	oFolder := files[0].Folder
	gcount := 1
	j := 0
	for i, f := range files {
		if oFolder != f.Folder {
			oFolder = f.Folder
			tp.PrintRow("", "Sum: "+cast.ToString(j)+" files.")
			tp.PrintMiddleSepLine()
			j = 1
			gcount++
		} else {
			j++
		}
		if j == 1 {
			// folder := paw.TrimPrefix(f.Folder, sourceFolder)
			// tp.PrintRow("", fmt.Sprintf("folder %d: %q", gcount, "./"+folder))
			tp.PrintRow("", fmt.Sprintf("folder %d: %q", gcount, f.ShortFolder))
		}
		tp.PrintRow(j, f.File)
		// path := paw.TrimPrefix(f.FullPath, sourceFolder)
		// tp.PrintRow(j, path)
		if i == len(files)-1 {
			tp.PrintRow("", "Sum: "+cast.ToString(j)+" files.")
		}
	}
	tp.SetAfterMessage("Total: " + cast.ToString(gcount) + " subfolders and " + cast.ToString(len(files)) + " files. ")
	tp.PrintEnd()
}

func exGrouppingFiles1() {
	paw.Logger.Info("")
	// sourceFolder := "../"
	sourceFolder := "/Users/shyang/go/src/rover/opcc/"
	isRecursive := true
	sourceFolder, err := filepath.Abs(sourceFolder)
	if err != nil {
		paw.Logger.Error(err)
	}
	sourceFolder += "/"
	hsb := strings.Builder{}
	hsb.WriteString("GetFilesFuncString:\n")
	hsb.WriteString("  sourceFolder: " + `"../" <- "` + sourceFolder + "\"\n")
	hsb.WriteString("   isRecursive: " + strconv.FormatBool(isRecursive) + "\n")
	prefix := "."
	regexPattern := `\.git`
	re := regexp.MustCompile(regexPattern)
	hsb.WriteString("  Exculde:" + "\n")
	hsb.WriteString(`          prefix: "` + prefix + `"` + "\n")
	hsb.WriteString(`    regexPattern: "` + regexPattern + `"`)

	tp := &paw.TableFormat{
		Fields:    []string{"No.", "File", "Sorted Files"},
		LenFields: []int{5, 40, 40},
		Aligns:    []paw.Align{paw.AlignRight, paw.AlignLeft, paw.AlignLeft},
		Padding:   "# ",
	}
	tp.Prepare(os.Stdout)
	tp.SetBeforeMessage(hsb.String())
	tp.PrintSart()

	files, err := _junk.GetFilesFunc(sourceFolder, isRecursive,
		func(f _junk.File) bool {
			return len(f.FileName) == 0 || strings.HasPrefix(f.FileName, prefix) || re.MatchString(f.FullPath)
		})
	if err != nil {
		paw.Logger.Error(err)
	}
	sfiles := make([]_junk.File, len(files))
	copy(sfiles, files)
	_junk.GrouppingFiles(sfiles)

	oFolder := sfiles[0].Folder
	for i, f := range files {
		path := strings.TrimPrefix(f.FullPath, sourceFolder)
		spath := strings.TrimPrefix(sfiles[i].FullPath, sourceFolder)
		j := i + 1
		// if j%5 == 0 {
		if oFolder != sfiles[i].Folder {
			oFolder = sfiles[i].Folder
			tp.PrintMiddleSepLine()
		}
		tp.PrintRow(j, path, spath)
	}
	tp.PrintEnd()
}
