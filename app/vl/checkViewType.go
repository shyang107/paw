package main

import (
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
)

func (opt *option) checkViewType() {
	lg.Info(paw.Caller(1))
	if opt.isViewListTree {
		if opt.depth == 0 {
			opt.depth = -1
		}
		opt.viewType = vfs.ViewListTree
	} else if opt.isViewTree {
		if opt.depth == 0 {
			opt.depth = -1
		}
		opt.viewType = vfs.ViewTree
	} else if opt.isViewTable {
		opt.viewType = vfs.ViewTable
	} else if opt.isViewLevel {
		opt.viewType = vfs.ViewLevel
	} else if opt.isViewClassify {
		opt.viewType = vfs.ViewClassify
	} else if opt.isViewList {
		opt.viewType = vfs.ViewList
	}
	lg.WithField("viewType", opt.viewType).Trace()

	if opt.isViewX {
		lg.WithField("isViewX", opt.isViewX).Trace()
		if opt.viewType&vfs.ViewClassify == 0 {
			opt.viewType |= vfs.ViewExtended
		}
		lg.WithField("> viewType", opt.viewType).Trace()
	}

	if opt.isDepthRecurse {
		lg.WithField("isDepthRecurse", opt.isDepthRecurse).Trace()
		opt.depth = -1
	}

	lg.WithFields(logrus.Fields{
		"isViewNoDirs":  opt.isViewNoDirs,
		"isViewNoFiles": opt.isViewNoFiles,
	}).Trace()
	if opt.isViewNoDirs && !opt.isViewNoFiles {
		switch opt.viewType {
		case vfs.ViewList, vfs.ViewLevel, vfs.ViewTable, vfs.ViewClassify,
			vfs.ViewListX, vfs.ViewLevelX, vfs.ViewTableX:
			opt.viewType |= vfs.ViewNoDirs
		}
		lg.WithField("> viewType", opt.viewType).Trace()
	}

	if !opt.isViewNoDirs && opt.isViewNoFiles {
		switch opt.viewType {
		case vfs.ViewList, vfs.ViewLevel, vfs.ViewTable, vfs.ViewClassify,
			vfs.ViewListX, vfs.ViewLevelX, vfs.ViewTableX:
			opt.viewType |= vfs.ViewNoFiles
		}
		lg.WithField("> viewType", opt.viewType).Trace()
	}
	// lg.Debugf("viewType: %v [%d]; ViewLevelXNoFiles: %v [%d]", opt.viewType, opt.viewType, vfs.ViewLevelXNoFiles, vfs.ViewLevelXNoFiles)

	lg.WithFields(logrus.Fields{
		"view":  opt.viewType,
		"depth": opt.depth,
	}).Trace()
}
