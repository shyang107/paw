package filetree

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/shyang107/paw"
	"github.com/spf13/cast"
)

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
			date := DateString(file.ModifiedTime())
			cdate := cdap.Sprintf("%-[1]*[2]s", fd.Width, date)
			fd.SetColorfulValue(cdate)
		case PFieldCreated: //"Date Created",
			date := DateString(file.CreatedTime())
			cdate := cdap.Sprintf("%-[1]*[2]s", fd.Width, date)
			fd.SetColorfulValue(cdate)
		case PFieldAccessed: //"Date Accessed",
			date := DateString(file.AccessedTime())
			cdate := cdap.Sprintf("%-[1]*[2]s", fd.Width, date)
			fd.SetColorfulValue(cdate)
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
			if i == f.Count()-1 {
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

// Widths will return Field.Width slice of FieldSlice
func (f *FieldSlice) Widths() []int {
	widths := make([]int, f.Count())
	for i, fd := range f.fds {
		widths[i] = fd.Width
	}
	return widths
}

// // HeadWidths will return the int slice from Field.Width of FieldSlie
// func (f *FieldSlice) HeadWidths() []int {
// 	hds := make([]int, f.Count())
// 	for i := 0; i < f.Count(); i++ {
// 		fd := f.fds[i]
// 		hds[i] = fd.Width
// 	}
// 	return hds
// }

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
func (f *FieldSlice) Heads(isColor bool) []string {
	hds := make([]string, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		if isColor {
			hds[i] = fd.ColorHeadString()
		} else {
			hds[i] = fd.HeadString()
		}
	}
	return hds
}

// HeadsInterface will return the string slice from Field.Name of FieldSlie
func (f *FieldSlice) HeadsInterface(isColor bool) []interface{} {
	hds := make([]interface{}, f.Count())
	for i, v := range f.Heads(isColor) {
		hds[i] = v
	}
	return hds
}

// HeadsString will return string join by a space of FieldSlice.Head()
func (f *FieldSlice) HeadsString() string {
	return strings.Join(f.Heads(false), " ")
}

// HeadsStringWidth will return width of FieldSlice.HeadString() as you see
func (f *FieldSlice) HeadsStringWidth() int {
	hds := f.HeadsString()
	return paw.StringWidth(hds)
}

// MetaHeadsStringWidth will return width of FieldSlice.HeadString() exclude `PFieldName` as you see
func (f *FieldSlice) MetaHeadsStringWidth() int {
	return f.HeadsStringWidth() - f.Get(PFieldName).Width - 1
	// wd := 0
	// for _, fd := range f.fds {
	// 	if fd.Key == PFieldName {
	// 		continue
	// 	}
	// 	wd += fd.Width + 1
	// }
	// return wd - 1
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
			vals = append(vals, fd.ColorValueString())
		}
	}
	return vals
}

// ColorHeadsString will return colorful string join by a space of FieldSlice.ColorMetaValues() exclude `PFieldName`
func (f *FieldSlice) ColorMetaValuesString() string {
	return strings.Join(f.ColorMetaValues(), " ")
}
