package filetree

import (
	"fmt"
	"io"
	"strings"

	"github.com/thoas/go-funk"

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
	if keys == nil || len(keys) == 0 {
		keys = pfieldKeysDefualt
	}
	for _, k := range keys {
		if k == PFieldGit && git.NoGit {
			continue
		}
		f.Add(k)
	}
	return f
}

func fdColorizedSize(size uint64, width int) (csize string) {
	ss := ByteSize(size)
	nss := len(ss)
	sn := fmt.Sprintf("%[1]*[2]s", width-1, ss[:nss-1])
	su := paw.ToLower(ss[nss-1:])
	cn := paw.NewEXAColor("sn")
	cu := paw.NewEXAColor("sb")
	csize = cn.Sprint(sn) + cu.Sprint(su)
	return csize
}

// Fields will return Fields of FieldSlice
func (f *FieldSlice) Fields() []*Field {
	return f.fds
}

// EmptyValues will empty all Field.Value[C] of FieldSlice with nil
func (f *FieldSlice) EmptyValues() {
	for _, fd := range f.fds {
		fd.Value = ""
		fd.ValueC = ""
	}
}

// SetValues sets up values of FieldSlice from File and GitStatus
func (f *FieldSlice) SetValues(file *File, git GitStatus) {
	for _, fd := range f.fds {
		switch fd.Key {
		case PFieldINode: //"inode",
			fd.SetValue(file.INode())
			// fd.SetValueColor(cinp)
		case PFieldPermissions: //"Permissions",
			fd.SetValue(file.Permission())
			fd.SetColorfulValue(file.ColorPermission())
			// fd.SetValueColor(cpmp)
		case PFieldLinks: //"Links",
			fd.SetValue(file.NLinks())
			// fd.SetValueColor(clkp)
		case PFieldSize: //"Size",
			if file.IsDir() {
				fd.SetValue("-")
				fd.SetColorfulValue(cdashp.Sprintf("%[1]*[2]v", fd.Width, "-"))
			} else {
				fd.SetValue(file.ByteSize())
				csize := fdColorizedSize(file.Size, fd.Width)
				fd.SetColorfulValue(csize)
			}
			// fd.SetValueColor(csnp)
		case PFieldBlocks: //"User",
			if file.IsDir() {
				fd.SetValue("-")
				fd.SetColorfulValue(cdashp.Sprintf("%[1]*[2]v", fd.Width, "-"))
			} else {
				fd.SetValue(file.Blocks())
				fd.SetColorfulValue(cbkp.Sprintf("%[1]*[2]v", fd.Width, file.Blocks()))
			}
			// fd.SetValueColor(cbkp)
		case PFieldUser: //"User",
			fd.SetValue(urname)
			fd.SetColorfulValue(cuup.Sprintf("%[1]*[2]v", fd.Width, urname))
			// fd.SetValueColor(cuup)
		case PFieldGroup: //"Group",
			fd.SetValue(gpname)
			fd.SetColorfulValue(cgup.Sprintf("%[1]*[2]v", fd.Width, gpname))
			// fd.SetValueColor(cgup)
		case PFieldModified: //"Date Modified",
			fd.SetValue(DateString(file.ModifiedTime()))
			// fd.SetColorfulValue(file.ColorModifyTime())
			// fd.SetValueColor(cdap)
		case PFieldCreated: //"Date Created",
			fd.SetValue(DateString(file.CreatedTime()))
			// fd.SetColorfulValue(file.ColorCreatedTime())
			// fd.SetValueColor(cdap)
		case PFieldAccessed: //"Date Accessed",
			fd.SetValue(DateString(file.AccessedTime()))
			// fd.SetColorfulValue(file.ColorAccessedTime())
			// fd.SetValueColor(cdap)
		case PFieldGit: //"Gid",
			if git.NoGit {
				continue
			} else {
				fd.SetValue(file.GitStatus(git))
				fd.SetColorfulValue(file.ColorGitStatus(git))
			}
			// fd.SetValueColor(cgitp)
		case PFieldName: //"Name",
			fd.SetValue(file.Name())
			fd.SetColorfulValue(file.ColorName())
			fd.SetValueColor(file.LSColor())
			// fd.SetValueColor(cfip)
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

// Add will append a Field to FieldSlice by Field ke
func (f *FieldSlice) Add(key PDFieldFlag) *FieldSlice {
	field := FieldsMap[key]
	if _, ok := pfieldAlignMap[key]; !ok {
		field.Align = paw.AlignLeft
	}
	f.AddByField(field)
	return f
}

// AddByField will append a Field to FieldSlice
func (f *FieldSlice) AddByField(field *Field) *FieldSlice {
	f.fds = append(f.fds, field)
	return f
}

// Remove will remove the first matched field according to PDFieldFlag
func (f *FieldSlice) Remove(key PDFieldFlag) *FieldSlice {
	for i, fd := range f.fds {
		if fd.Key == key {
			if i == f.Count()-1 {
				f.fds = f.fds[:i]
			} else {
				f.fds = append(f.fds[:i], f.fds[i+1:]...)
			}
		}
	}
	return f
}

// RemoveByName will remove the first matched field according to Field.Name
func (f *FieldSlice) RemoveByName(name string) *FieldSlice {
	for i, fd := range f.fds {
		if fd.Name == name {
			if i == f.Count() {
				f.fds = f.fds[:i]
			} else {
				f.fds = append(f.fds[:i], f.fds[i+1:]...)
			}
		}
	}
	return f
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

// Get will return *Field for first matched name (case-insensetive) in FieldSlice
func (f *FieldSlice) GetByName(name string) *Field {
	for _, fd := range f.fds {
		if paw.ToLower(fd.Name) == paw.ToLower(name) {
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

// ModifyWidth modifies Field.Width according to FileList and wdstty (maximum width on console).
func (f *FieldSlice) ModifyWidth(fl *FileList, wdstty int) {
	var (
		wdinode = 0
		wdlinks = 0
		wdsize  = 0
		wdblock = 0
	)

	for _, dir := range fl.Dirs() {
		for _, file := range fl.Map()[dir][1:] {
			ws := len(fmt.Sprint(file.INode()))
			if wdinode < ws {
				wdinode = ws
			}
			ws = len(fmt.Sprint(file.NLinks()))
			if wdlinks < ws {
				wdlinks = ws
			}
			ws = len(file.ByteSize())
			if wdsize < ws {
				wdsize = ws
			}
			ws = len(fmt.Sprint(file.Blocks()))
			if wdblock < ws {
				wdblock = ws
			}
		}
	}

	if fd := f.Get(PFieldINode); fd != nil {
		fd.Width = paw.MaxInt(wdinode, fd.Width)
	}
	if fd := f.Get(PFieldLinks); fd != nil {
		fd.Width = paw.MaxInt(wdlinks, fd.Width)
	}
	if fd := f.Get(PFieldSize); fd != nil {
		fd.Width = paw.MaxInt(wdsize, fd.Width)
	}
	if fd := f.Get(PFieldBlocks); fd != nil {
		fd.Width = paw.MaxInt(wdblock, fd.Width)
	}
	if fd := f.Get(PFieldName); fd != nil {
		wdmeta := f.MetaHeadsStringWidth() + 1
		fd.Width = wdstty - wdmeta
	}
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

// Colors will return a copy of the []*color.Color slice from Field.ValueColor of FieldSlie
func (f *FieldSlice) Colors() []*color.Color {
	vals := make([]*color.Color, f.Count())
	for i, fd := range f.fds {
		vals[i] = fd.ValueColor
	}
	return vals
}

// Values will return a copy of all interface{} slice from Field.Value of FieldSlie
func (f *FieldSlice) Values() []interface{} {
	vals := make([]interface{}, f.Count())
	for i, fd := range f.fds {
		vals[i] = fd.Value
	}
	return vals
}

// ValueCs will return a copy of the interface{} slice from Field.ValueC of FieldSlie
func (f *FieldSlice) ValueCs() []interface{} {
	vals := make([]interface{}, f.Count())
	for i, fd := range f.fds {
		vals[i] = fd.ValueC
	}
	return vals
}

// ValuesStrings will return a copy of the string slice from Field.Value of FieldSlie
func (f *FieldSlice) ValuesStrings() []string {
	vals := make([]string, f.Count())
	for i, fd := range f.fds {
		vals[i] = fd.ValueString()
	}
	return vals
}

// ValuesStringSlice will return a copy of the string slice according to idx from Field.Value of FieldSlie
func (f *FieldSlice) ValuesStringSlice(idxs ...int) []string {
	vs := f.ValuesStrings()
	out := make([]string, len(idxs))
	for _, i := range idxs {
		if err := paw.CheckIndex(vs, i, "f.ValuesStrings()"); err == nil {
			out[i] = vs[i]
		} else {
			out[i] = fmt.Sprint(err)
		}
	}
	return out
}

// ColorValueStrings will return a copy of the colorful string slice from Field.ColorValueString() of FieldSlie
func (f *FieldSlice) ColorValueStrings() []string {
	vals := make([]string, f.Count())
	for i, fd := range f.fds {
		vals[i] = fd.ColorValueString()
	}
	return vals
}

// ColorValueStringSlice will return a copy of the string slice according to idx from Field.Value of FieldSlie
func (f *FieldSlice) ColorValueStringSlice(idxs ...int) []string {
	vs := f.ColorValueStrings()
	out := make([]string, len(idxs))
	for _, i := range idxs {
		if err := paw.CheckIndex(vs, i, "f.ColorValueStrings()"); err == nil {
			out[i] = vs[i]
		} else {
			out[i] = fmt.Sprint(err)
		}
	}
	return out
}

// MetaValuesString will return a copy of string slice of Field.ValueString() exclude `PFieldName`
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

// PrintHeadRow prints out all Filed.Name to w
func (f *FieldSlice) PrintHeadRow(w io.Writer, pad string) {
	// print head
	fmt.Fprintln(w, pad+f.ColorHeadsString())
}

// PrintRow prints out all value of Field to w
func (f *FieldSlice) PrintRow(w io.Writer, pad string) {
	f.PrintRowPrefix(w, pad, "")
	// // print meta
	// fmt.Fprint(w, pad+f.ColorMetaValuesString()+" ")
	// // print Name field
	// var (
	// 	wpad   = f.MetaHeadsStringWidth()
	// 	fdName = f.Get(PFieldName)
	// 	width  = fdName.Width
	// 	value  = cast.ToString(fdName.Value)
	// 	wv     = paw.StringWidth(value)
	// 	cvalue = cast.ToString(fdName.ValueC)
	// 	c      = fdName.ValueColor
	// )
	// if wv <= width {
	// 	if fdName.ValueC != nil {
	// 		fmt.Fprintln(w, cvalue)
	// 	} else {
	// 		fmt.Fprintln(w, c.Sprint(value))
	// 	}
	// } else {
	// 	names := paw.WrapToSlice(value, width)
	// 	fmt.Fprintln(w, c.Sprint(names[0]))
	// 	sp := pad + paw.Spaces(wpad)
	// 	for i := 1; i < len(names); i++ {
	// 		fmt.Fprintln(w, sp, c.Sprint(names[i]))
	// 	}
	// }
}

// PrintRowPrefix prints out all value of Field to w
func (f *FieldSlice) PrintRowPrefix(w io.Writer, pad, prefix string) {
	// print meta
	fmt.Fprint(w, pad+f.ColorMetaValuesString()+" ")
	// print Name field
	var (
		wpad   = f.MetaHeadsStringWidth()
		wprf   = paw.StringWidth(paw.StripANSI(prefix))
		fdName = f.Get(PFieldName)
		width  = fdName.Width
		value  = cast.ToString(fdName.Value)
		wv     = paw.StringWidth(value)
		cvalue = cast.ToString(fdName.ValueC)
		c      = fdName.ValueColor
	)
	if wprf > 0 {
		prefix += " "
		width -= wprf - 1
	}
	if wv <= width {
		if fdName.ValueC != nil {
			fmt.Fprintln(w, prefix+cvalue)
		} else {
			fmt.Fprintln(w, prefix+c.Sprint(value))
		}
	} else {
		names := paw.WrapToSlice(value, width)
		fmt.Fprintln(w, prefix+c.Sprint(names[0]))
		sp := pad + paw.Spaces(wpad)
		for i := 1; i < len(names); i++ {
			fmt.Fprintln(w, sp, c.Sprint(names[i]))
		}
	}
}

// PrintRowButIgnoreln prints out all value of Field to w excluding ignoreIdx,ended with '\n'
func (f *FieldSlice) PrintRowButIgnoreln(w io.Writer, pad string, ignoreIdx ...int) {
	f.PrintRowButIgnore(w, pad, ignoreIdx...)
	fmt.Fprint(w, '\n')
}

// PrintRowButIgnore prints out all value of Field to w excluding ignoreIdx
func (f *FieldSlice) PrintRowButIgnore(w io.Writer, pad string, ignoreIdx ...int) {
	for i, fd := range f.Fields() {
		if funk.ContainsInt(ignoreIdx, i) {
			continue
		}
		var (
			width  = fd.Width
			value  = cast.ToString(fd.Value)
			wv     = paw.StringWidth(value)
			cvalue = cast.ToString(fd.ValueC)
			c      = fd.ValueColor
		)
		if i == 0 {
			fmt.Fprint(w, pad)
		}
		if wv <= width {
			if fd.ValueC != nil {
				fmt.Fprintln(w, cvalue)
			} else {
				fmt.Fprintln(w, c.Sprint(value))
			}
		} else {
			names := paw.WrapToSlice(value, width)
			fmt.Fprintln(w, c.Sprint(names[0]))
			for i := 1; i < len(names); i++ {
				fmt.Fprintln(w, c.Sprint(names[i]))
			}
		}
		fmt.Fprint(w, " ")
	}
}

// PrintRow prints out all value of Field to w
func (f *FieldSlice) PrintRowXattr(w io.Writer, pad string, xattrs []string, xsymb string) {

	if len(xattrs) == 0 {
		return
	}
	if len(xsymb) == 0 {
		xsymb = XattrSymbol
	}
	var (
		wmeta  = f.MetaHeadsStringWidth()
		spmeta = pad + paw.Spaces(wmeta)
		fdName = f.Get(PFieldName)
		csymb  = cxbp.Sprint(xsymb)
		wsymb  = paw.StringWidth(xsymb)
		cbsp   = cxbp.Sprint(paw.Spaces(wsymb))
		width  = fdName.Width - wsymb
	)
	for _, value := range xattrs {
		wv := paw.StringWidth(value)
		if wv <= width {
			fmt.Fprintln(w, spmeta, csymb+cxap.Sprint(value))
		} else {
			names := paw.WrapToSlice(value, width)
			fmt.Fprintln(w, spmeta, csymb+cxap.Sprint(names[0]))
			for i := 1; i < len(names); i++ {
				fmt.Fprintln(w, spmeta, cbsp+cxap.Sprint(names[i]))
			}
		}
	}
}
