package vfs

import (
	"errors"
	"io"
)

type ViewType int

const (
	// PExtendView is the option to add extended attributes view using in PrintDir
	ViewExtended ViewType = 1 << iota // 1 << 0 which is 00000001
	// ViewList is the option of list view using in PrintDir
	ViewList
	// ViewListTree is the option of combining list & tree view using in PrintDir
	ViewListTree
	// ViewTree is the option of tree view using in PrintDir
	ViewTree
	// ViewLevel is the option of level view using in PrintDir
	ViewLevel
	// ViewTable is the option of table view using in PrintDir
	ViewTable
	// ViewClassify display type indicator by file names (like as `exa -F` or `exa --classify`) in PrintDir
	ViewClassify

	// ViewListX is the option of list view icluding extend attributes using in PrintDir
	ViewListX = ViewList | ViewExtended
	// ViewListTreeX is the option of combining list & tree view including extend attribute using in PrintDir
	ViewListTreeX = ViewListTree | ViewExtended
	// ViewTreeX is the option of tree view icluding extend atrribute using in PrintDir
	ViewTreeX = ViewTree | ViewExtended

	// ViewLevelX is the option of level view icluding extend attributes using in PrintDir
	ViewLevelX = ViewLevel | ViewExtended

	// ViewTableX is the option of table view icluding extend attributes using in PrintDir
	ViewTableX = ViewTable | ViewExtended
)

var (
	ViewTypeNames = map[ViewType]string{
		ViewList:      "List view",
		ViewListX:     "Extended List view",
		ViewTree:      "Tree view",
		ViewTreeX:     "Extended Tree view",
		ViewLevel:     "Level view",
		ViewLevelX:    "Extended Level view",
		ViewTable:     "Table view",
		ViewTableX:    "Extended Table view",
		ViewListTree:  "List & Tree view",
		ViewListTreeX: "Extended List & Tree view",
		ViewClassify:  "Classify view",
	}

	isViewNoDirs  bool = false
	isViewNoFiles bool = false
)

// ViewDirAndFile enables ViewType showing directories and files (excluding ViewListTree and ViewListTreeX)
func (v ViewType) ViewDirAndFile() {
	isViewNoDirs = false
	isViewNoFiles = false
}

// NoDirs disables ViewType showing directories (excluding ViewListTree and ViewListTreeX)
// 	see examples/vfs
func (v ViewType) NoDirs() ViewType {
	isViewNoDirs = true
	return v
}

// NoFiles disables ViewType showing files (excluding ViewListTree and ViewListTreeX)
// 	see examples/vfs
func (v ViewType) NoFiles() ViewType {
	isViewNoFiles = true
	return v
}

func (v ViewType) String() string {
	if name, ok := ViewTypeNames[v]; ok {
		return name
	} else {
		return "Unknown"
	}
}

// Do  will print out VFS
func (v ViewType) Do(w io.Writer, vfs *VFS) error {
	if vfs == nil {
		return errors.New("not a valid VFS")
	}
	vfs.opt.ViewType = v
	vfs.View(w)
	return nil
}

// View excutes view operation of VFS and all needed arguments to view in VFS.opt.
func (v *VFS) View(w io.Writer) {
	var (
		vfields  = v.opt.ViewFields
		fields   = vfields.Fields()
		viewType = v.opt.ViewType
	)
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
	case ViewClassify:
		v.ViewClassify(w)
	default:
		v.ViewList(w, fields, false)
	}
}
