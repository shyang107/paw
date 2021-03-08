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

func (v *VFS) ViewTable(w io.Writer) {
	VFSViewTable(w, v)
}

func VFSViewTable(w io.Writer, v *VFS) {
	paw.Logger.WithFields(logrus.Fields{"View type": v.opt.ViewType}).Debug("view...")

	var (
		cur                               = v.RootDir()
		vfields                           = v.opt.ViewFields
		fields                            []ViewField
		hasX, isViewNoDirs, isViewNoFiles = v.hasX_NoDir_NoFiles()
	)

	if vfields&ViewFieldNo == 0 {
		vfields = ViewFieldNo | vfields
	}

	fields = checkFieldsHasGit(vfields.Fields(), cur.git.NoGit)

	modFieldWidths(cur, fields)

	viewTable(w, cur, fields, hasX, isViewNoDirs, isViewNoFiles)

	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
}

func viewTable(w io.Writer, cur *Dir, fields []ViewField, hasX, isViewNoDirs, isViewNoFiles bool) (totalsize int64) {
	var (
		wdstty   = sttyWidth - 2
		tnd, tnf = cur.NItems()
		wdidx    = ViewFieldNo.Width()
		nitems   = tnd + tnf
		nd, nf   int
		wdmeta   = 0
		roothead = GetRootHeadC(cur, wdstty)
		banner   = strings.Repeat("-", wdstty)
	)

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
		if cur.opt.IsNotViewRelPath(rp) {
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

		cdir, cname := filepath.Split(rp)
		cdir = cdirp.Sprint(cdir)
		cname = cdip.Sprint(cname)
		if rp != "." {
			cdir = cdirp.Sprint("./") + cdir
			tf.PrintLineln(cidx + cdir + cname)
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
			values := setTableValues(de, tf, fields)
			tf.PrintRow(values...)
			if hasX {
				xattrs := de.Xattibutes()
				if len(xattrs) > 0 {
					cxvalues := make([]string, len(fields))
					values := make([]interface{}, len(fields))
					for _, x := range xattrs {
						values[len(fields)-1] = paw.XAttrSymbol + x
						cxvalues[len(fields)-1] =
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
