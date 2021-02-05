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
	su := strings.ToLower(ss[nss-1:])
	csize = csnp.Sprint(sn) + csup.Sprint(su)
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
			fd.SetValueC(calign(cinp, fd.Align, fd.Width, file.INode()))
			fd.SetValueColor(cinp)
		case PFieldPermissions: //"Permissions",
			perm := file.Permission()
			fd.SetValue(perm)
			wp := len(perm)
			sp := ""
			if wp < fd.Width {
				sp = paw.Spaces(fd.Width - wp)
			}
			fd.SetValueC(file.PermissionC() + sp)
			fd.SetValueColor(cpms)
		case PFieldLinks: //"Links",
			fd.SetValue(file.NLinks())
			fd.SetValueC(calign(clkp, fd.Align, fd.Width, file.NLinks()))
			fd.SetValueColor(clkp)
		case PFieldSize: //"Size",
			sperm := file.Permission()
			c := string(sperm[0])
			switch c {
			case "c", "b": //file.IsChardev() || file.IsDev()
				major, minor := file.DevNumber()
				csj := csnp.Sprintf("%[1]*[2]v", fd.widthMajor, major)
				csn := csnp.Sprintf("%[1]*[2]v", fd.widthMinor, minor)
				cdev := csj + cdirp.Sprint(",") + csn
				fd.SetValue(file.DevNumberString())
				fd.SetValueC(cdev)
			case "d": //file.IsDir()
				fd.SetValue("-")
				fd.SetValueC(calign(cdashp, fd.Align, fd.Width, "-"))
			default:
				fd.SetValue(file.ByteSize())
				csize := fdColorizedSize(file.Size, fd.Width)
				fd.SetValueC(csize)
			}
			fd.SetValueColor(csnp)
		case PFieldBlocks: //"Block",
			if file.IsDir() {
				fd.SetValue("-")
				fd.SetValueC(calign(cdashp, fd.Align, fd.Width, "-"))
			} else {
				fd.SetValue(file.Blocks())
				fd.SetValueC(calign(cbkp, fd.Align, fd.Width, file.Blocks()))
				fd.SetValueColor(cbkp)
			}
		case PFieldUser: //"User",
			furname := file.User()
			fd.SetValue(furname)
			var c *color.Color
			if furname != urname {
				c = cunp
			} else {
				c = cuup
			}
			fd.SetValueC(calign(c, fd.Align, fd.Width, furname))
			fd.SetValueColor(c)
		case PFieldGroup: //"Group",
			fgpname := file.Group()
			fd.SetValue(fgpname)
			var c *color.Color
			if fgpname != gpname {
				c = cgnp
			} else {
				c = cgup
			}
			fd.SetValueC(calign(c, fd.Align, fd.Width, fgpname))
			fd.SetValueColor(c)
		case PFieldModified: //"Date Modified",
			sd := DateString(file.ModifiedTime())
			fd.SetValue(sd)
			fd.SetValueC(calign(cdap, fd.Align, fd.Width, sd))
			fd.SetValueColor(cdap)
		case PFieldCreated: //"Date Created",
			sd := DateString(file.CreatedTime())
			fd.SetValue(sd)
			fd.SetValueC(calign(cdap, fd.Align, fd.Width, sd))
			fd.SetValueColor(cdap)
		case PFieldAccessed: //"Date Accessed",
			sd := DateString(file.AccessedTime())
			fd.SetValue(sd)
			fd.SetValueC(calign(cdap, fd.Align, fd.Width, sd))
			fd.SetValueColor(cdap)
		case PFieldGit: //"Gid",
			if git.NoGit {
				continue
			} else {
				fd.SetValue(file.GitStatus(git))
				fd.SetValueC(file.GitStatusC(git))
				fd.SetValueColor(cgitp)
			}
		case PFieldName: //"Name",
			fd.SetValue(file.Name())
			fd.SetValueC(file.NameC())
			fd.SetValueColor(file.LSColor())
		}
	}
}

func calign(c *color.Color, al paw.Align, width int, value interface{}) string {
	var s string
	switch al {
	case paw.AlignLeft:
		s = c.Sprintf("%-[1]*[2]v", width, value)
	default:
		s = c.Sprintf("%[1]*[2]v", width, value)
	}
	return s
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
		if strings.ToLower(fd.Name) == strings.ToLower(name) {
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
	for _, dir := range fl.Dirs() {
		for _, file := range fl.Map()[dir][:] {
			for _, field := range pfieldKeys {
				var fd = f.Get(field)
				fd.Width = paw.MaxInt(fd.Width, file.WidthOf(field))
				switch field {
				case PFieldSize:
					_, wj, wn := file.widthOfSize()
					fd.widthMajor = paw.MaxInt(fd.widthMajor, wj)
					fd.widthMinor = paw.MaxInt(fd.widthMinor, wn)
					fd.Width = paw.MaxInt(fd.Width, fd.widthMajor+fd.widthMinor+1)

				}
			}
		}
	}
	if fd := f.Get(PFieldName); fd != nil {
		wdmeta := f.MetaHeadsStringWidth() + 1
		fd.Width = wdstty - wdmeta
	}
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
func (f *FieldSlice) Heads(isColor bool) []string {
	hds := make([]string, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		if isColor {
			hds[i] = fd.HeadStringC()
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

// HeadsC will return the colorful string slice from Field.Name of FieldSlie
func (f *FieldSlice) HeadsC() []string {
	hds := make([]string, f.Count())
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		hds[i] = fd.HeadStringC()
	}
	return hds
}

// HeadsStringC will return colorful string join by a space of FieldSlice.Head()
func (f *FieldSlice) HeadsStringC() string {
	return strings.Join(f.HeadsC(), " ")
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

// ValueStringCs will return a copy of the colorful string slice from Field.ValueStringC() of FieldSlie
func (f *FieldSlice) ValueStringCs() []string {
	vals := make([]string, f.Count())
	for i, fd := range f.fds {
		vals[i] = fd.ValueStringC()
	}
	return vals
}

// ValueStringCSlice will return a copy of the string slice according to idx from Field.Value of FieldSlie
func (f *FieldSlice) ValueStringCSlice(idxs ...int) []string {
	vs := f.ValueStringCs()
	out := make([]string, len(idxs))
	for _, i := range idxs {
		if err := paw.CheckIndex(vs, i, "f.ValueStringCs()"); err == nil {
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

// MetaValuesC will return string slice of Field.ValueStringC() exclude `PFieldName`
func (f *FieldSlice) MetaValuesC() []string {
	var vals []string
	for i := 0; i < f.Count(); i++ {
		fd := f.fds[i]
		switch fd.Key {
		case PFieldName:
			continue
		default:
			vals = append(vals, fd.ValueStringC())
		}
	}
	return vals
}

// MetaValuesStringC will return colorful string join by a space of FieldSlice.MetaValuesC() exclude `PFieldName`
func (f *FieldSlice) MetaValuesStringC() string {
	return strings.Join(f.MetaValuesC(), " ")
}

// PrintHeadRow prints out all Filed.Name to w
func (f *FieldSlice) PrintHeadRow(w io.Writer, pad string) {
	// print head
	fmt.Fprintln(w, pad+f.HeadsStringC())
}

// PrintRow prints out all value of Field to w
func (f *FieldSlice) PrintRow(w io.Writer, pad string) {
	f.PrintRowPrefix(w, pad, "")
}

// PrintRowPrefix prints out all value of Field to w
func (f *FieldSlice) PrintRowPrefix(w io.Writer, pad, prefix string) {
	// print meta
	fmt.Fprint(w, pad+f.MetaValuesStringC()+" ")
	// print Name field
	var (
		wpad   = f.MetaHeadsStringWidth()
		wprf   = paw.StringWidth(paw.StripANSI(prefix))
		fdName = f.Get(PFieldName)
		width  = fdName.Width
		value  = fmt.Sprint(fdName.Value)
		wv     = paw.StringWidth(value)
		cvalue = fmt.Sprint(fdName.ValueC)
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
