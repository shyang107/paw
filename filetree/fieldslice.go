package filetree

import (
	"fmt"
	"io"
	"path/filepath"
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
func NewFieldSliceFrom(flags []PDFieldFlag, git *GitStatus) (fds *FieldSlice) {
	f := NewFieldSlice()
	if len(flags) == 0 {
		flags = DefaultPDFieldKeys
	}
	f.fds = NewFieldsGit(git.NoGit, flags...)
	return f
}

func fdColorizedSize(size uint64, width int) (csize string) {
	ss := ByteSize(size)
	nss := len(ss)
	csn := csnp.Sprintf("%[1]*[2]s", width-1, ss[:nss-1])
	su := strings.ToLower(ss[nss-1:])
	csize = csn + csup.Sprint(su)
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
func (f *FieldSlice) SetValues(file *File, git *GitStatus) {
	for _, fd := range f.fds {
		cp := fd.Key.Color()
		switch fd.Key {
		case PFieldINode: //"inode",
			fd.SetValue(file.INode())
			fd.SetValueC(calign(cinp, fd.Align, fd.Width, file.INode()))
			fd.SetValueColor(cp)
			// fd.SetValueColor(cinp)
		case PFieldPermissions: //"Permissions",
			perm := file.Permission()
			fd.SetValue(perm)
			wp := len(perm)
			sp := ""
			if wp < fd.Width {
				sp = paw.Spaces(fd.Width - wp)
			}
			fd.SetValueC(file.PermissionC() + sp)
			fd.SetValueColor(cp)
			// fd.SetValueColor(cpms)
		case PFieldLinks: //"Links",
			fd.SetValue(file.NLinks())
			fd.SetValueC(calign(clkp, fd.Align, fd.Width, file.NLinks()))
			fd.SetValueColor(cp)
			// fd.SetValueColor(clkp)
		case PFieldSize: //"Size",
			kind := nodeTypeFromFileInfo(file.Info)
			switch kind {
			case kindChardev, kindDev:
				major, minor := file.DevNumber()
				csj := csnp.Sprintf("%[1]*[2]v", fd.widthMajor, major)
				csn := csnp.Sprintf("%[1]*[2]v", fd.widthMinor, minor)
				cdev := csj + cdirp.Sprint(",") + csn
				wdev := fd.widthMajor + fd.widthMinor + 1 //len(paw.StripANSI(cdev))
				if wdev < fd.Width {
					cdev = csj + cdirp.Sprint(",") + paw.Spaces(fd.Width-wdev) + csn
				}
				fd.SetValue(file.DevNumberString())
				fd.SetValueC(cdev)
			case kindDir:
				fd.SetValue("-")
				fd.SetValueC(calign(cdashp, fd.Align, fd.Width, "-"))
			default:
				fd.SetValue(file.ByteSize())
				csize := fdColorizedSize(file.Size, fd.Width)
				fd.SetValueC(csize)
			}
			fd.SetValueColor(cp)
			// fd.SetValueColor(csnp)
		case PFieldBlocks: //"Block",
			if file.IsDir() {
				fd.SetValue("-")
				fd.SetValueC(calign(cdashp, fd.Align, fd.Width, "-"))
			} else {
				fd.SetValue(file.Blocks())
				fd.SetValueC(calign(cbkp, fd.Align, fd.Width, file.Blocks()))
				fd.SetValueColor(cp)
				// fd.SetValueColor(cbkp)
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
			fd.SetValueColor(cp)
			// fd.SetValueColor(cdap)
		case PFieldCreated: //"Date Created",
			sd := DateString(file.CreatedTime())
			fd.SetValue(sd)
			fd.SetValueC(calign(cdap, fd.Align, fd.Width, sd))
			fd.SetValueColor(cp)
			// fd.SetValueColor(cdap)
		case PFieldAccessed: //"Date Accessed",
			sd := DateString(file.AccessedTime())
			fd.SetValue(sd)
			fd.SetValueC(calign(cdap, fd.Align, fd.Width, sd))
			fd.SetValueColor(cp)
			// fd.SetValueColor(cdap)
		case PFieldMd5: //"Date Accessed",
			md5 := file.GetMd5()
			fd.SetValue(md5)
			if md5 == "-" {
				fd.SetValueC(calign(cdashp, fd.Align, fd.Width, md5))
				fd.SetValueColor(cdashp)
			} else {
				fd.SetValueC(calign(cp, fd.Align, fd.Width, md5))
				fd.SetValueColor(cp)
			}
			// fd.SetValueC(calign(cmd5p, fd.Align, fd.Width, md5))
			// fd.SetValueColor(cmd5p)
		case PFieldGit: //"Gid",
			if git.NoGit {
				continue
			} else {
				var xy, xyc string
				xy = file.GitXYs(git)
				xyc = file.GitXYc(git)
				// var relpath fmt.Stringer
				// relpath := file.RelPath
				// xy = git.XYStatus(relpath)
				// xyc = git.XYStatusC(relpath)
				// if xy != GitUnChanged.String()+GitUnChanged.String() {
				// 	paw.Logger.WithFields(logrus.Fields{
				// 		"rp": relpath,
				// 		// "oXY": ,
				// 		"XY":  xy + "," + xyc,
				// 		"fXY": file.GitXY(git) + "," + file.GitXYC(git),
				// 	}).Debug(file.BaseNameC())
				// }
				fd.SetValue(" " + xy)
				fd.SetValueC(" " + xyc)
				// fd.SetValue(" " + file.GitXY(git))
				// fd.SetValueC(" " + file.GitXYC(git))
				fd.SetValueColor(cp)
				// fd.SetValueColor(cgitp)
			}
		case PFieldName: //"Name",
			fd.SetValue(file.Name())
			fd.SetValueC(file.NameC())
			fd.SetValueColor(file.LSColor())
			fd.SetIsLink(file.IsLink())
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
	field := key.Field()
	f.AddByField(field)
	return f
}

// AddByField will append a Field to FieldSlice
func (f *FieldSlice) AddByField(field *Field) *FieldSlice {
	if field != nil {
		f.fds = append(f.fds, field)
	}
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
			if i == f.Count()-1 {
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
	switch {
	case startIndex >= len(f.fds): // append to tail
		f.fds = append(f.fds, fds...)
	case startIndex <= 0: // insert to head
		f.fds = append(fds, f.fds...)
	default: // insert to middle
		rear := append([]*Field{}, f.fds[startIndex:]...)
		f.fds = append(f.fds[0:startIndex], fds...)
		f.fds = append(f.fds, rear...)
	}
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
		if strings.EqualFold(fd.Name, name) {
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
	paw.Logger.Trace()
	for _, dir := range fl.Dirs() {
		for _, file := range fl.Map()[dir][:] {
			for _, field := range pdOpt.FieldKeys() {
				var fd = f.Get(field)
				if fd == nil {
					continue
				}
				switch field {
				case PFieldSize:
					w, wj, wn := file.widthOfSize()
					fd.widthMajor = paw.MaxInt(fd.widthMajor, wj)
					fd.widthMinor = paw.MaxInt(fd.widthMinor, wn)
					wd := fd.widthMajor + fd.widthMinor + 1
					fd.Width = paw.MaxInts(fd.Width, w, wd)
				default:
					fd.Width = paw.MaxInt(fd.Width, file.WidthOf(field))
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
		isLink = fdName.isLink
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
		if isLink {
			values := strings.Split(value, " -> ")
			name, link := values[0], values[1]
			cname := c.Sprint(name)
			wbname := paw.StringWidth(name)
			carrow := cdashp.Sprint(" -> ")
			wbname += 4
			flink, _ := NewFile(link)
			clk := GetFileLSColor(flink)
			// L1.0
			fmt.Fprint(w, prefix+cname+carrow)
			dir, name := filepath.Split(link)
			wd, wn := paw.StringWidth(dir), paw.StringWidth(name)
			sp := pad + paw.Spaces(wpad)
			if wd+wn <= width-wbname {
				// L1.1End
				fmt.Fprintln(w, cdirp.Sprint(dir)+clk.Sprint(name))
			} else {
				if wd <= width-wbname {
					clink := cdirp.Sprint(dir) + clk.Sprint(name[:width-wbname-wd])
					// L1.1End
					fmt.Fprintln(w, clink)
					names := paw.WrapToSlice(name[width-wbname-wd:], width)
					for _, v := range names {
						clink = clk.Sprint(v)
						// L2...
						fmt.Fprintln(w, sp, clink)
					}
				} else { // wd > width-wbname
					// L1.1End
					var clink string
					wd1End := width - wbname
					clink = cdirp.Sprint(dir[:wd1End])
					fmt.Fprintln(w, clink)
					wdLast := wd - wd1End
					if wdLast <= width {
						// L2.0
						clink = cdirp.Sprint(dir[wd1End:])
						// fmt.Fprintln(w, sp, clink)
						wn2End := width - wdLast
						// L2
						if wn2End <= width {
							cname := clk.Sprint(name)
							fmt.Fprintln(w, sp, clink+cname)
						} else {
							cname := clk.Sprint(name[:wn2End])
							fmt.Fprintln(w, sp, clink+cname)
							names := paw.WrapToSlice(name[wn2End:], width)
							nn := len(names)
							for i := 0; i < nn; i++ {
								cname = clk.Sprint(name[i])
								fmt.Fprintln(w, cname)
							}
						}
					} else {
						// L2.0
						dirs := paw.WrapToSlice(dir, width)
						nd := len(dirs)
						for i := 0; i < nd-1; i++ {
							clink = cdirp.Sprint(dirs[i])
							fmt.Fprintln(w, sp, clink)
						}
						wnLast := width - paw.StringWidth(dirs[nd-1])
						if wn <= wnLast {
							cname := clk.Sprint(name)
							fmt.Fprintln(w, cname)
						} else {
							cname := clk.Sprint(name[:wnLast])
							fmt.Fprintln(w, cname)
							names := paw.WrapToSlice(name[wnLast:], width)
							nn := len(names)
							for i := 0; i < nn; i++ {
								cname = clk.Sprint(name[i])
								fmt.Fprintln(w, cname)
							}
						}
					}
				}
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
