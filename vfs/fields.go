package vfs

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/shyang107/paw"
	"github.com/shyang107/paw/cast"
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
	ViewFieldMajor
	ViewFieldMinor
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

	ViewFieldPSUGMN = DefaultViewField

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
		ViewFieldMajor:       "Major",
		ViewFieldMinor:       "Minor",
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
		ViewFieldNo:          len(ViewFieldNames[ViewFieldNo]),
		ViewFieldINode:       len(ViewFieldNames[ViewFieldINode]),
		ViewFieldPermissions: len(ViewFieldNames[ViewFieldPermissions]),
		ViewFieldLinks:       len(ViewFieldNames[ViewFieldLinks]),
		ViewFieldSize:        len(ViewFieldNames[ViewFieldSize]),
		ViewFieldMajor:       0,
		ViewFieldMinor:       0,
		ViewFieldBlocks:      len(ViewFieldNames[ViewFieldBlocks]),
		ViewFieldUser:        len(ViewFieldNames[ViewFieldUser]),
		ViewFieldGroup:       len(ViewFieldNames[ViewFieldGroup]),
		ViewFieldModified:    len(dateS(time.Now())),
		ViewFieldAccessed:    len(dateS(time.Now())),
		ViewFieldCreated:     len(dateS(time.Now())),
		ViewFieldGit:         paw.MaxInt(3, len(ViewFieldNames[ViewFieldGit])),
		ViewFieldMd5:         32,
		ViewFieldName:        len(ViewFieldNames[ViewFieldName]),
	}

	ViewFieldColors = map[ViewField]*Color{
		ViewFieldNo:          paw.Cnop,
		ViewFieldINode:       paw.Cinp,
		ViewFieldPermissions: paw.Cpms,
		ViewFieldLinks:       paw.Clkp,
		ViewFieldSize:        paw.Csnp,
		ViewFieldBlocks:      paw.Cbkp,
		ViewFieldUser:        paw.Cuup,
		ViewFieldGroup:       paw.Cgup,
		ViewFieldModified:    paw.Cdap,
		ViewFieldCreated:     paw.Cdap,
		ViewFieldAccessed:    paw.Cdap,
		ViewFieldGit:         paw.Cgitp,
		ViewFieldMd5:         paw.Cmd5p,
		ViewFieldName:        paw.Cnop,
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
func (f ViewField) AlignS() string {
	if a, ok := ViewFieldAligns[f]; ok {
		return a.String()
	} else {
		return paw.AlignLeft.String()
	}
}
func (f ViewField) AlignSA() []string {
	fields := f.Fields()
	aligns := make([]string, 0, len(fields))
	for _, f := range fields {
		aligns = append(aligns, f.AlignS())
	}
	return aligns
}

func (f ViewField) SetColor(color *Color) {
	ViewFieldColors[f] = color
}

func (f ViewField) Color() *Color {
	if c, ok := ViewFieldColors[f]; ok {
		return c
	} else {
		return paw.Cdashp
	}
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

func (f ViewField) Count() int {
	return len(f.Fields())
}

// IsOk returns true for effective and otherwise not. In genernal, use it in checking.
func (f ViewField) IsOk() (ok bool) {
	paw.Logger.Debug("checking ViewField..." + paw.Caller(1))

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

func (v ViewField) RemoveGit(isNoGit bool) (vfields ViewField) {
	if !isNoGit || v&ViewFieldGit == 0 {
		return v
	}
	DoRangeViewField(v, func(i int, fd ViewField) {
		if fd&ViewFieldGit == 0 {
			vfields |= fd
		}
	})
	// for _, f := range v.Fields() {
	// 	if f&ViewFieldGit != 0 {
	// 		continue
	// 	}
	// 	vfields |= f
	// }
	return vfields
}

func (v ViewField) GetAllValues(de DirEntryX) (values []interface{}, cvalues []string, colors []*Color) {
	fields := v.Fields()
	values = make([]interface{}, 0, len(fields))
	cvalues = make([]string, 0, len(fields))
	colors = make([]*Color, 0, len(fields))
	DoRangeFields(fields, func(i int, fd ViewField) {
		values = append(values, de.Field(fd))
		cvalues = append(cvalues, de.FieldC(fd))
		if fd&ViewFieldName == 0 {
			colors = append(colors, fd.Color())
		} else {
			colors = append(colors, de.LSColor())
		}
	})
	return values, cvalues, colors
}

func (v ViewField) GetValues(de DirEntryX) (values []interface{}) {
	fields := v.Fields()
	values = make([]interface{}, 0, len(fields))
	DoRangeFields(fields, func(i int, fd ViewField) {
		values = append(values, de.Field(fd))
	})
	return values
}

func (v ViewField) GetValuesC(de DirEntryX) (values []string) {
	fields := v.Fields()
	values = make([]string, 0, len(fields))
	DoRangeFields(fields, func(i int, fd ViewField) {
		values = append(values, de.FieldC(fd))
	})
	return values
}

func (v ViewField) GetValuesS(de DirEntryX) (values []string) {
	fields := v.Fields()
	values = make([]string, 0, len(fields))
	for _, field := range fields {
		values = append(values, de.Field(field))
	}
	return values
}

func (v ViewField) GetHead(c *Color) string {
	var sprintf func(string, ...interface{}) string
	if c != nil {
		sprintf = c.Sprintf
	} else {
		sprintf = fmt.Sprintf
	}

	hd := ""
	DoRangeViewField(v, func(i int, fd ViewField) {
		if fd&ViewFieldName == 0 {
			value := fd.AlignedS(fd.Name())
			hd += sprintf("%v", value) + " "
		} else {
			value := paw.AlignWithWidth(fd.Align(), fd.Name(), fd.Width())
			hd += sprintf("%v", value)
		}
	})
	return hd
}

func (v ViewField) GetHeadA(c *Color) (values []string) {
	var sprint func(...interface{}) string
	if c != nil {
		sprint = c.Sprint
	} else {
		sprint = fmt.Sprint
	}
	fields := v.Fields()
	values = make([]string, 0, len(fields))
	DoRangeFields(fields, func(i int, fd ViewField) {
		v := sprint(fd.AlignedS(fd.Name()))
		values = append(values, v)
	})
	return values
}

func (v ViewField) GetHeadFunc(fc func(i int) *Color) (head string) {
	heads := v.GetHeadFuncA(fc)
	return strings.Join(heads, " ")
}

func (v ViewField) GetHeadFuncA(fc func(i int) *Color) (values []string) {
	if fc == nil {
		fc = paw.ChoseColorH
	}
	var sprint func(...interface{}) string
	fields := v.Fields()
	values = make([]string, 0, len(fields))
	if fc == nil {
		sprint = fmt.Sprint
		DoRangeFields(fields, func(i int, fd ViewField) {
			v := sprint(fd.AlignedS(fd.Name()))
			values = append(values, v)
		})
	} else {
		DoRangeFields(fields, func(i int, fd ViewField) {
			v := fc(i).Sprint(fd.AlignedS(fd.Name()))
			values = append(values, v)
		})
	}
	return values
}

// AlignedS return aligned string of value according to ViewField.Align()
func (v ViewField) AlignedS(value interface{}) string {
	return v.Align().ToString(value, v.Width())
}

// AlignedSC returns the aligned colorful string
func (v ViewField) AlignedSC(cvalue interface{}) string {
	return v.Align().ToStringC(cvalue, v.Width())
}

func (v ViewField) RowString(de DirEntryX) string {
	sb := new(strings.Builder)
	DoRangeViewField(v, func(i int, fd ViewField) {
		if fd&ViewFieldName == 0 {
			sb.WriteString(fd.AlignedS(de.Field(fd)) + " ")
		} else {
			sb.WriteString(fd.AlignedS(de.Field(fd)))
		}
	})
	// for _, field := range v.Fields() {
	// 	if field&ViewFieldName != 0 {
	// 		sb.WriteString(field.AlignedS(de.Field(field)))
	// 		continue
	// 	}
	// 	sb.WriteString(field.AlignedS(de.Field(field)) + " ")
	// }
	return sb.String()
}

func (v ViewField) Rows(de DirEntryX) (values []string) {
	fields := v.Fields()
	values = make([]string, 0, len(fields))
	DoRangeFields(fields, func(i int, fd ViewField) {
		values = append(values, fd.AlignedS(de.Field(fd)))
	})
	// for _, field := range fields {
	// 	values = append(values, field.AlignedS(de.Field(field)))
	// }
	return values
}
func (v ViewField) XattibutesRowsC(de DirEntryX) (rows [][]string) {
	xattrs := de.Xattibutes()
	if len(xattrs) == 0 {
		return nil
	}
	fields := v.Fields()
	nfd := len(fields)
	rows = make([][]string, 0, len(xattrs))
	cxs := make([]string, nfd)
	for i := 0; i < nfd-1; i++ {
		cxs[i] = paw.Spaces(fields[i].Width())

	}
	idx := nfd - 1
	xsymb := "@ "
	wdname := ViewFieldName.Width() - len(xsymb)
	for _, x := range xattrs {
		sp := paw.Spaces(wdname - paw.StringWidth(x))
		cxs[idx] = paw.Cxbp.Sprint(xsymb) + paw.Cxap.Sprint(x) + sp
		rows = append(rows, cxs)
	}
	return rows
}

func (v ViewField) XattibutesRowsSC(de DirEntryX) (rows []string) {
	xrows := v.XattibutesRowsC(de)
	if len(xrows) == 0 {
		return nil
	}
	rows = make([]string, 0, len(xrows))
	for _, rowa := range xrows {
		rows = append(rows, strings.Join(rowa, " "))
	}
	return rows
}

func (v ViewField) Rows2D(des []DirEntryX, isViewNoFiles, isViewNoDirs bool) (values [][]string) {
	if len(des) < 1 {
		return nil
	}
	hasNo := false
	if v&ViewFieldNo != 0 {
		hasNo = true
	}
	values = make([][]string, 0, len(des))
	var (
		nd, nf int
		idx    string
	)
	for _, de := range des {
		if hasNo {
			if de.IsDir() {
				if isViewNoFiles {
					continue
				}
				nd++
				idx = fmt.Sprintf("D%d", nd)
			} else {
				if isViewNoDirs {
					continue
				}
				nf++
				idx = fmt.Sprintf("F%d", nf)
			}
			ViewFieldNo.SetValue(idx)
		}
		values = append(values, v.Rows(de))
	}
	return values
}

func (v ViewField) RowStringXName(de DirEntryX) string {
	sb := new(strings.Builder)
	DoRangeViewField(v, func(i int, fd ViewField) {
		if fd&ViewFieldName == 0 {
			sb.WriteString(fd.AlignedS(de.Field(fd)) + " ")
		}
	})
	// for _, field := range v.Fields() {
	// 	if field&ViewFieldName != 0 {
	// 		continue
	// 	}
	// 	sb.WriteString(field.AlignedS(de.Field(field)) + " ")
	// }
	return sb.String()
}

func (v ViewField) RowStringC(de DirEntryX) string {
	sb := new(strings.Builder)
	DoRangeViewField(v, func(i int, fd ViewField) {
		if fd&ViewFieldName == 0 {
			sb.WriteString(de.FieldC(fd) + " ")
		} else {
			sb.WriteString(de.FieldC(fd))
		}
	})
	// for _, field := range v.Fields() {
	// 	if field&ViewFieldName != 0 {
	// 		sb.WriteString(de.FieldC(field))
	// 		continue
	// 	}
	// 	sb.WriteString(de.FieldC(field) + " ")
	// }
	return sb.String()
}
func (v ViewField) RowStringFC(de DirEntryX, fields []ViewField) string {
	var s string
	DoRangeFields(fields, func(i int, fd ViewField) {
		if fd&ViewFieldName == 0 {
			s += de.FieldC(fd) + " "
		} else {
			s += de.FieldC(ViewFieldName)
		}
	})
	return s
}

func (v ViewField) RowsC(de DirEntryX) (values []string) {
	fields := v.Fields()
	values = make([]string, 0, len(fields))
	DoRangeFields(fields, func(i int, fd ViewField) {
		values = append(values, de.FieldC(fd))
	})
	return values
}

func (v ViewField) RowStringXNameC(de DirEntryX) string {
	sb := new(strings.Builder)
	DoRangeViewField(v, func(i int, fd ViewField) {
		if fd&ViewFieldName == 0 {
			sb.WriteString(de.FieldC(fd) + " ")
		}
	})
	return sb.String()
}

func (v ViewField) GetModifyWidthsNoGitFields(d *Dir) []ViewField {
	fields := v.RemoveGit(true).Fields()
	modFieldWidths(d, fields)
	return fields
}

func (v ViewField) ModifyWidths(d *Dir) {
	fields := v.Fields()
	modFieldWidths(d, fields)
}

func modFieldWidths(d *Dir, fields []ViewField) {
	childWidths(d, fields)
	hasFieldNo := false
	DoRangeFields(fields, func(i int, fd ViewField) {
		if fd&ViewFieldNo != 0 {
			hasFieldNo = true
			return
		}
	})
	if hasFieldNo {
		nd, nf, _ := d.NItems(true)
		wdidx := GetMaxWidthOf(nd, nf)
		ViewFieldNo.SetWidth(wdidx + 1)
	}
	ViewFieldName.SetWidth(GetViewFieldNameWidthOf(fields))
}

func childWidths(d *Dir, fields []ViewField) {
	ds, _ := d.ReadDirAll()
	for _, de := range ds {
		DoRangeFields(fields, func(i int, fd ViewField) {
			wd := de.WidthOf(fd)
			if !de.IsDir() && fd&ViewFieldSize != 0 {
				if de.IsCharDev() || de.IsDev() {
					fmajor := ViewFieldMajor.Width()
					fminor := ViewFieldMinor.Width()
					major, minor := de.DevNumber()
					wdmajor := len(cast.ToString(major))
					wdminor := len(cast.ToString(minor))
					ViewFieldMajor.SetWidth(paw.MaxInt(fmajor, wdmajor))
					ViewFieldMinor.SetWidth(paw.MaxInt(fminor, wdminor))
					wd = ViewFieldMajor.Width() +
						ViewFieldMinor.Width() + 1
				}
			}
			width := paw.MaxInt(fd.Width(), wd)
			fd.SetWidth(width)
		})
		if de.IsDir() {
			child := de.(*Dir)
			childWidths(child, fields)
		}
	}
}

func DoRangeFields(fields []ViewField, fc func(i int, fd ViewField)) {
	for i := 0; i < len(fields); i++ {
		fc(i, fields[i])
	}
}
func DoRangeViewField(vfields ViewField, fc func(i int, fd ViewField)) {
	fields := vfields.Fields()
	for i := 0; i < len(fields); i++ {
		fc(i, fields[i])
	}
}

// func GetPFHeadS(c *Color, fields ...ViewField) string {
// 	var sprintf func(string, ...interface{}) string
// 	if c != nil {
// 		sprintf = c.Sprintf
// 	} else {
// 		sprintf = fmt.Sprintf
// 	}

// 	hd := ""
// 	for _, f := range fields {
// 		if f&ViewFieldName != 0 {
// 			hd += sprintf("%-[1]*[2]s", f.Width(), f.Name())
// 			continue
// 		}
// 		value := f.AlignedS(f.Name())
// 		hd += sprintf("%v", value) + " "
// 	}
// 	return hd
// }
