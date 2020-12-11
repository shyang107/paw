package filetree

// filetree is tree structure of files
//
//  every thing is file, even directory just a special file!
//

// File will store information of a file
type File struct {
	// FullPath is the full path string of a file
	FullPath string
	// Dir is directory part of `FullPath`
	Dir string
	// File is file part (basebame) of `FullPath` (including extention)
	File string
	// Name is name part
	Name string
}
