package vfs

import (
	"sort"
	"strings"
)

// ByFunc is the type of a "less" function that defines the ordering of its File arguments.
//
// Example:
// 	lowerPathName := func(fi, fj *DirEntryX) bool {
// 		return paw.ToLower(fi.Path) < paw.ToLower(fj.Path)
// 	}
// 	ByFunc(lowerPathName).Sort(files)
type ByFunc func(fi, fj DirEntryX) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by ByFunc) Sort(files []DirEntryX) {
	// paw.Logger.Trace("sorting..." + paw.Caller(1))
	ps := &DirEntryXSorter{
		files: files,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// DirEntryXSorter joins a By function and a slice of Files to be sorted.
type DirEntryXSorter struct {
	files []DirEntryX
	by    ByFunc //func(p1, p2 DirEntryX) bool
}

// Len is part of sort.Interface.
func (s *DirEntryXSorter) Len() int {
	return len(s.files)
}

// Swap is part of sort.Interface.
func (s *DirEntryXSorter) Swap(i, j int) {
	s.files[i], s.files[j] = s.files[j], s.files[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *DirEntryXSorter) Less(i, j int) bool {
	return s.by(s.files[i], s.files[j])
}

type ByLowerString struct {
	values []string
}

func (a ByLowerString) Len() int      { return len(a.values) }
func (a ByLowerString) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByLowerString) Less(i, j int) bool {
	return strings.ToLower(a.values[i]) < strings.ToLower(a.values[j])
}

var (
	ByINodeFunc ByFunc = func(fi, fj DirEntryX) bool {
		return fi.INode() < fj.INode()
	}
	ByINodeFuncR ByFunc = func(fi, fj DirEntryX) bool {
		return ByINodeFunc(fj, fi)
	}

	ByHDLinksFunc ByFunc = func(fi, fj DirEntryX) bool {
		return fi.HDLinks() < fj.HDLinks()
	}
	ByHDLinksFuncR ByFunc = func(fi, fj DirEntryX) bool {
		return ByHDLinksFunc(fj, fi)
	}

	ByPathFunc ByFunc = func(fi, fj DirEntryX) bool {
		return strings.ToLower(fi.Path()) < strings.ToLower(fj.Path())
	}

	BySizeFunc ByFunc = func(fi, fj DirEntryX) bool {
		// if fi.IsDir() && fj.IsDir() {
		// 	return ByPathFunc(fi, fj)
		// }
		return fi.Size() < fj.Size()
	}
	BySizeFuncR ByFunc = func(fi, fj DirEntryX) bool {
		return BySizeFunc(fj, fi)
	}

	ByBlocksFunc ByFunc = func(fi, fj DirEntryX) bool {
		// if fi.IsDir() && fj.IsDir() {
		// 	return ByPathFunc(fi, fj)
		// }
		return fi.Blocks() < fj.Blocks()
	}
	ByBlocksFuncR ByFunc = func(fi, fj DirEntryX) bool {
		return ByBlocksFunc(fj, fi)
	}

	ByMTimeFunc ByFunc = func(fi, fj DirEntryX) bool {
		return fi.ModifiedTime().Before(fj.ModifiedTime())
	}
	ByMTimeFuncR ByFunc = func(fi, fj DirEntryX) bool {
		return ByMTimeFunc(fj, fi)
	}

	ByATimeFunc ByFunc = func(fi, fj DirEntryX) bool {
		return fi.AccessedTime().Before(fj.AccessedTime())
	}
	ByATimeFuncR ByFunc = func(fi, fj DirEntryX) bool {
		return ByATimeFunc(fj, fi)
	}
	ByCTimeFunc ByFunc = func(fi, fj DirEntryX) bool {
		return fi.CreatedTime().Before(fj.CreatedTime())
	}
	ByCTimeFuncR ByFunc = func(fi, fj DirEntryX) bool {
		return ByCTimeFunc(fj, fi)
	}

	ByNameFunc ByFunc = func(fi, fj DirEntryX) bool {
		return fi.Name() < fj.Name()
	}
	ByNameFuncR ByFunc = func(fi, fj DirEntryX) bool {
		return ByLowerNameFunc(fj, fi)
	}

	ByLowerNameFunc ByFunc = func(fi, fj DirEntryX) bool {
		return strings.ToLower(fi.Name()) < strings.ToLower(fj.Name())
	}
	ByLowerNameFuncR ByFunc = func(fi, fj DirEntryX) bool {
		return ByLowerNameFunc(fj, fi)
	}
)

type SortKey int

const (
	SortReverse SortKey = 1 << iota
	SortByINode
	SortByHDLinks
	SortBySize
	SortByBlocks
	SortByMTime
	SortByATime
	SortByCTime
	SortByName
	SortByLowerName
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
	SortFuncMap = map[SortKey]*ByFunc{
		SortByINode:      &ByINodeFunc,
		SortByHDLinks:    &ByHDLinksFunc,
		SortBySize:       &BySizeFunc,
		SortByBlocks:     &ByBlocksFunc,
		SortByMTime:      &ByMTimeFunc,
		SortByATime:      &ByATimeFunc,
		SortByCTime:      &ByCTimeFunc,
		SortByName:       &ByNameFunc,
		SortByLowerName:  &ByLowerNameFunc,
		SortByINodeR:     &ByINodeFuncR,
		SortByHDLinksR:   &ByHDLinksFuncR,
		SortBySizeR:      &BySizeFuncR,
		SortByBlocksR:    &ByBlocksFuncR,
		SortByMTimeR:     &ByMTimeFuncR,
		SortByATimeR:     &ByATimeFuncR,
		SortByCTimeR:     &ByCTimeFuncR,
		SortByNameR:      &ByNameFuncR,
		SortByLowerNameR: &ByLowerNameFuncR,
	}
)

type ByINode struct{ values []DirEntryX }

func (a ByINode) Len() int      { return len(a.values) }
func (a ByINode) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByINode) Less(i, j int) bool {
	return ByINodeFunc(a.values[i], a.values[i])
}

type ByHDLinks struct{ values []DirEntryX }

func (a ByHDLinks) Len() int      { return len(a.values) }
func (a ByHDLinks) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByHDLinks) Less(i, j int) bool {
	return ByHDLinksFunc(a.values[i], a.values[i])
}

type BySize struct{ values []DirEntryX }

func (a BySize) Len() int      { return len(a.values) }
func (a BySize) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a BySize) Less(i, j int) bool {
	return BySizeFunc(a.values[i], a.values[i])
}

type ByBlocks struct{ values []DirEntryX }

func (a ByBlocks) Len() int      { return len(a.values) }
func (a ByBlocks) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByBlocks) Less(i, j int) bool {
	return ByBlocksFunc(a.values[i], a.values[i])
}

type ByMTime struct{ values []DirEntryX }

func (a ByMTime) Len() int      { return len(a.values) }
func (a ByMTime) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByMTime) Less(i, j int) bool {
	return ByMTimeFunc(a.values[i], a.values[i])
}

type ByCTime struct{ values []DirEntryX }

func (a ByCTime) Len() int      { return len(a.values) }
func (a ByCTime) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByCTime) Less(i, j int) bool {
	return ByCTimeFunc(a.values[i], a.values[i])
}

type ByATime struct{ values []DirEntryX }

func (a ByATime) Len() int      { return len(a.values) }
func (a ByATime) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByATime) Less(i, j int) bool {
	return ByATimeFunc(a.values[i], a.values[i])
}

type ByName struct{ values []DirEntryX }

func (a ByName) Len() int      { return len(a.values) }
func (a ByName) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByName) Less(i, j int) bool {
	return ByNameFunc(a.values[i], a.values[i])
}

type ByLowerName struct{ values []DirEntryX }

func (a ByLowerName) Len() int      { return len(a.values) }
func (a ByLowerName) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByLowerName) Less(i, j int) bool {
	return ByLowerNameFunc(a.values[i], a.values[i])
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

var _DirEntryXALessFunc ByFunc

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s DirEntryXA) Less(i, j int) bool {
	return _DirEntryXALessFunc(s[i], s[j])
}

// SetLessFunc set less func using for sorting of DirEntryXA. If less is nil, less is ByLowerNameFunc.
func (s DirEntryXA) SetLessFunc(less ByFunc) *DirEntryXA {
	if less == nil {
		less = ByLowerNameFunc
	}
	_DirEntryXALessFunc = less
	return &s
}
