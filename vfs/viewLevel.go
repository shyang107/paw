package vfs

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func (v *VFS) ViewLevel(w io.Writer) {
	VFSViewLevel(w, v)
}

func VFSViewLevel(w io.Writer, v *VFS) {
	paw.Logger.WithFields(logrus.Fields{"View type": v.opt.ViewType}).Debug("view...")

	tmpfields := v.opt.ViewFields
	if v.opt.ViewFields&ViewFieldNo == 0 {
		v.opt.ViewFields = ViewFieldNo | v.opt.ViewFields
	}

	hasX, isViewNoDirs, isViewNoFiles := v.hasX_NoDir_NoFiles()
	viewLevel(w, v.RootDir(), hasX, isViewNoDirs, isViewNoFiles)

	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
	v.opt.ViewFields = tmpfields
}

func viewLevel(w io.Writer, cur *Dir, hasX, isViewNoDirs, isViewNoFiles bool) {
	var (
		vfields = cur.opt.ViewFields
		// fields           = vfields.GetModifyWidthsNoGitFields(cur, cur.git.NoGit)
		wdstty           = sttyWidth - 2
		tnd, tnf, nitems = cur.NItems()
		wdidx            = GetMaxWidthOf(tnd, tnf)
		nd, nf           int
		wdmeta           = 0
		roothead         = GetRootHeadC(cur, wdstty)
		totalsize        int64
	)
	vfields.ModifyWidths(cur)
	wdname := ViewFieldName.Width()
	// head := vfields.GetHead(paw.Chdp, cur.git.NoGit)

	fmt.Fprintf(w, "%v\n", roothead)
	FprintBanner(w, "", "=", wdstty)

	if hasX {
		wdmeta = GetViewFieldWidthWithoutName(cur.opt.ViewFields)
	}
	for _, rp := range cur.relpaths {
		if cur.opt.IsRelPathNotView(rp) {
			continue
		}
		var (
			level        int
			curnd, curnf int
			size         int64
		)
		if rp == "." {
			level = 0
		} else {
			level = len(strings.Split(rp, "/"))
		}
		wdpad := level * 3
		pad := paw.Spaces(wdpad)

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

		if level > 0 {
			slevel := fmt.Sprintf("L%d: ", level)
			FprintRelPath(w, pad, slevel, rp, false)
		}

		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, pad)
		}
		ViewFieldName.SetWidth(wdname - wdpad)
		head := vfields.GetHead(paw.Chdp, cur.git.NoGit)
		fmt.Fprintf(w, "%s%v\n", pad, head)
		for _, de := range des {
			var sidx string
			if de.IsDir() {
				if isViewNoDirs {
					nitems--
					continue
				}
				nd++
				curnd++
				sidx = fmt.Sprintf("D%-[1]*[2]d", wdidx, nd)
			} else {
				if isViewNoFiles {
					nitems--
					continue
				}
				nf++
				curnf++
				size += de.Size()
				sidx = fmt.Sprintf("F%-[1]*[2]d", wdidx, nf)
			}
			ViewFieldNo.SetValue(sidx)
			fmt.Fprintf(w, "%s", pad)

			// print fields of de
			fmt.Fprintf(w, "%v ", vfields.RowStringC(de))

			fmt.Println()
			if hasX {
				FprintXattrs(w, wdpad+wdmeta, de.Xattibutes())
			}
		}
		totalsize += size
		FprintDirSummary(w, pad, curnd, curnf, size, wdstty)
		if nd+nf < nitems {
			FprintBanner(w, "", "-", wdstty)
		}
		if cur.opt.Depth == 0 {
			break
		}
		ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
	}

	FprintBanner(w, "", "=", wdstty)
	FprintTotalSummary(w, "", nd, nf, totalsize, wdstty)
}
