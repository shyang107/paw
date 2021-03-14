package vfs

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func (v *VFS) ViewClassify(w io.Writer, fields []ViewField) {
	VFSViewClassify(w, v)
}
func VFSViewClassify(w io.Writer, v *VFS) {
	paw.Logger.WithFields(logrus.Fields{"View type": v.opt.ViewType}).Debug("view...")

	_, isViewNoDirs, isViewNoFiles := v.hasX_NoDir_NoFiles()

	viewClassify(w, v.RootDir(), isViewNoDirs, isViewNoFiles)

}

func viewClassify(w io.Writer, rootdir *Dir, isViewNoDirs, isViewNoFiles bool) {
	var (
		wdstty       = sttyWidth - 2
		_, _, nitems = rootdir.NItems(true)
		nd, nf       int
		// roothead     = GetRootHeadC(rootdir, wdstty)
		_, _, crootpath = GetPathC(rootdir.Path(), nil)
	)

	fmt.Fprintln(w, crootpath+":")
	// fmt.Fprintf(w, "%v\n", roothead)
	// FprintBanner(w, "", "=", wdstty)

	for _, rp := range rootdir.relpaths {
		if rootdir.opt.IsRelPathNotView(rp) {
			continue
		}
		paw.Logger.WithFields(logrus.Fields{"rp": rp}).Trace("getDir")
		cur, err := rootdir.getDir(rp)
		if err != nil {
			paw.Logger.WithFields(logrus.Fields{"rp": rp}).Fatal(err)
		}

		des, _ := cur.ReadDirAll()
		if len(des) < 1 {
			continue
		}

		if rp != "." {
			cur.FprintlnRelPathC(w, "", false)
			// fmt.Fprintln(w, cur.RelPathC("", false)+":")
			// FprintRelPath(w, "", "", "", rp, false)
		}

		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, "")
		}
		nfiles := len(des)
		names := make([]string, 0, nfiles)
		cnames := make([]string, 0, nfiles)
		for _, de := range des {
			name := de.Name()
			cname := de.LSColor().Sprint(strings.TrimSpace(name))
			isAppendName := true
			if de.IsDir() && isViewNoDirs {
				isAppendName = false
				nfiles--
				nitems--
			}
			if de.IsFile() && isViewNoFiles {
				isAppendName = false
				nfiles--
				nitems--
			}

			names, cnames = tailNames(de.Xattibutes(), names, cnames, name, cname, isAppendName)

			if de.IsDir() && !isViewNoDirs {
				nd++
			}
			if !de.IsDir() && !isViewNoFiles {
				nf++
			}
		}
		wdcols := vcGridWidths(names, wdstty)
		ncols := len(wdcols)
		if nfiles < 1 {
			continue
		}
		for i := 0; i < nfiles; i += ncols {
			idx := i
			for j := 0; j < ncols; j++ {
				if idx > nfiles-1 {
					break
				}
				wd := paw.StringWidth(names[idx])
				sp := paw.Spaces(wdcols[j] - wd)
				fmt.Fprintf(w, "%s%s", cnames[idx], sp)
				idx++
			}
			fmt.Fprintln(w)
		}
		// if rootdir.opt.Depth != 0 {
		// cur.FprintlnSummaryC(w, "", wdstty, false)
		// 	fmt.Fprintln(w, cur.SummaryC("", wdstty, false))
		// }
		if nd+nf < nitems {
			fmt.Fprintln(w)
			// FprintBanner(w, "", "-", wdstty)
		}
		if rootdir.opt.Depth == 0 {
			break
		}
	}

	fmt.Fprintln(w)
	// FprintBanner(w, "", "=", wdstty)
	rootdir.FprintlnSummaryC(w, "", wdstty, true)
	// fmt.Fprintln(w, rootdir.SummaryC("", wdstty, true))
}

func tailNames(xattrs, names, cnames []string, name, cname string, isAppendName bool) (tnames, tcnames []string) {
	tnames = make([]string, 0, len(names)+1)
	tcnames = make([]string, 0, len(cnames)+1)
	if xattrs == nil {
		if isAppendName {
			tnames = append(names, name+"?")
			tcnames = append(cnames, cname+paw.Cdashp.Sprint("?"))
		}
	} else {
		if len(xattrs) > 0 {
			if isAppendName {
				tnames = append(names, name+"@")
				tcnames = append(cnames, cname+paw.Cdashp.Sprint("@"))
			}
		} else {
			if isAppendName {
				tnames = append(names, name+" ")
				tcnames = append(cnames, cname+" ")
			}
		}
	}
	return tnames, tcnames
}

func vcGridWidths(names []string, wdstty int) (wdcols []int) {
	var (
		nf = len(names)
	)
	wds := make([]int, 0, nf)
	for _, name := range names {
		wds = append(wds, paw.StringWidth(name)+2)
	}
	wdcols = vcGridNcols(1, wds, wdstty)
	return wdcols
}

func vcGridNcols(nc int, wds []int, wdstty int) (wdcols []int) {
	wdcols = make([]int, nc)
	for i := 0; i < len(wds); i += nc {
		idx := i
		for j := 0; j < nc; j++ {
			if idx > len(wds)-1 {
				break
			}
			wdcols[j] = paw.MaxInt(wdcols[j], wds[idx])
			idx++
		}
	}
	if paw.SumIntA(wdcols...) < wdstty && nc < len(wds) {
		wdcols = vcGridNcols(nc+1, wds, wdstty)
	}
	if paw.SumIntA(wdcols...) > wdstty {
		for i := 0; i < len(wds); i += nc {
			idx := i
			for j := 0; j < nc; j++ {
				if idx > len(wds)-1 {
					break
				}
				wdcols[j] = paw.MaxInt(wdcols[j], wds[idx])
				idx++
			}
		}
		return wdcols[:nc]
	}
	return wdcols
}
