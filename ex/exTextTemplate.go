package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/shyang107/paw"
)

func exTextTemplate() {
	paw.Logger.Info("")
	type option struct {
		Function, SourceFolder, Prefix, RegexPattern string
		Sep, TopBanner, MidBanner, BottomBanner      string
		Head                                         string
		IsRecursive                                  bool
	}
	var (
		opt = option{
			Function:     "GetFilesFuncString",
			IsRecursive:  true,
			Prefix:       ".",
			RegexPattern: `\.git|\$RECYCLE\.BIN|desktop\.ini`,
			Sep:          " ",
		}
		lh1    = 5
		lh2    = 80
		field1 = "No."
		field2 = "File"
	)

	opt.SourceFolder = "../"
	// opt.SourceFolder, _ := homedir.Expand("~/Downloads/")
	// opt.SourceFolder := "/Users/shyang/go/src/rover/opcc/"

	opt.SourceFolder, _ = filepath.Abs(opt.SourceFolder)
	opt.SourceFolder += "/"
	opt.Head = fmt.Sprintf("%[1]*[2]s%[3]s%-[4]*[5]s", lh1, field1, opt.Sep, lh2, field2)
	opt.TopBanner = paw.Repeat("=", lh1+lh2+len(opt.Sep))
	opt.MidBanner = paw.Repeat("-", lh1) + opt.Sep + paw.Repeat("-", lh2)
	opt.BottomBanner = paw.Repeat("=", lh1+lh2+len(opt.Sep))
	const headSection = `
{{ .Function }}:
{{ printf "- sourceFolder:\t %q" .SourceFolder }}
{{ printf "- isRecursive:\t %t" .IsRecursive }}
- Excluding conditions:
{{ printf "\t- prefix:\t%q" .Prefix }}
{{ printf "\t- regexPattern:\t%q" .RegexPattern }}
{{ .TopBanner }}
{{ .Head }}
{{ .MidBanner }}
`
	const endSection = `{{ .BottomBanner }}
`

	tmplH, err := template.New("head").Parse(headSection)
	if err != nil {
		paw.Logger.Fatalf("parsing: %s", err)
	}
	err = tmplH.Execute(os.Stdout, opt)
	if err != nil {
		paw.Logger.Fatalf("execution: %s", err)
	}

	re := regexp.MustCompile(opt.RegexPattern)
	files, err := paw.GetFilesFunc(opt.SourceFolder, opt.IsRecursive,
		func(f paw.File) bool {
			return (len(f.FileName) == 0 || paw.HasPrefix(f.FileName, opt.Prefix) || re.MatchString(f.FullPath))
		})
	if err != nil {
		paw.Logger.Error(err)
	}

	paw.GrouppingFiles(files)

	const rowSection = `{{ range $i,$v := . }}{{ printf "%5d %-80s\n" $i $v.ShortPath }}{{ end }}`
	tmplR, err := template.New("rows").Parse(rowSection)
	if err != nil {
		paw.Logger.Fatalf("parsing: %s", err)
	}
	err = tmplR.Execute(os.Stdout, files)
	if err != nil {
		paw.Logger.Fatalf("execution: %s", err)
	}

	tmplE, err := template.New("end").Parse(endSection)
	if err != nil {
		paw.Logger.Fatalf("parsing: %s", err)
	}
	err = tmplE.Execute(os.Stdout, opt)
	if err != nil {
		paw.Logger.Fatalf("execution: %s", err)
	}
}
