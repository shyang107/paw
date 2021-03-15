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
		vfields        = rootdir.opt.ViewFields
		wdstty         = sttyWidth - 2
		tnd, _, nitems = rootdir.NItems(true)
		nd, nf         int
		roothead       = GetRootHeadC(rootdir, wdstty)
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
		des, _ := cur.ReadDirAll()
		if len(des) < 1 {
			tnd--
			continue
		}

		if rp != "." {
			cur.FprintlnRelPathC(w, "", false)
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
			} else {
				if isViewNoFiles {
					nitems--
					continue
				}
				nf++
			}

			// print fields of de
			fmt.Fprintf(w, "%v ", vfields.RowStringC(de))

			fmt.Println()
			if hasX {
				xrows := vfields.XattibutesRowsSC(de)
				for _, row := range xrows {
					fmt.Fprintln(w, row)
				}
			}
		}
		if rootdir.opt.Depth != 0 {
			cur.FprintlnSummaryC(w, "", wdstty, false)

		}
		if nd+nf < nitems {
			FprintBanner(w, "", "-", wdstty)
		}
		if rootdir.opt.Depth == 0 {
			break
		}
	}

	FprintBanner(w, "", "=", wdstty)
	rootdir.FprintlnSummaryC(w, "", wdstty, true)
}
