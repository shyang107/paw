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
	IsSkip(DirEntryX) bool
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
func (s *SkipConds) IsSkip(de DirEntryX) bool {
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

// SkipFunc is a func used to skip DirEntryX during building VFS
// 	see examples/vfs
type SkipFunc struct {
	name string
	skip func(DirEntryX) bool
}

// NewSkipFunc returns a new instance of SkipFunc
// 	see examples/vfs
func NewSkipFunc(name string, skip func(DirEntryX) bool) Skiper {
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
func (s *SkipFunc) IsSkip(de DirEntryX) bool {
	// paw.Logger.Trace(s.Name(), ": ", de.Name())
	return s.skip(de)
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

// DefaultSkiper is a deault SkipFunc used to skip DirEntryX
//
// Skip Condition, Name of DirEntryX is:
// 	1. with prefix of "."
// 	2. "_gsdata_"
// 	see examples/vfs
var DefaultSkiper = NewSkipFunc("«DefaultSkiper»", func(de DirEntryX) bool {
	if SkiperHiddens.IsSkip(de) {
		return true
	}

	if strings.EqualFold(strings.TrimSpace(de.Name()), "_gsdata_") {
		return true
	}

	return false
})

// SkiperHiddens is a SkipFunc used to skip Name of DirEntryX with prefix of "."
// 	see examples/vfs
var SkiperHiddens = NewSkipFunc("«SkiperHiddens»", func(de DirEntryX) bool {
	name := strings.TrimSpace(de.Name())
	if strings.HasPrefix(name, ".") {
		return true
	}
	return false
})

// SkiperHasNames is a SkipFunc used to skip Name of DirEntryX in SkipNames
// 	see examples/vfs
var SkiperHasNames = NewSkipFunc("«SkiperHasNames»", func(de DirEntryX) bool {
	name := strings.TrimSpace(de.Name())
	for _, skipname := range SkipNames {
		if strings.EqualFold(name, skipname) {
			return true
		}
	}
	return false
})

// SkipPrefixer is a SkipFunc used to skip Name of DirEntryX with prefix in SkipPrefix
// 	see examples/vfs
var SkipPrefixer = NewSkipFunc("«SkipPrefixer»", func(de DirEntryX) bool {
	name := strings.ToLower(strings.TrimSpace(de.Name()))
	for _, prefix := range SkipPrefix {
		if strings.HasPrefix(name, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
})

// SkipSuffixer is a SkipFunc used to skip Name of DirEntryX with prefix in SkipPrefix
// 	see examples/vfs
var SkipSuffixer = NewSkipFunc("«SkipSuffixer»", func(de DirEntryX) bool {
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
var SkipFiler = NewSkipFunc("«SkipFiler»", func(de DirEntryX) bool {
	if de.IsFile() {
		return true
	}
	return false
})

// SkipDirer skips directory file
// 	[Warning] If use SkipDirer, any directories under root do not be accessed.
//  Use VFSOption.ViewType.NoDirs() would be great; and use VFSOption.ViewType.ViewDirAndFile() back to default show directories and files (excluding ViewListTree and ViewListTreeX).
// 	see examples/vfs
var SkipDirer = NewSkipFunc("«SkipDirer»", func(de DirEntryX) bool {
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
func NewSkipFuncRe(name string, pattern string, skip func(d DirEntryX, re *regexp.Regexp) bool) Skiper {
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

// IsSkip return true to skip file, otherwise not.
func (s *SkipFuncRe) IsSkip(de DirEntryX) bool {
	// paw.Logger.Trace(s.Name(), ": ", de.Name())
	return s.skip(de, s.re)
}
