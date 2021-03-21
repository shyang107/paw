package main

import (
	"strings"
	"testing"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
)

var (
	root     = "/Users/shyang/DEVONthink"
	skipcond = vfs.NewSkipConds().Add(vfs.DefaultSkiper)
	vfields  = vfs.DefaultViewField
	vopt     = vfs.NewVFSOption()
)

func BenchmarkVFS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vfs01(root)
	}
}

func vfs01(root string) {
	vopt.Depth = -1
	fs, err := vfs.NewVFS(root, vopt)
	if err != nil {
		paw.Error.Println(err)
		return
	}
	fs.BuildFS()
	sb := new(strings.Builder)
	fs.View(sb)
}
