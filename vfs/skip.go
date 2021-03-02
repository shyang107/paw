package vfs

import (
	"regexp"
	"strings"

	"github.com/shyang107/paw"
	"github.com/sirupsen/logrus"
)

type SkipConds struct {
	skips []Skiper
}

// NewSkipConds creats a new instance of SkipConds
func NewSkipConds() *SkipConds {
	return &SkipConds{
		skips: []Skiper{},
	}
}

// Add adds a new skip to SkipConds
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

type Skiper interface {
	Name() string
	Skip(DirEntryX) bool
}

// SkipFunc is a func used to skip DirEntryX
type SkipFunc struct {
	name string
	skip func(DirEntryX) bool
}

func NewSkipFunc(name string, skip func(DirEntryX) bool) *SkipFunc {
	return &SkipFunc{
		name: name,
		skip: skip,
	}
}

func (s *SkipFunc) Name() string {
	return s.name
}
func (s *SkipFunc) Skip(de DirEntryX) bool {
	// paw.Logger.Debug(s.Name(), ": ", de.Name())
	return s.skip(de)
}

// DefaultSkip is a deault SkipFunc used to skip Name of DirEntryX with prefix of "." or equal to "_gsdata_"
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

// SkipRegularFile skips regular files
var SkipRegularFile = NewSkipFunc("skip regular file", func(de DirEntryX) bool {
	if de.Type().IsRegular() {
		return true
	}
	return false
})

// SkipDir skips directory file
var SkipDir = NewSkipFunc("skip directory", func(de DirEntryX) bool {
	if de.IsDir() {
		return true
	}
	return false
})

// SkipFuncRe is a func to skip DirEntryX using regex
type SkipFuncRe struct {
	name    string
	re      *regexp.Regexp
	pattern string
	skip    func(DirEntryX, *regexp.Regexp) bool
}

func NewSkipFuncRe(name string, pattern string, skip func(DirEntryX, *regexp.Regexp) bool) *SkipFuncRe {
	return &SkipFuncRe{
		name:    name,
		pattern: pattern,
		re:      regexp.MustCompile(pattern),
		skip:    skip,
	}
}

func (s *SkipFuncRe) Name() string {
	return s.name
}
func (s *SkipFuncRe) Skip(de DirEntryX) bool {
	// paw.Logger.Debug(s.Name(), ": ", de.Name())
	return s.skip(de, s.re)
}
