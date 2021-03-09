package vfs

import (
	"fmt"
	"io"
	"os"

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
	modFieldWidths(cur, fields)

	viewList(w, cur, fields, hasX, isViewNoDirs, isViewNoFiles)
	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
}

func viewList(w io.Writer, cur *Dir, fields []ViewField, hasX, isViewNoDirs, isViewNoFiles bool) {
	// paw.Logger.Debug()
	var (
		wdstty         = sttyWidth - 2
		tnd, _, nitems = cur.NItems()
		nd, nf         int
		wdmeta         = 0
		roothead       = GetRootHeadC(cur, wdstty)
		head           = GetPFHeadS(chdp, fields...)
		totalsize      int64
		vfields        = cur.opt.ViewFields
	)

	fmt.Fprintf(w, "%v\n", roothead)
	FprintBanner(w, "", "=", wdstty)

	// paw.Logger.Trace("hasX")
	if hasX {
		wdmeta = GetViewFieldWidthWithoutName(cur.opt.ViewFields)
	}
	// paw.Logger.Trace("cur.relpaths")
	for _, rp := range cur.relpaths {
		if cur.opt.IsNotViewRelPath(rp) {
			continue
		}
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

		cdir, cname, cpath := GetPathC(rp)
		if rp != "." {
			cdir = cdirp.Sprint("./") + cdir
			cpath = cdir + cname
			fmt.Fprintf(w, "%v\n", cpath)
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

			// print fields of de
			fmt.Fprintf(w, "%v ", vfields.RowStringC(de))

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
		if cur.opt.Depth == 0 {
			break
		}
	}

	FprintBanner(w, "", "=", wdstty)
	FprintTotalSummary(w, "", nd, nf, totalsize, wdstty)
}
