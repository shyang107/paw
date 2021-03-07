package main

import (
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
)

func (opt *option) setVFSOption() {
	lg.Debug(paw.Caller(1))
	opt.vopt = &vfs.VFSOption{
		Depth:        opt.depth,
		IsScanAllSub: opt.isDepthScanAllSub,
		Grouping:     opt.grouping,
		ByField:      opt.byField,
		Skips:        opt.skips,
		ViewFields:   opt.viewFields,
		ViewType:     opt.viewType,
	}

	lg.WithFields(logrus.Fields{
		"Depth":        opt.vopt.Depth,
		"IsScanAllSub": opt.vopt.IsScanAllSub,
		"Grouping":     opt.vopt.Grouping,
		"ByField":      opt.vopt.ByField,
		"Skips":        opt.vopt.Skips,
		"ViewFields":   opt.vopt.ViewFields,
		"ViewType":     opt.vopt.ViewType,
	}).Debug()
}
