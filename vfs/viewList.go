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

	hasX, isViewNoDirs, isViewNoFiles := v.hasX_NoDir_NoFiles()
	viewList(w, v.RootDir(), hasX, isViewNoDirs, isViewNoFiles)

	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
}

func viewList(w io.Writer, cur *Dir, hasX, isViewNoDirs, isViewNoFiles bool) {
	// paw.Logger.Debug()
	var (
		vfields        = cur.opt.ViewFields
		fields         = vfields.GetModifyWidthsNoGitFields(cur, cur.git.NoGit)
		wdstty         = sttyWidth - 2
		tnd, _, nitems = cur.NItems()
		nd, nf         int
		wdmeta         = 0
		roothead       = GetRootHeadC(cur, wdstty)
		head           = GetPFHeadS(paw.Chdp, fields...)
		totalsize      int64
	)

	fmt.Fprintf(w, "%v\n", roothead)
	FprintBanner(w, "", "=", wdstty)

	// paw.Logger.Trace("hasX")
	if hasX {
		wdmeta = GetViewFieldWidthWithoutName(cur.opt.ViewFields)
	}
	// paw.Logger.Trace("cur.relpaths")
	for _, rp := range cur.relpaths {
		if cur.opt.IsRelPathNotView(rp) {
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
			cpath = paw.Cdirp.Sprint("./") + cdir + cname
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
