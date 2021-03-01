package vfs

import "io"

func (v *VFS) View(w io.Writer, fields []ViewField, viewType ViewType) {
	switch viewType {
	case ViewList:
		v.ViewList(w, fields, false)
	case ViewListX:
		v.ViewList(w, fields, true)
	case ViewTree:
		v.ViewListTree(w, fields, false, false)
	case ViewTreeX:
		v.ViewListTree(w, fields, true, false)
	case ViewListTree:
		v.ViewListTree(w, fields, false, true)
	case ViewListTreeX:
		v.ViewListTree(w, fields, true, true)
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
