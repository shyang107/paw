package vfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func (v *VFS) ViewClassify(w io.Writer, fields []ViewField) {
	paw.Logger.Info("[vfs] LevelView...")

	cur := v.RootDir()

	if fields == nil {
		fields = DefaultViewFields
	}
	fields = checkFieldsHasGit(fields, cur.git.NoGit)

	modFieldWidths(v, fields)

	viewClassify(w, cur, 0, fields)

}

func viewClassify(w io.Writer, cur *Dir, wdidx int, fields []ViewField) {
	var (
		wdstty   = sttyWidth - 2
		tnd, tnf = cur.NItems()
		nitems   = tnd + tnf
		nd, nf   int
		// wdmeta    = 0
		roothead  = getRootHeadC(cur, wdstty)
		head      = getPFHeadS(chdp, fields...)
		totalsize int64
	)

	fmt.Fprintf(w, "%v\n", roothead)
	fprintBanner(w, "", "=", wdstty)

	// if hasX {
	// 	for _, fd := range fields {
	// 		if fd&ViewFieldName == ViewFieldName {
	// 			continue
	// 		}
	// 		wdmeta += fd.Width() + 1
	// 	}
	// }
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

		des, _ := cur.ReadDir(-1)
		cur.resetIdx()
		if len(des) < 1 {
			tnd--
			continue
		}

		cdir, cname := filepath.Split(rp)
		cdir = cdirp.Sprint(cdir)
		if rp != "." {
			cdir = cdirp.Sprint("./") + cdir
		}
		cname = cdip.Sprint(cname)

		fmt.Fprintf(w, "%v\n", cdir+cname)
		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, "")
		}

		fmt.Fprintf(w, "%v\n", head)
		for _, de := range des {
			if de.IsFile() {
				nf++
				curnf++
				size += de.Size()
			} else {
				nd++
				curnd++
			}
			for _, field := range fields {
				fmt.Fprintf(w, "%v ", de.FieldC(field))
			}
			fmt.Println()
			// if hasX {
			// 	fprintXattrs(w, wdmeta, de.Xattibutes())
			// }
		}
		totalsize += size
		fprintDirSummary(w, "", curnd, curnf, size, wdstty)
		if nd+nf < nitems {
			fprintBanner(w, "", "-", wdstty)
		}
	}

	fprintBanner(w, "", "=", wdstty)
	fprintTotalSummary(w, "", nd, nf, totalsize, wdstty)
}
