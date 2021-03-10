package main

import (
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
)

func (opt *option) setVFSOption() {
	lg.Debug(paw.Caller(1))
	opt.vopt = &vfs.VFSOption{
		Depth:          opt.depth,
		IsForceRecurse: opt.isForceRecurse,
		Grouping:       opt.grouping,
		ByField:        opt.byField,
		Skips:          opt.skips,
		ViewFields:     opt.viewFields,
		ViewType:       opt.viewType,
	}
	info("settings: {",
		paw.ValuePairA([]*paw.ValuePair{
			paw.NewValuePair("Depth", opt.vopt.Depth),
			paw.NewValuePair("IsForceRecurse", opt.vopt.IsForceRecurse),
			paw.NewValuePair("Grouping", opt.vopt.Grouping),
			paw.NewValuePair("ByField", opt.vopt.ByField),
			paw.NewValuePair("Skips", opt.vopt.Skips),
			paw.NewValuePair("ViewFields", opt.vopt.ViewFields),
			paw.NewValuePair("ViewType", opt.vopt.ViewType),
		}), "}")
}
