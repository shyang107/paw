package vfs

import (
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
)

type ViewField int

const (
	// ViewFieldINode is inode field
	ViewFieldINode ViewField = 1 << iota
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
	// ViewFieldNo is No. field
	ViewFieldNo

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
		return strings.Join(names, ", ")
	}
	// switch f {
	// case ViewFieldNo:
	// 	return "No"
	// case ViewFieldINode:
	// 	return "inode"
	// case ViewFieldPermissions:
	// 	return "Permissions"
	// case ViewFieldLinks:
	// 	return "Links"
	// case ViewFieldSize:
	// 	return "Size"
	// case ViewFieldBlocks:
	// 	return "Blocks"
	// case ViewFieldUser:
	// 	return "User"
	// case ViewFieldGroup:
	// 	return "Group"
	// case ViewFieldModified:
	// 	return "Modified"
	// case ViewFieldCreated:
	// 	return "Created"
	// case ViewFieldAccessed:
	// 	return "Accessed"
	// case ViewFieldGit:
	// 	return "Git"
	// case ViewFieldMd5:
	// 	return "md5"
	// case ViewFieldName:
	// 	return "Name"
	// default:
	// 	names := f.Names()
	// 	return strings.Join(names, ", ")
	// }
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
	// switch f {
	// case ViewFieldNo:
	// 	return paw.AlignLeft
	// case ViewFieldINode:
	// 	return paw.AlignRight
	// case ViewFieldPermissions:
	// 	return paw.AlignLeft
	// case ViewFieldLinks:
	// 	return paw.AlignRight
	// case ViewFieldSize:
	// 	return paw.AlignRight
	// case ViewFieldBlocks:
	// 	return paw.AlignRight
	// case ViewFieldUser:
	// 	return paw.AlignLeft
	// case ViewFieldGroup:
	// 	return paw.AlignLeft
	// case ViewFieldModified:
	// 	return paw.AlignLeft
	// case ViewFieldCreated:
	// 	return paw.AlignLeft
	// case ViewFieldAccessed:
	// 	return paw.AlignLeft
	// case ViewFieldGit:
	// 	return paw.AlignRight
	// case ViewFieldMd5:
	// 	return paw.AlignLeft
	// case ViewFieldName:
	// 	return paw.AlignLeft
	// default:
	// 	return paw.AlignLeft
	// }
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
	// if v, ok := ViewFieldValues[f]; ok {
	// 	return v
	// } else {
	// 	return nil
	// }
}

func (f ViewField) Slice() (fields []ViewField, names []string, nameWidths []int) {

	fields = []ViewField{}
	names = []string{}
	nameWidths = []int{}

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

	// if f&ViewFieldName != 0 {
	// 	fields = append(fields, ViewFieldName)
	// }
	fields = append(fields, ViewFieldName)
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

// var DefaultPDViewFields = NewViewFields(DefaultPDViewFieldKeys...)
// func (f ViewField) ViewField() *ViewField {
// 	return NewViewField(f)
// }
// // ViewField stores content of a field
// //
// // Elements:
// // 	Name: name of field
// // 	NameC: colorful name of field
// // 	Width: number of name on console
// // 	Value: value of the field
// // 	ValueC: colorfulString of value of the field
// // 	ValueColor: *color.Color use to create colorful srtring for value;no default color, use SetValueColor to setup
// // 	HeadColor: *color.Color use to create colorful srtring for head; has default color, use SetHeadColor to setup
// type ViewField struct {
// 	Key        ViewField
// 	Name       string
// 	Width      int
// 	widthMajor int // use in size field for Dev or CharDev
// 	widthMinor int // use in size field for Dev or CharDev
// 	Value      interface{}
// 	ValueC     interface{}
// 	Align      paw.Align
// 	ValueColor *color.Color
// 	HeadColor  *color.Color
// 	isLink     bool
// }

// // NewViewField will return *ViewField
// func NewViewField(flag ViewField) *ViewField {
// 	return &ViewField{
// 		Key:        flag,
// 		Name:       flag.Name(),  //pfieldsMap[flag],
// 		Width:      flag.Width(), //pfieldWidthsMap[flag],
// 		widthMajor: 0,
// 		widthMinor: 0,
// 		Value:      nil,
// 		ValueC:     nil,
// 		ValueColor: flag.Color(), // pfieldCPMap[flag],
// 		Align:      flag.Align(), //pfieldAlignMap[flag],
// 		HeadColor:  chdp,
// 		isLink:     false,
// 	}
// }

// // NewViewFields will return []*ViewField
// func NewViewFields(flags ...ViewField) []*ViewField {
// 	if len(flags) == 0 {
// 		return nil
// 	}
// 	dViewFields := make([]*ViewField, 0, len(flags))
// 	for _, f := range flags {
// 		dViewFields = append(dViewFields, NewViewField(f))
// 	}
// 	return dViewFields
// }

// // SetValue sets up ViewField.Value
// func (f *ViewField) SetValue(value interface{}) {
// 	f.Value = value
// }

// // SetIsLink sets up ViewField.isLink
// func (f *ViewField) SetIsLink(isLink bool) {
// 	f.isLink = isLink
// }

// // SetValueC sets up colorful value of ViewField.Value
// func (f *ViewField) SetValueC(value interface{}) {
// 	f.ValueC = value
// }

// // SetValueColor sets up color of ViewField.Value
// func (f *ViewField) SetValueColor(c *color.Color) {
// 	f.ValueColor = c
// }

// // GetValueColor returns color of ViewField.Value
// func (f *ViewField) GetValueColor(c *color.Color) *color.Color {
// 	return f.ValueColor
// }

// // SetHeadColor sets up color of ViewField.Name
// func (f *ViewField) SetHeadColor(c *color.Color) {
// 	f.HeadColor = c
// }

// // GetHeadColor returns color of ViewField.Name
// func (f *ViewField) GetHeadColor(c *color.Color) *color.Color {
// 	return f.HeadColor
// }

// // ValueString will return string of ViewField.Value
// func (f *ViewField) ValueString() string {
// 	s := alignedSring(f.Value, f.Align, f.Width)
// 	return s
// }

// func alignedSring(value interface{}, align paw.Align, width int) string {
// 	// wf := StringWidth(value)
// 	s := strings.TrimSpace(fmt.Sprintf("%v", value))
// 	ws := paw.StringWidth(s)
// 	if ws > width {
// 		return s
// 	}
// 	// fmt.Println("width =", width, "wf =", wf)
// 	switch align {
// 	case paw.AlignRight:
// 		// s = paw.Spaces(width-ws) + s
// 		s = fmt.Sprintf("%[1]*[2]s", width, s)
// 	case paw.AlignCenter:
// 		wsl := (width - ws) / 2
// 		wsr := width - ws - wsl
// 		s = paw.Spaces(wsl) + s + paw.Spaces(wsr)
// 	default: //AlignLeft
// 		// s = s + paw.Spaces(width-ws)
// 		s = fmt.Sprintf("%-[1]*[2]s", width, s)
// 	}

// 	return s
// }

// // ValueStringC will colorful string of ViewField.Value
// func (f *ViewField) ValueStringC() string {
// 	if f.ValueC != nil {
// 		// return fmt.Sprintf("%v", f.ValueC)
// 		return cast.ToString(f.ValueC)
// 	}

// 	s := f.ValueString()
// 	if f.ValueColor != nil {
// 		return f.ValueColor.Sprint(s)
// 	}
// 	return s
// }

// // HeadString will return string of ViewField.Name with width ViewField.Width
// func (f *ViewField) HeadString() string {
// 	return alignedSring(f.Name, f.Align, f.Width)
// }

// // HeadStringC will return colorful string of ViewField.Name with width ViewField.Width as see
// func (f *ViewField) HeadStringC() string {
// 	s := f.HeadString()
// 	return f.HeadColor.Sprint(s)
// }
