package vfs

import (
	"regexp"
	"strings"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

// SkipConds is skipping condtions during building VFS
// 	see examples/vfs
type SkipConds struct {
	skips []Skiper
}

// NewSkipConds creats a new instance of SkipConds
// 	see examples/vfs
func NewSkipConds() *SkipConds {
	return &SkipConds{
		skips: []Skiper{},
	}
}

func (s SkipConds) String() string {
	ss := make([]string, 0, len(s.skips))
	if len(s.skips) == 0 {
		return "nil"
	}
	for _, skip := range s.skips {
		ss = append(ss, skip.Name())
	}
	return strings.Join(ss, "|")
}

// Add adds a new skip to SkipConds
// 	see examples/vfs
func (s *SkipConds) Add(skips ...Skiper) *SkipConds {
	if skips == nil {
		return s
	}
	if s.skips == nil {
		s.skips = []Skiper{}
	}
	var caller string
	if paw.Logger.IsLevelEnabled(logrus.TraceLevel) {
		caller = paw.Caller(1)
	}
	for _, skip := range skips {
		paw.Logger.Trace("add skiper: ", skip.Name(), caller)
		s.skips = append(s.skips, skip)
	}
	return s
}

// Is returns true for skip
func (s *SkipConds) Is(de DirEntryX) bool {
	for _, s := range s.skips {
		if s.Skip(de) {
			return true
		}
	}
	return false
}

// Skipper is an interface of skipping function of SkipConds
type Skiper interface {
	Name() string
	Skip(DirEntryX) bool
}

// SkipFunc is a func used to skip DirEntryX during building VFS
// 	see examples/vfs
type SkipFunc struct {
	name string
	skip func(DirEntryX) bool
}

// NewSkipFunc returns a new instance of SkipFunc
// 	see examples/vfs
func NewSkipFunc(name string, skip func(DirEntryX) bool) *SkipFunc {
	return &SkipFunc{
		name: name,
		skip: skip,
	}
}

// Name return name of SkipFunc; in genral, message about this SkipFunc.
func (s *SkipFunc) Name() string {
	return s.name
}

// Skip return true to skip file, otherwise not.
func (s *SkipFunc) Skip(de DirEntryX) bool {
	// paw.Logger.Debug(s.Name(), ": ", de.Name())
	return s.skip(de)
}

// DefaultSkip is a deault SkipFunc used to skip Name of DirEntryX with prefix of "." or equal to "_gsdata_"
// 	see examples/vfs
var DefaultSkip = NewSkipFunc("default skip func", func(de DirEntryX) bool {
	name := strings.TrimSpace(de.Name())
	if strings.HasPrefix(name, ".") {
		return true
	}
	if strings.EqualFold(name, "_gsdata_") {
		return true
	}
	return false
})

// SkipFile skips regular files
// 	Another way, use VFSOption.ViewType.NoFiles(); and use VFSOption.ViewType.ViewDirAndFile() back to default show directories and files (excluding ViewListTree and ViewListTreeX).
// 	see examples/vfs
var SkipFile = NewSkipFunc("skip file", func(de DirEntryX) bool {
	if de.IsFile() {
		return true
	}
	return false
})

// SkipDir skips directory file
// 	[Warning] If use SkipDir, any directories under root do not be accessed.
//  Use VFSOption.ViewType.NoDirs() would be great; and use VFSOption.ViewType.ViewDirAndFile() back to default show directories and files (excluding ViewListTree and ViewListTreeX).
// 	see examples/vfs
var SkipDir = NewSkipFunc("skip directory", func(de DirEntryX) bool {
	if de.IsDir() {
		return true
	}
	return false
})

// SkipFuncRe is a func to skip DirEntryX using regex
// 	see examples/vfs
type SkipFuncRe struct {
	name    string
	re      *regexp.Regexp
	pattern string
	skip    func(DirEntryX, *regexp.Regexp) bool
}

// NewSkipFuncRe returns a new instance of SkipFuncRe
// 	see examples/vfs
func NewSkipFuncRe(name string, pattern string, skip func(DirEntryX, *regexp.Regexp) bool) *SkipFuncRe {
	return &SkipFuncRe{
		name:    name,
		pattern: pattern,
		re:      regexp.MustCompile(pattern),
		skip:    skip,
	}
}

// Name return name of SkipFuncRe; in genral, message about this SkipFuncRe.
func (s *SkipFuncRe) Name() string {
	return s.name
}

// Skip return true to skip file, otherwise not.
func (s *SkipFuncRe) Skip(de DirEntryX) bool {
	// paw.Logger.Debug(s.Name(), ": ", de.Name())
	return s.skip(de, s.re)
}
