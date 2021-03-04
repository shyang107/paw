package vfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func (v *VFS) ViewList(w io.Writer) {
	VFSViewList(w, v)
}

func VFSViewList(w io.Writer, v *VFS) {
	paw.Logger.Info("[vfs] " + v.opt.ViewType.String() + "...")

	var (
		fields                            = v.opt.ViewFields.Fields()
		hasX, isViewNoDirs, isViewNoFiles = v.hasX_NoDir_NoFiles()
	)

	cur := v.RootDir()

	if fields == nil {
		fields = DefaultViewFieldSlice
	}
	fields = checkFieldsHasGit(fields, cur.git.NoGit)
	modFieldWidths(v, fields)
	ViewFieldName.SetWidth(GetViewFieldNameWidthOf(fields))

	viewList(w, cur, fields, hasX, isViewNoDirs, isViewNoFiles)
	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
}

func viewList(w io.Writer, cur *Dir, fields []ViewField, hasX, isViewNoDirs, isViewNoFiles bool) {
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

		des, _ := cur.ReadDirAll()
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
			if de.IsDir() {
				if isViewNoDirs {
					nitems--
					continue
				}
				nd++
				curnd++
			} else {
				if isViewNoFiles {
					nitems--
					continue
				}
				nf++
				curnf++
				size += de.Size()
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
