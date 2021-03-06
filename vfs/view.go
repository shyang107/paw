package vfs

import "io"

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
