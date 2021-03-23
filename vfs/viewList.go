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

func viewList(w io.Writer, rootdir *Dir, hasX, isViewNoDirs, isViewNoFiles bool) {
	// paw.Logger.Debug()
	var (
		vfields      = rootdir.opt.ViewFields
		wdstty       = sttyWidth - 2
		_, _, nitems = rootdir.NItems(true)
		tnd, tnf     int
		tsize        int64
		count        int
		roothead     = GetRootHeadC(rootdir, wdstty)
		fields       = vfields.Fields()
	)

	vfields.ModifyWidths(rootdir)
	// head := vfields.GetHeadFunc(paw.ChoseColorH)
	head := vfields.GetHead(paw.Chdp)

	fmt.Fprintf(w, "%v\n", roothead)
	FprintBanner(w, "", "=", wdstty)

	// paw.Logger.Trace("cur.relpaths")
	for _, rp := range rootdir.relpaths {
		if rootdir.opt.IsRelPathNotView(rp) {
			continue
		}
		paw.Logger.WithFields(logrus.Fields{
			"rp": rp,
		}).Trace("getDir")
		cur, err := rootdir.getDir(rp)
		if err != nil {
			paw.Logger.WithFields(logrus.Fields{
				"rp": rp,
			}).Fatal(err)
		}

		if rp != "." {
			cur.FprintlnRelPathC(w, "", false)
		}

		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, "", false)
		}

		des, _ := cur.ReadDirAll()
		if len(des) < 1 {
			continue
		}

		var (
			curnd, curnf int
			size         int64
			vnitems      = nitems
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
		fmt.Fprintf(w, "%v\n", head)
		for _, de := range des {
			if isSkipViewItem(de, isViewNoDirs, isViewNoFiles, &nitems, &curnd, &curnf, &size) {
				continue
			}
			count++
			// print fields of de
			// fmt.Fprintf(w, "%v\n", vfields.RowStringC(de))
			fmt.Fprintf(w, "%v\n", vfields.RowStringFC(de, fields))
			xrows := vfields.XattibutesRowsSC(de)
			if hasX && len(xrows) > 0 {
				for _, row := range xrows {
					fmt.Fprintln(w, row)
				}
			}
		}
		tnd += curnd
		tnf += curnf
		tsize += size
		if rootdir.opt.Depth != 0 {
			// cur.FprintlnSummaryC(w, "", wdstty, false)
			fmt.Fprintln(w, dirSummary("", curnd, curnf, size, wdstty))
			if count < nitems {
				FprintBanner(w, "", "-", wdstty)
			}
		}
	BAN:
		if rootdir.opt.Depth == 0 {
			break
		}
	}

	FprintBanner(w, "", "=", wdstty)
	// rootdir.FprintlnSummaryC(w, "", wdstty, true)
	fmt.Fprintln(w, totalSummary("", tnd, tnf, tsize, wdstty))
}
