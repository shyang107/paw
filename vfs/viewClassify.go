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
		// nd, nf       int
		// roothead     = GetRootHeadC(rootdir, wdstty)
		_, _, crootpath = GetPathC(rootdir.Path(), nil)
		tnd, tnf        int
		tsize           int64
		count           int
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
		var (
			curnd, curnf    int
			size            int64
			vnitems         = nitems
			nfiles          = len(des)
			names           = make([]string, 0, nfiles)
			cnames          = make([]string, 0, nfiles)
			name, cname, sp string
			wd, idx, ncols  int
			wdcols          []int
		)
		if isViewNoDirs || isViewNoFiles {
			for _, de := range des {
				if isSkipViewItem(de, isViewNoDirs, isViewNoFiles, &nitems, &curnd, &curnf, &size) {
					continue
				}
			}
			if curnd+curnf == 0 {
				goto BAN
			}
			curnd, curnf, size, nitems = 0, 0, 0, vnitems
		}
		for _, de := range des {
			if isSkipViewItem(de, isViewNoDirs, isViewNoFiles, &nitems, &curnd, &curnf, &size) {
				continue
			}
			count++
			name = de.Name()
			cname = de.LSColor().Sprint(strings.TrimSpace(name))
			xattrs := de.Xattibutes()
			if xattrs == nil {
				names = append(names, name+"?")
				cnames = append(cnames, cname+paw.Cdashp.Sprint("?"))
			} else {
				if len(xattrs) > 0 {
					names = append(names, name+"@")
					cnames = append(cnames, cname+paw.Cdashp.Sprint("@"))
				} else {
					names = append(names, name+" ")
					cnames = append(cnames, cname+" ")
				}
			}
		}
		nfiles = len(names)
		wdcols = vcGridWidths(names, wdstty)
		ncols = len(wdcols)
		if nfiles < 1 {
			continue
		}
		for i := 0; i < nfiles; i += ncols {
			idx = i
			for j := 0; j < ncols; j++ {
				if idx > nfiles-1 {
					break
				}
				wd = paw.StringWidth(names[idx])
				sp = paw.Spaces(wdcols[j] - wd)
				fmt.Fprintf(w, "%s%s", cnames[idx], sp)
				idx++
			}
			fmt.Fprintln(w)
		}
		tnd += curnd
		tnf += curnf
		tsize += size
		// if rootdir.opt.Depth != 0 {
		// cur.FprintlnSummaryC(w, "", wdstty, false)
		// 	fmt.Fprintln(w, cur.SummaryC("", wdstty, false))
		// }
		if count < nitems {
			fmt.Fprintln(w)
			// FprintBanner(w, "", "-", wdstty)
		}
	BAN:
		if rootdir.opt.Depth == 0 {
			break
		}
	}

	fmt.Fprintln(w)
	// FprintBanner(w, "", "=", wdstty)
	fmt.Fprintln(w, totalSummary("", tnd, tnf, tsize, wdstty))
	// rootdir.FprintlnSummaryC(w, "", wdstty, true)
}

func appendNames(xattrs, names, cnames []string, name, cname string) {
	if xattrs == nil {
		names = append(names, name+"?")
		cnames = append(cnames, cname+paw.Cdashp.Sprint("?"))
	} else {
		if len(xattrs) > 0 {
			names = append(names, name+"@")
			cnames = append(cnames, cname+paw.Cdashp.Sprint("@"))
		} else {
			names = append(names, name+" ")
			cnames = append(cnames, cname+" ")
		}
	}
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
