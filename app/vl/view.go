package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"

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
	fs.BuildFS()
	fs.View(os.Stdout)

	return nil
}

var (
	sttyHeight, sttyWidth = paw.GetTerminalSize()
	viewPaths_errors      = []error{}
)

func (opt *option) viewPaths() error {
	lg.Debug()

	var (
		w               = os.Stdout
		wdstty          = sttyWidth - 2
		paths           = opt.paths
		vfields         = opt.viewFields.Fields()
		wdmeta          = 0
		totalsize, size int64
		nd, nf          int
		tnd, tnf        int
		dxs, rm, dirs   = createBasepaths(paths)
		c               *color.Color
	)
	for i, dir := range dirs {
		if len(dirs) == 1 {
			c = paw.Cdirp
		} else {
			c = rcolor(i)
		}
		de := rm[dir]
		rhead := c.Sprintf("«root%d»", i+1)
		rhead += " directory: "
		rhead += vfs.PathToLinkC(de, nil)
		fmt.Fprintf(w, "%s\n", rhead)
	}

	fields := make([]vfs.ViewField, 0, len(vfields))
	var tmpFields vfs.ViewField
	for _, f := range vfields {
		if f&vfs.ViewFieldGit != 0 {
			continue
		}
		tmpFields |= f
		fields = append(fields, f)
	}
	opt.vopt.ViewFields = tmpFields

	modFieldWidths(dxs, fields)
	vfs.ViewFieldName.SetWidth(vfs.GetViewFieldNameWidthOf(fields))
	wdmeta = wdstty - vfs.ViewFieldName.Width()

	if len(viewPaths_errors) > 0 {
		for _, err := range viewPaths_errors {
			fmt.Fprintf(w, "%v\n", paw.Cerror.Sprint(err))
		}
	}

	vfs.FprintBanner(w, "", "=", wdstty)
	head := opt.vopt.ViewFields.GetHead(paw.Chdp, false)
	fmt.Fprintln(w, head)
	for i, dir := range dirs {
		if len(dirs) == 1 {
			c = paw.Cdirp
		} else {
			c = rcolor(i)
		}
		rooti := c.Sprintf("«root%d»", i+1) + paw.Cdirp.Sprint("/")
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
			fmt.Fprintf(w, "%v", opt.viewFields.RowStringXNameC(de))
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
			nd, nf, _ = fs.RootDir().NItems()
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

func rcolor(i int) *color.Color {
	switch i % 2 {
	case 0:
		return paw.CEven
		// return paw.Cdirp
	default:
		return paw.COdd
		// return paw.Cdirp.Add(paw.EXAColors["bgprompt"]...)
	}
}

type pathinfo struct {
	shortroot string
	info      os.FileInfo
	git       *vfs.GitStatus
	de        vfs.DirEntryX
}

// pathmap is map[path]{*pathrep}
type pathmap map[string][]pathinfo

// demap is map[dir][]vfs
type demap map[string][]vfs.DirEntryX

// srmap is map[dir]path
type srmap map[string]vfs.DirEntryX

func createBasepaths(paths []string) (dxs demap, srm srmap, dirs []string) {
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

	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			viewPaths_errors = append(viewPaths_errors, err)
			continue
		}
		dir := filepath.Dir(path)
		shortroot, ok := sm[dir]
		if !ok {
			idx++
			shortroot = fmt.Sprintf("root%d", idx)
			sm[dir] = shortroot
			rde := vfs.NewDir(dir, "", nil)
			srm[dir] = rde
			dxs[dir] = make([]vfs.DirEntryX, 0, len(paths))
			dirs = append(dirs, dir)
		}
		var de vfs.DirEntryX
		if info.IsDir() {
			de = vfs.NewDir(path, "", nil)
		} else {
			de = vfs.NewFile(path, "", nil)
		}
		// lg.WithField("de", de.Path()).Debug()
		dxs[dir] = append(dxs[dir], de)
	}
	sort.Sort(vfs.ByLowerString(dirs))
	return dxs, srm, dirs
}

func modFieldWidths(dxm demap, fields []vfs.ViewField) {
	for _, des := range dxm {
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
