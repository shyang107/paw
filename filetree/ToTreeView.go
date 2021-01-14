package filetree

// ToTreeViewBytes will return the []byte of FileList in tree form
func (f *FileList) ToTreeViewBytes(pad string) []byte {
	return []byte(f.ToTreeView(pad))
}

// ToTreeExtendViewBytes will return the string of FileList in tree form
func (f *FileList) ToTreeExtendViewBytes(pad string) []byte {
	return []byte(f.ToTreeExtendView(pad))
}

// ToTreeView will return the string of FileList in tree form
func (f *FileList) ToTreeView(pad string) string {
	pdview = PTreeView
	return toListTreeView(f, pad, false)
}

// ToTreeExtendView will return the string of FileList icluding extend attribute in tree form
func (f *FileList) ToTreeExtendView(pad string) string {
	pdview = PTreeView
	return toListTreeView(f, pad, true)
}
