package vfs

import (
	"regexp"
	"strings"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

// Skipper is an interface of skipping function of SkipConds
type Skiper interface {
	Name() string
	IsSkip(DirEntry) bool
}

// SkipConds is skipping condtions during building VFS
// 	see examples/vfs
type SkipConds struct {
	skips []Skiper
}

// NewSkipConds creats a new instance of SkipConds
// 	see examples/vfs
func NewSkipConds(skips ...Skiper) *SkipConds {
	s := &SkipConds{
		skips: make([]Skiper, 0, len(skips)),
	}
	return s.Add(skips...)
}

func (s SkipConds) String() string {
	ss := make([]string, 0, len(s.skips))
	if len(s.skips) == 0 {
		return "«nil»"
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
		s.skips = make([]Skiper, 0, len(skips))
	}
	var caller string
	if paw.Logger.IsLevelEnabled(logrus.TraceLevel) {
		caller = paw.Caller(1)
	}
	for _, skip := range skips {
		if !s.isExistSkiper(skip) {
			paw.Logger.Trace("add skiper: ", skip.Name(), caller)
			s.skips = append(s.skips, skip)
		} else {
			paw.Logger.Warnf("ignore duplicate skiper: %s %s", skip.Name(), paw.Caller(1))
		}
	}
	return s
}

func (s *SkipConds) isExistSkiper(skip Skiper) bool {
	if s.skips == nil || len(s.skips) == 0 {
		return false
	}
	for _, v := range s.skips {
		if v.Name() == skip.Name() {
			return true
		}
	}
	return false
}

// AddToSkipNames appends name to SkipNames
func (s *SkipConds) AddToSkipNames(names ...string) *SkipConds {
	if !s.isExistSkiper(SkiperHasNames) {
		s.Add(SkiperHasNames)
	}
	SkipNames.Add(names...)
	return s
}

// AddToSkipPrefix appends prefix to SkipPrefix
func (s *SkipConds) AddToSkipPrefix(prefixs ...string) *SkipConds {
	if !s.isExistSkiper(SkipPrefixer) {
		s.Add(SkipPrefixer)
	}
	SkipPrefix.Add(prefixs...)
	return s
}

// AddToSkipSuffix appends suffix to SkipSuffix
func (s *SkipConds) AddToSkipSuffix(suffixs ...string) *SkipConds {
	if !s.isExistSkiper(SkipSuffixer) {
		s.Add(SkipSuffixer)
	}
	SkipSuffix.Add(suffixs...)
	return s
}

var _inside_skip bool

// IsSkip returns true for skip
func (s *SkipConds) IsSkip(de DirEntry) bool {
	_inside_skip = true
	if !s.IsOk() {
		goto END
	}
	for _, skipper := range s.skips {
		if skipper.IsSkip(de) {
			return true
		}
	}
END:
	_inside_skip = false
	return false
}

// IsOk returns true for effective and otherwise not. In genernal, use it in checking.
func (s *SkipConds) IsOk() bool {
	if !_inside_skip {
		paw.Logger.Trace("checking SkipConds..." + paw.Caller(1))
	}

	if s.skips == nil {
		return false
	}
	if len(s.skips) == 0 {
		return false
	}
	return true
}

type SkipFunc func(d DirEntry) bool

// Skipper is a func used to skip DirEntry during building VFS
// 	see examples/vfs
type Skipper struct {
	name string
	skip SkipFunc
}

// NewSkipper returns a new instance of Skipper
// 	see examples/vfs
func NewSkipper(name string, skip SkipFunc) Skiper {
	return &Skipper{
		name: name,
		skip: skip,
	}
}

// Name return name of Skipper; in genral, message about this Skipper.
func (s *Skipper) Name() string {
	return s.name
}

// Skip return true to skip file, otherwise not.
func (s *Skipper) IsSkip(d DirEntry) bool {
	// paw.Logger.Trace(s.Name(), ": ", de.Name())
	return s.skip(d)
}

type SkipNamesType []string
type SkipPrefixType []string
type SkipSuffixType []string

var (
	// SkipNames is names needed to skip
	SkipNames SkipNamesType = []string{}
	// SkipPrefix is name with prefix needed to skip
	SkipPrefix SkipPrefixType = []string{}
	// SkipSuffix is name with prefix needed to skip
	SkipSuffix SkipSuffixType = []string{}
)

// Add appends name to SkipNames
func (s *SkipNamesType) Add(names ...string) *SkipNamesType {
	for _, name := range names {
		if !paw.ContainsString((*s), name) {
			(*s) = append((*s), name)
		}
	}
	return s
}

// Add appends prefix to SkipPrefix
func (s *SkipPrefixType) Add(prefixs ...string) *SkipPrefixType {
	for _, name := range prefixs {
		if !paw.ContainsString((*s), name) {
			(*s) = append((*s), name)
		}
	}
	return s
}

// Add appends suffix to SkipSuffix
func (s *SkipSuffixType) Add(suffixs ...string) *SkipSuffixType {
	for _, name := range suffixs {
		if !paw.ContainsString((*s), name) {
			(*s) = append((*s), name)
		}
	}
	return s
}

// DefaultSkiper is a deault Skipper used to skip DirEntry
//
// Skip Condition, Name of DirEntry is:
// 	1. with prefix of "."
// 	2. "_gsdata_"
// 	see examples/vfs
var DefaultSkiper = NewSkipper("«DefaultSkiper»", func(de DirEntry) bool {
	if SkiperHiddens.IsSkip(de) {
		return true
	}
	tname := strings.TrimSpace(de.Name())
	if strings.EqualFold(tname, "_gsdata_") ||
		strings.EqualFold(tname, "$RECYCLE.BIN") {
		return true
	}

	return false
})

// SkiperHiddens is a Skipper used to skip Name of DirEntry with prefix of "."
// 	see examples/vfs
var SkiperHiddens = NewSkipper("«SkiperHiddens»", func(de DirEntry) bool {
	name := strings.TrimSpace(de.Name())
	if strings.HasPrefix(name, ".") {
		return true
	}
	return false
})

// SkiperHasNames is a Skipper used to skip Name of DirEntry in SkipNames
// 	see examples/vfs
var SkiperHasNames = NewSkipper("«SkiperHasNames»", func(de DirEntry) bool {
	name := strings.TrimSpace(de.Name())
	for _, skipname := range SkipNames {
		if strings.EqualFold(name, skipname) {
			return true
		}
	}
	return false
})

// SkipPrefixer is a Skipper used to skip Name of DirEntry with prefix in SkipPrefix
// 	see examples/vfs
var SkipPrefixer = NewSkipper("«SkipPrefixer»", func(de DirEntry) bool {
	name := strings.ToLower(strings.TrimSpace(de.Name()))
	for _, prefix := range SkipPrefix {
		if strings.HasPrefix(name, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
})

// SkipSuffixer is a Skipper used to skip Name of DirEntry with prefix in SkipPrefix
// 	see examples/vfs
var SkipSuffixer = NewSkipper("«SkipSuffixer»", func(de DirEntry) bool {
	name := strings.ToLower(strings.TrimSpace(de.Name()))
	for _, suffix := range SkipSuffix {
		if strings.HasSuffix(name, strings.ToLower(suffix)) {
			return true
		}
	}
	return false
})

// SkipFiler skips regular files
// 	Another way, use VFSOption.ViewType.NoFiles(); and use VFSOption.ViewType.ViewDirAndFile() back to default show directories and files (excluding ViewListTree and ViewListTreeX).
// 	see examples/vfs
var SkipFiler = NewSkipper("«SkipFiler»", func(de DirEntry) bool {
	if !de.IsDir() {
		return true
	}
	return false
})

// SkipDirer skips directory file
// 	[Warning] If use SkipDirer, any directories under root do not be accessed.
//  Use VFSOption.ViewType.NoDirs() would be great; and use VFSOption.ViewType.ViewDirAndFile() back to default show directories and files (excluding ViewListTree and ViewListTreeX).
// 	see examples/vfs
var SkipDirer = NewSkipper("«SkipDirer»", func(de DirEntry) bool {
	if de.IsDir() {
		return true
	}
	return false
})

// SkipperRe is a func to skip DirEntry using regex
// 	see examples/vfs
type SkipperRe struct {
	name    string
	re      *regexp.Regexp
	pattern string
	skip    func(DirEntry, *regexp.Regexp) bool
}

// NewSkipperRe returns a new instance of SkipperRe
// 	see examples/vfs
func NewSkipperRe(name string, pattern string, skip func(d DirEntry, re *regexp.Regexp) bool) Skiper {
	return &SkipperRe{
		name:    name,
		pattern: pattern,
		re:      regexp.MustCompile(pattern),
		skip:    skip,
	}
}

// Name return name of SkipperRe; in genral, message about this SkipperRe.
func (s *SkipperRe) Name() string {
	return s.name
}

// IsSkip return true to skip file, otherwise not.
func (s *SkipperRe) IsSkip(de DirEntry) bool {
	// paw.Logger.Trace(s.Name(), ": ", de.Name())
	return s.skip(de, s.re)
}
