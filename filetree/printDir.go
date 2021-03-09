package filetree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

var (
	// pdOpt  *PrintDirOption
	pdOpt  = NewPrintDirOption()
	pdview PDViewFlag
)

// PrintDir will find files using codintion `ignore` func
func PrintDir(w io.Writer, path string, isGrouped bool, opt *PrintDirOption, pad string) (error, *FileList) {

	var (
		err  error
		root string
	)
	if opt.isTrace {
		paw.Logger.SetLevel(logrus.TraceLevel)
	}

	paw.Logger.Infof("root: %q", opt.Root)

	// setup root
	root, err = filepath.Abs(opt.Root)
	if err != nil {
		paw.Logger.Panic(err)
	}

	// check opt, fields to view
	if opt != nil {
		pdOpt = opt
		pdOpt.ConfigFields()
	}

	pdview = pdOpt.ViewFlag
	pdOpt.File, _ = NewFileRelTo(root, root)

	// check sortOpt
	if pdOpt.SortOpt == nil {
		pdOpt.SortOpt = &PDSortOption{
			IsSort:   true,
			SortFlag: PDSortByName,
		}
	}

	// check filter and ignore function
	if opt.Ignore == nil {
		opt.Ignore = DefaultIgnoreFn
	}
	pdOpt.ConfigFilter()

	// setup FileList
	fl := setFileList(w, root, isGrouped, pdOpt)

	// NPath > 0
	if pdOpt.NPath() > 0 {
		// one path or mutiple paths
		// sort.Sort(ByLowerString(pdOpt.Paths))
		var (
			dirs []string
			// files []string
			files = pdOpt.Paths
		)
		if pdOpt.Depth != 0 {
			for _, path := range pdOpt.Paths {
				fi, err := os.Lstat(path)
				if err != nil {
					if pdOpt.isTrace {
						paw.Logger.Error(err)
					}
					continue
				}
				if fi.IsDir() {
					dirs = append(dirs, path)
				}
			}
		}
		// files
		if len(files) > 0 {
			if len(files) == 1 {
				listOneFile(fl.Writer(), files[0], pad)
			} else {
				for _, path := range files {
					file, err := NewFile(path)
					if err != nil {
						if pdOpt.isTrace {
							paw.Logger.Error(err)
						}
						fl.AddError(path, err)
						continue
					}
					if err := pdOpt.Ignore(file, nil); err == SkipThis {
						continue
					}
					fl.addFilePD(file)
				}
				if fl.IsSort {
					fl.Sort0()
				}
				listFiles(fl, pad, pdOpt)

				fl.dirs = []string{}
				fl.store = make(FileMap)
			}
		}
		// dirs
		if len(dirs) > 0 {
			err = listDirs(fl, dirs, pad, pdOpt)
			if err != nil {
				return err, nil
			}
		}
	} else { // NPath == 0
		// use root as default
		fl.SetRoot(root)
		fl.SetIgnoreFunc(pdOpt.Ignore)
		err := fl.FindFiles(pdOpt.Depth)
		if err != nil {
			return err, nil
		}

		if pdOpt.isGit {
			fl.ConfigGit()
		}

		err = fl.DoView(pdOpt.ViewFlag, pad)
		if err != nil {
			return err, nil
		}

		// showlogrus()
		// fl.dumpAll(paw.Logger.Level)
	}

	return nil, fl
}

func listOneFile(wr io.Writer, path string, pad string) {
	paw.Logger.Debug()

	var (
		w      = new(strings.Builder)
		wdstty = sttyWidth - 2 - paw.StringWidth(pad)
	)
	file, err := NewFile(path)
	if err != nil {
		paw.Error.Println(err)
		return
	}

	head := fmt.Sprintf("Full path: %v", file.PathC())
	fmt.Fprintln(w, head)
	printBanner(w, "", "=", wdstty)

	git := NewGitStatus(file.Dir)
	fields := PFieldAllKeys
	// remove name field
	fields = fields[:len(fields)-1]
	width := PFieldPermissions.Width()
	for i, fd := range fields {
		fmt.Fprintln(w, rowFile(i, fd, file.FieldC(fd, git), width, wdstty))
	}
	if len(file.XAttributes) > 0 {
		xfield := fmt.Sprintf("%[1]*[2]s", width, "Extended")
		wd := wdstty - width - 3
		sp := paw.Spaces(width)
		xsymb := "@"
		wsymb := paw.StringWidth(xsymb)
		csymb := cxbp.Sprint(xsymb)
		cbsp := cxbp.Sprint(paw.Spaces(wsymb))
		for i, value := range file.XAttributes {
			wv := paw.StringWidth(value)
			c := rowColor(i)
			if wv <= wd {
				if i == 0 {
					c.Fprint(w, xfield+" : ")
				} else {
					c.Fprint(w, sp+"   ")
				}
				fmt.Fprintln(w, csymb, cxap.Sprint(value))
			} else {
				names := paw.WrapToSlice(value, width)
				if i == 0 {
					c.Fprint(w, xfield)
				} else {
					c.Fprint(w, sp)
				}
				fmt.Fprintln(w, " : "+csymb, cxap.Sprint(names[0]))
				for i := 1; i < len(names); i++ {
					c = rowColor(i)
					c.Fprintln(w, sp+"   "+cbsp, cxap.Sprint(names[i]))
				}
			}
		}
	}

	printBanner(w, "", "=", wdstty)

	str := paw.PaddingString(w.String(), pad)
	fmt.Fprint(wr, str)
}

var (
	chdEven = color.New([]color.Attribute{38, 5, 251, 1, 48, 5, 236}...)
	chdOdd  = color.New([]color.Attribute{38, 5, 159, 1, 48, 5, 234}...)
)

func rowColor(row int) *color.Color {
	var c *color.Color
	switch row % 2 {
	case 0:
		c = chdEven
	case 1:
		c = chdOdd
	}
	return c
}

func rowFile(nline int, flag PDFieldFlag, valueC string, width, wdstty int) (row string) {
	field := flag.Name() //pfieldsMap[flag]
	field = paw.FillLeft(field, width)
	row = rowColor(nline).Sprintf("%s", field) + " : " + valueC
	return row
}

func listDirs(f *FileList, dirs []string, pad string, pdOpt *PrintDirOption) error {
	paw.Logger.Debug()

	for i, path := range dirs {
		pdOpt.SetRoot(path)
		f.SetRoot(path)
		f.SetIgnoreFunc(pdOpt.Ignore)
		err := f.FindFiles(pdOpt.Depth)
		if err != nil {
			return err
		}

		err = f.DoView(pdOpt.ViewFlag, pad)
		if err != nil {
			return err
		}
		if i < len(dirs)-1 {
			fmt.Fprintln(f.Writer())
		}
		f.dirs = []string{}
		f.store = make(FileMap)
	}
	return nil
}

func listFiles(f *FileList, pad string, pdOpt *PrintDirOption) {
	paw.Logger.Debug()

	var (
		w     = f.StringBuilder()
		dirs  = f.Dirs()
		fm    = f.Map()
		files = []*File{}
		git   = f.git
		fds   = NewFieldSliceFrom(pdOpt.FieldKeys(), git)
		// fdSize     = fds.Get(PFieldSize)
		fdName     = fds.Get(PFieldName)
		wdstty     = sttyWidth - 2 - paw.StringWidth(pad)
		isExtended = isExtendedView(pdOpt.ViewFlag)
	)

	fds.ModifyWidth(f, wdstty)
	for _, dir := range dirs {
		for _, file := range fm[dir] {
			files = append(files, file)
		}
	}

	w.Reset()
	f.FprintAllErrs(w, "")
	printBanner(w, "", "=", wdstty)
	fds.PrintHeadRow(w, "")
	var size uint64
	for _, file := range files {
		if !file.IsDir() {
			size += file.Size
		}
		fds.SetValues(file, git)
		fdName.Value = file.Path
		fdName.ValueC = file.PathC()
		fds.PrintRow(w, "")
		if isExtended && len(file.XAttributes) > 0 {
			fds.PrintRowXattr(w, "", file.XAttributes, "")
		}
	}
	printBanner(w, "", "=", wdstty)

	fmt.Fprintln(w, f.TotalSummary(wdstty))

	str := paw.PaddingString(w.String(), pad)
	fmt.Fprintln(f.Writer(), str)

}

func isExtendedView(viewFlag PDViewFlag) bool {
	switch viewFlag {
	case PListExtendView, PTreeExtendView, PListTreeExtendView, PLevelExtendView, PTableExtendView:
		return true
	default:
		return false
	}
}

func setFileList(w io.Writer, root string, isGrouped bool, opt *PrintDirOption) *FileList {
	paw.Logger.Debug()

	fl := NewFileList(root)
	fl.ResetWriters()
	if w != nil {
		fl.SetWriters(w)
	}

	fl.SetMd5(opt.FieldFlag&PFieldMd5 != 0)
	fl.IsGrouped = isGrouped
	fl.IsSort = opt.SortOpt.IsSort

	// set sorter
	if fl.IsSort {
		// setupFLSortOption(fl, opt)
		if opt.ViewFlag&PTreeView == 0 ||
			opt.ViewFlag&PListTreeView == 0 {
			fl.SetFilesSorter(opt.SortOpt.SortFlag.By())
		}
	}

	return fl
}

func setupFLSortOption(fl *FileList, opt *PrintDirOption) {
	paw.Logger.Debug()

	if opt.ViewFlag&PTreeView == 0 ||
		opt.ViewFlag&PListTreeView == 0 {
		if opt.SortOpt.IsSort {
			fl.SetFilesSorter(opt.SortOpt.SortFlag.By())
		}
	}
}
