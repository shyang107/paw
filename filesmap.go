package paw

import (
	"reflect"
	"sort"
	"strings"
)

// FilesMap store files ordered by folder, use folder name as the key
type FilesMap struct {
	// store files in specific folder into store["folder"]
	store map[string][]File
	// store uniq folder into keys
	keys []string
}

// NewFilesMap is constructor of `FilesMap`
func NewFilesMap() *FilesMap {
	return &FilesMap{
		store: map[string][]File{},
		keys:  []string{},
	}
}

// NewFilesMapFrom is constructor of `FilesMap` from `files`
func NewFilesMapFrom(files []File) *FilesMap {
	o := &FilesMap{
		store: map[string][]File{},
		keys:  []string{},
	}
	for _, f := range files {
		o.SetOne(f.ShortFolder, f)
	}
	return o
}

// Getter ans Setter

// Get will return the value (`[]File`) associated with the key (`folder`).
// If the key does not exist, the second return value will be `false`.
// 	Here, `key` is the name of specific folder.
func (m *FilesMap) Get(key string) ([]File, bool) {
	val, exist := m.store[key]
	return val, exist
}

// GetAll will return the folder-[]file pairs
func (m *FilesMap) GetAll() map[string][]File {
	return m.store
}

// GetFiles will return all files in store
func (m *FilesMap) GetFiles() []File {
	files := []File{}
	for _, v := range m.store {
		files = append(files, v...)
	}
	return files
}

// Set will store a key-value (folder-files) pair.
// If the key already exists, it will overwrite the existing key-value pair.
func (m *FilesMap) Set(key string, val []File) {
	if _, exist := m.store[key]; !exist {
		m.keys = append(m.keys, key)
	}
	m.store[key] = val
}

// SetOne will store a key-value (folder-one_file) pair.
// If the key already exists, it will overwrite the existing key-value pair.
func (m *FilesMap) SetOne(key string, val File) {
	if _, exist := m.store[key]; !exist {
		m.keys = append(m.keys, key)
	}
	for _, file := range m.store[key] {
		if reflect.DeepEqual(file, val) {
			return
		}
	}
	m.store[key] = append(m.store[key], val)
}

// Delete

// Delete will remove the key and its associated vale.
func (m *FilesMap) Delete(key string) {
	delete(m.store, key)

	// find key in slice
	idx := -1
	for i, val := range m.keys {
		if val == key {
			idx = i
			break
		}
	}
	if idx != -1 {
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
	}
}

// Iterator is used to loop through the stored key-value pairs.
// The returned anonymous function returns the index, key and value
//
// Example:
// 	m := NewFileMap()()
//
// 	m.Set("folder1", []File{file11, file12})
// 	m.Set("folder2", []File{file21, file22})
// 	m.Set("folder3", []File{file31, file32})
// 	m.Set("folder4", []File{file41, file42})
// 	m.Delete("folder3")
// 	m.Delete("folder8")
//
// 	iterator := m.Iterator()
//
// 	for {
// 		i, k, v := iterator()
// 		if i == nil {
// 			break
// 		}
// 		fmt.Println(*i, *k, v)
// 	}
func (m *FilesMap) Iterator() func() (*int, *string, []File) {
	var keys = m.keys

	j := 0

	return func() (_ *int, _ *string, _ []File) {
		if j > len(keys)-1 {
			return
		}

		row := keys[j]
		j++

		return &[]int{j - 1}[0], &row, m.store[row]

	}
}

// OrderedByFolder sort keys in increasing order (code number)
func (m *FilesMap) OrderedByFolder() {
	sort.Strings(m.keys)
}

// OrderedByFolderReverse sort keys in decreasing order (code number)
func (m *FilesMap) OrderedByFolderReverse() {
	sort.Sort(sort.Reverse(sort.StringSlice(m.keys)))
}

// OrderedAll sort keys and vals ([]File) in increasing order (code number)
func (m *FilesMap) OrderedAll() {
	m.OrderedByFolder()

	// byFolder := func(f1, f2 *File) bool {
	// 	return f1.Folder < f2.Folder
	// }
	byFileName := func(f1, f2 *File) bool {
		return f1.FileName < f2.FileName
	}
	// OrderedBy(byFolder, byFileName).Sort(m.Files)
	for _, files := range m.store {
		OrderedBy(byFileName).Sort(files)
	}
}

// OrderedAllReverse sort keys and vals ([]File) in decreasing order (code number)
func (m *FilesMap) OrderedAllReverse() {
	m.OrderedByFolder()

	// byFolder := func(f1, f2 *File) bool {
	// 	return f1.Folder > f2.Folder
	// }
	byFileName := func(f1, f2 *File) bool {
		return f1.FileName > f2.FileName
	}
	// OrderedBy(byFolder, byFileName).Sort(m.Files)
	for _, files := range m.store {
		OrderedBy(byFileName).Sort(files)
	}
}

// GetFilesFunc get files with codintion `exclude` func
func (m *FilesMap) GetFilesFunc(root string, isRecursive bool, exclude func(file File) bool) {
	files, err := GetFilesFunc(root, isRecursive, exclude)
	if err != nil {
		Logger.Error(err)
	}
	fdm := collectFilesMap(files)
	for k, v := range fdm {
		m.Set(k, v)
	}
}

func collectFilesMap(files []File) (fdm map[string][]File) {
	fdm = make(map[string][]File)
	sfd := ""
	for _, f := range files {
		if !strings.EqualFold(sfd, f.ShortFolder) {
			sfd = f.ShortFolder
			fdm[f.ShortFolder] = []File{}
		}
		fdm[f.ShortFolder] = append(fdm[f.ShortFolder], f)
	}
	return fdm
}
