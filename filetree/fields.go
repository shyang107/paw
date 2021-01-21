package filetree

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/shyang107/paw/cast"
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
	s := ""
	switch f.Align {
	case paw.AlignLeft:
		s = fmt.Sprintf("%-[1]*[2]v", f.Width, f.Value)
	default:
		s = fmt.Sprintf("%[1]*[2]v", f.Width, f.Value)
	}
	return s
}

// ColorValueString will colorful string of Field.Value
func (f *Field) ColorValueString() string {
	s := f.ValueString()
	if f.ValueC != nil {
		// return fmt.Sprintf("%v", f.ValueC)
		return cast.ToString(f.ValueC)
	}
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
	s := ""
	switch f.Align {
	case paw.AlignLeft:
		s = fmt.Sprintf("%-[1]*[2]v", f.Width, f.Name)
	default:
		s = fmt.Sprintf("%[1]*[2]v", f.Width, f.Name)
	}
	return s
}

// ColorHeadString will return colorful string of Field.Name with width Field.Width as see
func (f *Field) ColorHeadString() string {
	s := f.HeadString()
	return f.HeadColor.Sprint(s)
}

// FieldSlice is Field union
type FieldSlice struct {
	fds []*Field
}

// NewFieldSlice will return *fieldSlice
func NewFieldSlice() *FieldSlice {
	f := &FieldSlice{}
	f.fds = []*Field{}
	return f
}

// NewFieldSliceFrom will return *fieldSlice created from []PDFieldFlag and GitStatus
func NewFieldSliceFrom(keys []PDFieldFlag, git GitStatus) (fds *FieldSlice) {
	f := NewFieldSlice()
	for _, k := range keys {
		if k == PFieldGit && git.NoGit {
			continue
		}
		field := FieldsMap[k]
		if _, ok := pfieldAlignMap[k]; !ok {
			field.Align = paw.AlignLeft
		}
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

// Fields will return Fields of FieldSlice
func (f *FieldSlice) Fields() []*Field {
	return f.fds
}

// SetValues sets up values of FieldSlice from File and GitStatus
func (f *FieldSlice) SetValues(file *File, git GitStatus) {
	for _, fd := range f.fds {
		switch fd.Key {
		case PFieldINode: //"inode",
			fd.SetValue(file.INode())
			fd.SetValueColor(cinp)
		case PFieldPermissions: //"Permissions",
			perm := fmt.Sprintf("%v", file.Stat.Mode())
			if len(file.XAttributes) > 0 {
				perm += "@"
			} else {
				perm += " "
			}
			fd.SetValue(perm)
			fd.SetValueColor(cpmp)
			fd.SetColorfulValue(file.ColorPermission())
		case PFieldLinks: //"Links",
			fd.SetValue(file.NLinks())
			fd.SetValueColor(clkp)
		case PFieldSize: //"Size",
			if file.IsDir() {
				fd.SetValue("-")
				fd.SetColorfulValue(cdashp.Sprintf("%[1]*[2]v", fd.Width, "-"))
			} else {
				fd.SetValue(ByteSize(file.Size))
				csize := fdColorizedSize(file.Size, fd.Width)
				fd.SetColorfulValue(csize)
			}
			fd.SetValueColor(csnp)
		case PFieldBlocks: //"User",
			if file.IsDir() {
				fd.SetValue("-")
				fd.SetColorfulValue(cdashp.Sprintf("%[1]*[2]v", fd.Width, "-"))
			} else {
				fd.SetValue(file.Blocks())
				fd.SetColorfulValue(cbkp.Sprintf("%[1]*[2]v", fd.Width, file.Blocks()))
			}
			fd.SetValueColor(cbkp)
		case PFieldUser: //"User",
			fd.SetValue(urname)
			fd.SetValueColor(cuup)
			fd.SetColorfulValue(cuup.Sprintf("%[1]*[2]v", fd.Width, urname))
		case PFieldGroup: //"Group",
			fd.SetValue(gpname)
			fd.SetValueColor(cgup)
			fd.SetColorfulValue(cgup.Sprintf("%[1]*[2]v", fd.Width, gpname))
		case PFieldModified: //"Date Modified",
			fd.SetValue(DateString(file.ModifiedTime()))
			fd.SetValueColor(cdap)
			// fd.SetColorfulValue(file.ColorModifyTime())
		case PFieldCreated: //"Date Created",
			fd.SetValue(DateString(file.CreatedTime()))
			fd.SetValueColor(cdap)
			// fd.SetColorfulValue(file.ColorCreatedTime())
		case PFieldAccessed: //"Date Accessed",
			fd.SetValue(DateString(file.AccessedTime()))
			fd.SetValueColor(cdap)
			// fd.SetColorfulValue(file.ColorAccessedTime())
		case PFieldGit: //"Gid",
			if git.NoGit {
				continue
			} else {
				fd.SetValue(getGitStatus(git, file))
				fd.SetColorfulValue(file.ColorGitStatus(git))
			}
			fd.SetValueColor(cgtp)
		case PFieldName: //"Name",
			fd.SetValue(file.Name())
			fd.SetValueColor(cfip)
			fd.SetColorfulValue(file.ColorName())
		}
	}
}

// SetColorfulValues sets up colorful values of FieldSlice from File and GitStatus
func (f *FieldSlice) SetColorfulValues(file *File, git GitStatus) {
	for _, fd := range f.fds {
		switch fd.Key {
		case PFieldINode: //"inode",
		case PFieldPermissions: //"Permissions",
			perm := fmt.Sprintf("%v", file.Stat.Mode())
			if len(file.XAttributes) > 0 {
				perm += "@"
			} else {
				perm += " "
			}
			fd.SetColorfulValue(file.ColorPermission())
		case PFieldLinks: //"Links",
		case PFieldSize: //"Size",
			if file.IsDir() {
				fd.SetColorfulValue(cdashp.Sprintf("%[1]*[2]v", fd.Width, "-"))
			} else {
				csize := fdColorizedSize(file.Size, fd.Width)
				fd.SetColorfulValue(csize)
			}
		case PFieldBlocks: //"User",
			if file.IsDir() {
				fd.SetColorfulValue(cdashp.Sprintf("%[1]*[2]v", fd.Width, "-"))
			} else {
				fd.SetColorfulValue(cbkp.Sprintf("%[1]*[2]v", fd.Width, file.Blocks()))
			}
		case PFieldUser: //"User",
			fd.SetColorfulValue(cuup.Sprintf("%[1]*[2]v", fd.Width, urname))
		case PFieldGroup: //"Group",
			fd.SetColorfulValue(cgup.Sprintf("%[1]*[2]v", fd.Width, gpname))
		case PFieldModified: //"Date Modified",
			// fd.SetColorfulValue(file.ColorModifyTime())
		case PFieldCreated: //"Date Created",
			// fd.SetColorfulValue(file.ColorCreatedTime())
		case PFieldAccessed: //"Date Accessed",
			// fd.SetColorfulValue(file.ColorAccessedTime())
		case PFieldGit: //"Gid",
			if git.NoGit {
				continue
			} else {
				fd.SetColorfulValue(file.ColorGitStatus(git))
			}
		case PFieldName: //"Name",
			fd.SetColorfulValue(file.ColorName())
		}
	}
}

// Count will return number of fields in FieldSlice
func (f *FieldSlice) Count() int {
	return len(f.fds)
}

// Add will append a Field to FieldSlice
func (f *FieldSlice) Add(field *Field) {
	f.fds = append(f.fds, field)
}

// Remove will remove the first matched field according to PDFieldFlag
func (f *FieldSlice) Remove(key PDFieldFlag) {
	for i, fd := range f.fds {
		if fd.Key == key {
			if i == f.Count() {
				f.fds = f.fds[:i]
			} else {
				f.fds = append(f.fds[:i], f.fds[i+1:]...)
			}
		}
	}
}

// RemoveByName will remove the first matched field according to Field.Name
func (f *FieldSlice) RemoveByName(name string) {
	for i, fd := range f.fds {
		if fd.Name == name {
			if i == f.Count() {
				f.fds = f.fds[:i]
			} else {
				f.fds = append(f.fds[:i], f.fds[i+1:]...)
			}
		}
	}
}

// Insert will insert a field into the poisition of FieldSlice according to the index `startIndex`
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

// Get will return *Field for first matched key in FieldSlice
func (f *FieldSlice) Get(key PDFieldFlag) *Field {
	for _, fd := range f.fds {
		if fd.Key == key {
			return fd
		}
	}
	return nil
}

// Get will return *Field for first matched name in FieldSlice
func (f *FieldSlice) GetByName(name string) *Field {
	for _, fd := range f.fds {
		if fd.Name == name {
			return fd
		}
	}
	return nil
}

// Aligns will return the paw.Align slice from Field.Align of FieldSlie
func (f *FieldSlice) Aligns() []paw.Align {
	hds := make([]paw.Align, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		hds[i] = fd.Align
	}
	return hds
}

// Heads will return the string slice from Field.Name of FieldSlie
func (f *FieldSlice) Heads() []string {
	hds := make([]string, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		hds[i] = fd.HeadString()
	}
	return hds
}

// HeadWidths will return the int slice from Field.Width of FieldSlie
func (f *FieldSlice) HeadWidths() []int {
	hds := make([]int, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		hds[i] = fd.Width
	}
	return hds
}

// HeadsString will return string join by a space of FieldSlice.Head()
func (f *FieldSlice) HeadsString() string {
	return strings.Join(f.Heads(), " ")
	// sb := paw.NewStringBuilder()
	// for i := 0; i < f.Count(); i++ {
	// 	fd := f.fds[i]
	// 	if i < f.Count()-1 {
	// 		fmt.Fprintf(sb, "%s ", fd.HeadString())
	// 	} else {
	// 		fmt.Fprintf(sb, "%s", fd.HeadString())
	// 	}
	// }
	// return sb.String()
}

// HeadsStringWidth will return width of FieldSlice.HeadString() as you see
func (f *FieldSlice) HeadsStringWidth() int {
	hds := f.HeadsString()
	return paw.StringWidth(hds)
}

// MetaHeadsStringWidth will return width of FieldSlice.HeadString() exclude `PFieldName` as you see
func (f *FieldSlice) MetaHeadsStringWidth() int {
	wd := 0
	for _, fd := range f.fds {
		if fd.Key == PFieldName {
			continue
		}
		wd += fd.Width + 1
	}
	return wd - 1
}

// ColorHeads will return the colorful string slice from Field.Name of FieldSlie
func (f *FieldSlice) ColorHeads() []string {
	hds := make([]string, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		hds[i] = fd.ColorHeadString()
	}
	return hds
}

// ColorHeadsString will return colorful string join by a space of FieldSlice.Head()
func (f *FieldSlice) ColorHeadsString() string {
	return strings.Join(f.ColorHeads(), " ")
	// w := paw.NewStringBuilder()
	// for i, h := range f.ColorHeads() {
	// 	if i < f.Count() {
	// 		fmt.Fprintf(w, "%s ", h)
	// 	} else {
	// 		fmt.Fprintf(w, "%s", h)
	// 	}
	// }
	// return w.String()
}

// Colors will return the []*color.Color slice from Field.ValueColor of FieldSlie
func (f *FieldSlice) Colors() []*color.Color {
	vals := make([]*color.Color, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		vals[i] = fd.ValueColor
	}
	return vals
}

// Values will return all interface{} slice from Field.Value of FieldSlie
func (f *FieldSlice) Values() []interface{} {
	vals := make([]interface{}, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		vals[i] = fd.Value
	}
	return vals
}

// ValueCs will return the interface{} slice from Field.ValueC of FieldSlie
func (f *FieldSlice) ValueCs() []interface{} {
	vals := make([]interface{}, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		vals[i] = fd.ValueC
	}
	return vals
}

// ValuesStrings will return the string slice from Field.Value of FieldSlie
func (f *FieldSlice) ValuesStrings() []string {
	return cast.ToStringSlice(f.Values())
	// vals := make([]string, f.Count())
	// for i := 0; i < f.Count(); i++ {
	// 	fd := f.fds[i]
	// 	vals[i] = fd.ValueString()
	// }
	// return vals
}

// ValuesStringSlice will return the string slice according to idx from Field.Value of FieldSlie
func (f *FieldSlice) ValuesStringSlice(idxs ...int) []string {
	vs := f.ValuesStrings()
	out := make([]string, len(idxs))
	nidxs := len(idxs)
	for i := 0; i < nidxs; i++ {
		idx := idxs[i]
		if err := paw.CheckIndex(vs, idx, "f.ValuesStrings()"); err == nil {
			out[i] = vs[i]
		} else {
			out[i] = fmt.Sprint(err)
		}
	}
	return out
}

// ColorValueStrings will return the colorful string slice from Field.ColorValueString() of FieldSlie
func (f *FieldSlice) ColorValueStrings() []string {
	vals := make([]string, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		vals[i] = fd.ColorValueString()
	}
	return vals
}

// ColorValueStringSlice will return the string slice according to idx from Field.Value of FieldSlie
func (f *FieldSlice) ColorValueStringSlice(idxs ...int) []string {
	vs := f.ColorValueStrings()
	out := make([]string, len(idxs))
	nidxs := len(idxs)
	for i := 0; i < nidxs; i++ {
		idx := idxs[i]
		if err := paw.CheckIndex(vs, idx, "f.ColorValueStrings()"); err == nil {
			out[i] = vs[i]
		} else {
			out[i] = fmt.Sprint(err)
		}
	}
	return out
}

// MetaValuesString will return string slice of Field.ValueString() exclude `PFieldName`
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

// MetaValuesString will return string of FieldSlice.MetaValuesString() exclude `PFieldName` as you see
func (f *FieldSlice) MetaValuesString() string {
	return strings.Join(f.MetaValues(), " ")
	// sb := paw.NewStringBuilder()
	// for i, h := range f.MetaValues() {
	// 	if i < f.Count()-1 {
	// 		fmt.Fprintf(sb, "%s ", h)
	// 	} else {
	// 		fmt.Fprintf(sb, "%s", h)
	// 	}
	// }
	// return sb.String()
}

// MetaValuesStringWidth will return width of FieldSlice.MetaValuesString() exclude `PFieldName` as you see
func (f *FieldSlice) MetaValuesStringWidth() int {
	s := f.MetaValuesString()
	return paw.StringWidth(s)
}

// ColorMetaValues will return string slice of Field.ColorValueString() exclude `PFieldName`
func (f *FieldSlice) ColorMetaValues() []string {
	var vals []string
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		switch fd.Key {
		case PFieldName:
			continue
		default:
			value := fd.ColorValueString()
			// if fd.ValueC == nil {
			// 	switch fd.Align {
			// 	case paw.AlignLeft:
			// 		value = fmt.Sprintf("%-[1]*[2]v", fd.Width, fd.Value)
			// 	default:
			// 		value = fmt.Sprintf("%[1]*[2]v", fd.Width, fd.Value)
			// 	}
			// 	value = fd.ValueColor.Sprint(value)
			// }
			vals = append(vals, value)
		}
	}
	return vals
}

// ColorHeadsString will return colorful string join by a space of FieldSlice.ColorMetaValues() exclude `PFieldName`
func (f *FieldSlice) ColorMetaValuesString() string {
	// v := strings.Join(f.ColorMetaValues(), " ")
	// fmt.Printf("%d, %q\n", paw.StringWidth(v), v)
	return strings.Join(f.ColorMetaValues(), " ")
	// sb := paw.NewStringBuilder()
	// for i, h := range f.ColorMetaValues() {
	// 	if i < f.Count() {
	// 		fmt.Fprintf(sb, "%s ", h)
	// 	} else {
	// 		fmt.Fprintf(sb, "%s", h)
	// 	}
	// }
	// return sb.String()
}
