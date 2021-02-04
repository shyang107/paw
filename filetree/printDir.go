package filetree

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shyang107/paw"
)

var (
	pdOpt  *PrintDirOption
	pdview PDViewFlag
)

// PrintDir will find files using codintion `ignore` func
func PrintDir(w io.Writer, path string, isGrouped bool, opt *PrintDirOption, pad string) (error, *FileList) {

	var (
		err     error
		root    string
		sortOpt *PDSortOption
	)

	// paw.Logger.WithField("root", opt.Root).Info()

	// setup root
	root, err = filepath.Abs(opt.Root)
	if err != nil {
		paw.Error.Println(err)
		os.Exit(1)
	}

	// check opt
	if opt == nil {
		pdOpt = NewPrintDirOption()
	} else {
		pdOpt = opt
	}

	// check fields to view
	checkFieldFlag(pdOpt)

	// check sortOpt
	if pdOpt.SortOpt == nil {
		sortOpt = &PDSortOption{
			IsSort:  true,
			SortWay: PDSortByName,
		}
	} else {
		sortOpt = pdOpt.SortOpt
	}

	// check ignore function
	if opt.Ignore == nil {
		opt.Ignore = DefaultIgnoreFn
	}

	// get view option
	pdview = pdOpt.OutOpt

	// setup FileList
	fl := setFileList(w, root, isGrouped, sortOpt)

	if !fl.IsSort {
		goto FIND
	}

	// setup sort options of FileList
	setupFLSortOption(fl, pdOpt, sortOpt)

FIND:
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
					paw.Logger.Error(err)
					continue
				}
				if fi.IsDir() {
					dirs = append(dirs, path)
				}
				// else {
				// 	files = append(files, path)
				// }
			}
		}
		// files
		if len(files) > 0 {
			for _, path := range files {
				file, err := NewFile(path)
				if err != nil {
					paw.Logger.Error(err)
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
			cehckAndFiltPrintDirFiltOpt(fl, pdOpt.FiltOpt)
			listFiles(fl, pad, pdOpt)

			fl.dirs = []string{}
			fl.store = make(FileMap)
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
		fl.FindFiles(pdOpt.Depth, pdOpt.Ignore)
		cehckAndFiltPrintDirFiltOpt(fl, pdOpt.FiltOpt)
		err = switchFileListView(fl, pdOpt.OutOpt, pad)
		if err != nil {
			return err, nil
		}
	}

	return nil, fl
}

func listDirs(f *FileList, dirs []string, pad string, pdOpt *PrintDirOption) error {
	for _, path := range dirs {
		pdOpt.SetRoot(path)
		f.SetRoot(path)
		f.FindFiles(pdOpt.Depth, pdOpt.Ignore)
		cehckAndFiltPrintDirFiltOpt(f, pdOpt.FiltOpt)
		err := switchFileListView(f, pdOpt.OutOpt, pad)
		if err != nil {
			return err
		}
		f.dirs = []string{}
		f.store = make(FileMap)
	}
	return nil
}

func listFiles(f *FileList, pad string, pdOpt *PrintDirOption) {
	var (
		w     = f.stringBuilder
		dirs  = f.Dirs()
		fm    = f.Map()
		files = []*File{}
		git   = f.GetGitStatus()
		fds   = NewFieldSliceFrom(pfieldKeys, git)
		// fdSize     = fds.Get(PFieldSize)
		fdName     = fds.Get(PFieldName)
		wdstty     = sttyWidth - 2 - paw.StringWidth(pad)
		isExtended = isExtendedView(pdOpt.OutOpt)
	)

	fds.ModifyWidth(f, wdstty)
	for _, dir := range dirs {
		for _, file := range fm[dir] {
			files = append(files, file)
			// wsize, _, _ := file.widthOfSize()
			// fdSize.Width = paw.MaxInt(fdSize.Width, wsize)
		}
	}

	w.Reset()

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

	fmt.Fprintln(w, f.TotalSummary())

	str := paw.PaddingString(w.String(), pad)
	fmt.Fprintln(f.Writer(), str)

}

func isExtendedView(outOpt PDViewFlag) bool {
	switch outOpt {
	case PListExtendView, PTreeExtendView, PListTreeExtendView, PLevelExtendView, PTableExtendView:
		return true
	default:
		return false
	}
}

func switchFileListView(fl *FileList, outOpt PDViewFlag, pad string) error {
	switch outOpt {
	case PListView:
		fl.ToListView(pad)
	case PListExtendView:
		fl.ToListExtendView(pad)
	case PTreeView:
		fl.ToTreeView(pad)
	case PTreeExtendView:
		fl.ToTreeExtendView(pad)
	case PListTreeView:
		fl.ToListTreeView(pad)
	case PListTreeExtendView:
		fl.ToListTreeExtendView(pad)
	case PLevelView:
		fl.ToLevelView(pad, false)
	case PLevelExtendView:
		fl.ToLevelView(pad, true)
	case PTableView:
		fl.ToTableView(pad, false)
	case PTableExtendView:
		fl.ToTableView(pad, true)
	case PClassifyView:
		fl.ToClassifyView(pad)
	default:
		return errors.New("No this option of PrintDir")
	}
	return nil
}

func cehckAndFiltPrintDirFiltOpt(fl *FileList, filtOpt *PDFilterOption) {
	if filtOpt != nil && filtOpt.IsFilt {
		switch filtOpt.FiltWay {
		case PDFiltNoEmptyDir: // no empty dir
			flf := NewFileListFilter(fl, []Filter{FiltEmptyDirs})
			flf.Filt()
		case PDFiltJustDirs: // no files
			flf := NewFileListFilter(fl, []Filter{FiltJustDirs})
			flf.Filt()
		case PDFiltJustFiles: // PDFiltJustFilesButNoEmptyDir // no dirs
			flf := NewFileListFilter(fl, []Filter{FiltJustFiles})
			flf.Filt()
		case PDFiltJustDirsButNoEmpty: // no file and no empty dir
			flf := NewFileListFilter(fl, []Filter{FiltEmptyDirs, FiltJustDirs})
			flf.Filt()
		}
	}
}

func setupFLSortOption(fl *FileList, opt *PrintDirOption, sortOpt *PDSortOption) {

	if opt.OutOpt&PTreeView == 0 ||
		opt.OutOpt&PListTreeView == 0 {
		if sortOpt.IsSort {
			if _, ok := sortByField[sortOpt.SortWay]; ok {
				fl.SetFilesSorter(sortByField[sortOpt.SortWay])
			} else {
				fl.SetFilesSorter(sortByField[PDSortByName])
			}
			// 	switch sortOpt.SortWay {
			// 	case PDSortByINode:
			// 		fl.SetFilesSorter(byINode)
			// 	case PDSortByReverseINode:
			// 		fl.SetFilesSorter(byINodeR)
			// 	case PDSortByLinks:
			// 		fl.SetFilesSorter(byLinks)
			// 	case PDSortByReverseLinks:
			// 		fl.SetFilesSorter(byLinksR)
			// 	case PDSortBySize:
			// 		fl.SetFilesSorter(bySize)
			// 	case PDSortByReverseSize:
			// 		fl.SetFilesSorter(bySizeR)
			// 	case PDSortByBlocks:
			// 		fl.SetFilesSorter(byBlocks)
			// 	case PDSortByReverseBlocks:
			// 		fl.SetFilesSorter(byBlocksR)
			// 	case PDSortByMTime:
			// 		fl.SetFilesSorter(byMTime)
			// 	case PDSortByReverseMTime:
			// 		fl.SetFilesSorter(byMTimeR)
			// 	case PDSortByCTime:
			// 		fl.SetFilesSorter(byCTime)
			// 	case PDSortByReverseCTime:
			// 		fl.SetFilesSorter(byCTimeR)
			// 	case PDSortByATime:
			// 		fl.SetFilesSorter(byATime)
			// 	case PDSortByReverseATime:
			// 		fl.SetFilesSorter(byATimeR)
			// 	case PDSortByReverseName:
			// 		fl.SetFilesSorter(byNameR)
			// 	default: //case PDSortByName :
			// 		fl.SetFilesSorter(byName)
			// 	}
		}
	}
}

func setFileList(w io.Writer, root string, isGrouped bool, sortOpt *PDSortOption) *FileList {
	fl := NewFileList(root)
	// fl.IsSort = false
	fl.ResetWriters()
	if w != nil {
		fl.SetWriters(w)
	}
	fl.IsGrouped = isGrouped
	fl.IsSort = sortOpt.IsSort

	return fl
}

// func setIgnoreFn(opt *PrintDirOption) {
// 	if opt.Ignore == nil {
// 		opt.Ignore = DefaultIgnoreFn
// 	}
// }

func checkFieldFlag(opt *PrintDirOption) {
	if opt.FieldFlag&PFieldINode != 0 {
		pfieldKeys = append(pfieldKeys, PFieldINode)
	}

	pfieldKeys = append(pfieldKeys, PFieldPermissions)

	if opt.FieldFlag&PFieldLinks != 0 {
		pfieldKeys = append(pfieldKeys, PFieldLinks)
	}

	pfieldKeys = append(pfieldKeys, PFieldSize)

	if opt.FieldFlag&PFieldBlocks != 0 {
		pfieldKeys = append(pfieldKeys, PFieldBlocks)
	}

	pfieldKeys = append(pfieldKeys, PFieldUser)
	pfieldKeys = append(pfieldKeys, PFieldGroup)

	if opt.FieldFlag&PFieldModified != 0 {
		pfieldKeys = append(pfieldKeys, PFieldModified)
	}
	if opt.FieldFlag&PFieldCreated != 0 {
		pfieldKeys = append(pfieldKeys, PFieldCreated)
	}
	if opt.FieldFlag&PFieldAccessed != 0 {
		pfieldKeys = append(pfieldKeys, PFieldAccessed)
	}

	if opt.FieldFlag&PFieldGit != 0 {
		pfieldKeys = append(pfieldKeys, PFieldGit)
	}
	// pfieldKeys = append(pfieldKeys, PFieldGit)
	pfieldKeys = append(pfieldKeys, PFieldName)

	for _, k := range pfieldKeys {
		pfields = append(pfields, pfieldsMap[k])
		pfieldWidths = append(pfieldWidths, pfieldWidthsMap[k])
	}
}

// func checkPrintDirOption(opt *PrintDirOption) {
// 	if opt == nil {
// 		opt = NewPrintDirOption()
// 		// opt = &PrintDirOption{
// 		// 	Depth:  0,
// 		// 	OutOpt: PListView,
// 		// 	// OutOpt: PListExtendView,
// 		// 	// OutOpt: PTreeView,
// 		// 	// OutOpt: PListTreeView,
// 		// 	// OutOpt: PLevelView,
// 		// 	// OutOpt: PTableView,
// 		// 	// OutOpt: PClassifyView,
// 		// 	FieldFlag: PFieldModified,
// 		// 	Ignore:    DefaultIgnoreFn,
// 		// }
// 	}
// }

// func checkSortOpt(sortOpt *PDSortOption) *PDSortOption {
// 	if sortOpt == nil {
// 		return &PDSortOption{
// 			IsSort:  true,
// 			SortWay: PDSortByName,
// 		}
// 	}
// 	return sortOpt
// }

// var errBreak = errors.New("return nil")

// func checkAndPrintFile(w io.Writer, path string, pad string) error {

// 	// paw.Logger.WithField("path", path).Info()

// 	file, err := NewFile(path)
// 	if err != nil {
// 		return err
// 	}
// 	if !file.IsDir() {
// 		fmt.Fprintf(w, "%sDirectory: %v \n", pad, GetColorizedDirName(file.Dir, ""))
// 		git, _ := GetShortGitStatus(file.Dir)
// 		fds := NewFieldSliceFrom(pfieldKeysDefualt, git)
// 		fl := NewFileList(file.Dir)
// 		fds.ModifyWidth(fl, sttyWidth-2)
// 		fds.SetValues(file, git)
// 		fmt.Fprintln(w, fds.ColorHeadsString())
// 		// fmt.Fprint(w, rowWrapFileName(file, fds, pad, sttyWidth-2))
// 		fds.PrintRow(w, pad)
// 		return errBreak
// 	}
// 	return nil
// }

// func cleanPath(path string) string {
// 	paw.Logger.WithField("path", path).Info()

// 	tpath := path
// 	// if strings.Contains(tpath, "~") {
// 	// 	tpath = strings.ReplaceAll(tpath, "~", paw.GetHomeDir())
// 	// }
// 	// paw.Logger.WithField("~", tpath).Info()
// 	// tpath = filepath.Clean(tpath)
// 	// paw.Logger.WithField("clean", tpath).Info()

// 	tpath, err := filepath.Abs(tpath)
// 	if err != nil {
// 		paw.Logger.Error(err)
// 		return tpath
// 	}
// 	// if !filepath.IsAbs(tpath) {
// 	// 	tpath, err := filepath.Abs(tpath)
// 	// 	if err != nil {
// 	// 		paw.Logger.Error(err)
// 	// 		return tpath
// 	// 	}
// 	// }
// 	paw.Logger.WithField("abs", tpath).Info()

// 	return tpath
// }
