package filetree

import (
	"os"
	"path/filepath"
	"strings"

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

func removeEmpty(fl *FileList, noEmptyDir bool) {
	if noEmptyDir {
		return
	}
	for _, dir := range fl.dirs {
		var files []*File
		fm := fl.store[dir]
		for _, file := range fm[1:] {
			if file.IsDir() {
				fis, err := os.ReadDir(file.Path)
				if err != nil && len(fis) > 0 {
					files = append(files, file)
					noEmptyDir = false
				} else {
					noEmptyDir = true
				}
			} else {
				files = append(files, file)
			}
		}
		if len(files) > 0 {
			files = append(files, fm[0])
			fl.store[dir] = files
		}
	}
	removeEmpty(fl, noEmptyDir)
}

var (
	// FiltEmptyDirs Filter = func(fl *FileList) {
	// 	removeEmpty(fl, false)
	// 	var dirs []string
	// 	for dir, files := range fl.store {
	// 		if len(files) > 1 {
	// 			dirs = append(dirs, dir)
	// 		}
	// 	}
	// 	sort.Sort(ByLowerString(dirs))
	// 	fl.dirs = dirs
	// }

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
				tdirs := strings.Split(dir, PathSeparator)
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
		for _, dir := range fl.dirs {
			var files []*File
			// var dirs []*File
			fm := fl.store[dir]
			files = append(files, fm[0])
			for _, f := range fm[1:] {
				if f.IsDir() {
					files = append(files, f)
				}
			}
			fl.store[dir] = files
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
