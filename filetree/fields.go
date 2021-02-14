package filetree

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/spf13/cast"
)

type PDFieldFlag int

const (
	// PFieldINode is inode field
	PFieldINode PDFieldFlag = 1 << iota
	// PFieldPermissions is permission field
	PFieldPermissions
	// PFieldLinks is hard link field
	PFieldLinks
	// PFieldSize is size field
	PFieldSize
	// PFieldBlocks is blocks field
	PFieldBlocks
	// PFieldUser is user field
	PFieldUser
	// PFieldGroup is group field
	PFieldGroup
	// PFieldModified is date modified field
	PFieldModified
	// PFieldAccessed is date accessed field
	PFieldAccessed
	// PFieldCreated is date created field
	PFieldCreated
	// PFieldGit is git field
	PFieldGit
	// PFieldMd5 is md5 field
	PFieldMd5
	// PFieldName is name field
	PFieldName
	// PFieldNone is non-default field
	PFieldNone

	// PFieldDefault useas default fields
	PFieldDefault = PFieldPermissions | PFieldSize | PFieldUser | PFieldGroup | PFieldModified | PFieldName
)

func (f PDFieldFlag) String() string {
	switch f {
	case PFieldINode:
		return "inode"
	case PFieldPermissions:
		return "Permissions"
	case PFieldLinks:
		return "Links"
	case PFieldSize:
		return "Size"
	case PFieldBlocks:
		return "Blocks"
	case PFieldUser:
		return "User"
	case PFieldGroup:
		return "Group"
	case PFieldModified:
		return "Modified"
	case PFieldCreated:
		return "Created"
	case PFieldAccessed:
		return "Accessed"
	case PFieldGit:
		return "Git"
	case PFieldMd5:
		return "md5"
	case PFieldName:
		return "Name"
	case PFieldDefault:
		//PFieldPermissions | PFieldSize | PFieldUser |
		// PFieldGroup | PFieldModified | PFieldName
		return "Permissions, Size, User, Group, Modified, Name"
	default:
		return ""
	}
}

func (f PDFieldFlag) Name() string {
	return f.String()
}

func (f PDFieldFlag) Width() int {
	wd := len(f.String())
	switch f {
	case PFieldINode:
		return paw.MaxInt(5, wd)
	case PFieldPermissions:
		return paw.MaxInt(11, wd)
	case PFieldLinks:
		return paw.MaxInt(2, wd)
	case PFieldSize:
		return paw.MaxInt(4, wd)
	case PFieldBlocks:
		return paw.MaxInt(6, wd)
	case PFieldUser:
		return paw.MaxInt(4, wd)
	case PFieldGroup:
		return paw.MaxInt(5, wd)
	case PFieldModified:
		return paw.MaxInt(11, wd)
	case PFieldCreated:
		return paw.MaxInt(11, wd)
	case PFieldAccessed:
		return paw.MaxInt(11, wd)
	case PFieldGit:
		return paw.MaxInt(2, wd)
	case PFieldMd5:
		return paw.MaxInt(32, wd)
	case PFieldName:
		return paw.MaxInt(4, wd)
	default:
		return 0
	}
}

func (f PDFieldFlag) Color() *color.Color {
	switch f {
	case PFieldINode:
		return cinp
	case PFieldPermissions:
		return cpms
	case PFieldLinks:
		return clkp
	case PFieldSize:
		return csnp
	case PFieldBlocks:
		return cbkp
	case PFieldUser:
		return cuup
	case PFieldGroup:
		return cgup
	case PFieldModified:
		return cdap
	case PFieldCreated:
		return cdap
	case PFieldAccessed:
		return cdap
	case PFieldGit:
		return cgitp
	case PFieldMd5:
		return cmd5p
	case PFieldName:
		return cnop
	default:
		return cdashp
	}
}

func (f PDFieldFlag) Align() paw.Align {
	switch f {
	case PFieldINode:
		return paw.AlignRight
	case PFieldPermissions:
		return paw.AlignLeft
	case PFieldLinks:
		return paw.AlignRight
	case PFieldSize:
		return paw.AlignRight
	case PFieldBlocks:
		return paw.AlignRight
	case PFieldUser:
		return paw.AlignLeft
	case PFieldGroup:
		return paw.AlignLeft
	case PFieldModified:
		return paw.AlignLeft
	case PFieldCreated:
		return paw.AlignLeft
	case PFieldAccessed:
		return paw.AlignLeft
	case PFieldGit:
		return paw.AlignRight
	case PFieldMd5:
		return paw.AlignLeft
	case PFieldName:
		return paw.AlignLeft
	case PFieldNone:
		return paw.AlignRight
	default:
		return paw.AlignLeft
	}
}

func (f PDFieldFlag) Field() *Field {
	return NewField(f)
}

var (
	// pfields = []string{}
	// pfieldsMap = map[PDFieldFlag]string{
	// 	PFieldINode:       "inode",
	// 	PFieldPermissions: "Permissions",
	// 	PFieldLinks:       "Links",
	// 	PFieldSize:        "Size",
	// 	PFieldBlocks:      "Blocks",
	// 	PFieldUser:        "User",
	// 	PFieldGroup:       "Group",
	// 	PFieldModified:    "Modified",
	// 	PFieldCreated:     "Created",
	// 	PFieldAccessed:    "Accessed",
	// 	// PFieldModified:    "Date Modified",
	// 	// PFieldCreated:     "Date Created",
	// 	// PFieldAccessed:    "Date Accessed",
	// 	PFieldGit:  "Git",
	// 	PFieldName: "Name",
	// 	PFieldNone: "",
	// }
	// pfieldWidths = []int{}
	// pfieldWidthsMap = map[PDFieldFlag]int{
	// 	// PFieldINode:       paw.MaxInt(8, len(pfieldsMap[PFieldINode])),
	// 	PFieldINode:       paw.MaxInt(5, len(pfieldsMap[PFieldINode])),
	// 	PFieldPermissions: paw.MaxInt(11, len(pfieldsMap[PFieldPermissions])),
	// 	PFieldLinks:       paw.MaxInt(2, len(pfieldsMap[PFieldLinks])),
	// 	PFieldSize:        paw.MaxInt(4, len(pfieldsMap[PFieldSize])),
	// 	PFieldBlocks:      paw.MaxInt(6, len(pfieldsMap[PFieldBlocks])),
	// 	PFieldUser:        paw.MaxInt(4, len(pfieldsMap[PFieldUser])),
	// 	PFieldGroup:       paw.MaxInt(5, len(pfieldsMap[PFieldGroup])),
	// 	PFieldModified:    paw.MaxInt(11, len(pfieldsMap[PFieldModified])),
	// 	PFieldCreated:     paw.MaxInt(11, len(pfieldsMap[PFieldCreated])),
	// 	PFieldAccessed:    paw.MaxInt(11, len(pfieldsMap[PFieldAccessed])),
	// 	PFieldGit:         paw.MaxInt(2, len(pfieldsMap[PFieldGit])),
	// 	PFieldName:        paw.MaxInt(4, len(pfieldsMap[PFieldName])),
	// 	PFieldNone:        0,
	// }

	// pfieldKeys = []PDFieldFlag{}

	DefaultPDFieldKeys = []PDFieldFlag{PFieldPermissions, PFieldSize, PFieldUser, PFieldGroup, PFieldModified, PFieldName}

	DefaultPDFields = NewFields(DefaultPDFieldKeys...)

	// pfieldCPMap = map[PDFieldFlag]*color.Color{
	// 	PFieldINode:       cinp,
	// 	PFieldPermissions: cpms,
	// 	PFieldLinks:       clkp,
	// 	PFieldSize:        csnp,
	// 	// PFieldBlocks:      cbkp,
	// 	PFieldBlocks:   cbkp,
	// 	PFieldUser:     cuup,
	// 	PFieldGroup:    cgup,
	// 	PFieldModified: cdap,
	// 	PFieldCreated:  cdap,
	// 	PFieldAccessed: cdap,
	// 	PFieldGit:      cgitp,
	// 	PFieldName:     cfip,
	// 	PFieldNone:     cdashp,
	// }
	// pfieldAlignMap = map[PDFieldFlag]paw.Align{
	// 	PFieldINode:       paw.AlignRight,
	// 	PFieldPermissions: paw.AlignLeft,
	// 	PFieldLinks:       paw.AlignRight,
	// 	PFieldSize:        paw.AlignRight,
	// 	PFieldBlocks:      paw.AlignRight,
	// 	PFieldUser:        paw.AlignLeft,
	// 	PFieldGroup:       paw.AlignLeft,
	// 	// PFieldUser:        paw.AlignRight,
	// 	// PFieldGroup:       paw.AlignRight,
	// 	PFieldModified: paw.AlignLeft,
	// 	PFieldCreated:  paw.AlignLeft,
	// 	PFieldAccessed: paw.AlignLeft,
	// 	PFieldGit:      paw.AlignRight,
	// 	PFieldName:     paw.AlignLeft,
	// 	PFieldNone:     paw.AlignRight,
	// }
	// FieldsMap = map[PDFieldFlag]*Field{
	// 	PFieldINode:       NewField(PFieldINode),
	// 	PFieldPermissions: NewField(PFieldPermissions),
	// 	PFieldLinks:       NewField(PFieldLinks),
	// 	PFieldSize:        NewField(PFieldSize),
	// 	PFieldBlocks:      NewField(PFieldBlocks),
	// 	PFieldUser:        NewField(PFieldUser),
	// 	PFieldGroup:       NewField(PFieldGroup),
	// 	PFieldModified:    NewField(PFieldModified),
	// 	PFieldCreated:     NewField(PFieldCreated),
	// 	PFieldAccessed:    NewField(PFieldAccessed),
	// 	PFieldGit:         NewField(PFieldGit),
	// 	PFieldName:        NewField(PFieldName),
	// 	PFieldNone:        NewField(PFieldNone),
	// }
)

// Field stores content of a field
//
// Elements:
// 	Name: name of field
// 	NameC: colorful name of field
// 	Width: number of name on console
// 	Value: value of the field
// 	ValueC: colorfulString of value of the field
// 	ValueColor: *color.Color use to create colorful srtring for value;no default color, use SetValueColor to setup
// 	HeadColor: *color.Color use to create colorful srtring for head; has default color, use SetHeadColor to setup
type Field struct {
	Key        PDFieldFlag
	Name       string
	Width      int
	widthMajor int // use in size field for Dev or CharDev
	widthMinor int // use in size field for Dev or CharDev
	Value      interface{}
	ValueC     interface{}
	Align      paw.Align
	ValueColor *color.Color
	HeadColor  *color.Color
	isLink     bool
}

// NewField will return *Field
func NewField(flag PDFieldFlag) *Field {
	return &Field{
		Key:        flag,
		Name:       flag.Name(),  //pfieldsMap[flag],
		Width:      flag.Width(), //pfieldWidthsMap[flag],
		widthMajor: 0,
		widthMinor: 0,
		Value:      nil,
		ValueC:     nil,
		ValueColor: flag.Color(), // pfieldCPMap[flag],
		Align:      flag.Align(), //pfieldAlignMap[flag],
		HeadColor:  chdp,
		isLink:     false,
	}
}

// NewFields will return []*Field
func NewFields(flags ...PDFieldFlag) []*Field {
	if len(flags) == 0 {
		return nil
	}
	dFields := make([]*Field, 0, len(flags))
	for _, f := range flags {
		dFields = append(dFields, NewField(f))
	}
	return dFields
}

// NewFieldsGit will return []*Field w.r.t. git status
func NewFieldsGit(noGit bool, flags ...PDFieldFlag) []*Field {
	if len(flags) == 0 {
		return nil
	}
	if noGit {
		irmGit := -1
		for i, f := range flags {
			if f == PFieldGit {
				irmGit = i
			}
		}
		if irmGit != -1 {
			flags = append(flags[:irmGit], flags[irmGit+1:]...)
		}
	}

	dFields := make([]*Field, 0, len(flags))
	for _, f := range flags {
		dFields = append(dFields, NewField(f))
	}
	return dFields
}

// SetValue sets up Field.Value
func (f *Field) SetValue(value interface{}) {
	f.Value = value
}

// SetIsLink sets up Field.isLink
func (f *Field) SetIsLink(isLink bool) {
	f.isLink = isLink
}

// SetValueC sets up colorful value of Field.Value
func (f *Field) SetValueC(value interface{}) {
	f.ValueC = value
}

// SetValueColor sets up color of Field.Value
func (f *Field) SetValueColor(c *color.Color) {
	f.ValueColor = c
}

// GetValueColor returns color of Field.Value
func (f *Field) GetValueColor(c *color.Color) *color.Color {
	return f.ValueColor
}

// SetHeadColor sets up color of Field.Name
func (f *Field) SetHeadColor(c *color.Color) {
	f.HeadColor = c
}

// GetHeadColor returns color of Field.Name
func (f *Field) GetHeadColor(c *color.Color) *color.Color {
	return f.HeadColor
}

// ValueString will return string of Field.Value
func (f *Field) ValueString() string {
	s := alignedSring(f.Value, f.Align, f.Width)
	return s
}

func alignedSring(value interface{}, align paw.Align, width int) string {
	// wf := StringWidth(value)
	s := strings.TrimSpace(fmt.Sprintf("%v", value))
	ws := paw.StringWidth(s)
	if ws > width {
		return s
	}
	// fmt.Println("width =", width, "wf =", wf)
	switch align {
	case paw.AlignRight:
		// s = paw.Spaces(width-ws) + s
		s = fmt.Sprintf("%[1]*[2]s", width, s)
	case paw.AlignCenter:
		wsl := (width - ws) / 2
		wsr := width - ws - wsl
		s = paw.Spaces(wsl) + s + paw.Spaces(wsr)
	default: //AlignLeft
		// s = s + paw.Spaces(width-ws)
		s = fmt.Sprintf("%-[1]*[2]s", width, s)
	}

	return s
}

// ValueStringC will colorful string of Field.Value
func (f *Field) ValueStringC() string {
	if f.ValueC != nil {
		// return fmt.Sprintf("%v", f.ValueC)
		return cast.ToString(f.ValueC)
	}

	s := f.ValueString()
	if f.ValueColor != nil {
		return f.ValueColor.Sprint(s)
	}
	return s
}

// HeadString will return string of Field.Name with width Field.Width
func (f *Field) HeadString() string {
	return alignedSring(f.Name, f.Align, f.Width)
}

// HeadStringC will return colorful string of Field.Name with width Field.Width as see
func (f *Field) HeadStringC() string {
	s := f.HeadString()
	return f.HeadColor.Sprint(s)
}
