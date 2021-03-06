package main

import (
	"strings"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
)

func (opt *option) checkByField() {
	lg.Debug(paw.Caller(1))
	var (
		sflag string
		ok    bool
	)
	if opt.isSortNo {
		opt.byField = vfs.SortByNone
		goto END
	}

	opt.byField = vfs.SortByLowerName
	sflag = strings.ToLower(opt.sortByField)
	if len(sflag) == 0 {
		sflag = "lname"
	}
	if opt.isSortBySize &&
		!opt.isSortByMTime &&
		!opt.isSortByName {
		sflag = "size"
	} else if !opt.isSortBySize &&
		opt.isSortByMTime &&
		!opt.isSortByName {
		sflag = "mtime"
	} else if !opt.isSortBySize &&
		!opt.isSortByMTime &&
		opt.isSortByName {
		sflag = "name"
	}
	if opt.isSortReverse {
		sflag += "r"
	}

	lg.WithFields(logrus.Fields{
		"sflag":   sflag,
		"SortKey": vfs.SortShortNameKeys[sflag],
	}).Trace()

	if opt.byField, ok = vfs.SortShortNameKeys[sflag]; !ok {
		opt.byField = vfs.SortByLowerName
	}

END:
	lg.WithFields(logrus.Fields{
		"byField": opt.byField,
	}).Info()
}
