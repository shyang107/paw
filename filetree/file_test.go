package filetree

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestConstructFile(t *testing.T) {
	var (
		path = []string{
			"/Users/shyang/go/src/github.com/shyang107/paw/filetree",
			"/Users/shyang/go/src/github.com/shyang107/paw/filetree/",
			"/Users/shyang/go/src/github.com/shyang107/paw/filetree/file.go",
		}
		stat = []os.FileInfo{}
	)
	for _, p := range path {
		s, _ := os.Lstat(p)
		stat = append(stat, s)
	}
	type args struct {
		path string
	}
	type test struct {
		name string
		args args
		want *File
	}
	tests := []test{}
	for i := 0; i < len(path); i++ {
		p, _ := filepath.Abs(path[i])
		dir, basename := filepath.Split(p)
		tests = append(tests, test{
			name: p,
			args: args{p},
			want: &File{
				Path:     p,
				Dir:      dir,      //filepath.Dir(p),
				BaseName: basename, //filepath.Base(p),
				File:     strings.TrimSuffix(filepath.Base(p), filepath.Ext(p)),
				Ext:      filepath.Ext(p),
				Stat:     stat[i],
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := NewFile(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConstructFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
