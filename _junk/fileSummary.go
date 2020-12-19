package _junk

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shyang107/paw"
)

// FileSummary store the summary of file
type FileSummary struct {
	RootDir string
	AbsPath string
	RelDir  string
	Name    string
	Ext     string
	// Info    os.FileInfo
	// Mode          os.FileMode
	// Size          int64
	// ModTime       time.Time
	IsDir         bool
	IsRegularFile bool
}

// NewFileSummary is constructor of `FileSummary`
// 	if !strings.Contains(fullpath, rootdir) { return nil }
// 	if {fullpath} is {rootdir}/subdir then `Ext` = ""
func NewFileSummary(fullpath, rootdir string) *FileSummary {
	if !strings.Contains(fullpath, rootdir) {
		return nil
	}
	fullpath, _ = filepath.Abs(fullpath)
	rootdir, _ = filepath.Abs(rootdir)
	var reldir string
	if fullpath == rootdir {
		reldir = "-"
	} else {
		reldir = filepath.Dir(fullpath)
		reldir = strings.Replace(reldir, rootdir, "", 1)
	}
	name := filepath.Base(fullpath)
	ext := filepath.Ext(fullpath)
	fi, err := os.Lstat(fullpath)
	if err != nil {
		paw.Logger.Errorln(fullpath, err)
	}

	fs := &FileSummary{
		AbsPath:       fullpath,
		RootDir:       rootdir,
		RelDir:        reldir,
		Name:          name,
		Ext:           ext,
		IsDir:         fi.IsDir(),
		IsRegularFile: fi.Mode().IsRegular(),
	}
	return fs
}

// Info will return the `os.FileInfo` of the file from `os.Lstat()`
func (f *FileSummary) Info() os.FileInfo {
	fi, err := os.Lstat(f.AbsPath)
	if err != nil {
		paw.Logger.Errorln(f.AbsPath, err)
	}
	return fi
}

// Mode will return the `os.FileMode` of the file
func (f *FileSummary) Mode() os.FileMode {
	return f.Info().Mode()
}

// Size will return the size of the file
func (f *FileSummary) Size() int64 {
	return f.Info().Size()
}

// ModTime will return the modification time of the file
func (f *FileSummary) ModTime() time.Time {
	return f.Info().ModTime()
}

// /Users/shyang/go/src/github.com/shyang107/paw
// ([]string) (len=20 cap=32) {
//  (string) "",
//  (string) (len=6) "/afero",
//  (string) (len=10) "/afero/mem",
//  (string) (len=13) "/afero/sftpfs",
//  (string) (len=12) "/afero/tarfs",
//  (string) (len=21) "/afero/tarfs/testdata",
//  (string) (len=12) "/afero/zipfs",
//  (string) (len=21) "/afero/zipfs/testdata",
//  (string) (len=5) "/cast",
//  (string) (len=3) "/ex",
//  (string) (len=5) "/funk",
//  (string) (len=10) "/godirwalk",
//  (string) (len=19) "/godirwalk/examples",
//  (string) (len=29) "/godirwalk/examples/find-fast",
//  (string) (len=44) "/godirwalk/examples/remove-empty-directories",
//  (string) (len=27) "/godirwalk/examples/scanner",
//  (string) (len=25) "/godirwalk/examples/sizes",
//  (string) (len=29) "/godirwalk/examples/walk-fast",
//  (string) (len=31) "/godirwalk/examples/walk-stdlib",
//  (string) (len=10) "/treeprint"
// }

// level: 0  ""
// level: 1  "/afero"
// level: 2  "/afero/mem"
// level: 2  "/afero/sftpfs"
// level: 2  "/afero/tarfs"
// level: 3  "/afero/tarfs/testdata"
// level: 2  "/afero/zipfs"
// level: 3  "/afero/zipfs/testdata"
// level: 1  "/cast"
// level: 1  "/ex"
// level: 1  "/funk"
// level: 1  "/godirwalk"
// level: 2  "/godirwalk/examples"
// level: 3  "/godirwalk/examples/find-fast"
// level: 3  "/godirwalk/examples/remove-empty-directories"
// level: 3  "/godirwalk/examples/scanner"
// level: 3  "/godirwalk/examples/sizes"
// level: 3  "/godirwalk/examples/walk-fast"
// level: 3  "/godirwalk/examples/walk-stdlib"
// level: 1  "/treeprint"

// "" 6 [afero cast ex funk godirwalk treeprint]
//    "/afero" 4 [mem sftpfs tarfs zipfs]
//       "/afero/mem" 0 []
//       "/afero/sftpfs" 0 []
//       "/afero/tarfs" 1 [testdata]
//          "/afero/tarfs/testdata" 0 []
//       "/afero/zipfs" 1 [testdata]
//          "/afero/zipfs/testdata" 0 []
//    "/cast" 0 []
//    "/ex" 0 []
//    "/funk" 0 []
//    "/godirwalk" 1 [examples]
//       "/godirwalk/examples" 6 [find-fast remove-empty-directories scanner sizes walk-fast walk-stdlib]
//          "/godirwalk/examples/find-fast" 0 []
//          "/godirwalk/examples/remove-empty-directories" 0 []
//          "/godirwalk/examples/scanner" 0 []
//          "/godirwalk/examples/sizes" 0 []
//          "/godirwalk/examples/walk-fast" 0 []
//          "/godirwalk/examples/walk-stdlib" 0 []
//    "/treeprint" 0 []

//       root: "/Users/shyang/go/src/github.com/shyang107/paw"
//   fullpath: "/Users/shyang/go/src/github.com/shyang107/paw"
//    AbsPath: "/Users/shyang/go/src/github.com/shyang107/paw"
//     RelDir: "-"
//       Name: "paw"
//        Ext: ""
//      IsDir: true

//       root: "/Users/shyang/go/src/github.com/shyang107/paw"
//   fullpath: "/Users/shyang/go/src/github.com/shyang107/paw/afero"
//    AbsPath: "/Users/shyang/go/src/github.com/shyang107/paw/afero"
//     RelDir: ""
//       Name: "afero"
//        Ext: ""
//      IsDir: true

//       root: "/Users/shyang/go/src/github.com/shyang107/paw"
//   fullpath: "/Users/shyang/go/src/github.com/shyang107/paw/afero/afero.go"
//    AbsPath: "/Users/shyang/go/src/github.com/shyang107/paw/afero/afero.go"
//     RelDir: "/afero"
//       Name: "afero.go"
//        Ext: ".go"
//      IsDir: false

//       root: "/Users/shyang/go/src/github.com/shyang107/paw"
//   fullpath: "/Users/shyang/go/src/github.com/shyang107/paw/afero/mem"
//    AbsPath: "/Users/shyang/go/src/github.com/shyang107/paw/afero/mem"
//     RelDir: "/afero"
//       Name: "mem"
//        Ext: ""
//      IsDir: true

//       root: "/Users/shyang/go/src/github.com/shyang107/paw"
//   fullpath: "/Users/shyang/go/src/github.com/shyang107/paw/afero/mem/dir.go"
//    AbsPath: "/Users/shyang/go/src/github.com/shyang107/paw/afero/mem/dir.go"
//     RelDir: "/afero/mem"
//       Name: "dir.go"
//        Ext: ".go"
// 	 IsDir: false
