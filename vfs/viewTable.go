package vfs

import (
	"fmt"
	"os"
	"strings"

	"io"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/tabulate"
	"github.com/sirupsen/logrus"
)

func (v *VFS) ViewTable(w io.Writer) {
	VFSViewTable(w, v)
}

func VFSViewTable(w io.Writer, v *VFS) {
	paw.Logger.WithFields(logrus.Fields{"View type": v.opt.ViewType}).Debug("view...")

	tmpfields := v.opt.ViewFields
	if v.opt.ViewFields&ViewFieldNo == 0 {
		v.opt.ViewFields = ViewFieldNo | v.opt.ViewFields
	}

	hasX, isViewNoDirs, isViewNoFiles := v.hasX_NoDir_NoFiles()
	viewTableByTabulate(w, v.RootDir(), hasX, isViewNoDirs, isViewNoFiles)
	// viewTable(w, v.RootDir(), hasX, isViewNoDirs, isViewNoFiles)

	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
	v.opt.ViewFields = tmpfields
}

func viewTableByTabulate(w io.Writer, cur *Dir, hasX, isViewNoDirs, isViewNoFiles bool) {
	isNoColor := color.NoColor
	// color.NoColor = true
	var (
		vfields          = cur.opt.ViewFields
		wdstty           = sttyWidth - 2
		tnd, tnf, nitems = cur.NItems()
		wdidx            = GetMaxWidthOf(tnd, tnf)
		// wdidx          = ViewFieldNo.Width()
		nd, nf       int
		roothead     = GetRootHeadC(cur, wdstty)
		totalsize    int64
		_MIN_PADDING = tabulate.MIN_PADDING
	)
	tabulate.MIN_PADDING = 2

	fmt.Fprintf(w, "%v\n", roothead)
	FprintBanner(w, "", "=", wdstty)

	vfields.ModifyWidths(cur)
	ViewFieldNo.SetWidth(wdidx + 1)
	ViewFieldName.SetWidth(ViewFieldName.Width() - vfields.Count()*2)
	_Widths := vfields.Widths()
	heads := vfields.GetHeadA(nil)
	idxmap := make(map[string]string)
	for _, rp := range cur.RelPaths() {
		if cur.opt.IsRelPathNotView(rp) {
			continue
		}
		var (
			curnd, curnf int
			size         int64
			idx          = idxmap[rp]
			cidx         = paw.Cfield.Sprintf(idx) + " "
			// idx          = fmt.Sprintf("G%-[1]*[2]d ", wdidx, i)
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

		if rp != "." {
			FprintRelPath(w, "", "", cidx, rp, false)
		}
		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, "")
		}

		rows := make([][]string, 0)
		// rows := make([][]string, 0, len(des))
		for _, de := range des {
			jdx := ""
			if de.IsDir() {
				if isViewNoFiles {
					nitems--
					continue
				}
				nd++
				curnd++
				jdx = fmt.Sprintf("D%d", nd)
				idxmap[de.RelPath()] = jdx
			} else {
				if isViewNoDirs {
					nitems--
					continue
				}
				nf++
				curnf++
				size += de.Size()
				jdx = fmt.Sprintf("F%d", nf)
			}
			ViewFieldNo.SetValue(jdx)
			// values := vfields.GetValuesS(de)
			values := vfields.GetValuesC(de)
			wdname := paw.StringWidth(de.Field(ViewFieldName))
			if wdname < ViewFieldName.Width() {
				values[len(values)-1] += paw.Spaces(ViewFieldName.Width() - wdname)
			}
			rows = append(rows, values)
			if hasX {
				xattrs := de.Xattibutes()
				if len(xattrs) > 0 {
					nv := len(values)
					cxs := make([]string, nv)
					for _, x := range xattrs {
						sp := paw.Spaces(ViewFieldName.Width() - 2 - paw.StringWidth(x))
						cxs[nv-1] = paw.Cxbp.Sprint("@ ") + paw.Cxap.Sprint(x) + sp
						// cxs[nv-1] = "@ " + x + sp
						rows = append(rows, cxs)
					}
				}
			}
		}
		t := tabulate.Create(rows)
		t.EnableRawOut(_Widths)
		t.SetHeaders(heads)
		// t.SetAlign("left")
		t.SetDenseMode()

		renders := t.Render("simple")
		fmt.Fprint(w, renders)

		totalsize += size
		if cur.opt.Depth != 0 {
			FprintDirSummary(w, "", curnd, curnf, size, wdstty)
		}
		if nd+nf < nitems {
			FprintBanner(w, "", "-", wdstty)
		}
		if cur.opt.Depth == 0 {
			break
		}
	}

	FprintBanner(w, "", "=", wdstty)
	FprintTotalSummary(w, "", nd, nf, totalsize, wdstty)

	tabulate.MIN_PADDING = _MIN_PADDING
	color.NoColor = isNoColor
}

func viewTable(w io.Writer, cur *Dir, hasX, isViewNoDirs, isViewNoFiles bool) (totalsize int64) {
	var (
		vfields        = cur.opt.ViewFields
		fields         = vfields.GetModifyWidthsNoGitFields(cur)
		wdstty         = sttyWidth - 2
		tnd, _, nitems = cur.NItems()
		wdidx          = ViewFieldNo.Width()
		nd, nf         int
		// wdmeta         = 0
		roothead = GetRootHeadC(cur, wdstty)
		banner   = strings.Repeat("-", wdstty)
		tf       = &paw.TableFormat{
			Fields:            make([]string, 0, len(fields)),
			LenFields:         make([]int, 0, len(fields)),
			Aligns:            make([]paw.Align, 0, len(fields)),
			Padding:           "",
			IsWrapped:         false,
			IsColorful:        true,
			XAttributeSymbol:  paw.XAttrSymbol,
			XAttributeSymbol2: paw.XAttrSymbol2,
		}
		values []interface{}
	)

	for _, fd := range fields {
		tf.Fields = append(tf.Fields, fd.Name())
		tf.LenFields = append(tf.LenFields, fd.Width())
		tf.Aligns = append(tf.Aligns, fd.Align())
	}

	tf.Prepare(w)
	errmsg := cur.Errors("")
	if len(errmsg) > 0 {
		errmsg = strings.TrimSuffix(errmsg, "\n")
		tf.SetBeforeMessage(fmt.Sprintf("%v\n%v", roothead, errmsg))
	} else {
		tf.SetBeforeMessage(roothead)
	}

	tf.PrintSart()

	for i, rp := range cur.RelPaths() {
		if cur.opt.IsRelPathNotView(rp) {
			continue
		}
		var (
			curnd, curnf int
			size         int64
			idx          = fmt.Sprintf("D%-[1]*[2]d ", wdidx, i)
			// cidx         = paw.Cfield.Sprint(idx)
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

		if rp != "." {
			tf.PrintLineln(GetRelPath("", idx, rp, false))
		}
		if len(cur.errors) > 0 {
			errmsg := cur.Errors("")
			tf.PrintLine(errmsg)
		}
		if rp != "." {
			// tf.PrintMiddleSepLine()
			tf.PrintHeads()
		}
		for _, de := range des {
			jdx := ""
			if de.IsDir() {
				if isViewNoFiles {
					nitems--
					continue
				}
				nd++
				curnd++
				jdx = fmt.Sprintf("D%d", nd)
			} else {
				if isViewNoDirs {
					nitems--
					continue
				}
				nf++
				curnf++
				size += de.Size()
				jdx = fmt.Sprintf("F%d", nf)
			}
			ViewFieldNo.SetValue(jdx)
			values, tf.FieldsColorString, tf.Colors = vfields.GetAllValues(de)
			tf.PrintRow(values...)
			if hasX {
				xattrs := de.Xattibutes()
				if len(xattrs) > 0 {
					nfields := len(fields)
					cxvalues := make([]string, nfields)
					values := make([]interface{}, nfields)
					for _, x := range xattrs {
						values[nfields-1] = paw.XAttrSymbol + x
						cxvalues[nfields-1] =
							paw.Cxbp.Sprint(paw.XAttrSymbol) +
								paw.Cxap.Sprint(x)
						tf.FieldsColorString = cxvalues
						tf.PrintRow(values...)
					}
				}
			}
		}
		totalsize += size
		tf.PrintMiddleSepLine()
		tf.PrintLineln(dirSummary("", curnd, curnf, size, wdstty))
		if nd+nf < nitems {
			tf.PrintLineln(banner)
		}
		if cur.opt.Depth == 0 {
			break
		}
	}

	tf.SetAfterMessage(totalSummary("", nd, nf, totalsize, wdstty))
	tf.PrintEnd()
	return totalsize
}
