package filetree

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/spf13/cast"
)

type PDFieldFlag int

const (
	// PFieldINode uses inode field
	PFieldINode PDFieldFlag = 1 << iota
	// PFieldPermissions uses permission field
	PFieldPermissions
	// PFieldLinks uses hard link field
	PFieldLinks
	// PFieldSize uses size field
	PFieldSize
	// PFieldBlocks uses blocks field
	PFieldBlocks
	// PFieldUser uses user field
	PFieldUser
	// PFieldGroup uses group field
	PFieldGroup
	// PFieldModified uses date modified field
	PFieldModified
	// PFieldAccessed uses date accessed field
	PFieldAccessed
	// PFieldCreated uses date created field
	PFieldCreated
	// PFieldGit uses git field
	PFieldGit
	// PFieldName uses name field
	PFieldName
	// PFieldNone uses non-default field
	PFieldNone
)

var (
	pfields    = []string{}
	pfieldsMap = map[PDFieldFlag]string{
		PFieldINode:       "inode",
		PFieldPermissions: "Permissions",
		PFieldLinks:       "Links",
		PFieldSize:        "Size",
		PFieldBlocks:      "Blocks",
		PFieldUser:        "User",
		PFieldGroup:       "Group",
		PFieldModified:    "Modified",
		PFieldCreated:     "Created",
		PFieldAccessed:    "Accessed",
		// PFieldModified:    "Date Modified",
		// PFieldCreated:     "Date Created",
		// PFieldAccessed:    "Date Accessed",
		PFieldGit:  "Git",
		PFieldName: "Name",
		PFieldNone: "",
	}
	pfieldWidths    = []int{}
	pfieldWidthsMap = map[PDFieldFlag]int{
		// PFieldINode:       paw.MaxInt(8, len(pfieldsMap[PFieldINode])),
		PFieldINode:       paw.MaxInt(5, len(pfieldsMap[PFieldINode])),
		PFieldPermissions: paw.MaxInt(11, len(pfieldsMap[PFieldPermissions])),
		PFieldLinks:       paw.MaxInt(2, len(pfieldsMap[PFieldLinks])),
		PFieldSize:        paw.MaxInt(4, len(pfieldsMap[PFieldSize])),
		PFieldBlocks:      paw.MaxInt(6, len(pfieldsMap[PFieldBlocks])),
		PFieldUser:        paw.MaxInt(paw.StringWidth(urname), len(pfieldsMap[PFieldUser])),
		PFieldGroup:       paw.MaxInt(paw.StringWidth(gpname), len(pfieldsMap[PFieldGroup])),
		PFieldModified:    paw.MaxInt(11, len(pfieldsMap[PFieldModified])),
		PFieldCreated:     paw.MaxInt(11, len(pfieldsMap[PFieldCreated])),
		PFieldAccessed:    paw.MaxInt(11, len(pfieldsMap[PFieldAccessed])),
		PFieldGit:         paw.MaxInt(2, len(pfieldsMap[PFieldGit])),
		PFieldName:        paw.MaxInt(4, len(pfieldsMap[PFieldName])),
		PFieldNone:        0,
	}

	pfieldKeys = []PDFieldFlag{}

	pfieldKeysDefualt = []PDFieldFlag{PFieldPermissions, PFieldSize, PFieldUser, PFieldGroup, PFieldModified, PFieldName}

	pfieldCPMap = map[PDFieldFlag]*color.Color{
		PFieldINode:       cinp,
		PFieldPermissions: cpmp,
		PFieldLinks:       clkp,
		PFieldSize:        csnp,
		// PFieldBlocks:      cbkp,
		PFieldBlocks:   cbkp,
		PFieldUser:     cuup,
		PFieldGroup:    cgup,
		PFieldModified: cdap,
		PFieldCreated:  cdap,
		PFieldAccessed: cdap,
		PFieldGit:      cgitp,
		PFieldName:     cfip,
		PFieldNone:     cnop,
	}
	pfieldAlignMap = map[PDFieldFlag]paw.Align{
		PFieldINode:       paw.AlignRight,
		PFieldPermissions: paw.AlignLeft,
		PFieldLinks:       paw.AlignRight,
		PFieldSize:        paw.AlignRight,
		PFieldBlocks:      paw.AlignRight,
		PFieldUser:        paw.AlignLeft,
		PFieldGroup:       paw.AlignLeft,
		// PFieldUser:        paw.AlignRight,
		// PFieldGroup:       paw.AlignRight,
		PFieldModified: paw.AlignLeft,
		PFieldCreated:  paw.AlignLeft,
		PFieldAccessed: paw.AlignLeft,
		PFieldGit:      paw.AlignRight,
		PFieldName:     paw.AlignLeft,
		PFieldNone:     paw.AlignRight,
	}
	FieldsMap = map[PDFieldFlag]*Field{
		PFieldINode:       NewField(PFieldINode),
		PFieldPermissions: NewField(PFieldPermissions),
		PFieldLinks:       NewField(PFieldLinks),
		PFieldSize:        NewField(PFieldSize),
		PFieldBlocks:      NewField(PFieldBlocks),
		PFieldUser:        NewField(PFieldUser),
		PFieldGroup:       NewField(PFieldGroup),
		PFieldModified:    NewField(PFieldModified),
		PFieldCreated:     NewField(PFieldCreated),
		PFieldAccessed:    NewField(PFieldAccessed),
		PFieldGit:         NewField(PFieldGit),
		PFieldName:        NewField(PFieldName),
		PFieldNone:        NewField(PFieldNone),
	}
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
	Value      interface{}
	ValueC     interface{}
	Align      paw.Align
	ValueColor *color.Color
	HeadColor  *color.Color
}

// NewField will return *Field
func NewField(flag PDFieldFlag) *Field {
	return &Field{
		Key:        flag,
		Name:       pfieldsMap[flag],
		Width:      pfieldWidthsMap[flag],
		Value:      nil,
		ValueC:     nil,
		ValueColor: pfieldCPMap[flag],
		Align:      pfieldAlignMap[flag],
		HeadColor:  chdp,
	}
}

// SetValue sets up Field.Value
func (f *Field) SetValue(value interface{}) {
	f.Value = value
}

// SetColorfulValue sets up colorful value of Field.Value
func (f *Field) SetColorfulValue(value interface{}) {
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
	s := paw.TrimSpace(fmt.Sprintf("%v", value))
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

// ColorValueString will colorful string of Field.Value
func (f *Field) ColorValueString() string {
	if f.ValueC != nil {
		// return fmt.Sprintf("%v", f.ValueC)
		return cast.ToString(f.ValueC)
	}

	s := f.ValueString()
	if f.ValueColor != nil {
		return f.ValueColor.Sprint(s)
	}
	return s
	// if f.ValueC == nil {
	// 	return ""
	// }
	// return fmt.Sprintf("%v", f.ValueC)
}

// HeadString will return string of Field.Name with width Field.Width
func (f *Field) HeadString() string {
	return alignedSring(f.Name, f.Align, f.Width)
	// s := ""
	// switch f.Align {
	// case paw.AlignLeft:
	// 	s = fmt.Sprintf("%-[1]*[2]v", f.Width, f.Name)
	// default:
	// 	s = fmt.Sprintf("%[1]*[2]v", f.Width, f.Name)
	// }
	// return s
}

// ColorHeadString will return colorful string of Field.Name with width Field.Width as see
func (f *Field) ColorHeadString() string {
	s := f.HeadString()
	return f.HeadColor.Sprint(s)
}
