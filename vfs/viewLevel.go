package vfs

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
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
		wdidx            = GetMaxWidthOf(tnd, tnf)
		nd, nf           int
		wdmeta           = 0
		roothead         = GetRootHeadC(rootdir, wdstty)
		ceven            = paw.CloneColor(paw.CEven).Add(color.Underline)
		codd             = paw.CloneColor(paw.COdd).Add(color.Underline)
	)
	vfields.ModifyWidths(rootdir)
	wdname := ViewFieldName.Width()

	fmt.Fprintf(w, "%v\n", roothead)
	FprintBanner(w, "", "=", wdstty)

	if hasX {
		wdmeta = GetViewFieldWidthWithoutName(rootdir.opt.ViewFields)
	}
	idxmap := make(map[string]string)
	for _, rp := range rootdir.relpaths {
		if rootdir.opt.IsRelPathNotView(rp) {
			continue
		}
		var (
			level int
			idx   = idxmap[rp]
			cidx  = " [" + paw.Cvalue.Sprintf("%s", idx) + "] "
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

		des, _ := cur.ReadDirAll()
		if len(des) < 1 {
			tnd--
			continue
		}

		if level > 0 {
			slevel := paw.Cfield.Sprintf("L%d", level) + cidx
			cur.FprintlnRelPathC(w, pad+slevel, false)
			// fmt.Fprintln(w, cur.RelPathC(pad+slevel, false))
			// FprintRelPath(w, pad, slevel, "", rp, false)
		}

		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, pad)
		}
		ViewFieldName.SetWidth(wdname - wdpad)

		head := vfields.GetHeadFunc(func(i int) *Color {
			if i%2 == 0 {
				return ceven
			} else {
				return codd
			}
		})
		// head := vfields.GetHead(paw.Chdp)
		fmt.Fprintf(w, "%s%v\n", pad, head)
		for _, de := range des {
			var sidx string
			if de.IsDir() {
				if isViewNoDirs {
					nitems--
					continue
				}
				nd++
				sidx = fmt.Sprintf("D%-[1]*[2]d", wdidx, nd)
				idxmap[de.RelPath()] = "D" + cast.ToString(nd)
			} else {
				if isViewNoFiles {
					nitems--
					continue
				}
				nf++
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
		// totalsize += size
		if rootdir.opt.Depth != 0 {
			cur.FprintlnSummaryC(w, pad, wdstty, false)
			// fmt.Fprintln(w, cur.SummaryC(pad, wdstty, false))
		}
		if nd+nf < nitems {
			FprintBanner(w, "", "-", wdstty)
		}
		if rootdir.opt.Depth == 0 {
			break
		}
		ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
	}

	FprintBanner(w, "", "=", wdstty)
	rootdir.FprintlnSummaryC(w, "", wdstty, true)
	// fmt.Fprintln(w, rootdir.SummaryC("", wdstty, true))
}
