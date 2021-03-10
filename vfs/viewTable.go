package vfs

import (
	"fmt"

	"io"
	"strings"

	"github.com/shyang107/paw"
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
	viewTable(w, v.RootDir(), hasX, isViewNoDirs, isViewNoFiles)

	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
	v.opt.ViewFields = tmpfields
}

func viewTable(w io.Writer, cur *Dir, hasX, isViewNoDirs, isViewNoFiles bool) (totalsize int64) {
	var (
		vfields        = cur.opt.ViewFields
		fields         = vfields.GetModifyWidthsNoGitFields(cur, cur.git.NoGit)
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
			idx          = fmt.Sprintf("G%-[1]*[2]d ", wdidx, i)
			cidx         = cdip.Sprint(idx)
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
			cpath = cdirp.Sprint("./") + cdir + cname
			tf.PrintLineln(cidx + cpath)
		}
		if len(cur.errors) > 0 {
			errmsg := cur.Errors("")
			tf.PrintLine(errmsg)
		}
		if rp != "." {
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
							cxbp.Sprint(paw.XAttrSymbol) +
								cxap.Sprint(x)
						tf.FieldsColorString = cxvalues
						tf.PrintRow(values...)
					}
				}
			}
		}
		totalsize += size
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

// func setTableValues(de DirEntryX, tf *paw.TableFormat, fields []ViewField) (values []interface{}) {
// 	values = make([]interface{}, 0, len(fields))
// 	cvalues := make([]string, 0, len(fields))
// 	colors := make([]*color.Color, 0, len(fields))
// 	for _, field := range fields {
// 		values = append(values, de.Field(field))
// 		cvalues = append(cvalues, de.FieldC(field))
// 		if field&ViewFieldName != 0 {
// 			colors = append(colors, de.LSColor())
// 		} else {
// 			colors = append(colors, field.Color())
// 		}
// 	}
// 	tf.Colors = colors
// 	tf.FieldsColorString = cvalues
// 	return values
// }
