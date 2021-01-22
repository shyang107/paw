package paw

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/shyang107/paw/cast"
)

// TableFormat define the format used to print out
//
// 	Elements:
// `Fields []string{}`, string list of heading row
// `LenFields []int{}`, length of every field of heading row
// `Aligns []Align{}` : ,
// `Padding string`: paddign string befor every row, default "",
// `Sep string` : sepperating string between fields
// `TopChar string`: character of top sepperating row, default "=",
// `MiddleChar string`: character of middle sepperating row , default "-",
// `BottomChar string`: character of bottom sepperating row, default "=",
type TableFormat struct {
	Fields            []string
	LenFields         []int
	Aligns            []Align
	Colors            []*color.Color
	FieldsColorString []string
	Padding           string
	Sep               string
	TopChar           string
	topBanner         string
	MiddleChar        string
	midBanner         string
	BottomChar        string
	botBanner         string
	beforeMsg         string
	afterMsg          string
	writer            io.Writer
	isAbbrSymbol      bool
	isPrepare         bool
	isPrepareBefore   bool
	isPrepareAfter    bool
	IsWrapped         bool
	IsColorful        bool
	// chdEven         *color.Color
	// chdOdd          *color.Color
	XAttributeSymbol  string
	XAttributeSymbol2 string
}

// Align is id that indicate alignment of head-column
type Align int

const (
	space      = " "
	abbrSymbol = "»"
	// AlignLeft align left
	AlignLeft Align = iota
	// AlignCenter align center
	AlignCenter
	// AlignRight align right
	AlignRight
)

// NewTableFormat return a instance of TableFormat
func NewTableFormat() *TableFormat {
	return &TableFormat{
		Fields:          []string{},
		LenFields:       []int{},
		Aligns:          []Align{},
		Padding:         "",
		Sep:             space,
		TopChar:         "=",
		MiddleChar:      "-",
		BottomChar:      "=",
		writer:          os.Stdout,
		isAbbrSymbol:    false,
		isPrepare:       false,
		isPrepareBefore: false,
		isPrepareAfter:  false,
		IsWrapped:       false,
		IsColorful:      false,
		// chdEven:         tbChdEven,
		// chdOdd:          tbChdOdd,
		XAttributeSymbol:  tbxSymb,
		XAttributeSymbol2: tbxSymb2,
	}
}

var (
	// isAbbrSymbol    bool
	// isPrepare       bool
	// isPrepareBefore bool
	// isPrepareAfter  bool
	tbChdEven  = color.New([]color.Attribute{38, 5, 228, 1, 48, 5, 236}...)
	tbChdOdd   = color.New([]color.Attribute{38, 5, 156, 1, 48, 5, 234}...)
	tbCRowEven = color.New([]color.Attribute{38, 5, 253, 48, 5, 236}...)
	tbCRowOdd  = color.New([]color.Attribute{38, 5, 156, 48, 5, 234}...)
	// tbCxattrEven = color.New([]color.Attribute{38, 5, 249, 4, 48, 5, 236}...)
	// tbCxattrOdd  = color.New([]color.Attribute{38, 5, 249, 4, 48, 5, 234}...)
	// tbCxsymbEven = color.New([]color.Attribute{38, 5, 249, 48, 5, 236}...)
	// tbCxsymbOdd  = color.New([]color.Attribute{38, 5, 249, 48, 5, 234}...)
	tbCxattr = color.New([]color.Attribute{38, 5, 249, 4, 48, 5, 234}...)
	tbCxsymb = color.New([]color.Attribute{38, 5, 249, 48, 5, 234}...)
	tbxSymb  = " @ "
	tbxSymb2 = "-@-"
	tbxSp    = "   "
)

// func (t *TableFormat) setColor() {
// 	t.chdEven = tbChdEven
// 	t.chdOdd = tbChdOdd
// }

// NFields will return number of TableFormat.Fields
func (t *TableFormat) NFields() int {
	return len(t.Fields)
}

// SetBeforeMessage set message to show before table
func (t *TableFormat) SetBeforeMessage(msg string) {
	t.beforeMsg = msg
	if len(t.Padding) > 0 {
		// t.beforeMsg = t.Padding + t.beforeMsg
		t.beforeMsg = PaddingString(t.beforeMsg, t.Padding)
	}
	t.isPrepareBefore = true
}

// SetAfterMessage set message to show after table
func (t *TableFormat) SetAfterMessage(msg string) {
	t.afterMsg = msg
	if len(t.Padding) > 0 {
		// t.afterMsg = t.Padding + t.afterMsg
		t.afterMsg = PaddingString(t.afterMsg, t.Padding)
	}
	t.isPrepareAfter = true
}

// Prepare initialize `TableFormat`
func (t *TableFormat) Prepare(w io.Writer) {
	t.check()
	t.writer = w
	t.isPrepare = false
	// if t.IsColorful {
	// 	t.setColor()
	// }
}

// SetWrapFields set true to TableFormat.IsWrapped
func (t *TableFormat) SetWrapFields() {
	t.IsWrapped = true
}

func (t *TableFormat) check() {
	if len(t.Sep) == 0 {
		t.Sep = space
	}
	t.checkFields()
	t.checkAlign()
	t.setBanner()
}

func (t *TableFormat) checkAlign() {
	la := len(t.Aligns)
	lf := len(t.Fields)
	switch {
	case la < 1:
		t.Aligns = make([]Align, lf)
		for i := 0; i < len(t.Fields); i++ {
			t.Aligns[i] = AlignRight
		}
	case la > lf:
		t.Aligns = t.Aligns[:lf]
	case la < lf && la > 1:
		tmp := make([]Align, lf)
		copy(tmp, t.Aligns)
		for i := la; i < lf; i++ {
			tmp[i] = AlignRight
		}
		t.Aligns = tmp
	}
}

func (t *TableFormat) checkBannerChar() {
	if len(t.TopChar) != 1 {
		t.TopChar = "="
	}
	if len(t.MiddleChar) != 1 {
		t.MiddleChar = "-"
	}
	if len(t.BottomChar) != 1 {
		t.BottomChar = "="
	}
	if len(t.Sep) == 0 {
		t.Sep = space
	}
}

func (t *TableFormat) setBanner() {
	t.checkBannerChar()
	llf := len(t.LenFields)
	lsep := len(t.Sep)
	tlen := 0
	for i := 0; i < llf-1; i++ {
		tlen += t.LenFields[i]
		tlen += lsep
	}
	tlen += t.LenFields[llf-1]
	t.topBanner = Repeat(t.TopChar, tlen)
	t.botBanner = Repeat(t.BottomChar, tlen)
	sb.Reset()
	for i := 0; i < llf-1; i++ {
		sb.WriteString(Repeat(t.MiddleChar, t.LenFields[i]))
		sb.WriteString(t.Sep)
	}
	sb.WriteString(Repeat(t.MiddleChar, t.LenFields[llf-1]))
	t.midBanner = sb.String()
	sb.Reset()
	if len(t.Padding) > 0 {
		t.topBanner = t.Padding + t.topBanner
		t.midBanner = t.Padding + t.midBanner
		t.botBanner = t.Padding + t.botBanner
	}
}

func (t *TableFormat) checkFields() {
	lf := len(t.Fields)
	llf := len(t.LenFields)
	switch {
	case llf < lf:
		for i := llf; i < lf; i++ {
			hc, ac := CountPlaceHolder(t.Fields[i])
			t.LenFields = append(t.LenFields[:llf], hc+ac)
		}
	case llf > lf:
		t.LenFields = t.LenFields[:lf]
	}
	for i := 0; i < lf; i++ {
		if t.LenFields[i] < len(t.Fields[i]) {
			t.LenFields[i] = len(t.Fields[i])
		}
	}
}

func (t *TableFormat) getRowString(fields []string, widths []int, aligns []Align, sep string, padding string) string {
	sb.Reset()
	var (
		str     string
		nFields = len(fields)
	)
	if t.IsWrapped {
		goto WRAPFIELDS
	}
	for i := 0; i < nFields; i++ {
		v := fields[i]
		wd := widths[i]
		// v = GetAbbrString(v, wd, "»")
		v = Truncate(v, wd, "»")
		// v = Wrap(v, wd)
		// nh, na := CountPlaceHolder(v)
		al := aligns[i]
		s := t.getAlignString(i, al, wd, v)
		sb.WriteString(s + sep)
	}
	str = sb.String()
	if !t.isAbbrSymbol {
		t.isAbbrSymbol = Contains(str, abbrSymbol)
	}
	return padding + str

WRAPFIELDS:
	// wfields store [nfields][nlines]
	wfields := make([][]string, nFields)
	// niles number of lines of every wrapped field
	nlines := make([]int, nFields)
	// idx record index postion
	idx := make([]int, nFields)
	for i, v := range fields {
		wd := widths[i]
		wfields[i] = WrapToSlice(v, wd) //Split(Wrap(v, wd), "\n")
		nlines[i] = len(wfields[i])     // count number of lines: wfields[i]
		idx[i] = 0
	}
	maxlines := MaxInts(nlines...)
	for i := 0; i < maxlines; i++ { // ith line
		hasXattr := false
		for j, wrapfields := range wfields {
			if idx[j] < nlines[j] {
				v := wrapfields[idx[j]]
				if HasPrefix(v, t.XAttributeSymbol) ||
					HasPrefix(v, t.XAttributeSymbol2) {
					hasXattr = true
					break
				}
			}
		}
		for j, wrapfields := range wfields { // jth field
			v := ""
			if idx[j] < nlines[j] {
				v = wrapfields[idx[j]]
			}
			s := ""
			if i == 0 && !hasXattr &&
				t.FieldsColorString != nil &&
				len(t.FieldsColorString[j]) > 0 {
				s = t.FieldsColorString[j]
			} else {
				s = t.getAlignString(j, aligns[j], widths[j], v)
			}
			sb.WriteString(s + sep)
			idx[j]++
		}
		if i < maxlines-1 {
			sb.WriteByte('\n')
		}
	}
	str = sb.String()
	return PaddingString(str, padding)
}

func (t *TableFormat) getHeadColorString(col int, field string) string {
	var c *color.Color
	switch col % 2 {
	case 0:
		c = tbChdEven
	case 1:
		c = tbChdOdd
	}
	return c.Sprint(field)
}

func getColorField(value string, cf, ct *color.Color, align Align, width int) string {
	// wf := StringWidth(value)
	s := TrimSpace(value)
	ws := StringWidth(s)
	// fmt.Println("width =", width, "wf =", wf)
	switch align {
	case AlignRight:
		s = ct.Sprint(Spaces(width-ws)) + cf.Sprint(s)
		// s = ct.Sprint(Repeat("X", width-ws)) + cf.Sprint(s)
	case AlignCenter:
		wsl := (width - ws) / 2
		wsr := width - ws - wsl
		s = ct.Sprint(Spaces(wsl)) + cf.Sprint(s) + ct.Sprint(Spaces(wsr))
	default: //AlignLeft
		s = cf.Sprint(s) + ct.Sprint(Spaces(width-ws))
		// s = cf.Sprint(s) + ct.Sprint(Repeat("X", width-ws))
	}

	return s
}

func getColorxattr(t *TableFormat, value, xsymb string, cs, cx, r *color.Color, width int) string {
	xattr := TrimRight(TrimPrefix(value, xsymb), space)
	wd := StringWidth(xattr) + StringWidth(xsymb)
	tail := ""
	if width-wd > 0 {
		tail = Spaces(width - wd)
	}
	if xsymb == t.XAttributeSymbol2 {
		xsymb = tbxSp
	}
	return cs.Sprint(xsymb) + cx.Sprint(xattr) + r.Sprint(tail)
}

func (t *TableFormat) getAlignString(col int, al Align, width int, value string) string {
	if t.IsColorful {
		var (
			r *color.Color
			c *color.Color
		)

		switch col % 2 {
		case 0:
			r = tbCRowEven
		case 1:
			r = tbCRowOdd
		}
		if t.Colors != nil {
			c = t.Colors[col]
		} else {
			c = r
		}
		if HasPrefix(value, t.XAttributeSymbol) {
			return getColorxattr(t, value, t.XAttributeSymbol, tbCxsymb, tbCxattr, r, width)
		} else if HasPrefix(value, t.XAttributeSymbol2) {
			return getColorxattr(t, value, t.XAttributeSymbol2, tbCxsymb, tbCxattr, r, width)
		}
		// value = fmt.Sprintf("%[1]*[2]s", width, value)
		s := getColorField(value, c, r, al, width)
		return s
	} else {
		var s string
		switch al {
		case AlignLeft:
			s = FillRight(value, width)
		case AlignRight:
			s = FillLeft(value, width)
		case AlignCenter:
			// lv := StringWidth(value)
			// nr := (width - lv) / 2
			// nl := width - lv - nr
			// // s = Repeat(space, nl) + value + Repeat(space, nr)
			// s = Spaces(nl) + value + Spaces(nr)
			s = FillLeftRight(value, width)
		default:
			s = value
		}
		return s
	}
	return ""
}

// PrintSart print out head-section in `t.Writer`
func (t *TableFormat) PrintSart() error {
	if !t.isPrepare {
		t.Prepare(os.Stdout)
	}
	if t.isPrepareBefore {
		fmt.Fprintln(t.writer, t.beforeMsg)
	}
	fmt.Fprintln(t.writer, t.topBanner)

	fmt.Fprintln(t.writer, t.getRowString(t.Fields, t.LenFields, t.Aligns, t.Sep, t.Padding))

	fmt.Fprintln(t.writer, t.midBanner)
	return nil
}

// PrintRow print row into `t.writer`
func (t *TableFormat) PrintRow(rows ...interface{}) {
	sRows := make([]string, len(rows))
	for i, v := range rows {
		sRows[i] = cast.ToString(v) //fmt.Sprintf("%v", v)
	}

	fmt.Fprintln(t.writer, t.getRowString(sRows, t.LenFields, t.Aligns, t.Sep, t.Padding))
}

// PrintLine prints s without field speration into `t.writer` in default format
func (t *TableFormat) PrintLine(s interface{}) {
	ss := t.Padding + fmt.Sprint(s)
	fmt.Fprintf(t.writer, "%s", ss)
}

// PrintLineln print s without field speration into `t.writer` in default format, end with '\n'
func (t *TableFormat) PrintLineln(s interface{}) {
	ss := t.Padding + fmt.Sprint(s)
	fmt.Fprintf(t.writer, "%s\n", ss)
}

// PrintLinef print s with format and no field speration into `t.writer`
func (t *TableFormat) PrintLinef(format string, s ...interface{}) {
	ss := t.Padding + fmt.Sprintf(format, s...)
	fmt.Fprintf(t.writer, "%s", ss)
}

// PrintMiddleSepLine print middle sepperating line using `MiddleChar`
func (t *TableFormat) PrintMiddleSepLine() {
	fmt.Fprintln(t.writer, t.midBanner)
}

// PrintEnd print end-section into `t.writer`
func (t *TableFormat) PrintEnd() {
	if t.isAbbrSymbol {
		fmt.Fprintln(t.writer, ReplaceAll(t.botBanner, t.BottomChar, t.MiddleChar))
		fmt.Fprintln(t.writer, t.Padding+"* '"+abbrSymbol+"' : abbreviated symbol of a term")
	}
	fmt.Fprintln(t.writer, t.botBanner)
	if t.isPrepareAfter {
		fmt.Fprintln(t.writer, t.afterMsg)
	}
}
