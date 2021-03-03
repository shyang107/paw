package vfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func (v *VFS) ViewClassify(w io.Writer) {
	paw.Logger.Info("[vfs] ViewClassify...")

	cur := v.RootDir()

	fields := []ViewField{ViewFieldName}

	viewClassify(w, cur, 0, fields)

}

func viewClassify(w io.Writer, cur *Dir, wdidx int, fields []ViewField) {
	var (
		wdstty    = sttyWidth - 2
		tnd, tnf  = cur.NItems()
		nitems    = tnd + tnf
		nd, nf    int
		roothead  = getRootHeadC(cur, wdstty)
		totalsize int64
	)

	fmt.Fprintf(w, "%v\n", roothead)

	for _, rp := range cur.relpaths {
		var (
			curnd, curnf int
			size         int64
		)

		cur, err := cur.getDir(rp)
		if err != nil {
			paw.Logger.WithFields(logrus.Fields{
				"rp": rp,
			}).Fatal(err)
		}

		des, _ := cur.ReadDirAll()
		if len(des) < 1 {
			// nitems--
			continue
		}

		cdir, cname := filepath.Split(rp)
		cname = cdip.Sprint(cname)
		cdir = cdirp.Sprint(cdir)
		if rp != "." {
			cdir = cdirp.Sprint("./") + cdir
			fmt.Fprintf(w, "%v\n", cdir+cname)
		}

		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, "")
		}
		nfiles := len(des)
		names := make([]string, 0, nfiles)
		cnames := make([]string, 0, nfiles)
		for _, de := range des {
			name := de.Name()
			cname = de.LSColor().Sprint(strings.TrimSpace(name))
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
			if de.Xattibutes() == nil {
				if isAppendName {
					names = append(names, name+"?")
					cnames = append(cnames, cname+cdashp.Sprint("?"))
				}
			} else {
				if len(de.Xattibutes()) > 0 {
					if isAppendName {
						names = append(names, name+"@")
						cnames = append(cnames, cname+cdashp.Sprint("@"))
					}
				} else {
					if isAppendName {
						names = append(names, name+" ")
						cnames = append(cnames, cname+" ")
					}
				}
			}
			if de.IsDir() && !isViewNoDirs {
				nd++
				curnd++
			}
			if !de.IsDir() && !isViewNoFiles {
				size += de.Size()
				nf++
				curnf++
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
		totalsize += size
		fprintDirSummary(w, "", curnd, curnf, size, wdstty)

		if nfiles < nitems {
			fprintBanner(w, "", "-", wdstty)
		}
	}

	fprintBanner(w, "", "=", wdstty)
	fprintTotalSummary(w, "", nd, nf, totalsize, wdstty)
}

func vcGridWidths(names []string, wdstty int) (wdcols []int) {
	var (
		nf = len(names)
	)
	wds := make([]int, 0, nf)
	for _, name := range names {
		wds = append(wds, paw.StringWidth(name)+2)
	}
	wdcols = vcGrisNcols(1, wds, wdstty)
	return wdcols
}

func vcGrisNcols(nc int, wds []int, wdstty int) (wdcols []int) {
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
	if paw.SumInts(wdcols...) < wdstty && nc < len(wds) {
		wdcols = vcGrisNcols(nc+1, wds, wdstty)
	}
	if paw.SumInts(wdcols...) > wdstty {
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