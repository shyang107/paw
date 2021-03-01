package vfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func (v *VFS) ViewList(w io.Writer, fields []ViewField, hasX bool) {
	paw.Logger.Info("[vfs] ViewList...")

	cur := v.RootDir()

	if fields == nil {
		fields = DefaultViewFields
	}
	fields = checkFieldsHasGit(fields, cur.git.NoGit)

	modFieldWidths(v, fields)

	viewList(w, cur, fields, hasX)

}

func viewList(w io.Writer, cur *Dir, fields []ViewField, hasX bool) {
	var (
		wdstty    = sttyWidth - 2
		tnd, tnf  = cur.NItems()
		nitems    = tnd + tnf
		nd, nf    int
		wdmeta    = 0
		roothead  = getRootHeadC(cur, wdstty)
		head      = getPFHeadS(chdp, fields...)
		totalsize int64
	)

	fmt.Fprintf(w, "%v\n", roothead)
	fprintBanner(w, "", "=", wdstty)

	if hasX {
		for _, fd := range fields {
			if fd&ViewFieldName == ViewFieldName {
				continue
			}
			wdmeta += fd.Width() + 1
		}
	}
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
		cname = cdip.Sprint(cname)
		cdir = cdirp.Sprint(cdir)
		if rp != "." {
			cdir = cdirp.Sprint("./") + cdir
			fmt.Fprintf(w, "%v\n", cdir+cname)
		}

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
			if hasX {
				fprintXattrs(w, wdmeta, de.Xattibutes())
			}
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
