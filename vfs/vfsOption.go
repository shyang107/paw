package vfs

import (
	"fmt"
	"strings"

	"github.com/shyang107/paw"
)

// VFSOption uses in VFS
type VFSOption struct {
	// ScanDepth *ScanDepth
	Depth        int
	IsScanAllSub bool
	Grouping     Group
	ByField      SortKey
	Skips        *SkipConds
	ViewFields   ViewField
	ViewType     ViewType
}

// NewVFSOption creates a new instance of VFSOption
func NewVFSOption() *VFSOption {
	return &VFSOption{
		Depth:        0,
		IsScanAllSub: false,
		Grouping:     GroupNone,
		ByField:      SortByLowerName,
		Skips:        NewSkipConds().Add(DefaultSkiper),
		ViewFields:   DefaultViewField,
		ViewType:     ViewList,
	}
}

func (v VFSOption) String() string {
	depth := fmt.Sprint(v.Depth)
	if v.IsScanAllSub && v.Depth >= 0 {
		depth += "(but recurse to all directory)"
	}
	s := fmt.Sprintf("[Depth: %v]", depth)
	s += fmt.Sprintf("[Grouping: %q]", v.Grouping)
	s += fmt.Sprintf("[Sort: %q]", v.ByField)
	s += fmt.Sprintf("[Skips: %q]", v.Skips)
	s += fmt.Sprintf("[ViewFields: %q]", v.ViewFields)
	s += fmt.Sprintf("[ViewType: %q]", v.ViewType)
	return s
}

func (s *VFSOption) IsRelPathNotScan(relpath string) bool {
	if relpath == "." ||
		s.Depth <= 0 ||
		s.IsScanAllSub {
		return false
	}

	curlevel := len(strings.Split(relpath, "/"))
	return curlevel > s.Depth
}

func (s *VFSOption) IsRelPathNotView(relpath string) bool {
	if relpath == "." || s.Depth <= 0 {
		return false
	}
	curlevel := len(strings.Split(relpath, "/"))
	return curlevel > s.Depth
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

type ScanDependType int

const (
	ScanDependDepth = 1 << iota
	ScanRecurse
)

type ScanDepth struct {
	ScanDepend ScanDependType
	Depth      int
}

func NewScanDepth() *ScanDepth {
	return &ScanDepth{
		ScanDepend: ScanDependDepth,
		Depth:      0,
	}
}

func (s ScanDepth) String() string {
	switch s.ScanDepend {
	case ScanRecurse:
		return fmt.Sprintf("%d (but recurse to all directory)", s.Depth)
	// case ScanDependByDepth:
	default:
		return fmt.Sprintf("%d", s.Depth)
	}
}

func (s *ScanDepth) _IsNotScan(curlevel int) bool {
	isNoScan := true
	if s.Depth <= 0 {
		isNoScan = false
	} else { // s.Depth >0
		switch s.ScanDepend {
		case ScanRecurse:
			isNoScan = false
		default:
			return curlevel > s.Depth
		}
	}
	return isNoScan
}

func (s *ScanDepth) IsRelPathNotScan(relpath string) bool {
	if relpath == "." {
		return false
	}
	curlevel := len(strings.Split(relpath, "/"))
	return s._IsNotScan(curlevel)
}

func (s *ScanDepth) _IsNotView(curlevel int) bool {
	if s.Depth <= 0 {
		return false
	}
	return curlevel > s.Depth
}

func (s *ScanDepth) IsRelPathNotView(relpath string) bool {
	if relpath == "." {
		return false
	}
	curlevel := len(strings.Split(relpath, "/"))

	return s._IsNotView(curlevel)
}
