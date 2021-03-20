package vfs

import "github.com/shyang107/paw"

type Group int

const (
	// GroupNone is default order of []DirEntryX returned by Dir.ReadDir(n)
	// 	see example/vfs
	GroupNone Group = iota + 1
	// Grouped use in returned values of Dir.ReadDir(n). Arrange order of []DirEntryX by dir..., file...
	// 	see example/vfs
	Grouped
	// GroupedR use in returned values of Dir.ReadDir(n). Arrange reverse order of []DirEntryX by file..., dir...
	// 	see example/vfs
	GroupedR
)

func (g Group) String() string {
	return []string{"Not grouped", "Grouped", "Grouped reversely"}[g-1]
}

// IsOk returns true for effective and otherwise not. In genernal, use it in checking.
func (g Group) IsOk() bool {
	paw.Logger.Debug("checking Group..." + paw.Caller(1))
	if g > 0 && g < 4 {
		return true
	}
	return false
	// switch g {
	// case Grouped, GroupedR, GroupNone:
	// 	return true
	// default:
	// 	return false
	// }
}
