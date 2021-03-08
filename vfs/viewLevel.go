package vfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

func (v *VFS) ViewLevel(w io.Writer) {
	VFSViewLevel(w, v)
}

func VFSViewLevel(w io.Writer, v *VFS) {
	paw.Logger.WithFields(logrus.Fields{"View type": v.opt.ViewType}).Debug("view...")

	var (
		cur                               = v.RootDir()
		vfields                           = v.opt.ViewFields
		fields                            []ViewField
		hasX, isViewNoDirs, isViewNoFiles = v.hasX_NoDir_NoFiles()
		nd, nf                            = cur.NItems()
		snd, snf                          = fmt.Sprint(nd), fmt.Sprint(nf)
		wdidx                             = paw.MaxInt(len(snd), len(snf))
	)

	paw.Logger.WithFields(logrus.Fields{
		"nd":     nd,
		"nf":     nf,
		"wididx": wdidx,
	}).Debug()

	if vfields&ViewFieldNo == 0 {
		vfields = ViewFieldNo | vfields
	}

	fields = checkFieldsHasGit(vfields.Fields(), cur.git.NoGit)
	modFieldWidths(cur, fields)

	viewLevel(w, cur, fields, hasX, isViewNoDirs, isViewNoFiles)

	ViewFieldName.SetWidth(paw.StringWidth(ViewFieldName.Name()))
}

func viewLevel(w io.Writer, cur *Dir, fields []ViewField, hasX, isViewNoDirs, isViewNoFiles bool) {
	var (
		wdname     = ViewFieldName.Width()
		wdstty     = sttyWidth - 2
		tnd, tnf   = cur.NItems()
		stnd, stnf = fmt.Sprint(tnd), fmt.Sprint(tnf)
		wdidx      = paw.MaxInt(len(stnd), len(stnf))
		nitems     = tnd + tnf
		nd, nf     int
		wdmeta     = 0
		roothead   = GetRootHeadC(cur, wdstty)
		// head      = getPFHeadS(chdp, fields...)
		totalsize int64
	)

	fmt.Fprintf(w, "%v\n", roothead)
	FprintBanner(w, "", "=", wdstty)

	if hasX {
		for _, fd := range fields {
			if fd&ViewFieldName == ViewFieldName {
				continue
			}
			wdmeta += fd.Width() + 1
		}
	}
	for _, rp := range cur.relpaths {
		if cur.opt.IsNotViewRelPath(rp) {
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

		cdir, cname := filepath.Split(rp)
		cname = cdip.Sprint(cname)
		cdir = cdirp.Sprint(cdir)
		if level > 0 {
			cdir = cdirp.Sprint("./") + cdir
		}

		fmt.Fprintf(w, "%sL%d: %v\n", pad, level, cdir+cname)
		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, pad)
		}
		ViewFieldName.SetWidth(wdname - wdpad)
		head := GetPFHeadS(chdp, fields...)
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
			for _, field := range fields {
				fmt.Fprintf(w, "%v ", de.FieldC(field))
			}
			fmt.Println()
			if hasX {
				FprintXattrs(w, wdpad+wdmeta, de.Xattibutes())
			}
		}
		totalsize += size
		fprintDirSummary(w, pad, curnd, curnf, size, wdstty)
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
