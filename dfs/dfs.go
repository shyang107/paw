package dfs

import (
	"io/fs"
	"os"
	"runtime"

	"github.com/shyang107/paw"
)

type dirFS string

func (dir dirFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) || runtime.GOOS == "windows" && paw.ContainsAny(name, `\:`) {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrInvalid}
	}
	f, err := os.Open(string(dir) + "/" + name)
	if err != nil {
		return nil, err // nil fs.File
	}
	return f, nil
}

func DirFS(dir string) fs.FS {
	return dirFS(dir)
}
