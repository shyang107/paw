package vfs

import (
	"sort"
	"strings"
)

type ByLessX func(fi, fj DirEntryX) bool

// ByFunc is the type of a "less" function that defines the ordering of its File arguments.
//
// Example:
// 	lowerPathName := func(fi, fj *DirEntryX) bool {
// 		return paw.ToLower(fi.Path) < paw.ToLower(fj.Path)
// 	}
// 	ByFunc(lowerPathName).Sort(files)
type ByFunc struct {
	Name       string
	_IsReverse bool
	Less       ByLessX
}

func (b ByFunc) String() string {
	return b.Name
}

func (b ByFunc) SetReverse() ByFunc {
	b._IsReverse = true
	b.Name = "reversely " + b.Name
	return b
}

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (b *ByFunc) Sort(files []DirEntryX) {
	// paw.Logger.Trace("sorting..." + paw.Caller(1))
	ps := &DirEntryXSorter{
		files: files,
		by:    b, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	if b._IsReverse {
		sort.Sort(sort.Reverse(ps))
	} else {
		sort.Sort(ps)
	}
}

// // SortReverse like as Byfunc.Sort, but in reverse order.
// func (b ByFunc) SortReverse(files []DirEntryX) {
// 	// paw.Logger.Trace("sorting..." + paw.Caller(1))
// 	ps := &DirEntryXSorter{
// 		files: files,
// 		by:    b, // The Sort method's receiver is the function (closure) that defines the sort order.
// 	}
// 	sort.Sort(sort.Reverse(ps))
// }

// DirEntryXSorter joins a By function and a slice of Files to be sorted.
type DirEntryXSorter struct {
	files []DirEntryX
	by    *ByFunc //func(p1, p2 DirEntryX) bool
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
	return s.by.Less(s.files[i], s.files[j])
}

var (
	ByINodeFunc = ByFunc{
		Name:       "by Inode",
		_IsReverse: false,
		Less: func(fi, fj DirEntryX) bool {
			return fi.INode() < fj.INode()
		},
	}

	ByINodeFuncR = ByINodeFunc.SetReverse()

	ByHDLinksFunc = ByFunc{
		Name:       "by HDLinks",
		_IsReverse: false,
		Less: func(fi, fj DirEntryX) bool {
			return fi.HDLinks() < fj.HDLinks()
		},
	}

	ByHDLinksFuncR = ByHDLinksFunc.SetReverse()

	ByPathFunc = ByFunc{
		Name:       "by lower Path",
		_IsReverse: false,
		Less: func(fi, fj DirEntryX) bool {
			return strings.ToLower(fi.Path()) < strings.ToLower(fj.Path())
		},
	}

	BySizeFunc = ByFunc{
		Name:       "by Size",
		_IsReverse: false,
		Less: func(fi, fj DirEntryX) bool {
			// if fi.IsDir() && fj.IsDir() {
			// 	return ByPathFunc(fi, fj)
			// }
			return fi.Size() < fj.Size()
		},
	}

	BySizeFuncR = BySizeFunc.SetReverse()

	ByBlocksFunc = ByFunc{
		Name:       "by Blocks",
		_IsReverse: false,
		Less: func(fi, fj DirEntryX) bool {
			// if fi.IsDir() && fj.IsDir() {
			// 	return ByPathFunc(fi, fj)
			// }
			return fi.Blocks() < fj.Blocks()
		},
	}

	ByBlocksFuncR = ByBlocksFunc.SetReverse()

	ByMTimeFunc = ByFunc{
		Name:       "by MTime",
		_IsReverse: false,
		Less: func(fi, fj DirEntryX) bool {
			return fi.ModifiedTime().Before(fj.ModifiedTime())
		},
	}

	ByMTimeFuncR = ByMTimeFunc.SetReverse()

	ByATimeFunc = ByFunc{
		Name:       "by ATime",
		_IsReverse: false,
		Less: func(fi, fj DirEntryX) bool {
			return fi.AccessedTime().Before(fj.AccessedTime())
		},
	}

	ByATimeFuncR = ByATimeFunc.SetReverse()

	ByCTimeFunc = ByFunc{
		Name:       "by CTime",
		_IsReverse: false,
		Less: func(fi, fj DirEntryX) bool {
			return fi.CreatedTime().Before(fj.CreatedTime())
		},
	}

	ByCTimeFuncR = ByCTimeFunc.SetReverse()

	ByNameFunc = ByFunc{
		Name:       "by Name",
		_IsReverse: false,
		Less: func(fi, fj DirEntryX) bool {
			return fi.Name() < fj.Name()
		},
	}

	ByNameFuncR = ByNameFunc.SetReverse()

	ByLowerNameFunc = ByFunc{
		Name:       "by lower Name",
		_IsReverse: false,
		Less: func(fi, fj DirEntryX) bool {
			return strings.ToLower(fi.Name()) < strings.ToLower(fj.Name())
		},
	}

	ByLowerNameFuncR = ByLowerNameFunc.SetReverse()
)

type SortKey int

const (
	_SortReverse SortKey = 1 << iota
	SortByINode
	SortByHDLinks
	SortBySize
	SortByBlocks
	SortByMTime
	SortByATime
	SortByCTime
	SortByName
	SortByLowerName
	SortByINodeR     = _SortReverse | SortByINode
	SortByHDLinksR   = _SortReverse | SortByHDLinks
	SortBySizeR      = _SortReverse | SortBySize
	SortByBlocksR    = _SortReverse | SortByBlocks
	SortByMTimeR     = _SortReverse | SortByMTime
	SortByATimeR     = _SortReverse | SortByATime
	SortByCTimeR     = _SortReverse | SortByCTime
	SortByNameR      = _SortReverse | SortByName
	SortByLowerNameR = _SortReverse | SortByLowerName
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

func (a ByINode) String() string {
	return ByINodeFunc.Name
}
func (a ByINode) Len() int      { return len(a.values) }
func (a ByINode) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByINode) Less(i, j int) bool {
	return ByINodeFunc.Less(a.values[i], a.values[i])
}

type ByHDLinks struct{ values []DirEntryX }

func (a ByHDLinks) String() string {
	return ByHDLinksFunc.Name
}
func (a ByHDLinks) Len() int      { return len(a.values) }
func (a ByHDLinks) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByHDLinks) Less(i, j int) bool {
	return ByHDLinksFunc.Less(a.values[i], a.values[i])
}

type BySize struct{ values []DirEntryX }

func (a BySize) String() string {
	return BySizeFunc.Name
}
func (a BySize) Len() int      { return len(a.values) }
func (a BySize) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a BySize) Less(i, j int) bool {
	return BySizeFunc.Less(a.values[i], a.values[i])
}

type ByBlocks struct{ values []DirEntryX }

func (a ByBlocks) String() string {
	return ByBlocksFunc.Name
}
func (a ByBlocks) Len() int      { return len(a.values) }
func (a ByBlocks) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByBlocks) Less(i, j int) bool {
	return ByBlocksFunc.Less(a.values[i], a.values[i])
}

type ByMTime struct{ values []DirEntryX }

func (a ByMTime) String() string {
	return ByMTimeFunc.Name
}
func (a ByMTime) Len() int      { return len(a.values) }
func (a ByMTime) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByMTime) Less(i, j int) bool {
	return ByMTimeFunc.Less(a.values[i], a.values[i])
}

type ByCTime struct{ values []DirEntryX }

func (a ByCTime) String() string {
	return ByCTimeFunc.Name
}
func (a ByCTime) Len() int      { return len(a.values) }
func (a ByCTime) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByCTime) Less(i, j int) bool {
	return ByCTimeFunc.Less(a.values[i], a.values[i])
}

type ByATime struct{ values []DirEntryX }

func (a ByATime) String() string {
	return ByATimeFunc.Name
}
func (a ByATime) Len() int      { return len(a.values) }
func (a ByATime) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByATime) Less(i, j int) bool {
	return ByATimeFunc.Less(a.values[i], a.values[i])
}

type ByName struct{ values []DirEntryX }

func (a ByName) String() string {
	return ByNameFunc.Name
}
func (a ByName) Len() int      { return len(a.values) }
func (a ByName) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByName) Less(i, j int) bool {
	return ByNameFunc.Less(a.values[i], a.values[i])
}

type ByLowerName struct{ values []DirEntryX }

func (a ByLowerName) String() string {
	return ByLowerNameFunc.Name
}
func (a ByLowerName) Len() int      { return len(a.values) }
func (a ByLowerName) Swap(i, j int) { a.values[i], a.values[j] = a.values[j], a.values[i] }
func (a ByLowerName) Less(i, j int) bool {
	return ByLowerNameFunc.Less(a.values[i], a.values[i])
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

var _DirEntryXALessFunc = ByFunc{
	Name: "by DirEntryXA",
	Less: nil,
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s DirEntryXA) Less(i, j int) bool {
	return _DirEntryXALessFunc.Less(s[i], s[j])
}

// SetLessFunc set less func using for sorting of DirEntryXA. If less is nil, less is ByLowerNameFunc.
func (s DirEntryXA) SetLessFunc(by ByFunc) *DirEntryXA {
	_DirEntryXALessFunc = ByFunc{
		Name: by.Name,
		Less: by.Less,
	}
	return &s
}
