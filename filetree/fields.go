package filetree

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
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
		PFieldModified:    "Date Modified",
		PFieldCreated:     "Date Created",
		PFieldAccessed:    "Date Accessed",
		PFieldGit:         "Git",
		PFieldName:        "Name",
		PFieldNone:        "",
	}
	pfieldWidths    = []int{}
	pfieldWidthsMap = map[PDFieldFlag]int{
		PFieldINode:       paw.MaxInt(8, len(pfieldsMap[PFieldINode])),
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
	pfieldKeys  = []PDFieldFlag{}
	pfieldCPMap = map[PDFieldFlag]*color.Color{
		PFieldINode:       cinp,
		PFieldPermissions: nil,
		PFieldLinks:       clkp,
		PFieldSize:        nil,
		// PFieldBlocks:      cbkp,
		PFieldBlocks:   nil,
		PFieldUser:     cuup,
		PFieldGroup:    cgup,
		PFieldModified: cdap,
		PFieldCreated:  cdap,
		PFieldAccessed: cdap,
		PFieldGit:      nil,
		PFieldName:     nil,
		PFieldNone:     nil,
	}
	pfieldAlignMap = map[PDFieldFlag]paw.Align{
		PFieldINode:       paw.AlignRight,
		PFieldPermissions: paw.AlignLeft,
		PFieldLinks:       paw.AlignRight,
		PFieldSize:        paw.AlignRight,
		PFieldBlocks:      paw.AlignRight,
		PFieldUser:        paw.AlignRight,
		PFieldGroup:       paw.AlignRight,
		PFieldModified:    paw.AlignLeft,
		PFieldCreated:     paw.AlignLeft,
		PFieldAccessed:    paw.AlignLeft,
		PFieldGit:         paw.AlignRight,
		PFieldName:        paw.AlignLeft,
		PFieldNone:        paw.AlignRight,
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
type Field struct {
	Key     PDFieldFlag
	Name    string
	Width   int
	Value   interface{}
	valueC  interface{}
	Align   paw.Align
	valuecp *color.Color
	headcp  *color.Color
}

func NewField(flag PDFieldFlag) *Field {
	return &Field{
		Key:     flag,
		Name:    pfieldsMap[flag],
		Width:   pfieldWidthsMap[flag],
		Value:   nil,
		valueC:  nil,
		valuecp: pfieldCPMap[flag],
		Align:   pfieldAlignMap[flag],
		headcp:  chdp,
	}
}

func (f *Field) SetValue(value interface{}) {
	f.Value = value
}

func (f *Field) SetColorfulValue(value interface{}) {
	f.valueC = value
}

func (f *Field) SetValueColor(c *color.Color) {
	f.valuecp = c
}

func (f *Field) ValueString() string {
	s := ""
	switch f.Align {
	case paw.AlignLeft:
		s = fmt.Sprintf("%-[1]*[2]v", f.Width, f.Value)
	default:
		s = fmt.Sprintf("%[1]*[2]v", f.Width, f.Value)
	}
	return s
}

func (f *Field) ColorValueString() string {
	s := f.ValueString()
	if f.valuecp != nil {
		return f.valuecp.Sprint(s)
	}
	if f.valueC != nil {
		return fmt.Sprintf("%v", f.valueC)
	}
	return s
}

func (f *Field) HeadString() string {
	s := ""
	switch f.Align {
	case paw.AlignLeft:
		s = fmt.Sprintf("%-[1]*[2]v", f.Width, f.Name)
	default:
		s = fmt.Sprintf("%[1]*[2]v", f.Width, f.Name)
	}
	return s
}

func (f *Field) ColorHeadString() string {
	s := f.HeadString()
	return f.headcp.Sprint(s)
}

type FieldSlice struct {
	fds []*Field
}

func NewFieldSlice() *FieldSlice {
	f := &FieldSlice{}
	f.fds = []*Field{}
	return f
}

func NewFieldSliceFrom(keys []PDFieldFlag, git GitStatus) (fds *FieldSlice) {
	f := NewFieldSlice()
	for _, k := range keys {
		if k&PFieldGit != 0 && git.NoGit {
			continue
		}
		field := FieldsMap[k]
		f.Add(field)
	}
	return f
}

func fdColorizedSize(size uint64, width int) (csize string) {
	ss := ByteSize(size)
	nss := len(ss)
	sn := fmt.Sprintf("%[1]*[2]s", width-1, ss[:nss-1])
	su := paw.ToLower(ss[nss-1:])
	cn := NewEXAColor("sn")
	cu := NewEXAColor("sb")
	csize = cn.Sprint(sn) + cu.Sprint(su)
	return csize
}

func (f *FieldSlice) SetValues(file *File, git GitStatus) {
	for _, fd := range f.fds {
		switch fd.Key {
		case PFieldINode: //"inode",
			fd.SetValue(file.INode())
		case PFieldPermissions: //"Permissions",
			perm := fmt.Sprintf("%v", file.Stat.Mode())
			if len(file.XAttributes) > 0 {
				perm += "@"
			} else {
				perm += " "
			}
			fd.SetValue(perm)
			fd.SetColorfulValue(file.ColorPermission())
		case PFieldLinks: //"Links",
			fd.SetValue(file.NLinks())
		case PFieldSize: //"Size",
			if file.IsDir() {
				fd.SetValue("-")
				fd.SetColorfulValue(cdashp.Sprintf("%[1]*[2]v", fd.Width, "-"))
			} else {
				csize := fdColorizedSize(file.Size, fd.Width)
				fd.SetValue(ByteSize(file.Size))
				fd.SetColorfulValue(csize)
			}
		case PFieldBlocks: //"User",
			if file.IsDir() {
				fd.SetValue("-")
				fd.SetColorfulValue(cdashp.Sprintf("%[1]*[2]v", fd.Width, "-"))
			} else {
				fd.SetValue(file.Blocks())
				fd.SetColorfulValue(cbkp.Sprintf("%[1]*[2]v", fd.Width, file.Blocks()))
			}
		case PFieldUser: //"User",
			fd.SetValue(urname)
		case PFieldGroup: //"Group",
			fd.SetValue(gpname)
		case PFieldModified: //"Date Modified",
			fd.SetValue(DateString(file.ModifiedTime()))
		case PFieldCreated: //"Date Created",
			fd.SetValue(DateString(file.CreatedTime()))
		case PFieldAccessed: //"Date Accessed",
			fd.SetValue(DateString(file.AccessedTime()))
		case PFieldGit: //"Gid",
			if git.NoGit {
				continue
			} else {
				fd.SetValue(getGitStatus(git, file))
				fd.SetColorfulValue(file.ColorGitStatus(git))
			}
		case PFieldName: //"Name",
			fd.SetValue(file.Name())
			fd.SetColorfulValue(file.ColorName())
		}
	}
}

func (f *FieldSlice) Count() int {
	return len(f.fds)
}

func (f *FieldSlice) Add(field *Field) {
	f.fds = append(f.fds, field)
}

func (f *FieldSlice) Remove(key PDFieldFlag) {
	for i, fd := range f.fds {
		if fd.Key&key != 0 {
			if i == f.Count() {
				f.fds = f.fds[:i]
			} else {
				f.fds = append(f.fds[:i], f.fds[i+1:]...)
			}
		}
	}
}

func (f *FieldSlice) Insert(startIndex int, fds ...*Field) {
	if len(fds) == 0 {
		return
	}

	tmp := make([]*Field, f.Count()+len(fds))
	if startIndex < 0 || startIndex > len(f.fds)-1 { // append to tail
		copy(tmp[:f.Count()], f.fds)
		copy(tmp[f.Count():], fds)
		f.fds = make([]*Field, f.Count()+len(fds))
		copy(f.fds, tmp)
		return
	}
	if startIndex == 0 {
		copy(tmp[:len(fds)], fds)
		copy(tmp[len(fds):], f.fds)
		f.fds = make([]*Field, f.Count()+len(fds))
		copy(f.fds, tmp)
		return
	}
	copy(tmp[:startIndex], f.fds[:startIndex])
	copy(tmp[startIndex:startIndex+len(fds)], fds)
	copy(tmp[startIndex+len(fds):], f.fds[startIndex:])
	f.fds = make([]*Field, f.Count()+len(fds))
	copy(f.fds, tmp)
}

func (f *FieldSlice) Get(key PDFieldFlag) *Field {
	for _, fd := range f.fds {
		if fd.Key&key != 0 {
			return fd
		}
	}
	return nil
}

func (f *FieldSlice) Heads() []string {
	hds := make([]string, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		hds[i] = fd.HeadString()
	}
	return hds
}

func (f *FieldSlice) HeadWidths() []int {
	hds := make([]int, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		hds[i] = fd.Width
	}
	return hds
}

func (f *FieldSlice) HeadAligns() []paw.Align {
	hds := make([]paw.Align, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		hds[i] = fd.Align
	}
	return hds
}

func (f *FieldSlice) HeadsString() string {
	sb := new(strings.Builder)
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		if i < f.Count()-1 {
			fmt.Fprintf(sb, "%s ", fd.HeadString())
		} else {
			fmt.Fprintf(sb, "%s", fd.HeadString())
		}
	}
	return sb.String()
}

func (f *FieldSlice) HeadsStringWidth() int {
	hds := f.HeadsString()
	return paw.StringWidth(hds)
}

func (f *FieldSlice) MetaHeadsStringWidth() int {
	wd := 0
	for _, fd := range f.fds {
		if fd.Key&PFieldName != 0 {
			continue
		}
		wd += fd.Width + 1
	}
	return wd - 1
}

func (f *FieldSlice) ColorHeads() []string {
	hds := make([]string, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		hds[i] = fd.ColorHeadString()
	}
	return hds
}

func (f *FieldSlice) ColorHeadsString() string {
	w := new(strings.Builder)
	for i, h := range f.ColorHeads() {
		if i < f.Count() {
			fmt.Fprintf(w, "%s ", h)
		} else {
			fmt.Fprintf(w, "%s", h)
		}
	}
	return w.String()
}

func (f *FieldSlice) Values() []interface{} {
	vals := make([]interface{}, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		vals[i] = fd.Value
	}
	return vals
}

func (f *FieldSlice) ValuesString() []string {
	vals := make([]string, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		vals[i] = fd.ValueString()
	}
	return vals
}

func (f *FieldSlice) ColorValues() []string {
	vals := make([]string, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		vals[i] = fd.ColorValueString()
	}
	return vals
}

func (f *FieldSlice) MetaValues() []string {
	var vals []string
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		switch fd.Key {
		case PFieldName:
			continue
		default:
			vals = append(vals, fd.ValueString())
		}
	}
	return vals
}

func (f *FieldSlice) MetaValuesString() string {
	sb := new(strings.Builder)
	for i, h := range f.MetaValues() {
		if i < f.Count()-1 {
			fmt.Fprintf(sb, "%s ", h)
		} else {
			fmt.Fprintf(sb, "%s", h)
		}
	}
	return sb.String()
}

func (f *FieldSlice) MetaValuesStringWidth() int {
	s := f.MetaValuesString()
	return paw.StringWidth(s)
}

func (f *FieldSlice) ColorMetaValues() []string {
	var vals []string
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		switch fd.Key {
		case PFieldName:
			continue
		default:
			vals = append(vals, fd.ColorValueString())
		}
	}
	return vals
}

func (f *FieldSlice) ColorMetaValuesString() string {
	sb := new(strings.Builder)
	for i, h := range f.ColorMetaValues() {
		if i < f.Count() {
			fmt.Fprintf(sb, "%s ", h)
		} else {
			fmt.Fprintf(sb, "%s", h)
		}
	}
	return sb.String()
}