package main

import (
	"fmt"
	"os"
	"path/filepath"

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
		w      = os.Stdout
		wdstty = sttyWidth - 2
		paths  = opt.paths
		fields = opt.viewFields.Fields()
		// basepaths
		head            = vfs.GetPFHeadS(paw.Chdp, fields...)
		pm, rdes        = createBasepaths(paths)
		wdmeta          = 0
		totalsize, size int64
		nd, nf          int
		tnd, tnf        int
	)

	wdpad := len(fmt.Sprint(len(rdes))) + 7
	for i, de := range rdes {
		pad := paw.Cpmpt.Sprintf("«root%d» ", i+1)
		roothead := vfs.GetRootHeadC(de, wdstty-wdpad)
		fmt.Fprintf(w, "%s%v\n", pad, roothead)
	}
	vfs.FprintBanner(w, "", "=", wdstty)

	modFieldWidths(pm, fields)
	if hasX {
		for _, fd := range fields {
			if fd&vfs.ViewFieldName == vfs.ViewFieldName {
				continue
			}
			wdmeta += fd.Width() + 1
		}
	}

	if len(viewPaths_errors) > 0 {
		for _, err := range viewPaths_errors {
			fmt.Fprintf(w, "%v\n", paw.Cerror.Sprint(err))
		}
	}

	fmt.Fprintln(w, head)
	for _, path := range paths {
		pi, ok := pm[path]
		if !ok {
			continue
		}
		// var (
		// 	size int64
		// )
		if pi.de.IsDir() {
			nd++
			size += pi.de.Size()
		} else {
			nf++
		}
		for _, field := range fields {
			var value string
			if field&vfs.ViewFieldName != 0 {
				value = paw.Cdirp.Sprintf("«%s»/", pi.shortroot) + pi.de.LSColor().Sprint(pi.name)
				// value = paw.Cdirp.Sprintf("%s/", pi.dir) + pi.de.LSColor().Sprint(pi.name)
			} else {
				value = pi.de.FieldC(field) + " "
			}
			fmt.Fprintf(w, "%v", value)
		}
		fmt.Fprintln(w)
		// fmt.Fprintln(w, paw.Cdirp.Sprint(pi.shortroot+"/")+pi.de.LSColor().Sprint(pi.name))
		if hasX {
			vfs.FprintXattrs(w, wdmeta, pi.de.Xattibutes())
		}
		totalsize += pi.de.Size()
	}

	vfs.FprintBanner(w, "", "=", wdstty)
	vfs.FprintTotalSummary(w, "", nd, nf, size, wdstty)

	if opt.depth == 0 {
		return nil
	}
	tnd += nd
	tnf += nf
	for _, path := range paths {
		// nd, nf, totalsize = 0, 0, 0
		pi, ok := pm[path]
		if !ok {
			continue
		}
		if pi.de.IsDir() {
			fmt.Fprintln(w)
			opt.rootPath = path
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
	dir       string
	name      string
	shortroot string
	info      os.FileInfo
	de        vfs.DirEntryX
}

// pathmap is map[path]{*pathrep}
type pathmap map[string]*pathinfo

// shortrootmap is map[shortname]path
type shortrootmap []vfs.DirEntryX

func createBasepaths(paths []string) (pm pathmap, srm shortrootmap) {
	if len(paths) == 0 {
		return nil, nil
	}
	pm = make(pathmap)
	srm = make(shortrootmap, 0)
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
			paw.Logger.Debug(idx, dir)
			rde := vfs.NewDir(dir, "")
			srm = append(srm, rde)
			lg.WithField("rde", rde.Path()).Debug()
		}
		// lg.WithField("path", path).Debug()
		var de vfs.DirEntryX
		if info.IsDir() {
			de = vfs.NewDir(path, "")
		} else {
			de = vfs.NewFile(path, "")
		}
		// lg.WithField("de", de.Path()).Debug()
		pm[path] = &pathinfo{
			dir:       dir,
			name:      info.Name(),
			info:      info,
			shortroot: shortroot,
			de:        de,
		}
	}
	return pm, srm
}

func modFieldWidths(pm pathmap, fields []vfs.ViewField) {
	for _, p := range pm {
		de := p.de
		// f, isFile := c.(*vfs.File)
		if !de.IsDir() {
			for _, fd := range fields {
				wd := de.WidthOf(fd)
				dwd := fd.Width()
				if fd&vfs.ViewFieldSize == vfs.ViewFieldSize {
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
						dwd = fd.Width()
					}
				}
				width := paw.MaxInt(dwd, wd)
				fd.SetWidth(width)
			}
		} else {
			// child := de.(*vfs.Dir)
			for _, fd := range fields {
				wd := de.WidthOf(fd)
				dwd := fd.Width()
				width := paw.MaxInt(dwd, wd)
				fd.SetWidth(width)
			}
			// childWidths(child, fields)
		}
	}
}
