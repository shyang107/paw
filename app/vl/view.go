package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
)

func (opt *option) view() error {
	lg.Debug()

	lg.Debug(opt.vopt)
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
		pm, rm, dirs    = createBasepaths(paths)
	)
	for i, dir := range dirs {
		de := rm[dir]
		rhead := fmt.Sprintf("«root%d» directory: ", i+1)
		rhead += vfs.PathToLinkC(de, nil)
		fmt.Fprintf(w, "%s\n", rhead)
	}
	vfs.FprintBanner(w, "", "=", wdstty)

	fields := make([]vfs.ViewField, 0, len(vfields))
	for _, f := range vfields {
		if f&vfs.ViewFieldGit != 0 {
			continue
		}
		fields = append(fields, f)
	}

	modFieldWidths(pm, fields)

	head := ""
	for _, fd := range fields {
		if fd&vfs.ViewFieldName == 0 {
			wdmeta += fd.Width() + 1
			switch fd.Align() {
			case paw.AlignLeft:
				head += paw.Chdp.Sprintf("%-[1]*[2]s", fd.Width(), fd.Name())
			default:
				head += paw.Chdp.Sprintf("%[1]*[2]s", fd.Width(), fd.Name())
			}
			head += " "
		} else {
			vfs.ViewFieldName.SetWidth(wdstty - wdmeta)
			head += paw.Chdp.Sprintf("%-[1]*[2]s", wdstty-wdmeta, fd.Name())
		}
	}

	if len(viewPaths_errors) > 0 {
		for _, err := range viewPaths_errors {
			fmt.Fprintf(w, "%v\n", paw.Cerror.Sprint(err))
		}
	}

	// head := vfs.GetPFHeadS(paw.Chdp, fields...)
	fmt.Fprintln(w, head)
	for i, dir := range dirs {
		rooti := paw.Cdirp.Sprintf("«root%d»/", i+1)
		for _, pi := range pm[dir] {
			if pi.de.IsDir() {
				nd++
				size += pi.de.Size()
			} else {
				nf++
			}
			for _, field := range fields {
				var value string
				if field&vfs.ViewFieldName != 0 {
					value = rooti + pi.de.FieldC(field)
				} else {
					value = pi.de.FieldC(field) + " "
				}
				fmt.Fprintf(w, "%v", value)
			}
			fmt.Fprintln(w)
			if hasX {
				vfs.FprintXattrs(w, wdmeta, pi.de.Xattibutes())
			}
			totalsize += pi.de.Size()
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
		for _, pi := range pm[dir] {
			if !pi.de.IsDir() {
				continue
			}
			fmt.Fprintln(w)
			opt.rootPath = pi.de.Path()
			fs := vfs.NewVFS(opt.rootPath, opt.vopt)
			fs.BuildFS()
			fs.View(os.Stdout)
			nd, nf = fs.RootDir().NItems()
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

type pathinfo struct {
	shortroot string
	info      os.FileInfo
	de        vfs.DirEntryX
}

// pathmap is map[path]{*pathrep}
type pathmap map[string][]pathinfo

// shortrootmap is map[dir]path
type shortrootmap map[string]vfs.DirEntryX

func createBasepaths(paths []string) (pm pathmap, srm shortrootmap, dirs []string) {
	if len(paths) == 0 {
		return nil, nil, nil
	}
	pm = make(pathmap)
	srm = make(shortrootmap)
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
		// git := gm[dir]
		if !ok {
			idx++
			shortroot = fmt.Sprintf("root%d", idx)
			sm[dir] = shortroot
			// gm[dir] = vfs.NewGitStatus(dir)
			// lg.Debug(idx, dir)
			rde := vfs.NewDir(dir, "", nil)
			srm[dir] = rde
			pm[dir] = make([]pathinfo, 0, len(paths))
			dirs = append(dirs, dir)
			// lg.WithField("rde", rde.Path()).Debug()
		}
		// lg.WithField("path", path).Debug()
		var de vfs.DirEntryX
		if info.IsDir() {
			de = vfs.NewDir(path, "", nil)
		} else {
			de = vfs.NewFile(path, "", nil)
		}
		// lg.WithField("de", de.Path()).Debug()
		pm[dir] = append(pm[dir], pathinfo{
			info:      info,
			shortroot: shortroot,
			de:        de,
		})
	}
	sort.Sort(vfs.ByLowerString{Values: dirs})
	return pm, srm, dirs
}

func modFieldWidths(pm pathmap, fields []vfs.ViewField) {
	for _, pis := range pm {
		for _, p := range pis {
			de := p.de
			var (
				wd, dwd int
			)
			for _, fd := range fields {
				wd = de.WidthOf(fd)
				dwd = fd.Width()
				if !de.IsDir() && fd&vfs.ViewFieldSize == vfs.ViewFieldSize {
					if de.IsCharDev() || de.IsDev() {
						fmajor := vfs.ViewFieldMajor.Width()
						fminor := vfs.ViewFieldMinor.Width()
						major, minor := de.DevNumber()
						wdmajor := len(fmt.Sprint(major))
						wdminor := len(fmt.Sprint(minor))
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
