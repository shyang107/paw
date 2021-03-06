package main

import (
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
)

func (opt *option) setVFSOption() {
	lg.Debug(paw.Caller(1))
	vfsOpt := opt.vopt
	vfsOpt = &vfs.VFSOption{
		Depth:      opt.depth,
		Grouping:   opt.grouping,
		ByField:    opt.byField,
		Skips:      opt.skips,
		ViewFields: opt.viewFields,
		ViewType:   opt.viewType,
	}

	lg.WithFields(logrus.Fields{
		"Depth":      vfsOpt.Depth,
		"Grouping":   vfsOpt.Grouping,
		"ByField":    vfsOpt.ByField,
		"Skips":      vfsOpt.Skips,
		"ViewFields": vfsOpt.ViewFields,
		"ViewType":   vfsOpt.ViewType,
	}).Info()
}
