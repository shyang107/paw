package main

import (
	"regexp"
	"strings"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
)

func (opt *option) checkSkips() {
	lg.Debug(paw.Caller(1))

	opt.skips = vfs.NewSkipConds()

	// All files
	if !opt.isNoSkip {
		opt.skips = opt.skips.Add(vfs.DefaultSkiper)
	}

	// reInclude
	if len(opt.reIncludePattern) > 0 {
		pattern := opt.reIncludePattern
		if lg.IsLevelEnabled(logrus.TraceLevel) {
			re := regexp.MustCompile(pattern)
			lg.WithField("ri.pattern", re.String()).Trace()
		}
		reSkiper := vfs.NewSkipFuncRe("«re-include»", pattern, func(de vfs.DirEntryX, re *regexp.Regexp) bool {
			name := strings.TrimSpace(de.Name())
			if re.MatchString(name) || de.IsDir() {
				return false
			}
			return true
		})
		opt.skips.Add(reSkiper)
	}
	// reExclude
	if len(opt.reExcludePattern) > 0 {
		pattern := opt.reExcludePattern
		if lg.IsLevelEnabled(logrus.TraceLevel) {
			re := regexp.MustCompile(pattern)
			lg.WithField("rx.pattern", re.String()).Trace()
		}
		reSkiper := vfs.NewSkipFuncRe("«re-exclude»", pattern, func(de vfs.DirEntryX, re *regexp.Regexp) bool {
			name := strings.TrimSpace(de.Name())
			if !re.MatchString(name) || de.IsDir() {
				return false
			}
			return true
		})
		opt.skips.Add(reSkiper)
	}

	// prefix
	if len(opt.withNoPrefix) > 0 {
		prefixs := strings.Split(opt.withNoPrefix, opt.psDelimiter)
		for _, prefix := range prefixs {
			opt.skips.AddToSkipPrefix(prefix)
		}
		lg.WithFields(logrus.Fields{
			"preifx":    prefixs,
			"delimiter": opt.psDelimiter,
		}).Trace()
	}
	// suffix
	if len(opt.withNoSufix) > 0 {
		suffixs := strings.Split(opt.withNoSufix, opt.psDelimiter)
		for _, suffix := range suffixs {
			opt.skips.AddToSkipSuffix(suffix)
		}
		lg.WithFields(logrus.Fields{
			"suffix":    suffixs,
			"delimiter": opt.psDelimiter,
		}).Trace()
	}

	// paw.Logger.WithField("skips", opt.skips).Info()
	info(paw.MesageFieldAndValueC("Skiper", opt.skips, logrus.InfoLevel, paw.Cnop, nil))
}
