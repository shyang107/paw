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
	paw.Logger.WithFields(logrus.Fields{"View type": v.opt.ViewType}).Debug("view...")

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

	viewList(w, v, cur, fields, hasX, isViewNoDirs, isViewNoFiles)
	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
}

func viewList(w io.Writer, v *VFS, cur *Dir, fields []ViewField, hasX, isViewNoDirs, isViewNoFiles bool) {
	// paw.Logger.Debug()
	var (
		wdstty    = sttyWidth - 2
		tnd, tnf  = cur.NItems()
		nitems    = tnd + tnf
		nd, nf    int
		wdmeta    = 0
		roothead  = GetRootHeadC(cur, wdstty)
		head      = GetPFHeadS(chdp, fields...)
		totalsize int64
	)

	fmt.Fprintf(w, "%v\n", roothead)
	FprintBanner(w, "", "=", wdstty)

	// paw.Logger.Trace("hasX")
	if hasX {
		for _, fd := range fields {
			if fd&ViewFieldName == ViewFieldName {
				continue
			}
			wdmeta += fd.Width() + 1
		}
	}
	// paw.Logger.Trace("cur.relpaths")
	for _, rp := range cur.relpaths {
		var (
			curnd, curnf int
			size         int64
		)
		paw.Logger.WithFields(logrus.Fields{
			"rp": rp,
		}).Trace("getDir")
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
				FprintXattrs(w, wdmeta, de.Xattibutes())
			}
		}
		totalsize += size
		fprintDirSummary(w, "", curnd, curnf, size, wdstty)
		if nd+nf < nitems {
			FprintBanner(w, "", "-", wdstty)
		}
		if v.opt.Depth == 0 {
			break
		}
	}

	FprintBanner(w, "", "=", wdstty)
	FprintTotalSummary(w, "", nd, nf, totalsize, wdstty)
}
