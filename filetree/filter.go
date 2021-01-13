package filetree

import (
	"path/filepath"

	"github.com/shyang107/paw"
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
				if pdir == UpDirMark {
					pdir = RootMark
				}
			}
			if hasEmpty {
				fm := fl.store[pdir]
				jdx := paw.LastIndexOf(len(fm), func(i int) bool {
					return fm[i].IsDir() && fm[i].BaseName == name
				})
				if jdx != -1 {
					fl.store[pdir] = append(fm[:jdx], fm[jdx+1:]...)
				}
				hasEmpty = false
			}
		}
		for _, v := range emptyDirs {
			delete(fl.store, v)
			idx := paw.LastIndexOfString(fl.dirs, v)
			if idx != -1 {
				fl.dirs = append(fl.dirs[:idx], fl.dirs[idx+1:]...)
			}
		}
	}

	FiltJustDirs Filter = func(fl *FileList) {
		// nd, nf := 0, 0
		for _, dir := range fl.dirs {
			// var dirs []*File
			fm := fl.store[dir]
			for _, f := range fm {
				j := paw.LastIndexOf(len(fm), func(i int) bool {
					return !f.IsDir()
				})
				if j != -1 {
					fl.store[dir] = append(fm[:j], fm[j+1:]...)
				}
			}
		}
	}

	FiltJustFiles Filter = func(fl *FileList) {
		FiltEmptyDirs(fl)
		var dirs []string
		for _, dir := range fl.dirs {
			fm := fl.store[dir]
			if len(fm) <= 1 {
				continue
			}
			var files []*File
			files = append(files, fl.store[dir][0])
			for _, file := range fl.store[dir][1:] {
				if !file.IsDir() {
					files = append(files, file)
					if paw.IndexOfString(dirs, dir) == -1 {
						dirs = append(dirs, file.Dir)
					}
				}
			}
			fl.store[dir] = files
		}
		fl.dirs = dirs
	}
)
