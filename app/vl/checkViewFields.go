package main

import (
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/vfs"
	"github.com/sirupsen/logrus"
)

func (opt *option) checkViewFields() {
	lg.Debug(paw.Caller(1))

	var (
		viewFields vfs.ViewField
		isOk       bool = false
		hasBasic        = opt.hasBasicPSUGN
	)

	if opt.hasAll {
		opt.viewFields = vfs.DefaultViewFieldAll
		goto END
	}
	if opt.hasAllNoMd5 {
		opt.viewFields = vfs.DefaultViewFieldAllNoMd5
		goto END
	}
	if opt.hasAllNoGit {
		opt.viewFields = vfs.DefaultViewFieldAllNoGit
		goto END
	}
	if opt.hasAllNoGitMd5 {
		opt.viewFields = vfs.DefaultViewFieldAllNoGitMd5
		goto END
	}

	if opt.hasINode {
		isOk = true
		viewFields |= vfs.ViewFieldINode
	}
	if hasBasic {
		isOk = true
		viewFields |= vfs.ViewFieldPermissions
	} else {
		if opt.hasPermission {
			isOk = true
			viewFields |= vfs.ViewFieldPermissions
		}
	}
	if opt.hasHDLinks {
		isOk = true
		viewFields |= vfs.ViewFieldLinks
	}
	if hasBasic {
		isOk = true
		viewFields |= vfs.ViewFieldSize
	} else {
		if opt.hasSize {
			isOk = true
			viewFields |= vfs.ViewFieldSize
		}
	}
	if opt.hasBlocks {
		isOk = true
		viewFields |= vfs.ViewFieldBlocks
	}
	if hasBasic {
		isOk = true
		viewFields |= vfs.ViewFieldUser
		viewFields |= vfs.ViewFieldGroup
		viewFields |= vfs.ViewFieldModified
	} else {
		if opt.hasUser {
			isOk = true
			viewFields |= vfs.ViewFieldUser
		}
		if opt.hasGroup {
			isOk = true
			viewFields |= vfs.ViewFieldGroup
		}
		if opt.hasMTime {
			isOk = true
			viewFields |= vfs.ViewFieldModified
		}
	}
	if opt.hasCTime {
		isOk = true
		viewFields |= vfs.ViewFieldCreated
	}
	if opt.hasATime {
		isOk = true
		viewFields |= vfs.ViewFieldAccessed
	}
	if opt.hasMd5 {
		isOk = true
		viewFields |= vfs.ViewFieldMd5
	}
	if opt.hasGit {
		isOk = true
		viewFields |= vfs.ViewFieldGit
	}

	viewFields |= vfs.ViewFieldName
	lg.WithFields(logrus.Fields{
		"isOk":       isOk,
		"viewFields": viewFields,
	}).Debug()

	if isOk {
		opt.viewFields = viewFields
		// if viewFields == vfs.ViewFieldName {
		// 	opt.viewFields = vfs.DefaultViewField
		// } else {
		// 	opt.viewFields = viewFields
		// }
	} else {
		opt.viewFields = vfs.DefaultViewField
	}
END:
	lg.WithField("viewFields", opt.viewFields).Info()

}
