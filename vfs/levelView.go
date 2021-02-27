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

func (v *VFS) LevelView(w io.Writer, cur *Dir, fields []PDFieldFlag) {
	paw.Logger.Info("[vfs] LevelView...")

	if cur == nil {
		cur = v.RootDir()
	}

	if fields == nil {
		fields = DefaultPDFieldKeys
	}
	fields = checkFieldsHasGit(fields, cur.git.NoGit)

	modFieldWidths(v, fields)

	cdir, cname := filepath.Split(cur.Path())
	cdir = cdirp.Sprint(cdir)
	cname = cdip.Sprint(cname)
	fmt.Fprintf(w, "Root: %v\n", cdir+cname)
	fprintBanner(w, "", "=", sttyWidth-2)

	nd, nf := cur.NItems()
	wdidx := paw.MaxInt(len(fmt.Sprint(nd)), len(fmt.Sprint(nf)))
	head := chdp.Sprintf("%[1]*[2]s", wdidx+1, "No") + " "
	head += getPFHeadS(chdp, fields...)
	size := levelView(w, cur, head, wdidx, fields)
	// var tnd, tnf int
	// size := levelView(w, cur, cur.path, head, wdidx, fields, &tnd, &tnf)

	fprintBanner(w, "", "=", sttyWidth-2)
	fprintTotalSummary(w, "", nd, nf, size, sttyWidth-2)
}

func levelView(w io.Writer, cur *Dir, head string, wdidx int, fields []PDFieldFlag) (totalsize int64) {
	// func levelView(w io.Writer, cur *Dir, root, head string, wdidx int, fields []PDFieldFlag, nd, nf *int) (totalsize int64) {
	var (
		wdstty   = sttyWidth - 2
		tnd, tnf = cur.NItems()
		nitems   = tnd + tnf
		nd, nf   int
	)
	// paw.Logger.WithField("tnd", tnd).Debug()
	// paw.Logger.WithField("rp", cur.RelPath()).Debug()
	// spew.Dump(cur.relpaths)
	for _, rp := range cur.relpaths {
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
		pad := paw.Spaces(level * 3)
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
		if level > 0 {
			cdir = cdirp.Sprint("./") + cdir
		}
		cname = cdip.Sprint(cname)

		fmt.Fprintf(w, "%sL%d: %v\n", pad, level, cdir+cname)
		if len(cur.errors) > 0 {
			cur.FprintErrors(os.Stderr, pad)
		}
		// fmt.Fprintf(w, "%s%#v\n", pad, cur.rdirs)

		fmt.Fprintf(w, "%s%v\n", pad, head)
		for _, child := range des {
			f, isFile := child.(*File)
			if isFile {
				// (*nf)++
				nf++
				curnf++
				size += f.Size()
				sidx := cfip.Sprintf("F%-[1]*[2]d", wdidx, nf)
				fmt.Fprintf(w, "%s%s ", pad, sidx)
				for _, field := range fields {
					fmt.Fprintf(w, "%v ", f.FieldC(field, nil))
				}
				fmt.Println()
			} else {
				// (*nd)++
				nd++
				curnd++
				d := child.(*Dir)
				sidx := cdip.Sprintf("D%-[1]*[2]d", wdidx, nd)
				fmt.Fprintf(w, "%s%s ", pad, sidx)
				for _, field := range fields {
					fmt.Fprintf(w, "%v ", d.FieldC(field, nil))
				}
				fmt.Println()
				// if len(d.relpaths) > 0 {
				// 	paw.Logger.WithField("rp", d.RelPath()).Debug()
				// 	spew.Dump(d.relpaths)
				// }
			}
		}
		totalsize += size
		fprintDirSummary(w, pad, curnd, curnf, size, wdstty)
		if nd+nf < nitems {
			fprintBanner(w, "", "-", wdstty)
		}
	}
	return totalsize
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
