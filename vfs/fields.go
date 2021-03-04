package vfs

import (
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
)

type ViewField int

const (
	// ViewFieldNo is No. field
	ViewFieldNo ViewField = 1 << iota
	// ViewFieldINode is inode field
	ViewFieldINode
	// ViewFieldPermissions is permission field
	ViewFieldPermissions
	// ViewFieldLinks is hard link field
	ViewFieldLinks
	// ViewFieldSize is size field
	ViewFieldSize
	_ViewFieldMajor
	_ViewFieldMinor
	// ViewFieldBlocks is blocks field
	ViewFieldBlocks
	// ViewFieldUser is user field
	ViewFieldUser
	// ViewFieldGroup is group field
	ViewFieldGroup
	// ViewFieldModified is date modified field
	ViewFieldModified
	// ViewFieldAccessed is date accessed field
	ViewFieldAccessed
	// ViewFieldCreated is date created field
	ViewFieldCreated
	// ViewFieldGit is git field
	ViewFieldGit
	// ViewFieldMd5 is md5 field
	ViewFieldMd5
	// ViewFieldName is name field
	ViewFieldName

	// ViewFieldDefault useas default fields
	DefaultViewField = ViewFieldPermissions | ViewFieldSize | ViewFieldUser | ViewFieldGroup | ViewFieldModified | ViewFieldName

	DefaultViewFieldAll = ViewFieldINode | ViewFieldPermissions | ViewFieldLinks | ViewFieldSize | ViewFieldBlocks | ViewFieldUser | ViewFieldGroup | ViewFieldModified | ViewFieldAccessed | ViewFieldCreated | ViewFieldGit | ViewFieldMd5 | ViewFieldName

	DefaultViewFieldAllNoGit = ViewFieldINode | ViewFieldPermissions | ViewFieldLinks | ViewFieldSize | ViewFieldBlocks | ViewFieldUser | ViewFieldGroup | ViewFieldModified | ViewFieldAccessed | ViewFieldCreated | ViewFieldMd5 | ViewFieldName

	DefaultViewFieldAllNoMd5 = ViewFieldINode | ViewFieldPermissions | ViewFieldLinks | ViewFieldSize | ViewFieldBlocks | ViewFieldUser | ViewFieldGroup | ViewFieldModified | ViewFieldAccessed | ViewFieldCreated | ViewFieldGit | ViewFieldName

	DefaultViewFieldAllNoGitMd5 = ViewFieldINode | ViewFieldPermissions | ViewFieldLinks | ViewFieldSize | ViewFieldBlocks | ViewFieldUser | ViewFieldGroup | ViewFieldModified | ViewFieldAccessed | ViewFieldCreated | ViewFieldName
)

var (
	DefaultViewFieldSlice = DefaultViewField.Fields()

	DefaultViewFieldsAllSlice = DefaultViewFieldAll.Fields()

	DefaultViewFieldsAllNoGitSlice = DefaultViewFieldAllNoGit.Fields()

	DefaultViewFieldsAllNoMd5Slice = DefaultViewFieldAllNoMd5.Fields()

	DefaultViewFieldsAllNoGitMd5Slice = DefaultViewFieldAllNoGitMd5.Fields()

	ViewFieldNames = map[ViewField]string{
		ViewFieldNo:          "No.",
		ViewFieldINode:       "inode",
		ViewFieldPermissions: "Permissions",
		ViewFieldLinks:       "Links",
		ViewFieldSize:        "Size",
		_ViewFieldMajor:      "Major",
		_ViewFieldMinor:      "Minor",
		ViewFieldBlocks:      "Blocks",
		ViewFieldUser:        "User",
		ViewFieldGroup:       "Group",
		ViewFieldModified:    "Modified",
		ViewFieldCreated:     "Created",
		ViewFieldAccessed:    "Accessed",
		ViewFieldGit:         "Git",
		ViewFieldMd5:         "md5",
		ViewFieldName:        "Name",
	}

	ViewFieldWidths = map[ViewField]int{
		ViewFieldNo:          3,
		ViewFieldINode:       5,
		ViewFieldPermissions: 11,
		ViewFieldLinks:       2,
		ViewFieldSize:        4,
		_ViewFieldMajor:      0,
		_ViewFieldMinor:      0,
		ViewFieldBlocks:      6,
		ViewFieldUser:        4,
		ViewFieldGroup:       5,
		ViewFieldModified:    11,
		ViewFieldAccessed:    11,
		ViewFieldCreated:     11,
		ViewFieldGit:         3,
		ViewFieldMd5:         32,
		ViewFieldName:        4,
	}

	ViewFieldColors = map[ViewField]*color.Color{
		ViewFieldNo:          cnop,
		ViewFieldINode:       cinp,
		ViewFieldPermissions: cpms,
		ViewFieldLinks:       clkp,
		ViewFieldSize:        csnp,
		ViewFieldBlocks:      cbkp,
		ViewFieldUser:        cuup,
		ViewFieldGroup:       cgup,
		ViewFieldModified:    cdap,
		ViewFieldCreated:     cdap,
		ViewFieldAccessed:    cdap,
		ViewFieldGit:         cgitp,
		ViewFieldMd5:         cmd5p,
		ViewFieldName:        cnop,
	}

	ViewFieldAligns = map[ViewField]paw.Align{
		ViewFieldNo:          paw.AlignLeft,
		ViewFieldINode:       paw.AlignRight,
		ViewFieldPermissions: paw.AlignLeft,
		ViewFieldLinks:       paw.AlignRight,
		ViewFieldSize:        paw.AlignRight,
		ViewFieldBlocks:      paw.AlignRight,
		ViewFieldUser:        paw.AlignLeft,
		ViewFieldGroup:       paw.AlignLeft,
		ViewFieldModified:    paw.AlignLeft,
		ViewFieldCreated:     paw.AlignLeft,
		ViewFieldAccessed:    paw.AlignLeft,
		ViewFieldGit:         paw.AlignRight,
		ViewFieldMd5:         paw.AlignLeft,
		ViewFieldName:        paw.AlignLeft,
	}

	ViewFieldValues = map[ViewField]interface{}{
		ViewFieldNo:          "",
		ViewFieldINode:       "",
		ViewFieldPermissions: "",
		ViewFieldLinks:       "",
		ViewFieldSize:        "",
		ViewFieldBlocks:      "",
		ViewFieldUser:        "",
		ViewFieldGroup:       "",
		ViewFieldModified:    "",
		ViewFieldCreated:     "",
		ViewFieldAccessed:    "",
		ViewFieldGit:         "",
		ViewFieldMd5:         "",
		ViewFieldName:        "",
	}
)

func (f ViewField) String() string {
	if name, ok := ViewFieldNames[f]; ok {
		return name
	} else {
		names := f.Names()
		return strings.Join(names, "|")
	}
}

func (f ViewField) SetName(name string) {
	ViewFieldNames[f] = name
}

func (f ViewField) Name() string {
	return f.String()
}

func (f ViewField) SetWidth(wd int) {
	ViewFieldWidths[f] = wd
}

func (f ViewField) Width() int {
	wd := paw.StringWidth(f.String())
	if dwd, ok := ViewFieldWidths[f]; ok {
		return paw.MaxInt(dwd, wd)
	} else {
		return wd
	}
	// switch f {
	// case ViewFieldINode:
	// 	return paw.MaxInt(5, wd)
	// case ViewFieldPermissions:
	// 	return paw.MaxInt(11, wd)
	// case ViewFieldLinks:
	// 	return paw.MaxInt(2, wd)
	// case ViewFieldSize:
	// 	return paw.MaxInt(4, wd)
	// case ViewFieldBlocks:
	// 	return paw.MaxInt(6, wd)
	// case ViewFieldUser:
	// 	return paw.MaxInt(4, wd)
	// case ViewFieldGroup:
	// 	return paw.MaxInt(5, wd)
	// case ViewFieldModified:
	// 	return paw.MaxInt(11, wd)
	// case ViewFieldCreated:
	// 	return paw.MaxInt(11, wd)
	// case ViewFieldAccessed:
	// 	return paw.MaxInt(11, wd)
	// case ViewFieldGit:
	// 	return paw.MaxInt(2, wd)
	// case ViewFieldMd5:
	// 	return paw.MaxInt(32, wd)
	// case ViewFieldName:
	// 	return paw.MaxInt(4, wd)
	// default:
	// 	return 0
	// }
}

func (f ViewField) SetAlign(align paw.Align) {
	ViewFieldAligns[f] = align
}

func (f ViewField) Align() paw.Align {
	if a, ok := ViewFieldAligns[f]; ok {
		return a
	} else {
		return paw.AlignLeft
	}
}

func (f ViewField) SetColor(color *color.Color) {
	ViewFieldColors[f] = color
}

func (f ViewField) Color() *color.Color {
	if c, ok := ViewFieldColors[f]; ok {
		return c
	} else {
		return cdashp
	}
	// switch f {
	// case ViewFieldINode:
	// 	return cinp
	// case ViewFieldPermissions:
	// 	return cpms
	// case ViewFieldLinks:
	// 	return clkp
	// case ViewFieldSize:
	// 	return csnp
	// case ViewFieldBlocks:
	// 	return cbkp
	// case ViewFieldUser:
	// 	return cuup
	// case ViewFieldGroup:
	// 	return cgup
	// case ViewFieldModified:
	// 	return cdap
	// case ViewFieldCreated:
	// 	return cdap
	// case ViewFieldAccessed:
	// 	return cdap
	// case ViewFieldGit:
	// 	return cgitp
	// case ViewFieldMd5:
	// 	return cmd5p
	// case ViewFieldName:
	// 	return cnop
	// default:
	// 	return cdashp
	// }
}

func (f ViewField) SetValue(value interface{}) {
	ViewFieldValues[f] = value
}

func (f ViewField) Value() interface{} {
	return ViewFieldValues[f]
}

func (f ViewField) Slice() (fields []ViewField, names []string, nameWidths []int) {

	fields = []ViewField{}
	names = []string{}
	nameWidths = []int{}

	if f&ViewFieldNo != 0 {
		fields = append(fields, ViewFieldNo)
	}

	if f&ViewFieldINode != 0 {
		fields = append(fields, ViewFieldINode)
	}

	if f&ViewFieldPermissions != 0 {
		fields = append(fields, ViewFieldPermissions)
	}

	if f&ViewFieldLinks != 0 {
		fields = append(fields, ViewFieldLinks)
	}

	if f&ViewFieldSize != 0 {
		fields = append(fields, ViewFieldSize)
	}

	if f&ViewFieldBlocks != 0 {
		fields = append(fields, ViewFieldBlocks)
	}

	if f&ViewFieldUser != 0 {
		fields = append(fields, ViewFieldUser)
	}

	if f&ViewFieldGroup != 0 {
		fields = append(fields, ViewFieldGroup)
	}

	if f&ViewFieldModified != 0 {
		fields = append(fields, ViewFieldModified)
	}
	if f&ViewFieldCreated != 0 {
		fields = append(fields, ViewFieldCreated)
	}
	if f&ViewFieldAccessed != 0 {
		fields = append(fields, ViewFieldAccessed)
	}

	if f&ViewFieldMd5 != 0 {
		hasMd5 = true
		fields = append(fields, ViewFieldMd5)
	}

	if f&ViewFieldGit != 0 {
		fields = append(fields, ViewFieldGit)
	}

	fields = append(fields, ViewFieldName)
	// if f&ViewFieldName != 0 {
	// 	fields = append(fields, ViewFieldName)
	// }

	sort.Slice(fields, func(i, j int) bool {
		return int(fields[i]) < int(fields[i])
	})

	for _, k := range fields {
		names = append(names, k.Name())
		nameWidths = append(nameWidths, k.Width())
	}
	return fields, names, nameWidths
}

func (f ViewField) Fields() (fields []ViewField) {
	fields, _, _ = f.Slice()
	return fields
}

func (f ViewField) Names() (names []string) {
	_, names, _ = f.Slice()
	return names
}

func (f ViewField) Widths() (widths []int) {
	_, _, widths = f.Slice()
	return widths
}

// IsOk returns true for effective and otherwise not. In genernal, use it in checking.
func (f ViewField) IsOk() (ok bool) {
	paw.Logger.Trace("checking ViewField..." + paw.Caller(1))

	if f&ViewFieldINode != 0 ||
		f&ViewFieldPermissions != 0 ||
		f&ViewFieldLinks != 0 ||
		f&ViewFieldSize != 0 ||
		f&ViewFieldBlocks != 0 ||
		f&ViewFieldUser != 0 ||
		f&ViewFieldGroup != 0 ||
		f&ViewFieldModified != 0 ||
		f&ViewFieldCreated != 0 ||
		f&ViewFieldAccessed != 0 ||
		f&ViewFieldMd5 != 0 ||
		f&ViewFieldGit != 0 ||
		f&ViewFieldName != 0 ||
		f&ViewFieldNo != 0 {
		return true
	} else {
		return false
	}
}

func getPFHeadS(c *color.Color, fields ...ViewField) string {
	hd := ""
	for _, f := range fields {
		if f&ViewFieldName != 0 {
			hd += c.Sprintf("%-[1]*[2]s", f.Width(), f.Name())
			continue
		}
		value := aligned(f, f.Name())
		hd += c.Sprintf("%v", value) + " "
	}
	return hd
}
