package filetree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

var (
	pdOpt  *PrintDirOption
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
	if opt == nil {
		pdOpt = NewPrintDirOption()
	} else {
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
				fi, err := os.Stat(path)
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
				listOneFile(fl, files[0], pad)
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

func bmark(b bool) string {
	if b {
		return csup.Sprint("✓")
	}
	return cdashp.Sprint("✗")
}

func listOneFile(fl *FileList, path string, pad string) {
	paw.Logger.Info()

	var (
		w      = new(strings.Builder)
		wdstty = sttyWidth - 2 - paw.StringWidth(pad)
		width  = 0
	)
	file, err := NewFile(path)
	if err != nil {
		paw.Error.Println(err)
		return
	}

	head := cpmpt.Sprint("Directory: ") + pmptColorizedPath(file.Dir, "")
	fmt.Fprintln(w, head)
	printBanner(w, "", "=", wdstty)
	for _, wd := range pdOpt.fieldWidths {
		width = paw.MaxInt(width, wd)
	}
	// for _, field := range pdOpt.FieldKeys() {
	// 	width = paw.MaxInt(width, field.Width())
	// }
	// for _, field := range pfieldsMap {
	// 	width = paw.MaxInt(width, len(field))
	// }

	nline := 0
	fmt.Fprintln(w, rowFile(nline, PFieldName, file.BaseNameToLinkC(), width, wdstty))
	nline++
	fmt.Fprintln(w, rowFile(nline, PFieldPermissions, file.PermissionC(), width, wdstty))
	nline++
	fmt.Fprintln(w, rowFile(nline, PFieldINode, file.INodeC(), width, wdstty))
	nline++
	fmt.Fprintln(w, rowFile(nline, PFieldLinks, file.NLinksC(), width, wdstty))
	nline++
	fmt.Fprintln(w, rowFile(nline, PFieldSize, file.SizeC(), width, wdstty))
	nline++
	fmt.Fprintln(w, rowFile(nline, PFieldBlocks, file.BlocksC(), width, wdstty))
	nline++
	fmt.Fprintln(w, rowFile(nline, PFieldUser, file.UserC(), width, wdstty))
	nline++
	fmt.Fprintln(w, rowFile(nline, PFieldGroup, file.GroupC(), width, wdstty))
	nline++
	fmt.Fprintln(w, rowFile(nline, PFieldModified, file.ModifiedTimeC(), width, wdstty))
	nline++
	fmt.Fprintln(w, rowFile(nline, PFieldCreated, file.CreatedTimeC(), width, wdstty))
	nline++
	fmt.Fprintln(w, rowFile(nline, PFieldAccessed, file.AccessedTimeC(), width, wdstty))
	nline++
	git := fl.GetGitStatus()
	if !git.NoGit {
		cgit := file.GitXYc(git)
		fmt.Fprintln(w, rowFile(nline, PFieldGit, cgit, width, wdstty))
		nline++
	}
	fmt.Fprintln(w, rowFile(nline, PFieldMd5, file.GetMd5(), width, wdstty))
	nline++
	if len(file.XAttributes) > 0 {
		xfield := fmt.Sprintf("%[1]*[2]s%s", width, "Extended", " : ")
		wd := wdstty - width - 3
		sp := paw.Spaces(width + 3)
		xsymb := "@"
		wsymb := paw.StringWidth(xsymb)
		csymb := cxbp.Sprint(xsymb)
		cbsp := cxbp.Sprint(paw.Spaces(wsymb))
		for i, value := range file.XAttributes {
			wv := paw.StringWidth(value)
			if wv <= wd {
				if i == 0 {
					fmt.Fprintln(w, xfield+csymb, cxap.Sprint(value))
				} else {
					fmt.Fprintln(w, sp+csymb, cxap.Sprint(value))
				}
			} else {
				names := paw.WrapToSlice(value, width)
				if i == 0 {
					fmt.Fprintln(w, xfield+csymb, cxap.Sprint(names[0]))
				} else {
					fmt.Fprintln(w, sp+csymb, cxap.Sprint(names[0]))
				}
				for i := 1; i < len(names); i++ {
					fmt.Fprintln(w, sp+cbsp, cxap.Sprint(names[i]))
				}
			}
		}
	}

	printBanner(w, "", "=", wdstty)

	str := paw.PaddingString(w.String(), pad)
	fmt.Fprint(fl.Writer(), str)
}

func rowFile(nline int, flag PDFieldFlag, valueC string, width, wdstty int) (row string) {
	wvalueC := paw.StringWidth(paw.StripANSI(valueC))
	field := flag.Name() //pfieldsMap[flag]
	wfield := paw.StringWidth(field)
	sp := paw.Spaces(width - wfield)
	sptail := paw.Spaces(wdstty - width - 3 - wvalueC)
	if nline%2 == 0 {
		row = cpmpt.Sprintf("%s%s : %s", sp, field, valueC) + cpmpt.Sprint(sptail)
	} else {
		row = fmt.Sprint(sp + field + " : " + valueC + sptail)
	}
	return row
}

func listDirs(f *FileList, dirs []string, pad string, pdOpt *PrintDirOption) error {
	paw.Logger.Info()

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
	paw.Logger.Info()

	var (
		w     = f.stringBuilder
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
		// cdir := cdirp.Sprint(file.Dir + "/")
		// cname := file.NameC()
		fdName.ValueC = GetColorizedPath(file.Path, "") //cdir + cname
		fds.PrintRow(w, "")
		if isExtended && len(file.XAttributes) > 0 {
			// sp := paw.Spaces(wdmeta)
			// fmt.Fprint(w, xattrEdgeString(file, sp, wdmeta, wdstty))
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
	paw.Logger.Info()

	fl := NewFileList(root)
	// fl.IsSort = false
	fl.ResetWriters()
	if w != nil {
		fl.SetWriters(w)
	}
	fl.IsGrouped = isGrouped
	fl.IsSort = opt.SortOpt.IsSort

	// set sorter
	if fl.IsSort {
		setupFLSortOption(fl, pdOpt)
	}

	return fl
}

func setupFLSortOption(fl *FileList, opt *PrintDirOption) {
	paw.Logger.Info()

	if opt.ViewFlag&PTreeView == 0 ||
		opt.ViewFlag&PListTreeView == 0 {
		if opt.SortOpt.IsSort {
			fl.SetFilesSorter(opt.SortOpt.SortFlag.By())
		}
	}
}
