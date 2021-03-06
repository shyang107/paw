package vfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shyang107/paw"
)

// View excutes view operation of VFS and all needed arguments to view in VFS.opt.
func (v *VFS) View(w io.Writer) {
	if view, ok := ViewTypeFuncs[v.opt.ViewType]; ok {
		view(w, v)
	} else {
		VFSViewList(w, v)
	}
}

func (v *VFS) hasX_NoDir_NoFiles() (hasX, isViewNoDirs, isViewNoFiles bool) {
	var (
		vt = v.opt.ViewType
	)
	hasX = vt&ViewExtended != 0
	isViewNoDirs = vt&ViewNoDirs != 0
	isViewNoFiles = vt&ViewNoFiles != 0
	return hasX, isViewNoDirs, isViewNoFiles
}

func (v *VFS) hasList_hasX() (hasList, hasX bool) {
	var (
		vt = v.opt.ViewType
	)
	hasList = vt&_ViewList != 0
	hasX = vt&ViewExtended != 0
	return hasList, hasX
}

// func (v *VFS) DumpFS(w io.Writer) {
// 	color.NoColor = true
// 	v.View(w)
// 	color.NoColor = paw.NoColor
// }

func (v *VFS) Dump(w io.Writer) {
	paw.Logger.Debug()
	// color.NoColor = true
	var (
		cur     = v.RootDir()
		root    = cur.path
		vfields = cur.opt.ViewFields
		// vopt = *v.Option()
		wdstty   = sttyWidth - 2
		roothead = "Root: " + PathTo(cur, &PathToOption{true, nil, PRTPathToLink})
		tnd, tnf int
		tsize    int64
	)
	// vfields.ModifyWidths(cur)
	ViewFieldSize.SetWidth(7)
	head := vfields.GetHead(paw.Chdp)
	fmt.Fprintf(w, "%v\n", roothead)
	FprintBanner(w, "", "=", wdstty)

	hasX, isViewNoDirs, isViewNoFiles := v.hasX_NoDir_NoFiles()
	_dump(w, cur, root, 0, head, hasX, isViewNoDirs, isViewNoFiles, &tnd, &tnf, &tsize)
	// color.NoColor = paw.NoColor
	fmt.Fprintln(w, totalSummary("", tnd, tnf, tsize, wdstty))
}

func _dump(w io.Writer, cur *Dir, root string, level int, head string, hasX, isViewNoDirs, isViewNoFiles bool, nd, nf *int, size *int64) {
	var (
		dpath = cur.Path()
		git   = cur.git
		skip  = cur.opt.Skips
	)
	if !cur.opt.IsForceRecurse &&
		cur.opt.Depth > 0 &&
		level > cur.opt.Depth {
		return
	}

	// if cur.RelPath() != "." {
	// 	cur.FprintlnRelPathC(w, "", false)
	// }
	des, _ := os.ReadDir(dpath)
	if len(des) == 0 {
		return
	}
	// fmt.Fprintln(w, head)
	for _, de := range des {
		path := filepath.Join(dpath, de.Name())
		_, err := os.Lstat(path)
		if err != nil {
			if cur.errors == nil {
				cur.errors = []error{}
			}
			cur.errors = append(cur.errors, err)
			// cur.errors = append(cur.errors, &fs.PathError{
			// 	Op:   "os", // "buildFS",
			// 	Path: path,
			// 	Err:  err,
			// })
			continue
		}
		// relpath, _ := filepath.Rel(root, path)
		// xattrs, _ := GetXattr(path)
		var child DirEntryX
		if !de.IsDir() {
			child, err = NewFile(path, root, git)
		} else {
			child, err = NewDir(path, root, git, cur.opt)
		}
		if err != nil {
			cur.AddErrors(err)
			continue
		}
		if skip.IsSkip(child) {
			continue
		}
		if child.IsDir() {
			(*nd)++
		} else {
			(*nf)++
			if de.Type().IsRegular() {
				(*size) += child.Size()
			}
		}
		cur.children[de.Name()] = child
		_dumpPrint(w, child, cur.opt.ViewFields, hasX)
		// paw.Logger.WithFields(logrus.Fields{
		// 	"name":  child.Name(),
		// 	"IsDir": child.IsDir(),
		// 	"depth": cur.opt.Depth,
		// }).Trace()
		if cur.opt.IsForceRecurse {
			if child.IsDir() {
				_dump(w, child.(*Dir), root, 0, head, hasX, isViewNoDirs, isViewNoFiles, nd, nf, size)
				// buildFS(child.(*Dir), root, 0)
			}
		} else {
			if cur.opt.Depth != 0 && child.IsDir() {
				_dump(w, child.(*Dir), root, level+1, head, hasX, isViewNoDirs, isViewNoFiles, nd, nf, size)
				// buildFS(child.(*Dir), root, level+1)
			}
		}
	}
}

func _dumpPrint(w io.Writer, de DirEntryX, vfields ViewField, hasX bool) {
	cmeta := vfields.RowStringXNameC(de)
	fmt.Fprintf(w, "%v ", cmeta)
	cname := PathTo(de, &PathToOption{true, nil, PRTRelPathToLink})
	fmt.Fprintf(w, "%v\n", cname)
	xs := de.Xattibutes()
	if hasX && len(xs) > 0 {
		pad := paw.Spaces(paw.StringWidth(paw.StripANSI(cmeta)))
		for _, x := range xs {
			x = paw.Cxbp.Sprint("@ ") + paw.Cxap.Sprint(x)
			fmt.Fprintln(w, pad, x)
		}
	}
}
