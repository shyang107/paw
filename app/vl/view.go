package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/shyang107/paw/cast"
	"github.com/sirupsen/logrus"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
)

func (opt *option) view() error {
	lg.WithFields(logrus.Fields{
		"Depth":          opt.vopt.Depth,
		"IsForceRecurse": opt.vopt.IsForceRecurse,
		"Grouping":       opt.vopt.Grouping,
		"ByField":        opt.vopt.ByField,
		"Skips":          opt.vopt.Skips,
		"ViewFields":     opt.vopt.ViewFields,
		"ViewType":       opt.vopt.ViewType,
	}).Debug()
	fs := vfs.NewVFS(opt.rootPath, opt.vopt)
	if opt.isDump {
		fs.Dump(os.Stdout)
	} else {
		fs.BuildFS()
		fs.View(os.Stdout)
	}

	return nil
}

var (
	sttyHeight, sttyWidth = paw.GetTerminalSize()
	viewPaths_errors      = []error{}
)

func (opt *option) viewPaths() error {
	lg.Debug()

	var (
		w       = os.Stdout
		wdstty  = sttyWidth - 2
		paths   = opt.paths
		vfields = opt.viewFields
		// vfields         = opt.viewFields.Fields()
		wdmeta          = 0
		totalsize, size int64
		nd, nf          int
		tnd, tnf        int
		c               *color.Color
	)

	dxs, srm, dirs := createBasepaths(paths)
	for i, dir := range dirs {
		c = paw.ChoseColor(i)
		de := srm[dir]
		rhead := c.Sprintf("{ R%d }: ", i+1)
		// rhead += " directory: "
		// rhead += vfs.PathToLinkC(de, nil)
		rhead += vfs.PathTo(de, &vfs.PathToOption{
			IsColor: true, Bgc: nil, PathReturn: vfs.PRTPathToLink})
		fmt.Fprintf(w, "%s\n", rhead)
	}
	vfields.RemoveGit(true)
	fields := vfields.Fields()
	dxs.modFieldWidths(fields)
	vfs.ViewFieldName.SetWidth(vfs.GetViewFieldNameWidthOf(fields))
	wdmeta = wdstty - vfs.ViewFieldName.Width()

	if len(viewPaths_errors) > 0 {
		for _, err := range viewPaths_errors {
			fmt.Fprintf(w, "%v\n", paw.Cerror.Sprint(err))
		}
	}

	vfs.FprintBanner(w, "", "=", wdstty)
	// head := vfields.GetHeadFunc(paw.ChoseColorH)
	head := vfields.GetHead(paw.Chdp)
	fmt.Fprintln(w, head)
	for i, dir := range dirs {
		c = paw.ChoseColor(i)
		rooti := c.Sprintf("{ R%d }/", i+1)
		des := dxs[dir]
		opt.vopt.ByField.Sort(des)
		for _, de := range des {
			if de.IsDir() {
				nd++
				size += de.Size()
			} else {
				nf++
			}
			// print fields of de
			fmt.Fprintf(w, "%v", vfields.RowStringXNameC(de))
			fmt.Fprintf(w, "%v\n", rooti+de.FieldC(vfs.ViewFieldName))

			if hasX {
				vfs.FprintXattrs(w, wdmeta, de.Xattibutes())
			}
			totalsize += de.Size()
		}
	}

	vfs.FprintBanner(w, "", "=", wdstty)
	vfs.FprintTotalSummary(w, "", nd, nf, totalsize, wdstty)

	if opt.depth == 0 {
		return nil
	}
	tnd += nd
	tnf += nf
	for _, dir := range dirs {
		des := dxs[dir]
		for _, de := range des {
			if !de.IsDir() {
				continue
			}
			// lg.WithFields(logrus.Fields{
			// 	"path": de.Path(),
			// 	"dir":  dir,
			// }).Debug()
			fmt.Fprintln(w)
			opt.rootPath = de.Path()
			if opt.depth > 0 {
				opt.depth--
			}
			fs := vfs.NewVFS(de.Path(), opt.vopt)
			fs.BuildFS()
			fs.View(os.Stdout)
			nd, nf, _ = fs.RootDir().NItems(true)
			tnd += nd
			tnf += nf
			totalsize += fs.TotalSize()
		}
	}
	// fmt.Fprintln(w)
	vfs.FprintBanner(w, "", "=", wdstty)
	vfs.FprintTotalSummary(w, paw.Cpmpt.Sprint("For all, "), tnd, tnf, totalsize, wdstty)
	return nil
}

// var (
// 	cevenH = paw.CloneColor(paw.CEvenH).Add(color.Underline)
// 	coddH  = paw.CloneColor(paw.COddH).Add(color.Underline)
// 	ceven  = paw.CloneColor(paw.CEven)
// 	codd   = paw.CloneColor(paw.COdd)
// )

// func rcolor(i int) *color.Color {
// 	switch i % 2 {
// 	case 0:
// 		return ceven
// 	default:
// 		return codd
// 	}
// }
// func rcolorH(i int) *color.Color {
// 	switch i % 2 {
// 	case 0:
// 		return cevenH
// 	default:
// 		return coddH
// 	}
// }

type pathinfo struct {
	shortroot string
	info      os.FileInfo
	git       *vfs.GitStatus
	de        vfs.DirEntryX
}

// pathmap is map[path]{*pathrep}
type pathmap map[string][]pathinfo

// demap is map[dir][]vfs.DirEntryX
type demap map[string][]vfs.DirEntryX

// srmap is map[dir]path
type srmap map[string]vfs.DirEntryX

func createBasepaths(paths []string) (dxs demap, srm srmap, dirs []string) {
	// paw.Logger.Debug()
	if len(paths) == 0 {
		return nil, nil, nil
	}
	dxs = make(demap)
	srm = make(srmap)
	dirs = make([]string, 0, len(paths))
	var (
		idx = 0
		sm  = make(map[string]string)
	)

	paw.Logger.WithFields(logrus.Fields{
		"paths": paths,
	}).Trace()
	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			viewPaths_errors = append(viewPaths_errors, err)
			paw.Logger.Error(err)
			continue
		}
		dir := filepath.Dir(path)
		shortroot, ok := sm[dir]
		if !ok {
			idx++
			shortroot = fmt.Sprintf("R%d", idx)
			sm[dir] = shortroot
			rde := vfs.NewDir(dir, "", nil, nil)
			srm[dir] = rde
			dxs[dir] = make([]vfs.DirEntryX, 0, len(paths))
			dirs = append(dirs, dir)
		}
		var de vfs.DirEntryX
		if info.IsDir() {
			de = vfs.NewDir(path, "", nil, nil)
		} else {
			de = vfs.NewFile(path, "", nil)
		}
		// lg.WithField("de", de.Path()).Debug()
		dxs[dir] = append(dxs[dir], de)
	}
	sort.Sort(vfs.ByLowerString(dirs))
	return dxs, srm, dirs
}

func (d *demap) modFieldWidths(fields []vfs.ViewField) {
	// func (d *demap) modFieldWidths(dxm demap, fields []vfs.ViewField) {
	for _, des := range *d {
		var (
			wd, dwd int
		)
		for _, de := range des {
			for _, fd := range fields {
				wd = de.WidthOf(fd)
				dwd = fd.Width()
				if !de.IsDir() && fd&vfs.ViewFieldSize == vfs.ViewFieldSize {
					if de.IsCharDev() || de.IsDev() {
						fmajor := vfs.ViewFieldMajor.Width()
						fminor := vfs.ViewFieldMinor.Width()
						major, minor := de.DevNumber()
						wdmajor := len(cast.ToString(major))
						wdminor := len(cast.ToString(minor))
						vfs.ViewFieldMajor.SetWidth(paw.MaxInt(fmajor, wdmajor))
						vfs.ViewFieldMinor.SetWidth(paw.MaxInt(fminor, wdminor))
						wd = vfs.ViewFieldMajor.Width() +
							vfs.ViewFieldMinor.Width() + 1
					}
				}
				width := paw.MaxInt(dwd, wd)
				fd.SetWidth(width)
			}
		}
	}
}
