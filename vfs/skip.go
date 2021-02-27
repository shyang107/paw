package vfs

import (
	"io/fs"
	"strings"
)

type SkipConds struct {
	skips []*SkipFunc
}

func NewSkipConds(isSkipHiden bool) *SkipConds {
	s := &SkipConds{
		skips: []*SkipFunc{},
	}
	if isSkipHiden {
		s.Add(&defaultSkip)
	}
	return s
}

func (s *SkipConds) Add(skip *SkipFunc) {
	if s.skips == nil {
		s.skips = []*SkipFunc{}
	}
	s.skips = append(s.skips, skip)
}

func (s *SkipConds) Is(de fs.DirEntry) bool {
	for _, skip := range s.skips {
		if (*skip)(de) {
			return true
		}
	}
	return false
}

type SkipFunc func(de fs.DirEntry) bool

var defaultSkip SkipFunc = func(de fs.DirEntry) bool {
	name := de.Name()
	if strings.HasPrefix(name, ".") {
		return true
	}
	if strings.EqualFold(name, "_gsdata_") {
		return true
	}
	return false
}
