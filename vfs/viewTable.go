package vfs

import (
	"fmt"

	"io"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func (v *VFS) ViewTable(w io.Writer, fields []ViewField, hasX bool) {
	paw.Logger.Info("[vfs] LevelView...")

	cur := v.RootDir()

	if fields == nil {
		fields = DefaultViewFields
	}
	fields = checkFieldsHasGit(fields, cur.git.NoGit)

	modFieldWidths(v, fields)

	// cdir, cname := filepath.Split(cur.Path())
	// cdir = cdirp.Sprint(cdir)
	// cname = cdip.Sprint(cname)
	// fmt.Fprintf(w, "Root: %v\n", cdir+cname)
	// fprintBanner(w, "", "=", sttyWidth-2)

	nd, nf := cur.NItems()
	wdidx := paw.MaxInt(len(fmt.Sprint(nd)), len(fmt.Sprint(nf))) + 1

	ViewFieldNo.SetWidth(wdidx)
	fields = append([]ViewField{ViewFieldNo}, fields...)

	head := getPFHeadS(chdp, fields...)

	viewTable(w, cur, head, wdidx, fields, hasX)
	// size := viewTable(w, cur, head, wdidx, fields, hasX)

	// fprintBanner(w, "", "=", sttyWidth-2)
	// fprintTotalSummary(w, "", nd, nf, size, sttyWidth-2)
}

func viewTable(w io.Writer, cur *Dir, head string, wdidx int, fields []ViewField, hasX bool) (totalsize int64) {
	var (
		wdstty   = sttyWidth - 2
		tnd, tnf = cur.NItems()
		nitems   = tnd + tnf
		nd, nf   int
		wdmeta   = 0
		// cdir, cname = filepath.Split(cur.Path())
		roothead = getRootHeadC(cur, wdstty)
		spNo     = paw.Spaces(wdidx + 1)
		banner   = strings.Repeat("-", wdstty)
	)
	for _, fd := range fields {
		if fd&ViewFieldName == ViewFieldName {
			continue
		}
		wdmeta += fd.Width() + 1
	}
	ViewFieldName.SetWidth(wdstty - wdmeta)
	heads := make([]string, 0, len(fields))
	aligns := make([]paw.Align, 0, len(fields))
	widths := make([]int, 0, len(fields))
	for _, fd := range fields {
		heads = append(heads, fd.Name())
		widths = append(widths, fd.Width())
		aligns = append(aligns, fd.Align())
	}

	tf := &paw.TableFormat{
		Fields:            heads,
		LenFields:         widths,
		Aligns:            aligns,
		Padding:           "",
		IsWrapped:         false,
		IsColorful:        true,
		XAttributeSymbol:  paw.XAttrSymbol,
		XAttributeSymbol2: paw.XAttrSymbol2,
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

	if hasX {
		for _, fd := range fields {
			if fd&ViewFieldName == ViewFieldName {
				continue
			}
			wdmeta += fd.Width() + 1
		}
	}
	for i, rp := range cur.RelPaths() {
		var (
			curnd, curnf int
			size         int64
			idx          = fmt.Sprintf("G%-[1]*[2]d ", wdidx, i)
			// widx         = paw.StringWidth(idx)
			cidx = cdip.Sprint(idx)
		)

		cur, err := cur.getDir(rp)
		if err != nil {
			paw.Logger.WithFields(logrus.Fields{
				"rp": rp,
			}).Fatal(err)
		}

		des, _ := cur.ReadDir(-1)
		cur.resetIdx()
		if len(des) < 1 {
			tnd--
			continue
		}

		cdir, cname := filepath.Split(rp)
		cdir = cdirp.Sprint(cdir)
		if rp != "." {
			cdir = cdirp.Sprint("./") + cdir
		}
		cname = cdip.Sprint(cname)
		if rp != "." {
			tf.PrintLineln(cidx + cdir + cname)
		}
		if len(cur.errors) > 0 {
			errmsg := cur.Errors("")
			tf.PrintLine(errmsg)
		}
		if rp != "." {
			tf.PrintHeads()
		}
		for _, child := range des {
			jdx := ""
			// cjdx := ""
			f, isFile := child.(*File)
			if isFile {
				// (*nf)++
				nf++
				curnf++
				size += f.Size()
				jdx = fmt.Sprintf("F%d", nf)
				// cjdx = cfip.Sprintf("%[1]*[2]s", wdidx, jdx)
				ViewFieldNo.SetValue(jdx)
				values := setTableValues(child, tf, fields)
				tf.PrintRow(values...)
				if hasX {
					xattrs := f.Xattibutes()
					if len(xattrs) > 0 {
						cxvalues := make([]string, len(fields))
						for _, x := range xattrs {
							cxvalues[len(fields)-1] =
								cxbp.Sprint(tf.XAttributeSymbol) +
									cxap.Sprint(x)
							tf.FieldsColorString = cxvalues
							tf.PrintRow(nil)
						}
					}
				}
			} else {
				nd++
				curnd++
				d := child.(*Dir)
				jdx = fmt.Sprintf("D%d", nd)
				ViewFieldNo.SetValue(jdx)
				values := setTableValues(child, tf, fields)
				tf.PrintRow(values...)
				if hasX {
					xattrs := d.Xattibutes()
					if len(xattrs) > 0 {
						cxvalues := make([]string, len(fields))
						for _, x := range xattrs {
							cxvalues[len(fields)-1] =
								cxbp.Sprint(tf.XAttributeSymbol) +
									cxap.Sprint(x)
							tf.FieldsColorString = cxvalues
							tf.PrintRow(nil)
						}
					}
				}
			}
		}
		totalsize += size
		tf.PrintLineln(dirSummary(spNo, curnd, curnf, size, wdstty))
		if nd+nf < nitems {
			tf.PrintLineln(banner)
		}
	}

	tf.SetAfterMessage(totalSummary("", nd, nf, totalsize, wdstty))
	tf.PrintEnd()
	return totalsize
}

func setTableValues(de DirEntryX, tf *paw.TableFormat, fields []ViewField) (values []interface{}) {
	values = make([]interface{}, 0, len(fields))
	cvalues := make([]string, 0, len(fields))
	colors := make([]*color.Color, 0, len(fields))
	f, isFile := de.(*File)
	if isFile {
		for _, field := range fields {
			values = append(values, f.Field(field))
			cvalues = append(cvalues, f.FieldC(field))
			if field&ViewFieldName != 0 {
				colors = append(colors, f.LSColor())
			} else {
				colors = append(colors, field.Color())
			}
		}
	} else {
		d := de.(*Dir)
		for _, field := range fields {
			values = append(values, d.Field(field))
			cvalues = append(cvalues, d.FieldC(field))
			colors = append(colors, field.Color())
		}
	}
	tf.Colors = colors
	tf.FieldsColorString = cvalues
	return values
}

// func levelView(w io.Writer, cur *Dir, root, head string, level, wdidx int, fields []PDFieldFlag, nd, nf *int) {
// 	des, _ := cur.ReadDir(-1)
// 	cur.resetIdx()
// 	if len(des) == 0 {
// 		return
// 	}
// 	pad := paw.Spaces(level * 3)

// 	cdir, cname := filepath.Split(cur.RelPath())
// 	cdir = cdirp.Sprint(cdir)
// 	cname = cdip.Sprint(cname)
// 	fmt.Fprintf(w, "%sL%d: %v\n", pad, level, cdir+cname)
// 	fmt.Fprintf(w, "%s%#v\n", pad, cur.rdirs)

// 	fmt.Fprintf(w, "%s%v\n", pad, head)

// 	if len(cur.errors) > 0 {
// 		cur.FprintErrors(os.Stderr, pad)
// 	}
// 	for _, child := range des {
// 		f, isFile := child.(*File)
// 		if isFile {
// 			(*nf)++
// 			sidx := cfip.Sprintf("F%-[1]*[2]d", wdidx, *nf)
// 			fmt.Fprintf(w, "%s%s ", pad, sidx)
// 			for _, field := range fields {
// 				fmt.Fprintf(w, "%v ", f.FieldC(field, nil))
// 			}
// 			fmt.Println()
// 		} else {
// 			(*nd)++
// 			sidx := cdip.Sprintf("D%-[1]*[2]d", wdidx, *nd)
// 			d := child.(*Dir)
// 			fmt.Fprintf(w, "%s%s ", pad, sidx)
// 			for _, field := range fields {
// 				fmt.Fprintf(w, "%v ", d.FieldC(field, nil))
// 			}
// 			fmt.Println()
// 			// ndd, nff := d.NItems()
// 			if len(cur.rdirs) > 0 {
// 				levelView(w, d, root, head, level+1, wdidx, fields, nd, nf)
// 			}
// 		}
// 	}
// }
