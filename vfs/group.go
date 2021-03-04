package vfs

import "github.com/shyang107/paw"

type Group int

const (
	// GroupNone is default order of []DirEntryX returned by Dir.ReadDir(n)
	// 	see example/vfs
	GroupNone Group = 1 << iota
	// Grouped use in returned values of Dir.ReadDir(n). Arrange order of []DirEntryX by dir..., file...
	// 	see example/vfs
	Grouped
	// GroupedR use in returned values of Dir.ReadDir(n). Arrange reverse order of []DirEntryX by file..., dir...
	// 	see example/vfs
	GroupedR
)

func (g Group) String() string {
	switch g {
	case Grouped:
		return "Grouped"
	case GroupedR:
		return "Grouped reversely"
	default:
		return "Not grouped"
	}

}

// IsOk returns true for effective and otherwise not. In genernal, use it in checking.
func (g Group) IsOk() bool {
	paw.Logger.Trace("checking Group..." + paw.Caller(1))

	switch g {
	case Grouped, GroupedR, GroupNone:
		return true
	default:
		return false
	}
}
