package vfs

import "io"

func (v *VFS) View(w io.Writer, fields []ViewField, viewType ViewType) {
	switch viewType {
	case ViewList:
		v.ViewList(w, fields, false)
	case ViewListX:
		v.ViewList(w, fields, true)
	// case ViewTree:
	// case ViewTreeX:
	// case ViewListTree:
	// case ViewListTreeX:
	case ViewLevel:
		v.ViewLevel(w, fields, false)
	case ViewLevelX:
		v.ViewLevel(w, fields, true)
	case ViewTable:
		v.ViewTable(w, fields, false)
	case ViewTableX:
		v.ViewTable(w, fields, true)
	// case ViewClassify:
	default:
		v.ViewList(w, fields, false)
	}
}
