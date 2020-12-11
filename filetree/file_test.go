package filetree

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewFile(t *testing.T) {
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
		tests = append(tests, test{
			name: p,
			args: args{p},
			want: &File{
				Path:     p,
				Dir:      filepath.Dir(p),
				BaseName: filepath.Base(p),
				File:     strings.TrimSuffix(filepath.Base(p), filepath.Ext(p)),
				Ext:      filepath.Ext(p),
				Stat:     stat[i],
			},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConstructFile(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
