package vfs

import (
	"sort"
	"strings"

	"github.com/shyang107/paw"
)

// type Interface interface {
// 	Len() int
// 	Less(i, j int) bool
// 	Swap(i, j int)
// }

// // ByFunc is the type of a "less" function that defines the ordering of its File arguments.
// //
// // Example:
// // 	lowerPathName := func(fi, fj *DirEntryX) bool {
// // 		return paw.ToLower(fi.Path) < paw.ToLower(fj.Path)
// // 	}
// // 	ByFunc(lowerPathName).Sort(files)
// type ByFunc struct {
// 	key  SortKey
// 	Less ByLessFunc
// }

// func (b ByFunc) String() string {
// 	return SortFuncFields[b.key]
// }

// // Sort is a method on the function type, By, that sorts the argument slice according to the function.
// func (b *ByFunc) Sort(files []DirEntryX) {
// 	// paw.Logger.Debug("sorting..." + paw.Caller(1))
// 	ps := &DirEntryXSorter{
// 		files: files,
// 		by:    b, // The Sort method's receiver is the function (closure) that defines the sort order.
// 	}
// 	sort.Sort(ps)
// 	// sort.Sort(sort.Reverse(ps))
// }

// // DirEntryXSorter joins a By function and a slice of Files to be sorted.
// type DirEntryXSorter struct {
// 	files []DirEntryX
// 	by    *ByFunc //func(p1, p2 DirEntryX) bool
// }

// // Len is part of sort.Interface.
// func (s *DirEntryXSorter) Len() int { return len(s.files) }

// // Swap is part of sort.Interface.
// func (s *DirEntryXSorter) Swap(i, j int) { s.files[i], s.files[j] = s.files[j], s.files[i] }

// // Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
// func (s *DirEntryXSorter) Less(i, j int) bool { return s.by.Less(s.files[i], s.files[j]) }

// type reverse struct {
// 	// This embedded Interface permits Reverse to use the methods of
// 	// another Interface implementation.
// 	Interface
// }

// // Less returns the opposite of the embedded implementation's Less method.
// func (r reverse) Less(i, j int) bool {
// 	return r.Interface.Less(j, i)
// }

// // Reverse returns the reverse order for data.
// func Reverse(data Interface) Interface {
// 	return &reverse{data}
// }

type ByLessFunc func(fi, fj DirEntryX) bool

var (
	ByINodeLessFunc = ByLessFunc(func(fi, fj DirEntryX) bool {
		return fi.INode() < fj.INode()
	})

	ByHDLinksLessFunc = ByLessFunc(func(fi, fj DirEntryX) bool {
		return fi.HDLinks() < fj.HDLinks()
	})

	// ByHDLinksLessFuncR = ByHDLinksFunc.SetReverse()

	ByPathLessFunc = ByLessFunc(func(fi, fj DirEntryX) bool {
		return strings.ToLower(fi.Path()) < strings.ToLower(fj.Path())
	})

	BySizeLessFunc = ByLessFunc(func(fi, fj DirEntryX) bool {
		// if fi.IsDir() && fj.IsDir() {
		// 	return ByPathFunc(fi, fj)
		// }
		return fi.Size() < fj.Size()
	})

	// BySizeLessFuncR = BySizeFunc.SetReverse()
	// BySizeLessFuncR = ByLessFunc{
	// 	Name: "by Size",
	// 	// _IsReverse: true,
	// 	Less: func(fi, fj DirEntryX) bool {
	// 		return BySizeFunc.Less(fj, fi)
	// 	},
	// }

	ByBlocksLessFunc = ByLessFunc(func(fi, fj DirEntryX) bool {
		// if fi.IsDir() && fj.IsDir() {
		// 	return ByPathFunc(fi, fj)
		// }
		return fi.Blocks() < fj.Blocks()
	})

	// ByBlocksLessFuncR = ByBlocksFunc.SetReverse()

	ByMTimeLessFunc = ByLessFunc(func(fi, fj DirEntryX) bool {
		return fi.ModifiedTime().Before(fj.ModifiedTime())
	})

	// ByMTimeLessFuncR = ByMTimeFunc.SetReverse()

	ByATimeLessFunc = ByLessFunc(func(fi, fj DirEntryX) bool {
		return fi.AccessedTime().Before(fj.AccessedTime())
	})

	// ByATimeLessFuncR = ByATimeFunc.SetReverse()

	ByCTimeLessFunc = ByLessFunc(func(fi, fj DirEntryX) bool {
		return fi.CreatedTime().Before(fj.CreatedTime())
	})

	// ByCTimeLessFuncR = ByCTimeFunc.SetReverse()

	ByNameLessFunc = ByLessFunc(func(fi, fj DirEntryX) bool {
		return fi.Name() < fj.Name()
	})

	// ByNameLessFuncR = ByNameFunc.SetReverse()

	ByLowerNameLessFunc = ByLessFunc(func(fi, fj DirEntryX) bool {
		return strings.ToLower(fi.Name()) < strings.ToLower(fj.Name())
	})

	// ByLowerNameLessFuncR = ByLowerNameFunc.SetReverse()
)

type SortKey int

const (
	SortByINode SortKey = 1 << iota
	SortByHDLinks
	SortBySize
	SortByBlocks
	SortByMTime
	SortByATime
	SortByCTime
	SortByName
	SortByLowerName

	SortByNone
	SortReverse

	SortByINodeR     = SortReverse | SortByINode
	SortByHDLinksR   = SortReverse | SortByHDLinks
	SortBySizeR      = SortReverse | SortBySize
	SortByBlocksR    = SortReverse | SortByBlocks
	SortByMTimeR     = SortReverse | SortByMTime
	SortByATimeR     = SortReverse | SortByATime
	SortByCTimeR     = SortReverse | SortByCTime
	SortByNameR      = SortReverse | SortByName
	SortByLowerNameR = SortReverse | SortByLowerName
)

var (
	SortLessFuncMap = map[SortKey]ByLessFunc{
		SortByINode:      ByINodeLessFunc,
		SortByHDLinks:    ByHDLinksLessFunc,
		SortBySize:       BySizeLessFunc,
		SortByBlocks:     ByBlocksLessFunc,
		SortByMTime:      ByMTimeLessFunc,
		SortByATime:      ByATimeLessFunc,
		SortByCTime:      ByCTimeLessFunc,
		SortByName:       ByNameLessFunc,
		SortByLowerName:  ByLowerNameLessFunc,
		SortByINodeR:     ByINodeLessFunc,
		SortByHDLinksR:   ByHDLinksLessFunc,
		SortBySizeR:      BySizeLessFunc,
		SortByBlocksR:    ByBlocksLessFunc,
		SortByMTimeR:     ByMTimeLessFunc,
		SortByATimeR:     ByATimeLessFunc,
		SortByCTimeR:     ByCTimeLessFunc,
		SortByNameR:      ByNameLessFunc,
		SortByLowerNameR: ByLowerNameLessFunc,
	}

	SortFuncFields = map[SortKey]string{
		SortByNone:       "none",
		SortByINode:      "INode",
		SortByHDLinks:    "HDLinks",
		SortBySize:       "Size",
		SortByBlocks:     "Blocks",
		SortByMTime:      "MTime",
		SortByATime:      "ATime",
		SortByCTime:      "CTime",
		SortByName:       "Name",
		SortByLowerName:  "LowerName",
		SortByINodeR:     "reverse INode",
		SortByHDLinksR:   "reverse HDLinks",
		SortBySizeR:      "reverse Size",
		SortByBlocksR:    "reverse Blocks",
		SortByMTimeR:     "reverse MTime",
		SortByATimeR:     "reverse ATime",
		SortByCTimeR:     "reverse CTime",
		SortByNameR:      "reverse Name",
		SortByLowerNameR: "reverse LowerName",
	}
	SortKeyNames = map[SortKey]string{
		SortByNone:       "SortByNone",
		SortByINode:      "SortByINode",
		SortByHDLinks:    "SortByHDLinks",
		SortBySize:       "SortBySize",
		SortByBlocks:     "SortByBlocks",
		SortByMTime:      "SortByMTime",
		SortByATime:      "SortByATime",
		SortByCTime:      "SortByCTime",
		SortByName:       "SortByName",
		SortByLowerName:  "SortByLowerName",
		SortByINodeR:     "SortByINodeR",
		SortByHDLinksR:   "SortByHDLinksR",
		SortBySizeR:      "SortBySizeR",
		SortByBlocksR:    "SortByBlocksR",
		SortByMTimeR:     "SortByMTimeR",
		SortByATimeR:     "SortByATimeR",
		SortByCTimeR:     "SortByCTimeR",
		SortByNameR:      "SortByNameR",
		SortByLowerNameR: "SortByLowerNameR",
	}
	SortNameKeys = map[string]SortKey{
		"SortByNone":       SortByNone,
		"SortByINode":      SortByINode,
		"SortByHDLinks":    SortByHDLinks,
		"SortBySize":       SortBySize,
		"SortByBlocks":     SortByBlocks,
		"SortByMTime":      SortByMTime,
		"SortByATime":      SortByATime,
		"SortByCTime":      SortByCTime,
		"SortByName":       SortByName,
		"SortByLowerName":  SortByLowerName,
		"SortByINodeR":     SortByINodeR,
		"SortByHDLinksR":   SortByHDLinksR,
		"SortBySizeR":      SortBySizeR,
		"SortByBlocksR":    SortByBlocksR,
		"SortByMTimeR":     SortByMTimeR,
		"SortByATimeR":     SortByATimeR,
		"SortByCTimeR":     SortByCTimeR,
		"SortByNameR":      SortByNameR,
		"SortByLowerNameR": SortByLowerNameR,
	}

	SortShortNameKeys = map[string]SortKey{
		"none":    SortByNone,
		"inode":   SortByINode,
		"links":   SortByHDLinks,
		"size":    SortBySize,
		"blocks":  SortByBlocks,
		"mtime":   SortByMTime,
		"atime":   SortByATime,
		"ctime":   SortByCTime,
		"name":    SortByName,
		"lname":   SortByLowerName,
		"inoder":  SortByINodeR,
		"linksr":  SortByHDLinksR,
		"sizer":   SortBySizeR,
		"blocksr": SortByBlocksR,
		"mtimer":  SortByMTimeR,
		"atimer":  SortByATimeR,
		"ctimer":  SortByCTimeR,
		"namer":   SortByNameR,
		"lnamer":  SortByLowerNameR,
	}
)

func (s SortKey) String() string {
	return s.Name()
}

func (s SortKey) Name() string {
	if s&SortByNone != 0 {
		return "not sort"
	}
	field, ok := SortFuncFields[s]
	if !ok {
		return "[Error] use default sort field: by " + SortFuncFields[SortByLowerName]
	}
	return "by " + field
}

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (s SortKey) Sort(dxs []DirEntryX) {
	dxa := DirEntryXA(dxs).SetLessFunc(s)
	switch s {
	case SortByINodeR, SortByHDLinksR, SortBySizeR, SortByBlocksR, SortByMTimeR, SortByATimeR, SortByCTimeR, SortByNameR, SortByLowerNameR:
		sort.Sort(sort.Reverse(dxa))
	case SortByNone:
		return
	default:
		sort.Sort(dxa)
	}
}

// IsOk returns true for effective and otherwise not. In genernal, use it in checking.
func (s SortKey) IsOk() bool {
	paw.Logger.Debug("checking SortKey..." + paw.Caller(1))

	switch s {
	case SortByINode, SortByHDLinks, SortBySize, SortByBlocks, SortByMTime, SortByATime, SortByCTime, SortByName, SortByLowerName, SortByINodeR, SortByHDLinksR, SortBySizeR, SortByBlocksR, SortByMTimeR, SortByATimeR, SortByCTimeR, SortByNameR, SortByLowerNameR, SortByNone:
		return true
	default:
		return false
	}
}

type ByINode struct{ values []DirEntryX }

func (a ByINode) String() string {
	return "sort by " + SortFuncFields[SortByINode]
}

func (a ByINode) Len() int      { return len(a.values) }
func (a ByINode) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByINode) Less(i, j int) bool {
	return ByINodeLessFunc(a.values[i], a.values[i])
}

type ByHDLinks struct{ values []DirEntryX }

func (a ByHDLinks) String() string {
	return "sort by " + SortFuncFields[SortByHDLinks]
}
func (a ByHDLinks) Len() int      { return len(a.values) }
func (a ByHDLinks) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByHDLinks) Less(i, j int) bool {
	return ByHDLinksLessFunc(a.values[i], a.values[i])
}

type BySize struct{ values []DirEntryX }

func (a BySize) String() string {
	return "sort by " + SortFuncFields[SortBySize]
}
func (a BySize) Len() int      { return len(a.values) }
func (a BySize) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a BySize) Less(i, j int) bool {
	return BySizeLessFunc(a.values[i], a.values[i])
}

type ByBlocks struct{ values []DirEntryX }

func (a ByBlocks) String() string {
	return "sort by " + SortFuncFields[SortByBlocks]
}
func (a ByBlocks) Len() int      { return len(a.values) }
func (a ByBlocks) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByBlocks) Less(i, j int) bool {
	return ByBlocksLessFunc(a.values[i], a.values[i])
}

type ByMTime struct{ values []DirEntryX }

func (a ByMTime) String() string {
	return "sort by " + SortFuncFields[SortByMTime]
}
func (a ByMTime) Len() int      { return len(a.values) }
func (a ByMTime) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByMTime) Less(i, j int) bool {
	return ByMTimeLessFunc(a.values[i], a.values[i])
}

type ByCTime struct{ values []DirEntryX }

func (a ByCTime) String() string {
	return "sort by " + SortFuncFields[SortByCTime]
}
func (a ByCTime) Len() int      { return len(a.values) }
func (a ByCTime) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByCTime) Less(i, j int) bool {
	return ByCTimeLessFunc(a.values[i], a.values[i])
}

type ByATime struct{ values []DirEntryX }

func (a ByATime) String() string {
	return "sort by " + SortFuncFields[SortByATime]
}
func (a ByATime) Len() int      { return len(a.values) }
func (a ByATime) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByATime) Less(i, j int) bool {
	return ByATimeLessFunc(a.values[i], a.values[i])
}

type ByName struct{ values []DirEntryX }

func (a ByName) String() string {
	return "sort by " + SortFuncFields[SortByName]
}
func (a ByName) Len() int      { return len(a.values) }
func (a ByName) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByName) Less(i, j int) bool {
	return ByNameLessFunc(a.values[i], a.values[i])
}

type ByLowerName struct{ values []DirEntryX }

func (a ByLowerName) String() string {
	return "sort by " + SortFuncFields[SortByLowerName]
}
func (a ByLowerName) Len() int      { return len(a.values) }
func (a ByLowerName) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByLowerName) Less(i, j int) bool {
	return ByLowerNameLessFunc(a.values[i], a.values[i])
}

type ByLowerString struct {
	values []string
}

func (a ByLowerString) Len() int      { return len(a.values) }
func (a ByLowerString) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByLowerString) Less(i, j int) bool {
	return strings.ToLower(a.values[i]) < strings.ToLower(a.values[j])
}

// DirEntryXSorter joins a By function and a slice of Files to be sorted.
type DirEntryXA []DirEntryX

// Len is part of sort.Interface.
func (s DirEntryXA) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s DirEntryXA) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var _DirEntryXALessFunc ByLessFunc = SortLessFuncMap[SortByLowerName]

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s DirEntryXA) Less(i, j int) bool {
	if _DirEntryXALessFunc == nil {
		paw.Logger.Error("less is «nil», use default SortByLowerName")
		_DirEntryXALessFunc = SortLessFuncMap[SortByLowerName]
	}
	return _DirEntryXALessFunc(s[i], s[j])
}

// SetLessFunc set less func using for sorting of DirEntryXA. If less is nil, less is ByLowerNameFunc.
func (s DirEntryXA) SetLessFunc(byField SortKey) *DirEntryXA {
	_DirEntryXALessFunc = SortLessFuncMap[byField]
	return &s
}
