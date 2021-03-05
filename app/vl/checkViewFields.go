package main

import "github.com/shyang107/paw"

func (opt *option) checkViewFields() {
	lg.Info("TODO..." + paw.Caller(1))

	// var (
	// 	flag   filetree.PDFieldFlag
	// 	fields []string
	// )

	// if opt.isFieldINode {
	// 	flag = flag | filetree.PFieldINode
	// 	fields = append(fields, filetree.PFieldINode.Name())
	// }

	// flag = flag | filetree.PFieldPermissions
	// fields = append(fields, filetree.PFieldPermissions.Name())

	// if opt.isFieldLinks {
	// 	flag = flag | filetree.PFieldLinks
	// 	fields = append(fields, filetree.PFieldLinks.Name())
	// }

	// flag = flag | filetree.PFieldSize
	// fields = append(fields, filetree.PFieldSize.Name())

	// if opt.isFieldBlocks {
	// 	flag = flag | filetree.PFieldBlocks
	// 	fields = append(fields, filetree.PFieldBlocks.Name())
	// }

	// flag = flag | filetree.PFieldUser
	// fields = append(fields, filetree.PFieldUser.Name())

	// flag = flag | filetree.PFieldGroup
	// fields = append(fields, filetree.PFieldGroup.Name())

	// if opt.isFieldModified {
	// 	flag = flag | filetree.PFieldModified
	// 	fields = append(fields, filetree.PFieldModified.Name())
	// }
	// if opt.isFieldAccessed {
	// 	flag = flag | filetree.PFieldAccessed
	// 	fields = append(fields, filetree.PFieldAccessed.Name())
	// }
	// if opt.isFieldCreated {
	// 	flag = flag | filetree.PFieldCreated
	// 	fields = append(fields, filetree.PFieldCreated.Name())
	// }
	// if !opt.isFieldModified &&
	// 	!opt.isFieldAccessed &&
	// 	!opt.isFieldCreated {
	// 	flag = flag | filetree.PFieldModified
	// 	fields = append(fields, filetree.PFieldModified.Name())
	// }

	// if opt.isFieldMd5 {
	// 	flag = flag | filetree.PFieldMd5
	// 	fields = append(fields, filetree.PFieldMd5.Name())
	// }

	// if opt.isFieldGit {
	// 	flag = flag | filetree.PFieldGit
	// 	fields = append(fields, filetree.PFieldGit.Name())
	// }

	// fields = append(fields, filetree.PFieldName.Name())
	// lg.WithFields(logrus.Fields{
	// 	"N":      len(fields),
	// 	"fields": fields,
	// }).Trace("fields" + paw.Caller(1))

	// return flag
}
