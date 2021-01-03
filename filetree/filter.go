package filetree

import (
	"fmt"
	"path/filepath"

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

var (
	FiltEmptyDirs Filter = func(fl *FileList) {
		// paw.Info.Println("FiltEmptyDirs")
		var emptyDirs []string
		for i, dir := range fl.dirs {
			hasEmpty := false
			var name, pdir string
			if len(fl.store[dir]) <= 1 {
				emptyDirs = append(emptyDirs, dir)
				hasEmpty = true
				_, name = filepath.Split(dir)
				pdir = fl.dirs[i-1]
				// fmt.Println("empty:", dir, name, pdir, i)
			}
			if hasEmpty {
				jdx := -1
				for j, file := range fl.store[pdir] {
					// fmt.Println("    ", j, file.BaseName, name)
					if file.IsDir() && file.BaseName == name {
						jdx = j
						// fmt.Println("    del", pdir, file.BaseName, jdx)
						break
					}
				}
				fl.store[pdir] = removeFile(fl.store[pdir], jdx)
				hasEmpty = false
			}
		}
		for _, v := range emptyDirs {
			delete(fl.store, v)
			i := funk.IndexOfString(fl.dirs, v)
			if i != -1 {
				fl.dirs = removeString(fl.dirs, i)
			}
		}
	}

	FiltJustDirs Filter = func(fl *FileList) {
		nd, nf := 0, 0
		for _, dir := range fl.dirs {
			var dirs []*File
			for _, file := range fl.store[dir] {
				if file.IsDir() {
					// nd++
					dirs = append(dirs, file)
				}
			}
			fl.store[dir] = dirs
			nd += len(dirs) - 1
		}
		fmt.Println("ndirs:", len(fl.dirs), "nd:", nd, "NDirs:", fl.NDirs())
		fmt.Println("nstore:", len(fl.store), "nfiles:", nf, "NFiles:", fl.NFiles())
	}

	FiltJustFiles Filter = func(fl *FileList) {
		// spew.Dump(fl.dirs)
		FiltEmptyDirs(fl)
		// spew.Dump(fl.dirs)
		var dirs []string
		nd, nf := 0, 0
		for _, dir := range fl.dirs {
			if len(fl.store[dir]) <= 1 {
				continue
			}
			var files []*File
			for _, file := range fl.store[dir][1:] {
				if !file.IsDir() {
					nf++
					files = append(files, file)
					if funk.IndexOfString(dirs, dir) == -1 {
						nd++
						dirs = append(dirs, file.Dir)
					}
				}
			}
			fl.store[dir] = files
		}
		fl.dirs = dirs
		// spew.Dump(fl.dirs)
		fmt.Println("ndirs:", len(fl.dirs), "nd:", nd, "NDirs:", fl.NDirs())
		fmt.Println("nstore:", len(fl.store), "nfiles:", nf, "NFiles:", fl.NFiles())
	}
)
