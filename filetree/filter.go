package filetree

import (
	"path/filepath"

	"github.com/shyang107/paw"

	"github.com/thoas/go-funk"
)

// Filter is the type of filter function that define the filtering of its arguments
type Filter func(fl *FileList)

type fileListFilter struct {
	fileList *FileList
	filters  []Filter
}

func NewFileListFilter(fl *FileList, filters []Filter) *fileListFilter {
	return &fileListFilter{
		fileList: fl,
		filters:  filters,
	}
}

func (ft *fileListFilter) Filt() {
	// fmt.Println("nfilters:", len(ft.filters))
	for _, filt := range ft.filters {
		// fmt.Println("filt", i)
		// spew.Dump(ft.fileList.dirs)
		filt(ft.fileList)
		// spew.Dump(ft.fileList.dirs)
	}
}

func removeFile(s []*File, i int) []*File {
	return append(s[:i], s[i+1:]...)
}

func removeString(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
}

func indexOf(files []*File, cond func(f *File) bool) int {
	idx := -1
	for i, file := range files {
		if cond(file) {
			idx = i
			break
		}
	}
	return idx
}

var (
	FiltEmptyDirs Filter = func(fl *FileList) {
		// paw.Info.Println("FiltEmptyDirs")
		var emptyDirs []string
		// spew.Dump(fl.Dirs())
		for _, dir := range fl.dirs {
			hasEmpty := false
			var name, pdir string
			if len(fl.store[dir]) <= 1 {
				emptyDirs = append(emptyDirs, dir)
				hasEmpty = true
				tdirs := paw.Split(dir, PathSeparator)
				name = tdirs[len(tdirs)-1]
				pdir = filepath.Join(tdirs[:len(tdirs)-1]...)
				if pdir == ".." {
					pdir = "."
				}
				// pdir = fl.dirs[i-1]
				// fmt.Println("empty> dir:", dir, "name:", name, "pdir:", pdir)
			}
			if hasEmpty {
				jdx := indexOf(fl.store[pdir], func(f *File) bool {
					if f.IsDir() && f.BaseName == name {
						return true
					}
					return false
				})
				if jdx != -1 {
					fl.store[pdir] = removeFile(fl.store[pdir], jdx)
				}
				hasEmpty = false
			}
		}
		for _, v := range emptyDirs {
			delete(fl.store, v)
			idx := funk.IndexOfString(fl.dirs, v)
			if idx != -1 {
				fl.dirs = removeString(fl.dirs, idx)
			}
		}
	}

	FiltJustDirs Filter = func(fl *FileList) {
		// nd, nf := 0, 0
		for _, dir := range fl.dirs {
			var dirs []*File
			// dirs = append(dirs, )
			for _, file := range fl.store[dir][:] {
				if file.IsDir() {
					// nd++
					dirs = append(dirs, file)
				}
			}
			fl.store[dir] = dirs
			// nd += len(dirs) - 1
		}
		// fmt.Println("ndirs:", len(fl.dirs), "nd:", nd, "NDirs:", fl.NDirs())
		// fmt.Println("nstore:", len(fl.store), "nfiles:", nf, "NFiles:", fl.NFiles())
	}

	FiltJustFiles Filter = func(fl *FileList) {
		// // spew.Dump(fl.dirs)
		FiltEmptyDirs(fl)
		// spew.Dump(fl.dirs)
		var dirs []string
		// nd, nf := 0, 0
		for _, dir := range fl.dirs {
			if len(fl.store[dir]) <= 1 {
				continue
			}
			var files []*File
			files = append(files, fl.store[dir][0])
			for _, file := range fl.store[dir][1:] {
				if !file.IsDir() {
					// nf++
					files = append(files, file)
					if funk.IndexOfString(dirs, dir) == -1 {
						// nd++
						dirs = append(dirs, file.Dir)
					}
				}
			}
			fl.store[dir] = files
		}
		fl.dirs = dirs
		// spew.Dump(fl.dirs)
		// fmt.Println("ndirs:", len(fl.dirs), "nd:", nd, "NDirs:", fl.NDirs())
		// fmt.Println("nstore:", len(fl.store), "nfiles:", nf, "NFiles:", fl.NFiles())
	}
)
