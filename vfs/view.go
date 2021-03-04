package vfs

import (
	"errors"
	"io"

	"github.com/shyang107/paw"
)

type ViewType int

const (

	// ViewList is the option of list view using in PrintDir
	ViewList ViewType = 1 << iota // 1 << 0 which is 00000001

	// ViewTree is the option of tree view using in PrintDir
	ViewTree
	// ViewLevel is the option of level view using in PrintDir
	ViewLevel
	// ViewTable is the option of table view using in PrintDir
	ViewTable
	// ViewClassify display type indicator by file names (like as `exa -F` or `exa --classify`) in PrintDir
	ViewClassify
	// PExtendView is the option to add extended attributes view using in PrintDir

	_ViewList

	_ViewExtended
	_ViewNoDirs
	_ViewNoFiles

	// ViewListTree is the option of combining list & tree view using in PrintDir
	ViewListTree = ViewTree | _ViewList

	// ViewListX is the option of list view icluding extend attributes using in PrintDir
	ViewListX = ViewList | _ViewExtended
	// ViewListTreeX is the option of combining list & tree view including extend attribute using in PrintDir
	ViewListTreeX = ViewListTree | _ViewExtended
	// ViewTreeX is the option of tree view icluding extend atrribute using in PrintDir
	ViewTreeX = ViewTree | _ViewExtended

	// ViewLevelX is the option of level view icluding extend attributes using in PrintDir
	ViewLevelX = ViewLevel | _ViewExtended

	// ViewTableX is the option of table view icluding extend attributes using in PrintDir
	ViewTableX = ViewTable | _ViewExtended

	ViewListNoDirs     = ViewList | _ViewNoDirs
	ViewLevelNoDirs    = ViewLevel | _ViewNoDirs
	ViewTableNoDirs    = ViewTable | _ViewNoDirs
	ViewClassifyNoDirs = ViewClassify | _ViewNoDirs

	ViewListNoFiles     = ViewList | _ViewNoFiles
	ViewLevelNoFiles    = ViewLevel | _ViewNoFiles
	ViewTableNoFiles    = ViewTable | _ViewNoFiles
	ViewClassifyNoFiles = ViewClassify | _ViewNoFiles

	ViewListXNoDirs  = ViewList | _ViewExtended | _ViewNoDirs
	ViewLevelXNoDirs = ViewLevel | _ViewExtended | _ViewNoDirs
	ViewTableXNoDirs = ViewTable | _ViewExtended | _ViewNoDirs

	ViewListXNoFiles  = ViewList | _ViewExtended | _ViewNoFiles
	ViewLevelXNoFiles = ViewLevel | _ViewExtended | _ViewNoFiles
	ViewTableXNoFiles = ViewTable | _ViewExtended | _ViewNoFiles
)

var (
	ViewTypeNames = map[ViewType]string{
		ViewList:            "List view",
		ViewListX:           "Extended List view",
		ViewTree:            "Tree view",
		ViewTreeX:           "Extended Tree view",
		ViewLevel:           "Level view",
		ViewLevelX:          "Extended Level view",
		ViewTable:           "Table view",
		ViewTableX:          "Extended Table view",
		ViewListTree:        "List & Tree view",
		ViewListTreeX:       "Extended List & Tree view",
		ViewClassify:        "Classify view",
		ViewListNoDirs:      "List view (no dirs)",
		ViewListXNoDirs:     "Extended List view (no dirs)",
		ViewLevelNoDirs:     "Level view (no dirs)",
		ViewLevelXNoDirs:    "Extended Level view (no dirs)",
		ViewTableNoDirs:     "Table view (no dirs)",
		ViewTableXNoDirs:    "Extended Table view (no dirs)",
		ViewClassifyNoDirs:  "Classify view (no dirs)",
		ViewListNoFiles:     "List view (no files)",
		ViewListXNoFiles:    "Extended List view (no files)",
		ViewLevelNoFiles:    "Level view (no files)",
		ViewLevelXNoFiles:   "Extended Level view (no files)",
		ViewTableNoFiles:    "Table view (no files)",
		ViewTableXNoFiles:   "Extended Table view (no files)",
		ViewClassifyNoFiles: "Classify view (no files)",
	}

	ViewTypeFuncs = map[ViewType]func(io.Writer, *VFS){
		ViewList:            VFSViewList,
		ViewListX:           VFSViewList,
		ViewTree:            VFSViewListTree,
		ViewTreeX:           VFSViewListTree,
		ViewListTree:        VFSViewListTree,
		ViewListTreeX:       VFSViewListTree,
		ViewLevel:           VFSViewLevel,
		ViewLevelX:          VFSViewLevel,
		ViewTable:           VFSViewTable,
		ViewTableX:          VFSViewTable,
		ViewClassify:        VFSViewClassify,
		ViewListNoDirs:      VFSViewList,
		ViewListXNoDirs:     VFSViewList,
		ViewLevelNoDirs:     VFSViewLevel,
		ViewLevelXNoDirs:    VFSViewLevel,
		ViewTableNoDirs:     VFSViewTable,
		ViewTableXNoDirs:    VFSViewTable,
		ViewClassifyNoDirs:  VFSViewClassify,
		ViewListNoFiles:     VFSViewList,
		ViewListXNoFiles:    VFSViewList,
		ViewLevelNoFiles:    VFSViewLevel,
		ViewLevelXNoFiles:   VFSViewLevel,
		ViewTableNoFiles:    VFSViewTable,
		ViewTableXNoFiles:   VFSViewTable,
		ViewClassifyNoFiles: VFSViewClassify,
	}
)

func (v ViewType) String() string {
	vtname := ""
	if name, ok := ViewTypeNames[v]; ok {
		vtname = name
	} else {
		vtname = "Unknown"
	}
	return vtname
}

// IsOk returns true for effective and otherwise not. In genernal, use it in checking.
func (v ViewType) IsOk() bool {
	paw.Logger.Trace("checking ViewType..." + paw.Caller(1))
	if _, ok := ViewTypeNames[v]; ok {
		return true
	} else {
		return false
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
	hasX = vt&_ViewExtended != 0
	isViewNoDirs = vt&_ViewNoDirs != 0
	isViewNoFiles = vt&_ViewNoFiles != 0
	return hasX, isViewNoDirs, isViewNoFiles
}

func (v *VFS) hasList_hasX() (hasList, hasX bool) {
	var (
		vt = v.opt.ViewType
	)
	hasList = vt&_ViewList != 0
	hasX = vt&_ViewExtended != 0
	return hasList, hasX
}
