package vfs

import (
	"fmt"

	"github.com/shyang107/paw"
)

// VFSOption uses in VFS
type VFSOption struct {
	Depth      int
	Grouping   Group
	ByField    SortKey
	Skips      *SkipConds
	ViewFields ViewField
	ViewType   ViewType
}

// NewVFSOption creates a new instance of VFSOption
func NewVFSOption() *VFSOption {
	return &VFSOption{
		Depth:      0,
		Grouping:   GroupNone,
		ByField:    SortByLowerName,
		Skips:      NewSkipConds().Add(DefaultSkiper),
		ViewFields: DefaultViewField,
		ViewType:   ViewList,
	}
}

func (v VFSOption) String() string {
	s := fmt.Sprintf("[Depth: %d]", v.Depth)
	s += fmt.Sprintf("[Grouping: %q]", v.Grouping)
	s += fmt.Sprintf("[Sort: %q]", v.ByField)
	s += fmt.Sprintf("[Skips: %q]", v.Skips)
	s += fmt.Sprintf("[ViewFields: %q]", v.ViewFields)
	s += fmt.Sprintf("[ViewType: %q]", v.ViewType)
	return s
}

func (v *VFSOption) Sort(dxs []DirEntryX) {
	v.ByField.Sort(dxs)
}

func (opt *VFSOption) Check() {
	paw.Logger.Debug("checking VFSOption..." + paw.Caller(1))

	if opt == nil {
		opt = NewVFSOption()
	} else {
		if !opt.Grouping.IsOk() {
			opt.Grouping = GroupNone
		}

		if !opt.ByField.IsOk() {
			opt.ByField = SortByLowerName
		}

		if opt.Skips == nil {
			opt.Skips = NewSkipConds()
		}

		// if !opt.Skips.IsOk() {
		// 	opt.Skips = NewSkipConds().Add(DefaultSkip)
		// }

		if !opt.ViewType.IsOk() {
			opt.ViewType = ViewList
		}

		if !opt.ViewFields.IsOk() {
			opt.ViewFields = DefaultViewField
		}
	}
}
