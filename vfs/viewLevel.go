package vfs

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/cast"
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

func viewLevel(w io.Writer, rootdir *Dir, hasX, isViewNoDirs, isViewNoFiles bool) {
	var (
		vfields          = rootdir.opt.ViewFields
		wdstty           = sttyWidth - 2
		tnd, tnf, nitems = rootdir.NItems(true)
		tsize            int64
		count            int
		roothead         = GetRootHeadC(rootdir, wdstty)
	)
	tnd, tnf = 0, 0
	vfields.ModifyWidths(rootdir)
	wdname := ViewFieldName.Width()

	fmt.Fprintf(w, "%v\n", roothead)
	FprintBanner(w, "", "=", wdstty)

	idxmap := make(map[string]string)
	for _, rp := range rootdir.relpaths {
		if rootdir.opt.IsRelPathNotView(rp) {
			continue
		}
		var (
			level int
			idx   = idxmap[rp]
			cidx  = " [" + paw.Cvalue.Sprint(idx) + "] "
		)
		if rp == "." {
			level = 0
		} else {
			level = len(strings.Split(rp, "/"))
		}
		wdpad := level * 3
		pad := paw.Spaces(wdpad)

		paw.Logger.WithFields(logrus.Fields{"rp": rp}).Trace("getDir")
		cur, err := rootdir.getDir(rp)
		if err != nil {
			paw.Logger.WithFields(logrus.Fields{"rp": rp}).Fatal(err)
		}
		if level > 0 {
			slevel := paw.Cfield.Sprintf("L%d", level)
			cur.FprintlnRelPathC(w, pad+slevel+cidx, false)
		}

		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, pad, false)
		}

		des, _ := cur.ReadDirAll()
		if len(des) < 1 {
			continue
		}

		ViewFieldName.SetWidth(wdname - wdpad)
		var (
			curnd, curnf int
			size         int64
			vnitems      = nitems
			head         string
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
		// head := vfields.GetHeadFunc(paw.ChoseColorH)
		head = vfields.GetHead(paw.Chdp)
		fmt.Fprintf(w, "%s%v\n", pad, head)
		for _, de := range des {
			if isSkipViewItem(de, isViewNoDirs, isViewNoFiles, &nitems, &curnd, &curnf, &size) {
				continue
			}
			count++
			var sidx string
			if de.IsDir() {
				sidx = "D" + cast.ToString(tnd+curnd)
				idxmap[de.RelPath()] = sidx
			} else {
				sidx = "F" + cast.ToString(tnf+curnf)
			}
			ViewFieldNo.SetValue(sidx)
			fmt.Fprintf(w, "%s", pad)

			// print fields of de
			fmt.Fprintf(w, "%v \n", vfields.RowStringC(de))
			if hasX {
				xrows := vfields.XattibutesRowsSC(de)
				for _, row := range xrows {
					fmt.Fprintf(w, "%s%s\n", pad, row)
				}
			}
		}
		tnd += curnd
		tnf += curnf
		tsize += size
		if rootdir.opt.Depth != 0 {
			fmt.Fprintln(w, dirSummary(pad, curnd, curnf, size, wdstty))
			// cur.FprintlnSummaryC(w, pad, wdstty, false)
			if count < nitems {
				// fmt.Fprintln(w)
				// fmt.Fprintln(w, "count=", count, "nitems=", nitems)
				FprintBanner(w, "", "-", wdstty)
			}
		}
	BAN:
		if rootdir.opt.Depth == 0 {
			break
		}
		ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
	}

	FprintBanner(w, "", "=", wdstty)
	fmt.Fprintln(w, totalSummary("", tnd, tnf, tsize, wdstty))
	// rootdir.FprintlnSummaryC(w, "", wdstty, true)
}
