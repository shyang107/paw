package vfs

import (
	"fmt"
	"os"
	"strings"

	"io"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/cast"
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

func viewTableByTabulate(w io.Writer, rootdir *Dir, hasX, isViewNoDirs, isViewNoFiles bool) {
	isNoColor := color.NoColor
	// color.NoColor = true
	var (
		vfields          = rootdir.opt.ViewFields
		wdstty           = sttyWidth - 2
		tnd, tnf, nitems = rootdir.NItems(true)
		wdidx            = GetMaxWidthOf(tnd, tnf)
		// nd, nf           int
		tsize        int64
		count        int
		roothead     = GetRootHeadC(rootdir, wdstty)
		_MIN_PADDING = tabulate.MIN_PADDING
		// cevenH           = paw.CloneColor(paw.CEven).Add(paw.EXAColors["bgpmpt"]...)
		// coddH            = paw.CloneColor(paw.COdd).Add(paw.EXAColors["bgpmpt"]...)
	)
	tabulate.MIN_PADDING = 2
	tnd, tnf = 0, 0

	fmt.Fprintf(w, "%v\n", roothead)
	// FprintBanner(w, "", "=", wdstty)

	// ViewFieldNo.SetWidth(wdidx + 1)
	vfields.ModifyWidths(rootdir)
	ViewFieldName.SetWidth(ViewFieldName.Width() - vfields.Count()*2)
	_Widths := vfields.Widths()
	// heads := vfields.GetHeadFuncA(func(i int) *Color {
	// 	if i%2 == 0 {
	// 		return cevenH
	// 	} else {
	// 		return coddH
	// 	}
	// })
	heads := vfields.GetHeadA(paw.Cpmpt)
	idxmap := make(map[string]string)
	for _, rp := range rootdir.RelPaths() {
		if rootdir.opt.IsRelPathNotView(rp) {
			continue
		}
		var (
			idx  = idxmap[rp]
			cidx = "[" + paw.Cvalue.Sprintf(idx) + "] "
		)

		paw.Logger.WithFields(logrus.Fields{"rp": rp}).Trace("getDir")
		cur, err := rootdir.getDir(rp)
		if err != nil {
			paw.Logger.WithFields(logrus.Fields{"rp": rp}).Fatal(err)
		}

		if rp != "." {
			if isViewNoDirs {
				cur.FprintlnRelPathC(w, "", false)
			} else {
				cur.FprintlnRelPathC(w, cidx, false)
			}
		}
		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, "", false)
		}

		des, _ := cur.ReadDirAll()
		if len(des) < 1 {
			continue
		}

		var (
			curnd, curnf  int
			size          int64
			vnitems       = nitems
			sidx, renders string
			rows          = make([][]string, 0)
			xrows         [][]string
			values        []string
			wdname        int
			t             *tabulate.Tabulate
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
		// rows := make([][]string, 0, len(des))
		for _, de := range des {
			if isSkipViewItem(de, isViewNoDirs, isViewNoFiles, &nitems, &curnd, &curnf, &size) {
				continue
			}
			count++
			if de.IsDir() {
				sidx = fmt.Sprintf("D%-[1]*[2]d", wdidx, tnd+curnd)
				idxmap[de.RelPath()] = "D" + cast.ToString(tnd+curnd)
			} else {
				sidx = fmt.Sprintf("F%-[1]*[2]d", wdidx, tnf+curnf)
			}
			ViewFieldNo.SetValue(sidx)
			values = vfields.GetValuesC(de)
			wdname = paw.StringWidth(de.FieldC(ViewFieldName))
			if wdname < ViewFieldName.Width() {
				values[len(values)-1] += paw.Spaces(ViewFieldName.Width() - wdname)
			}
			rows = append(rows, values)
			if hasX {
				xrows = vfields.XattibutesRowsC(de)
				rows = append(rows, xrows...)
			}
		}
		t = tabulate.Create(rows)
		t.EnableRawOut(_Widths)
		t.SetHeaders(heads)
		// t.SetAlign("left")
		t.SetDenseMode()

		renders = t.Render("simple")
		fmt.Fprint(w, renders)

		tnd += curnd
		tnf += curnf
		tsize += size
		if rootdir.opt.Depth != 0 {
			fmt.Fprintln(w, dirSummary("", curnd, curnf, size, wdstty))
			// cur.FprintlnSummaryC(w, "", wdstty, false)
		}
		if count < nitems {
			FprintBanner(w, "", "-", wdstty)
		}
	BAN:
		if rootdir.opt.Depth == 0 {
			break
		}
	}

	FprintBanner(w, "", "=", wdstty)
	fmt.Fprintln(w, totalSummary("", tnd, tnf, tsize, wdstty))
	// rootdir.FprintlnSummaryC(w, "", wdstty, true)

	tabulate.MIN_PADDING = _MIN_PADDING
	color.NoColor = isNoColor
}

func viewTable(w io.Writer, rootdir *Dir, hasX, isViewNoDirs, isViewNoFiles bool) {
	var (
		vfields        = rootdir.opt.ViewFields
		fields         = vfields.GetModifyWidthsNoGitFields(rootdir)
		wdstty         = sttyWidth - 2
		tnd, _, nitems = rootdir.NItems(true)
		wdidx          = ViewFieldNo.Width()
		nd, nf         int
		// wdmeta         = 0
		roothead = GetRootHeadC(rootdir, wdstty)
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
	errmsg := rootdir.Errors("", false)
	if len(errmsg) > 0 {
		errmsg = strings.TrimSuffix(errmsg, "\n")
		tf.SetBeforeMessage(fmt.Sprintf("%v\n%v", roothead, errmsg))
	} else {
		tf.SetBeforeMessage(roothead)
	}

	tf.PrintSart()

	for i, rp := range rootdir.RelPaths() {
		if rootdir.opt.IsRelPathNotView(rp) {
			continue
		}
		var (
			idx = fmt.Sprintf("D%-[1]*[2]d ", wdidx, i)
			// cidx         = paw.Cfield.Sprint(idx)
		)

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
			tf.PrintLineln(cur.RelPathC(idx, false))
			// tf.PrintLineln(getRelPath("", idx, rp, false))
		}
		if len(cur.errors) > 0 {
			errmsg := cur.Errors("", false)
			tf.PrintLine(errmsg)
		}

		des, _ := cur.ReadDirAll()
		if len(des) < 1 {
			tnd--
			continue
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
				jdx = fmt.Sprintf("D%d", nd)
			} else {
				if isViewNoDirs {
					nitems--
					continue
				}
				nf++
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
		tf.PrintMiddleSepLine()
		tf.PrintLineln(cur.SummaryC("", wdstty, false))
		if nd+nf < nitems {
			tf.PrintLineln(banner)
		}
		if rootdir.opt.Depth == 0 {
			break
		}
	}

	tf.SetAfterMessage(rootdir.SummaryC("", wdstty, true))
	tf.PrintEnd()
}
