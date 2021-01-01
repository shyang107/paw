package filetree

import (
	"sort"

	"github.com/shyang107/paw"
)

// FilesBy is the type of a "less" function that defines the ordering of its File arguments.
//
// Example:
// 	lowerPathName := func(fi, fj *File) bool {
// 		return paw.ToLower(fi.Path) < paw.ToLower(fj.Path)
// 	}
// 	FilesBy(lowerPathName).Sort(files)
type FilesBy func(fi, fj *File) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by FilesBy) Sort(files []*File) {
	ps := &fileSorter{
		files: files,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// fileSorter joins a By function and a slice of Files to be sorted.
type fileSorter struct {
	files []*File
	by    func(p1, p2 *File) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *fileSorter) Len() int {
	return len(s.files)
}

// Swap is part of sort.Interface.
func (s *fileSorter) Swap(i, j int) {
	s.files[i], s.files[j] = s.files[j], s.files[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *fileSorter) Less(i, j int) bool {
	return s.by(s.files[i], s.files[j])
}

// DirsBy is the type of a "less" function that defines the ordering of its Dir arguments of FileList.
//
// Example:
// 	lowerDirhName := func(di, dj *string) bool {
// 		return paw.ToLower(di) < paw.ToLower(dj)
// 	}
// 	DirsBy(lowerDirName).Sort(dirs)
type DirsBy func(di, dj string) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by DirsBy) Sort(dirs []string) {
	ps := &dirSorter{
		dirs: dirs,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// dirSorter joins a By function and a slice of Files to be sorted.
type dirSorter struct {
	dirs []string
	by   func(p1, p2 string) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *dirSorter) Len() int {
	return len(s.dirs)
}

// Swap is part of sort.Interface.
func (s *dirSorter) Swap(i, j int) {
	s.dirs[i], s.dirs[j] = s.dirs[j], s.dirs[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *dirSorter) Less(i, j int) bool {
	return s.by(s.dirs[i], s.dirs[j])
}

// type FileSortByPathP []*File

// func (a FileSortByPathP) Len() int           { return len(a) }
// func (a FileSortByPathP) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a FileSortByPathP) Less(i, j int) bool { return a[i].BaseName < a[j].BaseName }

// ByLowerString is using in sort.Sort(data)
// 	paw.ToLower(a[i]) < paw.ToLower(a[j])
type ByLowerString []string

func (a ByLowerString) Len() int           { return len(a) }
func (a ByLowerString) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLowerString) Less(i, j int) bool { return paw.ToLower(a[i]) < paw.ToLower(a[j]) }

// type FileSortByPath []File

// func (a FileSortByPath) Len() int           { return len(a) }
// func (a FileSortByPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a FileSortByPath) Less(i, j int) bool { return a[i].Path < a[j].Path }

// ByLowerFilePath is using in sort.Sort(data).
// 	paw.ToLower(a[i].Path) < paw.ToLower(a[j].Path)

// type ByLowerFilePath []*File

// func (a ByLowerFilePath) Len() int           { return len(a) }
// func (a ByLowerFilePath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a ByLowerFilePath) Less(i, j int) bool { return paw.ToLower(a[i].Path) < paw.ToLower(a[j].Path) }

// type SortBy []*File

// func (a SortBy) Len() int           { return len(a) }
// func (a SortBy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a SortBy) Less(i, j int) bool { return paw.ToLower(a[i].Path) < paw.ToLower(a[j].Path) }
